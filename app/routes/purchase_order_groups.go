// This file implements purchase order groups (issue #168). A super_po holder
// fills out one purchase order as a template and uploads a CSV with one row per
// staff member. The preview route reports what each row would create, and the
// create route makes one ordinary Unapproved purchase order per row, owned by
// that staff member, all linked to one purchase_order_groups record. Each
// purchase order then goes through normal approval, which assigns its number.

package routes

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/notifications"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

const maxPurchaseOrderGroupRows = 1000

// purchaseOrderGroupUploadFields lists the template fields a CSV column may
// replace. The email column is always required and picks the owner.
var purchaseOrderGroupUploadFields = []string{"total", "description", "division", "branch"}

var purchaseOrderGroupIDPattern = regexp.MustCompile(`^[a-z0-9]{15}$`)

// Friendlier messages for validation codes whose standard wording assumes the
// signed-in user is creating their own purchase order.
var purchaseOrderGroupMessages = map[string]string{
	"invalid_approver_for_stage":                 "The approver can't approve this amount for this kind and division.",
	"priority_second_approver_required":          "This amount needs a second approver. Choose one in the template.",
	"invalid_priority_second_approver_for_stage": "The second approver can't give final approval for this amount.",
}

type purchaseOrderGroupTemplate struct {
	Type                   string  `json:"type"`
	Kind                   string  `json:"kind"`
	PaymentType            string  `json:"payment_type"`
	Date                   string  `json:"date"`
	Division               string  `json:"division"`
	Branch                 string  `json:"branch"`
	Description            string  `json:"description"`
	Total                  float64 `json:"total"`
	Currency               string  `json:"currency"`
	Vendor                 string  `json:"vendor"`
	Approver               string  `json:"approver"`
	PrioritySecondApprover string  `json:"priority_second_approver"`
}

type purchaseOrderGroupRequest struct {
	GroupID      string                     `json:"group_id"`
	Template     purchaseOrderGroupTemplate `json:"template"`
	UploadFields []string                   `json:"upload_fields"`
	Rows         []map[string]string        `json:"rows"`
}

type purchaseOrderGroupStaff struct {
	UID    string `db:"uid" json:"uid"`
	Email  string `db:"email" json:"email"`
	Name   string `db:"name" json:"name"`
	Active bool   `db:"active" json:"-"`
}

type purchaseOrderGroupRow struct {
	Index        int               `json:"index"`
	Email        string            `json:"email"`
	UID          string            `json:"uid"`
	Name         string            `json:"name"`
	Total        float64           `json:"total"`
	Description  string            `json:"description"`
	DivisionCode string            `json:"division_code"`
	BranchName   string            `json:"branch_name"`
	Errors       map[string]string `json:"errors,omitempty"`
}

type purchaseOrderGroupPreview struct {
	Rows               []purchaseOrderGroupRow   `json:"rows"`
	Count              int                       `json:"count"`
	ErrorCount         int                       `json:"error_count"`
	Total              float64                   `json:"total"`
	CurrencyCode       string                    `json:"currency_code"`
	MissingActiveStaff []purchaseOrderGroupStaff `json:"missing_active_staff"`
	// The people who could approve every row, for the page to offer. A
	// second approver is offered only when some row needs one.
	Approvers              []utilities.Approver `json:"approvers"`
	SecondApprovers        []utilities.Approver `json:"second_approvers"`
	SecondApprovalRequired bool                 `json:"second_approval_required"`
}

// purchaseOrderGroupApproverPools narrows the approvers to those who could
// approve every row, using the same policy as ValidatePurchaseOrder. A row's
// policy depends on its division, amount and kind, so rows that share those
// share a policy. The owner self-approval bypass is left out because the
// owner differs from row to row.
type purchaseOrderGroupApproverPools struct {
	policies       map[string]utilities.POApproverPolicy
	first          map[string]utilities.Approver
	second         map[string]utilities.Approver
	secondRequired bool
}

