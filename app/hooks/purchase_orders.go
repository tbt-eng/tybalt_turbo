// This file implements cleaning and validation rules for the purchase_orders
// collection.

package hooks

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
	"tybalt/constants"
	"tybalt/errs"
	"tybalt/notifications"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// CleanPurchaseOrder removes properties from the
// purchase_order record that are not allowed to be set based on the value of
// the record's type property. It is also used to set the approval_total field,
// which matches the total field unless the type is Recurring. It is called by
// ProcessPurchaseOrder to reduce the number of fields that need to be
// validated.
func CleanPurchaseOrder(app core.App, purchaseOrderRecord *core.Record) error {
	// initialize approval_total to total and occurrences to zero. These will be
	// changed if the PO is recurring.
	purchaseOrderRecord.Set("approval_total", purchaseOrderRecord.GetFloat("total"))
	purchaseOrderRecord.Set("occurrences", 0)
	homeCurrencyID := ""
	if homeCurrency, err := utilities.FindHomeCurrency(app); err == nil && homeCurrency != nil {
		homeCurrencyID = homeCurrency.Id
	}

	typeString := purchaseOrderRecord.GetString("type")

	// One-Time and Cumulative purchase orders both have empty values for
	// end_date and frequency
	if typeString == "One-Time" || typeString == "Cumulative" {
		purchaseOrderRecord.Set("end_date", "")
		purchaseOrderRecord.Set("frequency", "")
	}

	if typeString == "Recurring" {
		occurrences, calculatedTotal, err := utilities.CalculateRecurringPurchaseOrderTotalValue(app, purchaseOrderRecord)
		if err != nil {
			var hookErr *errs.HookError
			if errors.As(err, &hookErr) {
				return err
			} else {
				return &errs.HookError{
					Status:  http.StatusInternalServerError,
					Message: "hook error when calculating recurring PO total value",
					Data: map[string]errs.CodeError{
						"global": {
							Code:    "error_calculating_total_value",
							Message: fmt.Sprintf("error calculating recurring PO total value: %v", err),
						},
					},
				}
			}
		}
		purchaseOrderRecord.Set("approval_total", calculatedTotal)
		purchaseOrderRecord.Set("occurrences", occurrences)
	}

	if strings.TrimSpace(purchaseOrderRecord.GetString("currency")) == "" {
		parentPOID := strings.TrimSpace(purchaseOrderRecord.GetString("parent_po"))
		if parentPOID != "" {
			if parentPO, err := app.FindRecordById("purchase_orders", parentPOID); err == nil && parentPO != nil {
				if parentCurrencyID := strings.TrimSpace(parentPO.GetString("currency")); parentCurrencyID != "" {
					purchaseOrderRecord.Set("currency", parentCurrencyID)
				}
			}
		}
	}
	if strings.TrimSpace(purchaseOrderRecord.GetString("currency")) == "" && homeCurrencyID != "" {
		purchaseOrderRecord.Set("currency", homeCurrencyID)
	}

	currencyInfo, err := utilities.ResolveCurrencyInfo(app, purchaseOrderRecord.GetString("currency"))
	if err != nil {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when cleaning purchase order",
			Data: map[string]errs.CodeError{
				"currency": {Code: "not_found", Message: "referenced currency not found"},
			},
		}
	}
	if err := utilities.RequirePositiveForeignCurrencyRate(currencyInfo); err != nil {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when cleaning purchase order",
			Data: map[string]errs.CodeError{
				"currency": {
					Code:    "missing_rate",
					Message: "selected foreign currency is missing an exchange rate",
				},
			},
		}
	}
	purchaseOrderRecord.Set("approval_total_home", purchaseOrderRecord.GetFloat("approval_total")*utilities.CurrencyRateOrOne(currencyInfo))

	// Clear all rejection fields here. ProcessPurchaseOrder, which calls
	// cleanPurchaseOrder, is only ever called when a user is creating or
	// updating a PO. POs cannot be rejected upon creation, so clearing rejection
	// fields is idempotent. Upon update, rejection fields should be cleared if
	// any changes are made to the record, on the assumption that the user is
	// preparing to resubmit the PO after making changes.
	purchaseOrderRecord.Set("rejected", "")
	purchaseOrderRecord.Set("rejection_reason", "")
	purchaseOrderRecord.Set("rejector", "")

	// When a job is present, branch must always match the job's branch.
	// Without a job, derive from the owner's default branch only when blank, so
	// an owner without a default branch can still have a chosen branch.
	var derivedBranchID string
	jobId := purchaseOrderRecord.GetString("job")
	if jobId != "" {
		jobRecord, err := app.FindRecordById("jobs", jobId)
		if err != nil || jobRecord == nil {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when cleaning purchase order",
				Data: map[string]errs.CodeError{
					"job": {Code: "not_found", Message: "referenced job not found"},
				},
			}
		}
		derivedBranchID = jobRecord.GetString("branch")
		if derivedBranchID == "" {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when cleaning purchase order",
				Data: map[string]errs.CodeError{
					"job": {Code: "missing_branch", Message: "referenced job is missing a branch"},
				},
			}
		}
		// Enforce strict job->branch alignment.
		purchaseOrderRecord.Set("branch", derivedBranchID)
	} else if strings.TrimSpace(purchaseOrderRecord.GetString("branch")) == "" {
		uid := purchaseOrderRecord.GetString("uid")
		adminProfile, err := app.FindFirstRecordByFilter("admin_profiles", "uid={:uid}", dbx.Params{"uid": uid})
		if err != nil || adminProfile == nil {
			return &errs.HookError{
				Status:  http.StatusInternalServerError,
				Message: "hook error when cleaning purchase order",
				Data: map[string]errs.CodeError{
					"uid": {Code: "profile_lookup_error", Message: "error looking up admin profile"},
				},
			}
		}
		derivedBranchID = adminProfile.GetString("default_branch")
		if derivedBranchID == "" {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when cleaning purchase order",
				Data: map[string]errs.CodeError{
					"uid": {Code: "missing_default_branch", Message: "your admin_profiles record is missing a default_branch"},
				},
			}
		}
		purchaseOrderRecord.Set("branch", derivedBranchID)
	}

	return nil
}

