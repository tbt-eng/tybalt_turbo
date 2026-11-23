package hooks

import (
	"testing"
	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
)

// Two integration-style tests exercising ProcessClient with the real test DB.

func TestProcessClient_BusinessDevelopmentLeadActiveSuccess(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	// active user from test DB
	leadID := "4r70mfovf22m9uh"

	clientsColl, _ := app.FindCollectionByNameOrId("clients")
	rec := core.NewRecord(clientsColl)
	rec.Set("business_development_lead", leadID)

	evt := &core.RecordRequestEvent{Record: rec}
	if err := ProcessClient(app, evt); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestProcessClient_BusinessDevelopmentLeadInactiveFailure(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	// inactive user from test DB
	leadID := "u_inactive"

	clientsColl, _ := app.FindCollectionByNameOrId("clients")
	rec := core.NewRecord(clientsColl)
	rec.Set("business_development_lead", leadID)

	evt := &core.RecordRequestEvent{Record: rec}
	if err := ProcessClient(app, evt); err == nil {
		t.Fatalf("expected failure, got success")
	}
}
