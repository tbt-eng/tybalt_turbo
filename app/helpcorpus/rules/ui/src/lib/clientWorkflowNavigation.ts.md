<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkflowNavigation.ts","sha256":"00ead9fffdce51cad1288515d6af12809fda4d2ec2d8884995a1529533bbcbf7","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkflowNavigation.ts

## Draft restores only with the return token in this tab

Lines 2-2

Turbo restores a draft only when the page address has the client_draft token and the draft kind matches the form. Without the token, or with a different kind, the form opens without a draft.

## Client workspace drafts expire after four hours

Lines 12-29

A saved project form draft can be restored only within 4 hours of when it was saved (lifetime). A draft exactly 4 hours old is still usable. A draft whose saved time is in the future is not usable.

## Selected PDF survives in-tab navigation but not reload

Lines 13-19

Turbo keeps an in-memory copy of the draft, which keeps a selected PDF while the user navigates in the same tab. Session storage keeps only text, so after a browser reload the user must select the PDF again.

## A draft restores only for the same user and form

Lines 25-29

Turbo restores a draft only for the user who saved it and on the same form path. A different user or a different form path gets no draft.

## Only job and project request pages are valid return targets

Lines 31-46

Turbo accepts a return address only if it is a path on this site that starts with a single slash. It must be /job-requests/new, a /job-requests/ page, /jobs/add (optionally with /from/ and an id) or /jobs/<id>/edit. Addresses with a backslash or line break, with a double slash, or on any other page are rejected. Query and hash parts are kept.

## Return target is named project request or job form

Lines 49-52

A return address that starts with /job-requests/ is named a project request. Any other return address is named a job form.

## Leaving a form for client workspace needs a valid return form

Lines 65-68

If the current form page is not a valid return address, Turbo refuses to open the client workspace. It shows the error "This form cannot return from the client workspace."

## Opening the client workspace saves the form as a draft

Lines 69-98

Before the user goes to the client workspace, Turbo saves the form contents under a new random token, with the user, the form path and the time. It adds the token (client_draft) to the form address, which the user returns to later.

## A new draft replaces earlier drafts of the same form

Lines 79-95

Saving a draft deletes earlier drafts of the same kind, user and form path, so old workspace history entries cannot bring back a superseded draft. Turbo also deletes unreadable drafts and drafts older than 4 hours.

## Client workspace link carries return address and project contact

Lines 99-104

The link to the client details page includes the form address to return to (workflow_return). It also includes the project contact (project_contact) when one is set.

## Contact and invoicing editors return to the workspace

Lines 105-108

When the target is a contact or invoicing profile, Turbo opens the add page for contacts or invoicing under the client. That page returns to the client workspace link (return_to) and passes along the project contact if one is set.

## Draft restore requires a matching token and kind

Lines 111-119

When the address has no client_draft token, no draft is read. The draft is also ignored if it is expired, or belongs to another user, path or kind.

## Draft form kind is derived from the return path

Lines 122-127

A project request form uses the draft kind request:<id>. The new job form uses job:new, and an existing job form uses job:<id>.

## New or reused invoicing profile is selected only in the originating form

Lines 129-165

After a profile is created, copied or reused, Turbo sets it as the invoicing information in the saved draft. This works only if the return address holds a usable draft token for the same user and kind. It also requires that the draft's client matches the profile's client. Otherwise nothing is selected. Saving the profile does not save the job or request.

## Profile selection needs all identifiers

Lines 137-139

Turbo selects no profile if the return address, user, client or profile is missing.

## Storage error leaves draft unchanged when selecting profile

Lines 157-162

If saving the profile selection to session storage fails, Turbo selects no profile and leaves the draft as it was.

## Successful save clears the draft

Lines 167-177

After a form is saved, Turbo deletes the draft named by the client_draft token in the address. A blocked browser store does not stop the save.
