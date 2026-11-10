<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/oauth_onboarding.go","sha256":"fb2a0eafe987d1300219f9ccf1eb6fd13346e1e83f314a457056a524cb87f37a","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/oauth_onboarding.go

## Microsoft display name falls back to account name

Lines 87-94

Asked as: wrong name after login; name from Outlook; display name not showing; account name vs Microsoft name

Turbo takes the Microsoft name from the account's name (user.Name) first. It uses the Microsoft displayName only when the account name is empty. Turbo trims spaces from both.

## Microsoft mail and principal name come from the raw profile

Lines 97-113

Asked as: missing details from Microsoft; blank email from Microsoft; spaces in Microsoft profile; Microsoft profile field empty

Turbo reads the Microsoft given name, surname, mail and user principal name from the raw Microsoft profile. A value that is missing or not text becomes empty. Turbo trims spaces from each value.

## Microsoft name taken from given name and surname first

Lines 115-121

Asked as: first and last name; name split wrong; surname missing from Microsoft; full name from Microsoft

Turbo builds a person's full name from the Microsoft given name and surname when both exist. If either is missing, it uses the Microsoft display name instead. The display name is the name Microsoft supplies, or the Microsoft displayName value when that is empty.

## Microsoft email falls back to user principal name

Lines 123-125

Asked as: no email from Microsoft; login email blank; UPN used as email; email missing at sign in

Turbo uses the Microsoft mail address as the person's email. If mail is empty, it uses the Microsoft user principal name (userPrincipalName).

## Relink falls back to username from email local part

Lines 127-129

Asked as: link account by username; old account matched by email; reconnect old username; account found from email prefix

If no safe email match exists, Turbo takes the part of the Microsoft email or user principal name before the @ and normalises it into a username. It looks for an existing user with that username. If there is no username or no such user, no relink happens.

## Relink first matches an existing user by Microsoft email

Lines 150-166

Asked as: existing account not linked; link my old account; same email different login; old user not connected to Microsoft

When Microsoft presents a provider id Turbo does not recognise, Turbo first looks for a user whose app email equals the Microsoft mail or user principal name. It accepts that user only if the user has no Microsoft link yet, or is already linked to the same provider id.

## Relink search stops when no email or username is usable

Lines 150-178

Asked as: cannot relink account; no account to link; relink finds nothing; login without linking

If the email lookup finds no safe user and the Microsoft mail and principal name give no usable username, Turbo finds no relink candidate. The login then continues without a relink.

## Username match requires matching name

Lines 174-256

Asked as: username matches but name differs; account not linked same username; username found but not used; name must match to link

Turbo reuses a user found by username only if the Microsoft full name equals the user's name, ignoring case. Otherwise it needs both a Microsoft given name and surname that equal the profile's given name and surname, ignoring case. Without split names, or without a profile, the user is not matched and no relink happens.

## Relink username lookup finds no user

Lines 180-188

Asked as: no user with that username; relink username not found; username lookup fails login

If no user has the derived username, Turbo finds no relink candidate. If the database fails for another reason, the login fails.

## Username match with different linked Microsoft id needs identity check

Lines 200-210

Asked as: account linked to another Microsoft; new Microsoft login old user; different Microsoft id on account

A user found by username is safe to relink if it has no Microsoft link or the same provider id. If it is linked to a different provider id, Turbo relinks only when the name identity check passes.

## Identity matches a user by full name, ignoring case

Lines 218-221

Asked as: same name Microsoft login; name matches Microsoft account; login matched by full name

Turbo treats a Microsoft identity as the same person when the user's name equals the Microsoft full name, ignoring case and outer spaces.

## Profile name check needs both given name and surname

Lines 223-226

Asked as: first last name not matching; surname mismatch at login; given name check failing

If the full name does not match, Turbo checks the profile's given and family names. It runs that check only when Microsoft sent both a given name and a surname. Without both, the identity does not match.

