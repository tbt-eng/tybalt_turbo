package routes

import (
	"fmt"
	"testing"
	"time"
	"tybalt/internal/testseed"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// helper to fetch a time_type id by code
func getTimeTypeId(t *testing.T, app core.App, code string) string {
	record, err := app.FindFirstRecordByFilter("time_types", "code={:code}", dbx.Params{"code": code})
	if err != nil {
		t.Fatalf("failed to load time_type %s: %v", code, err)
	}
	return record.Id
}

func insertWorkRecordNumberForValidateTimeEntriesTest(t *testing.T, app core.App, id string, number string) {
	t.Helper()

	if _, err := app.DB().NewQuery(`
		INSERT INTO work_records (
			id,
			creator,
			date,
			job,
			location,
			number,
			type,
			work_description
		) VALUES (
			{:id},
			'f2j5a8vk006baub',
			'2026-12-31',
			'mg0sp9iyjzo4zw9',
			'Test Location',
			{:number},
			'oask5u50nka0km8',
			'Duplicate label coverage'
		)
	`).Bind(dbx.Params{
		"id":     id,
		"number": number,
	}).Execute(); err != nil {
		t.Fatalf("failed to insert work record: %v", err)
	}
	if _, err := app.DB().NewQuery(`
		INSERT INTO work_records_subjects (
			id,
			work_record,
			uid,
			hours_on_site,
			hours_travel_time,
			distance_travelled_km,
			vehicle_type,
			company_vehicle_unit_number,
			approved
		) VALUES (
			'wrsduplabel001',
			{:id},
			'f2j5a8vk006baub',
			8,
			0,
			0,
			'company',
			'7',
			0
		)
	`).Bind(dbx.Params{
		"id": id,
	}).Execute(); err != nil {
		t.Fatalf("failed to insert work record subject: %v", err)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Scenario 1: User has negative Vacation balance but is NOT making a new OV
//
//	claim on this bundle → validation should succeed.
//
// ──────────────────────────────────────────────────────────────────────────────
func TestValidateTimeEntries_NegativeBalanceNoOVClaims_Passes(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	uid := "u_no_claims"
	weekEnding := "2024-01-13" // Saturday

	// create current R entry (40h) without OV/OP claims
	rId := getTimeTypeId(t, app, "R")

	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	entry := core.NewRecord(timeEntriesCollection)
	entry.Set("uid", uid)
	entry.Set("time_type", rId)
	entry.Set("hours", 40.0)
	entry.Set("date", "2024-01-08")
	entry.Set("week_ending", weekEnding)

	// payroll year end date (must be before weekEnding)
	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	// fetch admin_profile record
	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}

	if err := validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{entry}); err != nil {
		t.Fatalf("validation should have passed, got error: %v", err)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Scenario 2: User has negative Vacation balance AND is claiming additional OV
//
//	→ validator must reject with ov_claim_exceeds_balance.
//
// ──────────────────────────────────────────────────────────────────────────────
func TestValidateTimeEntries_ClaimOVWithNegativeBalance_Fails(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	uid := "u_with_claim"
	weekEnding := "2024-01-13" // Saturday

	// create current OV entry (1h) which should now trigger error
	ovId := getTimeTypeId(t, app, "OV")
	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	entry := core.NewRecord(timeEntriesCollection)
	entry.Set("uid", uid)
	entry.Set("time_type", ovId)
	entry.Set("hours", 1.0)
	entry.Set("date", "2024-01-08")
	entry.Set("week_ending", weekEnding)

	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}

	if err := validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{entry}); err == nil {
		t.Fatalf("expected validation to fail but it passed")
	} else {
		if codeErr, ok := err.(*CodeError); ok {
			if codeErr.Code != "ov_claim_exceeds_balance" {
				t.Fatalf("expected error code ov_claim_exceeds_balance, got %s", codeErr.Code)
			}
		} else {
			t.Fatalf("expected CodeError, got %v", err)
		}
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Scenario 3: User has negative PPTO balance but is NOT making a new OP claim
//
//	on this bundle → validation should succeed.
//
// ──────────────────────────────────────────────────────────────────────────────
func TestValidateTimeEntries_NegativeBalanceNoOPClaims_Passes(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	uid := "u_no_ppto_claim"
	weekEnding := "2024-01-13"

	rId := getTimeTypeId(t, app, "R")

	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	entry := core.NewRecord(timeEntriesCollection)
	entry.Set("uid", uid)
	entry.Set("time_type", rId)
	entry.Set("hours", 40.0)
	entry.Set("date", "2024-01-08")
	entry.Set("week_ending", weekEnding)

	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}

	if err := validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{entry}); err != nil {
		t.Fatalf("validation should have passed, got error: %v", err)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Scenario 4: User has negative PPTO balance AND is claiming additional OP
//
//	→ validator must reject with ppto_claim_exceeds_balance.
//
// ──────────────────────────────────────────────────────────────────────────────
func TestValidateTimeEntries_ClaimOPWithNegativeBalance_Fails(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	uid := "u_with_ppto_claim"
	weekEnding := "2024-01-13"

	opId := getTimeTypeId(t, app, "OP")

	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	entry := core.NewRecord(timeEntriesCollection)
	entry.Set("uid", uid)
	entry.Set("time_type", opId)
	entry.Set("hours", 1.0)
	entry.Set("date", "2024-01-08")
	entry.Set("week_ending", weekEnding)

	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}

	if err := validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{entry}); err == nil {
		t.Fatalf("expected validation to fail but it passed")
	} else {
		if codeErr, ok := err.(*CodeError); ok {
			if codeErr.Code != "ppto_claim_exceeds_balance" {
				t.Fatalf("expected error code ppto_claim_exceeds_balance, got %s", codeErr.Code)
			}
		} else {
			t.Fatalf("expected CodeError, got %v", err)
		}
	}
}

func TestValidateTimeEntries_UntrackedTimeOffSkipsMinimumHoursCheck(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	uid := "u_no_claims"
	weekEnding := "2024-01-13"
	rId := getTimeTypeId(t, app, "R")

	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	entry := core.NewRecord(timeEntriesCollection)
	entry.Set("uid", uid)
	entry.Set("time_type", rId)
	entry.Set("hours", 32.0)
	entry.Set("date", "2024-01-08")
	entry.Set("week_ending", weekEnding)

	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}
	adminProfile.Set("salary", true)
	adminProfile.Set("skip_min_time_check", "no")
	adminProfile.Set("untracked_time_off", true)
	adminProfile.Set("work_week_hours", 40)

	if err := validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{entry}); err != nil {
		t.Fatalf("validation should have passed for salaried staff with untracked time off, got error: %v", err)
	}
}

