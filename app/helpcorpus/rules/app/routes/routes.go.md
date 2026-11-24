<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/routes.go","sha256":"b79358cdb69952c6923c57cb0e7470519b84630c3879c6ce1011ac58a77ae005","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/routes.go

## Expenses, purchase orders and vendors blocked when expense editing is off

Lines 14-31

Asked as: expense editing turned off; can't edit expenses; expenses locked; editing disabled error; can't change vendors

When the expenses editing feature flag is off, Turbo refuses actions on expenses, purchase orders and vendors with the message of ErrExpensesEditingDisabled. The check applies to these three collections only. Other collections skip it. In this file the check guards the purchase order approve, reject, cancel, close, make cumulative and revise routes, and the purchase order group preview and create routes.

## Time editing flag refuses actions when off

Lines 33-43

Asked as: time editing disabled; can't edit timesheet; time entries locked; time editing off; can't change hours

When the time editing feature flag is off, a check that uses requireTimeEditing refuses the action with the message of ErrTimeEditingDisabled. This file defines the check but binds it to no route itself.

## Reports require the report claim

Lines 45-54

Asked as: can't open report; report access denied; not authorized to view report; payroll report permission; can't see reports

A person without the report claim cannot open any report under /api/reports. Turbo refuses with the message "you are not authorized to view this report". This covers the payroll, weekly, receipts, payables spreadsheet, branch mismatch and active jobs reports.

## KPI reports require the kpi claim

Lines 56-65

Asked as: can't see KPI report; KPI access denied; KPI dashboard permission; KPI not authorized

A person without the kpi claim cannot open a KPI report under /api/kpi/reports. Turbo refuses with the message "you are not authorized to view KPI reports".

## Version endpoint is public

Lines 93-105

Asked as: app version check; which turbo version; build number endpoint; git commit version

The /api/version endpoint needs no sign-in. It returns the app name, version, build, full version, git commit, git branch and build time.

## Most Turbo API routes need a signed-in user

Lines 107-439

Asked as: api login required; unauthorized api call; 401 on timesheets; token needed for api; session expired api

The time sheet, expense, expense uncommit request, absorb proposal, time amendment, time entry, user, nav, admin profile, purchase order, client, vendor, job, WIP, report, stats, machine secret, attachment audit, currency, rate sheet, claim and work record routes all require a signed-in user. Rules for who may act on a given record are enforced in each handler, not here.

## Legacy export routes use machine secrets, not a user sign-in

Lines 133-136

Asked as: legacy export without login; export no sign in; machine secret export; export timesheets jobs expenses

The legacy export routes for time sheets, jobs and expenses do not require a user sign-in. They use their own authorization through machine_secrets.

## Radiator reports use machine secrets

Lines 138-139

Asked as: radiator login report; radiator report access; radiator secret; radiator report no login

The radiator login report routes use a machine secret with the radiator role instead of a user sign-in.

## Computer routes need the computers or admin claim

Lines 141-143

Asked as: can't see computers list; assign computer permission; retire computer access; computer admin access

The routes to list, view, assign, retire and un-retire computers are bound separately. Per the comment, they need the computers or admin claim.

## Absorb proposals open to any signed-in user, review by absorb holders

Lines 185-193

Asked as: propose an absorb; who approves absorb; withdraw absorb proposal; absorb review access

Per the comment, anyone who is signed in can propose an absorb. People with absorb access review proposals. Withdraw and decline have their own routes.

## Purchase order mutation routes require expense editing enabled

Lines 239-255

Asked as: can't approve purchase order; PO locked; can't cancel purchase order; reject PO disabled; purchase order frozen

Approve, reject, cancel, close, make cumulative and revise on a purchase order, and preview or create of a purchase order group, are refused when expense editing is off. Read routes for purchase orders do not have this gate.

## Legacy purchase order routes skip the expense editing gate

Lines 257-261

Asked as: old purchase order form; legacy PO create; edit old PO; PO edit page

The legacy purchase order routes (edit view, create, update) are not behind the expense editing flag check in this file. They need only a signed-in user.

## Absorb routes for clients, contacts, vendors and rate sheets

Lines 272-286

Asked as: absorb a client; undo absorb vendor; merge rate sheet; absorb client contact

Absorb and undo absorb exist for clients, client contacts, vendors and rate sheets. Per the comments, the absorb routes for rate sheets need admin, or absorb with rate_sheet_revise. Claim checks are in the handlers.

## Branch hours report checks access by branch

Lines 364-372

Asked as: branch hours report denied; can't open branch hours; branch access report; branch report permission

The /api/branch_hours report checks access with requireBranchHoursAccess. It uses the signed-in user and the branch query parameter. The rule itself is defined in another file.

## Machine secrets and rate sheet admin actions are claim-gated in handlers

Lines 379-439

Asked as: delete rate sheet; roll projects to latest rate; create rate sheet permission; machine secret admin only; claim details admin

Per the comments, machine secret management is admin only. Deleting an inactive unused rate sheet and rolling projects out to the latest rate sheet revision are admin actions. Creating rate sheets needs the job claim. Updating rate sheet entries and viewing claim details need the admin claim. The handlers enforce these, not this file.
