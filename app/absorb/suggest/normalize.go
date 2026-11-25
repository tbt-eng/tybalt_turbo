// Package suggest finds probable duplicate records by name and alias, for the
// absorb suggestions page. It is pure: Find takes rows and returns buckets,
// with no database, so it can be tested with plain lists of names. The routes
// in app/routes/absorb_suggestions.go load the rows. See descriptions/Absorb.md,
// "Absorb suggestions".
package suggest

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// stopWords are dropped from every name: they say what kind of company it is,
// not which one.
var stopWords = map[string]bool{
	"inc": true, "ltd": true, "limited": true, "corp": true, "corporation": true,
	"co": true, "company": true, "the": true, "of": true, "and": true,
	"llc": true, "lp": true, "llp": true, "incorporated": true, "ltee": true,
}

// letters that do not decompose into a base letter and a mark.
var foldLetters = strings.NewReplacer(
	"ß", "ss", "æ", "ae", "œ", "oe", "ø", "o", "ł", "l", "đ", "d", "ð", "d", "þ", "th",
)

// deleted characters join the letters on each side: O'Brien is obrien and
// R.B.C. is rbc.
var deleted = strings.NewReplacer("'", "", "’", "", "`", "", ".", "")

// Tokens normalises a name or an alias into its tokens: lowercase, accents
// folded, apostrophes, backticks and full stops deleted, any other character
// that is not a letter or a digit a space, stop words dropped, and runs of
// single letters joined ("R & L" is rl).
func Tokens(s string) []string {
	s = foldLetters.Replace(strings.ToLower(s))
	s = stripMarks(s)
	s = deleted.Replace(s)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	})
	var out []string
	// run is true while the last token kept is a run of single letters. A
	// stop word between two single letters does not end the run.
	run := false
	for _, w := range words {
		if stopWords[w] {
			continue
		}
		single := len(w) == 1 && w[0] >= 'a' && w[0] <= 'z'
		if single && run {
			out[len(out)-1] += w
			continue
		}
		run = single
		out = append(out, w)
	}
	return out
}

// stripMarks folds accents: é is e.
func stripMarks(s string) string {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		return s
	}
	var b strings.Builder
	for _, r := range norm.NFD.String(s) {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func hasDigit(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}