func TestValidateTimeEntries_UntrackedTimeOffRejectsRestrictedLeaveTypes(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	uid := "u_no_claims"
	weekEnding := "2024-01-13"
	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}
	adminProfile.Set("salary", true)
	adminProfile.Set("skip_min_time_check", "no")
	adminProfile.Set("untracked_time_off", true)

	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	for _, code := range []string{"OB", "OH", "OP", "OV"} {
		entry := core.NewRecord(timeEntriesCollection)
		entry.Set("uid", uid)
		entry.Set("time_type", getTimeTypeId(t, app, code))
		entry.Set("hours", 8.0)
		entry.Set("date", "2024-01-08")
		entry.Set("week_ending", weekEnding)

		err := validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{entry})
		if err == nil {
			t.Fatalf("expected validation to fail for %s when untracked time off is enabled", code)
		}

		codeErr, ok := err.(*CodeError)
		if !ok {
			t.Fatalf("expected CodeError for %s, got %v", code, err)
		}
		if codeErr.Code != "untracked_time_off_restricted" {
			t.Fatalf("expected error code untracked_time_off_restricted for %s, got %s", code, codeErr.Code)
		}
	}
}

func TestValidateTimeEntries_RejectsDuplicateVisibleWorkRecordAcrossLegacyAndLinkedEntries(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	insertWorkRecordNumberForValidateTimeEntriesTest(t, app, "wrduplabel00001", "F23-137")

	uid := "f2j5a8vk006baub"
	weekEnding := "2024-01-13"
	rID := getTimeTypeId(t, app, "R")
	payrollYearEndDate, _ := time.Parse("2006-01-02", "2024-01-06")

	adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
	if err != nil {
		t.Fatalf("failed to fetch admin_profile: %v", err)
	}

	timeEntriesCollection, _ := app.FindCollectionByNameOrId("time_entries")
	legacyEntry := core.NewRecord(timeEntriesCollection)
	legacyEntry.Set("uid", uid)
	legacyEntry.Set("time_type", rID)
	legacyEntry.Set("hours", 4.0)
	legacyEntry.Set("date", "2024-01-08")
	legacyEntry.Set("week_ending", weekEnding)
	legacyEntry.Set("work_record", "F23-137")

	linkedEntry := core.NewRecord(timeEntriesCollection)
	linkedEntry.Set("uid", uid)
	linkedEntry.Set("time_type", rID)
	linkedEntry.Set("hours", 4.0)
	linkedEntry.Set("date", "2024-01-09")
	linkedEntry.Set("week_ending", weekEnding)
	linkedEntry.Set("work_record_subject_id", "wrsduplabel001")

	err = validateTimeEntries(app, adminProfile, payrollYearEndDate, []*core.Record{legacyEntry, linkedEntry})
	if err == nil {
		t.Fatal("expected duplicate visible work record labels to fail validation")
	}

	codeErr, ok := err.(*CodeError)
	if !ok {
		t.Fatalf("expected CodeError, got %T", err)
	}
	if codeErr.Code != "multiple_work_records" {
		t.Fatalf("expected multiple_work_records, got %s", codeErr.Code)
	}
	if got := codeErr.Message; got != fmt.Sprintf("work record %s appears in multiple entries", "F23-137") {
		t.Fatalf("unexpected error message: %q", got)
	}
}
