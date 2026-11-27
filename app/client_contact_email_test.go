package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"tybalt/errs"
	"tybalt/internal/testutils"
)

// The fixture contact 184qbpajciwrrp2 of Sepulchi & Co has ssecco@example.com.
const (
	takenContactEmail  = "ssecco@example.com"
	takenEmailMessage  = "This email is already a contact for Sepulchi & Co."
	contactEmailActor  = "author@soup.com"
	contactRecordsPath = "/api/collections/client_contacts/records"
)

// emailInUseError checks that a refused save names the client of the contact
// that has the email, at the top and on the email field.
func emailInUseError(t *testing.T, body []byte) {
	t.Helper()
	var result struct {
		Message string                    `json:"message"`
		Data    map[string]errs.CodeError `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("decode error body: %v; body=%s", err, body)
	}
	if result.Message != takenEmailMessage {
		t.Fatalf("message = %q, want %q", result.Message, takenEmailMessage)
	}
	field, ok := result.Data["email"]
	if !ok || field.Code != "email_in_use" || field.Message != takenEmailMessage {
		t.Fatalf("email field error = %+v, want code email_in_use and %q", result.Data["email"], takenEmailMessage)
	}
}

func contactEmailToken(t *testing.T) string {
	t.Helper()
	token, err := testutils.GenerateRecordToken("users", contactEmailActor)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

// A new contact with the email of another client's contact is refused, and the
// error names that client.
func TestClientContactEmailInUseOnCreate(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	response := performJSONRequest(t, app, http.MethodPost, contactRecordsPath, contactEmailToken(t), map[string]any{
		"client": "cwclient0000002", "given_name": "Moved", "surname": "Person", "email": takenContactEmail,
	})
	mustStatus(t, response, http.StatusBadRequest)
	emailInUseError(t, response.Body.Bytes())
}

// Changing a contact to an email that another contact has is refused, and the
// stored email does not change.
func TestClientContactEmailInUseOnUpdate(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	response := performJSONRequest(t, app, http.MethodPatch, contactRecordsPath+"/cwcontact000001", contactEmailToken(t), map[string]any{
		"email": takenContactEmail,
	})
	mustStatus(t, response, http.StatusBadRequest)
	emailInUseError(t, response.Body.Bytes())
	stored, err := app.FindRecordById("client_contacts", "cwcontact000001")
	if err != nil {
		t.Fatal(err)
	}
	if stored.GetString("email") != "alex.workspace@example.com" {
		t.Fatalf("email changed to %q", stored.GetString("email"))
	}
}

// Saves made outside requests get the same error from the model hook.
func TestClientContactEmailInUseOnModelSave(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	contact, err := app.FindRecordById("client_contacts", "cwcontact000002")
	if err != nil {
		t.Fatal(err)
	}
	contact.Set("email", takenContactEmail)
	err = app.Save(contact)
	var hookErr *errs.HookError
	if !errors.As(err, &hookErr) || hookErr.Message != takenEmailMessage {
		t.Fatalf("model save error = %v, want %q", err, takenEmailMessage)
	}
}

// A contact keeps its own email: an update that leaves the email unchanged
// still saves.
func TestClientContactUnchangedEmailSaves(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	updated := clientSetupJSON(t, app, http.MethodPatch, contactRecordsPath+"/cwcontact000001", contactEmailActor, map[string]any{
		"email": "alex.workspace@example.com", "surname": "Renamed",
	}, http.StatusOK)
	if updated["surname"] != "Renamed" || updated["email"] != "alex.workspace@example.com" {
		t.Fatalf("update did not save: %+v", updated)
	}

	// The model path also accepts the contact's own email.
	contact, err := app.FindRecordById("client_contacts", "cwcontact000001")
	if err != nil {
		t.Fatal(err)
	}
	contact.Set("phone", "807-555-0199")
	if err := app.Save(contact); err != nil {
		t.Fatalf("model save with unchanged email: %v", err)
	}
}

// Blank emails are exempt, as in the unique index. The fixtures have no contact
// without an email, so this test creates two through the API, the operation
// under test. They go to the empty fixture client and are discarded with the
// test app.
func TestClientContactBlankEmailSaves(t *testing.T) {
	app := newClientWorkspaceTestApp(t)
	for _, name := range []string{"First", "Second"} {
		created := clientSetupJSON(t, app, http.MethodPost, contactRecordsPath, contactEmailActor, map[string]any{
			"client": "cwclient0000002", "given_name": name, "surname": "No Email", "email": "", "no_email": true,
		}, http.StatusOK)
		if created["email"] != "" {
			t.Fatalf("%s contact email = %v, want blank", name, created["email"])
		}
	}

	// Clearing the email of a contact also saves while other contacts have
	// none.
	clientSetupJSON(t, app, http.MethodPatch, contactRecordsPath+"/cwcontact000003", contactEmailActor, map[string]any{
		"email": "", "no_email": true,
	}, http.StatusOK)
}