// ValidatePurchaseOrder performs validation, including cross-field validation.
// It is expected that the purchase_order record has already been cleaned by the
// CleanPurchaseOrder
// function. This ensures that only the fields that are allowed to be set are
// present in the record prior to validation. The function returns an error if
// the record is invalid, otherwise it returns nil.
func ValidatePurchaseOrder(app core.App, purchaseOrderRecord *core.Record, legacyMode bool) error {
	isRecurring := purchaseOrderRecord.GetString("type") == "Recurring"
	isChild := purchaseOrderRecord.GetString("parent_po") != ""
	var linkedJobRecord *core.Record

	// Provide friendly errors for required fields that PocketBase schema would
	// otherwise reject with a generic message. We have intentionally relaxed
	// the api rules in pocketbase to perform validation here instead. This allows
	// us to provide more helpful error messages since PocketBase performs API
	// rules validation prior to this hook.

	// TODO: Many of these validations return immediately, but we should return
	// them within the validation.Errors object instead so that the user can see
	// all the errors at once. This can probably be accomplished using a
	// combination of built-in validation.Validate() and custom
	// validation.RuleFunc.

	if purchaseOrderRecord.GetString("vendor") == "" {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating purchase order",
			Data: map[string]errs.CodeError{
				"vendor": {
					Code:    "required",
					Message: "vendor is required",
				},
			},
		}
	}

	kindID := strings.TrimSpace(purchaseOrderRecord.GetString("kind"))
	if kindID == "" {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating purchase order",
			Data: map[string]errs.CodeError{
				"kind": {
					Code:    "required",
					Message: "kind is required",
				},
			},
		}
	}

	kindRecord, err := app.FindRecordById("expenditure_kinds", kindID)
	if err != nil {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating purchase order",
			Data: map[string]errs.CodeError{
				"kind": {
					Code:    "not_found",
					Message: "invalid expenditure kind",
				},
			},
		}
	}
	kindAllowsJob := kindRecord.GetBool("allow_job")
	hasJob := strings.TrimSpace(purchaseOrderRecord.GetString("job")) != ""
	if hasJob && !kindAllowsJob {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating purchase order",
			Data: map[string]errs.CodeError{
				"kind": {
					Code:    "invalid_kind_for_job",
					Message: "selected kind does not allow job",
				},
			},
		}
	}
	if kindRecord.GetString("name") == utilities.ExpenditureKindNameProject && !hasJob {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating purchase order",
			Data: map[string]errs.CodeError{
				"job": {
					Code:    "job_required_for_kind",
					Message: "a job is required for the project expenditure kind",
				},
			},
		}
	}

	// If a vendor is provided, ensure it exists and is Active
	if vendorId := purchaseOrderRecord.GetString("vendor"); vendorId != "" {
		vendorRecord, err := app.FindRecordById("vendors", vendorId)
		if err != nil || vendorRecord == nil {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating vendor",
				Data: map[string]errs.CodeError{
					"vendor": {
						Code:    "not_found",
						Message: "vendor not found",
					},
				},
			}
		}
		if vendorRecord.GetString("status") != "Active" {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating vendor",
				Data: map[string]errs.CodeError{
					"vendor": {
						Code:    "inactive_vendor",
						Message: "provided vendor is not active",
					},
				},
			}
		}
	}

	// The total must be greater than 0.
	if purchaseOrderRecord.GetFloat("total") <= 0 {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating purchase order",
			Data: map[string]errs.CodeError{
				"total": {
					Code:    "required",
					Message: "total must be greater than 0",
				},
			},
		}
	}

	if isChild {
		// Validate parent PO is active and cumulative
		parentPO, err := app.FindRecordById("purchase_orders", purchaseOrderRecord.GetString("parent_po"))
		if err != nil {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when fetching parent PO",
				Data: map[string]errs.CodeError{
					"parent_po": {
						Code:    "not_found",
						Message: "parent PO not found",
					},
				},
			}
		}

		// Validate parent PO is not itself a child
		if parentPO.GetString("parent_po") != "" {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating parent PO",
				Data: map[string]errs.CodeError{
					"parent_po": {
						Code:    "child_po_cannot_be_parent",
						Message: "parent PO cannot be itself a child",
					},
				},
			}
		}

		if parentPO.GetString("status") != "Active" {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating parent PO",
				Data: map[string]errs.CodeError{
					"parent_po": {
						Code:    "invalid_status",
						Message: "parent PO must be active",
					},
				},
			}
		}

		if parentPO.GetString("type") != "Cumulative" {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating parent PO",
				Data: map[string]errs.CodeError{
					"parent_po": {
						Code:    "invalid_type",
						Message: "parent PO must be cumulative",
					},
				},
			}
		}

		if strings.TrimSpace(purchaseOrderRecord.GetString("currency")) == "" && strings.TrimSpace(parentPO.GetString("currency")) != "" {
			purchaseOrderRecord.Set("currency", parentPO.GetString("currency"))
		}

		// Validate that other child POs with status "Unapproved" do not exist.
		// Exclude this record so that an existing child can be edited or revised.
		otherChildPOs, err := app.FindRecordsByFilter("purchase_orders", "parent_po = {:parentId} && status != 'Closed' && status != 'Cancelled' && id != {:id}", "", 0, 0, dbx.Params{
			"parentId": parentPO.Id,
			"id":       purchaseOrderRecord.Id,
		})
		if err != nil {
			return &errs.HookError{
				Status:  http.StatusInternalServerError,
				Message: "hook error when fetching other child POs",
				Data: map[string]errs.CodeError{
					"parent_po": {
						Code:    "internal_server_error",
						Message: "error searching for other child POs",
					},
				},
			}
		}

		if len(otherChildPOs) > 0 {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating parent PO",
				Data: map[string]errs.CodeError{
					"parent_po": {
						Code:    "existing_children_with_blocking_status",
						Message: "other child POs that are not 'Closed' or 'Cancelled' already exist",
					},
				},
			}
		}

		// Validate fields match parent PO
		fieldsToMatch := []string{"job", "payment_type", "category", "description", "vendor", "kind"}
		for _, field := range fieldsToMatch {
			childValue := purchaseOrderRecord.GetString(field)
			parentValue := parentPO.GetString(field)
			if childValue != parentValue {
				return &errs.HookError{
					Status:  http.StatusBadRequest,
					Message: "hook error when validating parent PO",
					Data: map[string]errs.CodeError{
						field: {
							Code:    "value_mismatch",
							Message: fmt.Sprintf("field %s must match parent PO's %s", field, field),
						},
					},
				}
			}
		}

		parentCurrencyCode := utilities.EffectiveCurrencyCode(app, parentPO.GetString("currency"))
		childCurrencyCode := utilities.EffectiveCurrencyCode(app, purchaseOrderRecord.GetString("currency"))
		if parentCurrencyCode != childCurrencyCode {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating parent PO",
				Data: map[string]errs.CodeError{
					"currency": {
						Code:    "value_mismatch",
						Message: "child PO currency must match parent PO currency",
					},
				},
			}
		}
	}

	dateAsTime, parseErr := time.Parse("2006-01-02", purchaseOrderRecord.Get("date").(string))
	if parseErr != nil {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating date",
			Data: map[string]errs.CodeError{
				"date": {
					Code:    "invalid_date",
					Message: "date must be a valid date",
				},
			},
		}
	}

	// Validate that a selected user is active.
	userIsActive := func(app core.App, code string, message string) validation.RuleFunc {
		return func(value any) error {
			userID, _ := value.(string)
			if userID == "" {
				return nil // Let required validation handle empty
			}
			active, err := utilities.IsUserActive(app, userID)
			if err != nil {
				return &errs.HookError{
					Status:  http.StatusInternalServerError,
					Message: "failed to check user active status",
				}
			}
			if !active {
				return validation.NewError(code, message)
			}
			return nil
		}
	}

	// Validate that priority_second_approver is an active user if set.
	prioritySecondApproverIsActive := func(app core.App, purchaseOrderRecord *core.Record) validation.RuleFunc {
		return func(value any) error {
			prioritySecondApproverId := purchaseOrderRecord.GetString("priority_second_approver")
			if prioritySecondApproverId != "" {
				active, err := utilities.IsUserActive(app, prioritySecondApproverId)
				if err != nil {
					return &errs.HookError{
						Status:  http.StatusInternalServerError,
						Message: "failed to check priority second approver active status",
					}
				}
				if !active {
					return validation.NewError("priority_second_approver_not_active", "The selected priority second approver is not an active user")
				}
			}
			return nil
		}
	}

	validationsErrors := validation.Errors{
		"uid": validation.Validate(
			purchaseOrderRecord.GetString("uid"),
			validation.By(userIsActive(app, "uid_not_active", "the selected staff member is not an active user")),
		),
		"date": validation.Validate(
			purchaseOrderRecord.Get("date"),
			validation.Required.Error("date is required"),
			validation.Date("2006-01-02").Error("must be a valid date"),
		),
		"end_date": validation.Validate(
			purchaseOrderRecord.Get("end_date"),
			validation.When(isRecurring,
				validation.Required.Error("end_date is required for recurring purchase orders"),
				validation.Date("2006-01-02").Error("must be a valid date").Min(dateAsTime).RangeError("end date must be after start date").Max(dateAsTime.AddDate(0, 0, constants.RECURRING_MAX_DAYS)).RangeError(fmt.Sprintf("end date must be within %v days of the start date", constants.RECURRING_MAX_DAYS)),
			).Else(
				validation.In("").Error("end_date is not permitted for non-recurring purchase orders"),
			),
		),
		"priority_second_approver": validation.Validate(purchaseOrderRecord.Get("priority_second_approver"), validation.By(prioritySecondApproverIsActive(app, purchaseOrderRecord))),
		"frequency": validation.Validate(
			purchaseOrderRecord.Get("frequency"),
			validation.When(isRecurring,
				validation.Required.Error("frequency is required for recurring purchase orders"),
			).Else(
				validation.In("").Error("frequency is not permitted for non-recurring purchase orders"))),
		"description": validation.Validate(purchaseOrderRecord.Get("description"), validation.Length(4, 0).Error("must be at least 4 characters")),
		"approver": validation.Validate(purchaseOrderRecord.GetString("approver"),
			validation.By(userIsActive(app, "approver_not_active", "the selected approver is not an active user"))),
		"total": validation.Validate(purchaseOrderRecord.GetFloat("total"), validation.Max(constants.MAX_APPROVAL_TOTAL)),
		"type":  validation.Validate(purchaseOrderRecord.GetString("type"), validation.When(isChild, validation.In("One-Time").Error("child POs must be of type One-Time"))),
	}

	// Job-aware validation is performed in stages:
	// 1) job must exist and be Active
	// 2) PO branch must match job branch
	// 3) division must be allocated to the job, but only when creating a PO or
	//    when job/division is explicitly changed on update.
	//
	// Step (3) intentionally uses shouldValidateJobDivisionAllocationOnRecord so
	// historical rows can still be edited for unrelated fields without being
	// blocked by legacy allocation mismatches.
	if jobID := purchaseOrderRecord.GetString("job"); jobID != "" {
		jobRecord, err := app.FindRecordById("jobs", jobID)
		if err != nil || jobRecord == nil {
			validationsErrors["job"] = validation.NewError("invalid_reference", "invalid job reference")
		} else if jobRecord.GetString("status") != "Active" {
			validationsErrors["job"] = validation.NewError("not_active", "Job status must be Active")
		} else if purchaseOrderRecord.GetString("branch") != jobRecord.GetString("branch") {
			validationsErrors["branch"] = validation.NewError("value_mismatch", "branch must match the selected job")
		} else if shouldValidateJobDivisionAllocationOnRecord(app, purchaseOrderRecord) {
			if field, allocErr := validateDivisionAllocatedToJob(app, jobID, purchaseOrderRecord.GetString("division")); allocErr != nil {
				validationsErrors[field] = allocErr
			}
		}
		linkedJobRecord = jobRecord
	}

	if purchaseOrderRecord.GetBool("covered_within_project_budget") {
		switch {
		case strings.TrimSpace(purchaseOrderRecord.GetString("job")) == "":
			validationsErrors["covered_within_project_budget"] = validation.NewError(
				"project_job_required",
				"covered within project budget may only be confirmed for purchase orders linked to a project job",
			)
		case linkedJobRecord != nil && typeFromNumber(linkedJobRecord.GetString("number")) == jobTypeProposal:
			validationsErrors["covered_within_project_budget"] = validation.NewError(
				"project_job_required",
				"covered within project budget may only be confirmed for purchase orders linked to a project job",
			)
		}
	}

	// Policy computation depends on division and kind being valid. If either
	// already failed field-level validation, return those errors directly.
	if validationsErrors["division"] != nil || validationsErrors["kind"] != nil {
		return validationsErrors.Filter()
	}

	if !legacyMode {
		policy, err := utilities.GetPOApproverPolicy(
			app,
			purchaseOrderRecord.GetString("division"),
			utilities.EffectiveApprovalTotalHome(purchaseOrderRecord),
			purchaseOrderRecord.GetString("kind"),
			purchaseOrderRecord.GetString("job") != "",
		)
		if err != nil {
			if filteredErr := validationsErrors.Filter(); filteredErr != nil {
				return filteredErr
			}
			return &errs.HookError{
				Status:  http.StatusInternalServerError,
				Message: "hook error when computing purchase order approval policy",
				Data: map[string]errs.CodeError{
					"global": {
						Code:    "error_computing_approval_policy",
						Message: fmt.Sprintf("error computing purchase order approval policy: %v", err),
					},
				},
			}
		}

		setValidationError := func(field string, code string, message string) {
			if validationsErrors[field] == nil {
				validationsErrors[field] = validation.NewError(code, message)
			}
		}

		approverID := strings.TrimSpace(purchaseOrderRecord.GetString("approver"))
		prioritySecondApproverID := strings.TrimSpace(purchaseOrderRecord.GetString("priority_second_approver"))
		ownerID := strings.TrimSpace(purchaseOrderRecord.GetString("uid"))

		approverValidForStage := policy.IsFirstStageApprover(approverID) || policy.IsValidFirstStageViaBypass(approverID, ownerID)

		if !approverValidForStage {
			setValidationError("approver", "invalid_approver_for_stage", "selected approver is not valid for first-stage approval")
		}

		if policy.SecondApprovalRequired {
			if len(policy.FirstStageApprovers) == 0 && !approverValidForStage {
				setValidationError("approver", "first_pool_empty", "no configured first-stage approvers are available for this purchase order")
			}
			if len(policy.SecondStageApprovers) == 0 {
				setValidationError("priority_second_approver", "second_pool_empty", "no configured second-stage approvers are available for this purchase order")
			}
			if prioritySecondApproverID == "" {
				setValidationError("priority_second_approver", "priority_second_approver_required", "priority second approver is required when second approval is required")
			} else if !policy.IsSecondStageApprover(prioritySecondApproverID) {
				setValidationError("priority_second_approver", "invalid_priority_second_approver_for_stage", "selected priority second approver is not valid for second-stage approval")
			}
		} else if prioritySecondApproverID != "" {
			purchaseOrderRecord.Set("priority_second_approver", "")
		}
	}

	return validationsErrors.Filter()
}

