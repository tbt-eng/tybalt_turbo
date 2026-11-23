<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/clients.go","sha256":"fb12b393a76bd1e783dd5da3f075f06bbb9796cb7982d83256f61ae5a709e451","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/clients.go

## Business development lead on a client is optional

Lines 14-20

Asked as: client without lead; leave lead empty; lead field optional; save client no lead

A client (client record) can be saved with no business development lead (business_development_lead). Turbo then skips the lead check. The UI does not offer the lead for now.

## Lead is checked only when set or changed

Lines 14-23

Asked as: deactivated lead client save; inactive lead still saved; old lead error on edit; lead check when editing

On an existing client, Turbo checks the business development lead only if the lead differs from the saved value. A client whose lead has since been deactivated can still be saved with that same lead. A new client with a lead is always checked.

## Business development lead must be an active user

Lines 25-44

Asked as: lead not active user; inactive user as client lead; client lead must be active; business development lead error

When a client is created, or its lead is changed, the lead must be an active user. Otherwise Turbo refuses with "business development lead must be an active user". The field error says "the selected business development lead is not an active user" (code business_development_lead_not_active).

## Client save fails if the lead's active status cannot be checked

Lines 26-32

Asked as: failed to check lead status; client save lead check failed; cannot verify lead active; lead status error saving client

If Turbo cannot check whether the business development lead is active, it refuses the save with "failed to check business development lead active status".
