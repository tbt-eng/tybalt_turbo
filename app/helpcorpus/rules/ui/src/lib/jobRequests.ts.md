<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"475bcdc7a393b5518b8b8bfbf92824aac74badfba08702a4b8a55b75aacbb4c9","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Checklist starts with every answer false

Lines 14-16

Asked as: checklist all unanswered; checklist boxes unticked; why can't I continue; questions default to no; stuck on checklist

A new checklist for a request stage starts with each question unanswered (false). Each question must be answered true to continue.

## Approver sees reason for missing Approve button

Lines 19-21

Asked as: no approve button; approver can't approve; why approve button missing; approval blocked; approver limit too low

When the selected approver has no Approve button, for example because their limit is below the project value, the request detail gives a reason (approval_blocked_reason). Only that approver gets the reason.

## Next-step message for a draft project request

Lines 38-40

Asked as: draft request status message; not submitted yet; who has to submit; request still a draft; nobody submitted my request

For a Draft request, the creator sees "You have not submitted this request." Other viewers see "<creator> has not submitted this request." If the creator's name is missing, they see "Not submitted yet."

## Next-step message for a recalled project request

Lines 41-45

Asked as: recalled request what next; request was recalled; resubmit after recall; recalled not submitted again

For a Recalled request, the creator sees "Recalled. You have not submitted it again." Other viewers see the same sentence with the creator's name. If the name is missing, they see "Recalled. Not submitted again yet."

## Next-step message while awaiting approval

Lines 46-48

Asked as: who is approving my request; waiting for approval message; approver name missing; pending approval who; request stuck with approver

For an Awaiting approval request, the selected approver sees "Waiting for your approval." Other viewers see "Waiting for approval by <approver>." If the approver's name is missing, they see "Waiting for approval." The message never uses a role in place of a name.

## Next-step message while awaiting job creation

Lines 49-52

Asked as: waiting for job creation; who creates the job; approved but no job; procurement hasn't made job; job not created yet

For an Awaiting creation request, a viewer who can create jobs sees "Approved. Waiting for you or another procurement user to create the job." Other viewers see "Approved. Waiting for procurement to create the job." Only this status depends on a permission (canCreate), because procurement is a claim and not a person on the request.

## Next-step message for a rejected request

Lines 53-57

Asked as: rejected request what now; why was my request rejected; correct rejected request; waiting for correction; request sent back

For a Rejected request, the creator sees "Rejected. Correct the request, then submit it again." Other viewers see "Rejected. Waiting for <creator> to correct it." If the name is missing, they see "Rejected. Waiting for a correction."

## No next-step message for other statuses

Lines 59-59

Asked as: no next step shown; blank status message; why no message on request; empty next step

For any status that is not Draft, Recalled, Awaiting approval, Awaiting creation or Rejected, the next-step message is empty.

## Proposal statuses that accept a project request

Lines 61-65

Asked as: can't link proposal; proposal not available; which proposals can link; proposal status blocks save; link proposal error

A project request can link a proposal only when the proposal's status is In Progress, Submitted or Awarded. The server refuses the other statuses when the request is saved.

## Contact with no email shows a note

Lines 69-72

Asked as: contact has no email; no email address note; contact without email text; email missing on contact

If a contact has no email address but is marked as having none (no_email), Turbo shows "No email address". If neither is set, it shows nothing.

## Billing address falls back to the client's address

Lines 76-89

Asked as: billing address wrong; billing address missing; why client address shown; invoice address blank; address fallback

The billing address is the invoicing profile's address when it has any address parts. Otherwise Turbo uses the client's address. An address joins the address, city, province/state, postal code and country that are filled in, separated by commas.

## Project value display shows only positive whole dollars

Lines 90-98

Asked as: project value shows not provided; value showing cents; project value formatting; dollar amount display wrong; value not whole dollars

Turbo shows a project value as whole Canadian dollars with no cents. A value that is not a positive whole number shows "Not provided".

## Print is missing client name or description

Lines 103-104

Asked as: print missing client name; print warning description; can't print client; printout incomplete description

The missing-items list for printing includes "Client name" when the client has no name. It includes "Description" when the request's description is blank.

## Location must be a valid Plus Code

Lines 105-109

Asked as: plus code invalid; location format wrong; print location error; open location code; location won't print

The missing-items list includes "Valid Plus Code location" when the location is blank. It also includes it when the location is not 8 characters, a plus sign, then 2 or 3 characters from the Plus Code alphabet (23456789CFGHJMPQRVWX, any case).

## Contact name and email are needed for printing

Lines 110-117

Asked as: contact missing for print; print contact name email; billing contact incomplete; contact email required print

For both the project contact and the billing contact, the list includes "<label> name" when the name is blank. A name is the contact's name, or its given name and surname. It includes "<label> email" when the contact has no valid email address and is not marked as having no email.

## Project value must be a positive whole number

Lines 118-119

Asked as: project value error print; value not whole number; zero value print problem; value must be positive

The missing-items list includes "Positive whole-dollar project value" when the project value is not a whole number above zero.

## Fee basis must be set explicitly

Lines 120-121

Asked as: fee basis missing; time and materials not set; fixed or hourly unset; print fee warning

The missing-items list includes "Fee basis" unless time and materials has been explicitly set to a true or false value.

## Rate sheet readiness for printing

Lines 122-129

Asked as: rate sheet missing print; no rates for roles; rate sheet not active; proposed rates inactive; print rate warning

The missing-items list includes "Rate sheet" when the request has no rate sheet. For proposed rates it also includes it unless every role has a rate. For a chosen existing sheet it also includes it unless the sheet is active. Proposed rates stay inactive until the job is created.

## Time and materials needs client billing rates

Lines 130-139

Asked as: no billing rates time materials; overtime rate missing; client rates empty; t&m rates incomplete

For a time and materials request, the missing-items list includes "Client billing rates" when there are no rate sheet entries. It also includes it when any entry lacks a numeric rate or overtime rate. Other requests skip this check.

## Job creation queue visibility

Lines 142-149

Asked as: no job creation queue; queue not showing; can't see creation queue; job queue hidden; no count for queue

A person sees the job creation queue only when job requests are enabled, they hold the job claim, and the server has sent a count for the creation queue. The server sends that count only to an active job holder.
