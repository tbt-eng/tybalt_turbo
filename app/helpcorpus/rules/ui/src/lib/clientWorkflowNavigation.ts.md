<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkflowNavigation.ts","sha256":"bad05ae0c80670e9e15bb160b61b0038ae6ad59e4797e8ffa97db27d06ff8121","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkflowNavigation.ts

## Saved project form draft expires after four hours

Lines 12-29

Asked as: draft expired; old saved draft gone; form draft not restored; restore draft after hours

A kept draft of a project form or project request can be restored for 4 hours (lifetime). An older draft, or one with a creation time in the future, is not restored.

## Draft is stored in memory and in session storage

Lines 13-19

Asked as: attachment lost after reload; PDF gone after refresh; selected file missing; draft pdf not kept

Turbo keeps the draft both in memory and in the tab's session storage. Only the memory copy keeps a selected PDF. After a browser reload only the text comes back, so the user must select the PDF again.

## Only project and project request form addresses can be a return target

Lines 31-48

Asked as: back link not working; return to form missing; read only request page return; where does back go

Turbo accepts a workflow return address only if it is a local path (starts with one "/", not "//", with no backslash or line break). It must be a new or edit page of a project request (/job-requests/new or /job-requests/{id}/edit) or a job page (/jobs/add, /jobs/add/{id}, /jobs/add/from/{id} or /jobs/{id}/edit). Any other address, including a request's read-only page, gives no return (safeWorkflowReturn).

## Return form is named project request or job form

Lines 50-54

Asked as: what is the job form called; workspace return form name; which form am I returning to; request or job label

When Turbo names the form that a client workspace return leads back to, an address starting with /job-requests/ is called "project request". Every other address is called "job form" (workflowFormName).

## Leaving a form for the client workspace keeps a draft

Lines 58-105

Asked as: going to client page from form; open client workspace from form; form cannot return error; lost my form going to client

When a user goes from a project or request form to the client workspace, Turbo saves the form as a draft under a new random token. It adds client_draft={token} to the form address. The workspace link is /clients/{client}/details with workflow_return and, if there is one, project_contact. If the current address is not an allowed return form, Turbo throws the error "This form cannot return from the client workspace."

## A new draft replaces earlier drafts of the same form

Lines 79-97

Asked as: old draft came back; back button restores old form; previous draft reappeared; browser history old form data

When Turbo saves a draft, it deletes every earlier draft of the same kind, user and form path, in memory and in session storage. It also deletes unreadable drafts and drafts older than four hours. Old browser history entries therefore cannot bring back a superseded draft.

## Contact and invoicing targets open the add editor with a return link

Lines 106-110

Asked as: add contact opens wrong page; new contact link; invoicing add page redirect; workspace opens instead of add

For the target "contact" or "invoicing", Turbo opens /clients/{client}/contacts/add or /clients/{client}/invoicing/add instead of the workspace. The link has return_to set to the workspace address, and project_contact when there is one. The default target "workspace" opens the client details page.

## A form restores a draft only with its return token

Lines 113-121

Asked as: draft not loading on form; form came back empty; missing draft after returning; client_draft link not restoring

A form restores a draft only when its address has a client_draft token. The draft must be usable (same user, same form path, not expired) and of the form's kind. Otherwise Turbo restores nothing.

## Which draft kind a form uses

Lines 123-129

Asked as: where is my draft stored; draft for new job; which draft belongs to request; job draft key

Turbo finds the draft of a project request by the kind request:{id}. A new job (/jobs/add) uses job:new. Any other job form uses job:{id}.

## New profile is selected only in the originating form's draft

Lines 131-167

Asked as: new profile not selected; invoicing profile not added to form; profile created but not on job; reuse profile doesn't show

After a user creates, copies or reuses an invoicing profile, Turbo sets it as the invoicing information (invoicing_information) in the saved draft, and records selectedProfile. It does this only if the return address is valid and has a token, the draft is usable and of the right kind, and the draft's client equals the client of the profile. If any check fails, Turbo selects nothing. Saving the profile does not save the job or request.

## Profile selection leaves the draft intact if storage fails

Lines 159-164

Asked as: profile not chosen storage error; browser blocks saving selection; session storage full profile; selection failed draft unchanged

If session storage refuses the update, Turbo does not select the profile, and the draft stays as it was.

## Returning merges the draft with the reloaded form

Lines 186-214

Asked as: my edits lost after return; changes overwritten when back; merge changes on return; server changed my form

On return Turbo compares each part of the draft with the baseline (the form at its last save). A part that the user changed keeps the draft's value, even if the server changed it too. Any unchanged part takes the reloaded value. With no readable baseline, every draft value counts as changed.

## Fieldwise parts merge field by field and keep the server version

Lines 216-224

Asked as: version conflict after return; field changes merged; which values win on return; server version overwrite

For parts merged field by field, Turbo keeps the draft's value for each field the user changed, and takes the reloaded value for the others. The version field always comes from the server, never from the draft.

## Saving the form clears its draft

Lines 228-238

Asked as: draft still there after save; saved form keeps old draft; draft remains after submit; clear draft on save

After a form saves, Turbo deletes the draft named by the client_draft token in the address, from memory and from session storage. If the browser blocks storage, the save still succeeds.
