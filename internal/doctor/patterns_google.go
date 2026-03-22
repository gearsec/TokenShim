package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "Google API Key",
		Regex: regexp.MustCompile(`AIza[0-9A-Za-z_-]{35}`),
	})
}
