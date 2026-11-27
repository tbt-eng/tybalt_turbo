package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Find by rates in a project request (issue 207, step 3). The fixtures add
// four sheets with three roles each, created for these tests:
//
//	Sheet                         Principal        Administration  Senior Designer
//	rsfindrates0001 Alpha (00028)  4321 / 5617.30   111 / 144.30    222 / 288.60
//	rsfindrates0002 Beta  (00029)  4321 / 5617.30   111 / 150       223 / 289.90
//	rsfindrates0003 Gamma (00030)  4321 / 5617.30   (none)          222 / 288.60
//	rsfindrates0004 Inactive       4321 / 5617.30   111 / 144.30    225 / 292.50
//
// Alpha, Beta and Gamma are active; the fourth is not. The draft
// rsdraftsheet001 (a request's proposed rates) has the Administration Staff
// rate of the 2025 Standard Rates, 90 / 117.

const (
	findPath       = "/api/job_requests/find_rates"
	findAlpha      = "rsfindrates0001"
	findBeta       = "rsfindrates0002"
	findGamma      = "rsfindrates0003"
	findInactive   = "rsfindrates0004"
	adminRole      = "2bd6pxxof7p649x"
	designerRole   = "t7alk04t6yqsexe"
	findRatesEmail = workflowAuthor
)

func findRate(role string, rate float64) map[string]any {
	return map[string]any{"role": role, "rate": rate}
}

func findRateOvertime(role string, rate, overtime float64) map[string]any {
	return map[string]any{"role": role, "rate": rate, "overtime_rate": overtime}
}

// findMatchIDs returns the ids of the matches, in order.
func findMatchIDs(t *testing.T, result map[string]any) []string {
	t.Helper()
	ids := []string{}
	matches, ok := result["matches"].([]any)
	if !ok {
		t.Fatalf("no matches list: %#v", result)
	}
	for _, match := range matches {
		ids = append(ids, match.(map[string]any)["id"].(string))
	}
	return ids
}

func TestFindByRatesMatchesEveryEnteredRate(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	find := func(entries ...map[string]any) []string {
		t.Helper()
		result := requestJSON(t, app, http.MethodPost, findPath, findRatesEmail, map[string]any{"entries": entries}, http.StatusOK)
		return findMatchIDs(t, result)
	}
	expect := func(what string, got []string, want ...string) {
		t.Helper()
		if want == nil {
			want = []string{}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v, want %v", what, got, want)
		}
	}

	// One role matches several sheets, newest code first. The inactive sheet
	// has the same rate and is left out.
	expect("one role", find(findRate(principalRole, 4321)), findGamma, findBeta, findAlpha)

	// Each role narrows the list. Gamma has no Administration Staff rate, so
	// a rate for that role leaves it out.
	expect("a role the sheet lacks", find(findRate(principalRole, 4321), findRate(adminRole, 111)), findBeta, findAlpha)
	expect("several roles", find(findRate(principalRole, 4321), findRate(adminRole, 111), findRate(designerRole, 222)), findAlpha)

	// An overtime rate counts only when entered.
	expect("overtime entered and equal", find(findRateOvertime(adminRole, 111, 144.30)), findAlpha)
	expect("overtime entered and different", find(findRateOvertime(adminRole, 111, 150)), findBeta)
	expect("overtime one cent off", find(findRateOvertime(adminRole, 111, 144.31)))
	expect("rate one dollar off", find(findRate(designerRole, 221)))

	// A draft has the standard Administration Staff rate too; only the active
	// standard sheet is returned.
	found := find(findRate(adminRole, 90))
	for _, id := range found {
		if id == draftFixtureSheet || id == incompleteDraftSheet {
			t.Fatalf("a request's proposed rates were returned: %v", found)
		}
	}
	expect("a draft's rate", found, standardRateSheet)

	// Only the inactive sheet has a Senior Designer at 225.
	expect("an inactive sheet", find(findRate(designerRole, 225)))
}

