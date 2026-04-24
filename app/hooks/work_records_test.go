package hooks

import (
	"testing"
	"tybalt/errs"
	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
)

func requireHookError(t *testing.T, err error) *errs.HookError {
	t.Helper()

	if err == nil {
		t.Fatal("expected hook error")
	}

	normalizedErr, _ := NormalizeValidationError(err)
	hookErr, ok := normalizedErr.(*errs.HookError)
	if !ok {
		t.Fatalf("expected *errs.HookError, got %T", normalizedErr)
	}

	return hookErr
}

func newWorkRecordTestCollection() *core.Collection {
	collection := core.NewBaseCollection("work_records")
	collection.Fields.Add(&core.TextField{Name: "number"})
	collection.Fields.Add(&core.TextField{Name: "date"})
	collection.Fields.Add(&core.TextField{Name: "creator"})
	collection.Fields.Add(&core.TextField{Name: "job"})
	collection.Fields.Add(&core.TextField{Name: "type"})
	collection.Fields.Add(&core.TextField{Name: "location"})
	collection.Fields.Add(&core.TextField{Name: "work_description"})
	collection.Fields.Add(&core.TextField{Name: "parent_work_record"})
	collection.Fields.Add(&core.NumberField{Name: "lane_kms"})
	collection.Fields.Add(&core.TextField{Name: "attachment"})
	collection.Fields.Add(&core.TextField{Name: "attachment_hash"})
	return collection
}

func newWorkRecordSubjectTestCollection() *core.Collection {
	collection := core.NewBaseCollection("work_records_subjects")
	collection.Fields.Add(&core.TextField{Name: "work_record"})
	collection.Fields.Add(&core.TextField{Name: "uid"})
	collection.Fields.Add(&core.NumberField{Name: "hours_on_site"})
	collection.Fields.Add(&core.NumberField{Name: "hours_travel_time"})
	collection.Fields.Add(&core.NumberField{Name: "distance_travelled_km"})
	collection.Fields.Add(&core.TextField{Name: "vehicle_type"})
	collection.Fields.Add(&core.BoolField{Name: "is_passenger"})
	collection.Fields.Add(&core.TextField{Name: "company_vehicle_unit_number"})
	collection.Fields.Add(&core.BoolField{Name: "approved"})
	return collection
}

func newWorkRecordConsumableEntryTestCollection() *core.Collection {
	collection := core.NewBaseCollection("work_records_consumable_entries")
	collection.Fields.Add(&core.TextField{Name: "work_record"})
	collection.Fields.Add(&core.TextField{Name: "consumable"})
	collection.Fields.Add(&core.NumberField{Name: "quantity_number"})
	return collection
}

