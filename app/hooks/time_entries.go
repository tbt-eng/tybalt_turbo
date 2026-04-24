// this file implements cleaning and validation rules for the time_entries collection

package hooks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// The cleanTimeEntry function is used to remove properties from the time_entry
// record that are not allowed to be set based on the value of the record's
// time_type property. This is intended to reduce round trips to the database
// and to ensure that the record is in a valid state before it is created or
// updated. It is called by ProcessTimeEntry to reduce the number of fields
// that need to be validated.
func cleanTimeEntry(app core.App, timeEntryRecord *core.Record) ([]string, error) {
	timeTypeId := timeEntryRecord.GetString("time_type")

	// Load the allowed fields for the time_type from the time_types collection in
	// PocketBase. They are stored in the allowed_fields property as a JSON array
	// of strings.
	timeTypeRecord, findError := app.FindRecordById("time_types", timeTypeId)
	if findError != nil {
		return nil, findError
	}

	// Get the allowed fields from the time_type record with type assertion
	var allowedFields []string
	allowedFieldsJson := timeTypeRecord.Get("allowed_fields").(types.JSONRaw)

	// use json.Unmarshal to convert the JSON array to a Go slice of strings
	if unmarshalErrorAllowed := json.Unmarshal(allowedFieldsJson, &allowedFields); unmarshalErrorAllowed != nil {
		log.Fatalf("Error unmarshalling JSON: %v", unmarshalErrorAllowed)
		return nil, unmarshalErrorAllowed
	}

	// Get the required fields from the time_type record with type assertion
	var requiredFields []string
	requiredFieldsJson := timeTypeRecord.Get("required_fields").(types.JSONRaw)

	// if requiredFieldsJson has a value of null, "{}", or "[]" then all fields
	// are required so we set the requiredFields to the allowedFields.
	if requiredFieldsJson.String() == "null" || requiredFieldsJson.String() == "[]" || requiredFieldsJson.String() == "{}" {
		requiredFields = allowedFields
	} else if unmarshalErrorRequired := json.Unmarshal(requiredFieldsJson, &requiredFields); unmarshalErrorRequired != nil {
		log.Fatalf("Error unmarshalling JSON: %v", unmarshalErrorRequired)
		return nil, unmarshalErrorRequired
	}

	// Certain fields are always allowed to be set. We add them to the list of
	// allowed fields here. The role field is included because its requirement
	// depends on job assignment, not the time_type. The branch field is resolved
	// below and must survive cleanup.
	allowedFields = append(allowedFields, "id", "uid", "created", "updated", "role", "branch")
	if slices.Contains(allowedFields, "work_record") && !slices.Contains(allowedFields, "work_record_subject_id") {
		allowedFields = append(allowedFields, "work_record_subject_id")
	}

	// remove any fields from the time_entry record that are not in allowedFields.
	// I'm not sure if this is the best way to do this but let's try it.
	// FieldsData() is probably a drop in replacement for ColumnValueMap()
	for _, key := range timeEntryRecord.Collection().Fields.FieldNames() {
		if !slices.Contains(allowedFields, key) {
			//log.Println("Removing field: ", key)
			timeEntryRecord.Set(key, nil)
		}
	}

	// Clear role if no job is assigned
	if timeEntryRecord.GetString("job") == "" {
		timeEntryRecord.Set("role", nil)
	}

	// Match purchase order branch resolution:
	// - when a job is set, branch must always match the job's branch
	// - without a job, derive from the creator's default branch only when blank
	jobId := strings.TrimSpace(timeEntryRecord.GetString("job"))
	if jobId != "" {
		jobRecord, err := app.FindRecordById("jobs", jobId)
		if err != nil || jobRecord == nil {
			return nil, &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when cleaning time entry",
				Data: map[string]errs.CodeError{
					"job": {Code: "not_found", Message: "referenced job not found"},
				},
			}
		}
		jobBranchID := strings.TrimSpace(jobRecord.GetString("branch"))
		if jobBranchID == "" {
			return nil, &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when cleaning time entry",
				Data: map[string]errs.CodeError{
					"job": {Code: "missing_branch", Message: "referenced job is missing a branch"},
				},
			}
		}
		timeEntryRecord.Set("branch", jobBranchID)
	} else {
		uid := timeEntryRecord.GetString("uid")
		adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
		if err != nil {
			return nil, err
		}
		defaultBranchID := strings.TrimSpace(adminProfile.GetString("default_branch"))
		if defaultBranchID == "" {
			return nil, fmt.Errorf("your admin_profiles record is missing a default_branch")
		}
		if strings.TrimSpace(timeEntryRecord.GetString("branch")) == "" {
			timeEntryRecord.Set("branch", defaultBranchID)
		}
	}

	return requiredFields, nil
}

