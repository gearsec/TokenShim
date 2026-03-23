package doctor

import (
	"bufio"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// FileFinding represents a single secret detected in a file.
type FileFinding struct {
	File          string `json:"file"           yaml:"file"           xml:"File"`
	Line          int    `json:"line"           yaml:"line"           xml:"Line"`
	Pattern       string `json:"pattern"        yaml:"pattern"        xml:"Pattern"`
	Variable      string `json:"variable"       yaml:"variable"       xml:"Variable"`
	ValueRedacted string `json:"value_redacted" yaml:"value_redacted" xml:"ValueRedacted"`
}

// EnvFinding represents a single secret detected in an environment variable.
type EnvFinding struct {
	Variable      string `json:"variable"       yaml:"variable"      xml:"Variable"`
	Pattern       string `json:"pattern"        yaml:"pattern"       xml:"Pattern"`
	ValueRedacted string `json:"value_redacted" yaml:"value_redacted" xml:"ValueRedacted"`
}

// ScanFiles resolves all paths from cfg, reads each file, and returns findings.
// filesRead contains the paths that were successfully opened and scanned.
func ScanFiles(cfg Config) (findings []FileFinding, filesRead []string, filesWithFindings int, err error) {
	paths, err := resolvePaths(cfg)
	if err != nil {
		return nil, nil, 0, err
	}

	for _, path := range paths {
		fileFindings, scanErr := scanFile(path)
		if scanErr != nil {
			// Skip unreadable files silently — they may be owned by root, etc.
			continue
		}
		filesRead = append(filesRead, path)
		if len(fileFindings) > 0 {
			filesWithFindings++
			findings = append(findings, fileFindings...)
		}
	}
	return findings, filesRead, filesWithFindings, nil
}

// ScanEnv reads all process environment variables and returns findings.
func ScanEnv() (findings []EnvFinding, varsScanned int, varsWithFindings int) {
	env := os.Environ()
	varsScanned = len(env)

	seen := make(map[string]bool)
	for _, entry := range env {
		idx := strings.IndexByte(entry, '=')
		if idx < 0 {
			continue
		}
		key := entry[:idx]
		value := entry[idx+1:]

		if patternName, matched := matchPatterns(key, value); matched {
			if !seen[key] {
				seen[key] = true
				varsWithFindings++
			}
			findings = append(findings, EnvFinding{
				Variable:      key,
				Pattern:       patternName,
				ValueRedacted: redactValue(value),
			})
		}
	}
	return findings, varsScanned, varsWithFindings
}

// resolvePaths expands ~ and globs (including **) returning deduplicated existing file paths.
func resolvePaths(cfg Config) ([]string, error) {
	seen := make(map[string]bool)
	var result []string

	for _, raw := range cfg.ScanPaths {
		expanded, err := ExpandPath(raw)
		if err != nil {
			continue
		}

		matches, err := expandGlob(expanded)
		if err != nil || len(matches) == 0 {
			// Treat as a literal path with no wildcards.
			if info, statErr := os.Stat(expanded); statErr == nil && info.Mode().IsRegular() {
				if !seen[expanded] {
					seen[expanded] = true
					result = append(result, expanded)
				}
			}
			continue
		}

		for _, match := range matches {
			if info, statErr := os.Stat(match); statErr == nil && info.Mode().IsRegular() {
				if !seen[match] {
					seen[match] = true
					result = append(result, match)
				}
			}
		}
	}

	return result, nil
}

// expandGlob resolves a glob pattern, adding support for ** (recursive wildcard).
// For patterns without **, it delegates to filepath.Glob.
// For patterns with **, it walks the filesystem from the base directory and
// matches file names against the suffix pattern (e.g. "**/.env*" walks "."
// and matches any file whose name matches ".env*").
func expandGlob(pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return filepath.Glob(pattern)
	}

	// Split on the first **: everything before is the base dir to walk,
	// everything after is the filename pattern to match.
	idx := strings.Index(pattern, "**")
	baseDir := filepath.Clean(pattern[:idx])
	if baseDir == "" {
		baseDir = "."
	}
	suffix := strings.TrimPrefix(pattern[idx+2:], string(filepath.Separator))

	var matches []string
	if err := fs.WalkDir(os.DirFS(baseDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if suffix == "" {
			matches = append(matches, filepath.Join(baseDir, path))
			return nil
		}
		if ok, _ := filepath.Match(suffix, name); ok {
			matches = append(matches, filepath.Join(baseDir, path))
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return matches, nil
}

// scanFile reads a file line-by-line and checks each line for secret patterns.
func scanFile(path string) ([]FileFinding, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var findings []FileFinding
	scanner := bufio.NewScanner(f)
	// Raise the buffer limit to handle long lines without ErrTooLong.
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Try to parse as KEY=VALUE (shell assignment or .env format).
		if idx := strings.IndexByte(line, '='); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			// Strip export prefix common in shell files.
			key = strings.TrimPrefix(key, "export ")
			key = strings.TrimSpace(key)
			value := strings.Trim(line[idx+1:], `"' `)

			if patternName, matched := matchPatterns(key, value); matched {
				findings = append(findings, FileFinding{
					File:          path,
					Line:          lineNum,
					Pattern:       patternName,
					Variable:      key,
					ValueRedacted: redactValue(value),
				})
				continue
			}
		}

		// Fall back: scan the full line for value-only patterns (no NameRegex).
		for _, p := range Patterns {
			if p.NameRegex != nil {
				continue
			}
			if loc := p.Regex.FindStringIndex(line); loc != nil {
				match := line[loc[0]:loc[1]]
				findings = append(findings, FileFinding{
					File:          path,
					Line:          lineNum,
					Pattern:       p.Name,
					Variable:      "",
					ValueRedacted: redactValue(match),
				})
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return findings, err
	}
	return findings, nil
}

// matchPatterns returns the name of the first matching pattern for the given key/value pair.
func matchPatterns(key, value string) (string, bool) {
	if value == "" {
		return "", false
	}
	for _, p := range Patterns {
		if p.NameRegex != nil {
			if key == "" || !p.NameRegex.MatchString(key) {
				continue
			}
			// For the generic pattern, additionally gate on Shannon entropy.
			if p.Name == "Generic API Secret" && shannonEntropy(value) <= 3.5 {
				continue
			}
		}
		if p.Regex.MatchString(value) {
			return p.Name, true
		}
	}
	return "", false
}

// redactValue shows the first 4 and last 3 characters, masking the middle.
// Values shorter than 10 characters are fully masked.
func redactValue(v string) string {
	if len(v) < 10 {
		return "****"
	}
	return v[:4] + "****" + v[len(v)-3:]
}

// shannonEntropy computes the Shannon entropy (bits per character) of s.
func shannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	freq := make(map[rune]int)
	for _, c := range s {
		freq[c]++
	}
	n := float64(len(s))
	var h float64
	for _, count := range freq {
		p := float64(count) / n
		h -= p * math.Log2(p)
	}
	return h
}
