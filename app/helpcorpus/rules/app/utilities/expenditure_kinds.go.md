<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/expenditure_kinds.go","sha256":"0aa4c07e9954dcafea1efef3bab7b042cb8919d77b92755ee97deb282350fdce","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/expenditure_kinds.go

## Expenditure kinds and their approver limit columns

Lines 11-40

Turbo has six expenditure kinds: capital, project, sponsorship, staff_and_social, media_and_event and computer. Each kind maps to one limit column on the approver's PO limits (po_approver_props). The columns are max_amount (capital), project_max, sponsorship_max, staff_and_social_max, media_and_event_max and computer_max.

## Default expenditure kind depends on whether a job is present

Lines 76-83

When a purchase order has a job, the default expenditure kind is project. When it has no job, the default is capital.

## Missing kind on legacy record gets a default kind

Lines 85-92

When a record has no expenditure kind (blank or only spaces), Turbo uses project if the record has a job and capital if it has none.

## Capital kind with a job becomes project

Lines 93-98

When a record has the capital kind and also a job, Turbo changes the kind to project. Capital with a job is not allowed. Other kinds with a job stay as they are.

## Which approver limit applies for a second approval

Lines 101-123

For second-approval qualification, Turbo picks the approver limit column that matches the PO's expenditure kind. For example, sponsorship uses sponsorship_max and computer uses computer_max.

## Second-approval limit falls back when kind is blank or unknown

Lines 104-118

When the PO has no expenditure kind, or the kind is not recognised, Turbo picks the limit by job. With a job it uses project_max. Without a job it uses max_amount (the capital limit).

## Known kind without a limit column is an error

Lines 119-122

When a PO's kind is recognised but has no limit column, such as the legacy standard kind, Turbo returns the error: unsupported expenditure kind id with the id and name. It does not fall back to a default.

## Startup check for expenditure kinds and limit columns

Lines 128-170

Turbo loads the expenditure kinds from the database and checks that all six exist. It then checks that po_approver_props has all six limit columns. If a kind or column is missing, the check fails with a message such as: missing expenditure_kind with name "x" or missing po_approver_props column "x".
