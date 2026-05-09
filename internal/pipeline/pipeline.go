// Package pipeline wires together the logslice processing stages into a
// single reusable unit: scan → multiline-collapse → filter → sample →
// truncate → highlight → output.
//
// Callers construct a Pipeline with New, optionally configure it via the
// Options struct, and then call Run which reads from the supplied io.Reader
// and writes matching lines to the supplied io.Writer.
package pipeline

import (
	"context"
	"io"
	"regexp"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/highlight"
	"github.com/user/logslice/internal/multiline"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/sampler"
	"github.com/user/logslice/internal/scanner"
	"github.com/user/logslice/internal/truncate"
)

// Options controls every processing stage of the pipeline.
type Options struct {
	// Filter options.
	Filter filter.Options

	// MultilinePattern, when non-nil, groups continuation lines into a single
	// logical record before filtering.
	MultilinePattern *regexp.Regexp
	// MultilineMaxLines caps how many raw lines one logical record may span.
	MultilineMaxLines int

	// Sampler controls line-level down-sampling (e.g. every Nth line).
	Sampler sampler.Options

	// Truncate controls per-line length limiting.
	Truncate truncate.Options

	// HighlightPattern, when non-nil, wraps matching text in ANSI colour codes.
	HighlightPattern *regexp.Regexp
	// HighlightColor is the ANSI SGR code used for highlighting (default: "33").
	HighlightColor string

	// Scanner controls how the raw byte stream is tokenised.
	Scanner scanner.Options
}

// Result summarises what happened during a Run.
type Result struct {
	LinesRead    int64
	LinesMatched int64
	BytesWritten int64
}

// Pipeline is a configured processing chain ready to Run.
type Pipeline struct {
	opts Options
}

// New creates a Pipeline with the provided options.
func New(opts Options) *Pipeline {
	return &Pipeline{opts: opts}
}

// Run reads log lines from r, applies every configured stage in order, and
// writes matching output to w.  It respects ctx cancellation and returns any
// first non-EOF error encountered.
func (p *Pipeline) Run(ctx context.Context, r io.Reader, w io.Writer) (Result, error) {
	var res Result

	// --- output writer ---
	out := output.New(w)

	// --- highlight ---
	var hl *highlight.Highlighter
	if p.opts.HighlightPattern != nil {
		color := p.opts.HighlightColor
		if color == "" {
			color = "33"
		}
		hl = highlight.New(p.opts.HighlightPattern, color)
	}

	// --- truncator ---
	tr := truncate.New(p.opts.Truncate)

	// --- sampler ---
	smp, err := sampler.New(p.opts.Sampler)
	if err != nil {
		return res, err
	}

	// --- filter ---
	f := filter.New(p.opts.Filter)

	// --- multiline collector ---
	var ml *multiline.Collector
	if p.opts.MultilinePattern != nil {
		ml, err = multiline.New(multiline.Options{
			Pattern:  p.opts.MultilinePattern,
			MaxLines: p.opts.MultilineMaxLines,
		})
		if err != nil {
			return res, err
		}
	}

	// --- scanner ---
	sc := scanner.New(r, p.opts.Scanner)

	// processLine runs a fully assembled logical line through the later stages.
	processLine := func(line string) error {
		res.LinesRead++
		if !f.Match(line) {
			return nil
		}
		if !smp.Accept() {
			return nil
		}
		line = tr.Truncate(line)
		if hl != nil {
			line = hl.Apply(line)
		}
		if err := out.WriteLine(line); err != nil {
			return err
		}
		res.LinesMatched++
		return nil
	}

	for {
		if ctx.Err() != nil {
			break
		}
		raw, scanErr := sc.Next()
		if scanErr == io.EOF {
			// Flush any buffered multiline record.
			if ml != nil {
				if rec := ml.Flush(); rec != "" {
					if err := processLine(rec); err != nil {
						return res, err
					}
				}
			}
			break
		}
		if scanErr != nil {
			return res, scanErr
		}

		if ml != nil {
			if rec, complete := ml.Feed(raw); complete {
				if err := processLine(rec); err != nil {
					return res, err
				}
			}
			continue
		}

		if err := processLine(raw); err != nil {
			return res, err
		}
	}

	res.BytesWritten = out.BytesWritten()
	return res, nil
}
