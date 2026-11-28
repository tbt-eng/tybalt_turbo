<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/radiator.go","sha256":"0e763d3a84142fc86c198df3e5cc255c48a799fe25b2ebf4ed9e99d07213e97d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/radiator.go

## Radiator checks run before the body is read

Lines 3-11

Asked as: refused before body read; early refusal radiator request; body ignored wrong header; size check before reading; radiator request rejected early

Turbo checks the method, the machine secret and the Content-Type of a radiator report before it reads any of the body. A request that fails one of these checks is refused without its body being read. The body size limit applies after these checks.

## Radiator secret must belong to the radiator machine role

Lines 39-40

Asked as: radiator token wrong role; unauthorized radiator token; bearer token not accepted; machine token rejected radiator; radiator login refused unauthorized

The radiator route takes no user sign-in. The bearer token must be a valid machine token for the machine role "radiator" (radiatorMachineRole), as checked by ValidateMachineToken. Turbo answers "unauthorized" otherwise.

## Radiator body is capped at 64 KB

Lines 41-41

Asked as: radiator report too big; 64 kb limit radiator; request body too large; radiator size limit; large content-length refused

The radiator body may not exceed 64 KB (radiatorMaxBodySize, 64<<10). Turbo refuses a larger Content-Length at once, and stops reading a body without a length once it passes the limit. This replaces PocketBase's default 32 MB limit for this route.

## Radiator endpoint accepts only POST

Lines 45-78

Asked as: get radiator endpoint; method not allowed radiator; radiator rejects put; allow post header; wrong http verb radiator

The radiator login endpoint (POST /api/radiator) answers only POST. GET, PUT, PATCH and DELETE are refused with 'method not allowed' and an Allow: POST header. The check runs before any of the body is read.

## Radiator requires a valid machine secret

Lines 80-88

Asked as: radiator password missing; missing authorization header; blank radiator token; invalid radiator secret; radiator key not working

A radiator report must carry an Authorization header with a Bearer token that is valid for the machine role 'radiator'. Turbo trims white space around the token. A missing header, a blank token or an invalid token is refused with 'unauthorized'. No user sign-in is used.

## Radiator requires JSON content type

Lines 90-93

Asked as: content-type must be json; radiator wrong content type; send radiator report as json; charset content type error; json header missing radiator

The request Content-Type must be application/json, and parameters such as charset are allowed. Otherwise Turbo refuses with 'Content-Type must be application/json'. This check runs after the secret check.

## Radiator body read errors

Lines 100-108

Asked as: could not read request body; radiator upload failed; body read error radiator; report body unreadable; too large body error

If the body is larger than the 64 KB limit, Turbo refuses it as too large. If the body cannot be read for another reason, Turbo answers "could not read the request body".

## Radiator rejects empty or malformed bodies

Lines 109-137

Asked as: empty submission not saving; body not valid json; report must be json object; radiator byte order mark; blank or null radiator body

Turbo refuses with 'empty submission, not saving' a body that is blank, null, an empty list or an empty object. It refuses with 'the body is not valid JSON' a body that does not parse or has extra content after the JSON. It refuses any other non-object, such as a non-empty list, with 'the report must be a JSON object'. A leading UTF-8 byte order mark is ignored.

## Rejected radiator reports are saved as raw logins

Lines 139-158

Asked as: report accepted but not saved; raw login reasons; failed report still stored; why accepted but missing; validation errors stored radiator

When a report fails validation or cannot be stored, Turbo still answers 'accepted' and keeps the report as a raw login with the reasons joined by '; '. The reasons are cut to 4000 characters. A server failure while storing returns 'failed to store the report'.

## Radiator accepted report returns no content

Lines 151-158

Asked as: accepted with no body; empty radiator response; cannot tell if saved; radiator reply no content; success response empty radiator

After a report passes the guard checks and is processed, Turbo answers "accepted" with no body. It does this for a stored login and for a report saved as a raw login. A person cannot tell the two apart from the answer.

