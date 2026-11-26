package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func recordPublic(r *core.Record) map[string]any {
	if r == nil {
		return nil
	}
	return r.PublicExport()
}
func requestRelated(app core.App, r *core.Record) (map[string]any, error) {
	result := map[string]any{}
	relations := map[string]string{"client": "clients", "job_owner": "clients", "contact": "client_contacts", "manager": "users", "alternate_manager": "users", "creator": "users", "approver": "users", "branch": "branches", "rate_sheet": "rate_sheets", "parent": "jobs", "proposal": "jobs", "invoicing_information": "client_invoicing_information"}
	for key, collection := range relations {
		if id := r.GetString(key); id != "" {
			row, err := app.FindRecordById(collection, id)
			if err != nil {
				return nil, err
			}
			data := recordPublic(row)
			if collection == "users" {
				data = workflowPerson(app, row)
			}
			result[key] = data
		}
	}
	if profile, ok := result["invoicing_information"].(map[string]any); ok {
		if id, ok := profile["contact"].(string); ok && id != "" {
			row, err := app.FindRecordById("client_contacts", id)
			if err != nil {
				return nil, err
			}
			result["billing_contact"] = recordPublic(row)
		}
	}
	// Bill the client name that the invoicing profile chooses: official or alias.
	if id := r.GetString("client"); id != "" {
		client, err := app.FindRecordById("clients", id)
		if err != nil {
			return nil, err
		}
		var profile *core.Record
		if id := r.GetString("invoicing_information"); id != "" {
			if profile, err = app.FindRecordById("client_invoicing_information", id); err != nil {
				return nil, err
			}
		}
		result["bill_to_name"] = utilities.BillingName(client, profile)
	}
	rates := []map[string]any{}
	if id := r.GetString("rate_sheet"); id != "" {
		// Highest rate first, the order the rate sheet itself prints in.
		rows, err := app.FindRecordsByFilter("rate_sheet_entries", "rate_sheet={:id}", "-rate", 0, 0, dbx.Params{"id": id})
		if err != nil {
			return nil, err
		}
		for _, entry := range rows {
			role, err := app.FindRecordById("rate_roles", entry.GetString("role"))
			if err != nil {
				return nil, err
			}
			rates = append(rates, map[string]any{"id": entry.Id, "role": role.GetString("name"), "role_id": role.Id, "role_name": role.GetString("name"), "rate": entry.GetFloat("rate"), "overtime_rate": entry.GetFloat("overtime_rate"), "updated": entry.Get("updated")})
		}
	}
	if err := requestProposedRates(app, r, result, rates); err != nil {
		return nil, err
	}
	result["rate_sheet_entries"] = rates
	return result, nil
}
func workflowPerson(app core.App, user *core.Record) map[string]any {
	data := map[string]any{"id": user.Id, "name": user.GetString("name"), "email": user.GetString("email"), "phone": ""}
	if p, err := app.FindFirstRecordByData("profiles", "uid", user.Id); err == nil {
		data["given_name"] = p.GetString("given_name")
		data["surname"] = p.GetString("surname")
		data["name"] = strings.TrimSpace(p.GetString("given_name") + " " + p.GetString("surname"))
		data["phone"] = p.GetString("phone")
	}
	return data
}

// requestReviewToken binds a review page to both saved data and one submission.
// submission_token is a random identity issued only by submit, not a password
// or a sequence number. Without it, recall/resubmit of unchanged data and PDF
// would recreate the same digest and accept the old page's checked checklist.
// Approval stores this digest on the attachment; creation checks that stored
// digest and the caller's digest against the current submission. Permission and
// workflow-state checks remain separate. A draft can have an empty submission
// token for viewing/printing, but approve/create must reject an empty token.
func requestReviewToken(app core.App, r *core.Record) (string, error) {
	return requestReviewDigest(app, r, false)
}

