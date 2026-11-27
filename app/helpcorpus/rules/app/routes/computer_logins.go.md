<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_logins.go","sha256":"8a2fcc1e096445f09e328c6995ba8e49023feddf4ebf42f6df0d59860857b37b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_logins.go

## Logins list shows every sign-in, newest first, optional filters

Lines 28-81

Asked as: login history; sign-in list; filter logins by computer; who logged into a computer; logins page

The Logins page lists all sign-ins in one array, with no paging. A computer filter (?computer=) limits the list to one computer, and a search (?q=) narrows it further. Turbo trims spaces around the computer value. Each row includes the matched user's name (user_name), which a holder of the computers claim cannot read through the API. The code comment says Turbo deletes sign-ins after 30 days, but this file does not do that.

## Raw logins list leaves out the report

Lines 43-96

Asked as: rejected logins; failed login reports; raw logins page; errors from computers; login errors list

The Raw Logins page lists every rejected report in one array, newest first. Each row shows the id, created time, computer name and error, but not the report. A search (?q=) narrows the list.

## Search text matches anywhere and treats % _ \ literally

Lines 53-62

Asked as: search logins; percent sign in search; underscore in search; wildcard search; search not matching

The search text (?q=) on the Logins and Raw Logins lists matches anywhere in the searched values. Turbo trims spaces around it first. The characters %, _ and \ count as themselves, not as wildcards. Blank search text gives an empty pattern (computerSearchPattern), so the code applies no pattern.

## Viewing one raw login shows its report, or null

Lines 98-123

Asked as: view rejected report; raw login shows null; empty report; open raw login; raw login not found

Opening one raw login returns the row with the report as it was received. When the report is empty or missing, Turbo shows null. When the raw login does not exist, Turbo answers "Raw login not found."

## Deleting a raw login removes only that record

Lines 125-145

Asked as: delete rejected login; remove failed report; delete raw login; clear login error; remove raw login record

Deleting a raw login removes that one rejected report and changes nothing else, because nothing refers to a raw login. When the raw login does not exist, Turbo answers "Raw login not found."
