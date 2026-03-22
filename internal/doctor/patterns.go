package doctor

import "regexp"

// Pattern describes a single secret detection rule.
type Pattern struct {
	// Name is the human-readable label used in reports (e.g. "AWS Access Key ID").
	Name string
	// Regex is matched against the value (or full line for non-KV file lines).
	Regex *regexp.Regexp
	// NameRegex, if non-nil, is additionally matched against the variable name.
	// The pattern only fires when the name also matches. Use this for generic
	// rules that would produce too many false positives on value alone.
	NameRegex *regexp.Regexp
}

// Patterns is the ordered registry of all active secret detection rules.
// Specific patterns are registered before generic ones to avoid shadowing.
//
// DO NOT edit this slice directly. Call RegisterPattern() from your own
// patterns_<provider>.go file using an init() function — that way your
// contribution is fully self-contained and will never conflict with others.
var Patterns []Pattern

// RegisterPattern appends p to the global pattern registry.
// Call this from an init() function in a provider-specific file:
//
//	func init() {
//	    RegisterPattern(Pattern{
//	        Name:  "My Provider Token",
//	        Regex: regexp.MustCompile(`myp_[a-zA-Z0-9]{32}`),
//	    })
//	}
func RegisterPattern(p Pattern) {
	Patterns = append(Patterns, p)
}
