package helpcorpus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func testChunk(kind, path string, start int, heading, text string, requires ...[]string) Chunk {
	if requires == nil {
		requires = [][]string{}
	}
	return Chunk{
		ID:       fmt.Sprintf("%s:%d", path, start),
		Path:     path,
		Kind:     kind,
		Start:    start,
		End:      start + strings.Count(text, "\n"),
		Heading:  heading,
		Requires: requires,
		Text:     text,
	}
}

func testIndex(t *testing.T, chunks ...Chunk) *Index {
	t.Helper()
	index, err := NewIndex(&Corpus{Format: FormatVersion, Handbook: "the handbook", Hash: "abc", Chunks: chunks})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { index.Close() })
	return index
}

func search(t *testing.T, index *Index, claims Claims, query string, options SearchOptions) []Hit {
	t.Helper()
	hits, err := index.Search(claims, query, options)
	if err != nil {
		t.Fatalf("search %q: %v", query, err)
	}
	return hits
}

func hitIDs(hits []Hit) string {
	ids := make([]string, len(hits))
	for i, hit := range hits {
		ids[i] = hit.ID
	}
	return strings.Join(ids, " ")
}

var (
	docsChunk  = testChunk(KindDocs, "ui/src/docs/time.md", 20, "Time > Submit a week", "Select **Submit** on the week. Turbo sends the time sheet to your manager.")
	hookChunk  = testChunk(KindCode, "app/hooks/time_entries.go", 24, "", "func cleanTimeEntry(app core.App) {\n\t// remove the properties that the time type does not use\n}")
	limitChunk = testChunk(KindCode, "app/routes/purchase_orders.go", 40, "", "// approvers must have a limit that covers the total\nfunc approvalLimit() {}")
	otherChunk = testChunk(KindNotes, "descriptions/vendors.md", 1, "Vendors", "A vendor has a status. The week does not matter here, and neither does the manager.")

	reportChunk = testChunk(KindNotes, "descriptions/payroll.md", 1, "Payroll > Export", "The payroll export lists the hours for each week.", []string{"report"})
	adminChunk  = testChunk(KindCode, "app/routes/payroll_admin.go", 1, "", "// payroll export for one week\nfunc payrollAdmin() {}", []string{"admin"})
	bothChunk   = testChunk(KindCode, "app/reports/payroll.sql", 1, "", "SELECT hours FROM payroll -- for the week", []string{"report"}, []string{"admin", "commit"})
)

func TestSearchFindsAndRanks(t *testing.T) {
	index := testIndex(t, otherChunk, hookChunk, limitChunk, docsChunk)
	cases := []struct{ name, query, want string }{
		{"a heading match comes first", "submit week", docsChunk.ID},
		{"a question works as a query", "How do I submit my week?", docsChunk.ID},
		{"words find a camelCase name", "clean time entry", hookChunk.ID},
		{"a camelCase name finds itself", "cleanTimeEntry", hookChunk.ID},
		{"a name with an underscore finds its path", "time_entries", hookChunk.ID},
		{"other forms of a word match", "approving limits", limitChunk.ID},
		{"upper case is the same", "VENDOR STATUS", otherChunk.ID},
	}
	for _, c := range cases {
		hits := search(t, index, nil, c.query, SearchOptions{})
		if len(hits) == 0 || hits[0].ID != c.want {
			t.Errorf("%s: %q gave [%s], want %s first", c.name, c.query, hitIDs(hits), c.want)
		}
	}
	hits := search(t, index, nil, "week manager", SearchOptions{})
	if len(hits) != 2 || hits[0].Score < hits[1].Score {
		t.Errorf("hits must come best first, got [%s]", hitIDs(hits))
	}
	if hits := search(t, index, nil, "zebra", SearchOptions{}); len(hits) != 0 {
		t.Errorf("a word that is nowhere gave [%s]", hitIDs(hits))
	}
	if index.Handbook() != "the handbook" || index.Hash() != "abc" {
		t.Error("the index must hand back its corpus's handbook and hash")
	}
}

