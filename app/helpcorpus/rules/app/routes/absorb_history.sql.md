<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_history.sql","sha256":"ad3f1f65f473361cebc88f8e7399b4e63aa740781f82ffe0cb78f83c47dee2f7","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_history.sql

## Absorb history shows when each absorb was committed

Lines 10-12

Asked as: absorb committed time; when was absorb saved; absorb timestamp wrong; absorbed at versus committed

The absorb history list gives the time a record was created in the absorb_history table as the committed time (committed_at). It gives the absorbed time (absorbed_at) as a separate field.

## Absorb history shows names of the absorber and the committer

Lines 13-17

Asked as: who absorbed it; absorber name missing; committer name blank; absorb history user names

The absorb history list shows the full name of the person who absorbed (absorbed_by_name) and of the person who committed (committed_by_name). Each name is the given name and the surname from that person's profile, matched on their user ID. When no profile matches, the name is an empty string and the row still appears.
