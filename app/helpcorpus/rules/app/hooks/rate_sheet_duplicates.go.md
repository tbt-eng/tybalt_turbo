<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_duplicates.go","sha256":"8c1eab0031646fe4f9dcd3417c1190d58b4c49e23d11bdc7560b4e081b579bba","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_duplicates.go

## Near-match limits for similar rate sheets

Lines 15-20

Asked as: how many roles can differ; similar rate sheet limit; rate sheets too different; near match threshold

A rate sheet counts as a close match only if it differs from the entries in at most 5 roles (RateSheetNearMatchMaxRoles). Turbo shows at most 3 close matches (RateSheetNearMatchLimit).

## How differing roles are counted between rate sheets

Lines 30-39

Asked as: why role counted as different; role missing on one side; overtime rate difference counts; how roles get compared

Turbo counts a role as differing when its rate or overtime rate is different. It also counts a role that is on one side only (DifferingRoles). An exact match has 0 differing roles.

## Close matches are listed closest first and trimmed

Lines 49-63

Asked as: why only three matches shown; similar sheets list cut off; closest rate sheets list; missing similar sheets

The list of close rate sheets (Closest) stops at the given limit. It also stops at the first sheet that differs in more roles than the maximum. The list is sorted by fewest differing roles, so no later sheet is closer.

## A rate sheet with no entries still gets compared

Lines 65-69

Asked as: rate sheet with no roles; empty rate sheet duplicate; blank rate sheet check; rate sheet without entries

An active rate sheet with no entries is still compared, and it has no roles. It is an exact match only when the compared entries are also empty.

## Rates are compared in whole cents

Lines 77-89

Asked as: tiny rate difference ignored; fraction of a cent; rates differ by pennies; rounding rate differences

Turbo rounds the rate and the overtime rate to whole cents before it compares them. A difference smaller than rounding to the cent does not count.

## Only active rate sheets are compared, except the named one

Lines 91-120

Asked as: inactive sheet not checked; old revision not compared; why sheet skipped in check; archived rate sheet matches

The comparison covers active rate sheets only (active = true). It skips every sheet with the excluded name, so a revision is never compared with its own earlier revisions. Inactive sheets are never compared.

## Exact match means same roles and same rates

Lines 151-169

Asked as: same rates different roles; rate sheets look identical; identical rate sheet detection; same roles and rates

An active rate sheet is an exact match when it has the same set of roles and every role has an equal rate and overtime rate. Any other active sheet goes in the near list.

## Order of matching rate sheets

Lines 176-187

Asked as: order of similar sheets; which match listed first; sorting rate sheet matches; tie between rate sheets

Turbo sorts matches by fewest differing roles first. Ties go by sheet name, then by revision number, lowest first.

## Duplicate rate sheet error message

Lines 203-223

Asked as: rate sheet already exists error; identical rates refused; use existing rate sheet; duplicate rate sheet message; cannot save same rates

When a rate sheet has rates identical to an active sheet, Turbo refuses it with "a rate sheet with identical rates is already active". The field error says "the rates are identical to active rate sheet "<name>" (rev. <n>); use that sheet instead" and names the first match (code duplicate_rate_sheet). The error lists every match so the UI can link to the existing sheet.
