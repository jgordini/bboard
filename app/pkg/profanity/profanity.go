package profanity

import (
	"regexp"
	"strings"
)

// Word-boundary, case-insensitive regex for blocked terms.
// Compiled once at init.
var blockRegex *regexp.Regexp

func init() {
	// Seven dirty words (George Carlin) + common variants
	dirty := []string{
		`shit`, `piss`, `fuck`, `cunt`, `cocksucker`, `motherfucker`, `tits`,
		`shitty`, `fucking`, `fucker`, `fucked`, `fck`, `fuk`, `cunts`,
	}
	// Common racist/hate slurs (unambiguous)
	slurs := []string{
		`nigger`, `nigga`, `chink`, `gook`, `spic`, `kike`,
		`fag`, `faggot`, `dyke`, `tranny`, `retard`, `retarded`, `faggots`,
	}
	all := append(dirty, slurs...)
	// Match whole words only, case-insensitive
	pattern := `(?i)\b(` + strings.Join(all, `|`) + `)\b`
	blockRegex = regexp.MustCompile(pattern)
}

// ContainsProfanity reports whether text contains a blocked word.
// Returns (true, matchedWord) if found, (false, "") otherwise.
// matchedWord is the first match (lowercase) for logging; do not show to end users.
func ContainsProfanity(text string) (bool, string) {
	if text == "" {
		return false, ""
	}
	m := blockRegex.FindStringSubmatch(text)
	if len(m) < 2 {
		return false, ""
	}
	return true, strings.ToLower(m[1])
}