func appendValidationMessage(existing error, code string, message string) error {
	if existing == nil {
		return validation.NewError(code, message)
	}

	existingMessage := strings.TrimSpace(existing.Error())
	if existingMessage == "" {
		return validation.NewError(code, message)
	}

	return validation.NewError(code, existingMessage+"; "+message)
}

func validateLinkedWorkRecordSubject(app core.App, timeEntryRecord *core.Record) validation.Errors {
	subjectID := strings.TrimSpace(timeEntryRecord.GetString("work_record_subject_id"))
	if subjectID == "" {
		return nil
	}

	validations := validation.Errors{}

	var subject struct {
		UID             string  `db:"uid"`
		Job             string  `db:"job"`
		Date            string  `db:"date"`
		HoursOnSite     float64 `db:"hours_on_site"`
		HoursTravelTime float64 `db:"hours_travel_time"`
	}
	err := app.DB().NewQuery(`
		SELECT
			wrs.uid AS uid,
			COALESCE(wr.job, '') AS job,
			wr.date AS date,
			COALESCE(wrs.hours_on_site, 0) AS hours_on_site,
			COALESCE(wrs.hours_travel_time, 0) AS hours_travel_time
		FROM work_records_subjects wrs
		JOIN work_records wr
		  ON wr.id = wrs.work_record
		WHERE wrs.id = {:subjectId}
		LIMIT 1
	`).Bind(dbx.Params{
		"subjectId": subjectID,
	}).One(&subject)
	if err != nil {
		validations["work_record_subject_id"] = validation.NewError("not_found", "Referenced work record subject not found")
		return validations
	}

	var duplicate struct {
		ID string `db:"id"`
	}
	err = app.DB().NewQuery(`
		SELECT id
		FROM time_entries
		WHERE work_record_subject_id = {:subjectId}
		  AND id != {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"subjectId": subjectID,
		"id":        timeEntryRecord.Id,
	}).One(&duplicate)
	if err == nil && duplicate.ID != "" {
		validations["work_record_subject_id"] = appendValidationMessage(
			validations["work_record_subject_id"],
			"work_record_subject_already_linked",
			"This work record subject is already linked to another time entry",
		)
	}

	if subject.UID != timeEntryRecord.GetString("uid") {
		validations["work_record_subject_id"] = appendValidationMessage(
			validations["work_record_subject_id"],
			"work_record_uid_mismatch",
			"Selected work record subject belongs to a different user",
		)
	}

	if subject.Job != timeEntryRecord.GetString("job") {
		validations["work_record_subject_id"] = appendValidationMessage(
			validations["work_record_subject_id"],
			"work_record_job_mismatch",
			"Selected work record is for a different job",
		)
	}

	if subject.Date != timeEntryRecord.GetString("date") {
		validations["work_record_subject_id"] = appendValidationMessage(
			validations["work_record_subject_id"],
			"work_record_date_mismatch",
			"Selected work record is for a different date",
		)
	}

	expectedHours := subject.HoursOnSite + subject.HoursTravelTime
	if timeEntryRecord.GetFloat("hours") != expectedHours {
		validations["hours"] = validation.NewError(
			"work_record_hours_mismatch",
			"Hours must match the linked work record subject. Edit the work record subject or create a separate unlinked time entry.",
		)
	}

	return validations
}

func legacyWorkRecordMatchesExistingWorkRecordNumber(app core.App, workRecordText string) (bool, error) {
	workRecordText = strings.TrimSpace(workRecordText)
	if workRecordText == "" {
		return false, nil
	}

	var match struct {
		ID string `db:"id"`
	}
	err := app.DB().NewQuery(`
		SELECT id
		FROM work_records
		WHERE number = {:workRecord}
		LIMIT 1
	`).Bind(dbx.Params{
		"workRecord": workRecordText,
	}).One(&match)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return match.ID != "", nil
}

func legacyWorkRecordUpgradeRequired(timeEntryRecord *core.Record) bool {
	if timeEntryRecord.IsNew() {
		return false
	}

	original := timeEntryRecord.Original()
	if original == nil {
		return false
	}

	return strings.TrimSpace(original.GetString("work_record")) != "" &&
		strings.TrimSpace(original.GetString("work_record_subject_id")) == ""
}

// cross-field validation is performed in this function. It is expected that the
// time_entry record has already been cleaned by the cleanTimeEntry function.
// This ensures that only the fields that are allowed to be set are present in
// the record prior to validation. The requiredFields slice is used to validate
// the presence of required fields. The function returns an error if the record
// is invalid, otherwise it returns nil.
func validateTimeEntry(app core.App, timeEntryRecord *core.Record, requiredFields []string) error {
	jobIsPresent := timeEntryRecord.GetString("job") != ""
	totalHours := timeEntryRecord.GetFloat("hours") + timeEntryRecord.GetFloat("meals_hours")
	workRecordSubjectID := strings.TrimSpace(timeEntryRecord.GetString("work_record_subject_id"))
	workRecordText := strings.TrimSpace(timeEntryRecord.GetString("work_record"))
	linkedWorkRecordsEnabled, err := utilities.IsLinkedWorkRecordsEnabled(app)
	if err != nil {
		return err
	}

	if workRecordSubjectID != "" {
		timeEntryRecord.Set("work_record", "")
		workRecordText = ""
	}

	// validation is performed in two passes. The first pass is to validate the
	// presence of required fields. We do this by using the validation.Required
	// function on each required field from the requiredFields slice using
	// validation.Errors with the field name as the key.
	requiredValidationsErrors := validation.Errors{}
	for _, field := range requiredFields {
		requiredValidationsErrors[field] = validation.Validate(timeEntryRecord.Get(field), validation.Required.Error("Value required"))
	}

	// Check for errors in the first pass
	if err := requiredValidationsErrors.Filter(); err != nil {
		return err
	}

	// The second pass performs everything else (cross-field validation, field
	// values, etc.)
	otherValidationsErrors := validation.Errors{
		"hours":                 validation.Validate(timeEntryRecord.Get("hours"), validation.By(utilities.IsPositiveMultipleOfPointFive())),
		"date":                  validation.Validate(timeEntryRecord.Get("date"), validation.By(utilities.IsValidDate)),
		"global":                validation.Validate(totalHours, validation.Max(18.0).Error("Total hours must not exceed 18")),
		"meals_hours":           validation.Validate(timeEntryRecord.Get("meals_hours"), validation.Max(3.0).Error("Meals Hours must not exceed 3")),
		"description":           validation.Validate(timeEntryRecord.GetString("description"), validation.Length(5, 0).Error("must be at least 5 characters")),
		"work_record":           validation.Validate(timeEntryRecord.Get("work_record"), validation.When(workRecordSubjectID == "" && workRecordText != "" && jobIsPresent, validation.Match(regexp.MustCompile("^[FKQ][0-9]{2}-[0-9]{3,4}(-[0-9]+)?$")).Error("must be in the correct format")).Else(validation.In("").Error("Work Record must be empty when job is not provided"))),
		"payout_request_amount": validation.Validate(timeEntryRecord.Get("payout_request_amount"), validation.Min(0.0).Exclusive().Error("Amount must be greater than 0"), validation.By(utilities.IsPositiveMultipleOfPointZeroOne())),
	}

	if workRecordSubjectID != "" {
		for field, err := range validateLinkedWorkRecordSubject(app, timeEntryRecord) {
			otherValidationsErrors[field] = err
		}
		timeEntryRecord.Set("work_record", "")
	} else {
		legacyUpgradeRequired := linkedWorkRecordsEnabled && legacyWorkRecordUpgradeRequired(timeEntryRecord)
		if workRecordText != "" {
			// Today this check is effectively a no-op for standard app-created data:
			// legacy time entries only accept F/K/Q-prefixed values, while first-class
			// work records are generated with W-prefixed numbers. The overlap that
			// would make this branch matter can only happen if import/backfill flows
			// preserve a non-W work_records.number, or if future changes broaden the
			// legacy regex / first-class numbering format so both sources can emit the
			// same visible identifier. If that happens, this validation may start
			// blocking edits to existing legacy entries in hybrid mode and should be
			// revisited alongside those format changes.
			duplicateNumber, duplicateErr := legacyWorkRecordMatchesExistingWorkRecordNumber(app, workRecordText)
			if duplicateErr != nil {
				return duplicateErr
			}
			if duplicateNumber {
				otherValidationsErrors["work_record"] = validation.NewError(
					"work_record_number_reserved",
					"This work record number already belongs to a linked work record. Start from Work Records instead.",
				)
			}
		}
		if legacyUpgradeRequired {
			otherValidationsErrors["work_record_subject_id"] = validation.NewError(
				"work_record_upgrade_required",
				"Select an existing work record subject before saving this legacy time entry",
			)
		} else if linkedWorkRecordsEnabled && workRecordText != "" {
			otherValidationsErrors["work_record"] = validation.NewError(
				"linked_work_records_required",
				"Text-only work records are disabled. Start from Work Records to create a linked time entry.",
			)
		}
	}

	// If a job is present, verify the referenced job exists and has Active status
	if jobIsPresent {
		jobID := timeEntryRecord.GetString("job")
		jobRecord, err := app.FindRecordById("jobs", jobID)
		if err != nil || jobRecord == nil {
			otherValidationsErrors["job"] = validation.Validate(jobID, validation.Required.Error("invalid job reference"))
		} else {
			if err := validateJobAllowsTimeTracking(app, jobRecord); err != nil {
				otherValidationsErrors["job"] = err
			}
			// Role is required when a job is assigned
			if timeEntryRecord.GetString("role") == "" {
				otherValidationsErrors["role"] = validation.Validate(timeEntryRecord.Get("role"), validation.Required.Error("Role is required when a job is assigned"))
			}
			// Time entries have always enforced allocation membership for job+division.
			// Use the shared helper so the exact error codes/messages stay consistent
			// with PO/expense/time-amendment validation.
			if field, allocErr := validateDivisionAllocatedToJob(app, jobID, timeEntryRecord.GetString("division")); allocErr != nil {
				otherValidationsErrors[field] = allocErr
			}
		}
	}

	return otherValidationsErrors.Filter()
}

// The ProcessTimeEntry function is used to validate the time_entry record
// before it is created or updated. A lot of the work is done by PocketBase
// itself so this is for cross-field validation. If the time_entry record is
// invalid this function throws an error explaining which field(s) are invalid
// and why.
func ProcessTimeEntry(app core.App, e *core.RecordRequestEvent) error {
	record := e.Record
	// get the auth record from the context
	authRecord := e.Auth

	enabled, err := utilities.IsTimeEditingEnabled(app)
	if err != nil {
		return err
	}
	if !enabled {
		return utilities.ErrTimeEditingDisabled
	}

	// If the uid property is not equal to the authenticated user's uid, return an
	// error.
	if record.GetString("uid") != authRecord.Id {
		return apis.NewApiError(400, "uid property must be equal to the authenticated user's id", map[string]validation.Error{})
	}

	// set properties to nil if they are not allowed to be set based on the
	// time_type
	requiredFields, cleanErr := cleanTimeEntry(app, record)
	if cleanErr != nil {
		var hookErr *errs.HookError
		if errors.As(cleanErr, &hookErr) {
			return hookErr
		}
		return apis.NewBadRequestError("Error cleaning time_entry record", cleanErr)
	}

	if err := EnsureActiveDivision(app, record.GetString("division"), "division"); err != nil {
		if ve, ok := err.(validation.Errors); ok {
			return apis.NewBadRequestError("Validation error", ve)
		}
		return err
	}

	if err := utilities.EnsureUserCanUseRecordBranch(app, record, authRecord.Id); err != nil {
		return err
	}

	// write the week_ending property to the record. This is derived exclusively
	// from the date property.
	weekEnding, wkEndErr := utilities.GenerateWeekEnding(record.GetString("date"))
	if wkEndErr != nil {
		return apis.NewBadRequestError("Error generating week_ending", wkEndErr)
	}
	record.Set("week_ending", weekEnding)

	// perform the validation for the time_entry record. In this step we also
	// write the uid property to the record so that we can validate it against the
	if validationErr := validateTimeEntry(app, record, requiredFields); validationErr != nil {
		return apis.NewBadRequestError("Validation error", validationErr)
	}

	if strings.TrimSpace(record.GetString("work_record_subject_id")) != "" {
		record.Set("work_record", "")
	}

	return nil
}
