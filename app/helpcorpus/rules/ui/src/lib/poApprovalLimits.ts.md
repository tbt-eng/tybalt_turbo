<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/poApprovalLimits.ts","sha256":"3f67a6f7ac05140d8b97d38f03cee5df27a39249af46a0a65e317e4991d0cfe0","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/poApprovalLimits.ts

## Who can view purchase order approvers' limits

Lines 1-3

Turbo lets a person view purchase order approvers only if they hold the report claim or the po_approver claim (canViewPOApprovers). Staff with neither claim cannot view them.
