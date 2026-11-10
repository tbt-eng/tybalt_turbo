<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientWorkspace.ts","sha256":"da1d696e9b8049eda2f0a84da311b9e2f2aa5fdf1579b50b0166ee285781a83a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientWorkspace.ts

## Client list page number must be a positive whole number

Lines 53-56

Asked as: page number invalid; client list wrong page; page 0 or negative; list page not loading

In the client workspace lists, a page number (positivePage) counts only when it is a safe positive integer. Missing, zero, negative, fractional or non-numeric values fall back to page 1.

## Client list page size allows only 10, 20 or 50

Lines 57-60

Asked as: change items per page; client list page size; how many rows per page; page size options

A list page size (listPageSize) in the client workspace can be 10, 20 or 50. Any other value, or none, falls back to 20 per page.

## Saved record search and page position in client lists

Lines 62-85

Asked as: saved record not showing in list; search hides my saved contact; record missing after save; list jumps after saving

After a record is saved, Turbo finds it in a client's contact or invoicing profile list using the same search filter and page size as the list. The search is case-insensitive and matches any part of the text. If the filter hides the saved record, Turbo clears the search and goes to the record's page in the unfiltered list. If the record is not in the list, Turbo does nothing. Page sizes other than 10, 20 or 50 count as 20 (savedListPosition).

## Saved record page when search shows it

Lines 78-84

Asked as: which page is my saved record on; jump to saved record; record after save wrong page; search still shows saved entry

When the search still shows the saved record, Turbo moves to the page that holds it. The current page is limited to the range from 1 to the last page. The record counts as visible only when the current page is that page (savedListPosition).

## Client page view defaults to projects

Lines 87-90

Asked as: client page opens on projects; view setting in url; client proposals or owner view; wrong tab on client page

The client page view comes from the "view" setting in the address. Only "proposals" and "owner" are accepted. Any other value or a missing value gives "projects" (parseClientWorkspaceQuery).

## Client contact and invoicing list defaults

Lines 91-97

Asked as: contacts list default search; invoicing list page reset; bad page in address; client list url values

The contacts and invoicing profile lists each take a search text, a page and a page size from the address. The search text defaults to empty. The page defaults to 1 and the page size to 20 when the value is not valid. The recipient defaults to empty (parseClientWorkspaceQuery).

## Client page links drop empty settings

Lines 100-113

Asked as: empty filter in link; link keeps blank settings; client page link extra params; cleaned up client url

When Turbo builds a client page link, it removes any setting whose new value is empty or null. It keeps the other settings. The link opens the "details" tab unless another tab is given (workspaceHref).

## Return link must stay on the client's details page

Lines 114-131

Asked as: back link goes to wrong page; return link redirects; after save sent to other client; broken return url

Turbo follows a return link only if it points to this client's details page, with an optional query or fragment. Any other link, such as another client, another page or another site, sends the person to /clients/{id}/details. A link that cannot be read does the same (safeReturnTo).

## Contact name fallback

Lines 132-138

Asked as: contact name blank; contact shows unavailable; missing surname in contact; name shows wrong format

A contact's name is the given name and the surname, with a space between them, and an empty part is left out. When the contact is missing, Turbo shows "Contact unavailable" (contactName).

## Client label shows alias in parentheses

Lines 139-144

Asked as: alias in brackets; client name shows alias; trade name next to client; client label wrong

A client is labeled with its official name. When the client has an alias, Turbo adds it in parentheses. An alias of only spaces counts as no alias (clientLabel).

## Invoicing profile bills alias only when chosen

Lines 145-153

Asked as: invoice billed under alias; billing name wrong; invoice shows alias not legal name; which name is billed

An invoicing profile bills the client's alias only when the profile's billing name (billing_name) is "alias" and the client has a non-blank alias. In every other case it bills the client's official name (billingName).

## Count labels use singular only for one

Lines 154-157

Asked as: 1 job or 1 jobs; plural wording wrong; count says job instead of jobs; zero jobs text

A count shows the singular noun only when it equals 1, such as "1 job". Zero and all other numbers show the plural, such as "0 jobs" and "2 jobs" (countLabel).

## Contact search matches name, email and phone

Lines 158-162

Asked as: search contacts by phone; find contact by email; contact search not finding address; search contact by name

The contact search looks at the given name, surname, email and phone. It does not look at the address (contactSearchText).

## Invoicing profile search includes contact details

Lines 163-171

Asked as: search invoicing profiles by contact; find profile by instructions; profile search missing results; search billing profile by email

The invoicing profile search looks at the profile name, the invoicing instructions, and the name, email and phone of the profile's contact (invoicingSearchText).

## Address shown as one comma-separated line

Lines 172-182

Asked as: address display format; address with commas; address shows phone; address missing country

An address is shown as street address, city, province or state, postal code and country, with commas between them. Empty parts are left out. The phone number is not included (addressText).

## Client list filters kept when URL does not change

Lines 184-193

Asked as: lose filters going back; back button resets list; client list filters cleared; keep filters on client page

On a client's details page, Turbo uses the list address saved in the page state, when there is one, so links and Back keep the current filters. It first checks that address with the return link rule. On other pages it uses the normal address (currentWorkspaceUrl).

## Who can create invoicing profiles

Lines 195-199

Asked as: who can add invoicing profile; cannot create billing profile; no permission to make profile; time claim create profile

Staff with the time claim or the job claim (procurement) can create invoicing profiles from existing client contacts. Only holders of the job claim maintain clients and contacts (canCreateInvoicingProfile).

## Message when a client has no invoicing profiles

Lines 201-211

Asked as: no invoicing profiles message; why no create button; ask maintainer to add contact; client has no billing profile

The message begins "This client has no invoicing profiles." When the client has no contacts, a holder of the job claim sees "Add a client contact, then create a profile." Anyone else sees "Ask a client maintainer to add a contact first." When the client has contacts, staff with the time claim or the job claim see "Create one using an existing client contact." Anyone else sees "Ask a client maintainer to create one."

## Profile creator name hidden when unknown

Lines 213-217

Asked as: who created this profile; creator name missing; profile author blank; profile created by unknown

An invoicing profile shows its creator's name only when the profile has a creator. The name has spaces trimmed from both ends. Otherwise the name is empty (profileCreatorName).
