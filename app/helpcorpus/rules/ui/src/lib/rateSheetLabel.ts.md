<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetLabel.ts","sha256":"7f779cadc954d55455392892b2d4f17cfaefb089ff5fbeba50a2ed94ef43a95b","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetLabel.ts

## Rate sheet code and revision label

Lines 14-19

Asked as: rate sheet code; rev. 0 label; no revision shown; sheet shows name not code

Turbo shows a rate sheet as its code and revision, such as "2026-11-00004 (rev. 2)". A sheet with no revision shows "rev. 0". A sheet with no code, such as one that is not saved, shows its name in place of the code.

## Full rate sheet label with the name in quotes

Lines 21-30

Asked as: full rate sheet label; name in quotes; sheet name missing from label; unsaved sheet label

The full label of a rate sheet is the code and revision, then the name in quotes, such as 2026-11-00004 (rev. 2) "Client Rates". Turbo leaves out the name when the sheet has no name. It also leaves out the name when the sheet has no code, because the name already stands in for the code.

## Fields that a rate sheet search finds a sheet by

Lines 32-45

Asked as: search rate sheets by date; find rate sheet by job; search box hint; find sheet by client name

A rate sheet search finds a sheet by its code, name, effective date, job numbers and client names (rate_sheets_augmented). The search box hint reads "Code, name, date, job or client".

## Rate Sheets list search text and line breaks

Lines 47-64

Asked as: search matches two job numbers; rate sheets search line break; search across jobs wrong match; search finds mixed clients

The Rate Sheets list matches a search against the label, the effective date, the job numbers and the client names. A line break separates them, so a search matches inside one of them and never spans two job numbers or two client names.

## Rate Sheets list match is a case-insensitive substring

Lines 76-83

Asked as: case sensitive rate search; partial job number search; search lowercase client name; empty search rate sheets

In the Rate Sheets list, a job number or client name entry counts as matched when it contains the search text, ignoring case. An empty search matches no entry.

## Picker prefers entries that hold every search word

Lines 85-111

Asked as: picker search too many results; 91-89 finds 91-8801; all words must match picker; exact job number in picker

In a rate sheet picker (MiniSearch), a sheet is kept when every search word starts a word of the sheet. If a job number or client name has every search word as the start of one of its words, only such entries match. For example, "91-89" names 91-8901 and not 91-8801.

## Picker falls back to entries with any matched word

Lines 85-112

Asked as: picker no exact match; search returns nothing in picker; partial word match rate sheet; 24-321 matches both jobs

In a rate sheet picker, when no entry has every search word, an entry matches if one of its words was matched by the search for that field. For example, "24-321" on a sheet with 24-100 and 25-321 names both. If the search has no result for the sheet, no entry matches.

## Match line shows why a search found a sheet

Lines 114-132

Asked as: why sheet matched search; match line empty; job and client shown under sheet; plus more matches count

The match line lists the matched job numbers as "Job N", then the matched client names as "Client Name", joined with " • ". It shows three at most, then "+N more" for the rest. It is empty when nothing matched, such as for an empty search or one that matched only the code, name or date.
