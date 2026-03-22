package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "Twilio API Key",
		Regex: regexp.MustCompile(`SK[0-9a-fA-F]{32}`),
	})
}
