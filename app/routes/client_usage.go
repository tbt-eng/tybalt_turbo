package routes

import (
	_ "embed" // for go:embed
	"net/http"

	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed client_usage_jobs.sql
var clientUsageJobsQuery string

//go:embed client_usage_requests.sql
var clientUsageRequestsQuery string

// ClientUsageJob is a job that uses a contact or an invoicing profile.
// ProjectContact and InvoiceRecipient say how a contact is used; for an
// invoicing profile only InvoiceRecipient is true.
type ClientUsageJob struct {
	ID               string `db:"id" json:"id"`
	Number           string `db:"number" json:"number"`
	Description      string `db:"description" json:"description"`
	Status           string `db:"status" json:"status"`
	ProjectContact   bool   `db:"project_contact" json:"project_contact"`
	InvoiceRecipient bool   `db:"invoice_recipient" json:"invoice_recipient"`
}

// ClientUsageRequest is a job request that uses a contact or an invoicing profile.
type ClientUsageRequest struct {
	ID               string `db:"id" json:"id"`
	Description      string `db:"description" json:"description"`
	Status           string `db:"status" json:"status"`
	Creator          string `db:"creator" json:"-"`
	Approver         string `db:"approver" json:"-"`
	ProjectContact   bool   `db:"project_contact" json:"project_contact"`
	InvoiceRecipient bool   `db:"invoice_recipient" json:"invoice_recipient"`
}

// ClientUsage lists what uses a contact or an invoicing profile. HiddenRequests
// counts the requests left out because the caller may not read them, so the
// lists still add up to the counts in the client details.
type ClientUsage struct {
	Jobs           []ClientUsageJob     `json:"jobs"`
	Requests       []ClientUsageRequest `json:"requests"`
	HiddenRequests int                  `json:"hidden_requests"`
}

// createGetClientUsageHandler lists the jobs and job requests that use one of
// a client's contacts or invoicing profiles. collection is client_contacts or
// client_invoicing_information, and the record id comes from the path.
//
// Every user may read jobs. Job requests follow the job_requests read rule: the
// job claim, the creator, or the approver with the project_authorization claim.
func createGetClientUsageHandler(app core.App, collection string) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		clientID := e.Request.PathValue("id")
		record, err := app.FindRecordById(collection, e.Request.PathValue("record"))
		if err != nil || record.GetString("client") != clientID {
			return e.Error(http.StatusNotFound, "record not found for this client", nil)
		}

		params := dbx.Params{"client": clientID, "contact": "", "profile": ""}
		if collection == "client_contacts" {
			params["contact"] = record.Id
		} else {
			params["profile"] = record.Id
		}

		usage := ClientUsage{Jobs: []ClientUsageJob{}, Requests: []ClientUsageRequest{}}
		if err := app.DB().NewQuery(clientUsageJobsQuery).Bind(params).All(&usage.Jobs); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load jobs", err)
		}
		var requests []ClientUsageRequest
		if err := app.DB().NewQuery(clientUsageRequestsQuery).Bind(params).All(&requests); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load job requests", err)
		}

		readsAll, err := utilities.HasClaim(app, e.Auth, "job")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		approves := false
		if !readsAll {
			if approves, err = utilities.HasClaim(app, e.Auth, "project_authorization"); err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check claims", err)
			}
		}
		for _, request := range requests {
			if readsAll || request.Creator == e.Auth.Id || (approves && request.Approver == e.Auth.Id) {
				usage.Requests = append(usage.Requests, request)
			} else {
				usage.HiddenRequests++
			}
		}
		return e.JSON(http.StatusOK, usage)
	}
}
