<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_invoicing.go","sha256":"3125f9da05a7279051ea21e4d5cb8b44fcc8ca3b122108c413331a656d8a432a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_invoicing.go

## Who can set the invoicing profile of a project request

Lines 19-35

Asked as: who can change invoicing profile; set invoicing profile on request; procurement invoicing profile permission; can't change invoicing profile

Only a holder of the job claim (procurement) can set the invoicing profile of a project request. The request must have the status "Awaiting creation". Project requests (approval required) and job editing must both be on. Without a signed-in user, the answer is no.

## Invoicing change keeps status, version and approval

Lines 37-43

Asked as: clear invoicing profile; does profile change reset approval; invoicing change lose approval; approval still valid after profile edit

Changing or clearing the invoicing profile (an empty value clears it) changes only the profile. The status, the version and the approval stay. The review digest leaves the profile out, so the approval still covers the request.

## Invoicing profile change needs project requests and job editing on

Lines 51-65

Asked as: invoicing profile greyed out; jobs editing disabled error; can't edit profile project request; profile change refused job editing off

Turbo refuses the change when project requests are off, or when the user lacks the job claim. It also refuses when job editing is off, because the profile becomes the job's and follows the jobs editing switch. In that case Turbo returns the jobs-editing-disabled error (ErrJobsEditingDisabled).

## Invoicing profile changes only while awaiting job creation

Lines 66-72

Asked as: invoicing profile locked; change billing profile after created; procurement can't edit invoicing once approved; job create queue invoicing edit

Procurement can change a request's invoicing profile only when its status is "Awaiting creation". Otherwise Turbo says: "Procurement can change the invoicing profile only while the request is awaiting job creation." Procurement can do this on anyone's request in the Job Create Queue.

## Invoicing profile change checks the request version

Lines 73-75

Asked as: invoicing change version out of date; request version conflict; someone else edited request; stale request version error

The change must carry the request's current version. Turbo checks it with requestVersion and refuses the change if the version is out of date.

## Setting the same invoicing profile changes nothing

Lines 76-78

Asked as: save same invoicing profile; nothing happens when saving profile; invoicing profile unchanged; same profile selected again

When the submitted profile equals the one already on the request, Turbo makes no change and saves nothing. It simply returns the request.

## Old approval is upgraded before a profile change

Lines 79-123

Asked as: approval went stale after profile change; old approval digest; approval outdated after update; profile edit makes approval stale

Before a profile change, Turbo upgrades an "Approved" approval saved with an older review digest. It does so only when the old digest still matches the request's data. Turbo then stores the current digest, so the change does not make the approval stale. A saved digest that matches no older form stays stale, and so does an empty one.

## Profile must fit the client and billing contact email

Lines 85-91

Asked as: invoicing profile wrong client; billing contact missing email; profile rejected when saving; billing contact no email error

A new profile must belong to the request's client. Its billing contact needs an email address or the no-email mark, as job creation checks. Otherwise Turbo refuses the change and does not save it.

## Server start upgrades approvals of requests awaiting creation

Lines 125-152

Asked as: approval stale after deploy; server restart approvals; old approval after upgrade; job creation digest mismatch after deploy

When the server starts, Turbo upgrades the saved approval of every request that waits for job creation. A profile edit after a deploy therefore does not make an earlier approval stale. If the upgrade fails, Turbo logs the error and keeps running. Job creation still accepts an older matching digest.
