<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poApprovers.ts","sha256":"f129ffbd03c7e5352ac938d2a6217ce206575760580499524fc725c2c6ac0131","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poApprovers.ts

## Missing PO fields are sent as empty text in approver lookups

Lines 46-56

Asked as: empty approver lookup; blank PO fields approvers; missing division currency approvers; no default value PO

When Turbo looks up approvers for a purchase order, a missing division, currency, kind, type, start date (date), end date or frequency is sent as an empty string. Turbo does not substitute a default value for these fields.

## Missing PO total counts as zero when finding approvers

Lines 48-48

Asked as: PO total zero approvers; wrong approver amount; missing total approver lookup; PO amount approver wrong

When Turbo looks up approvers for a purchase order, it sends the total as a number converted to text. A missing total (total) is sent as "0". Turbo also converts a total typed as text into a number first. Which approvers a given amount returns depends on the server, not on this file.

## PO approver lookup flags whether the PO has a job

Lines 51-51

Asked as: PO has job flag; job missing approvers; PO without job approver; has_job true false

When Turbo looks up approvers for a purchase order, it sends has_job as "true" if the PO has a job that is not an empty string. It sends "false" if the job is empty or missing. How the server uses this flag is not shown in this file.
