package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "OpenAI API Key",
		Regex: regexp.MustCompile(`sk-proj-[a-zA-Z0-9_-]{20,}|sk-[a-zA-Z0-9]{20,}`),
	})
}