var purchaseOrderMeaningfulChangeSkipFields = []string{
	"approved",
	"second_approval",
	"second_approver",
	"rejector",
	"rejected",
	"rejection_reason",
	"cancelled",
	"canceller",
	"closed",
	"closer",
	"closed_by_system",
	"po_number",
	"status",
}

func purchaseOrderHasMeaningfulChanges(record *core.Record) bool {
	return utilities.RecordHasMeaningfulChanges(record, purchaseOrderMeaningfulChangeSkipFields...)
}

// purchaseOrderRevisionEditableFields are the fields the creator may change on
// a revision. The rest of the purchase order was approved before and keeps its
// number, so it must not change. Computed fields are listed because
// CleanPurchaseOrder overwrites whatever the request sends for them. The update
// rule refuses changes to the skip fields. A new attachment is possible only on
// a revision that awaits approval, because the revise route takes no file.
// attachment_hash is not listed: this check runs before ProcessPurchaseOrder
// computes it from the uploaded file, so a submitted hash is always refused.
var purchaseOrderRevisionEditableFields = append([]string{
	"total",
	"currency",
	"approver",
	"priority_second_approver",
	"description",
	"division",
	"category",
	"attachment",
	"approval_total",
	"approval_total_home",
	"occurrences",
}, purchaseOrderMeaningfulChangeSkipFields...)

