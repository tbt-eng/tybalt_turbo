<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/rate_sheets_rollout.go","sha256":"7a8b5d41c5c53d8ccc2038ac8734ee64a131bee6e11e7f76318916afaf2c0b25","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/rate_sheets_rollout.go

## Only an active admin can use rate sheet rollout

Lines 45-58

Asked as: not an admin; only admins can move projects; rollout permission denied; who can move projects to a revision

Turbo refuses both the rollout preview and the rollout itself unless the person is an active user with the admin claim. The message is "Only an active admin can move projects to a revision."

## Rollout target must be the latest active revision

Lines 63-79

Asked as: rate sheet not found; move to old revision; activate revision first; revision not latest

Projects move only to a rate sheet revision that exists, is the latest revision of its name, and is active. Turbo refuses with "Rate sheet not found.", "Projects can move only to the latest revision of a rate sheet." or "Activate this revision before you move projects to it."

## Rollout moves only active non-proposal projects on earlier revisions

Lines 92-110

Asked as: which projects get moved; proposal jobs excluded; jobs on older revision; projects eligible for rollout

A revision rollout looks at jobs whose current rate sheet has the same name and a lower revision (rs.revision < the chosen revision). Jobs whose number starts with P (proposals) are left out. Projects are counted by status and by their latest period.

## Rollout sorts projects into move and skipped groups

Lines 115-130

Asked as: why was project skipped; project not moved; cancelled project rollout; later period skipped

Turbo counts a Cancelled project as cancelled (cancelled). It counts any other project that is not Active (Closed) as closed (closed). An Active project whose latest dated period starts on or after the chosen date counts as later_period, so the equal date is skipped. All other Active projects move (move).

## Rollout is blocked while a rate sheet absorb is pending

Lines 137-146

Asked as: absorb pending; rate sheet absorb waiting; cannot move projects absorb; undo absorb before rollout

While any rate sheet absorb is waiting to be committed or undone, Turbo refuses the rollout with "A rate sheet absorb is waiting to be committed or undone. Finish it before moving projects." Undo would overwrite moved projects.

## Rollout preview date defaults to the revision's effective date

Lines 151-172

Asked as: rollout preview; how many projects would move; default effective date; preview rate sheet change

The preview (GET) counts the projects that would move and those skipped. When no effective_date is given or it is blank, Turbo uses the revision's own effective date. The date must pass the new-period date rules (CheckRateSheetPeriodDate and CheckRateSheetPeriodRevision). The preview changes nothing.

## Rollout trims spaces from date and reason

Lines 187-188

Asked as: extra spaces in date; trim reason spaces; whitespace in rollout reason; spaces around effective date

When an admin moves projects to a revision, Turbo trims leading and trailing spaces from the effective date (effective_date) and from the reason. It checks and saves the trimmed values.

## Rollout request has no default date

Lines 187-187

Asked as: blank date rollout; rollout date empty; no date when moving projects; date required to move projects

The rollout itself does not use the revision's effective date when the effective date is empty, unlike the preview. Turbo passes the empty date to the same date check as any other date (checkRolloutDate), so the admin should send a date.

## Rollout needs a confirmed project count

Lines 189-191

Asked as: confirm project count; expected count missing; number of projects required; rollout count not confirmed

A rollout request without expected_count is refused with "Confirm the number of projects to move." The admin must give the count that the preview showed.

## Rollout is refused when job editing is off

Lines 195-201

Asked as: job editing disabled; cannot edit jobs rollout; rollout blocked editing off; jobs editing turned off

The rollout checks inside its transaction that job editing is enabled. When it is disabled, Turbo refuses with the jobs-editing-disabled error (ErrJobsEditingDisabled).

## Rollout checks date and reason before moving

Lines 206-211

Asked as: rollout date invalid; reason rejected rollout; rollout failed date check; nothing moved after error

The rollout checks the effective date with the new-period date rules. It checks the reason with the rate sheet reason rule (checkRateSheetReason). A failure on either stops the rollout and moves nothing.

## Rollout refuses when the project count changed

Lines 215-221

Asked as: project count changed; numbers changed rollout; count mismatch projects; projects changed since preview

If the number of projects to move differs from the admin's expected count, nothing moves. Turbo answers "The projects to move have changed. Check the new count and try again." with the new counts.

## Rollout moves all projects in one transaction

Lines 222-246

Asked as: all or nothing rollout; projects added new period; rollout partly failed; rollout success counts

For each project to move, Turbo adds a new dated rate sheet period for the latest revision, from the chosen date, with the given reason and the admin as the creator. It uses the same history rules as a new period on the job page. If any project fails, none move. On success Turbo returns the counts.
