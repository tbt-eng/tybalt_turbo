<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/users.md","sha256":"cd9aa2242ade6cbfdf36ee2a79cb6dcd2723c16cad6790c3ab210af85898cca8","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/users.md

## Who can list user records

Lines 7-11

Asked as: see other users list; why can't I see everyone; staff directory empty; list of all staff

A signed-in person can list only their own user record (id = @request.auth.id). Other people's records do not appear in the list, and the list rule has no exceptions for approvers or the tapr claim.

## Who can view one user record

Lines 13-33

Asked as: can't open colleague profile; view someone's profile; why can't I see user details; profile access denied

A signed-in person can view one user record in these cases: it is their own; it belongs to someone with the po_approver claim or the tapr claim; the caller is the viewer (user_view_grants viewer) of that user, such as an approver of that user's committed expenses; the record has a view grant with reason 'active_po' (the user is the uid, approver or second approver of an active purchase order); or the caller holds the tapr claim. A person who is not signed in cannot view any record.

## Anyone can create a user record

Lines 36-38

Asked as: sign up without login; register new account anonymous; create account without signing in

The create rule is empty, so anyone can create a user record through the API, including a person who is not signed in.

## Only the account owner can update a user record

Lines 40-44

Asked as: edit colleague's details; change someone else's name; update another user's profile

A signed-in person can update only their own user record (id = @request.auth.id). Nobody can update another person's record through the API, not even approvers or tapr holders.

## Nobody can delete user records through the API

Lines 46-48

Asked as: remove user account; delete staff member; how to delete a user

The delete rule is empty of any filter (no rule), so nobody can delete a user record through the API. Only a PocketBase superuser can delete one.

## User id format

Lines 54-54

Asked as: user id rejected; id length wrong; custom user id invalid

A user id (id) is required and must be exactly 15 characters, using only lowercase letters and digits.

## Password minimum length

Lines 55-55

Asked as: password too short; minimum password characters; password rejected when saving

A password is required and must have at least 8 characters. The password field is hidden.

## Email must be unique when set

Lines 57-57

Asked as: email already taken; duplicate email address; email in use error

Two users cannot have the same email address. The unique index applies only when the email is not empty, so many users can have no email.

## Username format, length and uniqueness

Lines 60-60

Asked as: username not allowed characters; username already exists; username with spaces invalid; username case duplicate

A username (username) is required, 3 to 150 characters, and must start with a letter, digit or underscore. After the first character it can contain only letters, digits, underscores, dots and hyphens. Two users cannot have the same username, and the comparison ignores case.

## User name length limit

Lines 61-61

Asked as: display name too long; name character limit; name too many letters

The display name (name) can have at most 255 characters. It is optional.

## Avatar file limits

Lines 64-64

Asked as: upload profile picture; avatar file type; photo format not accepted; change profile picture

A user can have one avatar (maxSelect 1). It must be a JPEG, PNG, SVG, GIF or WebP image.
