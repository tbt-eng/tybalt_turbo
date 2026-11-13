<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_absorb.go","sha256":"afba3d56ac6c7629f7b213525e07cbba60b6b9f54741c28cba29fb1e721c3177","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_absorb.go

## Any signed-in user can list duplicate rate sheet groups

Lines 11-23

Asked as: duplicate rate sheets; list duplicate rates; rate sheets with same rates; who can see duplicate groups

The duplicates list is read-only, so any signed-in user can read it. It returns every group of rate sheets with identical rates, active and inactive. Each group has a target and the reason when a sheet can't be absorbed into it. The absorb route checks the claims, not this list.

## Absorb candidates for a target rate sheet

Lines 25-41

Asked as: absorb page empty; which sheets can merge; candidates to merge into target; why can't I absorb this sheet

Turbo returns the target rate sheet and the sheets with identical rates. Each candidate carries a reason when the target can't absorb it. The absorb page offers only the candidates without a reason. This route does not check claims in this file.

## Absorb candidates refuses an unknown rate sheet

Lines 34-36

Asked as: rate sheet not found; unknown rate sheet id; absorb error wrong id; rate sheet missing error

When the target rate sheet id is not found, Turbo refuses the request with the message "rate sheet not found".