// requestLegacyReviewToken is the digest as it was before the proposal's
// status and updated time were left out: it holds the whole proposal record.
// Job creation accepts it as the saved approval of a request approved before
// that change. Remove it once no such request is awaiting creation.
func requestLegacyReviewToken(app core.App, r *core.Record) (string, error) {
	return requestReviewDigest(app, r, true)
}

// RequestLegacyReviewTokenForTest lets a test save the legacy digest that an
// approval made before the change would hold. Remove it with the legacy digest.
func RequestLegacyReviewTokenForTest(app core.App, id string) (string, error) {
	r, err := app.FindRecordById("job_requests", id)
	if err != nil {
		return "", err
	}
	return requestLegacyReviewToken(app, r)
}

func requestReviewDigest(app core.App, r *core.Record, legacy bool) (string, error) {
	related, err := requestRelated(app, r)
	if err != nil {
		return "", err
	}
	delete(related, "approver")
	delete(related, "creator")
	// For proposed rates the digest already holds what the approval is bound
	// to: the draft's record, with rates_approver and based_on as ids, and each
	// rate beside the source's. The approver's contact details and the source
	// sheet's own record are for display, and a change to them (a profile edit,
	// a rename, a newer revision) does not change what was approved.
	delete(related, "rate_sheet_rates_approver")
	delete(related, "rate_sheet_source")
	// The missing roles follow from the rates, which the digest holds, and
	// from the list of roles, which is not part of the request. Their rows
	// add the source sheet's rates for display. A request with a missing
	// role can't be approved or turned into a job.
	delete(related, "rate_sheet_missing_roles")
	delete(related, "rate_sheet_missing_rates")
	// The request binds the sheet by its id, so the sheet's code, name and
	// updated time are left out: a rename or a code change does not change
	// what was approved. Every other field stays, so a change to the revision,
	// active, the effective date or, for proposed rates, the reason, based_on,
	// rates_approver or rates_approved makes the approval stale. The rates are
	// in related["rate_sheet_entries"], apart from the record, so a changed
	// rate does too.
	if sheet, ok := related["rate_sheet"].(map[string]any); ok {
		delete(sheet, "code")
		delete(sheet, "name")
		delete(sheet, "updated")
	}
	// Procurement may award the linked proposal after approval and before
	// creation, by design: creation needs an Awarded proposal, and approval
	// does not. So the award must not make the approval stale. The status and
	// updated time are left out; every other field of the proposal stays, so
	// another edit to the proposal still makes the approval stale. The first
	// save of a proposal that came from the import also clears its _imported
	// mark, so that mark is left out too.
	if proposal, ok := related["proposal"].(map[string]any); ok && !legacy {
		delete(proposal, "status")
		delete(proposal, "updated")
		delete(proposal, "_imported")
	}
	data := map[string]any{"version": r.GetInt("version"), "submission_token": r.GetString("submission_token"), "related": related}
	for _, key := range requestBusinessFields {
		data[key] = r.Get(key)
	}
	allocs, err := requestAllocations(app, r.Id)
	if err != nil {
		return "", err
	}
	values := []map[string]any{}
	for _, a := range allocs {
		d, err := app.FindRecordById("divisions", a.GetString("division"))
		if err != nil {
			return "", err
		}
		values = append(values, map[string]any{"division": d.Id, "hours": a.GetFloat("hours"), "active": d.GetBool("active")})
	}
	data["allocations"] = values
	a, err := requestAttachment(app, r.Id)
	if err != nil {
		return "", err
	}
	if a != nil {
		data["document_hash"] = a.GetString("document_hash")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
func requestDetailHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error { return respondRequestDetail(app, e, e.Request.PathValue("id")) }
}
func respondRequestDetail(app core.App, e *core.RequestEvent, id string) error {
	var result map[string]any
	err := app.RunInTransaction(func(tx core.App) error {
		r, err := loadVisibleRequest(tx, id, e.Auth)
		if err != nil {
			return err
		}
		a, err := requestAttachment(tx, id)
		if err != nil {
			return err
		}
		allocs, err := requestAllocations(tx, id)
		if err != nil {
			return err
		}
		// Expand each division so the page shows its code and name, not its id.
		if errs := tx.ExpandRecords(allocs, []string{"division"}, nil); len(errs) > 0 {
			return fmt.Errorf("expand allocation divisions: %v", errs)
		}
		related, err := requestRelated(tx, r)
		if err != nil {
			return err
		}
		token, err := requestReviewToken(tx, r)
		if err != nil {
			return err
		}
		enabled, err := utilities.IsJobRequestApprovalRequired(tx)
		if err != nil {
			return err
		}
		hasJob, err := utilities.HasClaim(tx, e.Auth, "job")
		if err != nil {
			return err
		}
		hasAdmin, err := utilities.HasClaim(tx, e.Auth, "admin")
		if err != nil {
			return err
		}
		canReject, err := canRejectRequest(tx, r, e.Auth)
		if err != nil {
			return err
		}
		editable := r.GetString("creator") == e.Auth.Id && (r.GetString("status") == "Draft" || r.GetString("status") == "Recalled" || r.GetString("status") == "Rejected")
		submitted := r.GetString("status") == "Awaiting approval" || r.GetString("status") == "Awaiting creation"
		_, authorityErr := utilities.ProjectAuthorizationLimit(tx, e.Auth.Id, r.GetFloat("project_value"))
		checklistVersion, checklistQuestions, err := utilities.CurrentRequestChecklist(tx, requestChecklistScope(r))
		if err != nil {
			return err
		}
		attachment := recordPublic(a)
		if a != nil {
			attachment["document_url"] = "/api/job_attachments/" + a.Id + "/document"
		}
		result = map[string]any{"request": r, "allocations": allocs, "attachment": attachment, "related": related, "review_token": token, "checklist": map[string]any{"version": checklistVersion, "questions": checklistQuestions}, "enabled": enabled, "permissions": map[string]bool{"can_edit": enabled && editable, "can_submit": enabled && editable && r.GetString("status") != "Rejected", "can_recall": enabled && submitted && r.GetString("creator") == e.Auth.Id, "can_reject": enabled && canReject, "can_approve": enabled && r.GetString("status") == "Awaiting approval" && r.GetString("approver") == e.Auth.Id && authorityErr == nil, "can_create": enabled && hasJob && r.GetString("status") == "Awaiting creation", "can_delete": (enabled && editable) || (hasAdmin && abandonedRequest(r))}}
		return nil
	})
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, result)
}
func requestQueueSelection(app core.App, auth *core.Record, queue string) (string, dbx.Params, error) {
	if err := workflowActive(app, auth); err != nil {
		return "", nil, err
	}
	params := dbx.Params{"uid": auth.Id}
	switch queue {
	case "creation":
		if err := requestEnabled(app); err != nil {
			return "", nil, err
		}
		if err := workflowClaim(app, auth, "job"); err != nil {
			return "", nil, err
		}
		return creationQueueFilter, params, nil
	case "approval":
		if err := workflowClaim(app, auth, "project_authorization"); err != nil {
			return "", nil, err
		}
		return "approver={:uid} && status='Awaiting approval'", params, nil
	case "approved":
		// Requests this approver approved that still wait for a job, like the
		// PO list of approved POs that wait for a second approval.
		if err := workflowClaim(app, auth, "project_authorization"); err != nil {
			return "", nil, err
		}
		return "approver={:uid} && status='Awaiting creation'", params, nil
	case "abandoned":
		if err := workflowClaim(app, auth, "admin"); err != nil {
			return "", nil, err
		}
		params["cutoff"] = time.Now().UTC().AddDate(0, 0, -JOB_REQUEST_ABANDONMENT_DAYS).Format("2006-01-02 15:04:05.000Z")
		return "last_edited <= {:cutoff}", params, nil
	case "all":
		if err := workflowClaim(app, auth, "job"); err != nil {
			return "", nil, err
		}
		return "id != ''", params, nil
	case "", "mine":
		return "creator={:uid}", params, nil
	default:
		return "", nil, workflowError("Unknown request queue.")
	}
}
func requestListHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		filter, params, err := requestQueueSelection(app, e.Auth, e.Request.URL.Query().Get("queue"))
		if err != nil {
			return err
		}
		page, _ := strconv.Atoi(e.Request.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(e.Request.URL.Query().Get("perPage"))
		if limit < 1 || limit > 100 {
			limit = 50
		}
		var response map[string]any
		err = app.RunInTransaction(func(tx core.App) error {
			all, err := tx.FindRecordsByFilter("job_requests", filter, "-submitted,-created,id", 0, 0, params)
			if err != nil {
				return err
			}
			total := len(all)
			start := min((page-1)*limit, total)
			end := min(start+limit, total)
			items := []map[string]any{}
			for _, r := range all[start:end] {
				row := recordPublic(r)
				if c, err := tx.FindRecordById("clients", r.GetString("client")); err == nil {
					row["client_name"] = c.GetString("name")
				}
				items = append(items, row)
			}
			enabled, err := utilities.IsJobRequestApprovalRequired(tx)
			if err != nil {
				return err
			}
			response = map[string]any{"items": items, "totalItems": total, "page": page, "perPage": limit, "totalPages": (total + limit - 1) / limit, "enabled": enabled}
			return nil
		})
		if err != nil {
			return err
		}
		return e.JSON(http.StatusOK, response)
	}
}
func requestOptionsHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := workflowActive(app, e.Auth); err != nil {
			return err
		}
		allowed := false
		for _, claim := range []string{"time", "job", "project_authorization", "admin"} {
			ok, err := utilities.HasClaim(app, e.Auth, claim)
			if err != nil {
				return err
			}
			allowed = allowed || ok
		}
		if !allowed {
			return apis.NewForbiddenError("You cannot use project requests.", nil)
		}
		// The request form loads this before it shows, so each list holds only
		// what the form reads. Clients and jobs come from the UI's own stores.
		result := map[string]any{}
		for key, col := range map[string]string{"invoicing_profiles": "client_invoicing_information", "branches": "branches", "rate_sheets": "rate_sheets", "divisions": "divisions"} {
			var rows []*core.Record
			var err error
			if col == "rate_sheets" {
				// A request's proposed rates are listed nowhere.
				rows, err = app.FindAllRecords(col, dbx.HashExp{"request": ""})
			} else {
				rows, err = app.FindAllRecords(col)
			}
			if err != nil {
				return err
			}
			result[key] = rows
		}
		contacts := []struct {
			Id        string `db:"id" json:"id"`
			Client    string `db:"client" json:"client"`
			GivenName string `db:"given_name" json:"given_name"`
			Surname   string `db:"surname" json:"surname"`
			Email     string `db:"email" json:"email"`
		}{}
		if err := app.DB().NewQuery("SELECT id, client, given_name, surname, email FROM client_contacts").All(&contacts); err != nil {
			return err
		}
		result["contacts"] = contacts
		// Active people, named from the profile when there is one, as in
		// workflowPerson.
		users := []struct {
			Id         string `db:"id"`
			Name       string `db:"name"`
			Email      string `db:"email"`
			HasProfile bool   `db:"has_profile"`
			GivenName  string `db:"given_name"`
			Surname    string `db:"surname"`
		}{}
		if err := app.DB().NewQuery(`SELECT u.id, u.name, u.email, p.id IS NOT NULL AS has_profile,
				COALESCE(p.given_name, '') AS given_name, COALESCE(p.surname, '') AS surname
			FROM users u
			JOIN admin_profiles a ON a.id = (SELECT id FROM admin_profiles WHERE uid = u.id LIMIT 1) AND a.active = TRUE
			LEFT JOIN profiles p ON p.id = (SELECT id FROM profiles WHERE uid = u.id LIMIT 1)`).All(&users); err != nil {
			return err
		}
		people := []map[string]any{}
		for _, u := range users {
			person := map[string]any{"id": u.Id, "name": u.Name, "email": u.Email, "phone": ""}
			if u.HasProfile {
				person["given_name"] = u.GivenName
				person["surname"] = u.Surname
				person["name"] = strings.TrimSpace(u.GivenName + " " + u.Surname)
			}
			people = append(people, person)
		}
		result["users"] = people
		// The proposals a request can name. The form shows each from the jobs
		// store, so the id, client and status are enough.
		proposals := []struct {
			Id     string `db:"id" json:"id"`
			Client string `db:"client" json:"client"`
			Status string `db:"status" json:"status"`
		}{}
		if err := app.DB().NewQuery("SELECT id, client, status FROM jobs WHERE number LIKE 'P%' AND status IN ('In Progress', 'Submitted', 'Awarded')").All(&proposals); err != nil {
			return err
		}
		result["proposals"] = proposals
		enabled, err := utilities.IsJobRequestApprovalRequired(app)
		if err != nil {
			return err
		}
		result["enabled"] = enabled
		return e.JSON(http.StatusOK, result)
	}
}
func authorizationApproversHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := workflowActive(app, e.Auth); err != nil {
			return err
		}
		value, err := strconv.ParseFloat(e.Request.URL.Query().Get("value"), 64)
		if err != nil || !validWholeAmount(value, true) {
			return workflowError("Enter a positive whole-dollar project value.")
		}
		rows, err := app.FindRecordsByFilter("user_claims", "cid.name='project_authorization'", "uid", 0, 0)
		if err != nil {
			return err
		}
		items := []map[string]any{}
		for _, row := range rows {
			limit, err := utilities.ProjectAuthorizationLimit(app, row.GetString("uid"), value)
			if err != nil {
				continue
			}
			u, err := app.FindRecordById("users", row.GetString("uid"))
			if err != nil {
				return err
			}
			person := workflowPerson(app, u)
			person["max_value"] = limit
			items = append(items, person)
		}
		return e.JSON(http.StatusOK, map[string]any{"items": items})
	}
}
func addProjectWorkflowBadges(app core.App, auth *core.Record, counts map[string]int) error {
	if err := workflowActive(app, auth); err != nil {
		return nil
	}
	enabled, err := utilities.IsJobRequestApprovalRequired(app)
	if err != nil {
		return err
	}
	if enabled {
		hasJob, err := utilities.HasClaim(app, auth, "job")
		if err != nil {
			return err
		}
		if hasJob {
			filter, params, err := requestQueueSelection(app, auth, "creation")
			if err != nil {
				return err
			}
			rows, err := app.FindRecordsByFilter("job_requests", filter, "", 0, 0, params)
			if err != nil {
				return err
			}
			counts["/job-requests?queue=creation"] = len(rows)
		}
	}
	hasPA, err := utilities.HasClaim(app, auth, "project_authorization")
	if err != nil {
		return err
	}
	if hasPA {
		if enabled {
			rows, err := app.FindRecordsByFilter("job_requests", "approver={:uid} && status='Awaiting approval'", "", 0, 0, dbx.Params{"uid": auth.Id})
			if err != nil {
				return err
			}
			counts["/job-requests?queue=approval"] = len(rows)
		}
		rows, err := app.FindRecordsByFilter("job_attachments", "assigned_approver={:uid} && job!='' && status='Pending'", "", 0, 0, dbx.Params{"uid": auth.Id})
		if err != nil {
			return err
		}
		counts["/jobs/value-approvals"] = len(rows)
	}
	return nil
}
