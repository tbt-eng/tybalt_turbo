<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetLabel.ts","sha256":"9a8895f9139b9b56ee29048287f455eacdbc99fdcdb1a97c5e34a8a7287b72e8","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetLabel.ts

## Rate sheet label shows code and revision

Lines 12-17

Asked as: rate sheet code and revision; rev number on rate sheet; rate sheet shows rev 0; missing revision on rate sheet

Turbo labels a rate sheet with its code and revision, such as "2026-11-00004 (rev. 2)". Turbo treats a missing revision as 0 and shows "(rev. 0)" (rateSheetCodeRevision).

## Unsaved rate sheet shows its name instead of a code

Lines 12-17

Asked as: rate sheet no code; rate sheet shows name not code; blank code on rate sheet; unsaved rate sheet label

When a rate sheet has no code, or its code is only spaces, Turbo shows its name in the code's place, followed by the revision. A sheet with neither code nor name shows only " (rev. N)" with no text before it.

## Full rate sheet label adds the name in quotes

Lines 19-28

Asked as: rate sheet name in quotes; full rate sheet label; rate sheet dropdown option text; rate sheet button label

The plain-text label of a rate sheet is its code and revision, then the name in quotes, such as 2026-11-00004 (rev. 2) "Client Rates". Dropdown options and buttons use this label (rateSheetLabel).

## Rate sheet label omits the name when name or code is missing

Lines 23-28

Asked as: rate sheet name repeated; rate sheet label without name; rate sheet no name shown; name appears twice rate sheet

When a rate sheet has no name, or no code, its plain-text label is only the code-and-revision form with no quoted name. A sheet without a code shows its name once, in the code's place, and does not repeat it in quotes.

## Rate sheet list search matches label and effective date

Lines 30-36

Asked as: search rate sheets by date; find rate sheet by effective date; rate sheet list search not working; search rate sheet by name

A search in a rate sheet list matches the full label (code, revision and name) and the effective date (effective_date). A missing effective date adds nothing to the searched text.
