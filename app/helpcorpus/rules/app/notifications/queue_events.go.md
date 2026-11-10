<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/queue_events.go","sha256":"ecea7a18343c86159b39058897f6c8cd11080ebfdf22a198b117ed554bff28b5","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/queue_events.go

## Expense rejection notice content and amount format

Lines 15-18

Asked as: expense rejected email content; rejection notice expense amount format; CAD 12.50 currency display; expense rejection reason link

The expense rejection notification includes the owner's name, the expense date, the rejector's name, the rejection reason and a link to the expense details page. The amount is shown as the currency code followed by the total with 2 decimals, such as "CAD 12.50". The currency code comes from utilities.EffectiveCurrencyCode, which this file does not define.

## Time sheet rejection notice needs both profiles

Lines 30-48

Asked as: timesheet rejection no notification; cannot find employee profile rejection; error finding rejector profile; timesheet rejected error no emails

If Turbo cannot find the profile of the time sheet owner or of the rejector, it sends no rejection notifications and returns an error ("error finding employee profile" or "error finding rejector profile").

## Who is notified when a time sheet is rejected

Lines 50-76

Asked as: who gets notified timesheet rejected; timesheet rejection manager notified; rejected timesheet notify owner manager; timesheet rejection duplicate notifications

When a time sheet is rejected, Turbo notifies the time sheet's owner. It also notifies the person who rejected it, unless that person is the owner. It also notifies the owner's manager, unless the manager is empty or is the owner or the rejector. Each person receives one notification (template timesheet_rejected).

## Time sheet rejection notice content

Lines 52-58

Asked as: timesheet rejection email contents; week ending rejected timesheet notice; timesheet rejected message details; rejection notice timesheet link

The time sheet rejection notification includes the owner's name, the week ending, the rejector's name and the rejection reason. It also includes a link to the time sheet details page.

## Expense rejection notice needs both profiles

Lines 99-117

Asked as: expense rejected missing profile; error finding employee profile expense; expense rejection no emails sent; expense rejector profile not found

If Turbo cannot find the profile of the expense owner or of the rejector, it sends no rejection notifications and returns an error ("error finding employee profile" or "error finding rejector profile").

## Who is notified when an expense is rejected

Lines 119-146

Asked as: who gets told expense rejected; expense rejection manager copied; rejected expense notify manager; expense rejected notification recipients

When an expense is rejected, Turbo notifies the expense's owner. It also notifies the rejector, unless that person is the owner. It also notifies the owner's manager, unless the manager is empty or is the owner or the rejector (template expense_rejected).

## No notice when requester decides own uncommit request

Lines 161-166

Asked as: uncommit own expense no email; self uncommit no notification; I uncommitted my own expense; no alert when I uncommit

When the person who requested an expense uncommit is also the person who decided it, Turbo sends no notification. A direct uncommit has no one else to tell.

## Uncommit decision notice goes to the requester

Lines 168-207

Asked as: uncommit request decision email; who hears about uncommit decision; uncommit declined notice requester; decider comment uncommit message

When someone else decides an expense uncommit request, Turbo notifies only the requester. The notice names the decider and gives the reason, the decider's comment and a link to the expense details page. If the decider's profile is not found, Turbo returns the error "error finding decider profile" and sends nothing.

## Uncommit notice expense label

Lines 173-180

Asked as: uncommit email expense label; PO number vendor in uncommit notice; uncommit notice expense name format; expense description missing in uncommit email

The expense label in the uncommit notice lists the PO number as "PO <number>" if there is one. It then lists the vendor name if there is one. It always ends with the currency code and the total with 2 decimals. The parts are joined with commas.

## Uncommit notice outcome wording

Lines 182-185

Asked as: uncommit notice wrong wording; uncommitted and rejected message; back in commit queue text; uncommit message says rejected when declined

If the request's outcome is uncommit_and_reject, the notice says "The expense was uncommitted and rejected." For any other outcome, it says "The expense was uncommitted and is back in the commit queue." The wording does not depend on whether the request was declined.

## Uncommit notice template depends on decision

Lines 187-190

Asked as: uncommit declined email template; uncommit approved email template; which email sent uncommit request; uncommit notice template choice

If the uncommit request's status is declined, Turbo uses the template expense_uncommit_declined. For any other status, it uses expense_uncommit_approved.

## Time sheet sharing notifies each new viewer

Lines 216-259

Asked as: timesheet shared no email; new viewer gets notice timesheet; sharing timesheet notification sent; timesheet share alert recipients

When a time sheet is shared, Turbo sends one notification to each newly added viewer (template timesheet_shared). It says who shared it, whose time sheet it is, the week ending and a link to the details page. If there are no new viewers, Turbo sends nothing.

## Time sheet sharing notice needs sharer and owner profiles

Lines 224-242

Asked as: share timesheet error sharer profile; timesheet sharing no notifications sent; error finding sharer profile; share notification missing profile error

If Turbo cannot find the profile of the sharer or of the time sheet owner, it sends no sharing notifications and returns an error ("error finding sharer profile" or "error finding employee profile").
