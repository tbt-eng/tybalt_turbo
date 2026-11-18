<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetPermissions.ts","sha256":"468d3754952c2f15ea3c905256649c27efb68e30c665169de10abd6f1fa15f56","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetPermissions.ts

## Who can create or copy a rate sheet

Lines 5-9

Asked as: create rate sheet; copy a rate sheet; can't make new rate sheet; add entry to incomplete sheet

Only a user with the job claim (procurement) can create a new rate sheet or copy one (a rate sheet with revision 0). The same claim lets a user add a new entry to an incomplete rate sheet. The pages show these controls only to users with the job claim.

## Who can activate or deactivate a rate sheet

Lines 11-15

Asked as: activate rate sheet; deactivate rate sheet; turn off rate sheet; switch on the sheet

Only a user with the job claim (procurement) can activate or deactivate a rate sheet. The pages show these controls only to users with that claim.

## Who can revise or rename a rate sheet

Lines 17-21

Asked as: revise rate sheet; rename rate sheet; can't rename sheet; revision button missing

A user with the admin claim or the rate_sheet_revise claim can revise a rate sheet (a revision above 0) and rename it. Procurement (the job claim) alone cannot do this.

## Who can correct entry rates or roll out a revision

Lines 23-28

Asked as: correct entry rates; fix rate on entry; roll out revision; push revision to projects

Only a user with the admin claim can correct the rates of a rate sheet entry. Only an admin can move projects to a revision (rollout). The server also refuses a rollout by an inactive user, and this file does not check that.
