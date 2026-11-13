package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"

	"tybalt/hooks"
)

// createGetRateSheetDuplicatesHandler returns every group of rate sheets with
// identical rates, active and inactive, with the target of each group and
// the reason when a sheet can't be absorbed into it. It is read-only, so any
// signed-in user can read it; the absorb route checks the claims.
func createGetRateSheetDuplicatesHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		groups, err := hooks.FindRateSheetDuplicates(app)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to find duplicate rate sheets", err)
		}
		return e.JSON(http.StatusOK, map[string]any{"groups": groups})
	}
}

// createGetRateSheetAbsorbCandidatesHandler returns a target sheet and the
// sheets with identical rates, each with the reason when the target can't
// absorb it. The absorb page offers only the sheets without a reason.
func createGetRateSheetAbsorbCandidatesHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		target, candidates, ok, err := hooks.RateSheetAbsorbCandidates(app, e.Request.PathValue("id"))
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load rate sheets", err)
		}
		if !ok {
			return e.Error(http.StatusNotFound, "rate sheet not found", nil)
		}
		return e.JSON(http.StatusOK, map[string]any{
			"target":     target,
			"candidates": candidates,
		})
	}
}