func TestValidateWorkRecordSubject_RejectsDuplicateDay(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	// This test needs to get past the max-workers check so it can assert the
	// duplicate-day validation branch. Keep the mutation local to this ephemeral
	// test app instead of adding a global seed row that changes list/report totals.
	if _, err := app.DB().NewQuery(`
		UPDATE work_records
		SET type = {:type}
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id":   "wr260400000010a",
		"type": "rqayp9yoz1qe6ui",
	}).Execute(); err != nil {
		t.Fatalf("failed to relax seeded work record max-workers limit: %v", err)
	}

	collection := newWorkRecordSubjectTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":                          "wrsubdupday01",
		"work_record":                 "wr260400000010a",
		"uid":                         "f2j5a8vk006baub",
		"hours_on_site":               8,
		"hours_travel_time":           1,
		"distance_travelled_km":       10,
		"vehicle_type":                "company",
		"company_vehicle_unit_number": "7",
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed original duplicate-day state: %v", err)
	}

	err := ValidateWorkRecordSubject(app, record)
	hookErr := requireHookError(t, err)
	if hookErr.Data["date"].Code != "duplicate_work_record_day" {
		t.Fatalf("expected duplicate day error, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordSubject_RejectsMaxSubjectsExceeded(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	collection := newWorkRecordSubjectTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":                          "wrsubmaxlimit01",
		"work_record":                 "wr260400000010a",
		"uid":                         "rzr98oadsp9qc11",
		"hours_on_site":               8,
		"hours_travel_time":           1,
		"distance_travelled_km":       10,
		"vehicle_type":                "company",
		"company_vehicle_unit_number": "7",
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed max-workers state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecordSubject(app, record))
	if hookErr.Data["uid"].Code != "max_subjects_exceeded" {
		t.Fatalf("expected max subjects error, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordSubject_RejectsApprovedSubjectEvenWhenApprovalIsCleared(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	record := buildRecordFromMap(newWorkRecordSubjectTestCollection(), map[string]any{
		"id":                          "wrsubapproved01",
		"work_record":                 "wr260400000010a",
		"uid":                         "f2j5a8vk006baub",
		"hours_on_site":               8,
		"hours_travel_time":           1,
		"distance_travelled_km":       10,
		"vehicle_type":                "company",
		"company_vehicle_unit_number": "7",
		"approved":                    true,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed approved subject state: %v", err)
	}

	record.Set("approved", false)
	record.Set("hours_on_site", 9)

	hookErr := requireHookError(t, ValidateWorkRecordSubject(app, record))
	if hookErr.Data["global"].Code != "approved_subject_locked" {
		t.Fatalf("expected approved subject lock, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordSubject_EnforcesVehicleRules(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	collection := newWorkRecordSubjectTestCollection()
	companyRecord := buildRecordFromMap(collection, map[string]any{
		"id":                    "wrvehiclerule1",
		"work_record":           "wr260400000010a",
		"uid":                   "rzr98oadsp9qc11",
		"hours_on_site":         8,
		"hours_travel_time":     1,
		"distance_travelled_km": 10,
		"vehicle_type":          "company",
	})
	if err := companyRecord.PostScan(); err != nil {
		t.Fatalf("failed to seed original company vehicle state: %v", err)
	}
	companyRecord.Set("company_vehicle_unit_number", "")
	hookErr := requireHookError(t, ValidateWorkRecordSubject(app, companyRecord))
	if hookErr.Data["company_vehicle_unit_number"].Code != "required_for_company_vehicle" {
		t.Fatalf("expected company vehicle unit-number error, got %#v", hookErr.Data)
	}

	personalRecord := buildRecordFromMap(collection, map[string]any{
		"id":                          "wrvehiclerule02",
		"work_record":                 "wr260400000010a",
		"uid":                         "rzr98oadsp9qc11",
		"hours_on_site":               8,
		"hours_travel_time":           1,
		"distance_travelled_km":       10,
		"vehicle_type":                "personal",
		"company_vehicle_unit_number": "7",
	})
	if err := personalRecord.PostScan(); err != nil {
		t.Fatalf("failed to seed original personal vehicle state: %v", err)
	}
	hookErr = requireHookError(t, ValidateWorkRecordSubject(app, personalRecord))
	if hookErr.Data["company_vehicle_unit_number"].Code != "must_be_blank_for_personal_vehicle" {
		t.Fatalf("expected personal vehicle unit-number error, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecord_RejectsLaneKmsWhenTypeDisallowsThem(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	record := buildRecordFromMap(newWorkRecordTestCollection(), map[string]any{
		"id":               "wrlanekmsdeny01",
		"number":           "W2604-9010",
		"date":             "2026-04-24",
		"creator":          "f2j5a8vk006baub",
		"job":              "cjf0kt0defhq480",
		"type":             "oask5u50nka0km8",
		"location":         "Lane kms disallowed",
		"work_description": "Lane kms disallowed",
		"lane_kms":         12.5,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed lane-kms state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecord(app, record))
	if hookErr.Data["lane_kms"].Code != "lane_kms_not_allowed" {
		t.Fatalf("expected lane kms rejection, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecord_RejectsNegativeLaneKms(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	record := buildRecordFromMap(newWorkRecordTestCollection(), map[string]any{
		"id":               "wrlanekmsmin001",
		"number":           "W2604-9011",
		"date":             "2026-04-24",
		"creator":          "f2j5a8vk006baub",
		"job":              "cjf0kt0defhq480",
		"type":             "lanekmstype0001",
		"location":         "Lane kms negative",
		"work_description": "Lane kms negative",
		"lane_kms":         -1.5,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed lane-kms state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecord(app, record))
	if hookErr.Data["lane_kms"].Code != "lane_kms_min" {
		t.Fatalf("expected lane kms min rejection, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecord_RejectsTypeChangeWhenConsumablesBecomeInvalid(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	record := buildRecordFromMap(newWorkRecordTestCollection(), map[string]any{
		"id":               "wr260400000010a",
		"number":           "W2604-0010",
		"date":             "2026-04-10",
		"creator":          "f2j5a8vk006baub",
		"job":              "cjf0kt0defhq480",
		"type":             "oask5u50nka0km8",
		"location":         "Seeded location",
		"work_description": "Seeded work description",
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed original work record state: %v", err)
	}
	record.Set("type", "rqayp9yoz1qe6ui")

	hookErr := requireHookError(t, ValidateWorkRecord(app, record))
	if hookErr.Data["type"].Code != "incompatible_consumable_entries" {
		t.Fatalf("expected incompatible consumable type-change error, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordConsumableEntry_EnforcesInputKindRules(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	collection := newWorkRecordConsumableEntryTestCollection()
	numberRecord := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentrytest001",
		"work_record":     "wr260400000003a",
		"consumable":      "wrconsmbl000001",
		"quantity_number": 0,
	})
	if err := numberRecord.PostScan(); err != nil {
		t.Fatalf("failed to seed original numeric consumable state: %v", err)
	}
	hookErr := requireHookError(t, ValidateWorkRecordConsumableEntry(app, numberRecord))
	if hookErr.Data["quantity_number"].Code != "quantity_required" {
		t.Fatalf("expected numeric quantity error, got %#v", hookErr.Data)
	}

	booleanRecord := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentrytest002",
		"work_record":     "wr260400000003a",
		"consumable":      "wrconsmbl000002",
		"quantity_number": 1,
	})
	if err := booleanRecord.PostScan(); err != nil {
		t.Fatalf("failed to seed original boolean consumable state: %v", err)
	}
	hookErr = requireHookError(t, ValidateWorkRecordConsumableEntry(app, booleanRecord))
	if hookErr.Data["quantity_number"].Code != "boolean_consumable_quantity_forbidden" {
		t.Fatalf("expected boolean quantity error, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordConsumableEntry_RejectsApprovedWorkRecordMutations(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	if _, err := app.DB().NewQuery(`
		UPDATE work_records_subjects
		SET approved = 1
		WHERE work_record = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000003a",
	}).Execute(); err != nil {
		t.Fatalf("failed to approve seeded work record: %v", err)
	}

	collection := newWorkRecordConsumableEntryTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentrytest003",
		"work_record":     "wr260400000003a",
		"consumable":      "wrconsmbl000001",
		"quantity_number": 1,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed original approved-lock consumable state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecordConsumableEntry(app, record))
	if hookErr.Data["global"].Code != "approved_work_record_locked" {
		t.Fatalf("expected approved work record lock error, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordConsumableEntry_RejectsConsumablesNotAllowedForType(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	collection := newWorkRecordConsumableEntryTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentrytest004",
		"work_record":     "wr260400000003a",
		"consumable":      "wrconsmbl000003",
		"quantity_number": 2,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed allowlist consumable state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecordConsumableEntry(app, record))
	if hookErr.Data["consumable"].Code != "consumable_not_allowed_for_type" {
		t.Fatalf("expected allowlist rejection, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordConsumableEntry_RejectsInactiveConsumables(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	if _, err := app.DB().NewQuery(`
		UPDATE work_records_consumables
		SET active = 0
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": "wrconsmbl000001",
	}).Execute(); err != nil {
		t.Fatalf("failed to deactivate seeded consumable: %v", err)
	}

	collection := newWorkRecordConsumableEntryTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentrytest005",
		"work_record":     "wr260400000003a",
		"consumable":      "wrconsmbl000001",
		"quantity_number": 2,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed inactive consumable state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecordConsumableEntry(app, record))
	if hookErr.Data["consumable"].Code != "inactive_consumable" {
		t.Fatalf("expected inactive consumable rejection, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordConsumableEntry_RejectsImmutableFieldChanges(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	collection := newWorkRecordConsumableEntryTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentry0000001",
		"work_record":     "wr260400000010a",
		"consumable":      "wrconsmbl000001",
		"quantity_number": 12,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed immutable original state: %v", err)
	}

	record.Set("consumable", "wrconsmbl000002")
	hookErr := requireHookError(t, ValidateWorkRecordConsumableEntry(app, record))
	if hookErr.Data["consumable"].Code != "immutable_field" {
		t.Fatalf("expected immutable consumable rejection, got %#v", hookErr.Data)
	}
}

func TestValidateWorkRecordConsumableEntryDelete_RejectsApprovedWorkRecordMutations(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	if _, err := app.DB().NewQuery(`
		UPDATE work_records_subjects
		SET approved = 1
		WHERE work_record = {:id}
	`).Bind(map[string]any{
		"id": "wr260400000010a",
	}).Execute(); err != nil {
		t.Fatalf("failed to approve seeded work record: %v", err)
	}

	collection := newWorkRecordConsumableEntryTestCollection()
	record := buildRecordFromMap(collection, map[string]any{
		"id":              "wrcentry0000001",
		"work_record":     "wr260400000010a",
		"consumable":      "wrconsmbl000001",
		"quantity_number": 12,
	})
	if err := record.PostScan(); err != nil {
		t.Fatalf("failed to seed delete-lock consumable state: %v", err)
	}

	hookErr := requireHookError(t, ValidateWorkRecordConsumableEntryDelete(app, record))
	if hookErr.Data["global"].Code != "approved_work_record_locked" {
		t.Fatalf("expected approved work record lock on delete, got %#v", hookErr.Data)
	}
}
