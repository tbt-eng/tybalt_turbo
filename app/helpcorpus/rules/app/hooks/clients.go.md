<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/clients.go","sha256":"7e0da5737210f8817b241385d7b03194aef275be09205c566ac1fe8df9adf1ff","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/clients.go

## Business development lead is optional on a client

Lines 14-20

Asked as: client without busdev lead; leave lead empty; optional business development lead; no lead on client

A client can be saved with no business development lead (business_development_lead). When the field is empty, Turbo skips every lead check.

## Lead is checked only when set or changed

Lines 14-23

Asked as: old lead lost claim save; save client lead deactivated; existing client lead unchanged; can't save client old lead

On an existing client, Turbo checks the business development lead only if the lead is different from the saved one. A client whose lead has since lost the busdev claim or been deactivated can still be saved while the lead stays the same. A new client with a lead is always checked.

## Business development lead must be an active user

Lines 25-44

Asked as: lead inactive user error; lead not active; deactivated lead client save; busdev lead must be active

When a lead is set or changed on a client, the lead must be an active user. Otherwise Turbo refuses the save with the message "business development lead must be an active user". The field error is "the selected business development lead is not an active user" (business_development_lead_not_active).

## Lead check failures return a server error message

Lines 26-58

Asked as: lead status check failed; error validating lead claim; can't verify lead claim; server error saving client lead

If Turbo cannot check the lead's active status, it refuses the save with "failed to check business development lead active status". If it cannot check the busdev claim, it refuses with "error validating business development lead claim". The field error for the claim failure is "unable to verify business development lead claim" (claim_check_failed).

## Business development lead must hold the busdev claim

Lines 46-71

Asked as: lead lacks busdev role; invalid business development lead; lead not allowed busdev; wrong person as busdev lead

When a lead is set or changed on a client, the lead must have the busdev claim. Otherwise Turbo refuses the save with the message "invalid business development lead". The field error is "business development lead must have busdev claim" (missing_claim). This check runs after the active-user check.

## Client create and update run the lead validation

Lines 76-85

Asked as: client save validation runs; creating client lead check; updating client validation; when lead validation happens

Turbo runs the business development lead validation (validateClient) whenever a client is created or updated (ProcessClient). It applies no other client checks in this file.
