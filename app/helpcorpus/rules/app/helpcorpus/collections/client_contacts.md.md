<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/client_contacts.md","sha256":"1d250f41b8a0e21911c1e7bf0dec70d0f591848448bcf59621de5de889adc1db","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/client_contacts.md

## Who can list or view client contacts

Lines 7-17

Asked as: see client contacts list; who can view contacts; open a client contact; contacts visible to everyone

Any signed-in person can list client contacts and view one contact (list and view rules). Signing in is the only check; no claim is needed.

## Who can create or update a client contact

Lines 19-31

Asked as: edit client contact; change contact details; add new client contact; procurement cannot update contact

Only procurement (the job claim) can create or update a client contact through the API. The person must also have an active admin profile (active = true). Turbo's own server code is not limited by this rule.

## Who can delete a client contact, and when deletion is blocked

Lines 33-40

Asked as: cannot delete contact; remove client contact; contact still in use; delete contact blocked

Only procurement (the job claim) with an active admin profile can delete a client contact. Turbo refuses the deletion while any job, any project request or any client invoicing profile still uses the contact as its contact. Turbo's own server code is not limited by this rule.
