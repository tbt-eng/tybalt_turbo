package routes

import (
	"net/http"
	"testing"
	"tybalt/hooks"
	"tybalt/internal/testseed"
	"tybalt/utilities"

	"github.com/pocketbase/pocketbase/tests"
)

// heldTokenProbe reports what a token gets from the users view rule, which lets
// a person view only their own record: 200 for a valid token, 404 otherwise.
func heldTokenProbe(t *testing.T, app *tests.TestApp, userID, token string) int {
	t.Helper()
	return performClaimsJSONRequest(t, app, http.MethodGet, "/api/collections/users/records/"+userID, token, nil).Code
}

func userTokenKey(t *testing.T, app *tests.TestApp, userID string) string {
	t.Helper()
	user, err := app.FindRecordById("users", userID)
	if err != nil {
		t.Fatalf("failed to load user %s: %v", userID, err)
	}
	return user.TokenKey()
}

func superuserToken(t *testing.T, app *tests.TestApp) string {
	t.Helper()
	superuser, err := app.FindAuthRecordByEmail("_superusers", "test@example.com")
	if err != nil {
		t.Fatalf("failed to load superuser: %v", err)
	}
	token, err := superuser.NewAuthToken()
	if err != nil {
		t.Fatalf("failed to mint superuser token: %v", err)
	}
	return token
}

// Every way to make a person inactive makes the tokens they hold invalid. Each
// scenario changes one fixture profile through the path under test (that is
// the operation being tested), and the change is discarded with the test app.
// u_admin_only holds admin, so the active route refuses it; an admin is made
// inactive through the admin profile editor, HR's limited editor or the
// collection API instead.
func TestEveryDeactivationPathRevokesTokens(t *testing.T) {
	scenarios := []struct {
		name       string
		targetID   string
		actorEmail string // empty means the superuser
		method     string
		path       string
		body       map[string]any
		wantStatus int
	}{
		{
			name:       "active route",
			targetID:   "u_corp_noclaim",
			actorEmail: "it.identity@example.com",
			method:     http.MethodPost,
			path:       "/api/admin_profiles/ap_corp_noclaim/active",
			body:       map[string]any{"active": false},
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin profile editor deactivates an admin",
			targetID:   "u_admin_only",
			actorEmail: "author@soup.com",
			method:     http.MethodPost,
			path:       "/api/admin_profiles/ap_admin_only/save_with_claims",
			body: map[string]any{
				"admin_profile": map[string]any{"uid": "u_admin_only", "active": false},
				"claim_ids":     []string{adminClaimID},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "hr limited editor deactivates an admin",
			targetID:   "u_admin_only",
			actorEmail: "hr@example.com",
			method:     http.MethodPost,
			path:       "/api/admin_profiles/ap_admin_only/save_limited",
			body:       map[string]any{"admin_profile": map[string]any{"active": false}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "collection api deactivates an admin",
			targetID:   "u_admin_only",
			actorEmail: "author@soup.com",
			method:     http.MethodPatch,
			path:       "/api/collections/admin_profiles/records/ap_admin_only",
			body:       map[string]any{"active": false},
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin profile editor creates an inactive profile",
			targetID:   "uvaluation00002",
			actorEmail: "author@soup.com",
			method:     http.MethodPost,
			path:       "/api/admin_profiles/save_with_claims",
			body: map[string]any{
				"admin_profile": map[string]any{
					"uid":                     "uvaluation00002",
					"active":                  false,
					"work_week_hours":         40,
					"default_charge_out_rate": 50,
					"skip_min_time_check":     "no",
					"payroll_id":              "990501",
					"default_branch":          defaultBranchID,
				},
				"claim_ids": []string{},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "superuser deletes the profile",
			targetID:   "u_corp_noclaim",
			method:     http.MethodDelete,
			path:       "/api/collections/admin_profiles/records/ap_corp_noclaim",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			hooks.AddHooks(app)
			AddRoutes(app)

			target, err := app.FindRecordById("users", scenario.targetID)
			if err != nil {
				t.Fatalf("failed to load target: %v", err)
			}
			heldToken, err := target.NewAuthToken()
			if err != nil {
				t.Fatalf("failed to mint target token: %v", err)
			}
			if code := heldTokenProbe(t, app, scenario.targetID, heldToken); code != http.StatusOK {
				t.Fatalf("probe before the change = %d, want %d", code, http.StatusOK)
			}
			keyBefore := userTokenKey(t, app, scenario.targetID)

			actor := ""
			if scenario.actorEmail == "" {
				actor = superuserToken(t, app)
			} else {
				actor = authTokenForEmail(t, app, scenario.actorEmail)
			}
			var body any
			if scenario.body != nil {
				body = scenario.body
			}
			rec := performClaimsJSONRequest(t, app, scenario.method, scenario.path, actor, body)
			if rec.Code != scenario.wantStatus {
				t.Fatalf("change status = %d, want %d; body=%s", rec.Code, scenario.wantStatus, rec.Body.String())
			}

			active, err := utilities.IsUserActive(app, scenario.targetID)
			if err != nil {
				t.Fatal(err)
			}
			if active {
				t.Fatal("the target is still active after the change")
			}
			if userTokenKey(t, app, scenario.targetID) == keyBefore {
				t.Fatal("the change did not give the target a new token key")
			}
			if code := heldTokenProbe(t, app, scenario.targetID, heldToken); code != http.StatusNotFound {
				t.Fatalf("probe with the held token after the change = %d, want %d", code, http.StatusNotFound)
			}
		})
	}
}

// After deactivation the token the person held fails everywhere: a custom
// route, a collection API request governed by a rule that does not look at
// account status, and token refresh. The test deactivates the fixture profile
// ap_corp_noclaim through the active route, which is the operation under test.
func TestDeactivatedTokenFailsOnEverySurface(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	hooks.AddHooks(app)
	AddRoutes(app)

	const targetID = "u_corp_noclaim"
	heldToken := authTokenForEmail(t, app, "corp.noclaim@example.com")

	// The person lacks the claims this route needs, so a valid token gets 403
	// and an invalid one 401.
	const customRoute = "/api/admin_profiles/active_toggle_targets"
	const refresh = "/api/collections/users/auth-refresh"

	if rec := performClaimsJSONRequest(t, app, http.MethodGet, customRoute, heldToken, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("custom route before = %d, want %d; body=%s", rec.Code, http.StatusForbidden, rec.Body.String())
	}
	if code := heldTokenProbe(t, app, targetID, heldToken); code != http.StatusOK {
		t.Fatalf("collection rule before = %d, want %d", code, http.StatusOK)
	}
	if rec := performClaimsJSONRequest(t, app, http.MethodPost, refresh, heldToken, nil); rec.Code != http.StatusOK {
		t.Fatalf("auth-refresh before = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	actor := authTokenForEmail(t, app, "it.identity@example.com")
	if rec := performClaimsJSONRequest(t, app, http.MethodPost, "/api/admin_profiles/ap_corp_noclaim/active", actor, map[string]any{"active": false}); rec.Code != http.StatusOK {
		t.Fatalf("deactivate status = %d; body=%s", rec.Code, rec.Body.String())
	}

	if rec := performClaimsJSONRequest(t, app, http.MethodGet, customRoute, heldToken, nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("custom route after = %d, want %d; body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
	if code := heldTokenProbe(t, app, targetID, heldToken); code != http.StatusNotFound {
		t.Fatalf("collection rule after = %d, want %d", code, http.StatusNotFound)
	}
	if rec := performClaimsJSONRequest(t, app, http.MethodPost, refresh, heldToken, nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("auth-refresh after = %d, want %d; body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}
