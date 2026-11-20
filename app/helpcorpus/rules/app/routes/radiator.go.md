<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/radiator.go","sha256":"0950e6bb433376015e697bb70c621bc0652a4099708a2ec76a21ccb85ee7edd9","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/radiator.go

## Radiator report body limited to 64 KB

Lines 38-43

Asked as: radiator report too big; file size limit radiator; request body too large; report upload size limit

Turbo caps a radiator report body at 64 KB (radiatorMaxBodySize, 64<<10). This replaces PocketBase's default 32 MB limit for this route. A larger body is refused as too large.

## Radiator report endpoint accepts only POST

Lines 45-78

Asked as: radiator endpoint GET request; method not allowed radiator; radiator wrong http method; can't PUT to radiator

The radiator endpoint (/api/radiator) answers only POST requests. Any other method (GET, PUT, PATCH, DELETE) is refused with 'method not allowed' and an Allow: POST header.

## Radiator report needs the machine secret

Lines 80-88

Asked as: radiator unauthorized error; radiator token rejected; machine token radiator login; radiator bearer token missing

A radiator report must carry an Authorization header 'Bearer <token>'. The token must be a valid machine token for the 'radiator' role. A missing header, an empty token or an invalid token is refused with 'unauthorized'. No user login is used.

## Radiator report must be sent as JSON

Lines 90-95

Asked as: radiator content type error; radiator report not JSON; Content-Type must be application/json; radiator wrong content type

The Content-Type of a radiator report must be application/json. Otherwise Turbo refuses it with 'Content-Type must be application/json'. The method, secret and Content-Type are checked before the body is read.

## Unreadable radiator request body is refused

Lines 102-108

Asked as: cannot read radiator body; radiator request body unreadable; radiator upload failed to read; could not read request body

When Turbo cannot read the body of a radiator report, it refuses it with "could not read the request body". A body larger than the limit gets Turbo's too-large-request refusal instead.

## Empty or malformed radiator submissions are refused

Lines 109-137

Asked as: empty radiator submission; radiator report rejected empty; radiator body not valid JSON; radiator report is a list

Turbo refuses an empty body, JSON null, an empty list or an empty object with 'empty submission, not saving'. It refuses text that is not valid JSON, or has extra content after the JSON, with 'the body is not valid JSON'. It refuses any other non-object, such as a non-empty list or a string, with 'the report must be a JSON object'. A leading UTF-8 byte order mark is ignored.

## Radiator report may start with a byte order mark

Lines 109-110

Asked as: radiator BOM error; radiator PowerShell script byte order mark; radiator script adds hidden characters; radiator report starts with odd bytes

Turbo removes a UTF-8 byte order mark at the start of a radiator report body before it reads the JSON. Old PowerShell scripts add this mark.

## Radiator report must be one JSON value only

Lines 116-124

Asked as: radiator JSON parse error; radiator report extra text after JSON; radiator body malformed; radiator invalid JSON report

A radiator report that is not valid JSON, or has extra content after the JSON value, is refused with "the body is not valid JSON".

## Radiator report that is not an object is refused

Lines 125-137

Asked as: radiator report not an object; radiator report is an array; radiator empty object submitted; radiator null submission

A JSON list or value that is not an object is refused with "the report must be a JSON object". A null, an empty list or an empty object is refused with "empty submission, not saving".

## Accepted radiator reports answer with no content

Lines 139-158

Asked as: radiator report success response; radiator no content reply; radiator accepted but empty response; radiator 204 status

Turbo answers a radiator report with an accepted status and no content, whether the report is valid or is stored as a rejected raw login. A server failure while storing gives 'failed to store the report'.

## Radiator report storage failure gives a server error

Lines 140-156

Asked as: radiator report server error; failed to store radiator report; radiator storage failed; radiator 500 error on save

If Turbo fails to store a valid report, or fails to save a raw login, it answers "failed to store the report". It logs the cause.

## Invalid report is saved as a raw login

Lines 151-156

Asked as: bad radiator report saved where; radiator rejected report stored; raw logins radiator invalid reports; radiator failed report kept

When a report fails validation or cannot be stored, Turbo saves the body as received in the raw logins (computer_raw_logins). It saves the reasons joined with '; ' and the computer name if the report has a valid one. The reasons are cut to 4000 characters.

## Valid report updates the computer and records a login

Lines 169-239

Asked as: radiator login recorded for computer; radiator computer not updated; radiator last login time; radiator report creates login record

A valid report finds the computer by its slug (serial and manufacturer) or creates it. It saves the computer, adds a login record and sets the last login time. All of this happens in one transaction.

## Radiator email and UPN are lowercased

Lines 179-184

Asked as: radiator email uppercase; radiator UPN case mismatch; radiator email lowercase; radiator upn letter case

