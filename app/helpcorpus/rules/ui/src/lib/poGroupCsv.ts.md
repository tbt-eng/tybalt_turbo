<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poGroupCsv.ts","sha256":"0005ed72497240db8377465684889ccfda55ac83c3edb0fad57e056673e612fb","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poGroupCsv.ts

## Upload field hints for purchase order group CSV

Lines 5-16

Asked as: upload field hints; what format for total; description minimum length; division code example

The purchase order group upload offers four fields: total (an amount, such as 250 or 250.00), description (at least 4 characters), division (the division code, such as CI) and branch (the branch name or code). The hints only guide the user; this file does not check the values.

## CSV columns required for a purchase order group

Lines 18-24

Asked as: csv column order; which columns needed; email column missing; purchase order group csv headers

A purchase order group CSV always needs an email column, which picks the owner. It then needs one column for each upload field the creator chose. The columns come in a fixed order: total, description, division, branch.

## CSV parsing accepts quotes, line breaks, CRLF and byte order mark

Lines 26-66

Asked as: quoted commas in csv; csv line breaks in cell; windows line endings csv; byte order mark error

Turbo reads quoted cells that contain commas, doubled quotes and line breaks. It accepts both CRLF and LF line endings. It ignores a leading byte order mark. A final row without a trailing line break is still read.

## Empty CSV reports every column as missing

Lines 81-84

Asked as: empty csv upload; blank file all columns missing; csv with nothing in it; no rows imported

If the uploaded CSV has no non-blank line, Turbo returns no rows and lists all the needed columns (email plus the chosen fields) as missing.

## Header row is the first non-blank line; names ignore case and spaces

Lines 81-88

Asked as: header row not found; column name case sensitive; csv header spaces; column reported missing

Turbo takes the first non-blank line as the header row. It matches column names without regard to case or surrounding spaces. It reports each needed column that is not in the header as missing.

## CSV rows skip blank lines and ignore extra columns

Lines 90-101

Asked as: blank lines in csv; extra columns in csv; short row empty cells; csv cells trimmed

Turbo skips blank lines after the header. It keeps only the needed columns and ignores any others. It trims every cell. A needed column missing from the header, or a cell missing from a short row, becomes an empty string.

## Error line numbers count from the file's first line

Lines 91-101

Asked as: wrong line number error; which line is the error; csv error row number; line numbers off

Each row keeps its line number in the file, so errors can point at it. Blank lines and lines above the header count in the numbering. The first line of the file is line 1.

## New purchase order group id is made before submitting

Lines 105-111

Asked as: duplicate purchase order group; submit button created twice; retry purchase order group; group id generated

Turbo makes a random 15-character id of lowercase letters and digits, in the shape PocketBase uses, before it submits a purchase order group. A repeated request then finds the group it already created.
