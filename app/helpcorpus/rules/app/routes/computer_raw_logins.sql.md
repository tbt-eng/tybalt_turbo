<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_raw_logins.sql","sha256":"b1a9002f8ad910ebf0666b15d4985dcd4e833c74f2cc3a8d43379774028b1ef2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_raw_logins.sql

## Raw Logins list leaves out the report itself

Lines 6-11

Asked as: raw logins list missing report; report not shown in list; list does not show report body; why is the report not in the list; raw login list slow to load

Each row shows only its id, created time, computer name and error. The report is not in the list because it can be 64 KB. The page loads the report for one row at a time.

## Missing computer name or error shows as blank

Lines 9-10

Asked as: raw login blank computer name; empty error column raw login; null computer name shows blank; no error shown on raw login; raw login missing error text

When a raw login has no computer name (computer_name) or no error, the list shows an empty text for it instead of a null value.

## Raw Logins search matches computer name or error

Lines 13-15

Asked as: search raw logins by computer name; raw logins search not finding; search raw login error text; wildcard search raw logins; raw login search returns nothing

The search on the Raw Logins page is a LIKE pattern. A raw login is listed when its computer name or its error matches the pattern. A backslash escapes the wildcard characters in the pattern. A blank search lists every row.

## Raw Logins page lists newest reports first

Lines 16-17

Asked as: raw logins sorted order; newest raw logins first; raw logins page order; raw login paging limit offset; rejected computer reports list order

The Raw Logins page lists raw logins (rejected computer reports) ordered by created date, newest first. Rows with the same created time are ordered by id, highest first. The page shows one page of rows at a time, set by the limit and offset.
