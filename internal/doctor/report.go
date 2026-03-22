package doctor

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Summary holds aggregate scan counts.
type Summary struct {
	FilesScanned        int `json:"files_scanned"          yaml:"files_scanned"          xml:"FilesScanned"`
	FilesWithFindings   int `json:"files_with_findings"    yaml:"files_with_findings"    xml:"FilesWithFindings"`
	EnvVarsScanned      int `json:"env_vars_scanned"       yaml:"env_vars_scanned"       xml:"EnvVarsScanned"`
	EnvVarsWithFindings int `json:"env_vars_with_findings" yaml:"env_vars_with_findings" xml:"EnvVarsWithFindings"`
	TotalFindings       int `json:"total_findings"         yaml:"total_findings"         xml:"TotalFindings"`
}

// Report is the top-level structure for all output formats.
type Report struct {
	GeneratedAt string        `json:"generated_at" yaml:"generated_at" xml:"GeneratedAt"`
	Summary     Summary       `json:"summary"      yaml:"summary"      xml:"Summary"`
	FilesRead   []string      `json:"files_read"   yaml:"files_read"   xml:"FilesRead>Path"`
	FileScan    []FileFinding `json:"file_scan"    yaml:"file_scan"    xml:"FileScan>Finding"`
	EnvScan     []EnvFinding  `json:"env_scan"     yaml:"env_scan"     xml:"EnvScan>Finding"`
}

// NewReport constructs a Report from scan results.
func NewReport(
	fileFindings []FileFinding, filesRead []string, filesWithFindings int,
	envFindings []EnvFinding, varsScanned, varsWithFindings int,
) Report {
	return Report{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Summary: Summary{
			FilesScanned:        len(filesRead),
			FilesWithFindings:   filesWithFindings,
			EnvVarsScanned:      varsScanned,
			EnvVarsWithFindings: varsWithFindings,
			TotalFindings:       len(fileFindings) + len(envFindings),
		},
		FilesRead: filesRead,
		FileScan:  fileFindings,
		EnvScan:   envFindings,
	}
}

// WriteReport serializes r in the given format to w.
// Supported formats: json, yaml, xml, csv, html.
func WriteReport(w io.Writer, r Report, format string) error {
	switch format {
	case "json":
		return writeJSON(w, r)
	case "yaml":
		return writeYAML(w, r)
	case "xml":
		return writeXML(w, r)
	case "csv":
		return writeCSV(w, r)
	case "html":
		return writeHTML(w, r)
	default:
		return fmt.Errorf("unsupported output format %q: must be one of json, yaml, xml, csv, html", format)
	}
}

func writeJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func writeYAML(w io.Writer, r Report) error {
	return yaml.NewEncoder(w).Encode(r)
}

func writeXML(w io.Writer, r Report) error {
	if _, err := fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(struct {
		XMLName xml.Name `xml:"DoctorReport"`
		Report
	}{Report: r}); err != nil {
		return err
	}
	return enc.Flush()
}

func writeCSV(w io.Writer, r Report) error {
	cw := csv.NewWriter(w)
	// Files-read manifest
	if err := cw.Write([]string{"section", "path"}); err != nil {
		return err
	}
	for _, path := range r.FilesRead {
		if err := cw.Write([]string{"files_read", path}); err != nil {
			return err
		}
	}
	// Blank separator row then findings header
	if err := cw.Write([]string{}); err != nil {
		return err
	}
	if err := cw.Write([]string{"scan_type", "file", "line", "variable", "pattern", "value_redacted"}); err != nil {
		return err
	}
	for _, f := range r.FileScan {
		if err := cw.Write([]string{
			"file", f.File, strconv.Itoa(f.Line), f.Variable, f.Pattern, f.ValueRedacted,
		}); err != nil {
			return err
		}
	}
	for _, e := range r.EnvScan {
		if err := cw.Write([]string{
			"env", "", "", e.Variable, e.Pattern, e.ValueRedacted,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

var htmlTmpl = template.Must(template.New("report").Funcs(template.FuncMap{
	"inc": func(i int) int { return i + 1 },
}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>TokenShim Doctor Report</title>
<style>
  body { font-family: monospace; max-width: 960px; margin: 2rem auto; color: #222; }
  h1 { color: #c0392b; }
  h2 { border-bottom: 1px solid #ddd; padding-bottom: .3rem; }
  table { border-collapse: collapse; width: 100%; margin-bottom: 2rem; }
  th { background: #f5f5f5; text-align: left; padding: .4rem .6rem; border: 1px solid #ddd; }
  td { padding: .4rem .6rem; border: 1px solid #ddd; word-break: break-all; }
  .badge { display:inline-block; padding:.2rem .5rem; border-radius:3px; font-size:.85em; }
  .ok   { background:#d4edda; color:#155724; }
  .warn { background:#fff3cd; color:#856404; }
  .meta { color: #888; font-size: .9em; margin-bottom: 1rem; }
  .empty { color: #888; font-style: italic; }
</style>
</head>
<body>
<h1>TokenShim Doctor Report</h1>
<p class="meta">Generated: {{.GeneratedAt}}</p>

<h2>Summary</h2>
<table>
  <tr><th>Metric</th><th>Value</th></tr>
  <tr><td>Total Findings</td><td>{{if gt .Summary.TotalFindings 0}}<span class="badge warn">{{.Summary.TotalFindings}}</span>{{else}}<span class="badge ok">0</span>{{end}}</td></tr>
  <tr><td>Files Scanned</td><td>{{.Summary.FilesScanned}}</td></tr>
  <tr><td>Files with Findings</td><td>{{.Summary.FilesWithFindings}}</td></tr>
  <tr><td>Env Vars Scanned</td><td>{{.Summary.EnvVarsScanned}}</td></tr>
  <tr><td>Env Vars with Findings</td><td>{{.Summary.EnvVarsWithFindings}}</td></tr>
</table>

<h2>Files Read</h2>
{{if .FilesRead}}
<table>
  <tr><th>#</th><th>Path</th></tr>
  {{range $i, $p := .FilesRead}}
  <tr><td>{{inc $i}}</td><td>{{$p}}</td></tr>
  {{end}}
</table>
{{else}}<p class="empty">No files scanned.</p>{{end}}

<h2>File Scan Findings</h2>
{{if .FileScan}}
<table>
  <tr><th>File</th><th>Line</th><th>Variable</th><th>Pattern</th><th>Value (redacted)</th></tr>
  {{range .FileScan}}
  <tr><td>{{.File}}</td><td>{{.Line}}</td><td>{{.Variable}}</td><td>{{.Pattern}}</td><td>{{.ValueRedacted}}</td></tr>
  {{end}}
</table>
{{else}}<p class="empty">No findings.</p>{{end}}

<h2>Environment Variable Scan Findings</h2>
{{if .EnvScan}}
<table>
  <tr><th>Variable</th><th>Pattern</th><th>Value (redacted)</th></tr>
  {{range .EnvScan}}
  <tr><td>{{.Variable}}</td><td>{{.Pattern}}</td><td>{{.ValueRedacted}}</td></tr>
  {{end}}
</table>
{{else}}<p class="empty">No findings.</p>{{end}}

</body>
</html>
`))

func writeHTML(w io.Writer, r Report) error {
	return htmlTmpl.Execute(w, r)
}
