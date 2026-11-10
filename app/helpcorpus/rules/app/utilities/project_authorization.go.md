<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/project_authorization.go","sha256":"40725a529425e600fa4aa5837e695f07ac409b2316fbb7b69b3ae0d2da0d5577","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/project_authorization.go

## Abandonment period for job requests is 60 days

Lines 16-16

Turbo defines the job request abandonment period as 60 days (JOB_REQUEST_ABANDONMENT_DAYS). This file only sets the number. It does not show what happens after the period ends.

## Project value must be a positive whole dollar amount

Lines 19-30

When Turbo checks an approver's authority, it refuses a project value that is zero or negative, above 9007199254740991, not a number, infinite, or not a whole number. The message is "project value must be a positive whole dollar amount". This check runs before any check of the approver.

## Project approval errors show as bad request messages

Lines 23-27

Every refusal from the approver authority check reaches the user as a bad request error with the message of that refusal.

## Project approver must be active

Lines 31-37

Turbo refuses a project value approval when the chosen approver is not an active user. The message is "approver must be active".

## Approver needs the project authorization claim

Lines 38-44

Turbo refuses a project value approval when the chosen approver has no project_authorization claim. The message is "approver must have project_authorization authority".

## Approver needs a project value limit set

Lines 45-51

Turbo refuses a project value approval when the approver's project_authorization claim has no limit record (project_authorization_props). The message is "approver has no value limit".

## Approver's limit must cover the full project value

Lines 52-56

Turbo compares the approver's max_value with the full proposed value. A limit equal to the value is accepted. A limit below the value is refused with "approver limit does not cover the full project value". Turbo also refuses this when the stored limit is zero or negative, above 9007199254740991, not a number, infinite, or not a whole number. When the check passes, Turbo returns the approver's limit.

## Request contacts need an email or a no-email mark

Lines 62-79

Turbo checks the project contact of a project request, and the billing contact of its invoicing profile when that profile is found. Each contact must have a non-blank email address or be marked no_email. Otherwise Turbo refuses with "Add an email address for the project contact, or mark the contact as having no email address." The message says "billing contact" for the billing contact. A blank email made of spaces counts as missing.

## Billing contact check is skipped without an invoicing profile

Lines 65-67

When the invoicing profile of the request cannot be found, Turbo does not check a billing contact. It checks only the project contact.

## Missing request contact is refused

Lines 68-72

Turbo refuses a project request when a contact record cannot be found. The message is "The project contact was not found." or "The billing contact was not found."
