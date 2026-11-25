package suggest

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

// The names in these tests are inputs to a pure function, not database
// fixtures.

var normal = LevelParams(Normal, "clients")
var strict = LevelParams(Strict, "clients")

func records(names ...string) []Record {
	out := make([]Record, len(names))
	for i, name := range names {
		out[i] = Record{ID: fmt.Sprintf("r%02d", i), Name: name}
	}
	return out
}

// bucketNames lists each bucket's member names, target first.
func bucketNames(buckets []Bucket) [][]string {
	out := [][]string{}
	for _, b := range buckets {
		names := []string{}
		for _, m := range b.Members {
			names = append(names, m.Name)
		}
		out = append(out, names)
	}
	return out
}

func TestTokens(t *testing.T) {
	cases := map[string]string{
		"The Home Depot Canada Inc.": "home depot canada",
		"R & L Trucking Ltd":         "rl trucking",
		"A and B Supply":             "ab supply",
		"O'Brien’s `Bakery`":         "obriens bakery",
		"R.B.C. Dominion":            "rbc dominion",
		"Café Élan Ltée":             "cafe elan",
		"Keewatin-Aski Ltd.":         "keewatin aski",
		"Smith, Co., LLC":            "smith",
		"Unit 12 Holdings":           "unit 12 holdings",
		"Straße":                     "strasse",
		"":                           "",
	}
	for in, want := range cases {
		if got := strings.Join(Tokens(in), " "); got != want {
			t.Errorf("Tokens(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTokenSimilarity(t *testing.T) {
	cases := []struct {
		a, b string
		p    Params
		want float64
	}{
		{"depot", "depot", normal, 1},
		// A shorter token under four letters must be equal.
		{"abc", "abd", normal, 0},
		{"tel", "tell", normal, 0},
		// A token with a digit must be equal.
		{"unit12", "unit13", normal, 0},
		{"route66", "route66", normal, 1},
		// A transposition is one edit: 1 - 1/8.
		{"moncrief", "moncreif", normal, 0.875},
		// One edit in five letters is within Strict's 0.20.
		{"depot", "depat", strict, 0.8},
		// Two edits in eight letters is over Strict's 0.20 but within Normal's 0.25.
		{"armstrong", "armstorng", strict, 1 - 1.0/9},
		{"construction", "constructoin", normal, 1 - 1.0/12},
		{"keewatin", "kewatyin", normal, 0.75},
		{"keewatin", "kewatyin", strict, 0},
		// The prefix rule: contract is 8 letters, at least 4 and 2/3 of 11.
		{"contractor", "contracting", normal, prefixScore},
		{"contractor", "contracting", strict, 0},
		// A shared prefix shorter than 2/3 of the longer token is no match.
		{"construction", "constable", normal, 0},
	}
	for _, c := range cases {
		got := TokenSimilarity(c.a, c.b, c.p)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("TokenSimilarity(%q, %q, %+v) = %v, want %v", c.a, c.b, c.p, got, c.want)
		}
	}
}

func TestEditDistance(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"", "", 0}, {"abc", "", 3}, {"kitten", "sitting", 3}, {"ab", "ba", 1}, {"abcd", "acbd", 1},
	}
	for _, c := range cases {
		if got := editDistance(c.a, c.b); got != c.want {
			t.Errorf("editDistance(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

// score is the record score of the first two names among names, all of
// which set the token weights.
func score(p Params, names ...string) float64 {
	m := newMatcher(records(names...), p)
	return m.recordScore(0, 1)
}

func TestPairScore(t *testing.T) {
	filler := []string{"Other Place", "Northern Supply", "Lakehead Supply", "Kenora Rentals"}
	pair := func(a, b string) float64 {
		return score(normal, append([]string{a, b}, filler...)...)
	}
	if s := pair("George Armstrong Co. Ltd.", "George Armstrong Co"); s != 1 {
		t.Errorf("stop words: got %v, want 1", s)
	}
	// Two adjacent tokens joined equal one token.
	if s := pair("Finn Way Contractors", "Finnway Contractors"); s != 1 {
		t.Errorf("Finn Way: got %v, want 1", s)
	}
	// Or are within one edit of it, when both are six letters or more.
	if s := pair("Tbay Tel Mobility", "Tbaytell Mobility"); s < 0.9 || s >= 1 {
		t.Errorf("Tbay Tel: got %v, want 0.9 to 1", s)
	}
	// Words in another order still match.
	if s := pair("Ontario IO", "IO - Ontario"); s != 1 {
		t.Errorf("order: got %v, want 1", s)
	}
	// An unmatched token with a digit makes the score 0.
	if s := pair("Unit 12 Holdings", "Unit Holdings"); s != 0 {
		t.Errorf("unmatched digit: got %v, want 0", s)
	}
	if s := pair("Route 66 Gas Bar", "Route 61 Gas Bar"); s != 0 {
		t.Errorf("different digits: got %v, want 0", s)
	}
	// Unmatched words on both sides cost more than on one side.
	one := pair("Kestrel Surveying", "Kestrel Surveying Kenora")
	both := pair("Kestrel Surveying Dryden", "Kestrel Surveying Kenora")
	if !(both < one && one < 1) {
		t.Errorf("unmatched: one side %v, both sides %v", one, both)
	}
	// The alias counts: the best of name and alias.
	recs := []Record{{ID: "a", Name: "Richard & Sons Auto"}, {ID: "b", Name: "Dicks Auto Parts", Alias: "Richard and Sons Auto"}}
	m := newMatcher(append(recs, records(filler...)...), normal)
	if s := m.recordScore(0, 1); s != 1 {
		t.Errorf("alias: got %v, want 1", s)
	}
}

func TestRareWordsWeighMore(t *testing.T) {
	// "construction" is in every name, so a difference in the rare word
	// matters more than a difference in the common one.
	names := []string{
		"Moncrief Construction", "Moncrief Constructors",
		"Smith Construction", "Jones Construction", "Brown Construction", "Green Construction",
	}
	m := newMatcher(records(names...), normal)
	rare := m.recordScore(0, 2) // Moncrief vs Smith
	if rare >= normal.MinScore {
		t.Errorf("different rare words scored %v", rare)
	}
}

func TestStarsDoNotChain(t *testing.T) {
	// A chain 0-1-2-3-4. The centre with the most matches, record 1 (ties
	// go to the lower id), takes only its own matches; 3 and 4 form their
	// own bucket.
	m := &matcher{recs: records("a", "b", "c", "d", "e")}
	adj := [][]int{{1}, {0, 2}, {1, 3}, {2, 4}, {3}}
	got := m.stars(adj)
	want := [][]int{{1, 0, 2}, {3, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stars = %v, want %v", got, want)
	}
	// More references break a tie of matches. Then 0 and 4 have no
	// unassigned match left, so they are in no bucket.
	m.recs[2].Refs = 5
	got = m.stars(adj)
	want = [][]int{{2, 1, 3}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stars with refs = %v, want %v", got, want)
	}
}

func TestFindBuckets(t *testing.T) {
	recs := records(
		"Moncrief Construction", "Moncreif Construction Ltd", "Moncrief Construction Limited",
		"Keewatin Aski", "Keewatin-Aski Ltd.",
		"Lakehead Surveys", "Other Place",
		"Route 66 Gas Bar", "Route 61 Gas Bar",
	)
	recs[1].Refs = 12 // the most used Moncrief is the suggested target
	recs[4].Refs = 3
	got := bucketNames(Find(recs, normal))
	want := [][]string{
		{"Moncreif Construction Ltd", "Moncrief Construction", "Moncrief Construction Limited"},
		{"Keewatin-Aski Ltd.", "Keewatin Aski"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buckets = %v, want %v", got, want)
	}
}

func TestBucketScores(t *testing.T) {
	recs := records("Keewatin Aski", "Keewatin-Aski Ltd.", "Kewatin Aski", "Other Place", "Lakehead Surveys")
	recs[0].Refs, recs[1].Refs, recs[2].Refs = 28, 8, 1
	buckets := Find(recs, normal)
	if len(buckets) != 1 {
		t.Fatalf("got %d buckets", len(buckets))
	}
	b := buckets[0]
	if b.TargetID() != "r00" || b.Members[0].Score != 1 {
		t.Fatalf("target = %+v", b.Members[0])
	}
	if b.MovedRefs != 9 {
		t.Errorf("moved refs = %d, want 9", b.MovedRefs)
	}
	conf := (b.Members[1].Score + b.Members[2].Score) / 2
	if math.Abs(b.Confidence-conf) > 1e-12 {
		t.Errorf("confidence = %v, want %v", b.Confidence, conf)
	}
	rank := conf * conf * (2 + math.Log2(10))
	if math.Abs(b.Rank-rank) > 1e-12 {
		t.Errorf("rank = %v, want %v", b.Rank, rank)
	}
}

func TestRankOrder(t *testing.T) {
	recs := records(
		"Alpha Paving", "Alpha Paving Ltd",
		"Bravo Rentals", "Bravo Rentals Inc",
		"Charlie Electric", "Charlie Electric Co", "Charlie Electric Limited",
	)
	recs[2].Refs, recs[3].Refs = 40, 30 // Bravo moves the most references
	got := bucketNames(Find(recs, normal))
	if len(got) != 3 || got[0][0] != "Bravo Rentals" || got[1][0] != "Charlie Electric" || got[2][0] != "Alpha Paving" {
		t.Fatalf("order = %v", got)
	}
}

func TestFindIsDeterministic(t *testing.T) {
	names := []string{
		"Finnway General Contractor Inc", "Finn Way General Contractor", "Finnway general contractor",
		"Finnway General Contracting", "Finnway Contracting Ltd.", "Finnway Contracting",
		"George Armstrong Co.", "George Armstrong Co. Ltd", "George Armstong Co",
		"Red River Co-op", "Red River CO OP", "Red Rver Co-op", "Red River Co-op Kenora",
		"Tbaytel", "TBAY TEL", "Tbtaytel", "Tbaytell",
		"Home Depot", "The Home Depot", "HomeDepot", "Home Depot Canada Inc.",
		"Walmart", "Wal Mart", "Wal-Mart Kenora",
	}
	recs := records(names...)
	for i := range recs {
		recs[i].Refs = (i * 7) % 5 // ties on purpose
	}
	want := Find(recs, normal)
	if len(want) < 5 {
		t.Fatalf("expected several buckets, got %v", bucketNames(want))
	}
	r := rand.New(rand.NewSource(1))
	for n := 0; n < 20; n++ {
		shuffled := append([]Record(nil), recs...)
		r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		if got := Find(shuffled, normal); !reflect.DeepEqual(got, want) {
			t.Fatalf("order %d changed the result:\n%v\nwant\n%v", n, bucketNames(got), bucketNames(want))
		}
	}
}

func TestLevels(t *testing.T) {
	recs := records("Finnway General Contractor", "Finnway General Contracting", "Other Place", "Lakehead Surveys")
	if got := Find(recs, LevelParams(Normal, "clients")); len(got) != 1 {
		t.Errorf("normal: %v", bucketNames(got))
	}
	if got := Find(recs, LevelParams(Strict, "clients")); len(got) != 0 {
		t.Errorf("strict: %v", bucketNames(got))
	}
	if p := LevelParams(Loose, "vendors"); p.MinScore != 0.75 {
		t.Errorf("loose vendors = %v", p.MinScore)
	}
	if p := LevelParams(Loose, "clients"); p.MinScore != 0.70 {
		t.Errorf("loose clients = %v", p.MinScore)
	}
	for _, s := range []string{"", "normal", "strict", "loose"} {
		if _, err := ParseLevel(s); err != nil {
			t.Errorf("ParseLevel(%q): %v", s, err)
		}
	}
	if _, err := ParseLevel("tight"); err == nil {
		t.Error("ParseLevel(tight) passed")
	}
}

func TestDismissals(t *testing.T) {
	d := NewDismissals([][]string{{"a", "b", "c"}, {"d", "e"}})
	cases := []struct {
		ids  []string
		want bool
	}{
		{[]string{"a", "b"}, true},
		{[]string{"c", "a", "b"}, true},
		// A new record joins the bucket: it shows again.
		{[]string{"a", "b", "x"}, false},
		// Members of two dismissed sets are not one dismissed set.
		{[]string{"a", "d"}, false},
		{[]string{"e", "d"}, true},
		{nil, false},
	}
	for _, c := range cases {
		if got := d.Covers(c.ids); got != c.want {
			t.Errorf("Covers(%v) = %v, want %v", c.ids, got, c.want)
		}
	}
}

// BenchmarkFind runs 6,000 synthetic names shaped like the real lists: one or
// two rare words, often a common word, and some typos. On the real client
// list (about 5,700 names) Find takes about 0.2 s.
func BenchmarkFind(b *testing.B) {
	r := rand.New(rand.NewSource(7))
	syllables := strings.Fields("ka ri mo te lu sa no vi re ta pe lo mi da ke su ba ro ni ga")
	rare := make([]string, 4000)
	for i := range rare {
		w := ""
		for n := 2 + r.Intn(3); n > 0; n-- {
			w += syllables[r.Intn(len(syllables))]
		}
		rare[i] = w
	}
	common := strings.Fields("general contracting construction services supply rentals electric plumbing " +
		"paving surveys engineering trucking holdings marine timber forest mining hardware auto parts")
	recs := make([]Record, 6000)
	for i := range recs {
		parts := []string{}
		for n := 1 + r.Intn(2); n > 0; n-- {
			w := []byte(rare[r.Intn(len(rare))])
			if r.Intn(8) == 0 {
				w[r.Intn(len(w))] = byte('a' + r.Intn(26))
			}
			parts = append(parts, string(w))
		}
		if r.Intn(3) > 0 {
			parts = append(parts, common[r.Intn(len(common))])
		}
		recs[i] = Record{ID: fmt.Sprintf("r%05d", i), Name: strings.Join(parts, " "), Refs: r.Intn(20)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Find(recs, normal)
	}
}
