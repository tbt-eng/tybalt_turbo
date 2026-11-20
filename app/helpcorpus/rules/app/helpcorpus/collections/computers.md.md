<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/computers.md","sha256":"8c3c989ea728ebcadb6b9e8610099fd96fe9692000e0192db04c5647390b2582","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/computers.md

## Who can list or view computer records

Lines 7-23

Asked as: who can see computers; computer list permission; view computer records access; why can't I see computers

A signed-in person holding the computers claim or the admin claim can list computer records and view one record through the API. Everyone else is refused. The list rule and the view rule use the same filter.

## Nobody can create computer records through the API

Lines 25-27

Asked as: add a computer via API; create computer record; API computer creation refused; new computer entry not allowed

The computers collection has no create rule, so nobody can create a record through the API. Only a PocketBase superuser can. Turbo's own server code can still create records.

## Nobody can update computer records through the API

Lines 29-31

Asked as: edit computer record API; change computer details not allowed; update computer via API; computer record won't save changes

The computers collection has no update rule, so nobody can update a record through the API. Only a PocketBase superuser can. Turbo's own server code can still update records.

## Nobody can delete computer records through the API

Lines 33-35

Asked as: remove computer via API; delete computer record; can't delete a computer; computer deletion refused

The computers collection has no delete rule, so nobody can delete a record through the API. Only a PocketBase superuser can. Turbo's own server code can still delete records.

## Computer slug must be unique

Lines 41-44

Asked as: duplicate computer slug; slug already taken; computer slug error; serial or manufacturer missing

Each computer record needs a slug (slug), and no two computers can have the same slug because of a unique index (idx_computers_slug). The serial (serial) and the manufacturer (mfg) are also required.

## Computer record id format

Lines 41-41

Asked as: computer id invalid; computer record id length; wrong characters in computer id; 15 character id rule

A computer record id is required. It has exactly 15 characters, and only lowercase letters a-z and digits 0-9 are allowed.
