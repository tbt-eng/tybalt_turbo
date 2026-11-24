<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspace.ts","sha256":"cfc7cf50e3a7b69c165cf259caa5fee7c28dc2b03b36c58dbc4b95727d3d3798","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspace.ts

## Contact usage: project contact, invoice recipient or both

Lines 36-46

Asked as: contact used as project contact; invoice recipient on jobs; which jobs use this contact; contact usage list

A client contact's usage lists the jobs and project requests that use it. Each row says whether the contact is the project contact (project_contact), the invoice recipient through an invoicing profile (invoice_recipient), or both.

## Hidden project requests are counted, not listed

Lines 47-52

Asked as: hidden project requests count; requests I cannot see; missing requests in usage; why some requests not shown

In the usage of a contact or an invoicing profile, project requests that the signed-in person may not read do not appear as rows. Turbo shows only their count (hidden_requests).

## Current and past jobs on a client rate sheet

Lines 70-83

Asked as: jobs on old rate sheet; past jobs rate sheet history; how many jobs use revision; rate sheet revision job count

For each rate sheet revision that the client's jobs use, current_jobs counts the jobs on it now. past_jobs counts the jobs that only their rate sheet history put on it.

## Rate sheet usage label on a client's rate sheets

Lines 85-93

Asked as: rate sheet job count label; what does now mean rate sheet; jobs now and past label; rate sheet usage text

Turbo shows how many of the client's jobs use a rate sheet. It shows "3 jobs now · 1 past" when the sheet has current and past jobs, and "5 jobs" when it has no past jobs. It shows "2 past jobs" when it has no current jobs. The word "now" appears only when a past count contrasts with it.

## Client list page number defaults to 1

Lines 95-98

Asked as: client list wrong page; page number in URL; list opens on page 1; invalid page number client

The page number in a client workspace list is the value in the URL when it is a whole number above 0. Zero, negative numbers, decimals and text give page 1.

## Client list page size allows only 10, 20 or 50

Lines 99-102

Asked as: rows per page client list; page size wrong in list; list shows 20 rows; how many rows per page

A client workspace list shows 10, 20 or 50 rows per page. Any other value in the URL gives 20 rows per page.

## Locating a saved record in a client list

Lines 104-127

Asked as: saved client record disappeared; can't find record after save; search hides saved item; list jumps after saving

After a save, Turbo finds the record in the list using the current search text (case-insensitive) and page size. If the search hides the record, Turbo clears the search and moves to the page that holds it in the full list. If the record is not in the list, Turbo does nothing. A page size outside 10, 20, 50 counts as 20.

## Client workspace view defaults to projects

Lines 129-132

Asked as: client opens on wrong tab; workspace default tab; client page opens projects; proposals or owner view link

The client workspace opens the view in the URL when it is "proposals" or "owner". Any other value or no value opens "projects".

## Search and page settings for contacts and invoicing profiles

Lines 133-139

Asked as: contacts list search remembered; invoicing profiles list page size; contacts page settings reset; per page contacts list

The contacts list and the invoicing profiles list each keep their own search text, page and page size in the URL. Search text is empty by default. Page size is 10, 20 or 50, with 20 as the default.

## Workspace links keep filters and drop empty settings

Lines 142-155

Asked as: filters lost in client link; link drops search; clear filter from client url; client link tab opens details

A client workspace link keeps the current search and page settings. A setting that a change sets to null or an empty string is removed from the link. Other changes replace the value. The link opens the details tab unless another tab is given.

## Return link after editing a client stays on details page

Lines 156-173

Asked as: after editing client goes where; return to client details after edit; edit client redirect wrong page; back link after saving client

After an edit, Turbo returns to the client's details page, with its search and hash, only when the return address points at that page. Any other address, including one for another client or another site, sends the person to /clients/{id}/details.

## Contact name fallback

Lines 174-180

Asked as: contact shows blank name; contact unavailable message; missing contact name on client; contact name not showing

Turbo shows a contact's given name and surname. It shows "Contact unavailable" when the contact is missing.

## Client label shows the alias in parentheses

Lines 181-186

Asked as: client name alias brackets; alias in brackets after client name; why client name has alias; client display name alias

Turbo shows the client's official name, followed by its alias in parentheses when it has a non-blank alias.

## Invoicing profile bills the alias only when chosen

Lines 187-195

Asked as: invoice shows alias instead of name; billed under alias not official name; billing name alias setting; invoice wrong client name

An invoicing profile bills the client's alias only when its billing name (billing_name) is "alias" and the client has a non-blank alias. Otherwise it bills the official name.

## Search text for contacts and invoicing profiles

Lines 200-213

Asked as: search contacts by email or phone; find invoicing profile by instructions; contact search not finding; search profiles by contact

A contact search matches given name, surname, email and phone. An invoicing profile search matches the profile name, its invoicing instructions and the search text of its contact.

## Address display skips empty parts

Lines 214-224

Asked as: address shows extra commas; blank parts in address; address formatting empty fields; postal code spacing in address

Turbo joins the address, city, province or state, postal code and country with commas. It leaves out empty parts.

## Who can create an invoicing profile

Lines 237-241

Asked as: cannot create invoicing profile; no button to make profile; who can add invoicing profile; permission for invoicing profile

Staff with the time claim or the job claim can create an invoicing profile from an existing contact. Only staff with the job claim maintain clients and contacts.

## Empty invoicing profile message depends on claims and contacts

Lines 243-253

Asked as: no invoicing profiles message; why profile message says ask maintainer; empty profile list text; create profile button missing message

With no contacts, a person with the job claim sees "This client has no invoicing profiles. Add a client contact, then create a profile." Others see "...Ask a client maintainer to add a contact first." With contacts, a person with the time or job claim sees "...Create one using an existing client contact." Others see "...Ask a client maintainer to create one."

## Invoicing profile creator name is blank without a creator

Lines 255-259

Asked as: creator name blank on profile; who created invoicing profile; profile shows no creator; empty creator name

Turbo shows the creator's name on an invoicing profile only when the profile has a creator. The name is trimmed. It is blank when there is no creator or no name.
