package fnmatch_test

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/tzvetkoff-go/fnmatch"
)

func here() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", path.Base(file), line)
}

// This is a set of tests ported from a set of tests for C fnmatch
// found at http://www.mail-archive.com/bug-gnulib@gnu.org/msg14048.html
func TestMatch(t *testing.T) {
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		{here(), "", "", 0, true},
		{here(), "*", "", 0, true},
		{here(), "*", "foo", 0, true},
		{here(), "*", "bar", 0, true},
		{here(), "*", "*", 0, true},
		{here(), "**", "f", 0, true},
		{here(), "**", "foo.txt", 0, true},
		{here(), "*.*", "foo.txt", 0, true},
		{here(), "foo*.txt", "foobar.txt", 0, true},
		{here(), "foo.txt", "foo.txt", 0, true},
		{here(), "foo\\.txt", "foo.txt", 0, true},
		{here(), "foo\\.txt", "foo.txt", fnmatch.NoEscape, false},
		{here(), "*", "", fnmatch.Period, true},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestWildcard(t *testing.T) {
	// A wildcard pattern "*" should match anything
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		{here(), "*", "", 0, true},
		{here(), "*", "foo", 0, true},
		{here(), "*", "*", 0, true},
		{here(), "*", "   ", 0, true},
		{here(), "*", ".foo", 0, true},
		{here(), "*", "わたし", 0, true},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestWildcardSlash(t *testing.T) {
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		// Should match "/" when flags are 0
		{here(), "*", "foo/bar", 0, true},
		{here(), "*", "/", 0, true},
		{here(), "*", "/foo", 0, true},
		{here(), "*", "foo/", 0, true},
		// Shouldnt match "/" when flags include `fnmatch.Pathname`
		{here(), "*", "foo/bar", fnmatch.Pathname, false},
		{here(), "*", "/", fnmatch.Pathname, false},
		{here(), "*", "/foo", fnmatch.Pathname, false},
		{here(), "*", "foo/", fnmatch.Pathname, false},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestWildcardPeriod(t *testing.T) {
	// `fnmatch.Period` means that "." is not matched in some circumstances
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		{here(), "*", ".foo", fnmatch.Period, false},
		{here(), "/*", "/.foo", fnmatch.Period, true},
		{here(), "/*", "/.foo", fnmatch.Period | fnmatch.Pathname, false},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestQuestionMark(t *testing.T) {
	// A question mark pattern "?" should match a single character
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		{here(), "?", "", 0, false},
		{here(), "?", "f", 0, true},
		{here(), "?", ".", 0, true},
		{here(), "?", "?", 0, true},
		{here(), "?", "foo", 0, false},
		{here(), "?", "わ", 0, true},
		{here(), "?", "わた", 0, false},
		// Match "/" when flags are 0
		{here(), "?", "/", 0, true},
		// Don't match "/" when flags include `fnmatch.Pathname`
		{here(), "?", "/", fnmatch.Pathname, false},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestQuestionMarkExceptions(t *testing.T) {
	// When flags include `fnmatch.Period` a "?" should not match a "."
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		{here(), "?", ".", fnmatch.Period, false},
		{here(), "foo?", "foo.", fnmatch.Period, true},
		{here(), "/?", "/.", fnmatch.Period, true},
		{here(), "/?", "/.", fnmatch.Period | fnmatch.Pathname, false},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestRange(t *testing.T) {
	azPat := "[a-z]"
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		// Should match a single character inside its range
		{here(), azPat, "a", 0, true},
		{here(), azPat, "q", 0, true},
		{here(), azPat, "z", 0, true},
		{here(), "[わ]", "わ", 0, true},

		// Should not match characters outside its range
		{here(), azPat, "-", 0, false},
		{here(), azPat, " ", 0, false},
		{here(), azPat, "D", 0, false},
		{here(), azPat, "é", 0, false},

		// Should only match one character
		{here(), azPat, "ab", 0, false},
		{here(), azPat, "", 0, false},

		// Should not consume more of the pattern than necessary
		{here(), azPat + "foo", "afoo", 0, true},

		// Should match "-" if it is the first/last character or if backslash is escaped
		{here(), "[-az]", "-", 0, true},
		{here(), "[-az]", "a", 0, true},
		{here(), "[-az]", "b", 0, false},
		{here(), "[az-]", "-", 0, true},
		{here(), "[a\\-z]", "-", 0, true},
		{here(), "[a\\-z]", "b", 0, false},

		// Ignore "\\" when `fnmatch.NoEscape` is given
		{here(), "[a\\-z]", "\\", fnmatch.NoEscape, true},
		{here(), "[a\\-z]", "-", fnmatch.NoEscape, false},

		// Should be negated if starting with "^" or "!"
		{here(), "[^a-z]", "a", 0, false},
		{here(), "[!a-z]", "b", 0, false},
		{here(), "[!a-z]", "é", 0, true},
		{here(), "[!a-z]", "わ", 0, true},

		// Still match "-" if following the negation character
		{here(), "[^-az]", "-", 0, false},
		{here(), "[^-az]", "b", 0, true},

		// Should support multiple characters/ranges
		{here(), "[abc]", "a", 0, true},
		{here(), "[abc]", "c", 0, true},
		{here(), "[abc]", "d", 0, false},
		{here(), "[a-cg-z]", "c", 0, true},
		{here(), "[a-cg-z]", "h", 0, true},
		{here(), "[a-cg-z]", "d", 0, false},

		// Should not match "/" when flags include `fnmatch.Pathname`
		{here(), "[abc/def]", "/", 0, true},
		{here(), "[abc/def]", "/", fnmatch.Pathname, false},
		{here(), "[.-0]", "/", 0, true}, // The range [.-0] includes "/"
		{here(), "[.-0]", "/", fnmatch.Pathname, false},

		// Should normally be case-sensitive
		{here(), "[a-z]", "A", 0, false},
		{here(), "[A-Z]", "a", 0, false},
		// Except when `fnmatch.CaseFold` is given
		{here(), "[a-z]", "A", fnmatch.CaseFold, true},
		{here(), "[A-Z]", "a", fnmatch.CaseFold, true},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestBackslash(t *testing.T) {
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		// A backslash should escape the following characters
		{here(), "\\\\", "\\", 0, true},
		{here(), "\\*", "*", 0, true},
		{here(), "\\*", "foo", 0, false},
		{here(), "\\?", "?", 0, true},
		{here(), "\\?", "f", 0, false},
		{here(), "\\[a-z]", "[a-z]", 0, true},
		{here(), "\\[a-z]", "a", 0, false},
		{here(), "\\foo", "foo", 0, true},
		{here(), "\\わ", "わ", 0, true},

		// Unless `fnmatch.NoEscape` is given
		{here(), "\\\\", "\\", fnmatch.NoEscape, false},
		{here(), "\\\\", "\\\\", fnmatch.NoEscape, true},
		{here(), "\\*", "foo", fnmatch.NoEscape, false},
		{here(), "\\*", "\\*", fnmatch.NoEscape, true},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestLiteral(t *testing.T) {
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		// Literal characters should match themselves
		{here(), "foo", "foo", 0, true},
		{here(), "foo", "foobar", 0, false},
		{here(), "foobar", "foo", 0, false},
		{here(), "foo", "Foo", 0, false},
		{here(), "わたし", "わたし", 0, true},
		// And perform case-folding when FNM_CASEFOLD is given
		{here(), "foo", "FOO", fnmatch.CaseFold, true},
		{here(), "FoO", "fOo", fnmatch.CaseFold, true},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}

func TestLeadingDir(t *testing.T) {
	cases := []struct {
		where    string
		pattern  string
		input    string
		flags    int
		expected bool
	}{
		// `fnmatch.LeadingDir` should ignore trailing "/*"
		{here(), "foo", "foo/bar", 0, false},
		{here(), "foo", "foo/bar", fnmatch.LeadingDir, true},
		{here(), "*", "foo/bar", fnmatch.Pathname, false},
		{here(), "*", "foo/bar", fnmatch.Pathname | fnmatch.LeadingDir, true},
	}

	for _, c := range cases {
		got := fnmatch.Match(c.pattern, c.input, c.flags)
		if got != c.expected {
			t.Errorf(
				" %s:  fnmatch.Match(%q, %q, 0x%02x),  expected: %v,  got: %v",
				c.where, c.pattern, c.input, c.flags, c.expected, got,
			)
		}
	}
}
