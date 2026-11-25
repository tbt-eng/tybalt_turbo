<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/suggest/match.go","sha256":"aec3400c7eb78822c3b587eeffac2ef2db9327ec23ebccb4f5b11db1139c1a26","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/suggest/match.go

## Three levels for duplicate name suggestions

Lines 8-15

Asked as: duplicate suggestion level; strict or loose matching; how close names must be; duplicate name sensitivity

Duplicate suggestions run at one of three levels: strict, normal or loose (Level). The level sets how close two names must be to be suggested as duplicates.

## What each level threshold controls

Lines 17-27

Asked as: what the level threshold does; minimum match score setting; how many typos allowed; prefix setting for names

Each level has a minimum pair score that counts as a match (MinScore). It also has a token tolerance, the most edits per letter of the longer word for two words to match (TokenTolerance). A Prefix setting lets two words that share a long prefix match, such as contractor and contracting.

## Minimum token lengths for edit and joined matches

Lines 29-35

Asked as: minimum word length for matching; short words not matching; joined words too short; why short names don't match

A word must be at least 4 letters (minEditLength) to match a word it is not equal to. A joined match within one edit needs at least 6 letters (minJoinLength) on each side.

## Prefix match and joined match similarity scores

Lines 37-42

Asked as: prefix match score; joined match score value; how similar scores are set; score for shared beginning

A shared-prefix match between two words scores 0.8 (prefixScore). A joined match within one edit scores 0.9 (joinScore).

## Absorb suggestion level defaults to Normal when blank

Lines 44-53

Asked as: no level chosen for absorb; default absorb level; unknown level error; invalid duplicate level

When no level is given for absorb suggestions, Turbo uses Normal. It accepts only the levels strict, normal and loose. Any other level is refused with the error 'unknown level' and the value.

## Strict level thresholds for duplicate suggestions

Lines 59-60

Asked as: strict duplicate matching; strict level threshold; names must be almost identical; no prefix matching strict

At the Strict level, a pair of names must score at least 0.90 to be suggested as duplicates. Two words may differ by at most 0.20 edits per letter of the longer word. Words that only share a prefix do not match (Prefix off).

## Loose level thresholds, stricter for vendors

Lines 61-66

Asked as: loose duplicate matching; loose vendor threshold; vendors stricter duplicates; lenient name matching

At the Loose level, clients need a pair score of at least 0.70 and vendors need at least 0.75. Words may differ by at most 0.25 edits per letter of the longer word. Words that share a long prefix can match.

## Normal level thresholds for duplicate suggestions

Lines 67-69

Asked as: normal duplicate matching; default duplicate threshold; normal level score; standard name matching

At the Normal level, a pair of names must score at least 0.80 to be suggested as duplicates. Two words may differ by at most 0.25 edits per letter of the longer word. Words that share a long prefix can match (Prefix on).

## Typo tolerance between words and its score

Lines 72-127

Asked as: typo tolerance; spelling mistakes allowed; transposed letters count; word length difference too big

Two words match when their edits divided by the longer length are within the level's tolerance (equal values count). The word score is 1 minus that ratio. A transposition of two neighbouring letters counts as one edit. Words whose lengths differ by more than half never match this way.

## Words with digits must match exactly

Lines 138-140

Asked as: numbers in names must match; digits not matching; near number match; names with numbers different

A word containing a digit matches another word only if the two are identical. Similar spellings such as near numbers score 0.

## Short words must be equal to match

Lines 141-147

Asked as: short word not matching; three letter word match; words under four letters; exact match only short

A word with fewer than 4 letters matches another word only when the two are exactly equal. This applies when the shorter of the two words is under 4 letters (minEditLength).

## Shared-prefix match scores 0.8

Lines 152-160

Asked as: shared beginning of words; contractor contracting match; same start different ending; prefix matches words

When Prefix is on and the edit check fails, two words still match if they share a prefix of at least 4 letters that is at least two thirds of the longer word (contractor and contracting). Such a match scores 0.8.

## Joined words match a single word

Lines 164-180

Asked as: two words joined together; space missing in name; merged words match; joined name matching

Two adjacent words joined together can match one word on the other name. Equal scores 1. Within one edit scores 0.9, and only when both sides have at least 6 letters and differ in length by at most 1. The single word must be longer than each half, and no digits are allowed.

## How two names are scored by word alignment

Lines 198-284

Asked as: how names are scored; name similarity calculation; word alignment score; why score is zero

Turbo pairs words greedily, best similarity times weight first, in any order. A one-to-one match counts at twice the smaller weight. A joined match counts at the sum of its three weights. Unmatched weight on both sides is penalised once more. The score is 0 when no words match.

## Identical joined names score 1

Lines 203-206

Asked as: same name without spaces; names equal when joined; identical after joining words; spacing difference same name

When two names are the same once their words are joined, they score 1 without further comparison.

## Joined match weight is the sum of three weights

Lines 211-211

Asked as: joined match weight; how words are weighted; word weight in score; weight of joined words

A joined match counts at the sum of the weights of the two joined words and the single word. A one-to-one word match counts at twice the smaller of the two weights.

## Matches are chosen greedily, best score times weight first

Lines 218-266

Asked as: order of word matching; greedy word pairing; which words get matched; best match first

Turbo sorts the possible matches by similarity times weight, highest first. It takes each match only if none of its words is already used. A name pair with no possible matches scores 0.

## Unmatched word with a digit zeroes the score

Lines 267-283

Asked as: digit word unmatched score zero; number word left out; score zero because of number; unmatched numbered word

If any word with a digit is left unmatched on either name, the pair scores 0, however well the other words match.

## Unmatched words lower the score of a name pair

Lines 267-284

Asked as: unmatched words lower score; extra words reduce match; leftover words penalty; score drops for extra words

The score is the matched total divided by the matched weight plus the unmatched weight on both sides plus the smaller unmatched side again. Unmatched words on both sides therefore cost more than on one side.

## Joined match needs a longer word than each half

Lines 295-299

Asked as: joined word too short; joined match not added; word must be longer than halves; joined match rejected

When Turbo joins two adjacent words to match one word on the other side, that word must be longer than each half. Otherwise Turbo treats it as an ordinary word match and does not add a joined match.

## Joined-word candidate is dropped when it scores zero

Lines 300-310

Asked as: joined candidate dropped; joined similarity zero; no joined match found; joined words ignored

Turbo joins two adjacent words of one name and compares them with one word of the other name. It keeps the candidate only when the joined similarity is above 0. Its weight is the sum of the weights of the two joined words and the single word.

## Word similarity is looked up from precomputed pairs

Lines 315-322

Asked as: word similarity lookup; similar words list; word pair not found score; precomputed word similarity

Equal words score 1. For other words, Turbo uses the similarity it computed once for every similar pair of words in its vocabulary (similarTokens). A pair that is not in that list scores 0.

## Record score uses best of names and aliases

Lines 331-342

Asked as: alias matching score; record alias duplicates; score using aliases; other names of record

A record's score against another is the best score over every pairing of the name and aliases of each record.
