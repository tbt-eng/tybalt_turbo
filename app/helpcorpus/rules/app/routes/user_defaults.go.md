<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/user_defaults.go","sha256":"98eb4720b786b95e203a244dec3ff37b1ab088a2a2c5c261661ddd9d1ea5f79c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/user_defaults.go

## User defaults come from the signed-in user's own profile

Lines 20-37

Turbo returns the default division, default role, default branch, default expense payment type and the personal reimbursement setting for the signed-in user only. The user's profile supplies the division, role and payment type. The admin profile supplies the branch and the reimbursement setting.

## Missing profile gives empty defaults and no personal reimbursement

Lines 25-32

When a user has no profile or no admin profile, or a value is empty, Turbo returns an empty text for the default division, role, expense payment type (default_expense_payment_type) and branch (default_branch). It returns false for allow_personal_reimbursement, so personal reimbursement is not allowed by default.

## Loading user defaults can fail

Lines 38-40

When Turbo cannot load the defaults, it returns the error "failed to load user defaults".
