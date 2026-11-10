<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/project_workflow_notices.go","sha256":"add6fe6bd5d37066f50f159f46c65cef1137cb0c28f6720ac792c0dc0d18a256","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/project_workflow_notices.go

## Workflow emails never go to the person who acted

Lines 22-41

Turbo queues each project request and value-change email for every distinct recipient except the person who performed the action. It skips blank recipients and sends a recipient only one copy of each notice.

## Workflow emails are sent after the action is committed

Lines 43-52

Turbo sends project workflow notices only after the workflow action has committed. A notice that fails to send is logged and does not undo the action.

## Email amounts show as whole dollars with commas

Lines 54-62

Turbo rounds amounts in project workflow emails to the nearest whole dollar. It shows them with a $ sign and comma separators, such as $12,345.

## Email client name is blank if the client is missing

Lines 64-70

When Turbo cannot find the client record, the client name in the workflow email is left empty.

## Request email title is the first non-blank line, max 60 characters

Lines 72-85

Turbo builds the email subject title from the first non-blank line of the request's scope of work (description). It collapses extra whitespace. A title longer than 60 characters is cut to 59 characters plus an ellipsis (…).

## Earlier request notices that failed do not count as sent

Lines 87-97

Turbo treats a request notice as already sent to a recipient only if an earlier notice of the same template about the same request exists and its status is not error.

## Value-change emails go to the approver or the uploader

Lines 99-123

A submitted value change emails the selected approver (assigned_approver). An approved or rejected value change emails the person who uploaded the change's PDF (uploader). The emails use the templates job_value_submitted, job_value_approved and job_value_rejected.

## Value-change email contents

Lines 103-118

The value-change email gives the job number, description and client. It shows the job's current value, the proposed value and the previous value, plus the rejection reason and a link to the job details. It carries a Revised flag when the submission replaces an earlier version with the same approver.

## Request emails go to the approver or the requester

Lines 126-144

A submitted project request emails the selected approver. An approved or rejected request emails the requester (creator). The emails use the templates job_request_submitted, job_request_approved and job_request_rejected.

## Request email contents and requester flag

Lines 130-161

The request email gives the client, title, description, project value, rejection reason and a link to the request. It marks a rejection of an approved request as Withdrawn. Each recipient's copy sets ForRequester to true only when the recipient is the request's creator.

## Resubmitted request is marked as revised

Lines 143-146

For a submitted request, Turbo marks the email Revised when the approver already received a job_request_submitted notice about this request that did not fail. Recall and rejection keep the request, so it can go to an approver again.

## Created-project email goes to requester and project manager

Lines 147-156

When a project is created from a request (job_request_created), Turbo emails the requester and the new project's manager. The email gives the job number, whether the project is Active, the requester's name, and a link to the job details. The link goes to the job instead of the request.
