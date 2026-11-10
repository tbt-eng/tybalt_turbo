<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/user_po_approver_profile.md","sha256":"427e14553b63f6e2054c7af7822e85e2785637f954eeb5b81ed29370a079e88b","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/user_po_approver_profile.md

## PO approver profile is read-only through the API

Lines 5-5

Asked as: edit approver profile; change approver limits API; can't update profile record; delete approver profile

The view user_po_approver_profile has no create, update or delete rule in this file. Nobody can add, change or delete its records through the API, because Turbo only reads this view.

## Who can list PO approver profiles

Lines 7-11

Asked as: see other approvers' profiles; list all approver profiles; who are the PO approvers; everyone's approval limits

A signed-in person can list only the record whose id equals their own user id (`@request.auth.id = id`). Nobody sees another person's PO approver profile in the list through the API.

## Who can view a PO approver profile

Lines 13-17

Asked as: open colleague's approver profile; view someone else's limits; 403 on profile; check another person's approval limit

A signed-in person can view only the record whose id equals their own user id. A request for another person's profile is refused.

## Every user has a profile, with limits of 0 by default

Lines 22-55

Asked as: why are my limits zero; no approval limit showing; profile shows 0 everywhere; not a PO approver zero limits

The view has one record for every user, including users who are not PO approvers. A user with no PO approver claim, or with no approver limits set, gets 0 for max_amount, project_max, sponsorship_max, staff_and_social_max, media_and_event_max and computer_max. Their divisions default to an empty list ('[]').

## Limits come only from the PO approver claim

Lines 22-35

Asked as: where do approval limits come from; limits missing after claim added; approver limit not updating; divisions empty for approver

Turbo takes the limits and divisions from the po_approver_props record of a user's 'po_approver' claim. Users without that claim get the defaults.

## Profile lists all claims of the user

Lines 45-53

Asked as: which claims do I have; list of my user roles; claims field empty; show my permissions

The claims field lists the names of all claims the user holds, not only po_approver. A user with no claims gets an empty list ('[]').