// purchaseOrderRevisionJobFields can change on a revision only together with
// its job. The budget confirmation was given for one job, and
// CleanPurchaseOrder sets the branch from the job.
var purchaseOrderRevisionJobFields = []string{
	"job",
	"covered_within_project_budget",
	"branch",
}

// revisionFieldLockedMessage says what can change on a revision that awaits
// approval. The update hook returns it with code revision_field_locked.
const revisionFieldLockedMessage = "only the total, currency, approvers, description, division, category, job and attachment can change on a revised purchase order, and the budget confirmation only with a new job"

// lockedRevisionFieldChange returns the first field outside
// purchaseOrderRevisionEditableFields that differs from the original record,
// or "" when only editable fields changed. When the job changes, the fields in
// purchaseOrderRevisionJobFields may change too.
func lockedRevisionFieldChange(record *core.Record) string {
	original := record.Original()
	jobChanged := strings.TrimSpace(record.GetString("job")) != strings.TrimSpace(original.GetString("job"))
	for _, fieldName := range record.Collection().Fields.FieldNames() {
		if fieldName == "created" || fieldName == "updated" || slices.Contains(purchaseOrderRevisionEditableFields, fieldName) {
			continue
		}
		if jobChanged && slices.Contains(purchaseOrderRevisionJobFields, fieldName) {
			continue
		}
		if fmt.Sprintf("%v", record.Get(fieldName)) != fmt.Sprintf("%v", original.Get(fieldName)) {
			return fieldName
		}
	}
	return ""
}

