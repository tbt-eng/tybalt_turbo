package main

import (
	"net/http"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tests"
)

// TestBundleTimesheet_InactiveManagerFails verifies that bundling a timesheet
// fails when the user's manager (who becomes the approver) is inactive.
func TestBundleTimesheet_InactiveManagerFails(t *testing.T) {
	noClaimsToken, err := testutils.GenerateRecordToken("users", "u_no_claims@example.com")
	if err != nil {
		t.Fatal(err)
	}
	// User has_inactive_mgr@test.com has a profile with manager = u_inactive
	recordToken, err := testutils.GenerateRecordToken("users", "has_inactive_mgr@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "bundle timesheet requires time claim",
			Method:         http.MethodPost,
			URL:            "/api/time_sheets/2024-01-13/bundle",
			Headers:        map[string]string{"Authorization": noClaimsToken},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"message":"Time claim required."`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:           "bundle timesheet fails when manager is inactive",
			Method:         http.MethodPost,
			URL:            "/api/time_sheets/2024-09-14/bundle", // A Saturday (week ending)
			Headers:        map[string]string{"Authorization": recordToken},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"code":"approver_not_active"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// TestBundleTimesheet_SelfManagerWithTaprAutoApproves verifies that when a user
// is their own manager and holds the tapr claim, bundling their timesheet sets
// the `approved` timestamp automatically (no separate approve call required).
//
// Fixture: self_apv_yes@test.com — uid u_self_apv_yes whose profile.manager
// points at themselves and who holds the tapr claim.
func TestBundleTimesheet_SelfManagerWithTaprAutoApproves(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "self_apv_yes@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:           "self-managed user with tapr auto-approves on bundle",
		Method:         http.MethodPost,
		URL:            "/api/time_sheets/2024-09-14/bundle",
		Headers:        map[string]string{"Authorization": recordToken},
		ExpectedStatus: 200,
		ExpectedContent: []string{
			`"message":"Time sheet processed successfully"`,
		},
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
			ts, err := app.FindFirstRecordByFilter(
				"time_sheets",
				"uid = {:uid} && week_ending = {:weekEnding}",
				dbx.Params{"uid": "u_self_apv_yes", "weekEnding": "2024-09-14"},
			)
			if err != nil {
				tb.Fatalf("failed to find newly bundled time sheet: %v", err)
			}
			if ts.GetString("approver") != "u_self_apv_yes" {
				tb.Fatalf("approver = %q, want %q", ts.GetString("approver"), "u_self_apv_yes")
			}
			if ts.GetDateTime("approved").IsZero() {
				tb.Fatalf("expected approved timestamp to be set on auto-approved time sheet")
			}
		},
	}

	scenario.Test(t)
}

// TestBundleTimesheet_UnlinkedWorkRecordsInWeekFail verifies that a user cannot
// bundle a week while any of their work record subject rows for that week has no
// linked time entry. The fixtures include missing links on both week boundaries
// to cover the Sunday-through-Saturday date window.
func TestBundleTimesheet_UnlinkedWorkRecordsInWeekFail(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "self_apv_yes@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:           "bundle timesheet fails with unlinked work record subject in week",
		Method:         http.MethodPost,
		URL:            "/api/time_sheets/2024-09-21/bundle",
		Headers:        map[string]string{"Authorization": recordToken},
		ExpectedStatus: http.StatusUnprocessableEntity,
		ExpectedContent: []string{
			`"code":"unlinked_work_records"`,
			`W2409-9001`,
			`W2409-9002`,
		},
		TestAppFactory: testutils.SetupTestApp,
	}

	scenario.Test(t)
}

func TestBundleTimesheet_LinkedWorkRecordInWeekSucceeds(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "self_apv_yes@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:           "bundle timesheet succeeds with linked work record subject in week",
		Method:         http.MethodPost,
		URL:            "/api/time_sheets/2024-09-28/bundle",
		Headers:        map[string]string{"Authorization": recordToken},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"message":"Time sheet processed successfully"`,
		},
		TestAppFactory: testutils.SetupTestApp,
		AfterTestFunc: func(tb testing.TB, app *tests.TestApp, _ *http.Response) {
			ts, err := app.FindFirstRecordByFilter(
				"time_sheets",
				"uid = {:uid} && week_ending = {:weekEnding}",
				dbx.Params{"uid": "u_self_apv_yes", "weekEnding": "2024-09-28"},
			)
			if err != nil {
				tb.Fatalf("failed to find newly bundled time sheet: %v", err)
			}
			if ts.GetString("approver") != "u_self_apv_yes" {
				tb.Fatalf("approver = %q, want %q", ts.GetString("approver"), "u_self_apv_yes")
			}
			subject, err := app.FindRecordById("work_records_subjects", "wrsubbundlelnk1")
			if err != nil {
				tb.Fatalf("failed to load linked work record subject: %v", err)
			}
			if !subject.GetBool("approved") {
				tb.Fatalf("expected linked work record subject to be approved after auto-approved bundle")
			}
		},
	}

	scenario.Test(t)
}

func TestBundleTimesheet_UnlinkedWorkRecordOutsideWeekSucceeds(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "self_apv_yes@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:           "bundle timesheet ignores unlinked work record outside week",
		Method:         http.MethodPost,
		URL:            "/api/time_sheets/2024-10-05/bundle",
		Headers:        map[string]string{"Authorization": recordToken},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"message":"Time sheet processed successfully"`,
		},
		TestAppFactory: testutils.SetupTestApp,
	}

	scenario.Test(t)
}

func TestBundleTimesheet_OtherUsersUnlinkedWorkRecordInWeekSucceeds(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "self_apv_yes@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:           "bundle timesheet ignores another user's unlinked work record in week",
		Method:         http.MethodPost,
		URL:            "/api/time_sheets/2024-10-12/bundle",
		Headers:        map[string]string{"Authorization": recordToken},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			`"message":"Time sheet processed successfully"`,
		},
		TestAppFactory: testutils.SetupTestApp,
	}

	scenario.Test(t)
}

// TestBundleTimesheet_SelfManagerWithoutTaprFails verifies that a user who is
// their own manager but does NOT hold the tapr claim cannot bundle (and thus
// cannot self-approve).
//
// Fixture: self_apv_no@test.com — uid u_self_apv_no whose profile.manager
// points at themselves but who does NOT hold the tapr claim.
func TestBundleTimesheet_SelfManagerWithoutTaprFails(t *testing.T) {
	recordToken, err := testutils.GenerateRecordToken("users", "self_apv_no@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenario := tests.ApiScenario{
		Name:           "self-managed user without tapr cannot bundle",
		Method:         http.MethodPost,
		URL:            "/api/time_sheets/2024-09-14/bundle",
		Headers:        map[string]string{"Authorization": recordToken},
		ExpectedStatus: 400,
		ExpectedContent: []string{
			`"code":"unqualified_approver"`,
		},
		TestAppFactory: testutils.SetupTestApp,
	}

	scenario.Test(t)
}
