<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/managers.md","sha256":"87dd4f07b77ca7db005229c39173faa5f2c1a007e53c49898c871a6a129e31a5","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/managers.md

## Managers view is read-only through the API

Lines 5-5

Asked as: edit manager record; cannot change managers list; managers view read only; delete manager through api

The managers view has no create, update or delete rules in this file. Nobody can create, update or delete its records through the API. Turbo's own server code is not limited by these rules. The view is a saved query that Turbo reads and never writes.

## Any signed-in person can list managers

Lines 7-11

Asked as: list managers api; who can see managers; managers list not loading; signed in list managers

Any signed-in person can list records of the managers view through the API. A request from someone who is not signed in fails the filter (`@request.auth.id != ""`).

## Any signed-in person can view one manager

Lines 13-17

Asked as: view single manager; open manager details; manager record fails unauthenticated; get one manager api

Any signed-in person can view one record of the managers view through the API. A request from someone who is not signed in fails the filter.

## Who counts as a manager: holders of the tapr claim

Lines 19-27

Asked as: who is a manager; why manager missing from list; tapr claim holders; manager appears twice; how do I become manager

The managers view lists each person whose profile is linked to the claim named 'tapr' (through user_claims and claims). Each record shows the person's id (the profile uid), surname and given name. A person without the tapr claim does not appear in the view. A person who holds the claim appears once for each link to it.

## Name field limits in the managers view

Lines 29-35

Asked as: name too long manager; surname hyphen not allowed; apostrophe in name rejected; manager name characters allowed; given name validation error

Surname and given name each need 2 to 48 characters. They allow letters only. A surname may also have hyphens or apostrophes between letter groups, such as O'Neil or Smith-Jones. A given name may have hyphens but no apostrophes. The id allows only lowercase letters and digits.
