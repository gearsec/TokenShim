package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "GitHub Token",
		Regex: regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}|gho_[a-zA-Z0-9]{36}|ghs_[a-zA-Z0-9]{36}|github_pat_[a-zA-Z0-9_]{82}`),
	})
}
