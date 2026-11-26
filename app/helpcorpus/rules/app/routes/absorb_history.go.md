<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_history.go","sha256":"46f55b60b6706f0948e03102d81a0e6a07345b50b6a21c162507fa34d5f50e08","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_history.go

## Absorb history is kept for a limited time

Lines 17-28

Asked as: how long absorb history is kept; absorb history retention period; old absorbs disappear; how far back absorb log goes

The absorb history lists committed absorbs only for ABSORB_HISTORY_RETENTION_DAYS, so that a bad absorb can be traced. Each entry shows who absorbed, who committed, the kept record, the names of the absorbed records and the number of references moved.

## Absorb history list leaves out saved copies and moved references

Lines 33-53

Asked as: absorb list missing saved copies; why no references in absorb list; absorb overview hides details; where are the moved references

The absorb history list shows each entry's id, collection, kept record (target), proposal, who absorbed and when, and who committed and when. It also shows the names of the absorbed records and the number of references moved. It does not include the saved copies of the absorbed records or the list of references that moved. Those appear only when a person opens one entry.

## Who can see the absorb history

Lines 61-70

Asked as: who can view absorb history; can't open absorb log; absorb history access denied; need absorb claim to see

A holder of the absorb claim sees every row of the absorb history. An admin without the absorb claim sees only the rate sheet rows. Anyone else gets the message "you may not see the absorb history".

## Absorb history entry names absorbed records

Lines 72-90

Asked as: absorbed record shows no name; absorb entry name blank; alias in brackets absorb; absorbed contact shows id

Turbo names each absorbed record from its saved copy. It uses the name. If the name is empty, it uses the given name and surname of a contact. If both are empty, it uses the record id. When the record has an alias, Turbo adds it in parentheses after the name.

## Absorb history counts references moved

Lines 92-112

Asked as: how many references were moved; absorb reference count wrong; count of updated rows absorb; which names were absorbed

Turbo counts the references moved as the number of saved rows that the absorb updated, summed over all tables and columns. It lists the names of the absorbed records from the saved copies.

## Listing the absorb history

Lines 114-137

Asked as: absorb log newest first; admin only sees rate sheets; list of committed absorbs; who absorbed and committed

The list shows committed absorbs, newest first. Each row has who absorbed, who committed, the kept record, the names of the absorbed records and the number of references moved. An admin without the absorb claim sees only rows whose collection is rate_sheets.

## Viewing one absorb history entry

Lines 139-166

Asked as: open one absorb entry; absorb entry not found; absorb record missing from history; can't see this absorb details

A single entry adds the saved copies of the absorbed records and the references that moved. A user who is not allowed to see the entry's collection gets the same message as for a missing entry: "this absorb is not in the history". A user with no access at all gets "you may not see the absorb history".

## Absorb history entry shows raw saved data

Lines 159-174

Asked as: absorb detail shows null; broken json in absorb; raw data of absorb; saved copy shows null

The detail of one absorb history entry includes the saved copies of the absorbed records and the references that moved. When the stored JSON is not valid, Turbo shows null in its place.
