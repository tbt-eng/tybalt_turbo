<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/job_request_allocations.md","sha256":"e565e63c9cf4965786e94e0ad42b9abf1492e3614bb0504bc83467438e1a8a60","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/job_request_allocations.md

## Who can list job request allocations

Lines 9-11

Asked as: who can see allocations; list job request allocations; see allocations for project requests; why can't I see allocations; allocations access for approvers

A signed-in person with an active admin profile (admin_profiles active = true) can list allocations. They must also be the creator of the project request (request.creator), or hold the job claim (procurement). Or they must be the request's approver (request.approver) and hold the project_authorization claim. Everyone else gets nothing.

## Who can view one job request allocation

Lines 15-17

Asked as: open one allocation; view single allocation; can't open allocation record; allocation details not showing; see one allocation

A signed-in person with an active admin profile (admin_profiles active = true) can view one allocation. They must also be the creator of the project request (request.creator), or hold the job claim (procurement). Or they must be the request's approver (request.approver) and hold the project_authorization claim. This is the same filter as for listing.

## Nobody creates job request allocations through the API

Lines 19-21

Asked as: create allocation via API; can't add new allocation; add allocation to request; API allocation creation not working; make a new allocation

No rule allows creating an allocation through the API, so only a PocketBase superuser can. Turbo's own server code can still create them.

## Nobody updates job request allocations through the API

Lines 23-25

Asked as: edit allocation through API; change allocation hours; update allocation not saving; allocation edit denied; modify allocation record

No rule allows updating an allocation through the API, so only a PocketBase superuser can. Turbo's own server code can still update them.

## Nobody deletes job request allocations through the API

Lines 27-29

Asked as: delete job request allocation; remove allocation; allocations gone after deleting request; can't delete allocation; erase allocation record

No rule allows deleting an allocation through the API, so only a PocketBase superuser can. Turbo's own server code can still delete them. Deleting a project request (job_requests) also deletes its allocations (cascadeDelete).

## One allocation per division on a project request

Lines 38-46

Asked as: duplicate allocation for division; two allocations same division; division already allocated; hours not accepting decimals; allocation division error

A unique index (idx_request_division) means a project request can have only one allocation for each division. The request and the division are required. The hours field takes whole numbers only (onlyInt).
