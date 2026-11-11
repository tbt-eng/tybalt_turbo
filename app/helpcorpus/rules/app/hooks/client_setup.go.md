<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/client_setup.go","sha256":"eb06d4893a913223a02022861b10dca53a941b2ca48f4052040eb174f3fc1eb3","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/client_setup.go

## What blocks deleting a client

Lines 15-29

Asked as: cannot delete client; client won't delete; delete client blocked; client still linked to jobs

A client cannot be deleted while any of these records refer to it: jobs (client or job owner), project requests (job_requests, client or job owner), invoicing profiles, client notes, or invoices. Its contacts are not on this list; they are deleted along with the client, and each contact's own uses are checked as it is deleted.

## What blocks deleting a contact

Lines 30-34

Asked as: cannot delete contact; contact won't delete; remove client contact; contact still in use

A client contact cannot be deleted while a job, a project request (job_requests) or an invoicing profile uses it as its contact.

## What blocks deleting an invoicing profile

Lines 35-38

Asked as: cannot delete invoicing profile; billing profile won't delete; remove invoicing info; invoicing profile in use

An invoicing profile cannot be deleted while a job or a project request (job_requests) uses it as its invoicing information.

## Invoicing profile creator is set by Turbo

Lines 44-52

Asked as: who created invoicing profile; creator of billing profile; set profile creator; creator empty after request

When an invoicing profile (client_invoicing_information) is created through a request, Turbo sets its creator (creator) to the signed-in user's ID. The creator stays empty if the requester is not in the users collection. A value sent in the form is ignored.

## Invoicing profile creator cannot be changed

Lines 54-65

Asked as: change profile creator; edit invoicing creator; creator field locked; reassign billing profile owner

Turbo refuses an update to an invoicing profile that changes its creator (creator). The message is "the creator of an invoicing profile cannot be changed".

## Contact with an email is not marked as having no email

Lines 67-71

Asked as: contact has email but no email flag; no email checkbox still ticked; email ticked wrongly; unmark no email

When a client contact is created or updated with a non-blank email (email), Turbo sets "no_email" to false. Spaces alone count as blank.

## Contacts and invoicing profiles must have valid client references

Lines 72-78

Asked as: client reference error saving; contact save refused client; invalid client on contact; profile client link rejected

When a client contact or an invoicing profile is created or updated, Turbo checks its client references with utilities.ValidateClientReferences. Turbo refuses the save if that check fails. The rules of that check are in another file.

## Invoicing profile billing name is validated on create and update

Lines 82-102

Asked as: billing name rejected; invoicing name error; bad billing name; billing name not saving

Turbo checks the billing name of an invoicing profile with utilities.ValidateBillingName when it is created or updated. A failed check refuses the save. The rules of that check are in another file.

## Client records in use cannot be deleted

Lines 105-142

Asked as: record still used error; why can't I delete this; still used by collection; delete refused in use

Turbo refuses to delete a client, a client contact or an invoicing profile while another record still uses it. The message is "this record is still used by <collection>", with the name of the first collection found that uses it. This check also runs on deletes made outside requests and on cascaded deletes.

## Client alias cannot be cleared while profiles bill it

Lines 144-166

Asked as: cannot clear client alias; remove alias error; alias still billed; alias blank not allowed

Turbo refuses to clear a client's alias while invoicing profiles bill that alias. The message is "An invoicing profile bills this alias. Change it to bill the client name first." For several profiles it reads "N invoicing profiles bill this alias. Change them to bill the client name first." The check applies only to updates of an existing client that had an alias, and the count comes from utilities.CountProfilesBillingAlias.