// RevisionJobChangeRefusal returns the code and message that refuse a
// revision's change of job from po's job to newJob, or two empty strings when
// the change is allowed. A revision can move a purchase order from one job to
// another. It can't add or remove a job, and a child purchase order must keep
// its parent's job.
func RevisionJobChangeRefusal(po *core.Record, newJob string) (string, string) {
	oldJob := strings.TrimSpace(po.GetString("job"))
	newJob = strings.TrimSpace(newJob)
	switch {
	case newJob == oldJob:
		return "", ""
	case strings.TrimSpace(po.GetString("parent_po")) != "":
		return "child_po_job_locked", "a child purchase order must keep the job of its parent"
	case oldJob == "":
		return "revision_job_add_not_allowed", "a revision cannot add a job to a purchase order that has none"
	case newJob == "":
		return "revision_job_removal_not_allowed", "a revision cannot remove the job"
	}
	return "", ""
}

// ValidatePurchaseOrderCategory checks that a purchase order's category belongs
// to its job, so a purchase order without a job can't have one. The collection
// API rules check this on a create or an edit, but the revise route doesn't go
// through them.
func ValidatePurchaseOrderCategory(app core.App, po *core.Record) error {
	categoryID := strings.TrimSpace(po.GetString("category"))
	if categoryID == "" {
		return nil
	}
	category, err := app.FindRecordById("categories", categoryID)
	if err != nil || category == nil {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating category",
			Data: map[string]errs.CodeError{
				"category": {Code: "not_found", Message: "category not found"},
			},
		}
	}
	if category.GetString("job") != po.GetString("job") {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating category",
			Data: map[string]errs.CodeError{
				"category": {Code: "must_match_job", Message: "category must belong to the selected job"},
			},
		}
	}
	return nil
}

// requestBodyHas reports whether the request that fired the hook sent field.
// It reports false when there is no request.
func requestBodyHas(e *core.RecordRequestEvent, field string) bool {
	if e.RequestEvent == nil || e.Request == nil {
		return false
	}
	info, err := e.RequestInfo()
	if err != nil {
		return false
	}
	_, ok := info.Body[field]
	return ok
}

