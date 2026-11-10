<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/po_approver_props.md","sha256":"d3a0840c8c43d7eec032b9f08e407e509c4f3175b86a83dad5be22b152134ad6","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/po_approver_props.md

## Nobody can list PO approver limits through the API

Lines 7-9

Asked as: list approver limits API; see all PO limits; who can read approver limits; fetch po_approver_props list

The list action on po_approver_props (purchase order approver limits) has no rule. Nobody can list these records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one PO approver limit record through the API

Lines 11-13

Asked as: view one approver limit; open single PO limit record; get approver limit by id; read one po_approver_props

The view action on po_approver_props has no rule. Nobody can view a single record through the API, except a PocketBase superuser. Turbo's own server code can still read it.

## Only admins can create or update PO approver limits

Lines 15-27

Asked as: who can edit PO limits; change approver limit amount; only admins update limits; create approver limit permission

A signed-in person with the admin claim can create and update po_approver_props records through the API. Everyone else is refused. The create rule and the update rule are the same filter.

## Nobody can delete PO approver limits through the API

Lines 29-31

Asked as: delete PO approver limit; remove approver limit record; limit disappears when user deleted; erase po_approver_props entry

The delete action on po_approver_props has no rule. Nobody can delete these records through the API, except a PocketBase superuser. Turbo's own code can still delete them. A record is also deleted automatically when its user claim (user_claim) is deleted (cascadeDelete).

## PO approver limit amounts must be whole numbers

Lines 37-44

Asked as: approver limit decimal amount; PO limit cents not allowed; max amount must be integer; too many divisions approver

max_amount is required and must be a whole number (onlyInt). The optional limits project_max, sponsorship_max, staff_and_social_max, media_and_event_max and computer_max must also be whole numbers. The divisions field can hold up to 999 divisions.

## Each user claim has at most one PO approver limit record

Lines 40-53

Asked as: duplicate approver limit for user; two PO limits same person; user already has approver limit; unique approver limit per user

A unique index on user_claim means two po_approver_props records cannot point to the same user claim. A user claim is required, and each record allows only one (maxSelect 1).
