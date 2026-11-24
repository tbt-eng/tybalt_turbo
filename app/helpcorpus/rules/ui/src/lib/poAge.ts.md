<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poAge.ts","sha256":"a7c6531728cef9a4611be45a022f72556fccd2edd2de8b2b7346fbe0150b2601","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poAge.ts

## Active purchase order becomes stale after 30 days

Lines 7-24

Asked as: PO approval too old; purchase order out of date; old PO still active; approval older than 30 days

An Active purchase order is stale when its approval is more than 30 days old (STALE_AFTER_DAYS = 30). The cutoff is the current time minus 30 days, in UTC. Only POs with status Active can be stale.

## Recurring PO is expired when end date is before today

Lines 26-32

Asked as: recurring PO ended; PO past end date; purchase order expired; recurring order no longer valid

An Active Recurring purchase order is expired when its end date is before today's date. Turbo takes today's date from the user's own device, not from UTC. A PO that ends today is not expired.

## Recurring PO is expiring within 30 days, today included

Lines 34-56

Asked as: PO ending soon; recurring PO expiry warning; purchase order running out; PO end date approaching

An Active Recurring purchase order is expiring when its end date is today or no later than 30 days from today (EXPIRING_WITHIN_DAYS = 30). A PO that ends exactly 30 days from today is expiring. A PO that ends later has no label. The cutoff uses the user's local calendar date.

## Stale check uses the older of two approvals

Lines 42-48

Asked as: second approval date old; PO flagged stale; which approval counts for stale; stale PO second approver

For an Active purchase order, Turbo checks the approval date (approved) and the second approval date (second_approval). The PO is stale when either one is before the cutoff. Turbo skips the second approval when it is empty.

## Only Active Recurring POs with an end date can expire

Lines 50-53

Asked as: why no expiry label; one-time PO expiring; PO has no end date; cumulative PO expired label

Turbo gives an end state only to a purchase order that is Active, has type Recurring and has an end date. Every other PO has no expiring or expired label. A One-Time PO, a Cumulative PO, a PO that is not Active and a PO with no end date never expire.
