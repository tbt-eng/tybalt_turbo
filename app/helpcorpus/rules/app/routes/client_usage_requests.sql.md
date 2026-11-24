<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_usage_requests.sql","sha256":"6628aaa0b92a81dea16e00538089ff1ef9da4046ebc038afec98a49e0998662b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_usage_requests.sql

## Which invoicing profiles count for the usage check

Lines 4-8

Asked as: invoicing profile for client; which billing profile counts; profile linked to contact; client billing profiles check

The profiles considered are those of the given client. A profile counts if it is the chosen profile, or if a contact is given and the profile's contact is that contact. Profiles of other clients never count.

## Usage rows show an empty approver when none is set

Lines 10-14

Asked as: approver blank in list; no approver shown; empty approver field; request missing approver

Each row gives the request's id, description, status, creator and approver. A request with no approver shows an empty string for approver. The handler uses the creator and approver to leave out requests the caller may not read.

## Each usage row shows the contact and invoice recipient roles

Lines 15-16

Asked as: project contact label; invoice recipient tag; which role request has; contact and recipient both

Each request is marked as project contact (project_contact) when a contact is given and the request's contact is that contact. It is marked as invoice recipient (invoice_recipient) when its invoicing profile is one of the matching profiles. A request can have both roles.

## Which project requests count as using a client's contact or profile

Lines 18-19

Asked as: requests using client contact; job requests for contact; requests tied to billing profile; which jobs use this client

Turbo lists a project request (job_requests) when its contact is the chosen contact. It also lists a request whose invoicing profile (invoicing_information) is one of the chosen profiles. The contact test applies only when a contact is given.

## Usage requests are sorted newest first

Lines 20-20

Asked as: order of usage requests; newest requests first; sort requests by date; list order of requests

Turbo sorts the requests by creation date, newest first. Requests created at the same time are sorted by id.