## Identity does not match when the user has no profile

Lines 228-255

Asked as: user has no profile; profile names missing login; identity check profile fail

If the user has no profile record, the identity does not match. A match needs both the profile given_name and surname to equal Microsoft's, ignoring case.

## Existing Microsoft link is looked up per user

Lines 258-275

Asked as: where Microsoft link is stored; Microsoft link not found error; look up Microsoft connection

Turbo looks for a user's existing Microsoft link by the users collection, the user and the Microsoft provider. If none exists, that is not an error.

## Email match is rejected if another Microsoft id is linked

Lines 277-307

Asked as: email belongs to other Microsoft; wrong account by email; email linked to different login

If the user found by email is already linked to a different Microsoft provider id, Turbo does not reuse that user by email. It then tries the username match instead.

## Username relink is safe when no Microsoft link or same id

Lines 309-330

Asked as: account already linked same Microsoft; relink old username account; missing user record relink

When Turbo considers relinking a user found by username, it accepts the user if no Microsoft login is linked to them yet. It also accepts the user if the linked Microsoft id already equals the current one, ignoring case. A missing user record is never safe.

## Username relink with a different linked Microsoft id rechecks identity

Lines 332-339

Asked as: linked to different Microsoft login; name check before relink; relink rejected different Microsoft

If the user found by username has a different Microsoft id linked, Turbo does not accept them automatically. Turbo accepts the relink only if the user's name matches the Microsoft name, or, when Microsoft gives both given name and surname, their profile name matches. Otherwise Turbo rejects the relink.

## Relink replaces the old Microsoft link

Lines 344-371

Asked as: replace old Microsoft link; missing provider id relink; relink provider id error

When relinking, Turbo deletes the user's existing Microsoft link if it has a different provider id. If the link already has the current provider id, it leaves it alone. A relink with an empty provider id fails with "missing provider id for Microsoft relink".

## Email existence check ignores the user's own record

Lines 373-388

Asked as: email already in use; email exists for another user; duplicate email error

Turbo treats an email as taken when another user has it. The user's own record is excluded. An empty email counts as not taken.

## Username is not updated at later logins

Lines 390-394

Asked as: username changed after login; why username stays same; username not updating login

Turbo sets the username only when the user is created. Later Microsoft logins sync name and email but leave the username unchanged.

## Login syncs name from Microsoft

Lines 395-408

Asked as: name updated at login; name changed from Microsoft; name sync on sign in

On a returning Microsoft login, Turbo updates the user's name to the Microsoft full name when it differs, ignoring case. It does not change the username.

## Login name sync ignores case

Lines 403-408

Asked as: capital letters name change; name case update login; name only changes case

At login Turbo updates the user's name to the Microsoft full name only when the two differ, ignoring case. It does nothing when Microsoft gives no name.

## Login syncs email unless another user owns it

Lines 410-425

Asked as: email updated after login; email not changing at login; new email from Microsoft

On a returning Microsoft login, Turbo changes the user's email to the Microsoft email when they differ. It skips the change if another user already has that email.

## Login marks user verified when email matches Microsoft

Lines 427-440

Asked as: account not verified; verified after Microsoft login; email verification at sign in

If the user is not verified and the stored email equals the Microsoft email (ignoring case), Turbo marks the user verified. The record is saved only if something changed.

## Onboarding creates only an admin profile

Lines 461-467

Asked as: no business profile created; manager field missing onboarding; only admin profile made

After Microsoft sign-in, Turbo ensures the user has an admin profile. It does not create the business profile, because Microsoft does not provide fields such as manager.

## New users get the time claim once

Lines 469-507

Asked as: time claim missing; duplicate time claim; time claim added twice

Turbo gives a user the time claim if the user does not already have it. It does not add a duplicate.

## Default admin profile values for new users

Lines 522-553

Asked as: default profile settings new user; skip min time check default; new user settings off

