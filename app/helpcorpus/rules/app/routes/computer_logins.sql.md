<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computer_logins.sql","sha256":"0c58de34b7cdedecd7c2f8c27ecbde3fbc71e3cedc4dcfef93a053d754b61663","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computer_logins.sql

## Missing login values show as empty text

Lines 14-18

Asked as: blank computer name on login; login shows empty serial; missing given name in login; null values on logins page

Turbo shows an empty string, not a blank value, for a missing computer name, serial, reported given name, reported surname or user on a login row.

## Login row shows matched user name with fallbacks

Lines 19-28

Asked as: user name shows wrong on login; login shows email instead of name; name missing on computer login; who logged in name display

The user name (user_name) on a login row is the profile name (given name and surname). If the profile name is blank, Turbo uses the user name, then the email, then an empty string. This lets a holder of the computers claim see names without reading other people's profiles.

## Logins without a computer record are not listed

Lines 26-26

Asked as: login missing from list; deleted computer login disappeared; login not showing up; old login gone from logins

The logins list uses an inner join on computers. A login whose computer record no longer exists does not appear.

## Logins list can be limited to one computer

Lines 30-30

Asked as: filter logins by computer; only one computer's logins; show logins for a single machine; logins for one device

When the computer filter is not blank, the logins list shows only the logins of that computer. When it is blank, it shows the logins of every computer.

## Logins search fields

Lines 31-42

Asked as: search logins by serial; find login by user name; logins search not working; search computer name in logins

The search on the logins list matches the computer name, the serial, the reported name (given name and surname from the login) and the matched user's name as the row shows it. The search is a LIKE pattern with backslash as the escape character. A blank search returns every row.

## Logins search matches the whole value unless the pattern has wildcards

Lines 32-41

Asked as: search needs wildcard; partial match logins search; search only finds exact text; why does search miss part of name

The search uses LIKE with the pattern as given. Turbo does not add wildcards itself, so a pattern without wildcard characters must match the whole field.

## Logins list is newest first

Lines 43-44

Asked as: logins sorted oldest first; latest login at top; logins order wrong; paging through login history

The Logins page and the computer details page list logins with the newest first (created, then id, descending). Each page is cut by the limit and offset that the request passes.
