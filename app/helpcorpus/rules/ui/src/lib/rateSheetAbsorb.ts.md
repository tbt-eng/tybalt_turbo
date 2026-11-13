<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetAbsorb.ts","sha256":"87d622beb8f6aba13178b4d36f759b3691d08ed54902be471456aeb35883f3ba","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetAbsorb.ts

## How many rate sheets an absorb deletes

Lines 14-40

Asked as: how many rate sheets get deleted; absorb deletes which sheets; revisions deleted by absorb; count of sheets removed

An absorb into the target deletes either 1 sheet or every revision of that sheet's name (absorbs). The same count appears for absorb candidates. A sheet's revisions count is the number of revisions of its name.

## Target of a duplicate rate sheet group

Lines 18-21

Asked as: which duplicate rate sheet is the target; earliest effective date duplicate; which sheet gets kept; duplicate group target sheet

In a group of duplicate rate sheets, the target is the sheet with the earliest effective date (target). Duplicates are absorbed into this target.

## Why a duplicate rate sheet cannot be absorbed

Lines 22-34

Asked as: why can't I absorb this sheet; absorb button greyed out reason; rate sheet cannot be absorbed; why duplicate sheet blocked

The server decides which duplicate sheets the target can absorb. For a sheet that it cannot absorb, it gives a reason (reason), and the UI only shows that reason. A group can also have a reason why no sheet in it can be absorbed into its target.

## Who can absorb rate sheets

Lines 42-48

Asked as: who is allowed to absorb rate sheets; need which claim to absorb; absorb permission missing; can't absorb rate sheets access

A person can absorb rate sheets if they have the admin claim. A person who has both the absorb claim and the rate_sheet_revise claim can also absorb them. Someone with only one of those two claims cannot. The server checks the same claims.