func intersectPurchaseOrderGroupApprovers(current map[string]utilities.Approver, approvers []utilities.Approver) map[string]utilities.Approver {
	next := map[string]utilities.Approver{}
	for _, approver := range approvers {
		if _, ok := current[approver.ID]; current == nil || ok {
			next[approver.ID] = approver
		}
	}
	return next
}

// add narrows the pools by one cleaned purchase order.
func (p *purchaseOrderGroupApproverPools) add(app core.App, record *core.Record) error {
	division := record.GetString("division")
	amount := utilities.EffectiveApprovalTotalHome(record)
	kind := record.GetString("kind")
	key := division + "|" + strconv.FormatFloat(amount, 'f', -1, 64) + "|" + kind
	policy, ok := p.policies[key]
	if !ok {
		var err error
		policy, err = utilities.GetPOApproverPolicy(app, division, amount, kind, false)
		if err != nil {
			return err
		}
		p.policies[key] = policy
	}
	p.first = intersectPurchaseOrderGroupApprovers(p.first, policy.FirstStageApprovers)
	if policy.SecondApprovalRequired {
		p.secondRequired = true
		p.second = intersectPurchaseOrderGroupApprovers(p.second, policy.SecondStageApprovers)
	}
	return nil
}

func sortedPurchaseOrderGroupApprovers(pool map[string]utilities.Approver) []utilities.Approver {
	approvers := make([]utilities.Approver, 0, len(pool))
	for _, approver := range pool {
		approvers = append(approvers, approver)
	}
	slices.SortFunc(approvers, func(a, b utilities.Approver) int {
		if c := strings.Compare(a.Surname, b.Surname); c != 0 {
			return c
		}
		if c := strings.Compare(a.GivenName, b.GivenName); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return approvers
}

// purchaseOrderGroupRowsError aborts a create whose rows fail validation. The
// handler returns the preview so the UI can show which rows to fix.
type purchaseOrderGroupRowsError struct {
	preview purchaseOrderGroupPreview
}

func (e *purchaseOrderGroupRowsError) Error() string {
	return fmt.Sprintf("%d rows have errors", e.preview.ErrorCount)
}

type purchaseOrderGroupLookups struct {
	staff          []purchaseOrderGroupStaff
	staffByEmail   map[string]purchaseOrderGroupStaff
	divisionByCode map[string]string
	divisionCode   map[string]string
	branchByKey    map[string]string
	branchName     map[string]string
}

func loadPurchaseOrderGroupLookups(app core.App) (*purchaseOrderGroupLookups, error) {
	lookups := &purchaseOrderGroupLookups{
		staffByEmail:   map[string]purchaseOrderGroupStaff{},
		divisionByCode: map[string]string{},
		divisionCode:   map[string]string{},
		branchByKey:    map[string]string{},
		branchName:     map[string]string{},
	}

	if err := app.DB().NewQuery(`
		SELECT
			u.id AS uid,
			u.email AS email,
			TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')) AS name,
			COALESCE(ap.active, 0) AS active
		FROM users u
		LEFT JOIN profiles p ON p.uid = u.id
		LEFT JOIN admin_profiles ap ON ap.uid = u.id
		ORDER BY p.surname, p.given_name
	`).All(&lookups.staff); err != nil {
		return nil, err
	}
	for _, staff := range lookups.staff {
		lookups.staffByEmail[strings.ToLower(strings.TrimSpace(staff.Email))] = staff
	}

	divisions := []struct {
		ID   string `db:"id"`
		Code string `db:"code"`
	}{}
	if err := app.DB().NewQuery(`SELECT id, code FROM divisions`).All(&divisions); err != nil {
		return nil, err
	}
	for _, division := range divisions {
		lookups.divisionByCode[strings.ToLower(division.Code)] = division.ID
		lookups.divisionCode[division.ID] = division.Code
	}

	branches := []struct {
		ID   string `db:"id"`
		Code string `db:"code"`
		Name string `db:"name"`
	}{}
	if err := app.DB().NewQuery(`SELECT id, code, name FROM branches`).All(&branches); err != nil {
		return nil, err
	}
	for _, branch := range branches {
		lookups.branchByKey[strings.ToLower(branch.Name)] = branch.ID
		lookups.branchByKey[strings.ToLower(branch.Code)] = branch.ID
		lookups.branchName[branch.ID] = branch.Name
	}

	return lookups, nil
}

func requirePurchaseOrderGroupAccess(app core.App, authRecord *core.Record) error {
	hasClaim, err := utilities.HasClaim(app, authRecord, "super_po")
	if err != nil {
		return &errs.HookError{
			Status:  http.StatusInternalServerError,
			Message: "unable to check purchase order group permissions",
			Data: map[string]errs.CodeError{
				"global": {Code: "claim_check_failed", Message: "unable to check purchase order group permissions"},
			},
		}
	}
	if !hasClaim {
		return &errs.HookError{
			Status:  http.StatusForbidden,
			Message: "creating purchase orders for other staff requires the super_po claim",
			Data: map[string]errs.CodeError{
				"global": {Code: "super_po_claim_required", Message: "you do not have permission to create purchase orders for other staff"},
			},
		}
	}
	return nil
}

// requirePurchaseOrderGroupViewAccess lets admin holders look at groups as
// well as the super_po holders who create them.
func requirePurchaseOrderGroupViewAccess(app core.App, authRecord *core.Record) error {
	isAdmin, err := utilities.HasClaim(app, authRecord, "admin")
	if err == nil && isAdmin {
		return nil
	}
	return requirePurchaseOrderGroupAccess(app, authRecord)
}

// normalizePurchaseOrderGroupRequest lowercases CSV column names and checks the
// parts of the request that apply to every row.
func normalizePurchaseOrderGroupRequest(req *purchaseOrderGroupRequest, creating bool) error {
	fieldErrors := map[string]errs.CodeError{}
	add := func(field, code, message string) {
		if _, exists := fieldErrors[field]; !exists {
			fieldErrors[field] = errs.CodeError{Code: code, Message: message}
		}
	}

	if creating && !purchaseOrderGroupIDPattern.MatchString(req.GroupID) {
		add("group_id", "invalid_group_id", "group_id must be 15 lowercase letters or digits")
	}

	uploaded := map[string]bool{}
	for _, field := range req.UploadFields {
		if !slices.Contains(purchaseOrderGroupUploadFields, field) {
			add("upload_fields", "invalid_upload_field", fmt.Sprintf("%q can't come from the upload", field))
			continue
		}
		uploaded[field] = true
	}
	req.UploadFields = req.UploadFields[:0]
	for _, field := range purchaseOrderGroupUploadFields {
		if uploaded[field] {
			req.UploadFields = append(req.UploadFields, field)
		}
	}

	t := &req.Template
	if t.Type != "One-Time" && t.Type != "Cumulative" {
		add("type", "invalid_type", "Type must be One-Time or Cumulative.")
	}
	requireValue := func(field, value, message string) {
		if strings.TrimSpace(value) == "" {
			add(field, "required", message)
		}
	}
	requireValue("kind", t.Kind, "Kind is required.")
	requireValue("payment_type", t.PaymentType, "Payment is required.")
	requireValue("date", t.Date, "Date is required.")
	requireValue("vendor", t.Vendor, "Vendor is required.")
	// The page chooses the approver from the preview's list of people who can
	// approve every row, so only creating needs one.
	if creating {
		requireValue("approver", t.Approver, "Approver is required.")
	}
	if !uploaded["branch"] {
		requireValue("branch", t.Branch, "Branch is required.")
	}
	if !uploaded["division"] {
		requireValue("division", t.Division, "Division is required.")
	}
	if !uploaded["description"] {
		requireValue("description", t.Description, "Description is required.")
	}
	if !uploaded["total"] && t.Total <= 0 {
		add("total", "required", "Total must be greater than 0.")
	}

	if len(req.Rows) == 0 {
		add("rows", "required", "The upload has no rows.")
	} else if len(req.Rows) > maxPurchaseOrderGroupRows {
		add("rows", "too_many_rows", fmt.Sprintf("The upload has %d rows. The limit is %d.", len(req.Rows), maxPurchaseOrderGroupRows))
	}
	for i, row := range req.Rows {
		normalized := make(map[string]string, len(row))
		for key, value := range row {
			normalized[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
		}
		req.Rows[i] = normalized
	}

	if len(fieldErrors) > 0 {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "purchase order group template is invalid",
			Data:    fieldErrors,
		}
	}
	return nil
}

// addPurchaseOrderGroupErrors copies field errors from err into rowErrors. It
// returns err itself when the failure is a server error rather than a problem
// with the row.
func addPurchaseOrderGroupErrors(rowErrors map[string]string, err error) error {
	if err == nil {
		return nil
	}
	set := func(field, code, message string) {
		if friendly, ok := purchaseOrderGroupMessages[code]; ok {
			message = friendly
		}
		if _, exists := rowErrors[field]; !exists {
			rowErrors[field] = message
		}
	}

	var hookErr *errs.HookError
	var validationErrs validation.Errors
	switch {
	case errors.As(err, &hookErr):
		if hookErr.Status >= http.StatusInternalServerError {
			return err
		}
		for field, codeErr := range hookErr.Data {
			set(field, codeErr.Code, codeErr.Message)
		}
		if len(hookErr.Data) == 0 {
			set("global", "", hookErr.Message)
		}
	case errors.As(err, &validationErrs):
		for field, fieldErr := range validationErrs {
			var nestedHookErr *errs.HookError
			if errors.As(fieldErr, &nestedHookErr) && nestedHookErr.Status >= http.StatusInternalServerError {
				return nestedHookErr
			}
			code := ""
			var validationErr validation.Error
			if errors.As(fieldErr, &validationErr) {
				code = validationErr.Code()
			}
			set(field, code, fieldErr.Error())
		}
	default:
		return err
	}
	return nil
}

func parsePurchaseOrderGroupTotal(value string) (float64, error) {
	cleaned := strings.NewReplacer("$", "", ",", "", " ", "").Replace(value)
	if cleaned == "" {
		return 0, errors.New("empty")
	}
	return strconv.ParseFloat(cleaned, 64)
}

// buildPurchaseOrderGroupRow makes the purchase order for one CSV row and runs
// the same cleaning and validation as a normal create. The creator chooses the
// branch, so it is checked against the creator's branch claims rather than the
// owner's, and is never filled from the owner's default. Each row that gets
// as far as cleaning narrows pools. Until an approver is chosen, rows aren't
// faulted for the approvers. The record is returned only when the row has no
// errors.
func buildPurchaseOrderGroupRow(
	app core.App,
	collection *core.Collection,
	lookups *purchaseOrderGroupLookups,
	req *purchaseOrderGroupRequest,
	creatorID string,
	index int,
	emailCounts map[string]int,
	pools *purchaseOrderGroupApproverPools,
) (purchaseOrderGroupRow, *core.Record, error) {
	row := req.Rows[index]
	t := req.Template
	result := purchaseOrderGroupRow{Index: index, Errors: map[string]string{}}

	email := strings.ToLower(row["email"])
	result.Email = email
	staff, found := lookups.staffByEmail[email]
	switch {
	case email == "":
		result.Errors["email"] = "Email is required."
	case emailCounts[email] > 1:
		result.Errors["email"] = fmt.Sprintf("This email appears %d times in the upload.", emailCounts[email])
	case !found:
		result.Errors["email"] = "No staff member has this email."
	case !staff.Active:
		result.Errors["email"] = "This staff member is not active."
	}
	if found {
		result.UID = staff.UID
		result.Name = staff.Name
	}

	total := t.Total
	description := t.Description
	division := t.Division
	branch := t.Branch
	for _, field := range req.UploadFields {
		value := row[field]
		switch field {
		case "total":
			parsed, err := parsePurchaseOrderGroupTotal(value)
			switch {
			case value == "":
				result.Errors["total"] = "Total is required."
			case err != nil:
				result.Errors["total"] = fmt.Sprintf("Total %q is not a number.", value)
			}
			total = parsed
		case "description":
			if value == "" {
				result.Errors["description"] = "Description is required."
			}
			description = value
		case "division":
			id, ok := lookups.divisionByCode[strings.ToLower(value)]
			switch {
			case value == "":
				result.Errors["division"] = "Division is required."
			case !ok:
				result.Errors["division"] = fmt.Sprintf("No division has code %q.", value)
			}
			division = id
		case "branch":
			id, ok := lookups.branchByKey[strings.ToLower(value)]
			switch {
			case value == "":
				result.Errors["branch"] = "Branch is required."
			case !ok:
				result.Errors["branch"] = fmt.Sprintf("No branch has the name or code %q.", value)
			}
			branch = id
		}
	}
	result.Total = total
	result.Description = description
	result.DivisionCode = lookups.divisionCode[division]
	result.BranchName = lookups.branchName[branch]
	if len(result.Errors) > 0 {
		return result, nil, nil
	}

	record := core.NewRecord(collection)
	record.Set("uid", staff.UID)
	record.Set("status", "Unapproved")
	record.Set("type", t.Type)
	record.Set("kind", t.Kind)
	record.Set("payment_type", t.PaymentType)
	record.Set("date", t.Date)
	record.Set("division", division)
	record.Set("branch", branch)
	record.Set("description", description)
	record.Set("total", total)
	record.Set("currency", t.Currency)
	record.Set("vendor", t.Vendor)
	record.Set("approver", t.Approver)
	record.Set("priority_second_approver", t.PrioritySecondApprover)

	if err := addPurchaseOrderGroupErrors(result.Errors, hooks.CleanPurchaseOrder(app, record)); err != nil {
		return result, nil, err
	}
	if len(result.Errors) == 0 {
		if err := pools.add(app, record); err != nil {
			return result, nil, err
		}
		checks := []error{
			utilities.EnsureUserCanUseRecordBranch(app, record, creatorID),
			hooks.EnsureActiveDivision(app, record.GetString("division"), "division"),
			hooks.ValidatePurchaseOrder(app, record, false),
		}
		for _, checkErr := range checks {
			if err := addPurchaseOrderGroupErrors(result.Errors, checkErr); err != nil {
				return result, nil, err
			}
		}
	}
	if strings.TrimSpace(t.Approver) == "" {
		delete(result.Errors, "approver")
		delete(result.Errors, "priority_second_approver")
	}

	if len(result.Errors) > 0 {
		return result, nil, nil
	}
	return result, record, nil
}

// evaluatePurchaseOrderGroup builds every row's purchase order. It returns the
// preview and, when no row has errors, the unsaved records in row order. It
// stops early when ctx ends, such as when the page sends a newer preview.
func evaluatePurchaseOrderGroup(ctx context.Context, app core.App, req *purchaseOrderGroupRequest, creatorID string) (purchaseOrderGroupPreview, []*core.Record, error) {
	preview := purchaseOrderGroupPreview{
		Rows:               []purchaseOrderGroupRow{},
		MissingActiveStaff: []purchaseOrderGroupStaff{},
		CurrencyCode:       utilities.EffectiveCurrencyCode(app, req.Template.Currency),
	}

	collection, err := app.FindCollectionByNameOrId("purchase_orders")
	if err != nil {
		return preview, nil, err
	}
	lookups, err := loadPurchaseOrderGroupLookups(app)
	if err != nil {
		return preview, nil, err
	}

	emailCounts := map[string]int{}
	for _, row := range req.Rows {
		if email := strings.ToLower(row["email"]); email != "" {
			emailCounts[email]++
		}
	}

	records := []*core.Record{}
	listed := map[string]bool{}
	pools := &purchaseOrderGroupApproverPools{policies: map[string]utilities.POApproverPolicy{}}
	for i := range req.Rows {
		if err := ctx.Err(); err != nil {
			return preview, nil, err
		}
		row, record, err := buildPurchaseOrderGroupRow(app, collection, lookups, req, creatorID, i, emailCounts, pools)
		if err != nil {
			return preview, nil, err
		}
		if len(row.Errors) > 0 {
			preview.ErrorCount++
		} else {
			records = append(records, record)
		}
		if row.UID != "" {
			listed[row.UID] = true
		}
		preview.Total += row.Total
		preview.Rows = append(preview.Rows, row)
	}
	preview.Count = len(preview.Rows)
	preview.Approvers = sortedPurchaseOrderGroupApprovers(pools.first)
	preview.SecondApprovers = sortedPurchaseOrderGroupApprovers(pools.second)
	preview.SecondApprovalRequired = pools.secondRequired

	for _, staff := range lookups.staff {
		if staff.Active && !listed[staff.UID] {
			preview.MissingActiveStaff = append(preview.MissingActiveStaff, staff)
		}
	}

	if preview.ErrorCount > 0 {
		return preview, nil, nil
	}
	return preview, records, nil
}

func bindPurchaseOrderGroupRequest(e *core.RequestEvent, creating bool) (*purchaseOrderGroupRequest, error) {
	req := &purchaseOrderGroupRequest{}
	if err := e.BindBody(req); err != nil {
		return nil, &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Data: map[string]errs.CodeError{
				"global": {Code: "invalid_request_body", Message: "invalid request body"},
			},
		}
	}
	if err := normalizePurchaseOrderGroupRequest(req, creating); err != nil {
		return nil, err
	}
	return req, nil
}

func respondPurchaseOrderGroupError(e *core.RequestEvent, err error) error {
	var rowsErr *purchaseOrderGroupRowsError
	if errors.As(err, &rowsErr) {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"code":    "rows_invalid",
			"message": rowsErr.Error(),
			"preview": rowsErr.preview,
		})
	}
	var hookErr *errs.HookError
	if errors.As(err, &hookErr) {
		return e.JSON(hookErr.Status, hookErr)
	}
	var codeErr *CodeError
	if errors.As(err, &codeErr) {
		return e.JSON(http.StatusConflict, codeErr)
	}
	return e.JSON(http.StatusInternalServerError, map[string]any{
		"code":    "purchase_order_group_error",
		"message": err.Error(),
	})
}

func createPreviewPurchaseOrderGroupHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requirePurchaseOrderGroupAccess(app, e.Auth); err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		req, err := bindPurchaseOrderGroupRequest(e, false)
		if err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		preview, _, err := evaluatePurchaseOrderGroup(e.Request.Context(), app, req, e.Auth.Id)
		if err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		return e.JSON(http.StatusOK, preview)
	}
}

// purchaseOrderGroupAmount formats a group total for the approver's email, for
// example "12,345.60 CAD".
func purchaseOrderGroupAmount(value float64, currencyCode string) string {
	cents := int64(math.Round(value * 100))
	whole := strconv.FormatInt(cents/100, 10)
	for i := len(whole) - 3; i > 0; i -= 3 {
		whole = whole[:i] + "," + whole[i:]
	}
	return fmt.Sprintf("%s.%02d %s", whole, cents%100, currencyCode)
}

// purchaseOrderGroupEmailDescription is the description an approver's email
// shows for a group.
func purchaseOrderGroupEmailDescription(description string) string {
	if description == "" {
		return "Varies by staff member"
	}
	return description
}

// createPurchaseOrderGroupNotification creates the one approval email for a
// group. The approver gets it instead of a po_approval_required email for each
// purchase order, which would be one email per staff member.
func createPurchaseOrderGroupNotification(app core.App, req *purchaseOrderGroupRequest, group *core.Record, preview purchaseOrderGroupPreview) (string, error) {
	vendorName := ""
	if vendor, err := app.FindRecordById("vendors", req.Template.Vendor); err == nil {
		vendorName = vendor.GetString("name")
	}
	return notifications.DispatchNotification(app, notifications.DispatchArgs{
		TemplateCode: "po_group_approval_required",
		RecipientUID: req.Template.Approver,
		Data: map[string]any{
			"POGroupId":   group.Id,
			"Count":       preview.Count,
			"Description": purchaseOrderGroupEmailDescription(group.GetString("description")),
			"VendorName":  vendorName,
			"Total":       purchaseOrderGroupAmount(preview.Total, preview.CurrencyCode),
			"ActionURL":   notifications.BuildActionURL(app, "/pos/pending"),
		},
		ActorUID: group.GetString("creator"),
		Mode:     notifications.DeliveryDeferred,
	})
}