## Valid report updates computer, adds login, clears raw logins

Lines 169-243

Asked as: report saves computer and login; old raw logins deleted; new login record created; valid radiator report saved; computer updated from report

A valid report is saved in one transaction. Turbo creates or updates the computer, adds a login record, and deletes the raw logins with the same computer name.

## Missing report fields keep stored computer values

Lines 207-225

Asked as: missing field erased computer; fields not sent keep old; report changed assigned to; radiator keeps stored values; retired flag changed by report

A computer field missing from a report keeps its stored value. The user fields (given name, surname, email, UPN, source anchor, last user) and the last login time are always replaced. A report never changes assigned_to, assigned_time or retired.

## Unsavable field values make a raw login

Lines 245-250

Asked as: upn too long error; field value too long; could not be saved raw login; valid value still not saved; storage rejected field value

If a field refuses a value that passed validation, such as a very long UPN, Turbo stores the report as a raw login with 'the report could not be saved: ...' and not as a server error.

## Last user is matched by anchor, then email, then UPN

Lines 254-292

Asked as: last user not found; which user last logged in; last user matched by email; user lookup by upn; last user empty radiator

Turbo finds the last user first from the admin profile with the same user source anchor. If none matches, it looks for a user whose email equals the report's email, then its UPN, ignoring case. If nothing matches the last user is empty. Turbo never creates a user.

## Raw logins are replaced per computer name

Lines 294-338

Asked as: duplicate raw logins computer; blank computer name raw logins; raw logins cleanup 30 days; old raw logins piling up; raw login name case matters

A new raw login deletes older raw logins with the exact same computer name (case and spaces count), so each computer keeps one. A blank name matches nothing, so those raw logins stay until the retention job deletes them after 30 days. A valid report also deletes raw logins of its computer name.

## Report fields copied to the computer record

Lines 340-358

Asked as: computer record not updated; which fields update computer; ram and os version stored; radiator fields copied computer; hardware info not updating

A valid report sets these computer fields: serial, mfg, model, computerName, systemType, ram, osVersion, osSku, osArch, bootDrive, bootDriveFS, bootDriveCap, bootDriveFree, radiatorVersion and networkConfig. A field that is present in the report replaces the stored value.

## Radiator report validation lists every problem and ignores unknown fields

Lines 360-363

Asked as: all validation errors at once; unknown fields ignored report; legacy schema validation; report errors listed together; extra report fields dropped

Turbo checks a radiator report against the schema of the legacy system. It collects all the problems and does not stop at the first. It drops any property that the schema does not list. It converts values to the expected type where it can.

## Radiator UPN may not contain line breaks

Lines 385-386

Asked as: upn newline rejected; upn with line break; carriage return in upn; upn must have at sign; upn line separator error

A radiator report's UPN (upn) needs an @ with at least one character before it and one after it. These characters cannot be line breaks: \n, \r, or the Unicode line and paragraph separators. A UPN with a line break fails the pattern check, and the report becomes a raw login.

## Radiator network config keys must be MAC addresses

Lines 389-389

Asked as: network config key format; mac address keys invalid; networkconfig names colons; mac key rejected report; hex pairs network key

In a radiator report, each property name in the network config (networkConfig) must be a MAC address. That is six pairs of hex digits separated by colons, such as 00:1A:2B:3C:4D:5E. Upper and lower case both pass. A name in any other form makes the report invalid.

## Email must look like an email address

Lines 390-391

Asked as: invalid email stored blank; email dropped from report; bad email address computer; email format not accepted; computer email empty

A report's email must match a simple email-address pattern (the ajv 6 "fast" format). An email that does not match is dropped from the report and does not cause a rejection. Turbo then stores the computer's email as blank.

## Radiator field length and pattern limits

Lines 397-418

Asked as: serial too short; field max length radiator; osversion format wrong; bootdrive format error; osarch length limit

