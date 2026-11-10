<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/user_claims.md","sha256":"7a9e79718db1cc58d9a5b668a6dea395d0d381d5e14aca95c156ac8116f4b28e","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/user_claims.md

## Who can list user claim records

Lines 9-15

Asked as: see my claims list; list user claims; why can't I see other people's claims; view all claim records

A person can list a user claim record (user_claims) if it is their own record (uid is their ID). An admin (the admin claim) can list every user claim record. Anyone else sees none.

## Who can view one user claim record

Lines 19-25

Asked as: open one claim record; see a single user claim; who can view a claim; view someone else's claim

A person can view a user claim record (user_claims) if it is their own (uid is their ID). An admin (the admin claim) can view any user claim record. Anyone else cannot.

## Only admins can create user claim records

Lines 29-32

Asked as: give a claim to someone; add a claim to a user; grant permission to a person; create user claim via api

Only an admin (the admin claim) can create a user claim record (user_claims) through the API. This is how a claim is given to a person. Turbo's own server code can still create these records.

## Only admins can update user claim records

Lines 36-39

Asked as: change my own claims; edit a user claim; modify claim record; update claim permissions

Only an admin (the admin claim) can update a user claim record (user_claims) through the API. Everyone else is refused, including the person whose record it is.

## Only admins can delete user claim records

Lines 43-46

Asked as: remove a claim from someone; take away permission; delete user claim; revoke a person's claim

Only an admin (the admin claim) can delete a user claim record (user_claims) through the API, which takes a claim away from a person. Everyone else is refused, including the person whose record it is.

## A person cannot hold the same claim twice

Lines 52-63

Asked as: duplicate claim error; user already has this claim; same claim added twice; unique index claim conflict

A unique index on the person (uid) and the claim (cid) means a user claim record (user_claims) cannot repeat the same pair. Each person can hold a given claim only once. Both fields are required, and each points to one person and one claim.
