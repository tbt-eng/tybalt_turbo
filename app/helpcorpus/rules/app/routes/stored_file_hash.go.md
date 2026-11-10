<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/stored_file_hash.go","sha256":"69011f46355d22d075692109fefa8777692e86531850b4252fa29d8e2e16d7cd","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/stored_file_hash.go

## Only admins can audit or repair stored file hashes

Lines 94-106

Turbo refuses a stored-file hash audit or replace to a caller who is not signed in ('unauthorized'). It refuses a signed-in caller without the admin claim ('admin claim required'). Both the audit and the replace routes apply this check first.

## Hash audit returns the stored and calculated hashes

Lines 108-123

An admin's audit of a stored file returns the target record's details, the stored hash and the hash that Turbo calculates from the file. It also says whether the two match. The audit changes no data.

## Hash replace requires an updated value

Lines 130-148

The replace route trims the 'updated' value from the request body. It refuses a missing or blank value with 'updated is required'. It refuses an unreadable body with 'invalid request body'.

## Hash comparison ignores case and surrounding spaces

Lines 152-166

Turbo calculates a SHA-256 hash of the stored file. It trims the stored hash and compares the two without regard to upper or lower case. A stored hash in capitals or with spaces around it still matches.

## Replace checks that the target has not changed

Lines 181-191

Before it replaces a hash, Turbo reads the target again inside a transaction. If the collection, record, storage path or hash field differs from the audit, it refuses with the 'target changed' message.

## Replace refuses when the record was updated meanwhile

Lines 192-194

The replace request must carry the record's current 'updated' value. If the record's 'updated' value differs, Turbo refuses with the 'updated changed' message. It gives the same message if the update then affects no row.

## Replace does nothing when the hash already matches

Lines 195-201

If the stored hash already equals the calculated hash, ignoring case and spaces, Turbo changes nothing. It leaves the 'updated' value as it is and reports the result as a no-op (noop true, replaced false).

## Replace writes the calculated hash and a new updated time

Lines 203-239

When the hashes differ, Turbo writes the calculated hash into the record's hash field and sets 'updated' to the current time. The write happens only if 'updated' still equals the expected value. The result reports replaced true, the previous hash and the new hash.

## Replace refuses a hash that duplicates another record

Lines 214-217

If the new hash breaks a uniqueness rule on the collection, Turbo refuses the replace with the 'unique conflict' message.

## Hash calculation fails when the file is missing or unreadable

Lines 251-277

Turbo refuses to calculate a hash when the storage path is empty or only spaces, and gives the 'empty storage path' message. It gives the 'file not found' message when it cannot read the file. If the filesystem cannot be opened or the file cannot be hashed, it gives the matching message for that failure. The exact wording of these messages comes from the caller's message set.
