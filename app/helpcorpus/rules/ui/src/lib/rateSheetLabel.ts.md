<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/rateSheetLabel.ts","sha256":"45f109fe9974d96fd7d0cfb119080a5f5ee193e9d242ca9cb649c8c9e0974f04","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/rateSheetLabel.ts

## Rate sheet label shows code and revision

Lines 14-19

Asked as: rate sheet revision number; how rate sheet code looks; rev 0 meaning; sheet version shown

Turbo names a rate sheet by its code and revision, such as "2026-11-00004 (rev. 2)". A sheet with no revision shows "rev. 0".

## Unsaved rate sheet shows its name instead of a code

Lines 14-19

Asked as: rate sheet not saved name; no code on rate sheet; unsaved sheet shows name; rate sheet missing code

A rate sheet with no code yet, such as one that is not saved, shows its name in the place of the code. Turbo trims spaces from the code and the name first. If both are empty, only " (rev. N)" remains.

## Copying a rate sheet code gives the code alone

Lines 21-26

Asked as: copy rate sheet code; copied code only; copy button rate sheet; rate sheet code clipboard

The copy text of a rate sheet is its code only, without the revision or the name, trimmed of spaces. It is empty when the sheet has no code yet. It does not fall back to the name.

## Full rate sheet label adds the name in quotes

Lines 28-37

Asked as: full rate sheet name dropdown; rate sheet label with name; rate sheet quoted name; dropdown shows rate sheet title

The plain-text label of a rate sheet is the code and revision, then the name in quotes, such as 2026-11-00004 (rev. 2) "Client Rates". Dropdown options and buttons use it. It matches the server's label (hooks.RateSheetLabel).

## Name is left out of the label when the sheet has no code

Lines 32-37

Asked as: rate sheet without name; no quotes on rate sheet; rate sheet name missing label; sheet label missing name

A rate sheet with no name, or with no code, shows only the code and revision label without a quoted name. A sheet with no code already shows its name in place of the code.

## Fields that a rate sheet search finds a sheet by

Lines 39-49

Asked as: find rate sheet by job; search rate sheet by client; rate sheet search what works; find old rate sheet revision

A rate sheet search finds a sheet by its code, name, effective date, job numbers and client names (rateSheetSearchFields). A job number or client name finds every revision of the code that the job used, now or in its rate sheet history.

## Rate sheet search box hint

Lines 51-52

Asked as: rate sheet search box placeholder; what to type in rate sheet search; search box text rate sheet

The hint in a rate sheet search box reads "Code, name, date, job or client".

## Rate Sheets list search text and line breaks

Lines 54-71

Asked as: rate sheets list search not matching; search spans two entries; list search line break; rate sheets list search wrong results

The Rate Sheets list matches a search against the label (code, revision, name), the effective date, the job numbers and the client names. A line break separates them, so a search matches inside one entry and never spans two.

## Rate Sheets list matches by substring, ignoring case

Lines 83-90

Asked as: rate sheets list job search; case sensitive search rate sheet; partial job number search; empty search rate sheets list

In the Rate Sheets list, a job number or client name entry matches when it contains the search text, ignoring case. An empty search matches no entry.

## Pickers match search words by word start

Lines 92-96

Asked as: rate sheet picker search words; picker search not finding sheet; typing start of word picker; rate sheet dropdown search

In a rate sheet picker, MiniSearch splits the search and each field into words. It keeps a sheet when every search word starts a word of the sheet.

## Picker highlights only entries with every search word

Lines 96-98

Asked as: picker finds wrong job number; job number search too many results; exact job number picker; similar job numbers picker

When a job number or client name entry has every search word as the start of one of its words, only such entries match. "91-89" names 91-8901, not 91-8801.

## Picker falls back to entries with any matched word

Lines 97-119

Asked as: picker shows unrelated sheets; partial search picker results; no exact match picker; search one word matches sheet

When no entry has every search word, an entry matches when one of its words was matched by the search in that field. "24-321" on a sheet with 24-100 and 25-321 names both. With no search result, nothing matches.

## Why a search found a sheet: Jobs and Clients groups

Lines 130-154

Asked as: why sheet matched search; search result reason group; matched by job or client; search match explanation

Turbo shows why a search found a sheet as one group per kind that matched, jobs first, then clients. The label is singular (Job, Client) for one entry and plural (Jobs, Clients) for more.

## Match groups show three entries and a count of more

Lines 135-154

Asked as: more matches hidden; plus more matches count; only three results shown; match list truncated

Each match group shows the first three matching entries and the number of further matches, such as "+5 more". There are no groups when nothing matched, such as an empty search or a match only on the code, name or date.
