package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Find a rate sheet by job or client (issue 207, step 2). The fixtures are
// rows of the seed CSV files (app/test_seed_data/data); no test changes them.
//
//   - Rollout Rates (2026-11-00026): rev 0 (rsrollout0r0001) is the current
//     sheet of 91-8930 to 91-8934; rev 1 (rsrollout1r0001) of 91-8935. All
//     are BN Rail Corporation jobs.
//   - Absorb History Revised (2026-11-00022): only a past period of 91-8916
//     is on rev 0 (rsabhistrev0r01); rev 1 (rsabhistrev1r01) has no jobs.
//   - 2025 Standard Rates (c41ofep525bcacj) is the current sheet of jobs of
//     three clients, and a past one of 91-8934, whose opening period is on it.
//   - Sheet A rev 1 (rs_sheeta_r1) has no jobs in any revision.
//   - rsdraftsheet001 is a request's proposed rates.

type rateSheetFindFields struct {
	JobCount    int    `db:"job_count"`
	JobNumbers  string `db:"job_numbers"`
	ClientNames string `db:"client_names"`
}

func rateSheetFindRow(t *testing.T, app core.App, id string) rateSheetFindFields {
	t.Helper()
	var row rateSheetFindFields
	err := app.DB().NewQuery("SELECT job_count, job_numbers, client_names FROM rate_sheets_augmented WHERE id = {:id}").
		Bind(dbx.Params{"id": id}).One(&row)
	if err != nil {
		t.Fatalf("%s: %v", id, err)
	}
	return row
}

func commaList(s string) []string {
	if s == "" {
		return []string{}
	}
	list := strings.Split(s, ",")
	slices.Sort(list)
	return list
}

// The view lists the job numbers and client names of every revision of the
// code, current and past, and job_count still counts this revision only.
func TestRateSheetsAugmentedFindByJob(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	rollout := []string{"91-8930", "91-8931", "91-8932", "91-8933", "91-8934", "91-8935"}
	for id, want := range map[string]struct {
		count   int
		numbers []string
		clients []string
	}{
		// A job on rev 0 finds rev 1, and rev 1's job finds rev 0.
		"rsrollout0r0001": {5, rollout, []string{"BN Rail Corporation"}},
		"rsrollout1r0001": {1, rollout, []string{"BN Rail Corporation"}},
		// A past period counts, on its revision and on the others.
		"rsabhistrev0r01": {1, []string{"91-8916"}, []string{"BN Rail Corporation"}},
		"rsabhistrev1r01": {0, []string{"91-8916"}, []string{"BN Rail Corporation"}},
		// No jobs gives empty text.
		"rs_sheeta_r1": {0, []string{}, []string{}},
	} {
		row := rateSheetFindRow(t, app, id)
		if row.JobCount != want.count {
			t.Errorf("%s job_count = %d; want %d", id, row.JobCount, want.count)
		}
		if got := commaList(row.JobNumbers); !slices.Equal(got, want.numbers) {
			t.Errorf("%s job_numbers = %v; want %v", id, got, want.numbers)
		}
		if got := commaList(row.ClientNames); !slices.Equal(got, want.clients) {
			t.Errorf("%s client_names = %v; want %v", id, got, want.clients)
		}
	}

	// A sheet with jobs of several clients lists each client once.
	standard := rateSheetFindRow(t, app, "c41ofep525bcacj")
	if got, want := commaList(standard.ClientNames), []string{"BN Rail Corporation", "Mobilia SA", "Sepulchi & Co"}; !slices.Equal(got, want) {
		t.Errorf("standard client_names = %v; want %v", got, want)
	}
	if !slices.Contains(commaList(standard.JobNumbers), "91-8934") {
		t.Errorf("standard job_numbers = %q; want the past job 91-8934", standard.JobNumbers)
	}
}

// The collection API returns the new fields as text, empty when there are no
// jobs, and still leaves out a request's proposed rates.
func TestRateSheetsAugmentedFindByJobAPI(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	status, body := absorbRequest(t, app, http.MethodGet, "/api/collections/rate_sheets_augmented/records?perPage=500", "author@soup.com", nil)
	if status != http.StatusOK {
		t.Fatalf("status = %d; body %s", status, body)
	}
	var list struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(body), &list); err != nil {
		t.Fatal(err)
	}
	byID := map[string]map[string]any{}
	for _, item := range list.Items {
		byID[item["id"].(string)] = item
	}
	if _, ok := byID["rsdraftsheet001"]; ok {
		t.Fatal("the view lists a request's proposed rates")
	}
	for id, want := range map[string]string{"rs_sheeta_r1": "", "rsabhistrev1r01": "91-8916"} {
		item, ok := byID[id]
		if !ok {
			t.Fatalf("%s is not listed", id)
		}
		if got, ok := item["job_numbers"].(string); !ok || got != want {
			t.Errorf("%s job_numbers = %#v; want %q", id, item["job_numbers"], want)
		}
		if _, ok := item["client_names"].(string); !ok {
			t.Errorf("%s client_names = %#v; want text", id, item["client_names"])
		}
	}
}
