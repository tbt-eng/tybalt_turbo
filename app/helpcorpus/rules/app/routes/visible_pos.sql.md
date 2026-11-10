<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/visible_pos.sql","sha256":"485d26c9f5f88fd10ff6323d3de28a497cad4cdebe5364b58efc678a5bb21f4d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/visible_pos.sql

## PO list shows approval, currency and spending details

Lines 4-77

Asked as: PO list columns; purchase order totals shown; remaining amount on PO list; who approved the PO list; PO print max amount

The visible purchase orders list returns, for each PO, its status, total, approval total (approval_total), remaining amount (remaining_amount) and maximum print amount (print_max_amount). It also returns the committed expenses count and expenses total, the recurring expected and remaining occurrences, and the names of the approver, second approver, rejector, canceller and closer. Turbo takes these values from the base visibility query, and this file does not calculate them.

## Only purchase orders passing a visibility check are listed

Lines 78-87

Asked as: why can't I see a PO; missing purchase orders in list; PO not showing up; hidden POs; which POs appear in list

The purchase order list returns a PO only if at least one visibility flag is 1: active, closed/cancelled, unapproved direct, unapproved second-stage eligible, revision, legacy, or group. The flags are computed in the visibility base (__PO_VISIBILITY_BASE__), which this file does not show. The chosen scope then narrows the list further.

## Scope 'all' adds no further filter

Lines 88-89

Asked as: show every PO; see all purchase orders; no filter on POs; all POs for everyone

With scope 'all', Turbo applies no extra filter beyond the visibility check. The user sees every PO the visibility flags allow.

## Scope 'mine' shows only the user's own POs

Lines 90-90

Asked as: only my POs; POs I created; my purchase orders; POs under my name

With scope 'mine', Turbo lists only POs whose creator (uid) is the current user.

## Scope 'group' shows POs in the given PO group

Lines 91-91

Asked as: POs in a group; filter POs by group; group purchase orders; PO group filter

With scope 'group', Turbo lists only POs whose po_group equals the requested group ID.

## Scope 'active' includes POs under revision

Lines 92-92

Asked as: active POs under revision; POs being revised; revision POs in active list; revised purchase orders showing

With scope 'active', Turbo lists POs with status 'Active'. It also lists POs the visibility base flags as revision-visible (is_revision_visible = 1), whatever their status.

## Scope 'rejected' shows the user's own rejected unapproved POs

Lines 93-98

Asked as: my rejected POs; POs I got rejected; rejected unapproved purchase orders; why was my PO rejected

With scope 'rejected', Turbo lists POs with status 'Unapproved' that have a rejection date (rejected not empty) and were created by the current user.

## Scope 'stale' finds active POs approved before a cutoff

Lines 99-106

Asked as: old active POs; POs approved long ago; stale purchase orders; POs approved before date; outdated approved POs

With scope 'stale', Turbo lists 'Active' POs. If the PO has a second approval date, that date must be before the cutoff (staleBefore). The check on the approval date (approved < staleBefore) is an alternative, so a PO whose first approval is older than the cutoff also matches even when its second approval is recent. The cutoff is a strict comparison, so a PO approved exactly at the cutoff is not stale.

## Scope 'expiring' lists active recurring POs ending soon

Lines 107-113

Asked as: recurring POs ending soon; POs about to expire; expiring purchase orders; recurring PO end date; POs running out

With scope 'expiring', Turbo lists POs with status 'Active' and type 'Recurring' that have an end date. The end date must be on or before the cutoff (expiringBefore). A PO ending exactly on the cutoff is included. POs with no end date are excluded.

## Scope 'approved_by_me_awaiting_second' lists POs waiting for a second approver

Lines 114-121

Asked as: POs waiting for second approver; I approved, need second sign-off; approved by me pending; awaiting second approval; my approvals still waiting

With scope 'approved_by_me_awaiting_second', Turbo lists POs with status 'Unapproved' that are not rejected and whose first approver is the current user. The user has approved them (approved not empty), but no second approval exists yet.

## List sort order depends on scope

Lines 123-127

Asked as: PO list order; sort purchase orders by date; newest POs first; why POs in this order; expiring POs sorted

The 'expiring' scope sorts by end date, earliest first. The 'approved_by_me_awaiting_second' scope sorts by approval date, newest first. After that, and for all other scopes, POs sort by date newest first, then by last update newest first.

## A limit of zero or less returns all rows

Lines 128-128

Asked as: PO list limit zero; show all rows no limit; how many POs returned; limit on purchase orders; unlimited PO results

If the limit is greater than 0, Turbo returns at most that many POs. A limit of 0 or less means no limit.
