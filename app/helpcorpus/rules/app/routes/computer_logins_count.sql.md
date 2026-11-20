<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_logins_count.sql","sha256":"d7a6a7e8194681ced74f4868f798d5fb6d52a9f0d4b7a9863ec7155d8c6b34f5","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_logins_count.sql

## Computer logins count uses the same filter as the list

Lines 3-7

Asked as: login count doesn't match list; logins missing from count; deleted computer logins counted; login with no user count

This query counts the computer login rows that computer_logins.sql pages through, with the same joins and the same filter. Each login must match an existing computer (inner join on computers), so a login whose computer record is gone is not counted. A login with no user is still counted (left joins on users and profiles).

## Computer logins filter by computer, empty means all

Lines 9-9

Asked as: count logins for one computer; filter logins by computer; all computers logins total; computer filter empty

When the computer filter is an empty string, Turbo counts logins from every computer. Otherwise it counts only the logins of that one computer (l.computer).

## Computer logins search fields and empty search

Lines 10-20

Asked as: search computer logins by serial; search login by computer name; login search returns nothing; search matches user name logins

An empty search counts every login that passes the computer filter. Otherwise a login is counted when the search pattern matches the computer name, the computer serial, the login's own name (given name and surname from the login report), or the user's name. The match uses LIKE with the backslash as the escape character. A missing value counts as an empty string.

## User name used by the computer logins search

Lines 14-20

Asked as: user name blank in login search; profile name missing search; search by user email logins; login search wrong user name

For the user's name, Turbo uses the first available of these: the profile's given name and surname, the user's name, the user's email, then an empty string. A blank profile name is skipped (NULLIF). The login's own reported name is searched separately and does not replace this one.
