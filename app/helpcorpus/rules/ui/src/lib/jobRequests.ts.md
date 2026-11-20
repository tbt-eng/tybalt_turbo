<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"975f550fb33fe526060340227d774792f888f52b09485d3afc1022244fbe49fe","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Checklist answers start as false

Lines 14-27

Asked as: checklist questions unanswered; request checklist won't continue; checkbox already ticked; can't continue new request

A new request checklist starts with every question set to false (unanswered). Each question must be answered true to continue.

## Statuses of a proposal that can be linked to a request

Lines 20-24

Asked as: proposal won't link to request; link proposal to project request; proposal status not allowed; request save refused proposal

A project request can link a proposal only when the proposal's status is In Progress, Submitted or Awarded. The server refuses other statuses when the request is saved.

## Contact email shows a note when client has none

Lines 28-31

Asked as: contact has no email; no email address showing; email blank on contact; contact marked no email

A contact shows its email address. If the contact has no email but is marked no email (no_email), Turbo shows "No email address". Otherwise it shows nothing.

## Name display falls back to given name and surname

Lines 32-34

Asked as: name missing on record; shows first and last name; blank name on record; only given name and surname

A record shows its name. If it has no name, Turbo joins its given name and surname with a space.

## Billing address falls back to the client address

Lines 35-48

Asked as: billing address wrong; invoice address not showing; address taken from client; empty billing address fields

Turbo shows the invoicing profile's address when it has any address part. Otherwise it shows the client's address. Address parts (address, city, province/state, postal code, country) are joined with commas, and empty parts are left out.

## Project value shows as whole Canadian dollars

Lines 49-57

Asked as: project value shows not provided; value displays with cents; dollar amount looks wrong; zero value showing

A project value displays as whole Canadian dollars with no cents, only when it is a positive whole number. Any other value, including zero, a decimal or a negative number, displays "Not provided".

## Print readiness: client name and description required

Lines 62-63

Asked as: request won't print; missing client name for print; blank description printing; print says incomplete

A project request is missing information for printing when the client name or the description is blank or only spaces. Turbo lists "Client name" or "Description".

## Print readiness: location must be a valid Plus Code

Lines 64-68

Asked as: location invalid for printing; plus code rejected; open location code wrong; map code not accepted

The request location must be a valid Plus Code: 8 characters, a plus sign, then 2 or 3 characters, from the Plus Code alphabet (23456789CFGHJMPQRVWX), in either case. Otherwise Turbo lists "Valid Plus Code location".

## Print readiness: project and billing contacts

Lines 69-76

Asked as: contact missing for print; billing contact not complete; contact email invalid printing; project contact incomplete

The project contact and the billing contact each need a name. Each also needs a valid-looking email unless it is marked no email (no_email). Turbo lists the missing items such as "Project contact name" or "Billing contact email".

## Print readiness: positive whole-dollar project value

Lines 77-78

Asked as: project value blocks print; value must be whole dollars; print value error; zero value cannot print

The project value must be a whole number above zero. Otherwise Turbo lists "Positive whole-dollar project value".

## Print readiness: fee basis must be set explicitly

Lines 79-80

Asked as: fee basis not set; time and materials unset; can't print fee type; choose billing type

The fee basis (time and materials or not) must be explicitly set and must be true or false. Otherwise Turbo lists "Fee basis".

## Print readiness: rate sheet must be ready

Lines 81-88

Asked as: rate sheet problem printing; missing role rate; rate sheet inactive; proposed rates incomplete

The request needs a rate sheet. For proposed rates, the server must find a rate for every role (no missing roles). For a chosen existing rate sheet, the sheet must be active. Otherwise Turbo lists "Rate sheet".

## Print readiness: billing rates for time and materials

Lines 89-98

Asked as: billing rates missing; overtime rate blank; no hourly rates for client; time and materials rates empty

When the fee basis is time and materials, the rate sheet must have entries, and every entry must have a numeric rate and overtime rate. Otherwise Turbo lists "Client billing rates".

## Job creation queue visibility

Lines 101-107

Asked as: no job creation queue; can't see creation queue; job requests menu hidden; queue count missing

The job creation queue shows only when requests are enabled, the person has the job claim, and the server sent a count for the creation queue (/job-requests?queue=creation). The server adds this count only for an active job holder while requests are enabled.
