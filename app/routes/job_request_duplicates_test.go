package routes

import (
	"slices"
	"testing"
	"time"
)

func TestSimilarDescriptions(t *testing.T) {
	cases := []struct {
		name string
		a, b string
		want bool
	}{
		{"same words in another order", "Highway 17 culvert replacement design", "Culvert replacement on Highway 17 north", true},
		{"case and punctuation do not count", "HIGHWAY-17 Culvert, Replacement", "highway 17 culvert replacement", true},
		{"one shared word is not enough", "Arena roof replacement", "Culvert replacement on Highway 17", false},
		{"two shared words of a long description are not enough", "Water treatment plant expansion design and review", "Water treatment chemical audit at Elm", false},
		{"two shared words of three are enough", "Bridge deck survey", "Bridge deck rehabilitation design and tender", true},
		{"stopwords and short words are not shared words", "The project for a new phase of work", "The project for the new phase of all works", false},
		{"a one-word description must be shared", "Sidewalks", "Sidewalks inventory and condition", true},
		{"a one-word description that is not shared", "Sidewalks", "Arena roof", false},
		{"a two-character word with a digit counts", "Highway 17 lighting", "Highway 17 resurfacing", true},
		{"other two-character words do not count", "ab cd ef", "ab cd ef", false},
		{"blank descriptions are never similar", "", "", false},
		{"one blank description", "Highway 17 culvert", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := similarDescriptions(c.a, c.b); got != c.want {
				t.Fatalf("similarDescriptions(%q, %q) = %v, want %v (words %v / %v)", c.a, c.b, got, c.want, descriptionWords(c.a), descriptionWords(c.b))
			}
			if got := similarDescriptions(c.b, c.a); got != c.want {
				t.Fatalf("similarDescriptions is not symmetric for %q and %q", c.a, c.b)
			}
		})
	}
}

func TestRecentProject(t *testing.T) {
	now := time.Date(2026, time.November, 20, 15, 0, 0, 0, time.UTC)
	cases := []struct {
		name     string
		award    string
		lastTime string
		want     bool
		wantLast string
	}{
		{"recent award", "2026-03-02", "", true, "2026-03-02"},
		{"old award but recent time", "2023-01-10", "2026-08-04", true, "2026-08-04"},
		{"old award and old time", "2022-01-03", "2024-05-06", false, "2024-05-06"},
		{"no award and no time", "", "", false, ""},
		{"first day of the window", "2025-11-20", "", true, "2025-11-20"},
		{"day before the window", "2025-11-19", "", false, "2025-11-19"},
		{"award in the future", "2027-01-15", "", true, "2027-01-15"},
		{"an award that is not a date is ignored", "unknown", "2026-01-05", true, "2026-01-05"},
		{"a timestamp keeps its day", "2026-02-01 00:00:00.000Z", "", true, "2026-02-01"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := recentProject(c.award, c.lastTime, now); got != c.want {
				t.Fatalf("recentProject(%q, %q) = %v, want %v", c.award, c.lastTime, got, c.want)
			}
			if got := lastActivity(c.award, c.lastTime); got != c.wantLast {
				t.Fatalf("lastActivity(%q, %q) = %q, want %q", c.award, c.lastTime, got, c.wantLast)
			}
		})
	}
}

func TestDuplicateReasons(t *testing.T) {
	request := duplicateSubject{ID: "request", Client: "c1", Proposal: "p1", Parent: "parent", ClientPO: "PO-1", ClientReference: "REF-1", Location: "LOC", Description: "Highway 17 culvert replacement"}
	cases := []struct {
		name           string
		other          duplicateSubject
		wantReasons    []string
		wantExact      bool
		wantSameClient bool
	}{
		{"location alone is not exact and not the same client", duplicateSubject{Client: "c2", Location: "LOC", Description: "Highway 17 culvert replacement"}, []string{"Same location"}, false, false},
		{"another client's description is not compared", duplicateSubject{Client: "c2", Description: "Highway 17 culvert replacement"}, []string{}, false, false},
		{"same client with a similar description", duplicateSubject{Client: "c1", Description: "Culvert replacement Highway 17"}, []string{"Same client", "Similar description"}, false, true},
		{"client PO ignores case and spaces", duplicateSubject{Client: "c2", ClientPO: " po-1 "}, []string{"Same client PO"}, true, false},
		{"client reference", duplicateSubject{Client: "c2", ClientReference: "REF-1"}, []string{"Same client reference"}, true, false},
		{"same proposal", duplicateSubject{Proposal: "p1"}, []string{"Same proposal"}, true, false},
		{"sibling subjob", duplicateSubject{Client: "c1", Parent: "parent"}, []string{"Same client", "Same parent"}, true, true},
		{"the parent project", duplicateSubject{ID: "parent", Client: "c1"}, []string{"Same client", "Parent project"}, true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reasons, exact, sameClient := duplicateReasons(request, c.other)
			if !slices.Equal(reasons, c.wantReasons) || exact != c.wantExact || sameClient != c.wantSameClient {
				t.Fatalf("got %v exact=%v sameClient=%v, want %v exact=%v sameClient=%v", reasons, exact, sameClient, c.wantReasons, c.wantExact, c.wantSameClient)
			}
		})
	}
	// Blank fields never match.
	reasons, exact, sameClient := duplicateReasons(duplicateSubject{}, duplicateSubject{})
	if len(reasons) != 0 || exact || sameClient {
		t.Fatalf("blank fields matched: %v exact=%v sameClient=%v", reasons, exact, sameClient)
	}
}
