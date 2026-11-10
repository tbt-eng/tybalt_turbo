<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/client_notes.md","sha256":"89b5a4cc19ff88247aed87a0cc35c89510666dc415a4bff2ccea42cb2d279cc4","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/client_notes.md

## Who can list or view client notes

Lines 7-17

Asked as: who can see client notes; view client notes list; notes visible to whom; read a client note

Any signed-in person can list client notes and view one client note. A request without a signed-in user is refused.

## Who can create a client note

Lines 19-23

Asked as: add a client note; write a note on client; create note for client; post a comment about client

Any signed-in person can create a client note (the create rule only checks that the request is signed in).

## Nobody can edit a client note through the API

Lines 25-27

Asked as: cannot edit my client note; change note after saving; fix typo in client note; update note through api

The update rule is empty, so nobody can update a client note through the API, except a PocketBase superuser. Turbo's own server code can still update one.

## Only admins can delete a client note

Lines 29-34

Asked as: delete client note; remove a note; who can delete notes; can't delete my note

Only a signed-in person who holds the admin claim (user_claims_via_uid with cid.name 'admin') can delete a client note through the API. Other signed-in people cannot.

## Client note text length limit

Lines 41-41

Asked as: note too short; note text too long; how many characters in note; character limit client note

A client note's text (note) is required and must be 10 to 1000 characters long.

## Client note must name a client and a user

Lines 42-45

Asked as: note without client; client note needs user; link note to job; note requires client and user

Every client note requires one client (client) and one user (uid). The job (job) is optional and can hold at most one job.

## Status a note can record for a job

Lines 46-46

Asked as: job status on note; note marked cancelled; no bid option note; job changed to cancelled

The field job_status_changed_to accepts only one value: "Cancelled" or "No Bid". It can also be left empty.