// notifyPurchaseOrderGroupSecondApprover emails the priority second approver
// of a group the first time one of its purchase orders is first-approved and
// needs second approval. They get this one email for the whole group instead
// of a po_priority_second_approval_required email for each purchase order.
func notifyPurchaseOrderGroupSecondApprover(app core.App, po *core.Record, actorID string) error {
	groupID := po.GetString("po_group")
	recipientID := po.GetString("priority_second_approver")
	sent := struct {
		Count int `db:"count"`
	}{}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM notifications n
		JOIN notification_templates t ON t.id = n.template
		WHERE t.code = 'po_group_second_approval_required'
		  AND n.recipient = {:recipient}
		  AND json_extract(n.data, '$.POGroupId') = {:group}
	`).Bind(dbx.Params{"recipient": recipientID, "group": groupID}).One(&sent); err != nil {
		return err
	}
	if sent.Count > 0 {
		return nil
	}

	summary := purchaseOrderGroupSummary{}
	if err := app.DB().NewQuery(purchaseOrderGroupSummaryQuery).Bind(dbx.Params{"id": groupID}).One(&summary); err != nil {
		return err
	}
	_, err := notifications.DispatchNotification(app, notifications.DispatchArgs{
		TemplateCode: "po_group_second_approval_required",
		RecipientUID: recipientID,
		Data: map[string]any{
			"POGroupId":   groupID,
			"Count":       summary.Count,
			"CreatorName": summary.CreatorName,
			"Description": purchaseOrderGroupEmailDescription(summary.Description),
			"VendorName":  summary.VendorName,
			"ActionURL":   notifications.BuildActionURL(app, "/pos/pending"),
		},
		ActorUID: actorID,
		Mode:     notifications.DeliveryImmediate,
	})
	return err
}

// createCreatePurchaseOrderGroupHandler creates every purchase order in one
// transaction, so a request either creates the whole group or nothing. The
// client generates group_id before submitting; a repeat of the same request
// finds the group already exists and creates nothing.
//
// The rows are checked before the transaction starts so that checking a large
// upload doesn't hold the database's write lock.
func createCreatePurchaseOrderGroupHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requirePurchaseOrderGroupAccess(app, e.Auth); err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		req, err := bindPurchaseOrderGroupRequest(e, true)
		if err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}

		groupExists := func(app core.App) error {
			if existing, _ := app.FindRecordById("purchase_order_groups", req.GroupID); existing != nil {
				return &CodeError{Code: "group_exists", Message: "this purchase order group has already been created"}
			}
			return nil
		}
		if err := groupExists(app); err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}

		preview, records, err := evaluatePurchaseOrderGroup(e.Request.Context(), app, req, e.Auth.Id)
		if err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		if preview.ErrorCount > 0 {
			return respondPurchaseOrderGroupError(e, &purchaseOrderGroupRowsError{preview: preview})
		}

		notificationID := ""
		err = app.RunInTransaction(func(txApp core.App) error {
			// A repeat of this request may have finished while the rows
			// were being checked.
			if err := groupExists(txApp); err != nil {
				return err
			}

			groups, err := txApp.FindCollectionByNameOrId("purchase_order_groups")
			if err != nil {
				return err
			}
			group := core.NewRecord(groups)
			group.Id = req.GroupID
			group.Set("creator", e.Auth.Id)
			if !slices.Contains(req.UploadFields, "description") {
				group.Set("description", req.Template.Description)
			}
			if err := txApp.Save(group); err != nil {
				return err
			}

			for _, record := range records {
				record.Set("po_group", group.Id)
				if err := txApp.Save(record); err != nil {
					return err
				}
			}

			notificationID, err = createPurchaseOrderGroupNotification(txApp, req, group, preview)
			return err
		})
		if err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}

		if notificationID != "" {
			if err := notifications.SendNotificationByID(app, notificationID); err != nil {
				app.Logger().Error("failed to send purchase order group notification", "notification_id", notificationID, "error", err)
			}
		}

		return e.JSON(http.StatusOK, map[string]any{
			"group_id":      req.GroupID,
			"count":         preview.Count,
			"total":         preview.Total,
			"currency_code": preview.CurrencyCode,
		})
	}
}

type purchaseOrderGroupSummary struct {
	ID              string  `db:"id" json:"id"`
	Created         string  `db:"created" json:"created"`
	Creator         string  `db:"creator" json:"creator"`
	CreatorName     string  `db:"creator_name" json:"creator_name"`
	Description     string  `db:"description" json:"description"`
	VendorName      string  `db:"vendor_name" json:"vendor_name"`
	Count           int     `db:"count" json:"count"`
	UnapprovedCount int     `db:"unapproved_count" json:"unapproved_count"`
	Total           float64 `db:"total" json:"total"`
	CurrencyCode    string  `db:"currency_code" json:"currency_code"`
}

// purchaseOrderGroupSummaryQuery summarizes groups for super_po and admin
// holders, who can see every purchase order in a group. The description is
// the template's, blank when it came from the upload. Every purchase order in
// a group has the template's vendor and currency.
const purchaseOrderGroupSummaryQuery = `
	SELECT
		g.id,
		g.created,
		g.creator,
		TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')) AS creator_name,
		COALESCE(g.description, '') AS description,
		COALESCE(MIN(v.name), '') AS vendor_name,
		COUNT(po.id) AS count,
		COALESCE(SUM(CASE WHEN po.status = 'Unapproved' THEN 1 ELSE 0 END), 0) AS unapproved_count,
		COALESCE(SUM(po.total), 0) AS total,
		COALESCE(MIN(cur.code), 'CAD') AS currency_code
	FROM purchase_order_groups g
	LEFT JOIN profiles p ON p.uid = g.creator
	LEFT JOIN purchase_orders po ON po.po_group = g.id
	LEFT JOIN vendors v ON v.id = po.vendor
	LEFT JOIN currencies cur ON cur.id = po.currency
	WHERE {:id} = '' OR g.id = {:id}
	GROUP BY g.id
	ORDER BY g.created DESC
`

func createGetPurchaseOrderGroupsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requirePurchaseOrderGroupViewAccess(app, e.Auth); err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		groups := []purchaseOrderGroupSummary{}
		if err := app.DB().NewQuery(purchaseOrderGroupSummaryQuery).Bind(dbx.Params{"id": ""}).All(&groups); err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		return e.JSON(http.StatusOK, groups)
	}
}

func createGetPurchaseOrderGroupHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requirePurchaseOrderGroupViewAccess(app, e.Auth); err != nil {
			return respondPurchaseOrderGroupError(e, err)
		}
		id := strings.TrimSpace(e.Request.PathValue("id"))
		summary := purchaseOrderGroupSummary{}
		if id == "" {
			return e.JSON(http.StatusNotFound, map[string]string{
				"code":    "purchase_order_group_not_found",
				"message": "purchase order group not found",
			})
		}
		if err := app.DB().NewQuery(purchaseOrderGroupSummaryQuery).Bind(dbx.Params{"id": id}).One(&summary); err != nil {
			return e.JSON(http.StatusNotFound, map[string]string{
				"code":    "purchase_order_group_not_found",
				"message": "purchase order group not found",
			})
		}
		return e.JSON(http.StatusOK, summary)
	}
}
