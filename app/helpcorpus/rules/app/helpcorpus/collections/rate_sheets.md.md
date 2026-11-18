<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets.md","sha256":"a808c9bd44b31b224775b503fe1795dfb808deb5ac2379d38414ca7307b62184","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets.md

## Who can list rate sheets

Lines 9-11

Asked as: list all rate sheets; see every rate sheet; rate sheets not showing; get rate sheets from API

Any signed-in person can list rate sheets through the API. The list rule also requires the request field to be empty (`request = ""`). A signed-out person cannot list them.

## Who can view a rate sheet

Lines 15-17

Asked as: open a rate sheet; view rate sheet details; can't see rate sheet; fetch one rate sheet

Any signed-in person can view one rate sheet through the API. The view rule also requires the request field to be empty (`request = ""`). A signed-out person cannot view one.

## Who can create a rate sheet through the API

Lines 22-33

Asked as: who can add rate sheet; new rate sheet via API; procurement create rates; make revision 0 sheet

Only procurement (the job claim) can create a rate sheet through the API, and only a new sheet with revision 0. A revision goes through the create route, not through this API action.

## Create cannot set request, based-on, reason or approval fields

Lines 23-28

Asked as: set approver on new sheet; create sheet already approved; pass reason when creating; draft rates without request

A create request must not send request, based_on, reason, rates_approver or rates_approved. A draft is made and approved only through a project request.

## Server gives a new rate sheet its code

Lines 29-30

Asked as: set my own rate sheet code; code rejected on create; choose sheet number; manual code for rates

A create request must not send a code. The server gives each new rate sheet its code.

## Nobody can update or delete a rate sheet through the API

Lines 38-42

Asked as: edit saved rate sheet; change rates after saving; delete a rate sheet; rate sheet locked from changes

The API has no update rule and no delete rule for rate sheets. Nobody can update or delete a rate sheet through the API except a PocketBase superuser. Turbo's own server code can still do both.

## Rate sheet field limits

Lines 49-56

Asked as: rate sheet date invalid; effective date error; reason too long; revision must be number

The effective date (effective_date) is required and must look like YYYY-MM-DD. The reason (reason) has a maximum of 500 characters. The revision (revision) must be a whole number.

## Rate sheet code format and uniqueness

Lines 60-67

Asked as: duplicate rate sheet code; code already exists; wrong code format; same code and revision error

The code is required and must look like 2026-11-00004: four digits, a month from 01 to 12, then five digits. Two rate sheets cannot have the same code and revision together (unique index idx_rate_sheets_code_revision).
