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
	"net/http"
	"strings"
	"time"
	"tybalt/errs"
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
		{Name: "creator", Message: "Work record creator is required"},
		{Name: "job", Message: "Job is required"},
		{Name: "type", Message: "Work record type is required"},
		{Name: "location", Message: "Location is required"},
		{Name: "work_description", Message: "Work description is required"},
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

func workRecordHasApprovedSubject(app core.App, workRecordID string) (bool, error) {
	var row struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT EXISTS(
			SELECT 1
			FROM work_records_subjects
			WHERE work_record = {:workRecordId}
			  AND COALESCE(approved, 0) = 1
		) AS approved
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).One(&row); err != nil {
		return false, err
	}
	return row.Approved, nil
}

func validateWorkRecordLaneKms(app core.App, record *core.Record) error {
	if record.Get("lane_kms") == nil || record.GetFloat("lane_kms") == 0 {
		return nil
	}
	if record.GetFloat("lane_kms") < 0 {
		return validationFieldError(
			"lane_kms",
			"lane_kms_min",
			"Lane kilometres cannot be negative",
		)
	}

	var row struct {
		AllowLaneKms bool `db:"allow_lane_kms"`
	}
	err := app.DB().NewQuery(`
		SELECT COALESCE(allow_lane_kms, 0) AS allow_lane_kms
		FROM work_records_types
		WHERE id = {:type}
		LIMIT 1
	`).Bind(dbx.Params{
		"type": record.GetString("type"),
	}).One(&row)
	if err != nil {
		if err == sql.ErrNoRows {
			return validationFieldError("type", "not_found", "Work record type was not found")
		}
		return err
	}
	if !row.AllowLaneKms {
		return validationFieldError(
			"lane_kms",
			"lane_kms_not_allowed",
			"Lane kilometres are not allowed for this work record type",
		)
	}
	return nil
}

