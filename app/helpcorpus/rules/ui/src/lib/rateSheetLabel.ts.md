<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetLabel.ts","sha256":"31d391ac3900ee53e1208dc000b16d957ae3cfc71dc6860277be50567d64a43c","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetLabel.ts

## Rate sheet code and revision label shows name if no code

Lines 12-17

Asked as: rate sheet shows name instead of code; rate sheet code missing; blank rate sheet code; rev. 0 rate sheet label; rate sheet revision shows 0

Turbo labels a rate sheet as its code followed by the revision, such as "2026-11-00004 (rev. 2)". A sheet with no code, or a blank code, shows its name in place of the code. A missing revision shows as 0 (rev. 0).

## Rate sheet full text label adds name in quotes

Lines 19-28

Asked as: rate sheet name in quotes; full rate sheet label dropdown; rate sheet name missing from dropdown; rate sheet button text name; quoted name not showing

The full text label of a rate sheet is the code and revision, then the name in quotes, such as 2026-11-00004 (rev. 2) "Client Rates". Turbo adds no quoted name when the name is empty. It also adds none when the sheet has no code, because the name already stands in the code's place. Dropdown options and buttons use this label.

## Fields that a rate sheet search matches

Lines 30-40

Asked as: find rate sheet by client; search rate sheet by job number; rate sheet not found by name; old revision of rate sheet search; find rate sheet from job

A search for a rate sheet looks in its code, name, effective date, job numbers and client names (rate_sheets_augmented fields). A job number or client name finds every revision of the code that the job used, now or in its rate sheet history.

## Rate sheet list search text contents

Lines 42-57

Asked as: rate sheet list search results; search rate sheet by date; rate sheet search with no client; rate sheet list filter text; find rate sheet by effective date

The text that a rate sheet list search matches is the full label (code, revision and name), the effective date, the job numbers and the client names. A missing effective date, job numbers or client names counts as empty text.