Text limits: serial 4-48 characters, mfg 2-32, model up to 48, bootDriveFS 2-32, osArch 2-16, userGivenName 1-24, userSurname 2-35. userSourceAnchor must be 32 hex characters, bootDrive a letter and colon such as C:, and osVersion like 10.0.19045. upn must have an @ with text on both sides.

## Bad computerName or email is dropped, not rejected

Lines 400-400

Asked as: computer name too long dropped; bad computername ignored; report valid but name blank; computername field removed; invalid name still accepted

If computerName (1-48 characters) or email fails its check, Turbo drops that field and the report can still be valid. Other failing fields make the report invalid.

## Required report fields

Lines 420-461

Asked as: required field missing report; is required error radiator; missing serial or upn; report needs mfg field; required radiator fields list

A report must have serial, upn, networkConfig, userSourceAnchor, radiatorVersion, mfg, userGivenName and userSurname. A missing one gives the reason '<name>: is required'. Unknown properties are dropped.

## Field error messages in rejected reports

Lines 462-496

Asked as: must be at least characters; must match pattern error; must be a string message; radiator error message wording; every failing field reported

A rejected report's reasons use fixed messages. They read "name: must be at least N characters", "name: must be at most N characters", "name: must match the pattern ...", "name: must be an email address", "name: must be a string", and "name: must be an integer". Turbo reports every failing field, not just the first. Unknown properties are dropped. A minimum length of 0 is never checked.

## Network config must be MAC addresses

Lines 497-511

Asked as: networkconfig empty rejected; network config must be object; network config needs property; mac names network config error; networkconfig rejected report

networkConfig must be an object with at least 1 property, and each property name must be a MAC address such as 00:1A:2B:3C:4D:5E. Otherwise the report is rejected.

## Radiator values are coerced to the expected type

Lines 516-613

Asked as: number sent as text; null becomes empty string; boolean in text field; radiator value type conversion; numeric string integer field

Numbers and booleans in text fields become text, and null becomes an empty string. Integer fields (systemType, osSku, bootDriveFree, radiatorVersion, ram, bootDriveCap) accept whole numbers, numeric strings including hex, true as 1, false as 0 and null as 0. Fractions, infinity and an empty string are refused with '<name>: must be an integer'.

## Whole numbers required for integer fields

Lines 539-574

Asked as: must be an integer error; fraction ram rejected; blank integer field; decimal in ram field; whole number required radiator

The fields systemType, osSku, bootDriveFree, radiatorVersion, ram and bootDriveCap must be whole numbers. A value with a fraction, an infinite number or text that is not a number fails with "must be an integer". A blank string also fails. A boolean, null or numeric text is converted to a number.

## Numbers in text fields are converted like JavaScript

Lines 615-637

Asked as: number to text 1e21; zero becomes 0 string; non-breaking space trimmed; javascript number formatting; unicode spaces in text

A number sent for a text field becomes text in JavaScript style. Zero becomes "0". Numbers from 1e-6 up to below 1e21 are written in plain form, and others in exponent form such as 1e+21. Whitespace is judged by JavaScript's rules, including the non-breaking space and the Unicode spaces from U+2000 to U+200A.

## Computer is identified by serial and manufacturer

Lines 639-662

Asked as: how computer is matched; serial and manufacturer slug; manufacturer inc removed; computer identity lookup; slashes removed from serial

Turbo finds a computer by a slug built from serial and manufacturer (makeSlug). The serial loses white space, slashes and commas. The manufacturer is lowercased and loses slashes, periods, commas, 'inc' and 'ltd'. It is then trimmed and spaces become underscores. The slug is serial, a comma, then manufacturer.

## Slug needs serial of 4 and manufacturer of 2 characters

Lines 658-661

Asked as: serial too short raw login; manufacturer too short error; short serial number stored; serial or manufacturer too short; computer match too short

After cleaning, the serial must have at least 4 characters and the manufacturer at least 2, counted in UTF-16 units. Otherwise the report is stored as a raw login with the error 'serial ... or manufacturer ... too short'.
