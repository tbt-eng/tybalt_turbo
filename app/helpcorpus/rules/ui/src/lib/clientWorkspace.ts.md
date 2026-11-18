<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspace.ts","sha256":"1badd2439ce92398860e0d2f3ec862ccc7a1e5f3fbb703775e56d2c534592188","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspace.ts

## Rate sheet usage label on a client's rate sheets

Lines 68-75

Asked as: rate sheet job count; how many jobs use rate sheet; jobs now and past label; rate sheet usage text

Turbo describes a rate sheet's use by the client's jobs as "N jobs now · M past". It leaves out a part that is zero, such as "1 job now" or "2 past jobs". It shows nothing when both counts are zero (current_jobs, past_jobs).

## Client list page number and page size defaults

Lines 77-84

Asked as: wrong page number in list; page size not working; list shows wrong number of rows; client list paging defaults

A page number that is not a positive whole number becomes page 1. A page size is only 10, 20 or 50, and any other value becomes 20. This applies to the contacts and invoicing profiles lists.

## Locating a saved record in a client list

Lines 86-109

Asked as: saved client not showing in list; where did my record go after save; search cleared after saving; contact not visible after edit

After a save, Turbo finds the record in the list by the same search filter and page size as the list. The search is case-insensitive. When the saved record does not match the current search, Turbo clears the search and moves to the page where the record sits in the full list. Otherwise it goes to the record's page. A page size other than 10, 20 or 50 counts as 20.

## Client workspace URL options and defaults

Lines 111-123

Asked as: workspace tab wrong; client page opens wrong tab; workspace link options in url; proposals or owner tab

The client workspace view is "proposals" or "owner", and any other value shows "projects". The search text for contacts and invoicing profiles defaults to empty. The recipient defaults to empty. The page and page size options follow the list defaults.

## Workspace links drop empty options

Lines 124-137

Asked as: empty filter in link; link has blank parameters; client link extra query text; workspace url clutter

When Turbo builds a client workspace link, it removes an option whose new value is null or empty. It keeps the other options in the current URL. The link opens the details tab unless another tab is given.

## Return link falls back to client details

Lines 138-155

Asked as: back button goes wrong client; return link sends me elsewhere; back from client goes to other site; client back link broken

A return link is only used when it points to the same client's details page (/clients/{id}/details), with or without a query or a hash. Any other link, including a link to another site or another client, sends the person back to the client's details page.

## Contact name fallback

Lines 156-162

Asked as: contact name missing surname; blank contact name shown; no name on contact; contact shows unavailable

Turbo shows a contact's name as given name and surname, leaving out a part that is empty. It shows "Contact unavailable" when there is no contact.

## Client label shows alias in parentheses

Lines 163-168

Asked as: alias not showing with client; client name missing bracket alias; other name next to client; client display name alias

Turbo shows a client's official name, followed by its alias in parentheses when the client has an alias. An alias of only spaces does not count.

## Billing name: alias or official name

Lines 169-177

Asked as: wrong name on invoice; invoice billed under alias; billing shows trading name; which name goes on invoice

An invoicing profile bills the client's alias only when the profile's billing name (billing_name) is "alias" and the client has an alias. Otherwise it bills the client's official name.

## Count labels use singular for one

Lines 178-181

Asked as: says 1 jobs; plural wrong for one; count wording singular; 0 job text grammar

Turbo writes a count with the singular noun only when the count is exactly 1, such as "1 job". Any other count, including 0, uses the plural, such as "0 jobs".

## Fields that contact and profile searches cover

Lines 182-195

Asked as: search contacts by email; find profile by phone; profile search doesn't find instructions; what can I search for contacts

A contact search matches the given name, surname, email and phone. An invoicing profile search matches the profile name, its invoicing instructions, and the name, email and phone of its contact.

## Address text leaves out empty parts

Lines 196-206

Asked as: address shows extra commas; empty comma in address; address formatting gaps; postal code missing from address

Turbo writes an address as address, city, province or state, postal code and country, separated by commas. It leaves out any empty part. The phone is not part of it.

## Client list filters stay in the URL on details page

Lines 208-217

Asked as: filters lost after opening client; back loses my search; client list filter reset; keep filters when returning to list

On a client's details page, Turbo uses the list URL kept in the page state, so links and Back keep the current filters. It checks this URL with the safe return rule. It uses the real URL on other pages or when no list URL is saved.

## Who can create an invoicing profile

Lines 219-223

Asked as: who can add invoicing profile; cannot create billing profile; why no create profile button; permission to make invoicing profile

A person with the time claim or the job claim (procurement) can create an invoicing profile, using an existing client contact. Only holders of the job claim maintain clients and contacts.

## Message when a client has no invoicing profiles

Lines 225-235

Asked as: no profiles message for client; why client has no billing profile; empty profile list text; who should create profile message

The empty list says "This client has no invoicing profiles." When the client has no contacts, it adds "Add a client contact, then create a profile." for a job claim holder, and "Ask a client maintainer to add a contact first." for others. When the client has contacts, it adds "Create one using an existing client contact." for a time or job claim holder, and "Ask a client maintainer to create one." for others.

## Invoicing profile creator name display

Lines 237-241

Asked as: who created this profile; profile creator blank; creator name missing on profile; show profile author

Turbo shows the creator's name of an invoicing profile only when the profile has a creator. The name is trimmed. It shows an empty name when the creator is missing or the name is blank.
