package doctor

import "regexp"

// This file is intentionally named with the zzz_ prefix so that Go's
// alphabetical init() ordering guarantees it registers after all specific
// provider patterns. Any future catch-all patterns must follow the same
// naming convention (zzz_<name>.go) to preserve evaluation order.

func init() {
	RegisterPattern(Pattern{
		// Fires only when the variable name looks like a secret AND the value has
		// high Shannon entropy (>3.5 bits/char, enforced in matchPatterns).
		Name:      "Generic API Secret",
		Regex:     regexp.MustCompile(`.{8,}`),
		NameRegex: regexp.MustCompile(`(?i)(_KEY|_TOKEN|_SECRET|_PASSWORD|API_KEY)$`),
	})
}
