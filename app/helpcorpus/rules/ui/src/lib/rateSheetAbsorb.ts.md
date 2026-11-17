<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetAbsorb.ts","sha256":"279118536ed95e5adbc84bbb30dc96ad5cf91926819f6750a663273dfa235ec5","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetAbsorb.ts

## Why a sheet cannot be absorbed into the target

Lines 1-41

Asked as: why can't I absorb this sheet; absorb button greyed out; reason shown for rate sheet; sheet cannot be merged into target

The server decides which rate sheets can be absorbed and says why not in a reason. A duplicate sheet shows a reason when the group's target cannot absorb it. A group shows a reason when no sheet in it can be absorbed into the target. An absorb candidate has no reason when the target can absorb it. The UI only shows these reasons.

## Target of a duplicate rate sheet group

Lines 19-22

Asked as: which rate sheet is the target; earliest effective date duplicate; which duplicate sheet to keep; main sheet in duplicate group

Each group of duplicate rate sheets has a target (target). The target is the sheet with the earliest effective date (effective_date). The other sheets of the group are the ones that can be absorbed into it.

## How many sheets an absorb deletes

Lines 26-28

Asked as: how many rate sheets get deleted; does absorb delete all revisions; absorb removes every revision; how many sheets does absorb remove

The absorbs count says how many rate sheets an absorb into the target deletes. It is 1, or it is every revision of the sheet's name.

## Who can absorb rate sheets

Lines 43-49

Asked as: who is allowed to absorb rate sheets; need admin to merge sheets; absorb permission missing; rate_sheet_revise and absorb claim

A person can absorb rate sheets if they have the admin claim. A person who has both the absorb claim and the rate_sheet_revise claim can also do it. The absorb claim alone is not enough. The UI checks this with canAbsorbRateSheets, and the file says the server checks the same claims.
