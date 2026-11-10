<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/client_invoicing_information.md","sha256":"047a05ce27c9c00927980952c4a0a2eed2075fc97047838a3d1f4e90bc3ff333","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/client_invoicing_information.md

## Who can list or view invoicing profiles

Lines 7-17

Asked as: who can see invoicing profiles; list invoicing profiles; view client billing details; read invoicing info without claim

Any signed-in person with an active admin profile (admin_profiles active = true) can list invoicing profiles (client_invoicing_information) and view one. No claim is needed to read them.

## Who can create an invoicing profile

Lines 19-23

Asked as: add new invoicing profile; create billing profile for client; why can't I make invoice profile; time claim create invoicing

A signed-in person with an active admin profile can create an invoicing profile (client_invoicing_information) if they hold the job claim (procurement) or the time claim. Anyone else is refused.

## Who can update an invoicing profile

Lines 25-29

Asked as: edit invoicing profile; change billing details; update client invoicing info; can't save invoicing changes

Only a signed-in person with an active admin profile and the job claim (procurement) can update an invoicing profile (client_invoicing_information). The time claim does not allow updates, although it allows creating.

## Delete an invoicing profile only when no job or project request uses it

Lines 31-35

Asked as: delete invoicing profile; remove billing profile; cannot delete invoicing profile in use; invoice profile still linked to job

Only a signed-in person with an active admin profile and the job claim (procurement) can delete an invoicing profile. The delete is refused when any job (jobs.invoicing_information) or any project request (job_requests.invoicing_information) still uses the profile.

## Required fields on an invoicing profile

Lines 44-47

Asked as: invoicing profile required fields; missing client or contact; name too long billing; character limit invoicing name

Every invoicing profile must name a client and a contact (client, contact), one each. The name field holds at most 120 characters.

## Billing name choice on an invoicing profile

Lines 50-50

Asked as: billing name options; billing name name or alias; which billing name to pick; billing name value rejected

The billing name (billing_name) accepts only one value: 'name' or 'alias'.
