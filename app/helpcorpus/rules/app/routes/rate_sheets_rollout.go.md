<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rollout.go","sha256":"50e7666078f666b0c695961fde20feb7df2ce8eb3e80f715690f6a7dfb37bb12","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rollout.go

## Only admins can move projects to a rate sheet revision

Lines 45-54

Asked as: only admins can move projects; non-admin rollout denied; move projects permission; admin claim needed for rollout

Both the preview count and the rollout itself need the admin claim. Anyone else gets "Only an admin can move projects to a revision."

## Rollout target must be the latest revision of its sheet

Lines 59-70

Asked as: rate sheet not latest revision; newer revision exists; move to old revision; rate sheet not found

Projects can move only to the latest revision of a rate sheet (same code). If a newer revision exists, Turbo refuses with "Projects can move only to the latest revision of a rate sheet." If the sheet id is not found, Turbo says "Rate sheet not found."

## Rollout target revision must be active

Lines 71-73

Asked as: revision not active; activate revision first; inactive rate sheet revision; cannot move to inactive revision

If the revision is not active, Turbo refuses with "Activate this revision before you move projects to it."

## Rollout date follows the new-period date rules

Lines 77-83

Asked as: rollout date rejected; wrong date for rate sheet move; period date check failed; effective date not allowed

The rollout date must pass the same date check as a new rate sheet period (CheckRateSheetPeriodDate) and the revision check for that date (CheckRateSheetPeriodRevision). The rules and messages for these checks are in the hooks, not in this file.

## Which projects move and which are skipped

Lines 88-128

Asked as: which projects get skipped; closed projects not moved; cancelled jobs excluded; project skipped later period; why project did not move

Turbo looks at non-proposal jobs (number does not start with P) whose current sheet has the same code and an earlier revision. Cancelled projects are counted as Cancelled and Closed (any other non-Active) as Closed; neither moves. An Active project whose latest dated period starts on or after the rollout date is skipped as LaterPeriod, so a project with a period starting exactly on that date is skipped. All other Active projects move.

## Rollout is blocked while a rate sheet absorb is pending

Lines 130-142

Asked as: absorb pending blocks rollout; cannot move projects absorb; undo absorb first; absorb waiting to commit

If any rate sheet absorb action is still waiting to be committed or undone, Turbo refuses the rollout with "A rate sheet absorb is waiting to be committed or undone. Finish it before moving projects." Undoing the absorb would overwrite moved projects.

## Preview count defaults to the revision's effective date

Lines 144-169

Asked as: preview how many projects move; dry run rate sheet change; default effective date; count projects before moving

The preview shows how many projects would move and how many are skipped (closed, cancelled, later period). If no effective_date is given, Turbo uses the revision's own effective date. The preview changes nothing.

## Rollout requires a confirmed project count

Lines 185-187

Asked as: rollout without count; confirm project number; expected count missing; must enter project count

A rollout request without the expected count is refused with "Confirm the number of projects to move."

## Rollout is refused when job editing is disabled

Lines 191-197

Asked as: cannot move projects job editing off; jobs editing disabled rollout; rollout locked editing; editing turned off error

The rollout runs only while job editing is enabled. Otherwise Turbo returns the jobs-editing-disabled error (ErrJobsEditingDisabled); its message is defined elsewhere.

## Rollout requires a valid reason

Lines 205-207

Asked as: rollout reason missing; reason needed to move projects; why reason rejected; rollout note required

The rollout checks the reason with checkRateSheetReason after the date check. The rules for the reason are defined elsewhere.

## Rollout is refused if the count changed since confirmation

Lines 215-217

Asked as: project count changed; numbers changed since preview; count mismatch rollout; projects to move changed

If the number of projects to move differs from the count the admin confirmed, nothing moves. Turbo answers "The projects to move have changed. Check the new count and try again." with the new counts.

## Rollout moves all projects or none

Lines 218-242

Asked as: partial rollout failed; some projects not moved; all or nothing move; rollout rolled back; new period after rollout

In one transaction, each moving project gets a new dated rate sheet period on the new revision from the rollout date, with the reason and the admin recorded, using the same history rules as a new period on the job page. If any project fails, none move. On success Turbo returns the counts.
