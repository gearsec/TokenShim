package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "HuggingFace Token",
		Regex: regexp.MustCompile(`hf_[a-zA-Z0-9]{34,}`),
	})
}
