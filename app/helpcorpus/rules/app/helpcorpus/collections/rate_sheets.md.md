<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets.md","sha256":"5b0cb97b295f7bd32acb2fd76723ed50754bdec32acb8f4262e5fa649572c3d3","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets.md

## Who can list rate sheets

Lines 7-11

Asked as: list rate sheets; see all rate sheets; rate sheets empty; why can't I see rate sheets

Any signed-in person can list rate sheets through the API. A request that is not signed in cannot (rule: `@request.auth.id != ""`). The filter also requires `request = ""`, so the list passes only for records with no project request set.

## Who can view one rate sheet

Lines 13-17

Asked as: open a rate sheet; view one rate sheet; rate sheet not showing; can't open rate sheet

Any signed-in person can view one rate sheet through the API. The filter also requires `request = ""`, so a sheet with a project request set does not pass.

## Who can create a rate sheet

Lines 19-35

Asked as: make a new rate sheet; add rate sheet; who can create rate sheets; rate sheet permission denied; procurement new sheet revision

A person who is signed in and holds one of these claims can create a rate sheet: rate sheet revise (`rate_sheet_revise`), admin (`admin`), or procurement (the `job` claim). Procurement can create only a sheet whose revision (`revision`) is 0. The other two claims have no revision limit.

## Create request cannot set draft or approval fields

Lines 23-28

Asked as: can't set approver on sheet; set rates approved date; based on field rejected; rate sheet request field error

A create request must not send the project request (`request`), the base sheet (`based_on`), the reason (`reason`), the rates approver (`rates_approver`) or the rates approval date (`rates_approved`). A draft is made and approved only through a project request.

## Create request cannot set the code

Lines 29-30

Asked as: set rate sheet code; custom code rejected; rate sheet number wrong; assign own code

A create request must not send a code (`code`). The server gives a new sheet its code.

## Nobody can update or delete rate sheets through the API

Lines 37-43

Asked as: edit rate sheet; change saved rate sheet; delete rate sheet; update rates after saving

No one can update or delete a rate sheet through the API, except a PocketBase superuser. Turbo's own server code can still do both.

## Rate sheet field limits

Lines 47-61

Asked as: effective date invalid; rate sheet code format; reason too long; revision not a number

The effective date (`effective_date`) is required and must look like YYYY-MM-DD. The code (`code`) is required and must look like 2026-11-00004: a four-digit year, a month from 01 to 12, and five digits. The reason (`reason`) has a maximum of 500 characters. The revision (`revision`) must be a whole number.

## A code and revision pair is unique

Lines 63-68

Asked as: duplicate rate sheet code; code already used; same revision twice; rate sheet code taken

Two rate sheets cannot have the same code (`code`) and revision (`revision`) together. A revision of a sheet keeps its code and has a different revision number.
