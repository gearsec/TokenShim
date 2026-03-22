package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "Slack Token",
		Regex: regexp.MustCompile(`xox[baprs]-[a-zA-Z0-9-]{10,}`),
	})
}
