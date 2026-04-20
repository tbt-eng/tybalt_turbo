// Package hooks owns PocketBase record hooks that enforce backend-side business
// rules during record persistence.
//
// This file exists to keep work record numbering server-assigned and consistent
// regardless of which UI or internal caller creates the record. Work records
// are intended to use a generated W-number sequence, and the transition in this
// branch makes that number part of the user-visible workflow because linked time
// entries and compatibility reads resolve through the first-class work_records
// table.
//
// The hook deliberately does one small job: if a work record arrives without a
// number, it reserves the next valid W-number and writes it before save. If the
// number is already present, the hook preserves it so controlled import/backfill
// flows remain possible. Keeping the logic here means the numbering policy is
// enforced even when callers bypass the normal UI and talk directly to the
// backend.
package hooks

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"tybalt/numbering"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

// GenerateWorkRecordNumber reserves the next work-record W-number using the
// shared numbering package and the work record's saved date.
func GenerateWorkRecordNumber(app core.App, workDate string) (string, error) {
	dateValue := time.Now()
	if strings.TrimSpace(workDate) != "" {
		parsed, err := time.Parse("2006-01-02", workDate)
		if err != nil {
			return "", validationFieldError(
				"date",
				"invalid_work_record_date",
				fmt.Sprintf("invalid work record date: %v", err),
			)
		}
		dateValue = parsed
	}

	return numbering.ReserveWorkRecordNumber(app, dateValue)
}

func validateWorkRecordRequiredFields(record *core.Record) error {
	requiredStrings := []struct {
		Name    string
		Message string
	}{
		{Name: "date", Message: "Work record date is required"},
		{Name: "uid", Message: "Work record owner is required"},
		{Name: "creator", Message: "Work record creator is required"},
		{Name: "job", Message: "Job is required"},
		{Name: "type", Message: "Work record type is required"},
		{Name: "location", Message: "Location is required"},
		{Name: "work_description", Message: "Work description is required"},
		{Name: "vehicle_type", Message: "Vehicle type is required"},
	}

	fieldErrors := validation.Errors{}
	for _, field := range requiredStrings {
		if strings.TrimSpace(record.GetString(field.Name)) == "" {
			fieldErrors[field.Name] = validation.NewError("required", field.Message)
		}
	}

	return fieldErrors.Filter()
}

func validateWorkRecordCoreFieldFormats(record *core.Record) error {
	workDate := strings.TrimSpace(record.GetString("date"))
	if workDate == "" {
		return nil
	}
	if _, err := time.Parse("2006-01-02", workDate); err != nil {
		return validationFieldError(
			"date",
			"invalid_work_record_date",
			fmt.Sprintf("invalid work record date: %v", err),
		)
	}

	switch strings.TrimSpace(record.GetString("vehicle_type")) {
	case "company", "personal":
		return nil
	default:
		return validationFieldError(
			"vehicle_type",
			"invalid_vehicle_type",
			"Vehicle type must be company or personal",
		)
	}
}

func validateWorkRecordVehicleFields(record *core.Record) error {
	vehicleType := strings.TrimSpace(record.GetString("vehicle_type"))
	unitNumber := strings.TrimSpace(record.GetString("company_vehicle_unit_number"))

	if vehicleType == "company" && unitNumber == "" {
		return validationFieldError(
			"company_vehicle_unit_number",
			"required_for_company_vehicle",
			"Company vehicle unit number is required when vehicle type is company",
		)
	}

	if vehicleType == "personal" && unitNumber != "" {
		return validationFieldError(
			"company_vehicle_unit_number",
			"must_be_blank_for_personal_vehicle",
			"Company vehicle unit number must be empty when vehicle type is personal",
		)
	}

	return nil
}

