<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/client_setup.go","sha256":"a035de1baa2c4ad87f50e5d7ada5ec2bc26953a12713d76452ad404c81e0a45a","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/client_setup.go

## What blocks deleting each record

Lines 15-39

Asked as: cannot delete client; delete blocked by jobs; why can't I delete a record; still used by jobs; delete contact error

A client is blocked by jobs (client, job_owner), job requests (client, job_owner), invoicing profiles, client notes and invoices. A contact is blocked by jobs, job requests and invoicing profiles. An invoicing profile is blocked by jobs and job requests. Contacts cascade with their client, and their own uses are checked as they are deleted.

## Invoicing profile creator is set by Turbo

Lines 44-52

Asked as: who created invoicing profile; creator field empty; set creator automatically; creator not saved; invoicing profile owner

When someone creates an invoicing profile (client_invoicing_information) through a request, Turbo sets its creator (creator) to the signed-in user. If the signed-in account is not a user, the creator is left blank. A value sent in the form is ignored.

## Invoicing profile creator cannot be changed

Lines 54-65

Asked as: change creator invoicing profile; creator cannot be changed; edit creator error; update creator refused; reassign invoicing profile

Turbo refuses an update of an invoicing profile that changes its creator (creator). The message is "the creator of an invoicing profile cannot be changed".

## Contact with an email is not marked as having none

Lines 68-71

Asked as: no email checkbox; contact marked no email; no_email still ticked; email not saved contact; contact has email but flag

When a client contact (client_contacts) is created or updated with a non-blank email, Turbo sets no_email to false. Whitespace-only emails count as blank.

## Client references are validated on contacts and invoicing profiles

Lines 72-78

Asked as: client reference invalid; wrong client on contact; contact client check fails; profile client not valid; client link error

Turbo checks the client references of a client contact or invoicing profile on create and update with utilities.ValidateClientReferences. The rules of that check are in another file.

## Invoicing profile billing name is validated

Lines 102-102

Asked as: billing name rejected; invalid billing name; billing name error; profile name won't save; bill name validation

Turbo checks the billing name of an invoicing profile on create and update with utilities.ValidateBillingName. This also applies to saves made outside requests. The rules of that check are in another file.

## Client alias cannot be cleared while profiles bill it

Lines 103-103

Asked as: remove client alias; alias won't clear; cannot delete alias; clear alias error; alias still billed

Turbo refuses to clear a client's alias while invoicing profiles bill that alias. The message is "An invoicing profile bills this alias. Change it to bill the client name first." When more than one profile bills it, the message is "N invoicing profiles bill this alias. Change them to bill the client name first." This check applies only to updates of an existing client that had an alias.

## Records still in use cannot be deleted

Lines 109-143

Asked as: cannot delete still used; record in use delete; why can't remove client; delete fails related records; still used by error

Turbo refuses to delete a client, client contact or invoicing profile while other records refer to it. The message is "this record is still used by <collection>". The check applies to requests, saves outside requests and cascaded deletes.

## Contact email must be unique across clients

Lines 145-175

Asked as: duplicate contact email; email already taken; same email two clients; contact email exists; email used by another client

Turbo refuses to create or update a contact whose email another contact already has. The message is "This email is already a contact for <client name>.", or "another client" if the client is not found. The comparison is exact, and a blank email is not checked.

## Alias removal check details

Lines 177-199

Asked as: alias removal rules; when can alias be removed; alias check skipped; new client alias check; profiles billing alias count

The alias check skips new clients, clients that still have an alias, and clients that had no alias before. Otherwise it counts the invoicing profiles that bill the alias.