func shouldResetPurchaseOrderApprovals(record *core.Record) bool {
	if record.IsNew() {
		return false
	}

	original := record.Original()
	if original == nil {
		return false
	}

	// Fully approved/terminal records are not reset via update-path editing.
	if original.GetString("status") != "Unapproved" || strings.TrimSpace(original.GetString("second_approval")) != "" {
		return false
	}
	// Suppress draft-update spam: only first-approved records reset + re-notify.
	if strings.TrimSpace(original.GetString("approved")) == "" {
		return false
	}
	// Resubmitting a rejected revision starts its approval again, even with
	// nothing changed, so that its first approver is asked and told.
	if utilities.IsPurchaseOrderRevision(original) && strings.TrimSpace(original.GetString("rejected")) != "" {
		return true
	}

	return purchaseOrderHasMeaningfulChanges(record)
}

// revisionChangeNeedsApproverEmail reports whether an update to a revision that
// has no first approval changes what its approver must do or was told: a
// resubmission after rejection, a new approver, or a new total, currency, job,
// description, division, category or attachment. A blank currency means the
// home currency, so filling it in changes nothing.
func revisionChangeNeedsApproverEmail(app core.App, record *core.Record) bool {
	original := record.Original()
	if record.IsNew() || original == nil || !utilities.IsPurchaseOrderRevision(original) {
		return false
	}
	if strings.TrimSpace(original.GetString("rejected")) != "" ||
		original.GetFloat("total") != record.GetFloat("total") ||
		utilities.EffectiveCurrencyCode(app, original.GetString("currency")) != utilities.EffectiveCurrencyCode(app, record.GetString("currency")) {
		return true
	}
	for _, field := range []string{"approver", "job", "description", "division", "category"} {
		if original.GetString(field) != record.GetString(field) {
			return true
		}
	}
	return PurchaseOrderAttachmentChange(record) != ""
}

// PurchaseOrderAttachmentChange says how an update changes a purchase order's
// attachment: "added", "replaced" or "removed", or "" when it keeps it. A new
// upload is a file, not a name, until the record is saved, so GetString can't
// see it.
func PurchaseOrderAttachmentChange(record *core.Record) string {
	original := record.Original()
	if original == nil {
		return ""
	}
	before := strings.TrimSpace(original.GetString("attachment"))
	uploaded := len(record.GetUnsavedFiles("attachment")) > 0
	after := strings.TrimSpace(record.GetString("attachment"))
	switch {
	case !uploaded && before == after:
		return ""
	case !uploaded && after == "":
		return "removed"
	case before == "":
		return "added"
	}
	return "replaced"
}

// PurchaseOrderApprovalTotalText is what the approver approves, with its
// currency. On a recurring purchase order, total is only one occurrence.
func PurchaseOrderApprovalTotalText(app core.App, po *core.Record) string {
	return utilities.FormatCurrencyAmount(po.GetFloat("approval_total"), utilities.EffectiveCurrencyCode(app, po.GetString("currency")))
}

// purchaseOrderJobText is a job's number and description, or "" when there is
// no job.
func purchaseOrderJobText(app core.App, jobID string) string {
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return ""
	}
	job, err := app.FindRecordById("jobs", jobID)
	if err != nil || job == nil {
		return jobID
	}
	return strings.TrimSpace(job.GetString("number") + " " + job.GetString("description"))
}

// purchaseOrderDivisionText is a division's code and name, or "" when there is
// no division.
func purchaseOrderDivisionText(app core.App, divisionID string) string {
	divisionID = strings.TrimSpace(divisionID)
	if divisionID == "" {
		return ""
	}
	division, err := app.FindRecordById("divisions", divisionID)
	if err != nil || division == nil {
		return divisionID
	}
	return strings.TrimSpace(division.GetString("code") + " " + division.GetString("name"))
}

// purchaseOrderCategoryText is a category's name, or "" when there is no
// category.
func purchaseOrderCategoryText(app core.App, categoryID string) string {
	categoryID = strings.TrimSpace(categoryID)
	if categoryID == "" {
		return ""
	}
	category, err := app.FindRecordById("categories", categoryID)
	if err != nil || category == nil {
		return categoryID
	}
	return category.GetString("name")
}

// revisionPreviousValue is the text of a value before a revision changed it, or
// "" when it didn't change. A value that was blank before shows as "none", so
// that the change still shows.
func revisionPreviousValue(before, after string) string {
	if before == after {
		return ""
	}
	if strings.TrimSpace(before) == "" {
		return "none"
	}
	return before
}