Turbo lowercases the source anchor, UPN and email of a report before it stores or matches them. Email is empty when the report has none.

## Radiator new computer is created when slug is unknown

Lines 195-205

Asked as: radiator new laptop appears; radiator computer not showing up; radiator unknown computer added; radiator first report from machine

If no computer has the report's slug, Turbo creates one. Otherwise it updates the existing computer.

## Missing computer fields keep stored values

Lines 207-222

Asked as: radiator blank field overwritten; radiator missing field lost data; radiator keeps old values; radiator report partial update

A field missing from a report keeps its stored value on the computer. These fields are serial, mfg, model, computerName, systemType, ram, osVersion, osSku, osArch, the boot drive fields, radiatorVersion and networkConfig. The user fields (given name, surname, email, UPN, source anchor, last user) are always replaced by the latest report. The report never writes assigned_to, assigned_time or retired.

## Radiator user fields replace stored values

Lines 215-222

Asked as: radiator user changed on computer; radiator last user wrong; radiator user details not updating; radiator logged in user field

Each valid report replaces the computer's user given name, surname, email, UPN, source anchor and last user, and sets the last login time to now. It never writes assigned_to, assigned_time or retired. If no user matches, last user is empty.

## Valid report deletes earlier raw logins of that computer name

Lines 241-248

Asked as: radiator old raw logins remain; radiator duplicate raw login; radiator clear old login entries; radiator computer name raw logins

After a valid report with a computer name is stored, Turbo deletes all raw logins with that computer name. A report without a computer name deletes none.

## Report values refused by a field become raw logins

Lines 252-257

Asked as: radiator report saved as raw login; radiator value too long for field; radiator could not be saved; radiator field refused value

When the schema accepts a value but a computer field refuses it, such as a very long UPN, the report is saved as a raw login. The reason starts with 'the report could not be saved: '. It does not give a server error.

## Report user is matched by anchor, then email, then UPN

Lines 261-299

Asked as: radiator user not matched; radiator wrong user linked; radiator last user empty; radiator finds no user

Turbo finds the user of a report first by the admin profile with the same user source anchor. If none matches, it finds a user whose email equals the report's email, then the report's UPN, ignoring case. If nothing matches, the last user is left empty. Turbo never creates a user.

## Radiator report validation failures are all collected

Lines 303-312

Asked as: radiator all errors listed; radiator several problems in report; radiator error message long; radiator validation errors combined

Turbo checks every property of a report and lists all the problems, not only the first. It joins them with "; " in the raw login's error. It cuts the error text to 4000 characters (radiatorMaxErrorRunes).

## Raw login keeps the report and computer name

Lines 314-318

Asked as: raw login computer name empty; radiator rejected report name; raw login shows error reason; radiator bad report record

A rejected report is saved as received in a raw login (computer_raw_logins), with the error text. Its computer name (computer_name) comes from the report only when that name is valid. Otherwise it is empty.

## Radiator computer fields map to columns

Lines 321-339

Asked as: radiator field stored where; radiator computer columns mapping; radiator bootDrive saved to; radiator report field location

A report's computer fields (serial, mfg, model, computerName, systemType, ram, osVersion, osSku, osArch, bootDrive, bootDriveFS, bootDriveCap, bootDriveFree, radiatorVersion, networkConfig) are stored in the matching computer columns.

## Radiator pattern rules for report fields

Lines 366-391

Asked as: radiator must match the pattern; radiator wrong format anchor; radiator boot drive format error; radiator osVersion format

userSourceAnchor must be 32 hex characters. bootDrive must be a letter and a colon, such as C:. osVersion must be dotted numbers, such as 10.0.19045. upn must have text, an @ and more text on one line. A mismatch gives 'must match the pattern ...'.

## Radiator string field length limits

Lines 378-399

Asked as: radiator field too short; radiator text too long; radiator serial length error; radiator name characters limit

Text fields have limits. bootDriveFS 2-32 characters. serial 4-48. computerName 1-48. model up to 48. userGivenName 1-24. osArch 2-16. mfg 2-32. userSurname 2-35. A failing field is reported as 'must be at least N characters' or 'must be at most N characters'.

## Bad computer name or email is dropped, not rejected

Lines 381-394

Asked as: radiator email dropped; radiator computer name ignored; radiator invalid email stored blank; radiator bad computer name kept old

If computerName or email fails its checks, Turbo drops that property and the report stays valid. For email the check is that it is an email address. A dropped email is stored as empty. A dropped computerName means the report deletes no raw logins and the stored name stays.

## Required radiator report properties

Lines 401-442

Asked as: radiator is required error; radiator missing serial; radiator required fields list; radiator report incomplete