func validateWorkRecordUniqueDay(app core.App, record *core.Record) error {
	var duplicate struct {
		ID string `db:"id"`
	}
	err := app.DB().NewQuery(`
		SELECT id
		FROM work_records
		WHERE uid = {:uid}
		  AND job = {:job}
		  AND date = {:date}
		  AND id != {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"uid":  record.GetString("uid"),
		"job":  record.GetString("job"),
		"date": record.GetString("date"),
		"id":   record.Id,
	}).One(&duplicate)
	if err == nil && duplicate.ID != "" {
		return validationFieldError(
			"date",
			"duplicate_work_record_day",
			"A work record already exists for this worker, job, and date",
		)
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func validateWorkRecordTypeChange(app core.App, record *core.Record) error {
	original := record.Original()
	if original == nil || original.Id == "" || original.GetString("type") == record.GetString("type") || record.Id == "" {
		return nil
	}

	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_consumable_entries wrce
		LEFT JOIN work_records_type_consumables wrtc
		  ON wrtc.type = {:type}
		 AND wrtc.consumable = wrce.consumable
		 AND COALESCE(wrtc.active, 0) = 1
		WHERE wrce.work_record = {:workRecordId}
		  AND wrtc.id IS NULL
	`).Bind(dbx.Params{
		"type":         record.GetString("type"),
		"workRecordId": record.Id,
	}).One(&row); err != nil {
		return err
	}

	if row.Count > 0 {
		return validationFieldError(
			"type",
			"incompatible_consumable_entries",
			"Remove or adjust incompatible consumable entries before changing the work record type",
		)
	}

	return nil
}

func validateWorkRecordImmutableFields(record *core.Record) error {
	original := record.Original()
	if original == nil || original.Id == "" {
		return nil
	}

	if original.GetBool("approved") {
		return validationFieldError(
			"global",
			"approved_work_record_locked",
			"Approved work records cannot be edited",
		)
	}

	immutableStrings := []struct {
		Name    string
		Code    string
		Message string
	}{
		{Name: "uid", Code: "immutable_field", Message: "Work record owner cannot be changed"},
		{Name: "creator", Code: "immutable_field", Message: "Work record creator cannot be changed"},
		{Name: "number", Code: "immutable_field", Message: "Work record number cannot be changed"},
		{Name: "parent_work_record", Code: "immutable_field", Message: "Parent work record cannot be changed"},
	}
	for _, field := range immutableStrings {
		if original.GetString(field.Name) != record.GetString(field.Name) {
			return validationFieldError(field.Name, field.Code, field.Message)
		}
	}

	if original.GetBool("approved") != record.GetBool("approved") {
		return validationFieldError(
			"approved",
			"immutable_field",
			"Approved state is managed by timesheet workflow",
		)
	}

	return nil
}

func ValidateWorkRecord(app core.App, record *core.Record) error {
	if err := validateWorkRecordImmutableFields(record); err != nil {
		return err
	}
	if err := validateWorkRecordRequiredFields(record); err != nil {
		return err
	}
	if err := validateWorkRecordCoreFieldFormats(record); err != nil {
		return err
	}
	if err := validateWorkRecordVehicleFields(record); err != nil {
		return err
	}
	if err := validateWorkRecordUniqueDay(app, record); err != nil {
		return err
	}
	if err := validateWorkRecordTypeChange(app, record); err != nil {
		return err
	}
	return nil
}

func ValidateWorkRecordDelete(app core.App, record *core.Record) error {
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM time_entries
		WHERE work_record_id = {:workRecordId}
	`).Bind(dbx.Params{
		"workRecordId": record.Id,
	}).One(&row); err != nil {
		return err
	}

	if row.Count > 0 {
		return validationFieldError(
			"global",
			"linked_time_entry_exists",
			"Work records linked to time entries cannot be deleted",
		)
	}

	return nil
}

// ProcessWorkRecord applies work-record-specific save behavior before the
// record is persisted.
//
// Today that means assigning a generated number when the caller did not provide
// one. The hook does not overwrite an existing number because some internal
// import and repair flows may need to preserve seeded values explicitly.
func EnsureWorkRecordNumber(app core.App, record *core.Record) error {
	// Status quo:
	// - route-based callers normally leave number blank and let the hook assign it;
	// - internal/non-request callers may pre-seed a number and this hook
	//   preserves it, which leaves room for controlled import/backfill flows.
	//
	// Future consideration:
	// - if we later want all non-import creates to be strictly backend-numbered,
	//   tighten this branch to require an explicit import/backfill signal.
	if strings.TrimSpace(record.GetString("number")) != "" {
		return nil
	}

	number, err := GenerateWorkRecordNumber(app, record.GetString("date"))
	if err != nil {
		if validationErrs, ok := err.(validation.Errors); ok {
			normalizedErr, _ := NormalizeValidationError(validationErrs)
			return normalizedErr
		}
		return validationFieldError(
			"number",
			"error_generating_work_record_number",
			fmt.Sprintf("error generating work record number: %v", err),
		)
	}

	record.Set("number", number)
	return nil
}

func ProcessWorkRecord(app core.App, e *core.RecordEvent) error {
	record := e.Record
	if err := ValidateWorkRecord(app, record); err != nil {
		return err
	}

	return EnsureWorkRecordNumber(app, record)
}