// RevisionApprovalNotification builds the email that asks a revision's approver
// to approve it again. po_approval_required says the purchase order was just
// created, which would mislead an approver about a number that already exists.
// previous is the purchase order before this change, and attachmentChange is
// what PurchaseOrderAttachmentChange returned for it. The email shows each
// previous value only if it differs, so that an approver can't miss a change:
// a new currency when the amount stays the same, a new job, division, category
// or description, or a new attachment.
func RevisionApprovalNotification(app core.App, po *core.Record, previous *core.Record, attachmentChange string, actorID string, mode notifications.DeliveryMode) notifications.DispatchArgs {
	vendorName := ""
	if vendor, err := app.FindRecordById("vendors", po.GetString("vendor")); err == nil {
		vendorName = vendor.GetString("name")
	}
	total := PurchaseOrderApprovalTotalText(app, po)
	job := purchaseOrderJobText(app, po.GetString("job"))
	division := purchaseOrderDivisionText(app, po.GetString("division"))
	category := purchaseOrderCategoryText(app, po.GetString("category"))
	description := po.GetString("description")
	if previous == nil {
		previous = po
	}
	data := map[string]any{
		"POId":                po.Id,
		"PONumber":            po.GetString("po_number"),
		"Description":         description,
		"PreviousDescription": revisionPreviousValue(previous.GetString("description"), description),
		"VendorName":          vendorName,
		"Job":                 job,
		"PreviousJob":         revisionPreviousValue(purchaseOrderJobText(app, previous.GetString("job")), job),
		"Division":            division,
		"PreviousDivision":    revisionPreviousValue(purchaseOrderDivisionText(app, previous.GetString("division")), division),
		"Category":            category,
		"PreviousCategory":    revisionPreviousValue(purchaseOrderCategoryText(app, previous.GetString("category")), category),
		"Total":               total,
		"PreviousTotal":       revisionPreviousValue(PurchaseOrderApprovalTotalText(app, previous), total),
		// Attachment is "added", "replaced" or "removed", or "" when it is kept.
		"Attachment": attachmentChange,
		"ActionURL":  notifications.BuildActionURL(app, fmt.Sprintf("/pos/%s/details", po.Id)),
	}
	return notifications.DispatchArgs{
		TemplateCode: "po_revision_approval_required",
		RecipientUID: strings.TrimSpace(po.GetString("approver")),
		Data:         data,
		System:       false,
		ActorUID:     actorID,
		Mode:         mode,
	}
}

// createPOApprovalRequiredNotification creates a notification record for the
// approver but does NOT send it immediately. It returns the notification ID so
// the caller can trigger the send after the PO save succeeds (avoiding
// premature email delivery if the save pipeline fails).
func createPOApprovalRequiredNotification(app core.App, purchaseOrderRecord *core.Record, actorID string) (string, error) {
	approverID := strings.TrimSpace(purchaseOrderRecord.GetString("approver"))
	if approverID == "" {
		return "", nil
	}
	if utilities.IsPurchaseOrderRevision(purchaseOrderRecord) {
		return notifications.DispatchNotification(app, RevisionApprovalNotification(app, purchaseOrderRecord, purchaseOrderRecord.Original(), PurchaseOrderAttachmentChange(purchaseOrderRecord), actorID, notifications.DeliveryDeferred))
	}

	return notifications.DispatchNotification(app, notifications.DispatchArgs{
		TemplateCode: "po_approval_required",
		RecipientUID: approverID,
		Data: map[string]any{
			"POId":      purchaseOrderRecord.Id,
			"ActionURL": notifications.BuildActionURL(app, fmt.Sprintf("/pos/%s/details", purchaseOrderRecord.Id)),
		},
		System:   false,
		ActorUID: actorID,
		Mode:     notifications.DeliveryDeferred,
	})
}

