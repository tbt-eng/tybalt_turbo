<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/project_authorization.go","sha256":"985d5964b194f19a7a5154de6b40e0d8656568f19086cab89c7f7b6fb86f5bd7","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/project_authorization.go

## Abandonment period for project requests is 60 days

Lines 16-16

Asked as: abandoned project request; how long until request expires; project request auto cancel days; when do stale requests close

Turbo sets the constant JOB_REQUEST_ABANDONMENT_DAYS to 60. This file only defines the number. The code that uses it is in another file.

## Duplicate-work review looks back 365 days

Lines 18-21

Asked as: duplicate work warning; old projects showing up for client; which past projects are listed; duplicate project review range

On a project request, the duplicate-work review lists a project of the request's client only when its award date or its latest time entry falls in the last 365 days (JOB_REQUEST_DUPLICATE_RECENT_DAYS).

## Project value must be a positive whole dollar amount up to a maximum

Lines 23-35

Asked as: project value rejected; approval amount not accepted; decimal dollar value error; project amount too large

When Turbo checks an approver's authority, the project value must be greater than 0, a whole dollar amount, and no more than 9007199254740991 (MaxProjectAuthorizationValue). Otherwise Turbo refuses with "project value must be a positive whole dollar amount".

## Project approver must be active

Lines 36-42

Asked as: inactive approver; approver account disabled; why can't this person approve; deactivated approver

Turbo refuses a project approver who is not active. The message is "approver must be active".

## Project approver needs the project authorization claim

Lines 43-49

Asked as: approver lacks permission; no project approval rights; cannot approve projects; missing approval privilege

Turbo refuses an approver who does not hold the project authorization claim (project_authorization). The message is "approver must have project_authorization authority".

## Project approver needs a value limit

Lines 50-56

Asked as: approver has no spending cap; no approval limit set; approver limit missing; no dollar limit for approver

Turbo refuses an approver who has the project authorization claim but no value limit record (project_authorization_props). The message is "approver has no value limit".

## Approver limit must cover the full project value

Lines 57-61

Asked as: approval limit too low; approver cap below project; limit not enough for project; exceeds approver spending limit

The approver's limit (max_value) must cover the full proposed value. A limit equal to the value passes, and a lower limit fails. Turbo also refuses a limit that is zero, negative or not a whole number, or above 9007199254740991. The message is "approver limit does not cover the full project value". When the check passes, it returns the limit.

## Project request needs email or no-email mark on contacts

Lines 64-84

Asked as: contact missing email; project request email required; no email address on contact; tick no email box

A project request's project contact and its billing contact (the contact on the invoicing profile) must each have an email address or be marked as having no email (no_email). A blank or whitespace-only email counts as missing. Otherwise Turbo refuses with "Add an email address for the project contact, or mark the contact as having no email address." (or "billing contact"). Contact email is otherwise optional.

## Billing contact is checked only when the invoicing profile exists

Lines 70-72

Asked as: billing contact not checked; no invoicing profile error; missing billing info on request; invoicing details absent

Turbo checks the billing contact only if the request's invoicing profile (invoicing_information) can be found. If the profile is not found, Turbo checks only the project contact and does not report an error.

## Missing contact stops the project request

Lines 73-77

Asked as: contact record not found; project contact deleted; billing contact missing from request; cannot find contact for request

If the project contact or the billing contact record cannot be found, Turbo refuses with "The project contact was not found." or "The billing contact was not found.".
