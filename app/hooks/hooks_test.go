package hooks

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
	"tybalt/errs"
	"tybalt/internal/testseed"
	"tybalt/utilities"

	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// We need to instantiate a Collection object to be part of the Record object
// so everything works as expected
var timeEntriesCollection = core.NewBaseCollection("time_entries")

func buildRecordFromMap(collection *core.Collection, m map[string]any) *core.Record {
	if collection != nil {
		switch collection.Name {
		case "expenses", "purchase_orders":
			if _, exists := m["kind"]; !exists {
				m["kind"] = utilities.DefaultCapitalExpenditureKindID()
			}
		}
		if collection.Name == "expenses" {
			if _, hasOldAttachmentFixture := m["attachment"]; hasOldAttachmentFixture {
				if _, hasDocument := m["attachment_document"]; !hasDocument {
					m["attachment_document"] = "test_attachment_doc"
				}
				delete(m, "attachment")
			}
		}
	}
	record := core.NewRecord(collection)
	record.Load(m)
	return record
}

func upsertLinkedWorkRecordsConfigForHooksTest(t *testing.T, app *tests.TestApp, enabled bool) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatalf("failed to find app_config collection: %v", err)
	}

	record, err := app.FindFirstRecordByData("app_config", "key", "linked_work_records")
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", "linked_work_records")
	}

	record.Set("value", fmt.Sprintf(`{"enabled":%t}`, enabled))
	if err := app.Save(record); err != nil {
		t.Fatalf("failed to save linked_work_records config: %v", err)
	}
}

func insertWorkRecordNumberForHooksTest(t *testing.T, app *tests.TestApp, id string, number string) {
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
			'Collision coverage'
		)
	`).Bind(map[string]any{
		"id":     id,
		"number": number,
	}).Execute(); err != nil {
		t.Fatalf("failed to insert collision work record: %v", err)
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
			{:subjectId},
			{:id},
			'f2j5a8vk006baub',
			8,
			0,
			0,
			'company',
			'7',
			0
		)
	`).Bind(map[string]any{
		"id":        id,
		"subjectId": strings.TrimSuffix(id, "1") + "s",
	}).Execute(); err != nil {
		t.Fatalf("failed to insert collision work record subject: %v", err)
	}
}

func insertLegacyTimeEntryForHooksTest(t *testing.T, app *tests.TestApp, id string, workRecord string) {
	t.Helper()

	if _, err := app.DB().NewQuery(`
		INSERT INTO time_entries (
			id,
			uid,
			date,
			time_type,
			division,
			description,
			job,
			work_record,
			work_record_subject_id,
			hours,
			meals_hours,
			payout_request_amount,
			category,
			week_ending,
			role,
			tsid,
			branch
		) VALUES (
			{:id},
			'f2j5a8vk006baub',
			'2024-01-22',
			'sdyfl3q7j7ap849',
			'fy4i9poneukvq9u',
			'Legacy edit coverage',
			'mg0sp9iyjzo4zw9',
			{:workRecord},
			'',
			5,
			0,
			0,
			'',
			'2024-01-27',
			'tbgoiwwwfj8cvju',
			'',
			''
		)
	`).Bind(map[string]any{
		"id":         id,
		"workRecord": workRecord,
	}).Execute(); err != nil {
		t.Fatalf("failed to insert legacy time entry: %v", err)
	}
}

func linkSeededWorkRecordSubjectForHooksTest(t *testing.T, app *tests.TestApp, timeSheetID string, approved bool) {
	t.Helper()

	if _, err := app.DB().NewQuery(`
		UPDATE time_entries
		SET tsid = {:timeSheetId},
		    work_record_subject_id = 'wrsub260400003',
		    work_record = ''
		WHERE id = 'wrkrcdentry0001'
	`).Bind(map[string]any{
		"timeSheetId": timeSheetID,
	}).Execute(); err != nil {
		t.Fatalf("failed to link seeded time entry to work record subject: %v", err)
	}

	if _, err := app.DB().NewQuery(`
		UPDATE work_records_subjects
		SET approved = {:approved}
		WHERE id = 'wrsub260400003'
	`).Bind(map[string]any{
		"approved": approved,
	}).Execute(); err != nil {
		t.Fatalf("failed to seed work record subject approval: %v", err)
	}
}

