<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkflowNavigation.ts","sha256":"1256d6d38af1111eb9ac5eb8619150a2edd7868efc96ad30d78ae6447f849373","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkflowNavigation.ts

## Client workspace draft lasts four hours

Lines 12-29

Asked as: draft too old; expired draft; draft expired after hours; saved form gone; draft not restored

Turbo keeps a project form draft for 4 hours (lifetime) after it is stored. An older draft cannot be restored. A draft with a creation time in the future is also refused.

## Selected PDF is lost on browser reload

Lines 13-15

Asked as: PDF gone after reload; lost attached PDF; selected file cleared; reattach PDF after refresh; PDF missing after refresh

The in-memory copy of a draft keeps a selected PDF while the user moves around in the same tab. Session storage keeps only text. After a browser reload the user must select the PDF again.

## Draft restores only for same user and same form

Lines 25-29

Asked as: draft not showing; someone else's draft; form loads empty; draft missing on form; wrong user draft

Turbo restores a kept draft only for the user who stored it and only on the same form path. It also needs the matching draft kind. Otherwise the form loads without the draft.

## Only safe in-app form addresses count as workflow returns

Lines 31-46

Asked as: invalid return link; return address rejected; which pages can return; bad redirect back; return url not allowed

Turbo accepts a return address only if it starts with a single "/" and has no backslash or line break. It must stay on the same site. It must be a project request page (/job-requests/new or /job-requests/<id>) or a job page (/jobs/add, /jobs/add/<id>, /jobs/add/from/<id> or /jobs/<id>/edit). Any other address is rejected.

## Return form is named project request or job form

Lines 49-52

Asked as: project request or job form; what is a project request; return form type name; which form am I returning to; job form versus request

A return address starting with /job-requests/ is called a "project request". Any other return address is called a "job form".

## Form that is not a valid return cannot open client workspace

Lines 65-68

Asked as: cannot return from client workspace; client workspace button greyed out; can't open client workspace; client workspace won't open; workspace link refused

When the current form page is not an accepted return address, Turbo refuses to leave for the client workspace. It shows the error "This form cannot return from the client workspace."

## Leaving for client workspace saves the draft under a new token

Lines 69-103

Asked as: draft saved before workspace; client_draft in address; where does draft token come from; workflow_return link; keeps my form when leaving

Before the user goes to the client workspace, Turbo stores the form draft with a new random token. It adds the token as client_draft to the form address, and the return link carries that address (workflow_return).

## A new draft replaces earlier drafts of the same form

Lines 79-95

Asked as: old draft replaced; earlier draft deleted; previous draft overwritten; workspace history old version; superseded draft comes back

Saving a draft deletes earlier drafts of the same kind, user and form path, so old workspace history cannot bring back a superseded draft. Turbo also deletes unreadable drafts and drafts older than 4 hours.

## Client workspace link passes the project contact

Lines 99-104

Asked as: client link missing contact; project contact not passed; contact not carried to client page; client page link details; open client details from form

The link goes to the client's details page. It carries the project contact (project_contact) only when the form has one.

## Contact and invoicing targets open the add page with return

Lines 104-108

Asked as: add contact from client workspace; add invoicing from workspace; return to client after adding contact; new invoicing address page; contact page doesn't return

When the target is a contact or invoicing, Turbo opens the add page for contacts or invoicing under that client. That page returns to the client workspace link (return_to) and keeps the project contact if there is one.

## Draft restores only with the client_draft token in the URL

Lines 111-119

Asked as: draft not restored after return; no draft on coming back; missing client_draft in url; workspace return empty form; draft not appearing on return

Turbo restores a draft only when the page address has a client_draft token that matches a stored draft. Without the token, it restores nothing.

## Draft form is chosen from the page path

Lines 121-127

Asked as: which form does draft belong to; draft for new job; draft on request page; how draft is matched to job; path decides draft form

A draft belongs to a request when the path is /job-requests/<id>. For /jobs/add it is a new job. For any other job path it is that job.

## Profile selection leaves the rest of the draft and the saved record unchanged

Lines 129-165

Asked as: selecting profile saves job; profile choice changes record; choosing invoicing profile saves; profile selection side effects; profile pick and draft storage fails

Selecting a profile changes only the draft's invoicing information and its selected profile. It does not save the job, request or value change. If browser storage fails, the draft stays as it was.

## Profile selection applies only to the originating form's draft

Lines 131-148

Asked as: profile not selected after returning; invoicing profile not applied; new profile missing on form; wrong client profile selected; profile not set in draft

After the user creates, copies or reuses an invoicing profile, Turbo sets it in the kept draft only if the return address is valid and carries a usable client_draft token. The draft must belong to the same user and the same form. Its client must equal the client of the profile. Otherwise Turbo selects nothing.

## Returning merges the draft with the freshly loaded form

Lines 184-212

Asked as: my edits overwritten on return; changes lost after workspace; server overwrote my changes; form reload loses edits; merge draft after client page

When the form reloads after the client workspace, a value the user changed since the last save keeps the draft's value, even if the server changed it too. Every other value takes the freshly loaded one. If there is no baseline, every draft value counts as changed.

## Record parts merge field by field and keep the server version

Lines 214-222

Asked as: parts merge conflict; allocations overwritten; version conflict on save; field edits lost in parts; which field value wins

For parts merged field by field, each field the user changed keeps the draft's value and the others take the loaded value. The version always comes from the server. Other parts, such as allocations, are merged whole.

## Saving the form clears its draft

Lines 226-235

Asked as: draft still there after save; form draft not cleared; old draft after saving; remove draft on save; saved but draft remains

Clearing removes the draft named by the client_draft token from memory and session storage. A blocked browser store does not stop a successful save.
