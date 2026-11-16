<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets.md","sha256":"087e230634f2b30207bd84d607de3389720c5d3e61910b7066abbb7405d2a2a5","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets.md

## Who can list rate sheets

Lines 7-11

Asked as: list rate sheets; see all rate sheets; who can see rate sheets; rate sheet index access

Any signed-in person can list rate sheets through the API. The request must be a normal request (request = ""). People who are not signed in cannot list them.

## Who can view one rate sheet

Lines 13-17

Asked as: open one rate sheet; view a single rate sheet; rate sheet detail access; can't open rate sheet

Any signed-in person can view a single rate sheet through the API. The request must be a normal request (request = ""). People who are not signed in cannot view it.

## Who can create a rate sheet and which fields are blocked

Lines 19-33

Asked as: who can make a rate sheet; create new rate sheet permission; procurement add rate sheet; can't create rate sheet; reviser claim for rates

A signed-in person can create a rate sheet only if they hold one of three claims: a rate sheet reviser (the rate_sheet_revise claim), an admin (the admin claim), or procurement (the job claim). Procurement can create one only when revision is 0. A reviser or an admin can create one with any revision.

## Created rate sheet cannot set project request fields

Lines 22-28

Asked as: set approved flag on rate sheet; rate sheet created already approved; fill approver when creating; hand-enter rate sheet approval

A create request must not send request, based_on, reason, rates_approver or rates_approved. Turbo makes and approves a draft rate sheet only through a project request. Nobody can set these fields when creating a rate sheet through the API, whatever their claim.

## Nobody can update a rate sheet through the API

Lines 35-37

Asked as: edit saved rate sheet; change rates after saving; rate sheet won't update; modify existing rate sheet

The update rule is empty, so nobody can update a rate sheet through the API, except a PocketBase superuser. Turbo's own server code can still change one.

## Nobody can delete a rate sheet through the API

Lines 39-41

Asked as: remove a rate sheet; rate sheet won't delete; get rid of old rate sheet; erase rate sheet

The delete rule is empty, so nobody can delete a rate sheet through the API, except a PocketBase superuser. Turbo's own server code can still delete one.

## Rate sheet field limits

Lines 47-56

Asked as: rate sheet date format; rate sheet effective date error; rate sheet reason too long; revision must be number

The effective date (effective_date) is required and must look like YYYY-MM-DD. The name is required. The revision (revision) must be a whole number. The reason (reason) has a maximum of 500 characters.

## Rate sheet name and revision must be unique together

Lines 64-64

Asked as: duplicate rate sheet name; same rate sheet twice; rate sheet already exists; new revision number clash

Two rate sheets cannot have the same name and revision (unique index on name and revision). A new revision of a name needs a different revision number.
