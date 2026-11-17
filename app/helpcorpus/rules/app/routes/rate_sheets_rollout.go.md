<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rollout.go","sha256":"a39ed61dec3c9cff4527d425d0a851c4fc7a01c391c1b8fc0795c30745c0cc43","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rollout.go

## Rollout moves active projects from one date to the latest revision

Lines 15-20

Asked as: update jobs to new rates; projects still on old rate sheet; move jobs to new price list; reprice jobs after rate change

A revision rollout gives every active project on an earlier revision of a rate sheet a new rate sheet period on the latest revision, from the chosen date. Each project gets the period in the same way as a new period on the job details page. Time from that date is repriced.

## Only active admins can move projects to a rate sheet revision

Lines 45-58

Asked as: can't move projects to revision; not allowed to move jobs; need admin rights for rollout; permission denied moving projects

Both the preview count and the rollout need an active admin (the admin claim). Anyone else gets "Only an active admin can move projects to a revision."

## Rollout target must be the latest, active revision

Lines 63-79

Asked as: move projects to old revision; rate sheet revision not active; newer version of rate sheet exists; activate rate sheet before moving

The target revision must exist, or Turbo says "Rate sheet not found." If a newer revision of the same code exists, Turbo says "Projects can move only to the latest revision of a rate sheet." If the revision is not active, Turbo says "Activate this revision before you move projects to it."

## Rollout skips proposals, cancelled, closed and later-period projects

Lines 99-130

Asked as: why some jobs were skipped; cancelled jobs not moved; closed projects left behind; project skipped in rollout

The rollout looks at projects (not proposals, whose number starts with P) whose current sheet has the same code and an earlier revision. A Cancelled project is counted as cancelled (cancelled). Any other project that is not Active is counted as closed (closed). An Active project whose latest period starts on or after the date is counted as later_period, because a period from the date would not be the latest. All other projects move (move).

## Rollout is refused while a rate sheet absorb is pending

Lines 134-146

Asked as: rollout blocked by absorb; rate sheet absorb waiting; can't move projects pending absorb; undo or commit absorb first

If any rate sheet absorb can still be undone or committed, Turbo refuses the rollout with "A rate sheet absorb is waiting to be committed or undone. Finish it before moving projects." Undo would overwrite a moved project.

## Rollout preview defaults the date to the revision's effective date

Lines 148-173

Asked as: preview projects to move; count before rollout; which effective date is used; dry run rate sheet change

The preview counts the projects that would move and those skipped, without changing anything. Only an active admin can use it. If no effective_date is given, Turbo uses the effective date of the revision. The date must pass the new-period date and revision checks.

## Rollout checks date and reason before moving projects

Lines 187-211

Asked as: rollout reason rejected; invalid date for rate change; rollout date error; reason field validation rate sheet

Turbo checks that the date passes the new-period date rules and the revision rule for the sheet. It also checks the reason with the rate sheet reason check (checkRateSheetReason). The reason and date are trimmed of spaces first.

## Rollout needs a confirmed project count

Lines 189-191

Asked as: confirm number of projects; rollout missing count; expected count required; how many projects to confirm

A rollout request without expected_count is refused with "Confirm the number of projects to move."

## Rollout moves all projects or none

Lines 194-246

Asked as: some projects failed rollout; rollout failed halfway; partial rollout left behind; all projects move or nothing

The rollout runs in one transaction. Each project gets a new dated period on the sheet, with the reason and the admin who did it. If any project fails, none move. On success Turbo returns the counts.

## Rollout is refused when job editing is disabled

Lines 195-201

Asked as: jobs editing disabled rollout; can't edit jobs during rollout; job editing turned off; ErrJobsEditingDisabled

The rollout runs only while job editing is enabled. If it is disabled, Turbo refuses the rollout with the jobs-editing-disabled error (ErrJobsEditingDisabled).

## Rollout refuses when the project count differs from the confirmed count

Lines 215-221

Asked as: projects to move have changed; count changed since preview; new count try again; number of jobs changed rollout

Turbo recounts the projects to move inside the transaction. If the count differs from expected_count, nothing moves. Turbo answers "The projects to move have changed. Check the new count and try again." with the new counts.
