package doctor

import "regexp"

func init() {
	RegisterPattern(Pattern{
		Name:  "Stripe Secret Key",
		Regex: regexp.MustCompile(`sk_live_[a-zA-Z0-9]{24,}|sk_test_[a-zA-Z0-9]{24,}`),
	})
}
