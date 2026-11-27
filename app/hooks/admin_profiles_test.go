package hooks

import (
	"errors"
	"testing"
	"time"
	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// These tests deactivate the fixture profile ap_corp_noclaim (user
// u_corp_noclaim), which is the operation under test. Changes are discarded
// with the test app.

func revokeTestTokenKey(t *testing.T, app core.App, userID string) string {
	t.Helper()
	user, err := app.FindRecordById("users", userID)
	if err != nil {
		t.Fatalf("failed to load user %s: %v", userID, err)
	}
	return user.TokenKey()
}

// runWithTimeout fails the test if fn does not return, so a hook that writes
// through the outer app inside a transaction shows up as a failure instead of
// a hung test.
func runWithTimeout(t *testing.T, fn func() error) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- fn() }()
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
		t.Fatal("the save did not finish: the hook may be deadlocked")
		return nil
	}
}

func newRevokeTestApp(t *testing.T) *tests.TestApp {
	t.Helper()
	app := testseed.NewSeededTestApp(t)
	AddHooks(app)
	return app
}

func TestRevokeTokensOnDeactivation_InsideATransaction(t *testing.T) {
	sentinel := errors.New("roll back")
	for _, commit := range []bool{true, false} {
		name := "commit"
		if !commit {
			name = "rollback"
		}
		t.Run(name, func(t *testing.T) {
			app := newRevokeTestApp(t)
			keyBefore := revokeTestTokenKey(t, app, "u_corp_noclaim")

			err := runWithTimeout(t, func() error {
				return app.RunInTransaction(func(txApp core.App) error {
					profile, err := txApp.FindRecordById("admin_profiles", "ap_corp_noclaim")
					if err != nil {
						return err
					}
					profile.Set("active", false)
					if err := txApp.Save(profile); err != nil {
						return err
					}
					if revokeTestTokenKey(t, txApp, "u_corp_noclaim") == keyBefore {
						t.Error("the transaction does not see the new token key")
					}
					if !commit {
						return sentinel
					}
					return nil
				})
			})
			if commit && err != nil {
				t.Fatalf("transaction failed: %v", err)
			}
			if !commit && !errors.Is(err, sentinel) {
				t.Fatalf("transaction error = %v, want the rollback sentinel", err)
			}

			profile, err := app.FindRecordById("admin_profiles", "ap_corp_noclaim")
			if err != nil {
				t.Fatal(err)
			}
			keyAfter := revokeTestTokenKey(t, app, "u_corp_noclaim")
			if commit {
				if profile.GetBool("active") || keyAfter == keyBefore {
					t.Fatalf("after commit: active=%v, key changed=%v; want false, true", profile.GetBool("active"), keyAfter != keyBefore)
				}
				return
			}
			if !profile.GetBool("active") || keyAfter != keyBefore {
				t.Fatalf("after rollback: active=%v, key changed=%v; want true, false", profile.GetBool("active"), keyAfter != keyBefore)
			}
		})
	}
}

// A profile whose user row is gone can still be deactivated. The test deletes
// the users row of u_corp_noclaim with SQL to make that state, and saves the
// profile without validation, because its uid relation no longer resolves.
func TestRevokeTokensOnDeactivation_MissingUser(t *testing.T) {
	app := newRevokeTestApp(t)
	if _, err := app.DB().NewQuery("DELETE FROM users WHERE id = 'u_corp_noclaim'").Execute(); err != nil {
		t.Fatal(err)
	}
	profile, err := app.FindRecordById("admin_profiles", "ap_corp_noclaim")
	if err != nil {
		t.Fatal(err)
	}
	profile.Set("active", false)
	if err := app.SaveNoValidate(profile); err != nil {
		t.Fatalf("deactivating a profile with no user failed: %v", err)
	}
	if err := app.Delete(profile); err != nil {
		t.Fatalf("deleting a profile with no user failed: %v", err)
	}
}

// A superuser can move an active profile to another user. The old user is no
// longer active and is signed out; the new user keeps their tokens. The test
// moves ap_corp_noclaim to uvaluation00002, a fixture user with no profile.
func TestRevokeTokensOnDeactivation_UIDChange(t *testing.T) {
	app := newRevokeTestApp(t)
	oldKey := revokeTestTokenKey(t, app, "u_corp_noclaim")
	newKey := revokeTestTokenKey(t, app, "uvaluation00002")

	profile, err := app.FindRecordById("admin_profiles", "ap_corp_noclaim")
	if err != nil {
		t.Fatal(err)
	}
	profile.Set("uid", "uvaluation00002")
	if err := app.SaveNoValidate(profile); err != nil {
		t.Fatal(err)
	}
	if revokeTestTokenKey(t, app, "u_corp_noclaim") == oldKey {
		t.Error("the old user kept their token key")
	}
	if revokeTestTokenKey(t, app, "uvaluation00002") != newKey {
		t.Error("the new user's token key changed")
	}
}
