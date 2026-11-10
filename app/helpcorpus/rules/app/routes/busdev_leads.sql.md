<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/busdev_leads.sql","sha256":"453912cbcf451df1efa6fb48a98ac205427f632b8b05b19cfa03fcc31061342d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/busdev_leads.sql

## Business development lead list includes holders of the busdev claim

Lines 1-8

Turbo builds the list of business development leads from users who hold the 'busdev' claim. Each user appears once (DISTINCT), even if the user has several matching claim rows.

## Missing name or email shows as empty text in lead list

Lines 3-9

In the business development lead list, a missing given name (given_name), surname or email shows as an empty string, not as a blank value. A user with no profile still appears in the list, with empty name fields.

## Business development leads are sorted by surname, then given name

Lines 10-10

Turbo sorts the business development lead list by surname, then by given name. The sort uses the raw profile values, so users with no profile are placed by the database's default handling of missing values, not as empty strings.
