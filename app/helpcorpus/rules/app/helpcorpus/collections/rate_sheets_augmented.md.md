<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/rate_sheets_augmented.md","sha256":"bf97fa69b6341aed09a4a8f639443b4c49a6bfeee713020e0c95225f2cef60d2","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/rate_sheets_augmented.md

## Rate sheets view is read-only through the API

Lines 5-5

Asked as: cannot edit rate sheet; change rate sheet via api; rate sheet read only; delete rate sheet api

rate_sheets_augmented is a view, and it has no create, update or delete rules. Nobody can create, update or delete its records through the API. Turbo reads it and never writes to it.

## Who can list rate sheets with job counts

Lines 7-11

Asked as: list rate sheets; show all rate sheets; rate sheets not loading; see rate sheets without login

Any signed-in person can list the records of the rate_sheets_augmented view (the list rule is `@request.auth.id != ""`). A request without a signed-in user is refused.

## Who can view one rate sheet with its job count

Lines 13-17

Asked as: open one rate sheet; view single rate sheet; rate sheet details not showing; rate sheet access denied

Any signed-in person can view one record of the rate_sheets_augmented view (the view rule is `@request.auth.id != ""`). A request without a signed-in user is refused.

## Job count of a rate sheet counts the jobs that use it

Lines 21-34

Asked as: job count zero; how many jobs use rate sheet; rate sheet unused jobs; count of jobs per rate sheet

Turbo calculates job_count as the number of jobs whose rate sheet (rate_sheet) is this rate sheet. A rate sheet that no job uses shows a job_count of 0, because the view joins jobs with a LEFT JOIN. The view lists every rate sheet, active or not.