A missing admin profile is created as active, with default work week hours, default charge-out rate and default branch from Turbo's constants. It has skip_min_time_check set to "no", and salary, untracked time off and time sheet expected all off.

## New admin profile is created only when none exists

Lines 522-565

Asked as: duplicate admin profile; profile already exists error; admin profile created twice

Turbo creates an admin profile only when the user has none. If a concurrent save already created one, Turbo treats it as success.

## Admin profile already existing is not an error

Lines 527-565

Asked as: admin profile exists at login; profile already there; save failed profile exists

If an admin profile already exists for the user, Turbo makes no change. If a save fails because another admin profile for the user was created in the meantime, Turbo treats that as success.

## Placeholder payroll id is generated per user

Lines 540-614

Asked as: payroll id generated; payroll id nine digits; failed unique payroll id

A new admin profile gets a placeholder payroll id in the 900000000 to 999999999 range, derived from the user id. If that id is already used, Turbo tries the next ids, up to 100 attempts. After that it fails with "failed to create unique payroll id placeholder".

## Payroll placeholder retries up to 100 times

Lines 540-578

Asked as: payroll id taken retry; cannot create admin profile payroll; payroll number already used

If a generated payroll id is already taken, Turbo tries the next value, up to 100 attempts. After that, admin profile creation fails.

## Placeholder payroll ids fall in the 9xxxxxxxx range

Lines 609-614

Asked as: payroll id starts with 9; placeholder payroll number range; payroll number format

The placeholder payroll id is 900000000 plus a number from a hash of the user id. Each retry adds 1 and wraps within the same range.

## New username comes from Microsoft email local part

Lines 631-770

Asked as: username from email address; how username is chosen; username based on email

For a new user, Turbo builds the username from the part of the Microsoft email (or user principal name) before the @. It lowercases it and keeps letters, digits, _, . and -. Other characters become -, and leading or trailing - and . are removed.

## Username fallbacks when email gives nothing usable

Lines 642-648

Asked as: username starts with u-; invalid first character username; username too long

If the email gives no usable username, Turbo uses "u-" plus a short form of the provider id, up to 12 characters. If that also fails, it uses "u-user". A username whose first character is invalid gets a "u" prefix. The length limit is the username field maximum, or 150 if none is set.

## Taken usernames get a numeric suffix

Lines 650-665

Asked as: username already taken; username ends with -2; unable to allocate username

If the username is taken (compared ignoring case), Turbo appends -2, -3 and so on, shortening the base to fit the length limit. After 100 attempts it fails with "unable to allocate unique username".

## Email local part rules

Lines 670-685

Asked as: text before at sign; email with no at sign; email starting with at

Turbo takes the text before the first @. An address with no @ is used whole, and an address that starts with @ gives an empty result.

## Fallback username from the Microsoft id

Lines 687-694

Asked as: username from Microsoft id; u- username id; username from user id

If the email part gives no usable username, Turbo uses u- plus the Microsoft id cut to 12 characters. That part is user when the id gives nothing. As a last resort it uses u-user.

## Username is trimmed and cleaned to valid characters

Lines 708-774

Asked as: special characters in username; username lowercase cleaned; dots and dashes username

Turbo lowercases the base and keeps letters, digits, _, . and -. It replaces other characters with a single -. It strips leading and trailing - and ., and prefixes u if the first character is invalid. It cuts the result to the maximum length.

## Username is cut to fit the length with its suffix

Lines 779-802

Asked as: long username suffix cut; username truncated with number; suffix makes username long

For a collision suffix (-2, -3 ...), Turbo shortens the base to make room, then removes trailing - and .. If the base is empty, it uses u.

## Username collision check ignores case

Lines 806-824

Asked as: username taken different case; capital letters username exists; username clash ignoring case

Turbo treats a username as taken when another user has it, ignoring upper and lower case. It ignores the user's own record.
