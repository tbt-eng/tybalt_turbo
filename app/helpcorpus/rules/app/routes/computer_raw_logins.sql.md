<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_raw_logins.sql","sha256":"f9bbeba15bf9eefb7c72d2a0bcb88844945ed21c41dcfc52d24a205ba99641ca","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_raw_logins.sql

## Raw Logins list leaves out the report

Lines 1-10

Asked as: raw logins list slow; raw login report missing; raw logins page loads slowly; why report column empty

The Raw Logins list shows only the id, the created date, the computer name and the error of each raw login. It leaves out the report, which can be 64 KB. The page loads the report for one row at a time.

## Raw Logins page lists every raw login, newest first

Lines 5-15

Asked as: raw logins order; newest raw logins first; raw logins list all rows; raw login page no limit

The Raw Logins page lists every raw login (computer_raw_logins) from the newest to the oldest, ordered by created date and then by id. The list has no limit on the number of rows.

## Missing computer name or error shows as blank

Lines 8-14

Asked as: raw login blank computer name; raw login null error; empty error column raw login; raw login missing computer name

A raw login with no computer name or no error shows an empty text in that column and not a null value. A search treats these empty values as empty text.

## Raw Logins search matches computer name or error

Lines 11-14

Asked as: search raw logins partial match; raw login search wildcard; raw logins search percent sign; raw login search not finding; raw logins LIKE pattern

The search on the Raw Logins page matches a raw login when its computer name or its error fits the LIKE pattern. The match uses the pattern as given, so a person must add % for a partial match. A backslash escapes the wildcard characters. A blank search shows every raw login.
