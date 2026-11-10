<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/client_billing.go","sha256":"4d71cbb6a71214120e732545ba9582c02ad79d355fad174342d4602ef315a811","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/client_billing.go

## Which records get their selections checked against the client

Lines 20-32

Turbo checks selected contacts and invoicing profiles for jobs and project requests (job_requests). It checks the selected contact for invoicing profiles (client_invoicing_information). Other collections are not checked.

## Which records are checked when a contact or profile changes client

Lines 34-45

When Turbo saves an existing client contact (client_contacts), it checks the jobs, project requests (job_requests) and invoicing profiles that select it. When it saves an existing invoicing profile (client_invoicing_information), it checks the jobs and project requests that select it. Each of those records must belong to the same client as the contact or profile.

## Contact and invoicing profile must belong to the same client

Lines 63-75

When a job, project request, or invoicing profile is saved, Turbo checks the contact (contact) and invoicing profile (invoicing_information) it selects. Each selected record must belong to the same client as the record being saved. Otherwise Turbo refuses with "The selected record must belong to this client." Empty selections are skipped here and are left to required-field validation.

## Empty contact or invoicing profile is not checked for client

Lines 63-67

When a job, project request or invoicing profile leaves a contact or invoicing profile blank, Turbo skips the client check for that field. Turbo leaves a blank required field to required-field validation.

## Selected contact or invoicing profile not found

Lines 68-71

When a job, project request, or invoicing profile selects a contact or invoicing profile that Turbo cannot find, Turbo refuses the save. The message is "The selected record was not found."

## Cannot change client of a contact or profile used by another client

Lines 77-88

When an existing client contact or invoicing profile is saved, Turbo looks for records that use it but belong to a different client. A contact is checked against jobs, project requests and invoicing profiles. An invoicing profile is checked against jobs and project requests. If any such record exists, Turbo refuses with "This record is used by <collection> for another client." This check is skipped for new records.

## Client label shows alias in parentheses

Lines 92-104

Client pickers and lists show a client's official name. If the client has an alias, it follows in parentheses, as in "Name (Alias)". Turbo trims spaces around the alias, so a blank alias counts as none.

## Invoicing profile billing name: alias or official name

Lines 106-116

An invoicing profile bills the client's alias when its billing name (billing_name) is "alias" and the client has an alias. In every other case it bills the client's official name. That includes a blank choice and an alias choice for a client with no alias.

## Alias billing requires the client to have an alias

Lines 120-133

When an invoicing profile is saved with billing name (billing_name) "alias", the client must have an alias. Otherwise Turbo refuses with "This client has no alias. Add one to the client or bill its name." Turbo skips this check if the client cannot be found, and leaves that to other validation.

## Count of invoicing profiles billing a client's alias

Lines 135-138

Turbo counts the invoicing profiles of a client whose billing name (billing_name) is "alias".
