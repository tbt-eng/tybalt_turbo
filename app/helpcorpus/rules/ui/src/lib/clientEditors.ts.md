<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditors.ts","sha256":"de966e8dba28eb2877b34e435ad8273ea9dd71b938bdecf02a4c698949ead104","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditors.ts

## Client editor fields

Lines 4-29

Asked as: client editor fields; edit client address; client phone and alias; postal code country fields

The client editor works with the name, alias, phone and the postal fields: street, city, province/state, postal code and country.

## Business development lead is hidden in the client editor

Lines 19-29

Asked as: business development lead missing; busdev lead not showing; client lead hidden; saving client changes lead

The client editor does not show the business development lead (business_development_lead) while showBusdevLead is false, which is the current setting. The editor does not send the field, so saving a client leaves its existing lead unchanged. The field, the busdev claim and existing values stay.

## Editor fields default to empty text

Lines 32-34

Asked as: empty field shows blank; missing field in editor; editor blank values; null field client form

When the editor opens a record, any empty or missing field becomes an empty string.

## Contact first and last name are required

Lines 35-39

Asked as: contact name required; cannot save contact no name; contact first name error; contact last name missing

Turbo refuses to save a contact without a first name (given_name) or a last name (surname). A name of only spaces counts as empty. The message is "This field is required."

## Contact email format check

Lines 40-42

Asked as: invalid contact email; email rejected on contact; contact email error; email address not accepted

When a contact has an email, Turbo checks it for text, an @, and a domain with a dot, with no spaces, after trimming. If the email fails, Turbo shows "Enter a valid email address." An empty email passes, because it is optional.

## Contact counts as in use

Lines 45-49

Asked as: contact cannot be deleted; contact in use; contact has jobs; contact linked to requests

A contact is in use when its job count, request count or invoicing profile count is above 0. A count that is missing or not a number does not count as in use.

## Error display joins message and field messages

Lines 50-57

Asked as: error message details; save failed field errors; action failed message; which field has error

Turbo shows an action error as the API message followed by each field message, in the form "field name: message", with underscores in the field name shown as spaces. If there is no message, Turbo shows "The action failed."

## Return link after saving a contact or invoicing profile

Lines 58-66

Asked as: after saving contact returns; return to page after save; invoicing profile saved redirect; back link after editing contact

After a contact or invoicing profile is saved, Turbo returns to the earlier page without the edit parameter. It marks the saved record in the link and jumps to the contacts section for a contact or the invoicing section for an invoicing profile.
