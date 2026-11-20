<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_raw_logins_count.sql","sha256":"feb3b2d9cfb168dd80e56b7b3f780540297c9818141e3336002d1a7148c22d38","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_raw_logins_count.sql

## Raw login count: empty search counts every raw login

Lines 3-6

Asked as: raw logins total count; how many raw logins; count of all raw logins; raw login list count is wrong

The raw login count (computer_raw_logins) counts every row when the search text is empty. The count is the total that the raw logins list pages through.

## Raw login search matches computer name or error text

Lines 6-8

Asked as: search raw logins by computer name; raw login search not finding computer; filter raw logins by error text; raw logins search matches nothing; search for percent sign in logins

When search text is given, Turbo counts only the raw logins whose computer name (computer_name) or error message (error) matches the text with LIKE. A raw login with an empty name or error is treated as empty text, so it can still be counted. The backslash is the escape character, so the search can match a literal % or _ when escaped.
