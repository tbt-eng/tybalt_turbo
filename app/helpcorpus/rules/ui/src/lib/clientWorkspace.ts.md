<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspace.ts","sha256":"a10a132111d1b8fff96fa4b9422a59a309f0155fb4c49247ef8e490641e2a586","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspace.ts

## Client time divisions count committed time on its jobs

Lines 50-51

Asked as: client time by division; hours per division for client; committed time division; client division hours

A client's time division entry is a division with committed time on jobs whose client is this client. It shows the division's code, name and hours.

## Rate sheet counts on a client's rate sheets list

Lines 53-66

Asked as: rate sheet current jobs; past jobs on rate sheet; client rate sheet history; rate sheet revision count

A client's rate sheet list shows each rate sheet revision that the client's jobs use now or used in the past. The count of current jobs (current_jobs) is the jobs on the sheet now. The count of past jobs (past_jobs) is the jobs that only their rate sheet history put on it.

## Rate sheet usage label on a client's rate sheets

Lines 68-76

Asked as: rate sheet jobs label; jobs now and past count; how many jobs on rate sheet; past jobs label meaning

The usage label counts the client's jobs on a rate sheet. It reads "3 jobs now · 1 past" when the sheet has current and past jobs. It reads "5 jobs" with no past jobs, and "2 past jobs" with no current jobs. "now" appears only when a past count contrasts with it.

## Client list page number defaults to 1

Lines 78-81

Asked as: wrong page number contacts; page number reset to 1; client list page invalid; page 1 shown instead

The page number in a client contacts or invoicing profiles list must be a positive whole number (safe integer). Anything else, such as zero, a negative, a decimal or text, shows page 1.

## Client list page size allows only 10, 20 or 50

Lines 82-85

Asked as: rows per page client list; change how many rows; page size options contacts; list shows 20 rows

Contacts and invoicing profiles lists on a client page show 10, 20 or 50 rows per page. Any other value falls back to 20.

## Finding a saved record in a client list

Lines 88-110

Asked as: saved record not in list; can't find contact after saving; search hides saved record; record vanished after save

After a record is saved, Turbo locates it in the list using the same search filter and page size as the list (10, 20 or 50, otherwise 20). If the search hides the record, Turbo clears the search and goes to the page where the record sits in the unfiltered list. If the record is not in the list, Turbo does nothing.

## Client workspace view defaults to projects

Lines 112-124

Asked as: client page opens projects; wrong tab on client page; client view default; link opens wrong client view

The client page view is "proposals" or "owner" only when the link asks for it. Any other or missing view shows "projects". Search text, page and page size for contacts and invoicing profiles come from the link, and the recipient is empty by default.

## Client workspace links keep filters and drop empty ones

Lines 125-138

Asked as: client link keeps filters; clear filter from link; workspace link empty filter; link drops search settings

A client workspace link (workspaceHref) starts from the current URL's search settings. It removes a setting whose new value is empty or null and sets every other changed setting. The link opens the client's details tab unless another tab is given.

## Return link only goes to the same client's details page

Lines 139-156

Asked as: redirect after action client; return to client details; back link goes wrong page; client page return address

After an action, Turbo returns only to that client's details page (/clients/{id}/details), with its query or anchor. Any other, empty or off-site address sends the person to the client's details page.

## Missing contact shows "Contact unavailable"

Lines 157-163

Asked as: contact name blank; no contact shown; contact unavailable message; missing contact on profile

A contact's name is the given name and surname joined by a space, skipping empty parts. When the contact is missing, Turbo shows "Contact unavailable".

## Client label shows alias in parentheses

Lines 164-169

Asked as: client name alias brackets; show alias with client; client alias in name; official name and alias

A client appears as its official name, followed by its alias in parentheses when it has a non-blank alias.

## Invoicing profile bills the alias only when chosen

Lines 170-178

Asked as: invoice billed under alias; billing name alias wrong; which name goes on invoice; bill to alias instead

An invoicing profile bills the client's alias when its billing name (billing_name) is "alias" and the client has a non-blank alias. Otherwise it bills the official name.

## Search text for contacts and invoicing profiles

Lines 183-196

Asked as: search contacts by email; find contact by phone; search invoicing profile instructions; address not found in search

Searching contacts matches given name, surname, email and phone. Searching invoicing profiles matches the profile name, its invoicing instructions and its contact's name, email and phone. Address is not searched.

## Address shown as comma-separated parts

Lines 197-207

Asked as: client address format; address shows commas; address missing postal code; client address display

A client address shows street, city, province or state, postal code and country, separated by commas. Empty parts are skipped. Phone is not included.

## Who can create invoicing profiles

Lines 220-224

Asked as: who can add invoicing profile; can't create billing profile; need job claim for profile; permission to make profile

A person with the time claim or the job claim can create an invoicing profile from an existing contact. Only holders of the job claim (procurement) maintain clients and contacts.

## Empty invoicing profile list message

Lines 226-236

Asked as: no invoicing profiles message; empty profile list text; why no profiles shown; ask maintainer to add contact

With no profiles and no contacts, a job-claim holder sees "Add a client contact, then create a profile." Others see "Ask a client maintainer to add a contact first." With contacts, people with the time or job claim see "Create one using an existing client contact." Others see "Ask a client maintainer to create one."

## Profile creator name is blank without a creator

Lines 238-242

Asked as: creator name empty; profile created by blank; who made this profile; no creator on profile

An invoicing profile shows its creator's name only when it has a creator. Otherwise the name is empty, and a creator name that is only spaces is also empty.
