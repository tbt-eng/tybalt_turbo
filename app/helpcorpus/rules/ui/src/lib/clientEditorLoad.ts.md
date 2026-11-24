<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditorLoad.ts","sha256":"9b96cdd6bf7062bd590ec5e47fee5a8bfcc1258fbc324f1429ddb0a4b6c7e924","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditorLoad.ts

## Loading a client fails with a message when the server errors

Lines 13-25

Asked as: client not found error; unable to load client; client info fails to load; server error loading client

When Turbo cannot load a client's information and the server answers with an error, it shows "Client not found." if the client does not exist. For any other error it shows "Unable to load client information." An error that has no status code is passed on as it is.

## Creating an invoicing profile needs the time or job permission

Lines 35-42

Asked as: cannot create invoicing profile; invoicing profile permission; new invoicing profile refused; no access to add invoicing profile

To open the editor that creates a new invoicing profile, a person needs the time or job permission. Without it, Turbo refuses with "The time or job permission is required to create invoicing profiles." Turbo decides this with canCreateInvoicingProfile, and this file does not show its exact rule.

## Editing contacts or existing invoicing profiles needs the job permission

Lines 35-42

Asked as: job permission required; cannot edit contact; contact editor locked; manage clients access denied

To open the editor for a contact, or to edit an existing invoicing profile, a person needs the job permission (the job claim). Creating a contact also needs it. Without it, Turbo refuses with "The job permission is required to manage clients."

## Editor refuses a contact or invoicing profile from another client

Lines 44-54

Asked as: contact not found for client; invoicing profile not found; record from other client; wrong client contact opens

The editor looks for the record among the contacts or invoicing profiles of this client. If the record is not there, Turbo refuses with "Contact not found for this client." or "Invoicing profile not found for this client." When no record ID is given, the editor opens with an empty record to create a new one.

## Project contact in the link is kept only if it belongs to the client

Lines 55-62

Asked as: project contact cleared; project contact disappeared; link contact reset empty; wrong project contact in link

The editor takes a project contact from the project_contact value in the link. It keeps the value only if it is one of this client's contacts. Otherwise Turbo sets it to empty.

## Return link of the editor is checked

Lines 60-60

Asked as: return link redirect; return_to link check; where editor sends me back; redirect after saving client

The editor passes the return_to value of the link through safeReturnTo with the client ID, so the return destination is not used unchecked. This file does not show what safeReturnTo accepts.

## Any signed-in person can open the contact absorb page

Lines 65-69

Asked as: contact absorb page open; absorb proposal instead; merge contacts without permission; can't absorb contact

The page that absorbs a contact does not check the job or absorb permission. Without the absorb permission, the editor proposes the absorb (an absorb proposal) instead.

## Contact absorb page refuses a contact from another client

Lines 70-74

Asked as: absorb contact other client; merge contact wrong client; absorb targets list empty; contact not in this client

If the contact is not one of this client's contacts, Turbo refuses with "Contact not found for this client." Otherwise it loads all contacts of the client as possible targets.

## Contact absorb page returns to a checked link only when one is given

Lines 78-80

Asked as: no return after absorb; absorb page back link; return destination not set; absorb redirect missing

The contact absorb page passes return_to through safeReturnTo only when the link has a return_to value. Without it, no return destination is set.