A report needs serial, upn, networkConfig, userSourceAnchor, radiatorVersion, mfg, userGivenName and userSurname. A missing one is reported as '<name>: is required'. Unknown properties are dropped.

## Radiator unknown properties are dropped

Lines 416-436

Asked as: radiator extra properties ignored; radiator unknown fields report; radiator extra data not saved; radiator unrecognised keys

Turbo ignores report properties that are not in its list. It does not store them.

## Radiator optional string lengths and bounds

Lines 456-462

Asked as: radiator optional field length; radiator minimum length zero; radiator string too long; radiator characters counted

A minimum length of 0 or less is not enforced. A string whose length is below the minimum or above the maximum fails with a message such as "must be at least N characters". Lengths count characters.

## Radiator number fields are coerced to whole numbers

Lines 471-594

Asked as: radiator must be an integer; radiator number field error; radiator decimal value rejected; radiator ram not a number

The fields systemType, osSku, bootDriveFree, radiatorVersion, ram and bootDriveCap must be whole numbers. Turbo accepts true as 1, false and null as 0, and numeric strings, including hex, octal and binary forms. Decimals, infinity and other text fail with '<name>: must be an integer'. An empty string fails.

## Network config must map MAC addresses

Lines 478-492

Asked as: radiator network config error; radiator MAC address format; radiator networkConfig must be object; radiator MAC address invalid

networkConfig must be an object with at least 1 property. Every property name must be a MAC address in the form aa:bb:cc:dd:ee:ff. Otherwise the report is invalid, with 'must have at least 1 property' or 'the property name ... must be a MAC address'. A value that is not an object fails with 'must be an object'.

## Radiator text fields coerce numbers and null

Lines 497-518

Asked as: radiator must be a string; radiator number in text field; radiator null text field; radiator object in field

A text field given a number or boolean is read as its JavaScript text. A null becomes an empty string. An object or list fails with '<name>: must be a string'. An empty string is then too short for fields with a minimum length.

## Radiator number zero and boolean coercion

Lines 523-555

Asked as: radiator true false as number; radiator fraction number rejected; radiator null number field; radiator zero for empty number

A boolean becomes 1 or 0 and null becomes 0 for number fields. A string with only white space becomes 0 in the string reader, but an empty string is refused as not an integer. A number with a fraction, or an infinite number, fails with "must be an integer".

## Radiator numeric strings accept hex, octal and binary

Lines 559-594

Asked as: radiator hex number in report; radiator 0x value rejected; radiator binary number string; radiator numbers with spaces

A number field given as text may be a decimal, or start with 0x, 0o or 0b for hex, octal or binary. Signs or underscores in those digits are refused. Surrounding white space is trimmed.

## Number sent for a text field is written as JavaScript would

Lines 596-608

Asked as: radiator serial number as text; radiator number turned into string; radiator 1e21 format; radiator serial shows exponent

When a radiator report sends a number for a text field such as the serial, Turbo turns it into text the way JavaScript does. Zero becomes "0". A number from 1e-6 up to, but not including, 1e21 is written in plain decimal form. A smaller or larger number is written in exponent form, such as 1e+21.

## Radiator white space includes non-breaking and Unicode spaces

Lines 610-618

Asked as: radiator serial spaces removed; radiator non-breaking space in name; radiator unicode whitespace report; radiator trim manufacturer spaces

Turbo treats tab, line feed, vertical tab, form feed, carriage return, space, no-break space, the U+2000 to U+200A spaces, U+1680, U+2028, U+2029, U+202F, U+205F, U+3000 and U+FEFF as white space. It uses this definition to trim numeric strings in number fields. It also uses it to strip white space from the serial and to trim the manufacturer when it builds the computer slug.

## Radiator computer slug from serial and manufacturer

Lines 620-643

Asked as: radiator computer key format; radiator slug how built; radiator serial manufacturer key; radiator computer id from serial

The computer key (slug) is the cleaned serial, a comma and the cleaned manufacturer. The serial loses white space, slashes and commas. The manufacturer is lowercased, loses slashes, periods, commas, 'inc' and 'ltd', is trimmed and has spaces turned into underscores. The cleaned serial must be at least 4 characters and the manufacturer at least 2, counted in UTF-16 units. Otherwise the report is invalid and saved as a raw login with 'serial ... or manufacturer ... too short'.

## Radiator slug needs minimum serial and manufacturer length

Lines 627-643

Asked as: radiator serial too short; radiator manufacturer too short; radiator slug length error; radiator short serial number rejected

A computer's key (slug) is the cleaned serial, a comma and the cleaned manufacturer. After cleaning, the serial must have at least 4 characters and the manufacturer at least 2. Otherwise the report is saved as a raw login with "serial ... or manufacturer ... too short".
