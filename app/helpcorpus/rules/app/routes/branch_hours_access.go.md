<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/branch_hours_access.go","sha256":"e944ded3143fac8dead17c8b3e4985deaedf925366ce7fceeece105ba8f07c64","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/branch_hours_access.go

## KPI and admin claim holders view any branch's hours

Lines 17-25

Asked as: see hours for every branch; kpi claim view all branches; admin access branch hours; view other branches hours

A person who holds the kpi claim or the admin claim may view the branch hours report for any branch. Turbo does not check which branches they manage.

## Non-managers without kpi or admin cannot view branch hours

Lines 26-32

Asked as: branch hours report refused; need branch manager role; cannot open hours report; no access to branch hours

A person who holds neither the kpi nor the admin claim and who manages no branch is refused. Turbo shows: "You need to be a branch manager or hold the kpi or admin claim to view this report".

## Branch managers view only the branches they manage

Lines 33-36

Asked as: only my branch hours; other branch hours blocked; can only view my branches; hours of branch not mine

A branch manager without the kpi or admin claim may view the branch hours report only for a branch they manage. For any other branch Turbo shows: "You can only view the hours of branches you manage".

## Missing branch is not checked in the access check

Lines 33-37

Asked as: branch hours no branch selected; hours report blank branch; missing branch error hours; branch field empty hours

When no branch is given, a branch manager passes this access check. The handler reports the missing branch.
