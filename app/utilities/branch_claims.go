package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
	"tybalt/errs"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// EnsureUserCanUseRecordBranch applies EnsureUserCanUseBranch to a record's
// branch field unless the user did not choose that branch. Branch claims then
// do not apply:
//
//   - A job forces the record's branch to the job's branch. Charging time or
//     costs to a job on a restricted branch (for example a company-wide job on
//     the Corporate branch) must not require the branch claim. The branch is
//     compared with the job's rather than skipping whenever a job is set
//     because expenses linked to a purchase order take the purchase order's
//     branch, which can differ from the job's branch.
//   - The creator of a purchase order group chooses the branch of every
//     purchase order in it, and the group route checks it against the
//     creator's claims. The owner keeps that branch when editing the purchase
//     order, and expenses against it take it from the purchase order.
func EnsureUserCanUseRecordBranch(app core.App, record *core.Record, userID string) error {
	branchID := strings.TrimSpace(record.GetString("branch"))
	if branchID != "" && branchIsInherited(app, record, branchID) {
		return nil
	}
	return EnsureUserCanUseBranch(app, branchID, userID, "branch")
}

func branchIsInherited(app core.App, record *core.Record, branchID string) bool {
	if jobID := strings.TrimSpace(record.GetString("job")); jobID != "" {
		job, err := app.FindRecordById("jobs", jobID)
		if err == nil && job != nil && strings.TrimSpace(job.GetString("branch")) == branchID {
			return true
		}
	}
	// A purchase order in a group that still has the branch it was created with.
	if strings.TrimSpace(record.GetString("po_group")) != "" && !record.IsNew() {
		if original := record.Original(); original != nil && strings.TrimSpace(original.GetString("branch")) == branchID {
			return true
		}
	}
	// An expense against a purchase order in a group.
	if poID := strings.TrimSpace(record.GetString("purchase_order")); poID != "" {
		po, err := app.FindRecordById("purchase_orders", poID)
		if err == nil && po != nil && strings.TrimSpace(po.GetString("po_group")) != "" &&
			strings.TrimSpace(po.GetString("branch")) == branchID {
			return true
		}
	}
	return false
}

// EnsureUserCanUseBranch enforces branch-level claim restrictions for the
// target user. A branch with no allowed_claims is unrestricted.
func EnsureUserCanUseBranch(app core.App, branchID string, userID string, fieldName string) error {
	branchID = strings.TrimSpace(branchID)
	if branchID == "" {
		return nil
	}

	branch, err := app.FindRecordById("branches", branchID)
	if err != nil || branch == nil {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "branch lookup failed",
			Data: map[string]errs.CodeError{
				fieldName: {
					Code:    "invalid_branch",
					Message: "specified branch could not be found",
				},
			},
		}
	}

	allowedClaimIDs, err := normalizeStringList(branch.Get("allowed_claims"))
	if err != nil {
		return &errs.HookError{
			Status:  http.StatusInternalServerError,
			Message: "branch claim validation failed",
			Data: map[string]errs.CodeError{
				fieldName: {
					Code:    "branch_claim_config_invalid",
					Message: "branch claim requirements are misconfigured",
				},
			},
		}
	}
	if len(allowedClaimIDs) == 0 {
		return nil
	}

	userClaimRecords, err := app.FindRecordsByFilter(
		"user_claims",
		"uid={:uid}",
		"",
		0,
		0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return &errs.HookError{
			Status:  http.StatusInternalServerError,
			Message: "branch claim validation failed",
			Data: map[string]errs.CodeError{
				fieldName: {
					Code:    "branch_claim_check_failed",
					Message: "unable to validate branch claim requirements",
				},
			},
		}
	}

	for _, record := range userClaimRecords {
		if slices.Contains(allowedClaimIDs, strings.TrimSpace(record.GetString("cid"))) {
			return nil
		}
	}

	requiredClaims := make([]string, 0, len(allowedClaimIDs))
	for _, claimID := range allowedClaimIDs {
		claim, claimErr := app.FindRecordById("claims", claimID)
		if claimErr == nil && claim != nil && strings.TrimSpace(claim.GetString("name")) != "" {
			requiredClaims = append(requiredClaims, strings.TrimSpace(claim.GetString("name")))
			continue
		}
		requiredClaims = append(requiredClaims, claimID)
	}
	sort.Strings(requiredClaims)

	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: "branch claim requirement not met",
		Data: map[string]errs.CodeError{
			fieldName: {
				Code: "branch_claim_required",
				Message: fmt.Sprintf(
					"%s requires at least one of these claims: %s",
					branch.GetString("name"),
					strings.Join(requiredClaims, ", "),
				),
			},
		},
	}
}

func normalizeStringList(value any) ([]string, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case []string:
		return slices.DeleteFunc(slices.Clone(v), func(item string) bool {
			return strings.TrimSpace(item) == ""
		}), nil
	case string:
		if strings.TrimSpace(v) == "" {
			return nil, nil
		}
		var result []string
		if err := json.Unmarshal([]byte(v), &result); err == nil {
			return slices.DeleteFunc(result, func(item string) bool {
				return strings.TrimSpace(item) == ""
			}), nil
		}
		return nil, fmt.Errorf("unexpected string value for relation list")
	default:
		return nil, fmt.Errorf("unexpected relation list type %T", value)
	}
}
