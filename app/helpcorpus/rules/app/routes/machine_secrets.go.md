<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/machine_secrets.go","sha256":"d4c585d92caf89384ccb39d89892576d283d8dc34cda6f4cae33c80a42d2416a","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/machine_secrets.go

## Machine secret role must be a valid role

Lines 18-21

The only valid role for a machine secret is legacy_writeback. Any other role is refused with "invalid role".

## Only signed-in admins can create machine secrets

Lines 43-55

Creating a machine secret requires a signed-in user. Without a sign-in, Turbo refuses with "unauthorized". A signed-in user without the admin claim gets "admin claim required".

## Machine secret days must be a positive number

Lines 63-66

When an admin creates a machine secret, the number of days (days) must be greater than 0. Zero or a negative number is refused with "days must be a positive integer".

## Machine secret role check on creation

Lines 68-71

When an admin creates a machine secret, Turbo checks the role against the valid roles. A role that is not valid is refused with "invalid role".

## Machine secret is generated and stored only as a salted hash

Lines 73-104

Turbo generates a random 32-character secret and a random 21-character alphanumeric salt. It stores the SHA-256 hash of the salt followed by the secret, plus the salt, and does not store the secret itself.

## Machine secret expiry is days from now in UTC

Lines 91-116

The expiry of a new machine secret is the current UTC time plus the requested number of days. Turbo saves it as a date and time. The response gives the expiry in RFC3339 format.

## Secret is shown only once at creation

Lines 110-116

The response that creates a machine secret returns the plain secret together with the id, expiry and role. The list of secrets never shows the secret, so an admin cannot see it again.

## Only signed-in admins can list machine secrets

Lines 122-134

Listing machine secrets requires a signed-in user. Without a sign-in, Turbo refuses with "unauthorized". A signed-in user without the admin claim gets "admin claim required".

## Machine secret list shows all secrets, newest first, without hash or salt

Lines 136-158

The list returns every machine secret, with no limit, sorted by created date with the newest first. Each item shows only the id, role and expiry, never the hash or salt. Expired secrets are not filtered out.
