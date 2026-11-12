package main

import (
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"
)

// The update rule checks a category against the job only when the request
// sends one. The hook also checks it when the job changes, so a purchase order
// can't keep the category of its old job.
func TestPurchaseOrderCategoryFollowsJob(t *testing.T) {
	// PO 01897j210v01f69 is an Unapproved draft by author@soup.com on job
	// P24-487 with that job's category "case". Job 24-291 has no categories and
	// has division MD allocated.
	const (
		poPath      = "/api/collections/purchase_orders/records/01897j210v01f69"
		oldCategory = "ssivbnz8b2mqb01"
		newJobBody  = `"job": "u09fwwcg07y03m7", "division": "fy4i9poneukvq9u"`
	)
	author, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("a new job without a category is refused while the old category remains", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)

		body := reviseTestRequest(t, app, http.MethodPatch, poPath, `{`+newJobBody+`}`, author, http.StatusBadRequest)
		if !strings.Contains(body, `"category":{"code":"must_match_job"`) {
			t.Fatalf("stale category was not refused: %s", body)
		}
		if got := reviseTestPOField(t, app, "01897j210v01f69", "category"); got != oldCategory {
			t.Fatalf("category = %q, want it unchanged", got)
		}
	})

	t.Run("a new job with the category cleared is saved", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)

		reviseTestRequest(t, app, http.MethodPatch, poPath, `{`+newJobBody+`, "category": ""}`, author, http.StatusOK)
		if got := reviseTestPOField(t, app, "01897j210v01f69", "category"); got != "" {
			t.Fatalf("category = %q, want it cleared", got)
		}
	})

	t.Run("an edit that leaves the job and category alone skips the check", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)

		// Give the draft a category of another job, as older data can have.
		if _, err := app.DB().NewQuery("UPDATE purchase_orders SET category = 'kyed8fha9uaha27' WHERE id = '01897j210v01f69'").Execute(); err != nil {
			t.Fatal(err)
		}
		reviseTestRequest(t, app, http.MethodPatch, poPath, `{"description": "A draft with an older category"}`, author, http.StatusOK)
	})
}
