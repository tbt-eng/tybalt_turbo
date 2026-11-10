<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/expense_uncommit_requests.md","sha256":"5df93cfece981299b860c6212f85781bd66d393397a6d2c8b20d9d3a7077646c","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/expense_uncommit_requests.md

## Nobody can list expense uncommit requests through the API

Lines 7-9

Asked as: list uncommit requests API; cannot see any uncommit requests; fetch all expense uncommit records; API returns nothing for uncommit

The expense uncommit requests (expense_uncommit_requests) collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one expense uncommit request through the API

Lines 11-13

Asked as: view one uncommit request; open single uncommit record API; get uncommit request by id; cannot read uncommit request

The expense uncommit requests collection has no view rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create expense uncommit requests through the API

Lines 15-17

Asked as: create uncommit request via API; POST uncommit request forbidden; add new uncommit record API; API refuses new uncommit request

The expense uncommit requests collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update expense uncommit requests through the API

Lines 19-21

Asked as: edit uncommit request API; update uncommit record forbidden; change uncommit request status by API; PATCH uncommit request denied

The expense uncommit requests collection has no update rule. Nobody can change a record through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete expense uncommit requests through the API

Lines 23-25

Asked as: delete uncommit request API; remove uncommit record forbidden; erase uncommit request via API; cannot delete uncommit request

The expense uncommit requests collection has no delete rule. Nobody can delete a record through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Text length limits on uncommit requests

Lines 33-40

Asked as: uncommit note too long; character limit uncommit comment; decision comment too long; uncommit requester missing

The note (note) and the decision comment (decision_comment) of an uncommit request can each hold at most 1000 characters. The requester (requester) is required.

## Outcome and status accept only fixed values

Lines 36-37

Asked as: invalid uncommit outcome value; uncommit status not accepted; wrong status for uncommit request; uncommit outcome validation error

The outcome (outcome) of an uncommit request must be "uncommit_and_reject" or "uncommit_only". The status (status) must be "pending", "approved", "declined" or "cancelled". Both fields are required.

## Only one pending uncommit request per expense

Lines 59-59

Asked as: two pending uncommit requests; duplicate uncommit request same expense; already pending uncommit on expense; cannot submit second uncommit request

A unique index allows only one uncommit request with status "pending" for the same expense. A second pending request for that expense is refused. Requests that are approved, declined or cancelled do not count.
