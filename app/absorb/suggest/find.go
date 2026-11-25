package suggest

import (
	"math"
	"math/bits"
	"sort"
)

// Record is one row to compare: a client or a vendor, with the number of
// references that an absorb would move.
type Record struct {
	ID    string
	Name  string
	Alias string
	Refs  int
}

// Member is a record in a bucket, with its score against the bucket's
// suggested target (1 for the target itself).
type Member struct {
	Record
	Score float64
}

// Bucket is a group of probable duplicates. Members[0] is the suggested
// target: the member with the most references.
type Bucket struct {
	Members []Member
	// Confidence is the mean score of the other members against the target.
	Confidence float64
	// MovedRefs is the references of every member but the target.
	MovedRefs int
	// Rank orders the buckets, highest first:
	// Confidence² × ((members - 1) + log2(1 + MovedRefs)).
	Rank float64
}

// TargetID is the id of the suggested target.
func (b Bucket) TargetID() string {
	return b.Members[0].ID
}

// IDs are the ids of the members, target first.
func (b Bucket) IDs() []string {
	ids := make([]string, len(b.Members))
	for i, member := range b.Members {
		ids[i] = member.ID
	}
	return ids
}

// Find returns the buckets of probable duplicates among records, best first.
// The result depends only on the records and p, not on their order.
//
//  1. Each name and alias is normalised into tokens (Tokens). A token weighs
//     ln(N/df) + 0.5, where df is the number of names that have it.
//  2. Only records that share a matching token are scored (an inverted
//     index), and a pair matches when its score is at least p.MinScore.
//  3. Buckets are stars: the unassigned record with the most unassigned
//     matches, then the most references, then the lowest id, is a centre, and
//     its bucket is the centre with those matches. A match of a match does
//     not join, so buckets do not chain. Only buckets of two or more are kept.
func Find(records []Record, p Params) []Bucket {
	m := newMatcher(records, p)
	adj := m.matches()
	groups := m.stars(adj)
	buckets := make([]Bucket, 0, len(groups))
	for _, group := range groups {
		buckets = append(buckets, m.bucket(group))
	}
	sort.SliceStable(buckets, func(a, b int) bool {
		if buckets[a].Rank != buckets[b].Rank {
			return buckets[a].Rank > buckets[b].Rank
		}
		return buckets[a].TargetID() < buckets[b].TargetID()
	})
	return buckets
}

type matcher struct {
	p     Params
	recs  []Record
	keys  [][]key
	vocab []string
	// digits[t] is true when token t has a digit.
	digits []bool
	// similar holds the similarity of each pair of distinct similar tokens.
	similar map[uint64]float64
	// buffers that keyScore reuses
	units        []unit
	usedA, usedB []bool
}

func newMatcher(records []Record, p Params) *matcher {
	recs := append([]Record(nil), records...)
	sort.SliceStable(recs, func(a, b int) bool { return recs[a].ID < recs[b].ID })
	m := &matcher{p: p, recs: recs, keys: make([][]key, len(recs)), similar: map[uint64]float64{}}

	ids := map[string]int{}
	intern := func(t string) int {
		if id, ok := ids[t]; ok {
			return id
		}
		ids[t] = len(m.vocab)
		m.vocab = append(m.vocab, t)
		return ids[t]
	}
	type raw struct{ name, alias []int }
	raws := make([]raw, len(recs))
	df := map[int]int{}
	for i, r := range recs {
		for _, t := range Tokens(r.Name) {
			raws[i].name = append(raws[i].name, intern(t))
		}
		for _, t := range Tokens(r.Alias) {
			raws[i].alias = append(raws[i].alias, intern(t))
		}
		seen := map[int]bool{}
		for _, t := range raws[i].name {
			if !seen[t] {
				seen[t] = true
				df[t]++
			}
		}
	}
	m.digits = make([]bool, len(m.vocab))
	for t, token := range m.vocab {
		m.digits[t] = hasDigit(token)
	}
	n := float64(len(recs))
	weight := func(t int) float64 {
		d := df[t]
		if d < 1 {
			d = 1
		}
		return math.Log(n/float64(d)) + 0.5
	}
	makeKey := func(toks []int) key {
		k := key{toks: toks, w: make([]float64, len(toks))}
		for i, t := range toks {
			k.w[i] = weight(t)
			k.joined += m.vocab[t]
			if i > 0 {
				k.pairs = append(k.pairs, m.vocab[toks[i-1]]+m.vocab[t])
			}
		}
		return k
	}
	for i := range recs {
		if len(raws[i].name) > 0 {
			m.keys[i] = append(m.keys[i], makeKey(raws[i].name))
		}
		if len(raws[i].alias) > 0 && !sameTokens(raws[i].alias, raws[i].name) {
			m.keys[i] = append(m.keys[i], makeKey(raws[i].alias))
		}
	}
	m.similarTokens()
	return m
}

