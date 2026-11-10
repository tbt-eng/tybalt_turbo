<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clientEditors.ts","sha256":"84f5e92ed0bc8c2b986d164b3869f0eff218afa1aade999f653bac3760d13990","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clientEditors.ts

## Fields in the client, contact and invoicing profile editors

Lines 4-25

A contact has the fields given name, surname, email, phone and the address fields. An invoicing profile has name, contact, billing name (billing_name), fax and invoicing instructions. A client has name, alias, business development lead, phone and the address fields. The address fields are street, city, province/state, postal code and country.

## Editor fields that are missing show as empty

Lines 28-30

When an editor form loads a record, each listed field that is missing or empty is shown as an empty string. Only the listed fields are loaded.

## Contact requires given name and surname

Lines 33-35

When a contact is saved, Turbo requires a given name (given_name) and a surname (surname). A value that is empty or only spaces is refused with the message "This field is required.".

## Contact email must look like an address

Lines 36-38

Turbo checks a contact's email only when one is entered. An email that has no text before and after an @ and no dot in the domain is refused with the message "Enter a valid email address.". An empty email is accepted, so the email is optional.

## A contact is in use if linked to jobs, requests or profiles

Lines 41-45

Turbo treats a contact as in use when its job count (job_count), project request count (request_count) or invoicing profile count (profile_count) is greater than 0. A count of exactly 0 does not make it in use.

## Error display joins the message with field messages

Lines 47-53

When an action fails, Turbo shows the error message followed by one line per field in the form "field name: message", with underscores in the field name replaced by spaces. If there is no message, Turbo shows "The action failed.".

## After saving a contact or profile, Turbo returns to the page

Lines 54-62

After a contact or invoicing profile is saved, Turbo returns to the earlier page. It removes the edit marker and any earlier saved marker, and records the saved record's id. It then jumps to the contacts section for a contact, or the invoicing section for an invoicing profile.
