# logslice

Fast log file slicer and filter tool with time-range queries and regex support for large files.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build -o logslice .
```

---

## Usage

```
logslice [flags] <logfile>

Flags:
  --from      Start of time range (e.g. "2024-01-15T08:00:00")
  --to        End of time range   (e.g. "2024-01-15T09:00:00")
  --pattern   Filter lines by regex pattern
  --format    Timestamp format (default: RFC3339)
  --out       Output file (default: stdout)
```

### Examples

Slice logs between two timestamps:
```bash
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log
```

Filter by pattern within a time range:
```bash
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" --pattern "ERROR|WARN" app.log
```

Write results to a file:
```bash
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" --out slice.log app.log
```

---

## Features

- Binary search for fast time-range lookups in large files
- Regex filtering with full RE2 syntax support
- Supports custom timestamp formats
- Streams output — minimal memory footprint

---

## License

MIT © 2024 yourname