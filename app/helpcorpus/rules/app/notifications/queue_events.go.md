<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/queue_events.go","sha256":"1a151f236c84a6de6f6cf8776de4a88e253b8a5a0a91ee2e531e135b289498ea","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/queue_events.go

## Expense rejection notice content and amount format

Lines 16-19

Asked as: expense rejected email; expense rejection amount format; why expense rejected notice shows currency; rejected expense notification details

The expense rejection notification names the employee, expense date, rejector and reason, and links to the expense details. The amount appears as the currency code followed by the total with 2 decimals, such as 'CAD 12.50'. The currency code comes from the expense's currency (utilities.EffectiveCurrencyCode).

## Time sheet rejection notice needs both profiles

Lines 31-49

Asked as: no rejection notice timesheet; timesheet rejection error; missing profile timesheet reject; timesheet rejected but no email

If Turbo cannot find the employee's profile or the rejector's profile, it sends no rejection notification for the time sheet and returns an error.

## Who is notified when a time sheet is rejected

Lines 51-77

Asked as: who gets timesheet rejected alert; timesheet rejection notify manager; my manager told about rejected timesheet; timesheet rejected who receives message

Turbo notifies the time sheet's owner, the person who rejected it, and the owner's manager (the manager on the owner's profile). The rejector is skipped when they are the owner. The manager is skipped when empty or the same person as the rejector or owner. Each person gets one notification (template timesheet_rejected), which links to the time sheet details.

## Time sheet rejection notice content

Lines 53-59

Asked as: timesheet rejected message contents; what does rejected timesheet email say; week ending rejection notice; timesheet rejection reason email

The rejection notification (timesheet_rejected) names the employee, the week ending, the rejector and the rejection reason. It links to the time sheet details page.

## Expense rejection notice needs both profiles

Lines 100-118

Asked as: expense rejection no notification; expense rejected error missing profile; rejected expense email not sent; expense reject profile error

If Turbo cannot find the employee's profile or the rejector's profile, it sends no rejection notification for the expense and returns an error.

## Who is notified when an expense is rejected

Lines 131-147

Asked as: who is told expense was rejected; expense rejected manager alert; rejected expense who gets notified; expense rejection recipients

Turbo notifies the expense's owner, the person who rejected it, and the owner's manager. The rejector is skipped when they are the owner. The manager is skipped when empty or the same person as the rejector or owner. Each person gets one notification (template expense_rejected).

## Uncommit request decision notifies the requester

Lines 162-191

Asked as: uncommit request answer notification; uncommit approved email; uncommit declined alert to requester; I asked to uncommit expense result

When someone decides an expense uncommit request, Turbo notifies the requester. If the requester made the decision themself, Turbo sends nothing. The notice is expense_uncommit_declined when the request status is 'declined'. For any other status it is expense_uncommit_approved.

## Uncommit notice fails without decider profile

Lines 169-172

Asked as: uncommit decision no email; uncommit notice error; uncommit decider profile missing; uncommit alert not sent

If Turbo cannot find the decider's profile, it sends no uncommit decision notification and returns an error.

## Uncommit notice outcome text and expense label

Lines 174-208

Asked as: uncommit notice wording; back in commit queue message; uncommitted and rejected text; uncommit notification expense label

The outcome reads 'The expense was uncommitted and is back in the commit queue.' by default. It reads 'The expense was uncommitted and rejected.' when the outcome is uncommit_and_reject. The expense label joins 'PO <number>' (if any), the vendor name (if any), and the currency code with total to 2 decimals. The notice also carries the decider's name, the reason, the decision comment, and a link to the expense details.

## Sharing a time sheet notifies each new viewer

Lines 217-260

Asked as: timesheet shared with me email; added as viewer timesheet alert; someone shared my timesheet notice; viewer added to timesheet notification

When viewers are added to a time sheet, Turbo sends each new viewer one notification (timesheet_shared). It names the sharer and the employee, gives the week ending, and links to the time sheet details. If no new viewers are added, Turbo sends nothing.

## Share notice needs sharer and employee profiles

Lines 225-243

Asked as: timesheet share email not sent; share timesheet error; sharing timesheet profile error; no share notification received

If Turbo cannot find the sharer's profile or the employee's profile, it sends no sharing notification and returns an error.

## Absorb proposal notice content

Lines 271-339

Asked as: absorb proposal email text; absorb notice client contact vendor; proposal notification lists records; absorb decision comment email

The notice labels the target as Client, Contact or Vendor followed by the target name, plus ' of <client name>' when the proposal has a client name. It lists the names of the proposed records, separated by commas. It gives the decider's name (empty when there is no decider), the decision comment, and links to decided absorb proposals.

## Absorb proposal outcome notifies the requester

Lines 283-300

Asked as: absorb proposal result email; absorb accepted declined message; no absorb update sent; absorb proposal outcome alert

Turbo notifies the requester of an absorb proposal when its status is accepted, declined or obsolete. It uses the template absorb_proposal_accepted, absorb_proposal_declined or absorb_proposal_obsolete. For any other status Turbo sends nothing. If the requester decided it themself, Turbo also sends nothing.

## Absorb proposal notice fails without decider profile

Lines 302-309

Asked as: absorb notice error no profile; absorb proposal decider missing; absorb notification failed error; absorb decision not notified

When the proposal has a decider and Turbo cannot find their profile, it sends no notification and returns an error.
