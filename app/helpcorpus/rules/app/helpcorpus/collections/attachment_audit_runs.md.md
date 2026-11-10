<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/attachment_audit_runs.md","sha256":"539da00dc6ee9358330e078136af584329c33d648877612a795b29f1e49c56bb","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/attachment_audit_runs.md

## Who can list attachment audit runs

Lines 7-11

Asked as: list audit runs; see all attachment audits; non admin can't see audit list; why am I refused audit runs

Only an admin (the admin claim) can list attachment audit run records through the API. Everyone else is refused.

## Who can view one attachment audit run

Lines 13-17

Asked as: view one audit run; open a single audit record; audit run details forbidden; can't open audit run

Only an admin (the admin claim) can view a single attachment audit run record through the API. Everyone else is refused.

## Nobody can create attachment audit runs through the API

Lines 19-21

Asked as: create audit run via api; add new audit record; api post audit run; manually make an audit run

Nobody can create an attachment audit run record through the API, not even an admin. Only a PocketBase superuser can. Turbo's own server code can still create them.

## Nobody can update attachment audit runs through the API

Lines 23-25

Asked as: edit audit run; change audit run record; update audit status by api; modify audit via api

Nobody can update an attachment audit run record through the API, not even an admin. Only a PocketBase superuser can. Turbo's own server code can still update them.

## Nobody can delete attachment audit runs through the API

Lines 27-29

Asked as: remove audit run; delete audit record; get rid of old audit runs; erase audit run

Nobody can delete an attachment audit run record through the API, not even an admin. Only a PocketBase superuser can. Turbo's own server code can still delete them.

## Only one audit run per target

Lines 36-59

Asked as: duplicate audit for same target; two audits same target; target key already exists; audit already exists for target

Two attachment audit run records cannot have the same target key (target_key). The field is required and the unique index (idx_attachment_audit_runs_target_key) enforces this.

## Attachment audit run status values

Lines 40-40

Asked as: audit status invalid; what statuses can audit have; audit stuck in pending; audit status error

The status of an attachment audit run (status) is required and must be one of three values: running, completed or failed.

## Audit report files hold one file up to 512 MB

Lines 50-51

Asked as: upload audit report too big; report file size limit; missing report upload fails; orphaned report attachment limit

The missing report (missing_report) and the orphaned report (orphaned_report) each hold at most one file. Each file may be at most 536870912 bytes (512 MB).