// ProcessPurchaseOrder validates and cleans a purchase_order record before
// create or update. It returns the notification ID of a pending approval
// notification (if one was created) so the caller can trigger the send after
// the save pipeline succeeds. The notification record is created here (so it
// can reference the PO ID), but actual email delivery must be deferred.
func ProcessPurchaseOrder(app core.App, e *core.RecordRequestEvent) (notificationID string, err error) {
	record := e.Record
	if err := checkExpensesEditing(app); err != nil {
		return "", err
	}

	// get the auth record from the context
	authRecord := e.Auth
	if authRecord == nil || strings.TrimSpace(authRecord.Id) == "" {
		return "", &errs.HookError{
			Status:  http.StatusUnauthorized,
			Message: "hook error when validating uid",
			Data: map[string]errs.CodeError{
				"uid": {
					Code:    "authentication_required",
					Message: "authentication is required",
				},
			},
		}
	}

	requestUID := strings.TrimSpace(record.GetString("uid"))
	if record.IsNew() {
		if requestUID != authRecord.Id {
			return "", &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating uid",
				Data: map[string]errs.CodeError{
					"uid": {
						Code:    "value_mismatch",
						Message: "uid must be equal to the authenticated user's id",
					},
				},
			}
		}
	} else {
		original := record.Original()
		if original == nil {
			return "", &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating uid",
				Data: map[string]errs.CodeError{
					"uid": {
						Code:    "missing_original_record",
						Message: "cannot validate ownership for this purchase order update",
					},
				},
			}
		}

		originalUID := strings.TrimSpace(original.GetString("uid"))
		if originalUID != authRecord.Id {
			return "", &errs.HookError{
				Status:  http.StatusForbidden,
				Message: "hook error when validating uid",
				Data: map[string]errs.CodeError{
					"uid": {
						Code:    "not_owner",
						Message: "only the creator can edit this purchase order",
					},
				},
			}
		}
		if requestUID != originalUID {
			return "", &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating uid",
				Data: map[string]errs.CodeError{
					"uid": {
						Code:    "immutable_field",
						Message: "uid cannot be changed",
					},
				},
			}
		}

		if original.GetString("status") != "Unapproved" {
			return "", &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "hook error when validating status",
				Data: map[string]errs.CodeError{
					"status": {
						Code:    "invalid_status",
						Message: "only unapproved purchase orders can be edited",
					},
				},
			}
		}

		// This check runs before CleanPurchaseOrder so it compares the submitted
		// values, not the cleaned ones.
		if utilities.IsPurchaseOrderRevision(original) {
			if field := lockedRevisionFieldChange(record); field != "" {
				return "", &errs.HookError{
					Status:  http.StatusBadRequest,
					Message: "hook error when validating revised purchase order",
					Data: map[string]errs.CodeError{
						field: {
							Code:    "revision_field_locked",
							Message: revisionFieldLockedMessage,
						},
					},
				}
			}
			if code, message := RevisionJobChangeRefusal(original, record.GetString("job")); code != "" {
				return "", &errs.HookError{
					Status:  http.StatusBadRequest,
					Message: "hook error when validating revised purchase order",
					Data: map[string]errs.CodeError{
						"job": {Code: code, Message: message},
					},
				}
			}
			// A category and a budget confirmation belong to the old job. Keep
			// them only if the request sets them again for the new job.
			jobChanged := record.GetString("job") != original.GetString("job")
			if jobChanged {
				if !requestBodyHas(e, "category") {
					record.Set("category", "")
				}
				if !requestBodyHas(e, "covered_within_project_budget") {
					record.Set("covered_within_project_budget", false)
				}
			}
			// The update rule checks the category too. This gives the same
			// error as the revise route.
			if jobChanged || record.GetString("category") != original.GetString("category") {
				if err := ValidatePurchaseOrderCategory(app, record); err != nil {
					return "", err
				}
			}
		}
	}

	// Only the purchase order groups route assigns po_group. A collection
	// create can't join a group and an update can't move a PO between groups.
	requestGroup := strings.TrimSpace(record.GetString("po_group"))
	originalGroup := ""
	if original := record.Original(); !record.IsNew() && original != nil {
		originalGroup = strings.TrimSpace(original.GetString("po_group"))
	}
	if requestGroup != originalGroup {
		return "", &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating po_group",
			Data: map[string]errs.CodeError{
				"po_group": {
					Code:    "immutable_field",
					Message: "po_group can only be set when a purchase order group is created",
				},
			},
		}
	}

	// Purchase orders must remain Unapproved through collection create/update
	// requests. Transitioning to Active/Cancelled/Closed is only allowed via the
	// dedicated action endpoints.
	if record.GetString("status") != "Unapproved" {
		return "", &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "hook error when validating status",
			Data: map[string]errs.CodeError{
				"status": {
					Code:    "invalid_status",
					Message: "status must be Unapproved when creating or updating purchase orders",
				},
			},
		}
	}

	shouldResetApprovals := shouldResetPurchaseOrderApprovals(record)
	submittedApproverID := strings.TrimSpace(record.GetString("approver"))

	// set properties to nil if they are not allowed to be set based on the type
	cleanErr := CleanPurchaseOrder(app, record)
	if cleanErr != nil {
		return "", cleanErr
	}

	if err := utilities.EnsureUserCanUseRecordBranch(app, record, authRecord.Id); err != nil {
		return "", err
	}

	if shouldResetApprovals {
		// Clear approval state for meaningful edits to first-approved in-flight POs.
		record.Set("approved", "")
		record.Set("approver", "")
		record.Set("second_approval", "")
		record.Set("second_approver", "")
		// Keep any submitted assignee as the next first-stage approver.
		record.Set("approver", submittedApproverID)
	}

	NormalizePendingFileNames(record, "attachment")

	// if the purchase order has a new attachment, calculate the sha256 hash of the
	// file and set the attachment_hash property on the record
	attachmentHash, hashErr := CalculateFileFieldHash(e, "attachment")
	if hashErr != nil {
		return "", hashErr
	}
	if attachmentHash != "" {
		// New file uploaded - check for duplicate before setting hash
		// Look for existing purchase order with the same attachment hash (excluding this record)
		existingPO, _ := app.FindFirstRecordByFilter("purchase_orders", "attachment_hash = {:hash} && id != {:id}", dbx.Params{
			"hash": attachmentHash,
			"id":   record.Id,
		})
		if existingPO != nil {
			return "", &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "duplicate attachment detected",
				Data: map[string]errs.CodeError{
					"attachment": {
						Code:    "duplicate_file",
						Message: "This file has already been uploaded to another purchase order",
					},
				},
			}
		}

		record.Set("attachment_hash", attachmentHash)
	} else if record.GetString("attachment") == "" {
		// Attachment was explicitly removed - clear the hash
		record.Set("attachment_hash", "")
	}
	// Otherwise: no new file and attachment still exists - leave hash unchanged

	if err := EnsureActiveDivision(app, record.GetString("division"), "division"); err != nil {
		if ve, ok := err.(validation.Errors); ok {
			return "", apis.NewBadRequestError("Validation error", ve)
		}
		return "", err
	}

	// validate the purchase_order record
	if validationErr := ValidatePurchaseOrder(app, record, false); validationErr != nil {
		return "", validationErr
	}

	// If this is a new record, create a notification for the approver.
	// The notification record is created now (so it can reference the PO ID),
	// but the caller is responsible for triggering the send after the save
	// pipeline completes successfully.
	if record.IsNew() {

		// Generate a new id for the record here so that the notification can
		// reference it
		// https://github.com/pocketbase/pocketbase/discussions/6170
		// https://pocketbase.io/docs/collections/#textfield
		record.Set("id:autogenerate", "")

		nid, err := createPOApprovalRequiredNotification(app, record, authRecord.Id)
		if err != nil {
			return "", err
		}
		return nid, nil
	}
	if shouldResetApprovals || revisionChangeNeedsApproverEmail(app, record) {
		nid, err := createPOApprovalRequiredNotification(app, record, authRecord.Id)
		if err != nil {
			return "", err
		}
		return nid, nil
	}

	return "", nil
}
