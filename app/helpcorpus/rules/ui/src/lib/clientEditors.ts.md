<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditors.ts","sha256":"0071aca8add0691a67a52019b7b60dba18ac5ef2ed76c4eee03d96a57c58f20b","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditors.ts

## Fields in the client form

Lines 4-24

Asked as: client form fields; which fields client editor has; client editor fields list; fields on client record

The client editor works with these fields: name, alias, phone, street, city, province / state, postal code and country.

## Fields in the contact form

Lines 11-17

Asked as: contact form fields; contact editor fields; what fields for a contact; contact details I can enter

The contact editor works with given name, surname, email, phone, street, city, province / state, postal code and country.

## Fields in the invoicing profile form

Lines 18-18

Asked as: invoicing profile fields; billing name and fax; invoicing form inputs; invoicing instructions field

The invoicing profile editor works with these fields: name, contact, billing name, fax and invoicing instructions.

## Business development lead is not editable on clients

Lines 19-23

Asked as: cannot change busdev lead; busdev lead disappeared; business development lead not saving; lead missing from client page

The client editor does not send the business development lead (busdev lead) field, so saving a client leaves its existing lead unchanged. No screen edits the lead. The flag showBusdevLead is false, so the client page hides the stored lead.

## Empty fields show as blank in editors

Lines 27-29

Asked as: empty field shows blank; field shows nothing when opened; missing value in editor; why is the field blank

When an editor loads a record, it shows an empty string for any field that is missing or has a falsy value.

## Contact requires given name and surname

Lines 32-34

Asked as: contact name won't save; surname required error; given name is blank; contact save refused

Turbo refuses to save a contact when the given name or the surname is empty or only spaces. It shows "This field is required." beside that field.

## Contact email is optional but must look valid

Lines 35-37

Asked as: contact email invalid; email rejected on contact; email address error; contact email format wrong

A contact's email is optional. If one is entered, Turbo checks that it has text, an @, and a domain with a dot, with no spaces, after trimming. If not, it shows "Enter a valid email address."

## Contact counts as in use by jobs, requests or profiles

Lines 40-44

Asked as: cannot delete contact; contact linked to jobs; contact still in use; why contact can't be removed

A contact is in use when its job count, project request count or invoicing profile count is greater than 0. A count of 0 or a non-numeric value does not count as in use.

## Error message display combines API message and field messages

Lines 45-52

Asked as: error message after save; action failed message; what does the error mean; field error text shown

When an action fails, Turbo shows the API message followed by each field message as "field name: message", with underscores in field names replaced by spaces. If there is no message, it shows "The action failed."

## Return after saving a contact or invoicing profile

Lines 53-62

Asked as: where do I land after saving contact; page after saving invoicing profile; saved marker on contact; jump to contacts section

After a save, Turbo returns to the page it came from. It removes the edit marker and earlier saved or replaced markers, and adds a marker for the saved contact or invoicing profile. It jumps to the contacts section for a contact and the invoicing section for an invoicing profile.

## Return after replacing an invoicing profile

Lines 63-70

Asked as: profile replaced message; old invoicing profile replaced; return after replacing invoicing; which invoicing profile is active now

After an invoicing profile is replaced, Turbo returns to the client page and marks the profile that took over, which shows "Profile replaced." It does not keep the normal saved marker.
