<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/suggest/normalize.go","sha256":"8ae2ce72978136929ca7cda5e6e0ab10b0e479b5397212d66894338af6b9b6aa","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/suggest/normalize.go

## Company-type words are ignored when matching duplicate names

Lines 15-21

Asked as: company suffix ignored; inc ltd in names; acme inc vs acme company; duplicate suggestion same name different suffix; the in company name

When Turbo compares names and aliases to suggest duplicates for absorb, it drops these words from every name: inc, ltd, limited, corp, corporation, co, company, the, of, and, llc, lp, llp, incorporated, ltee. So "Acme Inc" and "The Acme Company" give the same tokens.

## Name matching ignores case and folds accents and special letters

Lines 23-38

Asked as: accented letters duplicate; case sensitive duplicate names; special letters matching names; é and e same name; ß ø letters in name

Turbo lowercases names and aliases before it compares them for duplicate suggestions. It removes accents, so é becomes e. It also replaces letters that have no accent form: ß becomes ss, æ ae, œ oe, ø o, ł l, đ and ð d, þ th.

## Apostrophes, backticks and full stops are deleted, not turned into spaces

Lines 28-39

Asked as: apostrophe in client name; o'brien duplicate; full stops in initials; dots removed from name; backtick in name matching

In duplicate matching, Turbo deletes the characters ' ’ ` and . so the letters on each side join. O'Brien becomes obrien and R.B.C. becomes rbc.

## Other punctuation splits names into words in duplicate matching

Lines 40-42

Asked as: hyphen in company name; ampersand duplicate suggestion; punctuation splits name; dash separates words; slash in name duplicates

Any other character that is not a letter a-z or a digit 0-9 acts as a space between words. For example, a hyphen or an ampersand separates words. Letters outside a-z that remain after accent folding also act as separators.

## Runs of single letters join into one word in duplicate matching

Lines 44-58

Asked as: single letters joined; R & L duplicate; initials merged into one word; letters separated by spaces match; acronym matching duplicates

When names have consecutive single-letter words, Turbo joins them into one word, so "R & L" becomes rl. A stop word between two single letters does not end the run. A single digit does not join a run.
