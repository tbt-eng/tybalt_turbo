<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets.md","sha256":"4624a1c3d9c7634b3f6f372dafb30899eb6e0033fdc4811147b6140e819b35dc","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets.md

## Who can list or view rate sheets

Lines 9-17

Asked as: list rate sheets; view rate sheet API; who can see rate sheets; rate sheet not showing

A signed-in person can list rate sheets and view one rate sheet through the API. The request must not carry a `request` value (the filter `request = ""`). Someone who is not signed in cannot list or view them.

## Who can create a rate sheet directly

Lines 22-30

Asked as: create rate sheet permission; procurement create rate sheet; which claim creates rate sheets; revision 0 only procurement

A signed-in person can create a rate sheet only if they hold one of three claims. The claims are rate sheet revise (`rate_sheet_revise`), admin (`admin`), or procurement (the job claim). A procurement person can create one only when the body sets `revision` to 0. A person with the rate sheet revise or admin claim can set any revision.

## Direct creation cannot set request, based on or reason

Lines 23-26

Asked as: create rate sheet with reason; rate sheet based on another; cannot send request field; make a draft directly

A create request must not send `request`, `based_on` or `reason`. A rate sheet draft is made only through a project request. So nobody can create a draft, a revision based on another sheet, or a sheet with a reason through this direct route.

## Nobody can update or delete a rate sheet through the API

Lines 35-39

Asked as: edit rate sheet API; delete rate sheet; update rate sheet forbidden; change rate sheet after saving

The API has no update rule and no delete rule for rate sheets. Nobody can update or delete a rate sheet there, except a PocketBase superuser. Turbo's own server code can still do both.

## Rate sheet field limits

Lines 45-54

Asked as: rate sheet required fields; effective date format; reason too long; revision must be number

A rate sheet needs an effective date (`effective_date`) in the form YYYY-MM-DD and a name. The revision (`revision`) must be a whole number. The reason (`reason`) holds at most 500 characters. The request links to one project request, and the based-on sheet (`based_on`) links to one rate sheet.

## Name and revision must be unique together

Lines 60-60

Asked as: duplicate rate sheet name; same revision number error; rate sheet already exists; new revision number taken

Two rate sheets cannot have the same name and revision number together (a unique index on `name` and `revision`). A new revision of a name needs a revision number that no sheet of that name has.
