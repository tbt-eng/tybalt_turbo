<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/job_requests.md","sha256":"52952ca416e35e529a50c6a0d06d2ccd56488ce9685a96050bd74a41986873e2","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/job_requests.md

## Who can list or view project requests

Lines 9-17

Asked as: who can see requests; open project request access; list job requests permission; who may view requests; which staff see project requests

Only a signed-in person with an active admin profile (admin_profiles_via_uid.active) can list or view a project request (job_requests record). They must also be one of three kinds of person. The first is the person who created the request (creator). The second is procurement (the job claim). The third is the request's chosen approver (approver) who also holds the project authorization claim (project_authorization).

## Inactive admin profile blocks project request access

Lines 10-10

Asked as: deactivated admin cannot see request; my request disappeared; disabled account no request access; why can't I open my request; not signed in view request

A person whose admin profile is not active cannot list or view any project request, even one they created. A person with no signed-in session cannot either.

## Approver needs the project authorization claim to see request

Lines 10-10

Asked as: named approver cannot see request; approver can't open project request; approver missing authorization; approver sees nothing on request; approver access denied

Being named as the approver on a project request is not enough to list or view it. The approver must also hold the project authorization claim (project_authorization).

## Nobody creates, updates or deletes project requests through the API

Lines 19-29

Asked as: cannot edit job request api; can't delete project request; create request through api; update request not working; api request changes refused

The API has no rule to create, update or delete a project request (job_requests record), so nobody can do these actions through it except a PocketBase superuser. Turbo's own server code can still do them.

## Project request creator is required

Lines 38-38

Asked as: request has no creator; missing creator field; who created the request; creator required error; project request without owner

Every project request must have a creator (creator), who is one user.

## Project request status values

Lines 39-40

Asked as: allowed request statuses; project status options; valid draft status values; request status choices; awaiting approval status meaning

A project request's status (status) is required and is one of Draft, Awaiting approval, Awaiting creation, Rejected or Recalled. Its project status (project_status) is optional and is Active or Closed.

## Project value and version must be whole numbers

Lines 41-65

Asked as: project value decimal error; version must be integer; project value cents rejected; whole number only value; version number invalid

A project request's project value (project_value) and version (version) accept whole numbers only.

## Submission token format

Lines 42-42

Asked as: submission token format; token length limit; invalid submission code; submission token pattern; token must be 32

A project request's submission token (submission_token) is at most 32 characters. It must match the pattern of exactly 32 lowercase letters and digits.

## Project request notes length limits

Lines 69-69

Asked as: notes too short error; notes character limit; request notes minimum length; notes too long message; how long notes must be

The notes field (notes) on a project request must be 10 to 1000 characters long when it has a value.