// The rates are compared in cents: 4321.00 and 5617.30 written with trailing
// zeros match the stored 4321 and 5617.3.
func TestFindByRatesComparesInCents(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	body := json.RawMessage(`{"entries":[{"role":"` + principalRole + `","rate":4321.00,"overtime_rate":5617.30}]}`)
	result := requestJSON(t, app, http.MethodPost, findPath, findRatesEmail, body, http.StatusOK)
	if got := findMatchIDs(t, result); !reflect.DeepEqual(got, []string{findGamma, findBeta, findAlpha}) {
		t.Fatalf("got %v", got)
	}
	match := result["matches"].([]any)[2].(map[string]any)
	for _, field := range []string{"id", "code", "name", "revision", "effective_date", "job_count"} {
		if _, ok := match[field]; !ok {
			t.Errorf("a match has no %s: %v", field, match)
		}
	}
	if match["code"] != "2026-11-00028" || match["name"] != "Find Rates Alpha" || match["job_count"] != float64(0) {
		t.Errorf("Alpha is %v", match)
	}
	if _, ok := match["differing_roles"]; ok {
		t.Errorf("a match has differing_roles: %v", match)
	}

	// The standard sheet has jobs.
	result = requestJSON(t, app, http.MethodPost, findPath, findRatesEmail, map[string]any{"entries": []any{findRate(principalRole, 250)}}, http.StatusOK)
	standard := result["matches"].([]any)[0].(map[string]any)
	if standard["id"] != standardRateSheet || standard["job_count"].(float64) < 1 {
		t.Errorf("the standard sheet is %v", standard)
	}
}

// Bad input is a field error under entries.N, and nothing is searched.
func TestFindByRatesRefusesBadInput(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	cases := []struct {
		what    string
		entries []any
		key     string
	}{
		{"no rate", []any{}, "entries"},
		{"no entries field", nil, "entries"},
		{"a negative rate", []any{findRate(principalRole, -250)}, "entries.0.rate"},
		{"a rate that is not whole", []any{findRate(principalRole, 250.5)}, "entries.0.rate"},
		{"a negative overtime rate", []any{findRateOvertime(principalRole, 250, -1)}, "entries.0.overtime_rate"},
		{"an overtime rate in tenths of a cent", []any{findRateOvertime(principalRole, 250, 325.001)}, "entries.0.overtime_rate"},
		{"an overtime rate without a rate", []any{map[string]any{"role": principalRole, "overtime_rate": 325}}, "entries.0.rate"},
		{"an unknown role", []any{findRate("notarole0000001", 250)}, "entries.0.role"},
		{"no role", []any{findRate("", 250)}, "entries.0.role"},
		{"a role twice", []any{findRate(principalRole, 250), findRate(principalRole, 260)}, "entries.1.role"},
	}
	for _, c := range cases {
		body := map[string]any{}
		if c.entries != nil {
			body["entries"] = c.entries
		}
		result := requestJSON(t, app, http.MethodPost, findPath, findRatesEmail, body, http.StatusBadRequest)
		data, _ := result["data"].(map[string]any)
		if _, ok := data[c.key]; !ok {
			t.Errorf("%s: no error under %s: %v", c.what, c.key, result)
		}
	}
}

// The route has the compare route's checks: a signed-in user with the time
// claim, while project requests are on. noclaims@example.com lacks the
// time claim.
func TestFindByRatesPermissions(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	body := map[string]any{"entries": []any{findRate(principalRole, 4321)}}
	mustStatus(t, performJSONRequest(t, app, http.MethodPost, findPath, "", body), http.StatusUnauthorized)
	refused := requestJSON(t, app, http.MethodPost, findPath, "noclaims@example.com", body, http.StatusForbidden)
	if !strings.Contains(refused["message"].(string), "time") {
		t.Errorf("refused for another reason: %v", refused)
	}
	requestJSON(t, app, http.MethodPost, findPath, draftFixtureCreator, body, http.StatusOK)

	off := newWorkflowTestApp(t, false)
	requestJSON(t, off, http.MethodPost, findPath, findRatesEmail, body, http.StatusForbidden)
}
