package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "Anthropic API Key",
		Regex: regexp.MustCompile(`sk-ant-[a-zA-Z0-9_-]{20,}`),
	})
}
