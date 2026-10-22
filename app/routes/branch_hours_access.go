package routes

import (
	"slices"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// requireBranchHoursAccess allows kpi and admin users to view any branch's
// hours. A branch manager may view only the branches they manage, because the
// report shows each person's hours without a job, which the kpi claim
// otherwise gates.
func requireBranchHoursAccess(app core.App, auth *core.Record, branchID string) error {
	for _, claim := range []string{"kpi", "admin"} {
		hasClaim, err := utilities.HasClaim(app, auth, claim)
		if err != nil {
			return apis.NewInternalServerError("Failed to check branch hours access", err)
		}
		if hasClaim {
			return nil
		}
	}
	managed, err := app.FindRecordsByFilter("branches", "manager = {:uid}", "", 0, 0, dbx.Params{"uid": auth.Id})
	if err != nil {
		return apis.NewInternalServerError("Failed to check branch hours access", err)
	}
	if len(managed) == 0 {
		return apis.NewForbiddenError("You need to be a branch manager or hold the kpi or admin claim to view this report", nil)
	}
	// The handler reports a missing branch.
	if branchID != "" && !slices.ContainsFunc(managed, func(branch *core.Record) bool { return branch.Id == branchID }) {
		return apis.NewForbiddenError("You can only view the hours of branches you manage", nil)
	}
	return nil
}
