<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditors.ts","sha256":"4dc04f0c4007aa1f81908db12c9c8f1a741682374df83cca206af6c0da439832","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditors.ts

## Fields a client editor saves

Lines 4-29

Asked as: blank client fields; client editor save; empty values saved; missing client address

The client editor works with these fields: name, alias, phone, and the postal fields (street, city, province/state, postal code, country). Blank or missing values become an empty string.

## Invoicing profile and contact editable fields

Lines 11-18

Asked as: invoicing profile fields; contact editor fields; billing name invoicing; invoicing instructions field

The contact editor works with given name, surname, email, phone and the postal fields. The invoicing profile editor works with name, contact, billing name, fax and invoicing instructions.

## Business development lead is not editable on clients

Lines 19-24

Asked as: business development lead; busdev lead missing; change client lead; client lead not saving

No screen edits a client's business development lead. The client editor does not send the lead field, so saving a client leaves its existing lead unchanged. The flag showBusdevLead is false, which hides the stored lead on the client page.

## Contact requires given name and surname

Lines 30-34

Asked as: contact name required; given name required; surname blank error; required field contact

Saving a contact fails on the given name or surname when it is empty or only spaces. Turbo shows 'This field is required.' beside that field.

## Contact email must look like an address

Lines 35-37

Asked as: invalid contact email; email address error; contact email rejected; email format wrong

When a contact has an email, it must have text, an @, and a domain with a dot, with no spaces. Otherwise Turbo shows 'Enter a valid email address.' An empty email is allowed.

## Contact counts as in use by jobs, requests or profiles

Lines 40-44

Asked as: cannot delete contact; contact still linked; contact used by jobs; contact has requests

Turbo treats a contact as in use when its job count, project request count or invoicing profile count is above 0.

## Error display joins message and field messages

Lines 45-52

Asked as: save error message; action failed message; error message details; field error text

A failed action shows the API message followed by each field's message, with underscores in field names replaced by spaces. If there is no message, Turbo shows 'The action failed.'

## Return link after saving a contact or invoicing profile

Lines 53-61

Asked as: after saving contact; returned to page; saved contact highlight; where did save go

After a save, Turbo returns to the previous page. It removes the edit and earlier saved_contact or saved_invoicing markers, and sets saved_contact or saved_invoicing to the saved record's ID. It jumps to the contacts section for a contact and the invoicing section for an invoicing profile.
