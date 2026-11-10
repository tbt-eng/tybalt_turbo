<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/machine_secrets.md","sha256":"9965ed1bf7adcb335c3618849cb2d0b3b334ba0a9035693cec6ea6b04a5e847a","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/machine_secrets.md

## Nobody can list machine secrets through the API

Lines 7-9

Asked as: list machine secrets; see all secrets api; who can read secrets list; secrets endpoint empty

The list rule for machine_secrets has no filter. Nobody can list these records through the API except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view one machine secret through the API

Lines 11-13

Asked as: view one machine secret; open a secret record; get secret by id; secret api access denied

The view rule for machine_secrets has no filter. Nobody can view a single record through the API except a PocketBase superuser. Turbo's own server code can still read it.

## Nobody can create machine secrets through the API

Lines 15-17

Asked as: create machine secret api; add new secret via api; post secret fails; can't make a secret

The create rule for machine_secrets has no filter. Nobody can create a record through the API except a PocketBase superuser. Turbo's own server code can still create one.

## Nobody can update machine secrets through the API

Lines 19-21

Asked as: edit machine secret; update secret api; change secret value; secret update forbidden

The update rule for machine_secrets has no filter. Nobody can update a record through the API except a PocketBase superuser. Turbo's own server code can still update one.

## Only admins can delete machine secrets through the API

Lines 23-27

Asked as: delete machine secret; remove a secret; who can delete secrets; secret delete not allowed

A person passes the delete rule for machine_secrets only when one of their user claims (user_claims_via_uid) has the claim name 'admin'. Other signed-in people and anonymous requests cannot delete these records.

## Machine secret record id format

Lines 33-33

Asked as: secret id wrong length; invalid secret record id; id rejected lowercase; secret id characters

The id of a machine secret is required. It is exactly 15 characters and matches ^[a-z0-9]+$, which is lowercase letters and digits only.

## Machine secret hash must be 64 hex characters

Lines 34-34

Asked as: secret hash invalid; sha256 not accepted; hash must be hex; secret hash length error

The sha256_hash field of a machine secret is required. It must match ^[a-fA-F0-9]{64}$, which is exactly 64 hexadecimal characters, upper or lower case.

## Machine secret salt must be 16 to 128 letters or digits

Lines 35-35

Asked as: salt rejected; secret salt too short; salt characters not allowed; salt length error

The salt field of a machine secret is required. It must match ^[a-zA-Z0-9]{16,128}$, which is 16 to 128 letters or digits and nothing else.

## Machine secret role allows only legacy_writeback

Lines 36-36

Asked as: secret role wrong; which role for secret; role value not accepted; legacy writeback role

The role field of a machine secret is required and takes one value. The only allowed value is 'legacy_writeback'.

## Machine secret requires an expiry date

Lines 37-37

Asked as: secret won't save; missing expiry on secret; secret expiry required; expiration date blank error

The expiry field of a machine secret is required. A record cannot be saved without an expiry date.
