<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"0851130dab80f6d88acd1289c430073e29cc19c7d1587d66f3b75f218d01c3bd","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Next-step message for a draft project request

Lines 41-43

Asked as: request still draft; says not submitted; never sent my request; why hasn't it gone out

For a Draft request, the creator sees "You have not submitted this request." Other viewers see "[creator] has not submitted this request." When the creator's name is missing, they see "Not submitted yet."

## Next-step message for a recalled project request

Lines 44-48

Asked as: I recalled my request; recalled request message; pulled back request; withdrawn request status

For a Recalled request, the creator sees "Recalled. You have not submitted it again." Other viewers see the creator's name in the message. When the name is missing, they see "Recalled. Not submitted again yet."

## Next-step message while awaiting approval

Lines 49-51

Asked as: who is approving my request; waiting on approver; stuck pending approval; approver name missing

For an Awaiting approval request, the selected approver sees "Waiting for your approval." Other viewers see "Waiting for approval by [approver]." When the approver's name is missing, they see "Waiting for approval." Turbo never puts a role in place of a missing name.

## Next-step message when awaiting job creation

Lines 52-59

Asked as: approved but no job; who creates the job; job not created yet; proposal needs awarding

For an Awaiting creation request, the message depends on whether procurement can create the job. If a linked proposal is not yet Awarded, the message names it: someone who can create is told to award it and then create the job, and others are told procurement must award it and create the job. With no unawarded proposal, a person who can create sees "Approved. Waiting for you or another procurement user to create the job." Others see "Approved. Waiting for procurement to create the job."

## Next-step message for a rejected request

Lines 60-64

Asked as: request was rejected; what to fix when rejected; sent back to me; why was it declined

For a Rejected request, the creator sees "Rejected. Correct the request, then submit it again." Other viewers see "Rejected. Waiting for [creator] to correct it." When the name is missing, they see "Rejected. Waiting for a correction."

## No next-step message for other statuses

Lines 65-66

Asked as: blank status line; no message under status; status text missing

For any status other than Draft, Recalled, Awaiting approval, Awaiting creation or Rejected, the next-step sentence is empty.

## Proposal statuses a request can link

Lines 68-72

Asked as: can't link proposal; proposal won't attach; link proposal error; wrong proposal status

A project request can link a proposal only when its status is In Progress, Submitted or Awarded. The server refuses other statuses when the request is saved.

## Copied allocations skip inactive divisions

Lines 77-94

Asked as: inactive division missing from copy; allocations missing when copying; copy from proposal lost rows; hours became zero

A request started from a proposal or a parent job copies that job's allocations (division and hours). Rows whose division is inactive are left out, because a save refuses an inactive division. Hours that are not a number become 0.

## Allocation rows without a division do not count

Lines 95-99

Asked as: empty allocation row; blank division row; allocation not counting; proposal not filling allocations

An allocation row with no division chosen does not count as an allocation. A proposal chosen later can still fill the list.

## New checklist starts with every answer false

Lines 100-102

Asked as: checklist all no; checklist won't continue; questions default to false; can't proceed checklist

A new checklist starts with each question answered false. Each question must be answered true to continue.

## Contact with no email shows a note

Lines 103-106

Asked as: no email address showing; contact missing email note; no_email flag; contact shows no email

A contact's email shows as its address. If the contact has no email and is marked no_email, Turbo shows "No email address". Otherwise it shows nothing.

## Billing address falls back to the client's address

Lines 110-123

Asked as: billing address blank; where billing address comes from; invoice address empty; address falls back to client

The billing address is the invoicing profile's address. When that is empty, Turbo uses the client's address. An address joins the non-empty parts of address, city, province/state, postal code and country with commas.

## Whole-dollar CAD display

Lines 124-132

Asked as: project value shows not provided; cents showing on value; CAD amount formatting; value displays wrong

A project value shows as Canadian dollars with no cents only when it is a positive whole number. Any other value (zero, negative, fractional or not a number) shows "Not provided".

## Print needs client name and description

Lines 137-138

Asked as: print check client name; print says description missing; print warning missing description; can't print client name

The print check lists "Client name" as missing when the client has no name, and "Description" when the request has none. Blank spaces alone count as empty.

## Print needs a valid Plus Code location

Lines 139-143

Asked as: plus code location error; print location invalid; location format plus code; print says location missing

The print check lists "Valid Plus Code location" as missing when the location is empty or does not match a Plus Code pattern: 8 characters, a +, then 2 or 3 characters from the Plus Code alphabet.

## Print needs contact names and emails

Lines 144-151

Asked as: print contact email error; billing contact missing name; print invalid contact email; contact without email print

For both the project contact and the billing contact, the print check lists a missing name. It also lists a missing email when the contact is not marked no_email and the email is not a valid address. A contact marked no_email passes without an email.

## Print needs a positive whole project value

Lines 152-153

Asked as: print project value error; value must be whole number; print value not positive; can't print project amount

The print check lists "Positive whole-dollar project value" as missing when the project value is not a whole number above 0.

## Print needs the fee basis set

Lines 154-155

Asked as: print fee basis missing; time and materials not set; fee basis error print; can't print fee

The print check lists "Fee basis" as missing unless time and materials has been set and is true or false.

## Print needs a ready rate sheet

Lines 156-163

Asked as: print rate sheet missing; rate sheet not chosen; inactive rate sheet; proposed rates missing role

The print check lists "Rate sheet" as missing when no rate sheet is chosen. Proposed rates are ready only when no role lacks a rate. A chosen existing sheet must be active.

## Time and materials needs client billing rates

Lines 164-173

Asked as: client billing rates missing; rate sheet no entries print; overtime rate missing; time materials rates error

When the fee basis is time and materials, the print check lists "Client billing rates" as missing if the rate sheet has no entries, or if any entry lacks a finite rate or overtime rate.

## Job creation queue visibility

Lines 176-183

Asked as: can't see creation queue; queue missing from menu; job holder queue hidden; creation queue not showing

A person sees the job creation queue only when requests are enabled, they hold the job claim, and the server has sent a count for the creation queue. The server adds that count only for an active job holder.
