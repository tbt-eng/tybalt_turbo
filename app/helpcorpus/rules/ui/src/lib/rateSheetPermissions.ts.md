<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetPermissions.ts","sha256":"a2dccd41c1cbc0cc0a902deaa0d98ce6b42995724731ee9e64581593a7a8ec18","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetPermissions.ts

## Who can create or copy a rate sheet

Lines 5-9

Asked as: create rate sheet; copy a rate sheet; who can make a new rate sheet; add entry to rate sheet

Procurement (the job claim) can create a new rate sheet or a copy (a request with revision 0). The same claim can add a new entry to an incomplete rate sheet. Others do not get the control.

## Who can activate or deactivate a rate sheet

Lines 11-15

Asked as: activate rate sheet; deactivate rate sheet; turn off rate sheet; switch rate sheet on or off

Procurement (the job claim) can activate and deactivate a rate sheet.

## Who can revise or rename a rate sheet

Lines 17-21

Asked as: revise rate sheet; rename rate sheet; change rate sheet name; who can edit a rate sheet

An admin (the admin claim) or a person with the rate sheet revise claim (rate_sheet_revise) can make a revision of a rate sheet (a request with a revision above 0). The same people can rename a rate sheet.

## Who can correct entry rates or roll out a revision

Lines 23-28

Asked as: correct entry rate; fix rate on rate sheet; roll out revision; rollout refused for my user

Only an admin (the admin claim) can correct the rates of a rate sheet entry. Only an admin can move projects to a revision (rollout). The server also refuses a rollout by an inactive user.

## Who can delete a rate sheet

Lines 30-34

Asked as: delete rate sheet; remove rate sheet; cannot delete rate sheet; rate sheet delete button missing

Only an admin (the admin claim) can delete a rate sheet. The sheet must be inactive and nothing may use it. The server checks the sheet.
