<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/machine_secrets.go","sha256":"33d3e7241407adb32c5f42d2f375c297e2caba74a5bd4f76da868936efc1ba23","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/machine_secrets.go

## Machine secret role must be legacy_writeback or radiator

Lines 18-22

Asked as: invalid role; wrong machine secret role; role not allowed; which roles can a secret have; legacy_writeback or radiator

A machine secret can have only the role legacy_writeback or radiator. Any other role gives "invalid role".

## Only admins can create a machine secret

Lines 44-56

Asked as: not an admin; admin claim required; who can create machine secret; no permission to make secret; unauthorized secret creation

Creating a machine secret requires a signed-in user who holds the admin claim. A user without the admin claim gets "admin claim required". A request with no signed-in user gets "unauthorized".

## Machine secret days must be a positive number

Lines 64-67

Asked as: days must be positive; zero days secret; negative expiry days; secret days number wrong; days must be a positive integer

Creating a machine secret refuses a number of days (days) of 0 or less. Turbo shows "days must be a positive integer".

## Machine secret role check on creation

Lines 69-72

Asked as: role rejected on create; secret creation role error; wrong role when making secret; can't create secret with this role; invalid role on new secret

Creating a machine secret refuses any role other than legacy_writeback or radiator. Turbo shows "invalid role". The days check runs before the role check.

## New machine secret is generated and stored only as a hash

Lines 74-109

Asked as: how secrets are stored; is the secret saved in plain text; secret hashed in database; where is the secret kept; salt and hash for secret

Turbo generates a random 32-character secret and a random 21-character alphanumeric salt. It stores the SHA-256 hash of the salt followed by the secret, with the salt, the role and the expiry. It does not store the secret itself.

## Machine secret expires after the given number of days

Lines 92-93

Asked as: when does secret expire; secret expiry time; how long secret stays valid; expiry in RFC3339; secret lifetime in days

Turbo sets the expiry of a new machine secret to the current UTC time plus the requested number of days. It returns the expiry in RFC3339 format.

## Secret is shown once at creation

Lines 111-117

Asked as: secret only shown once; lost my secret; can't see secret again; copy secret after creating; secret missing from list

The response to a successful creation returns the id, the unhashed secret, the expiry and the role, with a created result. The secret appears only in this response, and the list never shows it.

## Only admins can list machine secrets

Lines 123-135

Asked as: who can see secrets; can't view machine secrets; list secrets not allowed; admin claim needed to view; view secrets unauthorized

Listing machine secrets requires a signed-in user who holds the admin claim. A user without it gets "admin claim required". A request with no signed-in user gets "unauthorized".

## Machine secret list shows id, role and expiry, newest first

Lines 137-159

Asked as: list of machine secrets; secrets sorted by date; which fields show in secret list; secret list newest first; does list show the secret value

The list returns all machine secrets, with no limit, sorted by creation date with the newest first. Each item shows only the id, the role and the expiry. It never shows the hash, the salt or the secret.
