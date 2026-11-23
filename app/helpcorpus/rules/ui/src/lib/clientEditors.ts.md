<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditors.ts","sha256":"f8bbee17cc95b1762a0beea04033d39cb01158f816b2e90f2e08d18c195a3206","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditors.ts

## Fields that the editors offer

Lines 4-29

Asked as: fields in contact editor; what fields can I edit; postal address fields; invoicing profile fields

The contact editor has given name, surname, email, phone and the postal fields. The client editor has name, alias, phone and the postal fields. The invoicing profile editor has name, contact, billing name, fax and invoicing instructions. The postal fields are street, city, province / state, postal code and country.

## Business development lead is not editable on client screens

Lines 19-29

Asked as: cannot change business development lead; busdev lead locked; saving client keeps old lead; change client's busdev owner

Turbo no longer lets anyone edit a client's business development lead on any screen. The client editor does not send the lead field (clientFields omits it), so saving a client leaves its existing lead unchanged. Clients keep the lead they already have.

## Stored business development lead hidden on client page

Lines 19-23

Asked as: business development lead missing; busdev lead disappeared from client; lead not showing on client page; hidden busdev lead

Turbo hides a client's stored business development lead on the client page, because showBusdevLead is false. Hiding it does not remove or change the lead.

## Editor forms show empty text for missing fields

Lines 32-34

Asked as: editor fields blank; empty fields when editing; missing values show nothing; form shows blank text

When an editor opens a record, Turbo fills each editable field from the record. Any missing or empty value becomes an empty string.

## Contact requires given name and surname

Lines 35-39

Asked as: contact won't save; missing first name on contact; required field contact form; surname needed to save

Turbo refuses to save a contact whose given name (given_name) or surname is empty or only spaces. It shows "This field is required." beside each missing field.

## Contact email must look like an address when entered

Lines 40-42

Asked as: email rejected on contact; invalid email contact; email address not accepted; contact email error

Turbo checks a contact's email only when one is entered. Without spaces it must have text, an @, text, a dot and more text, or Turbo shows "Enter a valid email address." An empty email is accepted. Surrounding spaces are trimmed before the check.

## Contact counts as in use by jobs, requests or profiles

Lines 45-49

Asked as: contact can't be deleted; contact linked to jobs; contact used by requests; contact has profiles

Turbo treats a contact as in use when its job count (job_count), project request count (request_count) or invoicing profile count (profile_count) is above 0. A contact with all three at 0 or missing is not in use.

## Error message combines API message and field messages

Lines 50-57

Asked as: error message unclear; save failed what went wrong; field errors after action; action failed no details

When an action fails, Turbo shows the API message followed by each field message as "field name: message", with underscores in field names replaced by spaces. If the error has no message, Turbo shows "The action failed."

## Return link after saving a contact or invoicing profile

Lines 58-66

Asked as: where do I land after saving contact; page jumps after saving invoicing; saved contact redirect; back to client after saving profile

After a contact or invoicing profile is saved, Turbo returns to the previous page with the edit parameter and any earlier saved_contact or saved_invoicing marker removed. It adds saved_contact or saved_invoicing with the record ID. It jumps to the contacts section for a contact and the invoicing section for an invoicing profile.
