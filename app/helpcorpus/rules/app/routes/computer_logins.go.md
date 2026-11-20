<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_logins.go","sha256":"a91e4c1dd0c8d88bb0085b25ba507e5146e96cfe91dada6d2eb7f833035aa030","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_logins.go

## Who can open the Logins and Raw Logins pages

Lines 1-6

Asked as: who can see logins; access to logins page; raw logins permission; can't open logins page

Turbo binds the Logins and Raw Logins pages and the raw login delete behind the same guard as the other computer routes. The guard is the computers claim or admin. This file does not show the guard itself.

## Login lists show 20 rows a page by default

Lines 36-88

Asked as: how many rows per page; logins page size; default page length; list shows too few rows

The Logins and Raw Logins lists show 20 rows a page (computerListDefaultPerPage) when page or perPage is missing. A value that is not a number, zero or negative also gets the default. The default page is 1.

## Login lists never show more than 100 rows a page

Lines 36-86

Asked as: max rows per page; perPage limit; too many rows requested; page size capped

When perPage is above 100 (computerListMaxPerPage), Turbo lowers it to 100. This applies to the Logins list and the Raw Logins list.

## Logins list: newest first, by computer or search

Lines 41-54

Asked as: sign-ins by computer; latest logins first; filter logins by machine; search sign-in records

The Logins list shows sign-ins newest first. The ?computer= value limits it to one computer, and Turbo trims spaces around that value. The ?q= value searches the list. Each row has the computer, its serial, the given name and surname, and the matched user. The order and the search fields come from the query file, which this file does not show.

## Raw Logins list leaves out the report

Lines 56-64

Asked as: rejected reports list; raw login report missing; why no report in list; failed login reports

The Raw Logins list shows rejected reports newest first. Each row has the computer name and the error, and ?q= searches the list. The list does not include the report itself. A request for one raw login gives the report.

## Search text matches anywhere and treats symbols literally

Lines 90-99

Asked as: search with percent sign; underscore in search; wildcard in search; search special characters

The search text (?q=) matches anywhere in the searched fields. Turbo trims spaces around the text first. The characters %, _ and \ count as plain characters, not wildcards. Blank search text applies no search pattern.

## A page past the end shows the last page

Lines 101-116

Asked as: page number too high; empty page showing; last page instead; past the final page

When the requested page is beyond the last page, Turbo returns the last page. The response reports that page number. Turbo also returns the total items and total pages. When there are no rows, the page stays as requested and the items are empty.

## Logins list can be limited to one computer

Lines 118-133

Asked as: only one machine's logins; logins for a single computer; filter logins by device; show logins of one pc

The Logins list takes an optional computer value (?computer=). Turbo trims leading and trailing spaces from it before it filters. The list is paged and can also be searched with ?q=.

## Logins list failure message

Lines 128-130

Asked as: logins won't load; logins error message; failed to load logins; sign-in list broken

When Turbo cannot load the Logins list, it shows the error "failed to load logins".

## Raw Logins list is searched but not filtered by computer

Lines 135-149

Asked as: filter raw logins by computer; computer filter ignored; raw logins by machine; raw list not filtering

The Raw Logins list takes only the search text (?q=) and the paging values. It has no computer filter, so a computer value in the request does nothing. The list leaves out the reports.

## Raw Logins list failure message

Lines 144-146

Asked as: raw logins won't load; raw logins error; failed to load raw logins; rejected logins list broken

When Turbo cannot load the Raw Logins list, it shows the error "failed to load raw logins".

## Opening one raw login gives the full report

Lines 151-176

Asked as: view rejected report; see raw login details; raw login not found; empty report shows null

A request for one raw login returns its created time, computer name, error and the report as received. When the raw login does not exist, Turbo answers "Raw login not found." When the stored report is empty, Turbo returns null as the report. The page shows the report as text because it is untrusted.

## Delete a raw login

Lines 178-198

Asked as: remove rejected login; delete failed sign-in report; get rid of raw login; raw login already gone

Deleting a raw login removes that rejected report and nothing else, because nothing refers to a raw login. When the raw login does not exist, Turbo answers "Raw login not found."
