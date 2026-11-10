<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/queue_events.go","sha256":"ecea7a18343c86159b39058897f6c8cd11080ebfdf22a198b117ed554bff28b5","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/queue_events.go

## Expense rejection notice content and amount format

Lines 15-18

The expense rejection notification includes the owner's name, the expense date, the rejector's name, the rejection reason and a link to the expense details page. The amount is shown as the currency code followed by the total with 2 decimals, such as "CAD 12.50". The currency code comes from utilities.EffectiveCurrencyCode, which this file does not define.

## Time sheet rejection notice needs both profiles

Lines 30-48

If Turbo cannot find the profile of the time sheet owner or of the rejector, it sends no rejection notifications and returns an error ("error finding employee profile" or "error finding rejector profile").

## Who is notified when a time sheet is rejected

Lines 50-76

When a time sheet is rejected, Turbo notifies the time sheet's owner. It also notifies the person who rejected it, unless that person is the owner. It also notifies the owner's manager, unless the manager is empty or is the owner or the rejector. Each person receives one notification (template timesheet_rejected).

## Time sheet rejection notice content

Lines 52-58

The time sheet rejection notification includes the owner's name, the week ending, the rejector's name and the rejection reason. It also includes a link to the time sheet details page.

## Expense rejection notice needs both profiles

Lines 99-117

If Turbo cannot find the profile of the expense owner or of the rejector, it sends no rejection notifications and returns an error ("error finding employee profile" or "error finding rejector profile").

## Who is notified when an expense is rejected

Lines 119-146

When an expense is rejected, Turbo notifies the expense's owner. It also notifies the rejector, unless that person is the owner. It also notifies the owner's manager, unless the manager is empty or is the owner or the rejector (template expense_rejected).

## No notice when requester decides own uncommit request

Lines 161-166

When the person who requested an expense uncommit is also the person who decided it, Turbo sends no notification. A direct uncommit has no one else to tell.

## Uncommit decision notice goes to the requester

Lines 168-207

When someone else decides an expense uncommit request, Turbo notifies only the requester. The notice names the decider and gives the reason, the decider's comment and a link to the expense details page. If the decider's profile is not found, Turbo returns the error "error finding decider profile" and sends nothing.

## Uncommit notice expense label

Lines 173-180

The expense label in the uncommit notice lists the PO number as "PO <number>" if there is one. It then lists the vendor name if there is one. It always ends with the currency code and the total with 2 decimals. The parts are joined with commas.

## Uncommit notice outcome wording

Lines 182-185

If the request's outcome is uncommit_and_reject, the notice says "The expense was uncommitted and rejected." For any other outcome, it says "The expense was uncommitted and is back in the commit queue." The wording does not depend on whether the request was declined.

## Uncommit notice template depends on decision

Lines 187-190

If the uncommit request's status is declined, Turbo uses the template expense_uncommit_declined. For any other status, it uses expense_uncommit_approved.

## Time sheet sharing notifies each new viewer

Lines 216-259

When a time sheet is shared, Turbo sends one notification to each newly added viewer (template timesheet_shared). It says who shared it, whose time sheet it is, the week ending and a link to the details page. If there are no new viewers, Turbo sends nothing.

## Time sheet sharing notice needs sharer and owner profiles

Lines 224-242

If Turbo cannot find the profile of the sharer or of the time sheet owner, it sends no sharing notifications and returns an error ("error finding sharer profile" or "error finding employee profile").
