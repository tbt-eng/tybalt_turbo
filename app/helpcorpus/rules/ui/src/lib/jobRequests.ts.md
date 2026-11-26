<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"3a3954da3e017f9a41299a9b507e26c380a152fe450208dc79457933b20e6a3b","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Next-step message for a draft project request

Lines 41-43

Asked as: request saved but not sent; why does it say not submitted; draft request status; I never sent my request

For a Draft request, the creator (the viewer who made it) sees "You have not submitted this request." Other viewers see "{creator} has not submitted this request." If the creator's name is missing, they see "Not submitted yet."

## Next-step message for a recalled project request

Lines 44-48

Asked as: request was pulled back; recalled request message; resubmit withdrawn request; what happens after recall

For a Recalled request, the creator sees "Recalled. You have not submitted it again." Other viewers see the creator's name in the message. If the name is missing, they see "Recalled. Not submitted again yet."

## Next-step message while awaiting approval

Lines 49-51

Asked as: who is approving my request; pending approval message; request stuck with approver; waiting on approver text

For an Awaiting approval request, the selected approver sees "Waiting for your approval." Other viewers see "Waiting for approval by {approver}." If the approver's name is missing, they see "Waiting for approval."

## Next-step message when approved and awaiting job creation

Lines 52-59

Asked as: approved but no job; who creates the job; proposal not awarded yet; job not created after approval

For an Awaiting creation request, the message depends on whether the viewer can create jobs (a procurement user). If a linked proposal is not Awarded yet, it names the proposal. A procurement user is told to award the proposal and then create the job. Others are told that procurement must award it and create the job. With no unawarded proposal, a procurement user sees "Approved. Waiting for you or another procurement user to create the job." Others see "Approved. Waiting for procurement to create the job."

## Next-step message for a rejected request

Lines 60-64

Asked as: why was my request rejected; rejected request what now; how to fix rejected request; rejection correction message

For a Rejected request, the creator sees "Rejected. Correct the request, then submit it again." Other viewers see "Rejected. Waiting for {creator} to correct it." If the creator's name is missing, they see "Rejected. Waiting for a correction."

## No next-step message for other statuses

Lines 65-66

Asked as: blank next step under request; no message under request; why is next step empty; request shows no guidance

For any status other than Draft, Recalled, Awaiting approval, Awaiting creation or Rejected, the next-step message is empty. A missing name is left out and a role is never put in its place.

## Statuses of a proposal that a request can link

Lines 68-72

Asked as: cannot link proposal to request; proposal link refused; which proposals can I attach; proposal status blocks link

A project request can link a proposal only when the proposal's status is In Progress, Submitted or Awarded. The server refuses the other statuses when it saves the request.

## Contact email display falls back to a no-email note

Lines 76-79

Asked as: contact shows no email; no email address label; contact missing email address; email field empty on contact

A contact shows its email address. With no email but the no-email flag (no_email) set, it shows "No email address". Otherwise it shows nothing.

## Name display falls back to given name and surname

Lines 80-82

Asked as: name blank on record; name shows first and last; given name and surname display; client name missing

A record shows its name. With no name, it shows the given name and surname joined by a space, skipping any that is empty.

## Billing address falls back to the client's address

Lines 83-96

Asked as: billing address missing; invoice address empty; use client address for billing; address shows stray commas

The billing address is the invoicing profile's address. If that is empty, Turbo uses the client's address. An address joins address, city, province or state, postal code and country with commas and skips empty parts.

## Project value shows as whole Canadian dollars

Lines 97-105

Asked as: project value shows not provided; value displayed in dollars; why project value not provided; project value cents shown

A project value shows as whole Canadian dollars (CAD) with no cents only if it is a positive whole number. Zero, negative, fractional or non-number values show "Not provided".

## Print checklist: client, description and location

Lines 110-116

Asked as: cannot print request; print checklist missing items; plus code location invalid; request not ready to print

A request is not ready to print without a client name, a description, and a valid Plus Code location. The missing items are listed as "Client name", "Description" and "Valid Plus Code location". A Plus Code needs 8 characters, a plus sign, then 2 or 3 characters, from the Plus Code alphabet 23456789CFGHJMPQRVWX. Case does not matter.

## Print checklist: project and billing contacts

Lines 117-124

Asked as: project contact missing for print; billing contact email needed; contact not ready to print; contact name missing print

The project contact and the billing contact each need a name. Each also needs a valid-looking email address, unless the contact is marked as having no email. The missing items are listed as "Project contact name", "Billing contact email" and so on.

## Print checklist: project value and fee basis

Lines 125-128

Asked as: fee basis missing print; project value must be whole dollars; time and materials not set; value blocks printing

The project value must be a positive whole number, or "Positive whole-dollar project value" is listed as missing. The fee basis must be explicitly set to time and materials or not, or "Fee basis" is listed as missing.

## Print checklist: rate sheet readiness

Lines 129-136

Asked as: rate sheet missing print; proposed rate sheet role no rate; rate sheet not active; rate sheet blocks print

A request needs a rate sheet. If the rate sheet is a proposed one, it is ready only when no role lacks a rate. If a sheet was chosen, it must be active. Otherwise "Rate sheet" is listed as missing.

## Print checklist: client billing rates for time and materials

Lines 137-146

Asked as: time and materials rates missing; rate sheet entries empty; overtime rate missing; client billing rates not set

For a time and materials request, "Client billing rates" is listed as missing when there are no rate sheet entries. It is also listed when any entry lacks a numeric rate or overtime rate.

## Job creation queue visibility

Lines 149-156

Asked as: cannot see creation queue; job queue missing from menu; no job creation tab; job queue hidden

The job creation queue is visible only when project requests are enabled, the person holds the job claim, and the server has sent a count for the creation queue.
