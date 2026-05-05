// Package tail implements efficient tail-reading of large log files.
//
// It provides backward-seeking reads to extract the last N lines without
// loading the entire file into memory, making it suitable for very large
// log files.
//
// Basic usage:
//
//	r := tail.New(tail.Options{NumLines: 50})
//	lines, err := r.ReadFile("/var/log/app.log")
//
// Context-aware filtering:
//
//	f, _ := filter.New(filter.Options{Regex: "ERROR"})
//	cr := tail.NewContextReader(r, f, 2, 2)
//	results, err := cr.Read(file)
//	for _, res := range results {
//		fmt.Print(res.Format())
//	}
package tail
