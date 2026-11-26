<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_invoicing.go","sha256":"e569031a61fb5e7a6bb966ae46f8551223c82b2e4bdd22887d4320feed824d5e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_invoicing.go

## Who can set a request's invoicing profile

Lines 19-39

Asked as: who can set invoicing profile; invoicing profile permission; can't edit invoicing details; billing profile access

Only an active holder of the job claim (procurement) can set the invoicing profile of a project request. The request must have status "Awaiting creation". Project requests (approval) and job editing must both be switched on.

## Procurement can change invoicing profile on anyone's request

Lines 41-47

Asked as: change someone else's request billing; clear invoicing profile; procurement edit other request; remove billing details from request

Procurement can set, change or clear the invoicing profile (invoicing_information) of any project request in the Job Create Queue, not only its own. An empty profile clears it.

## Invoicing profile change keeps status and approval

Lines 41-47

Asked as: does changing billing reset approval; approval lost after billing change; status changes after invoicing edit; keep approval when editing profile

Only the invoicing profile changes. The status, the version and the approval stay as they were, because the review digest leaves the profile out. The approval still covers the request.

## Invoicing profile change needs project requests and job editing on

Lines 55-69

Asked as: invoicing profile not saving; billing change blocked jobs editing; job editing disabled error; can't change billing when jobs off

Turbo checks that project requests are enabled and that the person holds the job claim. It also refuses the change when job editing is switched off, because the profile becomes the job's. The error for that case is utilities.ErrJobsEditingDisabled.

## Invoicing profile changes only while awaiting job creation

Lines 74-76

Asked as: invoicing profile locked; can't edit billing after approval; billing change refused status; awaiting job creation billing

If the request's status is not "Awaiting creation", Turbo refuses with: "Procurement can change the invoicing profile only while the request is awaiting job creation."

## Request version must match when changing invoicing profile

Lines 77-79

Asked as: version mismatch invoicing save; request changed by someone else; conflict when saving billing; someone else edited request

Turbo checks the version sent against the request's version (requestVersion). A mismatch is refused. The version does not change after a successful edit.

## Same invoicing profile makes no change

Lines 80-82

Asked as: save with no changes; billing unchanged nothing saved; same invoicing profile selected; saving same billing details

When the profile sent equals the request's current profile, Turbo returns the request detail and saves nothing.

## Old approval gets the current digest before profile change

Lines 83-127

Asked as: approval out of date after update; stale approval digest; approval shows outdated; approval not matching request data

Before a profile change, Turbo upgrades an "Approved" approval saved with an older review digest, only if that digest still matches the request's data. A saved digest that matches no older form is left alone, so a stale approval stays stale. An approval with no saved digest is also left alone.

## Invoicing profile must fit the client and billing contact

Lines 89-96

Asked as: billing contact has no email; wrong client invoicing profile; profile rejected on save; billing contact missing email error

After the change, Turbo checks that the profile belongs to the request's client (ValidateClientReferences). It also checks that the profile's billing contact has an email address or the no-email mark (ValidateRequestContactEmails), as job creation does. A failure rejects the change.

## Approvals upgraded when the server starts

Lines 129-156

Asked as: approvals after server restart; server deploy approval stale; approval broken after update; startup approval upgrade

At server start, Turbo upgrades the saved approval of every request awaiting job creation, so a later profile edit does not make an approval saved before a deploy stale. A failure is logged and does not stop the server. Job creation still accepts an older digest that matches.
