<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/machine_secrets.md","sha256":"70960c078190bd490ae9bbc33ca7786ea760e37c6159511fcbab4fb52dde83cb","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/machine_secrets.md

## Nobody can list machine secrets through the API

Lines 7-9

Asked as: list machine secrets; see all secrets; fetch every secret via API; why can't I list secrets

The machine_secrets collection has no list rule. Nobody can list its records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view a machine secret through the API

Lines 11-13

Asked as: view a machine secret; open one secret record; read secret through API; can't see secret details

The machine_secrets collection has no view rule. Nobody can view one record through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create machine secrets through the API

Lines 15-17

Asked as: create a new machine secret; add secret via API; make a new secret record; POST secret refused

The machine_secrets collection has no create rule. Nobody can create a record through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update machine secrets through the API

Lines 19-21

Asked as: edit machine secret; change secret record; update secret via API; secret edit refused

The machine_secrets collection has no update rule. Nobody can update a record through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Only admins can delete machine secrets

Lines 23-27

Asked as: delete a secret; remove machine secret; who can delete secrets; delete secret permission denied

A person can delete a machine_secrets record through the API only if they hold the admin claim (a user_claims record whose claim name is 'admin'). Everyone else is refused.

## Machine secret hash and salt must match set patterns

Lines 34-35

Asked as: secret hash invalid; sha256 format error; salt rejected; 64 character hash

The sha256_hash field is required and must be exactly 64 hexadecimal characters. The salt field is required and must be 16 to 128 letters or digits.

## Machine secret role has two allowed values

Lines 36-37

Asked as: secret role value; which role for secret; expiry date missing; radiator or legacy role

The role field is required and takes one value: 'legacy_writeback' or 'radiator'. The expiry date is also required.
