package fnmatch

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Constants ...
const (
	NoEscape = 1 << iota
	Pathname
	Period
	LeadingDir
	CaseFold
)

// Searches a string for an ASCII char and returns its offset or -1
func strchr(s string, c rune) int {
	for i, sc := range s {
		if sc == c {
			return i
		}
	}

	return -1
}

// Unpacks a rune from string and advances the string pointer
func unpackRune(s *string) rune {
	r, size := utf8.DecodeRuneInString(*s)
	*s = (*s)[size:]
	return r
}

// Match ...
func Match(pattern string, str string, flags ...int) bool {
	combinedFlags := 0
	for _, flag := range flags {
		combinedFlags ^= flag
	}

	// Flags
	flagNoEscape := combinedFlags&NoEscape != 0
	flagPathname := combinedFlags&Pathname != 0
	flagPeriod := combinedFlags&Period != 0
	flagLeadingDir := combinedFlags&LeadingDir != 0
	flagCaseFold := combinedFlags&CaseFold != 0

	// Unpacks a rune from `str` and keeps some info about previous rune and string position
	strPrevAtStart := true
	strCurrAtStart := true
	strPrevUnpacked := rune(0)
	strCurrUnpacked := rune(0)
	unpackStr := func() rune {
		strPrevAtStart = strCurrAtStart
		strCurrAtStart = false
		strPrevUnpacked = strCurrUnpacked
		strCurrUnpacked = unpackRune(&str)
		return strCurrUnpacked
	}

	for len(pattern) > 0 {
		patternChar := unpackRune(&pattern)

		switch patternChar {
		case '?':
			// Any char
			if len(str) == 0 {
				return false
			}

			strChar := unpackStr()

			if strChar == '/' && flagPathname {
				return false
			}
			if strChar == '.' && flagPeriod && (strPrevAtStart || (flagPathname && strPrevUnpacked == '/')) {
				return false
			}

		case '*':
			// Collapse multiple "*"
			for len(pattern) > 0 && pattern[0] == '*' {
				pattern = pattern[1:]
			}

			if len(str) > 0 && str[0] == '.' && flagPeriod &&
				(strCurrAtStart || (flagPathname && strCurrUnpacked == '/')) {
				return false
			}

			// "*" at the end
			if len(pattern) == 0 {
				if flagPathname {
					if flagLeadingDir || !strings.ContainsRune(str, '/') {
						return true
					}

					return false
				}

				return true
			}

			// "*/"
			if pattern[0] == '/' && flagPathname {
				offset := strchr(str, '/')
				if offset == -1 {
					return false
				}

				str = str[offset:]
				unpackStr()
				pattern = pattern[1:]
				break
			}

			// General case - recurse
			for test := str; len(test) > 0; unpackRune(&test) {
				if Match(pattern, test, (combinedFlags & ^Period)) {
					return true
				}

				if flagPathname && test[0] == '/' {
					break
				}
			}

			return false

		case '[':
			// Range
			if len(str) == 0 {
				return false
			}

			if flagPathname && str[0] == '/' {
				return false
			}

			strChar := unpackStr()
			if !matchRange(&pattern, strChar, combinedFlags) {
				return false
			}

		case '\\':
			// Escape char
			if !flagNoEscape {
				if len(pattern) > 0 {
					patternChar = unpackRune(&pattern)
				}
			}

			fallthrough

		default:
			if len(str) == 0 {
				return false
			}
			strChar := unpackStr()
			switch {
			case strChar == patternChar:
			case flagCaseFold && unicode.ToLower(strChar) == unicode.ToLower(patternChar):
			default:
				return false
			}
		}
	}

	return len(str) == 0 || (flagLeadingDir && str[0] == '/')
}

func matchRange(pattern *string, test rune, flags int) bool {
	if len(*pattern) == 0 {
		return false
	}

	flagNoEscape := flags&NoEscape != 0
	flagCaseFold := flags&CaseFold != 0

	if flagCaseFold {
		test = unicode.ToLower(test)
	}

	var negate, matched bool
	if (*pattern)[0] == '^' || (*pattern)[0] == '!' {
		negate = true
		(*pattern) = (*pattern)[1:]
	}

	for !matched && len(*pattern) > 1 && (*pattern)[0] != ']' {
		c := unpackRune(pattern)

		if !flagNoEscape && c == '\\' {
			if len(*pattern) > 1 {
				c = unpackRune(pattern)
			} else {
				return false
			}
		}

		if flagCaseFold {
			c = unicode.ToLower(c)
		}

		if (*pattern)[0] == '-' && len(*pattern) > 1 && (*pattern)[1] != ']' {
			unpackRune(pattern) // skip the -
			c2 := unpackRune(pattern)

			if !flagNoEscape && c2 == '\\' {
				if len(*pattern) > 0 {
					c2 = unpackRune(pattern)
				} else {
					return false
				}
			}

			if flagCaseFold {
				c2 = unicode.ToLower(c2)
			}

			// This really should be more intelligent, but it looks like
			// fnmatch.c does simple int comparisons, therefore we will as well
			if c <= test && test <= c2 {
				matched = true
			}
		} else if c == test {
			matched = true
		}
	}

	// Skip past the rest of the pattern
	ok := false
	for !ok && len(*pattern) > 0 {
		c := unpackRune(pattern)
		if c == '\\' && len(*pattern) > 0 {
			unpackRune(pattern)
		} else if c == ']' {
			ok = true
		}
	}

	return ok && matched != negate
}
