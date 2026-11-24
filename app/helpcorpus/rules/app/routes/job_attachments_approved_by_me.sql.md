<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_attachments_approved_by_me.sql","sha256":"d889c324c41a1220e0a3335d3ec18986f32efcbf4ddc191980dd60ffcc629e93","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_attachments_approved_by_me.sql

## Only recent approvals within the cutoff are listed

Lines 1-22

Asked as: approvals older than 30 days; old approvals missing; approved list date cutoff; why approval not showing

Turbo lists only entries whose date is on or after the cutoff, which is the last 30 days according to the file comment. An entry dated exactly at the cutoff is included. The file does not set the cutoff value, which is passed in as {:cutoff}.

## Date of a new job approval is the job's creation date

Lines 12-13

Asked as: new job approval date; job created date approval; approval date wrong new job; when was job approved

For an initial approval (initial_approval), a new job, Turbo uses the date the job was created (jobs.created) and not the approval date. The request can be approved long before the job exists. The entry has kind 'new_job'.

## Date of a value change is its approval date

Lines 12-13

Asked as: value change approval date; date of changed value; approved value change date; project value update date

For an approval that is not initial, a later value change, Turbo uses the date it was approved (approved). The entry has kind 'value_change'.

## Which approvals appear in the approved-by-me list

Lines 14-18

Asked as: approvals I signed off; which approvals show for me; my approved project values; approved by me missing

The list shows only project values (job attachments) that the signed-in person approved as the approver. Each must have status 'Approved' and must be attached to a job. Approvals by other people do not appear.

## List order is newest first

Lines 20-23

Asked as: approved list sort order; newest approvals first; order of approved items; approved list ordering

Turbo sorts the entries by date, newest first, and by id when dates are equal. Each entry shows the job, job number, job description, project value, kind and date.
