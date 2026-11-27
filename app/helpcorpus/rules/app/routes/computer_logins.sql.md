<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_logins.sql","sha256":"58e46fc8fd6cfbcfcfdb6f968e0e02aa30e96a57e0d899dcc2f03a5cbefa62c1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_logins.sql

## Matched user name falls back to user name, then email

Lines 18-27

Asked as: user name blank on login; login shows email instead of name; can't see other people's names; computers claim names hidden

The user name on a login row is the matched user's profile name (given name and surname). If that is blank, Turbo uses the user's name, then the user's email, then an empty string. A person who holds the computers claim can see these names without being able to read other people's profiles.

## Logins list shows newest logins first

Lines 24-42

Asked as: order of logins list; latest login at top; oldest logins first; which login shows first

The computer logins query lists every login, newest first (created, then id, both descending). The Logins page and the computer details page use it.

## Logins with no matched user still appear

Lines 24-27

Asked as: login missing from list; login with no user; unknown user login not shown; logins disappearing from page

A login whose user is blank or not found still appears in the list, with an empty user name. The row is dropped only if its computer record is missing, because the computer join is not optional.

## Computer filter limits logins to one computer

Lines 29-29

Asked as: filter logins by computer; only one computer's logins; logins for a single machine; show all computers logins

When a computer is given, the logins list shows only that computer's logins. When the computer is blank, it shows logins for every computer.

## Search matches four fields of a login row

Lines 30-41

Asked as: search logins by serial; find login by computer name; search box not finding login; search by person's name on login

A search pattern (LIKE, with backslash as the escape) matches the computer name, the serial, the reported name and the matched user's name as the row shows it. The reported name is the login's given name and surname. A blank search matches every row.
