<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/po_approvers.go","sha256":"53aa57471fe5a43012383451058cb62c2edfef02416370e449bbaf8c666b886b","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/po_approvers.go

## First-stage approver check uses the first-stage pool

Lines 41-44

Asked as: am I a first approver; can I approve this PO first; first sign off rights; not on PO approver list

Turbo treats a person as a first-stage approver for a purchase order (PO) only if they are in that PO's first-stage approver pool (firstStageLimits). Anyone outside the pool is not a first-stage approver. The pool depends on the PO's division, amount, expenditure kind and whether it has a job.

## Second-stage approver check uses the second-stage pool

Lines 46-49

Asked as: second approver for this PO; can I give second sign off; who gives the second signature; second level approval rights

Turbo treats a person as a second-stage approver for a purchase order only if they are in that PO's second-stage pool (secondStageLimits). The pool is empty when the PO does not need a second approval.

## Non-zero limit check covers any eligible approver

Lines 51-54

Asked as: do I have a limit; limit shows zero on PO; why is my limit zero; no approval limit on kind

Turbo says a person has a non-zero limit for the PO's expenditure kind when they are among the eligible approvers (eligibleLimits). This is true even if their limit is below the PO amount, and false for anyone not eligible.

## Eligible approver lookup

Lines 56-59

Asked as: is this person an approver; find approver by name; lookup PO approver; is user eligible to approve

Turbo looks up a person among the eligible approvers for a PO (eligibleApprovers). It returns their name details if they are eligible and nothing if they are not.

## Owner can bypass first stage by approving as themselves

Lines 61-74

Asked as: approve my own PO; I own the PO can I skip; self approve purchase order; owner skip first approval

On a PO that needs a second approval, the first stage counts as valid via bypass when the approver is the PO owner and the owner has a non-zero limit for the kind. Both ids are trimmed, and a blank id fails. A PO that needs no second approval never uses the bypass.

## Final approval needs a second-stage limit covering the amount

Lines 76-84

Asked as: final approval over my limit; cannot give final approval; amount too big for my limit; final sign off rejected

An approver has sufficient final limit only if they are in the second-stage pool and their limit is at least the PO amount. Anyone outside that pool fails the check.

## Second approval needed only above the kind's threshold

Lines 86-112

Asked as: when does PO need second approval; why no second signature needed; second approval threshold amount; amount equal to threshold

A purchase order needs a second approval when its expenditure kind has a second approval threshold (second_approval_threshold) above 0 and the PO amount is greater than that threshold. An amount equal to the threshold needs no second approval. A kind with no threshold set counts as 0, so it never needs a second approval.

## Unknown or blank expenditure kind is refused

Lines 90-127

Asked as: unknown expenditure kind error; PO has no expenditure kind; kind missing so no approvers; expenditure kind not found

Turbo cannot work out the approvers for a purchase order when the expenditure kind is blank or does not exist. It returns the error 'unknown expenditure kind' (ErrUnknownExpenditureKind), followed by 'blank' or the kind id.

## Division is required to find PO approvers

Lines 121-123

Asked as: division is required error; no division on PO; cannot list approvers blank division; PO division missing

Turbo cannot list approvers for a purchase order when the division is blank. It returns the error 'division is required'.

## Which people can approve a purchase order

Lines 140-172

Asked as: who can approve purchase orders; approver missing from list; approver list empty for PO; who is allowed to sign off POs

An eligible PO approver holds the PO approver claim, has an active admin profile (active = 1) and has a limit above 0 for the limit column that applies to the kind and to whether the PO has a job. Their divisions list must be empty, which means all divisions, or must include the PO's division. The list is sorted by surname, then given name.

## Every eligible approver is recorded before pool split

Lines 174-189

Asked as: shows eligible but not approving; low limit approver still listed; eligible but not in pool; approver skipped from first stage

Turbo first records every eligible approver with their limit for the PO's kind, and only then puts them in the first-stage or second-stage pool. An approver whose limit is too low for a single-stage PO is skipped from the first-stage pool but stays recorded as eligible. The policy also keeps the kind's second-approval threshold, whether a second approval is required (SecondApprovalRequired) and the limit column used.

## Single-stage PO: approver's limit must cover the amount

Lines 190-199

Asked as: approver limit must cover amount; equal limit enough to approve; single approval PO who can approve; low limit approver not shown

When a PO needs no second approval, only eligible approvers whose limit is at least the PO amount (an equal limit is enough) can approve it. They form the first-stage pool, and this approval is final. Approvers with a lower limit are not listed.

## Two-stage PO: approvers split by limit into two pools

Lines 201-212

Asked as: split approvers by limit; lower limit first stage; higher limit final approval; same person in both stages

When a PO needs a second approval, the pools do not overlap. Eligible approvers whose limit is below the amount can do the first-stage approval. Approvers whose limit is at least the amount can give the final, second-stage approval. An approver is in only one pool.

## Purchase order under revision definition

Lines 217-225

Asked as: PO under revision; PO can be edited again; unapproved PO with number; legacy manual entry revision

A purchase order is under revision when its status is Unapproved, it already has a PO number and it is not a legacy manual entry (legacy_manual_entry). A PO with no PO number, or a legacy manual entry, is not under revision.
