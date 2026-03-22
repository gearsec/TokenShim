package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "AWS Access Key ID",
		Regex: regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	})
	RegisterPattern(Pattern{
		Name:      "AWS Secret Access Key",
		Regex:     regexp.MustCompile(`[A-Za-z0-9/+]{40}`),
		NameRegex: regexp.MustCompile(`(?i)AWS_SECRET`),
	})
}
