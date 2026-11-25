package suggest

import (
	"fmt"
	"sort"
)

// Level is how close two names must be to be suggested as duplicates.
type Level string

const (
	Strict Level = "strict"
	Normal Level = "normal"
	Loose  Level = "loose"
)

// Params are the thresholds of a level.
type Params struct {
	// MinScore is the lowest pair score that counts as a match.
	MinScore float64
	// TokenTolerance is the most edits per letter of the longer token for two
	// tokens to match.
	TokenTolerance float64
	// Prefix lets two tokens that share a long prefix match
	// (contractor and contracting).
	Prefix bool
}

// minEditLength is the shortest token that may match a token it is not equal
// to.
const minEditLength = 4

// minJoinLength is the shortest token, on each side, of a joined match within
// one edit.
const minJoinLength = 6

// prefixScore and joinScore are the similarity of a prefix match and of a
// joined match within one edit.
const (
	prefixScore = 0.8
	joinScore   = 0.9
)

// ParseLevel returns the level named by s. A blank s is Normal.
func ParseLevel(s string) (Level, error) {
	switch Level(s) {
	case "", Normal:
		return Normal, nil
	case Strict, Loose:
		return Level(s), nil
	}
	return "", fmt.Errorf("unknown level %q", s)
}

// LevelParams returns the thresholds of a level for a collection. Loose is a
// little stricter for vendors, whose names are shorter and more alike.
func LevelParams(level Level, collection string) Params {
	switch level {
	case Strict:
		return Params{MinScore: 0.90, TokenTolerance: 0.20, Prefix: false}
	case Loose:
		min := 0.70
		if collection == "vendors" {
			min = 0.75
		}
		return Params{MinScore: min, TokenTolerance: 0.25, Prefix: true}
	default:
		return Params{MinScore: 0.80, TokenTolerance: 0.25, Prefix: true}
	}
}

// editDistance is the Damerau-Levenshtein distance between a and b in its
// optimal string alignment form: a transposition of two neighbours is one
// edit. Tokens are ASCII, so it works on bytes.
func editDistance(a, b string) int {
	if a == b {
		return 0
	}
	prevPrev := make([]int, len(b)+1)
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			v := prev[j-1] + cost
			if prev[j]+1 < v {
				v = prev[j] + 1
			}
			if cur[j-1]+1 < v {
				v = cur[j-1] + 1
			}
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] && prevPrev[j-2]+1 < v {
				v = prevPrev[j-2] + 1
			}
			cur[j] = v
		}
		prevPrev, prev, cur = prev, cur, prevPrev
	}
	return prev[len(b)]
}

// editRatio is the edit distance over the longer length. Tokens whose lengths
// differ by more than half can't match, so it returns 1 for them without
// computing the distance.
func editRatio(a, b string) float64 {
	longer, diff := len(a), len(a)-len(b)
	if len(b) > longer {
		longer = len(b)
	}
	if longer == 0 {
		return 1
	}
	if diff < 0 {
		diff = -diff
	}
	if 2*diff > longer {
		return 1
	}
	return float64(editDistance(a, b)) / float64(longer)
}

// TokenSimilarity scores two tokens from 0 to 1. Equal tokens score 1. A
// token with a digit, or a shorter token under four letters, must be equal.
// Otherwise the tokens match when their edit ratio is within the tolerance,
// scoring 1 - ratio, or, with Prefix, when they share a prefix of at least
// four letters and two thirds of the longer token, scoring 0.8.
func TokenSimilarity(a, b string, p Params) float64 {
	if a == b {
		return 1
	}
	if hasDigit(a) || hasDigit(b) {
		return 0
	}
	shorter, longer := len(a), len(b)
	if shorter > longer {
		shorter, longer = longer, shorter
	}
	if shorter < minEditLength {
		return 0
	}
	ratio := editRatio(a, b)
	if ratio <= p.TokenTolerance {
		return 1 - ratio
	}
	if p.Prefix {
		n := 0
		for n < shorter && a[n] == b[n] {
			n++
		}
		if n >= minEditLength && 3*n >= 2*longer {
			return prefixScore
		}
	}
	return 0
}

// joinSimilarity scores two adjacent tokens joined (joined) against one token:
// 1 when equal, 0.9 within one edit when both are at least six letters.
func joinSimilarity(joined, token string) float64 {
	if joined == token {
		return 1
	}
	if hasDigit(joined) || hasDigit(token) || len(joined) < minJoinLength || len(token) < minJoinLength {
		return 0
	}
	if d := len(joined) - len(token); d < -1 || d > 1 {
		return 0
	}
	if editDistance(joined, token) > 1 {
		return 0
	}
	return joinScore
}

