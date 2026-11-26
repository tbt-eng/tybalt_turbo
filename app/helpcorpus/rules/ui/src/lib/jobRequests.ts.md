<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/jobRequests.ts","sha256":"e73caab2f42414fb18bf9898727ba24f30c1d0ca05a54761ce6c5f74cb7f7521","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/jobRequests.ts

## Next-step message for a draft project request

Lines 41-43

Asked as: request not submitted; draft says not submitted; creator hasn't submitted; why not submitted yet

A Draft request tells its creator "You have not submitted this request." Other viewers see "<creator> has not submitted this request." If the creator's name is blank, they see "Not submitted yet."

## Next-step message for a recalled project request

Lines 44-48

Asked as: recalled request; request was recalled; resubmit recalled request; recalled message wrong

A Recalled request tells its creator "Recalled. You have not submitted it again." Other viewers see the creator's name in that sentence. Without a name they see "Recalled. Not submitted again yet."

## Next-step message while awaiting approval

Lines 49-51

Asked as: waiting for approval; who is the approver; stuck waiting on approver; approver name missing

The selected approver sees "Waiting for your approval." Other viewers see "Waiting for approval by <approver>." If the approver's name is blank, they see "Waiting for approval."

## Next-step message when awaiting job creation

Lines 52-59

Asked as: approved but no job; who creates the job; waiting for job creation; approved request stuck

For an approved request in Awaiting creation, procurement who can create the job see "Approved. Waiting for you or another procurement user to create the job." Others see "Approved. Waiting for procurement to create the job."

## Linked proposal must be awarded before job creation

Lines 53-56

Asked as: proposal not awarded; can't create job yet; job blocked by proposal; award proposal first

If a linked proposal is not yet Awarded, the Awaiting creation message names that proposal. Procurement who can create are told "Award proposal <number>, then create the job." Others are told procurement must award the proposal and create the job.

## Next-step message for a rejected request

Lines 60-66

Asked as: request was rejected; why was it rejected; what to do rejected; rejected request next step

A Rejected request tells its creator "Rejected. Correct the request, then submit it again." Other viewers see "Rejected. Waiting for <creator> to correct it." Without a creator name they see "Rejected. Waiting for a correction." Any other status gives an empty message.

## Proposal statuses a request can link

Lines 68-72

Asked as: proposal won't link; can't link proposal; which proposal statuses allowed; proposal link refused

A project request can link a proposal only if its status is In Progress, Submitted or Awarded. The server refuses a proposal in any other status when the request is saved.

## Copied allocations skip inactive divisions

Lines 77-94

Asked as: allocations missing after copy; copied job lost divisions; inactive division removed; copied hours show zero

A request started from a proposal or a parent job copies that job's allocations (division and hours). Rows whose division is inactive are left out, because a save refuses an inactive division. Hours that are not a number become 0.

## Allocations rows without a division don't count

Lines 95-99

Asked as: allocation row ignored; blank division row; proposal not filling allocations; allocations list empty

A request has allocations only if at least one row has a division chosen. Rows with no division are ignored, so a proposal chosen later can still fill the list.

## Contact email display fallback

Lines 103-106

Asked as: contact shows no email; contact email missing; no email address shown; email field blank

A contact shows its email address. If it has none but is marked no email (no_email), it shows "No email address". Otherwise it shows nothing.

## Name display fallback

Lines 107-109

Asked as: name missing on record; blank name shows first last; contact name empty; name not displaying

A record shows its name. If it has no name, Turbo joins its given name and surname with a space.

## Billing address falls back to the client address

Lines 110-123

Asked as: billing address empty; which address for billing; address missing on invoice; billing uses client address

The billing address is the invoicing profile's address. If that is empty, Turbo uses the client's address. An address joins its non-empty street, city, province or state, postal code and country with commas.

## Whole-dollar CAD format

Lines 124-132

Asked as: project value not provided; value shows cents; dollar amount format; value displays wrong

A project value shows as Canadian dollars with no cents only if it is a whole number above 0. Zero, negative, fractional or non-number values show "Not provided".

## Print checklist: client, description and location

Lines 137-143

Asked as: plus code location error; location flagged on printout; print checklist location missing; invalid plus code location

A project request is missing Client name or Description if either is blank. It is missing "Valid Plus Code location" if the location is blank or is not a Plus Code of 8 characters, a plus sign, then 2 or 3 characters (letters 2-9, C, F, G, H, J, M, P, Q, R, V, W, X, any case).

## Print checklist: contacts

Lines 144-153

Asked as: contact flagged on printout; email invalid on checklist; billing contact check; printout contact problem

The project contact must have a name and an email. A billing contact is checked the same way, but only when the request has one. An invoicing profile is optional, so without a billing contact the printout leaves out the billing details. A contact marked no email needs no email address. Otherwise the email must look like name@domain.tld.

## Print checklist: project value

Lines 154-155

Asked as: project value flagged print; value missing from checklist; printout value error; positive value needed

A project request is missing "Positive whole-dollar project value" unless the project value is a whole number above 0.

## Print checklist: fee basis

Lines 156-157

Asked as: fee basis flagged; time and materials not set; print says fee basis; fee type not chosen

The request is missing "Fee basis" unless time and materials has been explicitly set to yes or no (time_and_materials_set is true and time_and_materials is a boolean).

## Print checklist: rate sheet readiness

Lines 158-165

Asked as: rate sheet flagged print; rate sheet not ready; missing roles on rate sheet; rate sheet inactive

With proposed rates, the rate sheet is ready when no role lacks a rate (rate_sheet_missing_roles is empty). With a chosen existing sheet, the sheet must be active. If no sheet is selected or it is not ready, the request is missing "Rate sheet".

## Print checklist: client billing rates

Lines 166-175

Asked as: client billing rates flagged; no rates on rate sheet; overtime rate missing; rate entries incomplete

For a time and materials request, "Client billing rates" is missing if the rate sheet has no entries or any entry lacks a finite rate or overtime rate. Fixed-fee requests skip this check.

## Who sees the job creation queue

Lines 178-185

Asked as: job creation queue hidden; can't see creation queue; no job queue tab; queue count missing

The job creation queue is visible only when requests are enabled, the user holds the job claim (procurement), and the server has sent a count for "/job-requests?queue=creation".
