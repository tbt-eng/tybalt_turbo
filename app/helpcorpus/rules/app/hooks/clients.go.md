<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/clients.go","sha256":"1d691678ceb92f2896c8e85192ab80c4a90c0855227688673c1a4758a9c1ad34","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/clients.go

## Business development lead is optional on a client

Lines 14-17

Asked as: client without busdev lead; leave lead blank; lead field optional; no business development lead

A client saves without a business development lead (business_development_lead). Turbo skips all lead checks when the field is empty.

## Business development lead must be an active user

Lines 19-38

Asked as: inactive lead on client; lead account disabled; can't save client lead; deactivated user as lead

When a client has a business development lead, Turbo checks that the lead is an active user. If not, it refuses the save with "business development lead must be an active user". The field error is "the selected business development lead is not an active user" (business_development_lead_not_active).

## Lead check failures give an error message

Lines 20-52

Asked as: lead check failed; can't verify lead status; claim validation error; couldn't check business development lead

If Turbo cannot check the lead's active status, it refuses the save with "failed to check business development lead active status". If it cannot check the lead's claim, it refuses with "error validating business development lead claim". The field error is "unable to verify business development lead claim" (claim_check_failed).

## Business development lead must hold the busdev claim

Lines 40-65

Asked as: lead lacks busdev permission; wrong role for client lead; invalid business development lead; user not busdev claim

A client's business development lead must have the busdev claim. Otherwise Turbo refuses the save with "invalid business development lead". The field error is "business development lead must have busdev claim" (missing_claim).

## Lead checks run when a client is created or updated

Lines 70-79

Asked as: lead checks on client update; when lead validation runs; client save validation order; processclient lead check

Turbo applies the business development lead checks every time a client is created or updated (ProcessClient). The active check runs before the claim check.