func validateWorkRecordImmutableFields(record *core.Record) error {
	original := record.Original()
	if original == nil || original.Id == "" {
		return nil
	}

	immutableStrings := []struct {
		Name    string
		Code    string
		Message string
	}{
		{Name: "creator", Code: "immutable_field", Message: "Work record creator cannot be changed"},
		{Name: "number", Code: "immutable_field", Message: "Work record number cannot be changed"},
		{Name: "parent_work_record", Code: "immutable_field", Message: "Parent work record cannot be changed"},
	}
	for _, field := range immutableStrings {
		if original.GetString(field.Name) != record.GetString(field.Name) {
			return validationFieldError(field.Name, field.Code, field.Message)
		}
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
	if err := validateWorkRecordLaneKms(app, record); err != nil {
		return err
	}
	if err := validateWorkRecordTypeChange(app, record); err != nil {
		return err
	}
	if record.Original() != nil && record.Original().Id != "" {
		locked, err := workRecordHasApprovedSubject(app, record.Id)
		if err != nil {
			return err
		}
		if locked {
			return validationFieldError(
				"global",
				"approved_work_record_locked",
				"Work records with approved workers cannot be edited",
			)
		}
	}
	return nil
}

func ValidateWorkRecordDelete(app core.App, record *core.Record) error {
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_subjects wrs
		JOIN time_entries te
		  ON te.work_record_subject_id = wrs.id
		WHERE wrs.work_record = {:workRecordId}
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

func validateWorkRecordSubjectRequiredFields(record *core.Record) error {
	return validation.Errors{
		"work_record":       validation.Validate(record.GetString("work_record"), validation.Required.Error("Work record is required")),
		"uid":               validation.Validate(record.GetString("uid"), validation.Required.Error("Worker is required")),
		"hours_on_site":     validation.Validate(record.Get("hours_on_site"), validation.Required.Error("Hours on site is required")),
		"hours_travel_time": validation.Validate(record.Get("hours_travel_time"), validation.Required.Error("Travel time is required")),
		"vehicle_type":      validation.Validate(record.GetString("vehicle_type"), validation.In("", "company", "personal").Error("Vehicle type must be company or personal")),
		"hours_on_site_min": validation.Validate(record.GetFloat("hours_on_site"), validation.Min(0.0).Error("Hours on site must be at least 0")),
		"travel_time_min":   validation.Validate(record.GetFloat("hours_travel_time"), validation.Min(0.0).Error("Travel time must be at least 0")),
		"distance_km_min":   validation.Validate(record.GetFloat("distance_travelled_km"), validation.Min(0.0).Error("Distance travelled must be at least 0")),
	}.Filter()
}

func validateWorkRecordSubjectVehicleFields(record *core.Record) error {
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

func countSubjectsForWorkRecord(app core.App, workRecordID string, excludeSubjectID string) (int, error) {
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_subjects
		WHERE work_record = {:workRecordId}
		  AND id != {:excludeSubjectId}
	`).Bind(dbx.Params{
		"workRecordId":     workRecordID,
		"excludeSubjectId": excludeSubjectID,
	}).One(&row); err != nil {
		return 0, err
	}
	return row.Count, nil
}

func validateWorkRecordSubjectMax(app core.App, record *core.Record) error {
	var row struct {
		MaxSubjects int `db:"max_subjects"`
	}
	err := app.DB().NewQuery(`
		SELECT COALESCE(wrt.max_subjects, 1) AS max_subjects
		FROM work_records wr
		JOIN work_records_types wrt
		  ON wrt.id = wr.type
		WHERE wr.id = {:workRecordId}
		LIMIT 1
	`).Bind(dbx.Params{
		"workRecordId": record.GetString("work_record"),
	}).One(&row)
	if err != nil {
		if err == sql.ErrNoRows {
			return validationFieldError("work_record", "not_found", "Work record was not found")
		}
		return err
	}
	count, err := countSubjectsForWorkRecord(app, record.GetString("work_record"), record.Id)
	if err != nil {
		return err
	}
	if count+1 > row.MaxSubjects {
		return validationFieldError(
			"uid",
			"max_subjects_exceeded",
			fmt.Sprintf("This work record type allows at most %d workers", row.MaxSubjects),
		)
	}
	return nil
}

func validateWorkRecordSubjectDuplicateDay(app core.App, record *core.Record) error {
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_subjects wrs
		JOIN work_records wr
		  ON wr.id = wrs.work_record
		JOIN work_records current_wr
		  ON current_wr.id = {:workRecordId}
		WHERE wrs.uid = {:uid}
		  AND wrs.id != {:subjectId}
		  AND wr.job = current_wr.job
		  AND wr.date = current_wr.date
	`).Bind(dbx.Params{
		"workRecordId": record.GetString("work_record"),
		"uid":          record.GetString("uid"),
		"subjectId":    record.Id,
	}).One(&row); err != nil {
		return err
	}
	if row.Count > 0 {
		return validationFieldError(
			"date",
			"duplicate_work_record_day",
			"A worker already has a work record for this job and date",
		)
	}
	return nil
}

func validateWorkRecordSubjectImmutableFields(record *core.Record) error {
	original := record.Original()
	if original == nil || original.Id == "" {
		return nil
	}

	if original.GetString("work_record") != record.GetString("work_record") {
		return validationFieldError("work_record", "immutable_field", "Work record cannot be changed")
	}
	if original.GetString("uid") != record.GetString("uid") {
		return validationFieldError("uid", "immutable_field", "Worker cannot be reassigned")
	}
	if original.GetBool("approved") {
		return validationFieldError("global", "approved_subject_locked", "Approved work record workers cannot be edited")
	}
	return nil
}

func ValidateWorkRecordSubject(app core.App, record *core.Record) error {
	if err := validateWorkRecordSubjectImmutableFields(record); err != nil {
		return err
	}
	if err := validateWorkRecordSubjectRequiredFields(record); err != nil {
		return err
	}
	if err := validateWorkRecordSubjectVehicleFields(record); err != nil {
		return err
	}
	if err := validateWorkRecordSubjectMax(app, record); err != nil {
		return err
	}
	if err := validateWorkRecordSubjectDuplicateDay(app, record); err != nil {
		return err
	}
	locked, err := workRecordHasApprovedSubject(app, record.GetString("work_record"))
	if err != nil {
		return err
	}
	if locked && (record.IsNew() || record.GetString("uid") != record.Original().GetString("uid")) {
		return validationFieldError("global", "approved_roster_locked", "Workers cannot be added or reassigned after approval")
	}
	return nil
}

func ValidateWorkRecordSubjectDelete(app core.App, record *core.Record) error {
	locked, err := workRecordHasApprovedSubject(app, record.GetString("work_record"))
	if err != nil {
		return err
	}
	if locked {
		return validationFieldError("global", "approved_roster_locked", "Workers cannot be removed after approval")
	}
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM time_entries
		WHERE work_record_subject_id = {:subjectId}
	`).Bind(dbx.Params{
		"subjectId": record.Id,
	}).One(&row); err != nil {
		return err
	}
	if row.Count > 0 {
		return validationFieldError("global", "linked_time_entry_exists", "Workers linked to time entries cannot be deleted")
	}
	return nil
}

// SyncWorkRecordAttachmentFields keeps the attachment file field and its
// derived hash in sync before a work record is persisted.
//
// This helper is deliberately side-effect-free outside the record being saved:
// it prepares attachment-related fields, but it does not perform the save
// itself. That lets custom routes compose it with their own access checks and
// transaction boundaries.
//
// Current responsibilities:
//   - normalize pending file names so the stored file keys are URL-safe;
//   - hash newly uploaded PDFs and store the digest in attachment_hash;
//   - reject duplicates across other work_records rows; and
//   - clear attachment_hash when the attachment itself is removed.
//
// Importantly, the dedicated attachment route calls this helper directly as part
// of the primary user-facing upload/remove flow, and the work_records request
// hooks also call it as an integrity backstop for collection-based writes such
// as admin edits. That split keeps the public write surface narrow while still
// centralizing attachment storage rules in one place.
func SyncWorkRecordAttachmentFields(app core.App, record *core.Record) error {
	NormalizePendingFileNames(record, "attachment")

	attachmentHash, hashErr := CalculateRecordFileFieldHash(record, "attachment")
	if hashErr != nil {
		return hashErr
	}

	if attachmentHash != "" {
		existingRecord, _ := app.FindFirstRecordByFilter(
			"work_records",
			"attachment_hash = {:hash} && id != {:id}",
			dbx.Params{
				"hash": attachmentHash,
				"id":   record.Id,
			},
		)
		if existingRecord != nil {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "duplicate attachment detected",
				Data: map[string]errs.CodeError{
					"attachment": {
						Code:    "duplicate_file",
						Message: "This file has already been uploaded to another work record",
					},
				},
			}
		}

		record.Set("attachment_hash", attachmentHash)
	} else if record.GetString("attachment") == "" {
		record.Set("attachment_hash", "")
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
