<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/vendors.md","sha256":"64d450871414a125851f4c3fea5e7e5736a3e4b8f29f7e8b05791f74122422d6","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/vendors.md

## Who can list vendors

Lines 9-11

Asked as: see vendor list; vendors not showing; who can browse suppliers; view all vendors

Any signed-in person can list vendor records. Someone who is not signed in cannot.

## Who can view one vendor

Lines 15-17

Asked as: open a vendor record; can't see vendor details; view one supplier; look up vendor

Any signed-in person can view a single vendor record. Someone who is not signed in cannot.

## Who can create a vendor

Lines 21-27

Asked as: add new supplier; can't add vendor; create supplier record; staff allowed to add vendors

A signed-in person can create a vendor only if they hold the time claim (time) or the payables admin claim (payables_admin). Other staff cannot create vendors.

## Who can update a vendor

Lines 31-33

Asked as: edit vendor details; can't change supplier; update vendor record; time staff editing vendor

Only a person with the payables admin claim (payables_admin) can update a vendor record. Staff with only the time claim can create vendors but cannot update them.

## Who can delete a vendor and when deletion is blocked

Lines 37-45

Asked as: remove supplier; delete vendor button greyed; vendor won't delete; why deletion blocked

Only a signed-in person with the payables admin claim (payables_admin) can delete a vendor. Turbo blocks the deletion if any purchase order or any expense still refers to the vendor.

## Vendor name and status requirements

Lines 52-55

Asked as: vendor name too short; valid vendor status; alias minimum length; vendor form won't save

A vendor must have a name (name) of at least 3 characters. It must have a status (status) of either "Active" or "Inactive". An alias (alias) is optional, but if given it must have at least 3 characters.

## Vendor names must be unique

Lines 64-64

Asked as: supplier name already taken; duplicate vendor name; same vendor twice; name exists error

Two vendors cannot have the same name (name). Turbo refuses a duplicate name.

## Vendor aliases must be unique when filled in

Lines 65-65

Asked as: alias already used; duplicate vendor alias; short name clash; two vendors same nickname

Two vendors cannot have the same alias (alias). The check ignores empty aliases, so many vendors can have no alias.