// key is one name or alias of a record, as token ids with their weights.
type key struct {
	toks []int
	w    []float64
	// joined is the whole key joined, and pairs[i] tokens i and i+1 joined.
	joined string
	pairs  []string
}

// unit is a possible match between tokens of two keys: one token each, or two
// adjacent tokens on one side joined against one token on the other.
type unit struct {
	ai, an, bi, bn int
	s, w           float64
}

// keyScore scores two keys. Tokens are aligned greedily, the best
// similarity times weight first, in any order. A one-to-one match counts at
// twice the smaller weight and a joined match at the sum of its three
// weights. Unmatched weight on both sides is penalised once more, and an
// unmatched token with a digit makes the score 0.
func (m *matcher) keyScore(a, b *key) float64 {
	if a.joined == b.joined {
		return 1
	}
	units := m.units[:0]
	for i, x := range a.toks {
		for j, y := range b.toks {
			if s := m.tokenSim(x, y); s > 0 {
				units = append(units, unit{i, 1, j, 1, s, 2 * min(a.w[i], b.w[j])})
			}
		}
	}
	units = m.joinUnits(units, a, b, false)
	units = m.joinUnits(units, b, a, true)
	m.units = units
	if len(units) == 0 {
		return 0
	}
	sort.Slice(units, func(p, q int) bool {
		up, uq := units[p], units[q]
		if vp, vq := up.s*up.w, uq.s*uq.w; vp != vq {
			return vp > vq
		}
		if up.ai != uq.ai {
			return up.ai < uq.ai
		}
		if up.bi != uq.bi {
			return up.bi < uq.bi
		}
		if up.an != uq.an {
			return up.an < uq.an
		}
		return up.bn < uq.bn
	})
	usedA := m.usedA[:0]
	for range a.toks {
		usedA = append(usedA, false)
	}
	usedB := m.usedB[:0]
	for range b.toks {
		usedB = append(usedB, false)
	}
	m.usedA, m.usedB = usedA, usedB
	got, matched := 0.0, 0.0
	for _, u := range units {
		free := true
		for i := u.ai; i < u.ai+u.an; i++ {
			free = free && !usedA[i]
		}
		for j := u.bi; j < u.bi+u.bn; j++ {
			free = free && !usedB[j]
		}
		if !free {
			continue
		}
		for i := u.ai; i < u.ai+u.an; i++ {
			usedA[i] = true
		}
		for j := u.bi; j < u.bi+u.bn; j++ {
			usedB[j] = true
		}
		got += u.s * u.w
		matched += u.w
	}
	var leftA, leftB float64
	for i, t := range a.toks {
		if !usedA[i] {
			if m.digits[t] {
				return 0
			}
			leftA += a.w[i]
		}
	}
	for j, t := range b.toks {
		if !usedB[j] {
			if m.digits[t] {
				return 0
			}
			leftB += b.w[j]
		}
	}
	return got / (matched + leftA + leftB + min(leftA, leftB))
}

// joinUnits adds the matches of two adjacent tokens of a joined against one
// token of b. swap puts them the other way round in the units, when a is the
// second key.
func (m *matcher) joinUnits(units []unit, a, b *key, swap bool) []unit {
	for i, joined := range a.pairs {
		first, second := m.vocab[a.toks[i]], m.vocab[a.toks[i+1]]
		for j, t := range b.toks {
			token := m.vocab[t]
			// The token must be longer than each half, or it is a plain
			// token match.
			if len(token) <= len(first) || len(token) <= len(second) {
				continue
			}
			s := joinSimilarity(joined, token)
			if s == 0 {
				continue
			}
			w := a.w[i] + a.w[i+1] + b.w[j]
			if swap {
				units = append(units, unit{j, 1, i, 2, s, w})
			} else {
				units = append(units, unit{i, 2, j, 1, s, w})
			}
		}
	}
	return units
}

// tokenSim looks up the similarity of two token ids. The matcher computes
// every similar pair of the vocabulary once, in similarTokens.
func (m *matcher) tokenSim(a, b int) float64 {
	if a == b {
		return 1
	}
	return m.similar[pairKey(a, b)]
}

func pairKey(a, b int) uint64 {
	if a > b {
		a, b = b, a
	}
	return uint64(a)<<32 | uint64(b)
}

// recordScore is the best score over the name and alias of each record.
func (m *matcher) recordScore(a, b int) float64 {
	best := 0.0
	for x := range m.keys[a] {
		for y := range m.keys[b] {
			if s := m.keyScore(&m.keys[a][x], &m.keys[b][y]); s > best {
				best = s
			}
		}
	}
	return best
}