func workRecordSubjectApprovedForHooksTest(t *testing.T, app *tests.TestApp) bool {
	t.Helper()

	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT COALESCE(approved, 0) AS approved
		FROM work_records_subjects
		WHERE id = 'wrsub260400003'
	`).One(&row); err != nil {
		t.Fatalf("failed to load work record subject approval: %v", err)
	}
	return row.Approved
}

func extractArbitraryKey(err error) (string, error) {
	errMap := err.(validation.Errors)
	keys := make([]string, len(errMap))
	i := 0
	for k := range errMap {
		keys[i] = k
		i++
	}

	// if the keys array length is not 1, return an error
	if len(keys) != 1 {
		return "", fmt.Errorf("multiple error keys: %v", strings.Join(keys, ", "))
	}

	// return the first key
	return keys[0], nil
}

func TestTimeSheetHooks_SyncWorkRecordSubjectOnDirectApprovalChange(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	AddHooks(app)

	linkSeededWorkRecordSubjectForHooksTest(t, app, "av32qwch9xrcb5n", false)

	timeSheet, err := app.FindRecordById("time_sheets", "av32qwch9xrcb5n")
	if err != nil {
		t.Fatalf("failed to load timesheet: %v", err)
	}
	timeSheet.Set("approved", "2026-04-27 12:00:00.000Z")
	if err := app.Save(timeSheet); err != nil {
		t.Fatalf("failed to approve timesheet through general save hook: %v", err)
	}

	if !workRecordSubjectApprovedForHooksTest(t, app) {
		t.Fatalf("expected direct timesheet approval to approve linked work record subject")
	}

	timeSheet, err = app.FindRecordById("time_sheets", "av32qwch9xrcb5n")
	if err != nil {
		t.Fatalf("failed to reload timesheet: %v", err)
	}
	timeSheet.Set("approved", "")
	if err := app.Save(timeSheet); err != nil {
		t.Fatalf("failed to unapprove timesheet through general save hook: %v", err)
	}

	if workRecordSubjectApprovedForHooksTest(t, app) {
		t.Fatalf("expected direct timesheet unapproval to clear linked work record subject")
	}
}

func TestTimeSheetHooks_SyncWorkRecordSubjectOnDirectRejection(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	AddHooks(app)

	linkSeededWorkRecordSubjectForHooksTest(t, app, "pts000000000001", true)

	timeSheet, err := app.FindRecordById("time_sheets", "pts000000000001")
	if err != nil {
		t.Fatalf("failed to load timesheet: %v", err)
	}
	timeSheet.Set("rejected", "2026-04-27 12:00:00.000Z")
	if err := app.Save(timeSheet); err != nil {
		t.Fatalf("failed to reject timesheet through general save hook: %v", err)
	}

	if workRecordSubjectApprovedForHooksTest(t, app) {
		t.Fatalf("expected direct timesheet rejection to clear linked work record subject")
	}
}

func TestTimeEntryHooks_SyncWorkRecordSubjectOnDirectUnbundle(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	AddHooks(app)

	linkSeededWorkRecordSubjectForHooksTest(t, app, "pts000000000001", true)

	entry, err := app.FindRecordById("time_entries", "wrkrcdentry0001")
	if err != nil {
		t.Fatalf("failed to load time entry: %v", err)
	}
	entry.Set("tsid", "")
	if err := app.Save(entry); err != nil {
		t.Fatalf("failed to unbundle time entry through general save hook: %v", err)
	}

	if workRecordSubjectApprovedForHooksTest(t, app) {
		t.Fatalf("expected direct time-entry unbundle to clear linked work record subject")
	}
}

func TestTimeSheetHooks_SyncWorkRecordSubjectOnDirectDelete(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)
	AddHooks(app)

	linkSeededWorkRecordSubjectForHooksTest(t, app, "pts000000000001", true)

	timeSheet, err := app.FindRecordById("time_sheets", "pts000000000001")
	if err != nil {
		t.Fatalf("failed to load timesheet: %v", err)
	}
	if err := app.Delete(timeSheet); err != nil {
		t.Fatalf("failed to delete timesheet through general delete hook: %v", err)
	}

	if workRecordSubjectApprovedForHooksTest(t, app) {
		t.Fatalf("expected direct timesheet delete to clear linked work record subject")
	}
}

func TestValidateTimeEntry(t *testing.T) {
	requiredFieldsByTimeTypeCode := map[string][]string{
		"OB":  {"date", "time_type", "hours"},
		"OH":  {"date", "time_type", "hours"},
		"OP":  {"date", "time_type", "hours"},
		"OR":  {"date", "time_type"},
		"OS":  {"date", "time_type", "hours"},
		"OTO": {"date", "time_type", "payout_request_amount"},
		"OV":  {"date", "time_type", "hours"},
		"OW":  {"date", "time_type"},
		"R":   {"date", "time_type", "division", "hours", "description"},
		"RB":  {"date", "time_type", "hours"},
		"RT":  {"date", "time_type", "division", "hours", "description"},
	}

	// Test cases
	cases := map[string]struct {
		timeTypeCode string
		valid        bool
		field        string
		record       *core.Record
	}{
		// Regular Time (R)
		// NB: the time_type value is arbitrary and is required but otherwise not
		// validated in this test. It is an id of a relation and is validated by the
		// pocketbase relation itself.
		"valid no job":              {timeTypeCode: "R", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"valid with job":            {timeTypeCode: "R", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "fy4i9poneukvq9u", "meals_hours": 0.5, "hours": 6.5, "job": "mg0sp9iyjzo4zw9", "description": "This is more than 5 chars", "work_record": "F23-137", "payout_request_amount": 0.0, "role": "tbgoiwwwfj8cvju"})},
		"valid leap year":           {timeTypeCode: "R", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-02-29", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"invalid date":              {timeTypeCode: "R", valid: false, field: "date", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-32", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"missing date":              {timeTypeCode: "R", valid: false, field: "date", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"invalid date string":       {timeTypeCode: "R", valid: false, field: "date", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "20240122", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"invalid leap year":         {timeTypeCode: "R", valid: false, field: "date", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2023-02-29", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"negative hours":            {timeTypeCode: "R", valid: false, field: "hours", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 1.0, "hours": -0.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"hours not multiple of 0.5": {timeTypeCode: "R", valid: false, field: "hours", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 1.0, "hours": 1.4, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"work record without job":   {timeTypeCode: "R", valid: false, field: "work_record", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "F23-137", "payout_request_amount": 0.0})},
		"missing division":          {timeTypeCode: "R", valid: false, field: "division", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "", "meals_hours": 1.0, "hours": 1.5, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"too many hours":            {timeTypeCode: "R", valid: false, field: "global", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 1.0, "hours": 18.0, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"no hours":                  {timeTypeCode: "R", valid: false, field: "hours", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 0.0, "hours": 0.0, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"bad work_record value":     {timeTypeCode: "R", valid: false, field: "work_record", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "fy4i9poneukvq9u", "meals_hours": 0.5, "hours": 5.0, "job": "mg0sp9iyjzo4zw9", "description": "This is more than 5 chars", "work_record": "F23-137-", "payout_request_amount": 0.0, "role": "tbgoiwwwfj8cvju"})},
		"no job, missing hours":     {timeTypeCode: "R", valid: false, field: "hours", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "DE", "meals_hours": 0.0, "hours": 0.0, "job": "", "description": "This is more than 5 chars", "work_record": "", "payout_request_amount": 0.0})},
		"description too short":     {timeTypeCode: "R", valid: false, field: "description", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "fy4i9poneukvq9u", "meals_hours": 0.5, "hours": 5.0, "job": "mg0sp9iyjzo4zw9", "description": "tiny", "work_record": "", "payout_request_amount": 0, "role": "tbgoiwwwfj8cvju"})},
		"closed job":                {timeTypeCode: "R", valid: false, field: "job", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "division": "fy4i9poneukvq9u", "hours": 1.0, "job": "zke3cs3yipplwtu", "description": "This is more than 5 chars", "role": "tbgoiwwwfj8cvju"})},

		// Vacation, Sick, Personal, Holiday, Bereavement
		"valid vacation":                {timeTypeCode: "OV", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "hours": 8, "description": "This is more than 5 chars"})},
		"valid sick":                    {timeTypeCode: "OS", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "hours": 8, "description": "This is more than 5 chars"})},
		"valid vacation no description": {timeTypeCode: "OV", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "hours": 8})},
		"vacation missing hours":        {timeTypeCode: "OV", valid: false, field: "hours", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "description": "This is more than 5 chars"})},
		"vac description too short":     {timeTypeCode: "OV", valid: false, field: "description", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "hours": 8, "description": "tiny"})},
		"ppto description too short":    {timeTypeCode: "OP", valid: false, field: "description", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "hours": 8, "description": "tiny"})},

		// Bank Time, Overtime Off Request
		"valid bank time":            {timeTypeCode: "RB", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "hours": 8})},
		"valid overtime off request": {timeTypeCode: "OTO", valid: true, record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "payout_request_amount": 100.0})},
		"bank time missing date":     {timeTypeCode: "RB", valid: false, field: "date", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "hours": 8})},
		"OTO no amount":              {timeTypeCode: "OTO", valid: false, field: "payout_request_amount", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22"})},
		"negative OTO amount":        {timeTypeCode: "OTO", valid: false, field: "payout_request_amount", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "payout_request_amount": -100.0})},
		"fractional OTO amount":      {timeTypeCode: "OTO", valid: false, field: "payout_request_amount", record: buildRecordFromMap(timeEntriesCollection, map[string]any{"time_type": "dummy", "date": "2024-01-22", "payout_request_amount": 132.001})},
	}

	// Initialize a PocketBase TestApp so validateTimeEntry can perform lookups
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	// Run tests
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			got := validateTimeEntry(app, tt.record, requiredFieldsByTimeTypeCode[tt.timeTypeCode])
			if got != nil {
				// we extract an arbitrary key from the error map. Key order doesn't
				// exist in Go maps but we expect only one error at a time so arbitrary
				// works here.
				field, keysError := extractArbitraryKey(got)
				if keysError != nil {
					t.Errorf("error extracting key from error map: %v", keysError)
				}
				if tt.valid {
					t.Errorf("failed validation (%v) but expected valid", got)
				}
				if !tt.valid {
					// If the field doesn't match expected field, the test has failed
					if field != string(tt.field) {
						t.Errorf("expected field: %s, got: %s", string(tt.field), field)
					}
				}
			}
			if got == nil && !tt.valid {
				t.Errorf("passed validation but expected invalid")
			}
		})
	}
}

func TestValidateTimeAmendment_AllowsSignedHalfHourHours(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	timeAmendmentsCollection, err := app.FindCollectionByNameOrId("time_amendments")
	if err != nil {
		t.Fatalf("failed to load time_amendments collection: %v", err)
	}

	tests := []struct {
		name    string
		hours   float64
		wantErr bool
	}{
		{name: "negative half hour", hours: -0.5},
		{name: "negative boundary", hours: -18},
		{name: "positive half hour", hours: 0.5},
		{name: "negative quarter hour", hours: -0.25, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := buildRecordFromMap(timeAmendmentsCollection, map[string]any{
				"time_type":             "dummy",
				"date":                  "2024-01-22",
				"division":              "DE",
				"hours":                 tt.hours,
				"description":           "Amendment correction",
				"payout_request_amount": 0.0,
			})

			got := validateTimeAmendment(app, record, []string{"date", "time_type", "division", "hours", "description"})
			if (got != nil) != tt.wantErr {
				t.Fatalf("validateTimeAmendment() error = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestNormalizeValidationError_Passthroughs(t *testing.T) {
	t.Run("returns existing hook error unchanged", func(t *testing.T) {
		original := &errs.HookError{
			Status:  http.StatusConflict,
			Message: "already normalized",
			Data: map[string]errs.CodeError{
				"global": {
					Code:    "already_normalized",
					Message: "already normalized",
				},
			},
		}

		got, handled := NormalizeValidationError(original)
		if !handled {
			t.Fatal("expected hook error to be handled")
		}
		if got != original {
			t.Fatalf("expected same hook error pointer, got %#v", got)
		}
	})

	t.Run("returns non validation error unchanged", func(t *testing.T) {
		original := errors.New("plain error")

		got, handled := NormalizeValidationError(original)
		if handled {
			t.Fatal("expected plain error to remain unhandled")
		}
		if got != original {
			t.Fatalf("expected same plain error, got %#v", got)
		}
	})
}

func TestValidateTimeEntry_AggregatesLinkedWorkRecordMismatchErrors(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	if err := app.RunAllMigrations(); err != nil {
		t.Fatalf("failed to apply migrations in test app: %v", err)
	}

	record := buildRecordFromMap(timeEntriesCollection, map[string]any{
		"time_type":              "dummy",
		"uid":                    "rzr98oadsp9qc11",
		"date":                   "2026-04-09",
		"division":               "fy4i9poneukvq9u",
		"hours":                  9.0,
		"job":                    "mg0sp9iyjzo4zw9",
		"description":            "Linked work record mismatch coverage",
		"role":                   "tbgoiwwwfj8cvju",
		"work_record_subject_id": "wrsub260400010",
	})

	err := validateTimeEntry(app, record, nil)
	if err == nil {
		t.Fatal("expected linked work record validation to fail")
	}

	errs, ok := err.(validation.Errors)
	if !ok {
		t.Fatalf("expected validation.Errors, got %T", err)
	}

	workRecordErr := errs["work_record_subject_id"]
	if workRecordErr == nil {
		t.Fatalf("expected work_record_subject_id error, got %#v", errs)
	}

	message := workRecordErr.Error()
	for _, expected := range []string{
		"Selected work record subject belongs to a different user",
		"Selected work record is for a different job",
		"Selected work record is for a different date",
	} {
		if !strings.Contains(message, expected) {
			t.Fatalf("expected aggregated work record error to contain %q, got %q", expected, message)
		}
	}
}

func TestValidateTimeEntry_HybridModeAllowsEditingLegacyTextOnlyEntry(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	insertLegacyTimeEntryForHooksTest(t, app, "legacyhybrid001", "F34-142")

	record, err := app.FindRecordById("time_entries", "legacyhybrid001")
	if err != nil {
		t.Fatalf("failed to load legacy time entry: %v", err)
	}

	record.Set("description", "Updated legacy work record entry")

	if err := validateTimeEntry(app, record, []string{"date", "time_type", "division", "hours", "description"}); err != nil {
		t.Fatalf("expected hybrid mode to allow editing legacy text-only entry, got: %v", err)
	}
}

func TestValidateTimeEntry_LinkedOnlyModeForcesLegacyUpgrade(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	upsertLinkedWorkRecordsConfigForHooksTest(t, app, true)
	insertLegacyTimeEntryForHooksTest(t, app, "legacylinked001", "F34-142")

	record, err := app.FindRecordById("time_entries", "legacylinked001")
	if err != nil {
		t.Fatalf("failed to load legacy time entry: %v", err)
	}

	record.Set("description", "Updated legacy work record entry")

	err = validateTimeEntry(app, record, []string{"date", "time_type", "division", "hours", "description"})
	if err == nil {
		t.Fatal("expected linked-only mode to require upgrading legacy entry")
	}

	errs, ok := err.(validation.Errors)
	if !ok {
		t.Fatalf("expected validation.Errors, got %T", err)
	}
	if errs["work_record_subject_id"] == nil {
		t.Fatalf("expected work_record_subject_id error, got %#v", errs)
	}
}

func TestValidateTimeEntry_LinkedOnlyModeRejectsNewTextOnlyWorkRecord(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	upsertLinkedWorkRecordsConfigForHooksTest(t, app, true)

	record := buildRecordFromMap(timeEntriesCollection, map[string]any{
		"time_type":             "dummy",
		"uid":                   "f2j5a8vk006baub",
		"date":                  "2024-01-22",
		"division":              "fy4i9poneukvq9u",
		"hours":                 5.0,
		"job":                   "mg0sp9iyjzo4zw9",
		"description":           "Linked-only rejection coverage",
		"role":                  "tbgoiwwwfj8cvju",
		"work_record":           "F23-137",
		"meals_hours":           0.0,
		"payout_request_amount": 0.0,
	})

	err := validateTimeEntry(app, record, []string{"date", "time_type", "division", "hours", "description"})
	if err == nil {
		t.Fatal("expected linked-only mode to reject new text-only work record")
	}

	errs, ok := err.(validation.Errors)
	if !ok {
		t.Fatalf("expected validation.Errors, got %T", err)
	}
	if errs["work_record"] == nil {
		t.Fatalf("expected work_record error, got %#v", errs)
	}
}

func TestValidateTimeEntry_RejectsLegacyTextThatMatchesLinkedWorkRecordNumber(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	insertWorkRecordNumberForHooksTest(t, app, "wrlegacyclash001", "F23-137")

	record := buildRecordFromMap(timeEntriesCollection, map[string]any{
		"time_type":             "dummy",
		"uid":                   "f2j5a8vk006baub",
		"date":                  "2024-01-22",
		"division":              "fy4i9poneukvq9u",
		"hours":                 5.0,
		"job":                   "mg0sp9iyjzo4zw9",
		"description":           "Legacy collision coverage",
		"role":                  "tbgoiwwwfj8cvju",
		"work_record":           "F23-137",
		"meals_hours":           0.0,
		"payout_request_amount": 0.0,
	})

	err := validateTimeEntry(app, record, []string{"date", "time_type", "division", "hours", "description"})
	if err == nil {
		t.Fatal("expected legacy collision with linked work record number to fail")
	}

	errs, ok := err.(validation.Errors)
	if !ok {
		t.Fatalf("expected validation.Errors, got %T", err)
	}
	if errs["work_record"] == nil {
		t.Fatalf("expected work_record error, got %#v", errs)
	}
}

func TestValidateTimeEntry_AllowsInProgressProposal(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	record := buildRecordFromMap(timeEntriesCollection, map[string]any{
		"time_type":             "dummy",
		"date":                  "2024-01-22",
		"division":              "fy4i9poneukvq9u",
		"hours":                 1.0,
		"job":                   "test_prop_inprog",
		"description":           "Proposal review and edits",
		"role":                  "tbgoiwwwfj8cvju",
		"work_record":           "F23-137",
		"payout_request_amount": 0.0,
	})

	if err := validateTimeEntry(app, record, []string{"date", "time_type", "division", "hours", "description"}); err != nil {
		t.Fatalf("expected in-progress proposal to accept time entry, got: %v", err)
	}
}

func TestValidateTimeAmendment_AllowsInProgressProposal(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	timeAmendmentsCollection, err := app.FindCollectionByNameOrId("time_amendments")
	if err != nil {
		t.Fatalf("failed to load time_amendments collection: %v", err)
	}

	record := buildRecordFromMap(timeAmendmentsCollection, map[string]any{
		"time_type":             "dummy",
		"date":                  "2024-01-22",
		"division":              "fy4i9poneukvq9u",
		"hours":                 1.0,
		"job":                   "test_prop_inprog",
		"description":           "Proposal review and edits",
		"work_record":           "F23-137",
		"payout_request_amount": 0.0,
	})

	if err := validateTimeAmendment(app, record, []string{"date", "time_type", "division", "hours", "description"}); err != nil {
		t.Fatalf("expected in-progress proposal to accept time amendment, got: %v", err)
	}
}

func TestValidateTimeAmendment_IgnoresProposalStatus(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	timeAmendmentsCollection, err := app.FindCollectionByNameOrId("time_amendments")
	if err != nil {
		t.Fatalf("failed to load time_amendments collection: %v", err)
	}

	record := buildRecordFromMap(timeAmendmentsCollection, map[string]any{
		"time_type":             "dummy",
		"date":                  "2024-01-22",
		"division":              "fy4i9poneukvq9u",
		"hours":                 1.0,
		"job":                   "awpropstale0001",
		"description":           "Proposal correction after award",
		"work_record":           "F23-137",
		"payout_request_amount": 0.0,
	})

	if err := validateTimeAmendment(app, record, []string{"date", "time_type", "division", "hours", "description"}); err != nil {
		t.Fatalf("expected time amendment to ignore awarded proposal status, got: %v", err)
	}
}

func TestValidateJobAllowsTimeTrackingAt_AllowsRecentlyAwardedProposal(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	proposal, err := app.FindRecordById("jobs", "awproprecent001")
	if err != nil {
		t.Fatalf("failed to load proposal: %v", err)
	}

	now := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)

	if err := validateJobAllowsTimeTrackingAt(app, proposal, now); err != nil {
		t.Fatalf("expected recently awarded proposal to allow time tracking, got: %v", err)
	}
}

func TestValidateJobAllowsTimeTrackingAt_RejectsStaleAwardedProposal(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	proposal, err := app.FindRecordById("jobs", "awpropstale0001")
	if err != nil {
		t.Fatalf("failed to load proposal: %v", err)
	}

	now := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)

	err = validateJobAllowsTimeTrackingAt(app, proposal, now)
	if err == nil {
		t.Fatal("expected awarded proposal older than 30 days to be rejected")
	}

	validationErr, ok := err.(validation.Error)
	if !ok {
		t.Fatalf("expected validation.Error, got %T", err)
	}
	if validationErr.Code() != "proposal_awarded_more_than_30_days_ago" {
		t.Fatalf("expected proposal_awarded_more_than_30_days_ago, got %s", validationErr.Code())
	}
}

func TestValidateJobAllowsTimeTrackingAt_RejectsAwardedProposalWithoutReferencingProject(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	proposal, err := app.FindRecordById("jobs", "awpropnoref0001")
	if err != nil {
		t.Fatalf("failed to load proposal: %v", err)
	}

	now := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)

	err = validateJobAllowsTimeTrackingAt(app, proposal, now)
	if err == nil {
		t.Fatal("expected awarded proposal without referencing project to be rejected")
	}

	validationErr, ok := err.(validation.Error)
	if !ok {
		t.Fatalf("expected validation.Error, got %T", err)
	}
	if validationErr.Code() != "proposal_awarded_without_referencing_project" {
		t.Fatalf("expected proposal_awarded_without_referencing_project, got %s", validationErr.Code())
	}
}

func TestValidateJobAllowsTimeTrackingAt_RejectsFutureProjectAwardDate(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	proposal, err := app.FindRecordById("jobs", "awpropfuture001")
	if err != nil {
		t.Fatalf("failed to load proposal: %v", err)
	}

	now := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)

	err = validateJobAllowsTimeTrackingAt(app, proposal, now)
	if err == nil {
		t.Fatal("expected future project award date to be rejected")
	}

	validationErr, ok := err.(validation.Error)
	if !ok {
		t.Fatalf("expected validation.Error, got %T", err)
	}
	if validationErr.Code() != "invalid_project_award_date" {
		t.Fatalf("expected invalid_project_award_date, got %s", validationErr.Code())
	}
}
