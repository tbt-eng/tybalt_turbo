<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"fb409b9143dc0120615cf402e6538c76c7b5a641ebb08975c002c3866502906e","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Request checklist starts with every question unanswered

Lines 14-22

A new checklist marks every question false. Each question must be answered true to continue at the request's stage, either approval or job creation once approved.

## Contact with no email shows a 'No email address' note

Lines 23-26

A contact's email address is shown when the contact has one. If the contact has no email but is marked as having none (no_email), Turbo shows "No email address". Otherwise it shows nothing.

## Contact or record name falls back to given name and surname

Lines 27-29

Turbo shows a record's name when it has one. If not, it joins the given name and surname (given_name, surname) with a space, leaving out any that is empty.

## Billing address falls back to the client's address

Lines 30-43

Turbo shows the invoicing profile's address when it has one. If that address is empty, it shows the client's address. An address joins address, city, province/state, postal code and country with commas, leaving out empty parts.

## Project value shows as whole Canadian dollars or 'Not provided'

Lines 44-52

Turbo shows a project value as Canadian dollars with no cents only when it is a whole number greater than 0. Zero, negative, fractional or non-number values show "Not provided".

## Print check: client name and description are required

Lines 57-58

The project authorization print check lists "Client name" as missing when the client's name is empty or only spaces. It lists "Description" as missing when the request's description is empty or only spaces.

## Print check: location must be a valid Plus Code

Lines 59-63

The print check lists "Valid Plus Code location" as missing when the location is empty or does not match a Plus Code. The code is 8 characters, then "+", then 2 or 3 characters, from the Plus Code alphabet 23456789CFGHJMPQRVWX. Case is ignored and surrounding spaces are trimmed.

## Print check: project and billing contacts need name and email

Lines 64-71

The print check requires a name for both the project contact and the billing contact. It lists "Project contact name" or "Billing contact name" when the name is empty. It lists the matching "email" item when the email does not look like a valid address, unless the contact is marked as having no email (no_email).

## Print check: project value must be a positive whole number

Lines 72-73

The print check lists "Positive whole-dollar project value" as missing when the project value is not a whole number or is 0 or less.

## Print check: fee basis must be set explicitly

Lines 74-75

The print check lists "Fee basis" as missing unless the time-and-materials choice has been set (time_and_materials_set is true) and is a true/false value.

## Print check: an active rate sheet is required

Lines 76-76

The print check lists "Rate sheet" as missing when the request has no rate sheet or the rate sheet is not active.

## Print check: time-and-materials needs complete client billing rates

Lines 77-86

When the fee basis is time and materials, the print check lists "Client billing rates" as missing if there are no rate sheet entries or any entry lacks a numeric rate or overtime rate. Other fee bases skip this check.

## Job creation queue is visible only to job holders when requests are enabled

Lines 89-96

Turbo shows the job creation queue only when job requests are enabled, the person holds the job claim, and the server has sent a count for the creation queue (/job-requests?queue=creation).
