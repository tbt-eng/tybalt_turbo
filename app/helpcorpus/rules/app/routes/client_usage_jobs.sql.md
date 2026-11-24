<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_usage_jobs.sql","sha256":"a476ed80a2b92e2f3864ca4d6eacce8037f1e0788fdcb21d978ea169813df733","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_usage_jobs.sql

## Jobs listed for a client contact or invoicing profile

Lines 6-21

Asked as: jobs for a client contact; which jobs use this contact; jobs linked to an invoicing profile; list client jobs by contact

Turbo lists the jobs that use one contact or one invoicing profile of a client. The caller supplies the contact or the profile and leaves the other empty. Each row shows the job's number, description and status.

## Invoicing profiles matched when a contact is given

Lines 6-10

Asked as: invoices sent to a contact; profiles that bill this contact; who receives invoices for jobs; invoice recipient jobs

When a contact is given, Turbo counts the client's invoicing profiles that send invoices to that contact. Jobs that use those profiles are listed as invoice recipients. Only profiles of the same client are considered.

## Empty contact matches no project contact

Lines 9-19

Asked as: blank contact no jobs; contact left empty; jobs missing project contact; why contact search shows nothing

An empty contact never matches a job as project contact. Jobs then appear only through the invoicing profile.

## Ways a job uses a contact

Lines 11-20

Asked as: job listed twice; project contact and invoice recipient; same job under two roles; contact role on job

A job uses a contact as its project contact (project_contact), or as the recipient of invoices through its invoicing profile (invoice_recipient), or both. A job with both uses is listed once, with both roles marked.

## Job list order

Lines 21-21

Asked as: job list sort order; jobs sorted by number; newest job first; order of job numbers

Turbo sorts the jobs by job number, highest first, then by id.
