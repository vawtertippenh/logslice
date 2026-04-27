// logslice is a fast log file slicer and filter tool that supports
// time-range queries and regex filtering for large log files.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/yourusername/logslice/internal/filter"
	"github.com/yourusername/logslice/internal/logline"
	"github.com/yourusername/logslice/internal/timeparse"
)

const usage = `logslice — Fast log file slicer and filter tool

Usage:
  logslice [options] <logfile>

Options:
  -from string
        Start of time range (e.g. "2024-01-15 10:00:00")
  -to string
        End of time range (e.g. "2024-01-15 11:00:00")
  -match string
        Regex pattern to filter log lines
  -tz string
        Timezone for time parsing (default: UTC)
  -h    Show this help message

Examples:
  logslice -from "2024-01-15 10:00:00" -to "2024-01-15 11:00:00" app.log
  logslice -match "ERROR|WARN" app.log
  logslice -from "2024-01-15 10:00:00" -match "ERROR" app.log
`

func main() {
	var (
		fromStr = flag.String("from", "", "Start of time range")
		toStr   = flag.String("to", "", "End of time range")
		match   = flag.String("match", "", "Regex pattern to filter log lines")
		tz      = flag.String("tz", "UTC", "Timezone for time parsing")
		help    = flag.Bool("h", false, "Show help")
	)

	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	if *help {
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: no log file specified")
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	loc, err := time.LoadLocation(*tz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid timezone %q: %v\n", *tz, err)
		os.Exit(1)
	}

	var from, to time.Time
	if *fromStr != "" {
		from, err = timeparse.ParseWithLocation(*fromStr, loc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot parse -from time %q: %v\n", *fromStr, err)
			os.Exit(1)
		}
	}
	if *toStr != "" {
		to, err = timeparse.ParseWithLocation(*toStr, loc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot parse -to time %q: %v\n", *toStr, err)
			os.Exit(1)
		}
	}

	var re *regexp.Regexp
	if *match != "" {
		re, err = regexp.Compile(*match)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid regex %q: %v\n", *match, err)
			os.Exit(1)
		}
	}

	f := filter.New(from, to, re)

	filePath := flag.Arg(0)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot open file %q: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	parser := logline.NewParser(loc)
	scanner := bufio.NewScanner(file)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for scanner.Scan() {
		raw := scanner.Bytes()
		line := parser.Parse(raw)
		if f.Match(line) {
			out.Write(raw)
			out.WriteByte('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: reading file: %v\n", err)
		os.Exit(1)
	}
}