func sameTokens(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// similarTokens finds every pair of distinct tokens of the vocabulary that
// match, and keeps their similarity. Only tokens of four or more letters and
// no digit can match a token they are not equal to, and only when the longer
// is at most half as long again as the shorter, so the tokens are sorted by
// length and each is compared with the next lengths only. A letter of one
// token that the other lacks needs an edit, which skips most pairs cheaply.
func (m *matcher) similarTokens() {
	cands := []int{}
	for id, t := range m.vocab {
		if len(t) >= minEditLength && !hasDigit(t) {
			cands = append(cands, id)
		}
	}
	sort.Slice(cands, func(a, b int) bool {
		ta, tb := m.vocab[cands[a]], m.vocab[cands[b]]
		if len(ta) != len(tb) {
			return len(ta) < len(tb)
		}
		return ta < tb
	})
	masks := make([]uint32, len(cands))
	for i, id := range cands {
		for _, c := range []byte(m.vocab[id]) {
			masks[i] |= 1 << (c - 'a')
		}
	}
	for i, ia := range cands {
		a := m.vocab[ia]
		for j := i + 1; j < len(cands); j++ {
			b := m.vocab[cands[j]]
			if 2*len(b) > 3*len(a) {
				break
			}
			maxEdits := int(m.p.TokenTolerance*float64(len(b)) + 1e-9)
			samePrefix := m.p.Prefix && a[:minEditLength] == b[:minEditLength]
			if !samePrefix && (maxEdits == 0 ||
				bits.OnesCount32(masks[i]&^masks[j]) > maxEdits ||
				bits.OnesCount32(masks[j]&^masks[i]) > maxEdits) {
				continue
			}
			if s := TokenSimilarity(a, b, m.p); s > 0 {
				m.similar[pairKey(ia, cands[j])] = s
			}
		}
	}
}

// deletions returns s with each one of its letters deleted in turn.
func deletions(s string) []string {
	out := make([]string, 0, len(s))
	for i := range s {
		out = append(out, s[:i]+s[i+1:])
	}
	return out
}

// matches returns, for each record, the records that it matches, sorted. A
// record is scored only against records that share an equal or similar
// token, a whole name joined, or a joined pair of tokens that is within one
// edit of a token of the other.
func (m *matcher) matches() [][]int {
	n := len(m.recs)
	// post: token id -> records with it
	post := map[int][]int{}
	// whole: a whole key joined -> records with it
	whole := map[string][]int{}
	// joins: a joined pair of adjacent tokens, and the pair with one letter
	// deleted -> records with it
	joins := map[string][]int{}
	addOnce := func(index map[string][]int, k string, i int) {
		if list := index[k]; len(list) == 0 || list[len(list)-1] != i {
			index[k] = append(list, i)
		}
	}
	for i := range m.recs {
		for _, k := range m.keys[i] {
			for _, t := range k.toks {
				if list := post[t]; len(list) == 0 || list[len(list)-1] != i {
					post[t] = append(list, i)
				}
			}
			addOnce(whole, k.joined, i)
			for _, joined := range k.pairs {
				addOnce(joins, joined, i)
				if len(joined) >= minJoinLength {
					for _, v := range deletions(joined) {
						addOnce(joins, v, i)
					}
				}
			}
		}
	}
	// near: token id -> similar token ids
	near := map[int][]int{}
	for pk := range m.similar {
		a, b := int(pk>>32), int(pk&0xffffffff)
		near[a] = append(near[a], b)
		near[b] = append(near[b], a)
	}
	// dels: a token, and the token with one letter deleted -> token ids
	dels := map[string][]int{}
	for id, t := range m.vocab {
		dels[t] = append(dels[t], id)
		if len(t) >= minJoinLength {
			for _, v := range deletions(t) {
				dels[v] = append(dels[v], id)
			}
		}
	}

	adj := make([][]int, n)
	stamp := make([]int, n)
	for i := range stamp {
		stamp[i] = -1
	}
	var cands []int
	for i := range m.recs {
		cands = cands[:0]
		add := func(list []int) {
			for _, j := range list {
				if j > i && stamp[j] != i {
					stamp[j] = i
					cands = append(cands, j)
				}
			}
		}
		for _, k := range m.keys[i] {
			add(whole[k.joined])
			for _, t := range k.toks {
				add(post[t])
				for _, u := range near[t] {
					add(post[u])
				}
				// a joined pair of the other record within one edit of t
				token := m.vocab[t]
				add(joins[token])
				if len(token) >= minJoinLength {
					for _, v := range deletions(token) {
						add(joins[v])
					}
				}
			}
			// a token of the other record within one edit of a joined pair
			for _, joined := range k.pairs {
				variants := []string{joined}
				if len(joined) >= minJoinLength {
					variants = append(variants, deletions(joined)...)
				}
				for _, v := range variants {
					for _, u := range dels[v] {
						add(post[u])
					}
				}
			}
		}
		for _, j := range cands {
			if m.recordScore(i, j) >= m.p.MinScore-1e-9 {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}
	for i := range adj {
		sort.Ints(adj[i])
	}
	return adj
}

// stars groups the matches into buckets. Each bucket is a centre, the
// unassigned record with the most unassigned matches (ties: more references,
// then the lower id), with its unassigned matches. Its first member is the
// centre.
func (m *matcher) stars(adj [][]int) [][]int {
	n := len(m.recs)
	used := make([]bool, n)
	free := make([]int, n)
	for i := range adj {
		free[i] = len(adj[i])
	}
	var groups [][]int
	for {
		best := -1
		for i := 0; i < n; i++ {
			if used[i] || free[i] == 0 {
				continue
			}
			if best < 0 || free[i] > free[best] || (free[i] == free[best] && m.recs[i].Refs > m.recs[best].Refs) {
				best = i
			}
		}
		if best < 0 {
			return groups
		}
		group := []int{best}
		for _, j := range adj[best] {
			if !used[j] {
				group = append(group, j)
			}
		}
		for _, u := range group {
			used[u] = true
		}
		for _, u := range group {
			for _, v := range adj[u] {
				free[v]--
			}
		}
		groups = append(groups, group)
	}
}

// bucket scores a group. The suggested target is the member with the most
// references (ties: the centre, then the lower id).
func (m *matcher) bucket(group []int) Bucket {
	target := group[0]
	total := 0
	for _, i := range group {
		total += m.recs[i].Refs
		if m.recs[i].Refs > m.recs[target].Refs {
			target = i
		}
	}
	b := Bucket{Members: make([]Member, 0, len(group))}
	b.Members = append(b.Members, Member{Record: m.recs[target], Score: 1})
	sum := 0.0
	for _, i := range group {
		if i == target {
			continue
		}
		s := m.recordScore(target, i)
		sum += s
		b.Members = append(b.Members, Member{Record: m.recs[i], Score: s})
	}
	others := b.Members[1:]
	sort.SliceStable(others, func(x, y int) bool {
		if others[x].Score != others[y].Score {
			return others[x].Score > others[y].Score
		}
		if others[x].Refs != others[y].Refs {
			return others[x].Refs > others[y].Refs
		}
		return others[x].ID < others[y].ID
	})
	b.Confidence = sum / float64(len(group)-1)
	b.MovedRefs = total - m.recs[target].Refs
	b.Rank = b.Confidence * b.Confidence * (float64(len(group)-1) + math.Log2(1+float64(b.MovedRefs)))
	return b
}
