package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"tybalt/hooks"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// An overtime rate has at most two decimal places on every path that writes
// rate_sheet_entries: 1.5 and 1.25 pass, 1.255 is refused with the same
// message. The browser limits the input with step="0.01"; the server checks
// too.

var overtimeDecimalCases = []struct {
	overtime float64
	ok       bool
}{{1.5, true}, {1.25, true}, {1.255, false}}

func TestHasAtMostTwoDecimals(t *testing.T) {
	for value, want := range map[float64]bool{1: true, 1.5: true, 1.15: true, 1.25: true, 344.5: true, 0.29: true, 1.255: false, 1.001: false} {
		if got := hooks.HasAtMostTwoDecimals(value); got != want {
			t.Errorf("HasAtMostTwoDecimals(%v) = %v; want %v", value, got, want)
		}
	}
}

// POST /api/rate_sheets refuses it under entries.N.overtime_rate.
func TestCreateRouteLimitsOvertimeDecimals(t *testing.T) {
	app := testutils.SetupTestApp(t)
	roles, err := app.FindRecordsByFilter("rate_roles", "1=1", "name", 1, 0, nil)
	if err != nil || len(roles) == 0 {
		t.Fatalf("find a rate role: %v", err)
	}
	for i, c := range overtimeDecimalCases {
		status, body := createRateSheetThroughRoute(t, app, "author@soup.com", map[string]any{
			"name":     fmt.Sprintf("Overtime Decimals %d", i),
			"revision": 0,
			// A rate no other sheet has, so the duplicate check passes.
			"entries": []map[string]any{{"role": roles[0].Id, "rate": 7001 + i, "overtime_rate": c.overtime}},
		})
		if c.ok {
			if status != http.StatusCreated {
				t.Fatalf("overtime %v: status %d %v", c.overtime, status, body)
			}
			continue
		}
		data, _ := body["data"].(map[string]any)
		fieldErr, _ := data["entries.0.overtime_rate"].(map[string]any)
		if status != http.StatusBadRequest || fieldErr["code"] != hooks.OvertimeRateDecimalsCode ||
			!strings.Contains(fmt.Sprint(fieldErr["message"]), hooks.OvertimeRateDecimalsMessage) {
			t.Fatalf("overtime %v: status %d %v", c.overtime, status, body)
		}
	}
}

// PUT /api/rate_sheet_entries/{id} refuses it with the message. The entry is
// on the inactive fixture sheet rs_partial_001.
func TestEntryCorrectionLimitsOvertimeDecimals(t *testing.T) {
	app := testutils.SetupTestApp(t)
	const entry = "rsepart_0l1llckilwz1mqz"
	for _, c := range overtimeDecimalCases {
		status, raw := absorbRequest(t, app, http.MethodPut, "/api/rate_sheet_entries/"+entry, "author@soup.com", map[string]any{"overtime_rate": c.overtime})
		if c.ok != (status == http.StatusOK) {
			t.Fatalf("overtime %v: status %d %s", c.overtime, status, raw)
		}
		if !c.ok && (status != http.StatusBadRequest || !strings.Contains(strings.ToLower(raw), hooks.OvertimeRateDecimalsMessage)) {
			t.Fatalf("overtime %v: status %d %s", c.overtime, status, raw)
		}
	}
	record, err := app.FindRecordById("rate_sheet_entries", entry)
	if err != nil {
		t.Fatal(err)
	}
	if got := record.GetFloat("overtime_rate"); got != 1.25 {
		t.Fatalf("overtime_rate = %v; want the last accepted 1.25", got)
	}
}

// The collection API refuses it under overtime_rate. The entries go on the
// fixture sheet test_empty_sheet, which has none.
func TestCollectionCreateLimitsOvertimeDecimals(t *testing.T) {
	app := testutils.SetupTestApp(t)
	token, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	roles, err := app.FindRecordsByFilter("rate_roles", "1=1", "name", len(overtimeDecimalCases), 0, nil)
	if err != nil || len(roles) < len(overtimeDecimalCases) {
		t.Fatalf("find rate roles: %v", err)
	}
	for i, c := range overtimeDecimalCases {
		scenario := tests.ApiScenario{
			Name:            fmt.Sprintf("overtime %v", c.overtime),
			Method:          http.MethodPost,
			URL:             "/api/collections/rate_sheet_entries/records",
			Body:            strings.NewReader(fmt.Sprintf(`{"rate_sheet":"test_empty_sheet","role":"%s","rate":100,"overtime_rate":%v}`, roles[i].Id, c.overtime)),
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{fmt.Sprintf(`"overtime_rate":%v`, c.overtime)},
		}
		if !c.ok {
			scenario.ExpectedStatus = http.StatusBadRequest
			scenario.ExpectedContent = []string{`"overtime_rate":{"code":"` + hooks.OvertimeRateDecimalsCode + `"`, hooks.OvertimeRateDecimalsMessage}
		}
		runOnApp(t, app, scenario)
	}
}

// Proposed rates refuse it under proposed_rates.entries.N.overtime_rate,
// which the request form shows on the role's row.
func TestProposedRatesLimitOvertimeDecimals(t *testing.T) {
	app := newWorkflowTestApp(t, true)
	for _, c := range overtimeDecimalCases {
		entries := standardRates(t, app, 265)
		entries[2]["overtime_rate"] = c.overtime
		status := http.StatusOK
		if !c.ok {
			status = http.StatusBadRequest
		}
		result := proposeRates(t, app, draftFixtureRequest, draftFixtureCreator, proposedRates("Proposed Fixture Rates", "Client rate schedule", entries), status)
		if c.ok {
			continue
		}
		data, _ := result["data"].(map[string]any)
		fieldErr, _ := data[hooks.ProposedRatesField+".entries.2.overtime_rate"].(map[string]any)
		if fieldErr["code"] != hooks.OvertimeRateDecimalsCode || len(data) != 1 {
			t.Fatalf("overtime %v: %#v", c.overtime, result)
		}
	}
}