func TestSearchReturnsOnlyWhatTheClaimsOpen(t *testing.T) {
	index := testIndex(t, docsChunk, reportChunk, adminChunk, bothChunk, otherChunk)
	cases := []struct {
		name   string
		claims Claims
		want   []string
	}{
		{"no claims", nil, nil},
		{"staff", Claims{"time": true}, nil},
		{"report", Claims{"report": true}, []string{reportChunk.ID}},
		{"admin", Claims{"admin": true}, []string{adminChunk.ID}},
		{"commit alone opens nothing", Claims{"commit": true}, nil},
		{"report and commit", Claims{"report": true, "commit": true}, []string{reportChunk.ID, bothChunk.ID}},
		{"report and admin", Claims{"report": true, "admin": true}, []string{reportChunk.ID, adminChunk.ID, bothChunk.ID}},
	}
	for _, c := range cases {
		got := map[string]bool{}
		for _, hit := range search(t, index, c.claims, "payroll", SearchOptions{}) {
			got[hit.ID] = true
		}
		if len(got) != len(c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
		for _, id := range c.want {
			if !got[id] {
				t.Errorf("%s: %s is missing", c.name, id)
			}
		}
		// "week" is in open and closed chunks. The open ones still come back.
		for _, hit := range search(t, index, c.claims, "week", SearchOptions{}) {
			if !hit.Allows(c.claims) {
				t.Errorf("%s: got %s, which these claims don't open", c.name, hit.ID)
			}
		}
	}
	if hits := search(t, index, Claims{"time": true}, "week", SearchOptions{}); len(hits) != 2 {
		t.Errorf("staff must still get the 2 open chunks, got [%s]", hitIDs(hits))
	}
}

func TestTheFilterRunsBeforeTheLimit(t *testing.T) {
	var chunks []Chunk
	for i := 0; i < 12; i++ {
		// The closed chunks match better: the word is in the heading.
		chunks = append(chunks, testChunk(KindCode, fmt.Sprintf("app/closed/%02d.go", i), 1, "zebra zebra", "zebra", []string{"admin"}))
	}
	for i := 0; i < 3; i++ {
		chunks = append(chunks, testChunk(KindCode, fmt.Sprintf("app/open/%02d.go", i), 1, "", "a long line of other words and one zebra"))
	}
	index := testIndex(t, chunks...)

	staff := search(t, index, Claims{"time": true}, "zebra", SearchOptions{Limit: 3})
	if len(staff) != 3 {
		t.Fatalf("staff got %d hits, want the 3 open chunks", len(staff))
	}
	for _, hit := range staff {
		if !strings.HasPrefix(hit.Path, "app/open/") {
			t.Errorf("staff got %s", hit.ID)
		}
	}
	admin := search(t, index, Claims{"admin": true}, "zebra", SearchOptions{Limit: 3})
	if len(admin) != 3 || !strings.HasPrefix(admin[0].Path, "app/closed/") {
		t.Errorf("an admin gets the best matches, got [%s]", hitIDs(admin))
	}
}

func TestSearchLimitsAndKinds(t *testing.T) {
	var chunks []Chunk
	for i := 0; i < 30; i++ {
		kind := KindCode
		if i%3 == 0 {
			kind = KindNotes
		}
		chunks = append(chunks, testChunk(kind, fmt.Sprintf("app/x/%02d.go", i), 1, "", "common words"))
	}
	index := testIndex(t, chunks...)
	for limit, want := range map[int]int{0: 8, -1: 8, 5: 5, 20: 20, 500: 20} {
		if got := len(search(t, index, nil, "common", SearchOptions{Limit: limit})); got != want {
			t.Errorf("limit %d gave %d hits, want %d", limit, got, want)
		}
	}
	notes := search(t, index, nil, "common", SearchOptions{Limit: 20, Kinds: []string{KindNotes}})
	if len(notes) != 10 {
		t.Fatalf("got %d notes, want 10", len(notes))
	}
	for _, hit := range notes {
		if hit.Kind != KindNotes {
			t.Errorf("got a %s chunk", hit.Kind)
		}
	}
}

func TestAnyTextIsASafeQuery(t *testing.T) {
	index := testIndex(t, docsChunk, hookChunk)
	queries := []string{
		"", "   ", `"`, `""`, `"unclosed`, "NEAR(", "a AND", "OR", "NOT week", "body:week", "path : x", "*", "week*",
		"-week", "^week", "(((", ")", "week) OR (1=1", "'; DROP TABLE chunks; --", "{col}: x", "a + b", "\x00",
		"café ☕ 時間", strings.Repeat("week ", 5000), strings.Repeat("a", 100000),
	}
	for _, query := range queries {
		if _, err := index.Search(nil, query, SearchOptions{}); err != nil {
			t.Errorf("%.30q: %v", query, err)
		}
	}
	// A query is words, not syntax: NOT is a word here, so the chunk with
	// "week" still matches.
	if hits := search(t, index, nil, "NOT week", SearchOptions{}); !strings.Contains(hitIDs(hits), docsChunk.ID) {
		t.Errorf("got [%s]", hitIDs(hits))
	}
	if hits := search(t, index, nil, "'; DROP TABLE chunks; --", SearchOptions{}); len(hits) != 0 {
		t.Errorf("got [%s]", hitIDs(hits))
	}
	if hits := search(t, index, nil, "week", SearchOptions{}); len(hits) != 1 {
		t.Error("the index must still work")
	}
}

func TestReadReturnsTheChunksAroundALine(t *testing.T) {
	path := "app/hooks/expenses.go"
	var chunks []Chunk
	for i := 0; i < 12; i++ {
		chunks = append(chunks, testChunk(KindCode, path, i*10+1, "", strings.Repeat("line\n", 8)+"line"))
	}
	chunks = append(chunks, reportChunk, bothChunk)
	index := testIndex(t, chunks...)

	starts := func(chunks []Chunk) string {
		var out []string
		for _, chunk := range chunks {
			out = append(out, fmt.Sprint(chunk.Start))
		}
		return strings.Join(out, " ")
	}
	cases := []struct {
		name     string
		from, to int
		want     string
	}{
		{"one line inside a chunk", 15, 15, "11"},
		{"a range over two chunks", 5, 12, "1 11"},
		{"the gap between chunks", 10, 10, ""},
		{"from a line to the end, at most 8 chunks", 31, 0, "31 41 51 61 71 81 91 101"},
		{"the whole file, at most 8 chunks", 0, 0, "1 11 21 31 41 51 61 71"},
		{"past the end", 500, 600, ""},
	}
	for _, c := range cases {
		if got := starts(index.Read(nil, path, c.from, c.to)); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
	if got := index.Read(nil, "app/missing.go", 0, 0); len(got) != 0 {
		t.Error("a file that is not in the corpus must give nothing")
	}
	if got := index.Read(Claims{"time": true}, reportChunk.Path, 0, 0); len(got) != 0 {
		t.Error("staff must not read a chunk for the report claim")
	}
	if got := index.Read(Claims{"report": true}, reportChunk.Path, 0, 0); len(got) != 1 {
		t.Error("the report claim opens its chunk")
	}
	if got := index.Read(Claims{"report": true}, bothChunk.Path, 0, 0); len(got) != 0 {
		t.Error("a chunk with two groups needs a claim from both")
	}
}

func TestSearchFromManyRequestsAtOnce(t *testing.T) {
	index := testIndex(t, docsChunk, hookChunk, limitChunk, otherChunk, reportChunk)
	var wait sync.WaitGroup
	failures := make(chan string, 64)
	for worker := 0; worker < 16; worker++ {
		wait.Add(1)
		go func(worker int) {
			defer wait.Done()
			claims := Claims{}
			if worker%2 == 0 {
				claims["report"] = true
			}
			for i := 0; i < 40; i++ {
				hits, err := index.Search(claims, "week payroll manager", SearchOptions{})
				want := 2
				if claims["report"] {
					want = 3
				}
				if err != nil || len(hits) != want {
					failures <- fmt.Sprintf("worker %d: %d hits, %v", worker, len(hits), err)
					return
				}
			}
		}(worker)
	}
	wait.Wait()
	close(failures)
	for failure := range failures {
		t.Error(failure)
	}
}

// With rules that restrict part of the real corpus, no search and no read gives
// staff a restricted chunk.
func TestTheFilterHoldsOnTheRealRepo(t *testing.T) {
	repo := filepath.Join("..", "..")
	if _, err := os.Stat(filepath.Join(repo, "ui", "src", "docs")); err != nil {
		t.Skip("the repo's UI sources are not here")
	}
	useRules(t,
		[]pathRule{
			{Prefix: "ui/src/", Requires: everyone},
			{Prefix: "app/", Requires: everyone},
			{Prefix: "app/reports/", Requires: []string{"report"}},
			{Prefix: "descriptions/", Requires: []string{"admin"}},
		},
		// A word that the code holds and the user docs don't.
		[]termRule{{Term: "litestream", Requires: []string{"it", "admin"}}},
	)
	corpus, _, err := Generate(Options{RepoDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	index, err := NewIndex(corpus)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()
	t.Logf("indexed %d chunks in %s", len(corpus.Chunks), time.Since(started))

	closed := 0
	paths := map[string]bool{}
	for _, chunk := range corpus.Chunks {
		paths[chunk.Path] = true
		if !chunk.ForEveryone() {
			closed++
		}
	}
	if closed < 100 {
		t.Fatalf("the rules closed only %d chunks, so this test proves little", closed)
	}

	staff := Claims{"time": true, "po_approver": true, "commit": true}
	queries := []string{
		"payroll", "payroll report export", "branch hours", "payables spreadsheet", "authority matrix",
		"expense commit", "purchase order approval", "time sheet", "admin", "report", "week ending", "SELECT",
		"how does the payroll report work", "who can see payroll", "uncommit", "job", "the", "a",
		"litestream replication check", "cron",
	}
	searched := 0
	for _, query := range queries {
		for _, kinds := range [][]string{nil, {KindNotes}, {KindCode}} {
			hits := search(t, index, staff, query, SearchOptions{Limit: 20, Kinds: kinds})
			searched += len(hits)
			for _, hit := range hits {
				if !hit.Allows(staff) {
					t.Fatalf("%q gave staff %s, which needs %v", query, hit.ID, hit.Requires)
				}
			}
		}
	}
	if searched == 0 {
		t.Fatal("the searches found nothing")
	}
	for path := range paths {
		for _, chunk := range index.Read(staff, path, 0, 0) {
			if !chunk.Allows(staff) {
				t.Fatalf("a read gave staff %s, which needs %v", chunk.ID, chunk.Requires)
			}
		}
	}
	// The same search does reach the closed chunks for a user with the claims.
	reached := 0
	for _, hit := range search(t, index, Claims{"admin": true, "report": true}, "payroll report export", SearchOptions{Limit: 20}) {
		if !hit.ForEveryone() {
			reached++
		}
	}
	if reached == 0 {
		t.Error("a user with the claims must get closed chunks")
	}
}

func TestSearchOnTheRealRepo(t *testing.T) {
	repo := filepath.Join("..", "..")
	corpus, err := Open(repo)
	if errors.Is(err, ErrNotBuilt) {
		t.Skip("the repo's sources are not here, and this binary holds no corpus")
	}
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewIndex(corpus)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	cases := []struct{ query, wantPath string }{
		{"ProcessTimeEntry", "app/hooks/time_entries.go"},
		{"purchase order approver limits", "po_approv"},
		{"copy time entry to tomorrow", "copy"},
		{"expense reimbursement mileage rate", "expense"},
	}
	for _, c := range cases {
		started := time.Now()
		hits := search(t, index, nil, c.query, SearchOptions{})
		elapsed := time.Since(started)
		found := false
		for _, hit := range hits {
			found = found || strings.Contains(strings.ToLower(hit.Path), c.wantPath)
		}
		if !found {
			t.Errorf("%q: no hit in a path with %q, got [%s]", c.query, c.wantPath, hitIDs(hits))
		}
		if elapsed > 250*time.Millisecond {
			t.Errorf("%q took %s", c.query, elapsed)
		}
		t.Logf("%-40q %8s  %s", c.query, elapsed.Round(10*time.Microsecond), hits[0].ID)
	}
}

func TestOpen(t *testing.T) {
	built, _, err := Generate(Options{RepoDir: writeRepo(t, fixture, false)})
	if err != nil {
		t.Fatal(err)
	}
	opened, err := Open(writeRepo(t, fixture, true))
	if err != nil || opened.Hash != built.Hash {
		t.Fatalf("Open with a repo must build from its files, got %v", err)
	}

	// With no repo, the answer is what the build put in the binary, if anything.
	embeddedCorpus, embeddedErr := Embedded()
	for _, dir := range []string{"", filepath.Join(t.TempDir(), "missing")} {
		corpus, err := Open(dir)
		switch {
		case embeddedErr != nil && !errors.Is(err, ErrNotBuilt):
			t.Errorf("Open(%q) with no corpus in the binary gave %v", dir, err)
		case embeddedErr == nil && (err != nil || corpus.Hash != embeddedCorpus.Hash):
			t.Errorf("Open(%q) must give the corpus in the binary, got %v", dir, err)
		}
	}
	if embeddedErr != nil && !errors.Is(embeddedErr, ErrNotBuilt) {
		t.Errorf("the corpus in the binary is bad: %v", embeddedErr)
	}
}
