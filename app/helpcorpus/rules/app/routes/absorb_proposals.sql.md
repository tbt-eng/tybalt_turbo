<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_proposals.sql","sha256":"035bd2bd67ea6690da2c74cb376e11ec7856fe05a2641e615101c9950bb62aa8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_proposals.sql

## Absorb proposal list shows requester and decider names

Lines 14-19

Asked as: who decided absorb proposal; requester name missing; absorb list blank name; decider name empty

The absorb proposals list shows the requester's name (requester_name) and the name of the person who decided (decided_by_name). Each name is the given name and surname from that person's profile. The name is an empty string when the profile is missing, or when nobody has decided yet.

## Absorb proposal shows an in-progress flag

Lines 16-16

Asked as: absorb proposal still running; is absorb in progress; absorb proposal flag wrong; absorb action started status

An absorb proposal is marked in progress (in_progress) when at least one absorb action (absorb_actions) refers to it. Otherwise the flag is false. The flag does not depend on the proposal's status.
