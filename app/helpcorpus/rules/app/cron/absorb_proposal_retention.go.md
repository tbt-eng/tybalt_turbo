<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/cron/absorb_proposal_retention.go","sha256":"ef30b65349ad1c82be51f59704d8b1351d05e38fb1d949384211101788ceaf37","model":"claude-sonnet-5-5"} -->

# Rules in app/cron/absorb_proposal_retention.go

## Decided absorb proposals are deleted after the retention period

Lines 13-17

Asked as: old absorb proposals deleted; how long proposals are kept; absorb proposal retention period; decided proposals removed automatically

A daily job deletes each absorb proposal that is no longer pending (accepted, declined, withdrawn or obsolete) when its decision date (decided_at) is older than ABSORB_PROPOSAL_RETENTION_DAYS days before now. The number of days is set in the constants file, not here. The cutoff is computed in UTC.

## Pending absorb proposals are never deleted

Lines 19-32

Asked as: pending proposal never deleted; old pending absorb proposal; undecided proposal kept forever; proposal without decision date

The retention job keeps every pending absorb proposal, however old it is. It also keeps any proposal with an empty decision date (decided_at).

## Closed proposals are not kept as history

Lines 25-32

Asked as: closed proposal history; where is absorb history kept; proposal exactly at cutoff; declined proposal disappears

A closed absorb proposal is not kept as history. The record of a committed absorb is its absorb_history row. A proposal decided exactly at the cutoff stays, because the comparison is strictly older than the cutoff.

## Deleting a proposal also deletes its records

Lines 30-49

Asked as: proposal records deleted too; absorb proposal records removed; deleting proposal cascades; partial deletion proposal

When the retention job deletes an expired absorb proposal, it first deletes that proposal's absorb_proposal_records rows. Both deletions run in one transaction, so a failure leaves nothing partly deleted.

## Retention job runs daily and logs the result

Lines 56-66

Asked as: daily cleanup job proposals; retention job log; how many proposals deleted; proposal cleanup failed error

The daily cron job (runAbsorbProposalRetention) deletes the expired proposals using the current time. It logs how many it deleted. If the deletion fails, it logs an error and deletes nothing.
