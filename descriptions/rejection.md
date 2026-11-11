# The Time Sheet Rejection System

## The RejectModal component

- implemented in ui/src/lib/components/ directory as RejectModal.svelte
- The RejectModal component was modelled on the ShareModal component
  which is found in ui/src/lib/components/ShareModal.svelte (now built on
  DSModal; RejectModal still has its own overlay). The RejectModal is
  used to reject a timesheet by prompting the user to enter a rejection reason.
  Once the user enters the reason and clicks Reject, the modal calls the reject
  route with the timesheet ID in the URL and the rejection reason in the body.
  There is also a cancel button that will close the modal without making any
  changes. The modal takes a collectionName prop, and the expense and purchase
  order pages use it too (POST /api/{collectionName}/{id}/reject).

## New API route

- registered in app/routes/routes.go; the handler is in
  app/routes/reject_record.go
- POST /api/time_sheets/{id}/reject
  - This route is used to reject a timesheet. The same handler rejects an
    expense at POST /api/expenses/{id}/reject.
  - The ID of the timesheet to reject is in the URL. The route takes a JSON
    body with the following structure:
    - rejection_reason (string): The reason for rejecting the timesheet. It
      must be at least 4 characters long.
  - The requestor must be the approver of the time_sheet record or hold the
    commit claim. Anyone else gets an error. A commit claim holder who is not
    the approver can reject only a timesheet that is already approved. The
    approver can reject a submitted timesheet before or after approval.
  - If the timesheet is already committed, it throws an error
  - If the timesheet is not submitted it throws an error
  - If the timesheet is already rejected, it throws an error
  - Everything is performed in a transaction. See the approve route
    (app/routes/approve_record.go) for details of this. If the checks pass, the
    time_sheets record is updated as follows: The rejected field is set to the
    current time The rejection_reason is set as specified in the JSON body The
    rejector column is set to the id of the user calling the function.
  - After the transaction, the route queues a rejection notification.

## updates to ui/src/routes/time/sheets/list/+page.svelte

- an openRejectModal function in this file behaves similarly to the approve
  function but instead of calling the rejection api directly it opens the
  RejectModal. The RejectModal is imported like the ShareModal.
- The Reject button opens the reject modal with the appropriate id. The button
  shows only for a submitted timesheet that is not committed and not rejected,
  and only to its approver, or to a commit claim holder when the timesheet is
  approved.

## other

If any other changes are needed to implement this system, create or update files
as needed. Add comments in the code as you see fit to make it clear what is
happening so it is maintainable and legible later
