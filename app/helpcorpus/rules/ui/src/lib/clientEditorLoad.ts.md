<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditorLoad.ts","sha256":"868ec5a95c0927d6f084569298a1cfd306e1f42bcd6d8bddccc30174a9753199","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditorLoad.ts

## Client page load errors: not found or unable to load

Lines 13-25

Asked as: client page error; client not found; unable to load client; client page won't open; 404 on client

When Turbo fails to load a client and the server answers with an error status (400 to 599), the page shows an error. A 404 shows "Client not found." Any other error status shows "Unable to load client information." An error without such a status is passed on unchanged.

## Creating an invoicing profile needs time or job permission

Lines 35-42

Asked as: can't create invoicing profile; new invoicing profile refused; time permission needed; add invoice profile denied; invoicing profile permission

To open the editor for a new invoicing profile (an invoicing editor with no record), a person must pass canCreateInvoicingProfile (the time or job permission). Otherwise Turbo refuses with "The time or job permission is required to create invoicing profiles."

## Managing client contacts or existing profiles needs job permission

Lines 36-42

Asked as: can't add client contact; job permission missing; manage client contacts denied; edit invoicing profile refused; no access to client contacts

To open the editor for a contact, or for an existing invoicing profile, a person needs the job permission (the job claim). Otherwise Turbo refuses with "The job permission is required to manage clients."

## Editing a missing contact or invoicing profile is refused

Lines 44-54

Asked as: contact missing for client; profile not found; contact deleted from client; invoicing profile gone; wrong contact id

When a record ID is given and the client has no contact or invoicing profile with that ID, Turbo refuses. For a contact it says "Contact not found for this client." For an invoicing profile it says "Invoicing profile not found for this client."

## New contact or profile editor starts with an empty record

Lines 45-47

Asked as: blank contact form; new contact empty form; empty invoicing profile; add new contact form; fresh profile form

When no record ID is given, the editor opens with an empty record, to create a new contact or invoicing profile.

## Project contact preselection must be a contact of this client

Lines 55-61

Asked as: project contact not preselected; wrong contact picked on project; contact dropdown empty; preselected contact ignored; project contact from link

The editor keeps the project_contact value from the URL only if it is the ID of one of the client's contacts. Otherwise it uses an empty value.

## Return link of the client editor is checked

Lines 60-60

Asked as: back link broken; return link client editor; go back after saving client; return button wrong page; where back button goes

The editor's return link (return_to) passes through safeReturnTo for this client before use. This file does not show what that function returns.

## Merging contacts needs absorb permission

Lines 67-68

Asked as: can't merge contacts; merge button missing; absorb permission needed; merge contacts denied; no merge access

To open the contact merge page, a person needs the absorb permission (the absorb claim). Otherwise Turbo refuses with "The absorb permission is required to merge contacts."

## Merge page refuses a contact from another client

Lines 70-71

Asked as: merge contact from other client; contact belongs to another client; merge page contact not found; cross client merge; wrong client contact merge

The merge page refuses when the contact is not one of this client's contacts. Turbo says "Contact not found for this client."

## Merge page lists all contacts of the client

Lines 72-74

Asked as: merge target list empty; which contact to keep; contacts to merge into; list of client contacts; choose contact to merge

The merge page loads every contact (client_contacts) of this client. Turbo uses the list to choose a merge target.

## Merge page return link is set only when given

Lines 78-80

Asked as: no back link merge page; merge page no return link; return link missing on merge; merge page back button absent; return after merging

The merge page passes the return_to value through safeReturnTo only when the URL has a return_to parameter. Without it, the page has no return link.
