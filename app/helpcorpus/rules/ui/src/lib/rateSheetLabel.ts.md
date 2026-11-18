<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetLabel.ts","sha256":"1c3dba4f71bab504852fb3bbd9685d868b3037839eafd20b7376acc84794c278","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetLabel.ts

## Rate sheet code and revision label format

Lines 14-19

Asked as: rate sheet code format; rate sheet shows no code; missing revision rev 0; unsaved rate sheet label

Turbo shows a rate sheet as its code and revision, such as "2026-11-00004 (rev. 2)". A sheet with no code, such as one that is not saved, shows its name in place of the code. A missing revision shows as 0 (rev. 0).

## Full rate sheet label adds the name in quotes

Lines 21-30

Asked as: rate sheet name in quotes; rate sheet label missing name; full label with sheet name; why no name on sheet

The text label of a rate sheet is its code and revision, then the name in quotes, such as 2026-11-00004 (rev. 2) "Client Rates". Turbo leaves out the name when it is empty. Turbo also leaves it out when the sheet has no code, because the name already stands in for the code.

## Fields a rate sheet search looks in

Lines 32-45

Asked as: search rate sheet by client; find rate sheet by job number; rate sheet search fields; search box hint rate sheets

A rate sheet search finds a sheet by its code, name, effective date, job numbers and client names. A job number or client name finds every revision of the code that the job used, now or in its rate sheet history. The search box hint reads "Code, name, date, job or client".

## Rate Sheets list search text and line-break separation

Lines 47-64

Asked as: rate sheets list search matches; search spans two entries; line break in rate sheet search; rate sheet list search text

The Rate Sheets list matches the search against the label (code, revision, name), the effective date, the job numbers and the client names. A line break separates these parts. A search therefore matches inside one entry and never spans two entries.

## Rate Sheets list match is a case-insensitive substring

Lines 76-83

Asked as: case sensitive rate sheet search; search job number in list; client name list search not matching; empty search rate sheets list

In the Rate Sheets list, a job number or client name entry matches when it contains the search text, ignoring case. An empty search matches no entry.

## Picker search matches by word prefix

Lines 85-112

Asked as: picker search wrong job number; rate sheet picker matches start of word; search 91-89 shows 91-8801; picker only matching words

In a rate sheet picker, a sheet matches when every search word starts a word of the sheet. When some job number or client name entry starts words with every search word, only those entries are shown as matches. For example "91-89" names 91-8901, not 91-8801.

## Picker falls back to partial entry matches

Lines 85-112

Asked as: picker no exact match; partial match rate sheet picker; search finds nothing highlighted; rate sheet picker highlights some entries

When no job number or client name entry has all the search words, a picker marks each entry that has one of the matched words in that field. For example "24-321" on a sheet with jobs 24-100 and 25-321 names both. A search with no result marks no entries.

## Why-found line groups matches by kind

Lines 114-148

Asked as: why found line rate sheet; why does this sheet match; matched jobs and clients order; why found says job or jobs

For a found sheet, Turbo lists the matching job numbers first and then the matching client names, one group per kind. A group is labeled "Job" or "Client" for one match, and "Jobs" or "Clients" for more. It is empty for an empty search or one that matched only the code, name or date.

## Why-found line shows three entries and a count

Lines 133-147

Asked as: plus more matches count; why found shows only three; how many more matched sheet; too many matching jobs listed

Each group of matches in the why-found line shows the first three entries. The group also gives how many more matched, such as +5 more. Turbo shows no count when three or fewer matched.
