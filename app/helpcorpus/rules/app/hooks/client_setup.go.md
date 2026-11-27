<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/client_setup.go","sha256":"2713280b100eadded70e7125e89fd987faa2c71a195c85034de61eca1585a5b6","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/client_setup.go

## Contact and invoicing profile deletion uses

Lines 15-38

Asked as: cannot delete contact; delete client contact blocked; contact still in use; delete invoicing profile error

A client contact cannot be deleted while a job, a project request or an invoicing profile uses it as contact. An invoicing profile cannot be deleted while a job or a project request uses it. Contacts are not checked when their client is deleted, because they are deleted with the client. Their own uses are checked as they are deleted.

## Invoicing profile creator is set by Turbo

Lines 43-51

Asked as: who created invoicing profile; creator set automatically; invoicing profile owner; creator field ignored

When someone creates an invoicing profile (client_invoicing_information) through a request, Turbo sets its creator (creator) to the signed-in staff user. The creator is empty if the request does not come from a staff user. A creator value in the form is ignored.

## Invoicing profile creator cannot be changed

Lines 53-64

Asked as: change creator invoicing profile; edit profile creator; creator cannot be changed; reassign invoicing profile

Turbo refuses an update of an invoicing profile that changes its creator (creator). The message is "the creator of an invoicing profile cannot be changed".

## Client contact and invoicing profile references are validated on save

Lines 66-77

Asked as: client contact save fails; reference check error; invalid client link; contact won't save

When a client contact or an invoicing profile is created or updated, Turbo checks its references with utilities.ValidateClientReferences and refuses the save if the check fails. This file does not show the rules of that check.

## Contact with an email is not marked as having no email

Lines 67-70

Asked as: no email box ticked; no email still checked; contact has email flag; email but marked none

When Turbo saves a client contact (client_contacts) that has a non-blank email, it sets no_email to false. Spaces alone count as blank, so they do not trigger this.

## Invoicing profile billing name is validated on save

Lines 81-101

Asked as: billing name rejected; invoicing name invalid; billing name error; profile name won't save

Turbo checks the billing name of an invoicing profile on create and update with utilities.ValidateBillingName. The check also runs on saves made outside requests, such as from the dashboard or code. This file does not show its rules.

## Alias check applies to client updates only

Lines 102-102

Asked as: alias removal error; client alias not saving; remove client alias; alias check on edit

The alias removal check runs when a client is updated, not when it is created. It is skipped when the client had no alias before.

## Records still in use cannot be deleted

Lines 105-142

Asked as: cannot delete client; record still in use; delete blocked by jobs; client has linked records

Turbo refuses to delete a client, a client contact or an invoicing profile while other records use it. The message is "this record is still used by <collection>". A client is used by jobs (client, job_owner), project requests (job_requests: client, job_owner), invoicing profiles and client notes.

## Contact email must be unique across clients

Lines 144-174

Asked as: email already used; duplicate contact email; same email other client; contact email taken

Turbo refuses to save a client contact whose email another contact already has. The message is "This email is already a contact for <client name>." It says "another client" if the client cannot be found. The check compares exact values and skips an empty email.

## Cannot clear a client alias while invoicing profiles bill it

Lines 176-198

Asked as: cannot clear alias; alias billed by invoicing profile; remove alias error; delete client alias blocked

Turbo refuses to clear the alias of an existing client when invoicing profiles bill that alias. For one profile the message is "An invoicing profile bills this alias. Change it to bill the client name first." For several it is "N invoicing profiles bill this alias. Change them to bill the client name first." Setting the alias for the first time, or changing it to another value, is not blocked.
