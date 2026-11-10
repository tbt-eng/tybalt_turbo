<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/client_setup.go","sha256":"34e317d366a04742a20c098ce5d3e7905379fb510536d3b15b0340d3fdd8c9f8","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/client_setup.go

## Which records block deleting a client

Lines 15-38

A client cannot be deleted while a job (client or job owner), a project request (job_requests, client or job owner), an invoicing profile (client_invoicing_information) or a client note uses it. Contacts are not on this list because they are deleted with their client, and each contact's own uses are checked as it is deleted.

## Which records block deleting a contact

Lines 29-33

A client contact cannot be deleted while a job, a project request or an invoicing profile uses it as its contact.

## Which records block deleting an invoicing profile

Lines 34-37

An invoicing profile cannot be deleted while a job or a project request uses it (invoicing_information).

## Invoicing profile creator is set by Turbo

Lines 43-51

When an invoicing profile is created through a request, Turbo sets its creator (creator) to the signed-in user's id. A value from the form is overwritten. The creator is empty when the requester is not a user (not in the users collection).

## Invoicing profile creator cannot be changed

Lines 53-64

Turbo refuses an update to an invoicing profile that changes its creator (creator). The message is "the creator of an invoicing profile cannot be changed".

## Contact with an email is not marked as having no email

Lines 66-70

When a client contact is created or updated with a non-blank email, Turbo sets no_email to false. Spaces alone count as blank. Turbo does not change no_email when the email is blank.

## Contacts and invoicing profiles have their client references validated

Lines 71-77

Turbo validates the client references of every client contact and invoicing profile on create and update, using utilities.ValidateClientReferences. It refuses the save if that check returns an error. The rules of that check are in another file.

## Invoicing profile billing name is validated on create and update

Lines 81-101

Turbo runs utilities.ValidateBillingName on an invoicing profile whenever it is created or updated, and refuses the save if the check fails. The request hooks return the error with its field codes. Saves made outside requests are checked too. The rules of that check are in another file.

## Client edits are checked for alias removal

Lines 81-102

Turbo checks every update to a client for alias removal. It does not run this check on create. It applies to requests and to saves made outside requests.

## Record in use cannot be deleted

Lines 104-139

Turbo refuses to delete a client, a client contact or an invoicing profile while another record still uses it. The message is "this record is still used by <collection>", with the name of the first collection found that uses it. The check applies to requests, to saves outside requests and to cascaded deletes.

## Alias cannot be cleared while profiles bill it

Lines 143-164

Turbo refuses to clear a client's alias if invoicing profiles bill that alias. The message is "An invoicing profile bills this alias. Change it to bill the client name first." For several profiles it is "<n> invoicing profiles bill this alias. Change them to bill the client name first."

## Alias check applies only when an existing alias is cleared

Lines 146-148

Turbo skips the alias check for a new client. It also skips it when the client still has an alias or had none before. The check only runs when an alias that was set is removed.
