# Absorb system

The parts make up the UI portion of the absorb system.

1. A full-view component (AbsorbEditor.svelte) that shows the target record (for
example a clients record) at the top using a snippet described in the containing
+page.svelte file and a heading just below that says "records to absorb". Under
this heading is a list of zero or more records to be absorbed. Initially there
are zero items in the list. There is a search box (using the DSAutoComplete.svelte
component when the page passes a search index, otherwise the DSSelector.svelte
component) to immediately beneath the heading to find records from the same
collection as the target record (for example a clients record) and pressing
enter will add the selected record to the list of records to be absorbed. There
is a button at the bottom to absorb and another to cancel. Clicking the absorb
button absorbs the records and opens the Absorb Actions page (/absorb/actions),
where the pending absorb action is committed or undone. A contact merge opened
from a client's page stays on the page and shows the pending action there. While
a collection has a pending absorb action, the component shows that action in
place of the editor.

2. A +page.svelte file for each collection that supports absorption (clients, vendors, client contacts and rate sheets), contained within the /absorb directory under the collection's id directory (i.e. \[cid\]/absorb/+page.svelte under the clients route). The client contacts and rate sheets routes also have a +page.ts file. The AbsorbEditor.svelte component is used in the +page.svelte file. Arguments are passed to the AbsorbEditor.svelte component, namely the collectionName, the targetRecordId (which is the id of the absorbing record pulled from the URL), the availableRecords and an optional autoCompleteIndex. A snippet is also provided to the AbsorbEditor.svelte component to render the record to be absorbed. The snippet is defined in the +page.svelte file on the route.

3. A button to absorb the records which appears on the regular list view for
clients and vendors, beside each contact on the client's page (for everyone:
without the `absorb` claim, the editor proposes instead, see Absorb
proposals), and beside each sheet that can absorb on the rate sheets Duplicates view. Options for
the icon are:

    - fluent-emoji-high-contrast:sponge
    - mdi:merge
    - material-symbols:merge

The buttons use the mdi:merge icon.

# Absorb proposals

Only `absorb` holders absorb. Anyone signed in may propose an absorb of
clients, client contacts or vendors, and an `absorb` holder carries it out
with the normal absorb. Rate sheets take no proposal.

- Data. `absorb_proposals` holds the collection, the target, the requester,
  the reason (required, at most 1000 characters), the status (`pending`,
  `accepted`, `declined`, `withdrawn`, `obsolete`), `decided_by`,
  `decided_at` and `decision_comment`. `absorb_proposal_records` holds the
  records to absorb. Both keep the names of the target and the records, and
  of the client for contacts, because a committed absorb deletes the
  records. `absorb_proposal_records.open` is true while the proposal is
  pending, and a partial unique index on (`collection_name`, `record_id`)
  where `open` keeps a record in one open proposal. Neither collection has
  API rules; the routes under `/api/absorb_proposals` read and write them
  (`app/routes/absorb_proposals.go`, checks in `app/absorb/proposals.go`).
- Routes. `GET ""` lists pending proposals, or decided ones with
  `?status=decided`: all of them for an `absorb` holder, otherwise the
  caller's own. `GET /{id}` returns one to a holder or its requester.
  `POST ""` creates one. `POST /{id}/withdraw` is for the requester and
  `POST /{id}/decline` for a holder, with a required comment. Neither works
  while an absorb action carries the proposal. `DELETE /{id}` lets the
  requester delete their proposal once it is no longer pending (withdrawn,
  declined, accepted or obsolete); its `absorb_proposal_records` rows go
  with it (the relation cascades), in one transaction. Anyone else gets 403
  `unauthorized`, a holder included; a pending proposal gets 400
  `proposal_pending` and must be withdrawn first; an unknown id gets 404
  `proposal_not_found`. The Decided tab shows a trashcan on the requester's
  own proposals, with no confirmation, because the delete changes no record
  and a committed absorb keeps no record of its own either. Errors are JSON with `code` and `message`.
- Checks, on create and again when a holder absorbs from a proposal: one
  of the three collections; every record exists; the ids are distinct; the
  target isn't in the list; contacts share one client (`ParentConstraint`);
  no record, target or item, is an item of another pending proposal; and no
  item is the target of another pending proposal. Two proposals may keep
  the same target.
- Review. The queue (`/absorb/proposals`) links a holder to the target's
  absorb page with `?proposal=<id>`. The editor starts with the proposal's
  records and no search. The holder may remove a record or keep another as
  the target, but not add one: `POST /api/{collection}/{id}/absorb` with
  `proposal` refuses a target or record that isn't in the proposal. The
  absorb action stores the id in `absorb_actions.proposal`.
- Commit. The `absorb_actions` delete hook runs
  `absorb.SettleProposalsOnCommit` in the commit transaction, for these
  three collections only. It accepts the proposal that the action carries
  (decided by the committing user), and makes obsolete, with an automatic
  comment, every other pending proposal that names an absorbed record as
  its target or an item. A committed client absorb also gives the pending
  contact proposals of the absorbed clients the kept client. Absorb alone
  changes no proposal, so an undo leaves it pending.
- Notifications. The requester is told when a proposal is accepted,
  declined or made obsolete (`absorb_proposal_accepted`, `_declined`,
  `_obsolete`), unless they decided it themself. Holders get no message per
  proposal; the menu badge counts the pending ones.

# Optional config hooks

`absorb.Config` (`app/absorb/config.go`) has optional fields. A collection that
leaves them unset behaves as clients, client contacts and vendors always have.

- `ClaimSets` replaces the `absorb` claim check for absorb, undo and commit. A
  user passes with every claim of one set. The `absorb_actions` collection
  rules require `absorb` to list, view and commit an action. They also admit
  an `admin` to rate sheet actions only, and the delete hook then checks the
  claim sets. A new collection whose claim sets admit a user without `absorb`
  needs the same change to those rules.
- `OwnedRows` are child rows of an absorbed record. They are saved in
  `absorbed_records` under the key `@owned_rows`, deleted with the record and
  recreated on undo. Absorb uses plain SQL, so a `cascadeDelete` never runs.
- `Prepare` runs in the absorb transaction before anything changes. It refuses
  the request with an `absorb.RefusedError` (HTTP 400 with its message) or
  returns the ids to absorb, which may include more than were asked for.
- `OnCommit` runs in the same transaction as the delete of the action.
- `ChangesJobs` refuses the absorb while job editing is turned off.

A collection whose config needs code from a package that imports `absorb`
adds itself with `absorb.Register` from an `init` function. Rate sheets do this
in `app/hooks/rate_sheet_absorb.go`.

An absorb that moves no reference still saves `updated_references`, with an
empty map for each configured table and column, because the field is required.

Absorb and undo both change references with plain SQL. In a table with an
`_imported` column, both set `_imported = false` and a new `updated` on each
record they change, so the writeback sends the record again: with the target's
id after an absorb, and with the restored id after an undo.

# Rate sheets

Duplicate rate sheets are absorbed like the other collections (issue #191).
The rule behind every check: an absorb never reprices time.

- Claims: `admin`, or both `absorb` and `rate_sheet_revise`. `absorb` alone is
  refused for absorb, undo and commit.
- References: `jobs.rate_sheet`, `job_rate_sheets.rate_sheet` (the periods
  of a project's rate sheet history, which price its time) and
  `job_requests.rate_sheet`. The entries (`rate_sheet_entries`) are owned
  rows.
- Use: a job uses a sheet when `jobs.rate_sheet` is the sheet or any of its
  periods is on it, so a sheet that only an old period uses is in use.
  `job_count` in the Duplicates view and the candidates counts those jobs
  once each.
- Rates: the target and the absorbed sheet must have the same set of roles,
  with the same rate and overtime rate in cents. A missing role is a
  difference. A sheet without entries is never a duplicate.
- Sheets: the test is on the sheet, which its code identifies (all its
  revisions share the code), not on the selected revision. A sheet that is
  its only revision can be absorbed. A sheet with revisions can be absorbed
  only when no job, no period and no project request uses any of its
  revisions; then every revision and its entries are absorbed in one action,
  and only the latest revision must match the target. Otherwise no revision
  of the sheet can be absorbed. A revision of the target itself can't be
  absorbed. Before issue 207 the revisions were linked by name.
- Codes: the target keeps its code. The code of an absorbed sheet is
  retired: `rate_sheet_codes` keeps it, so it is never issued again, and an
  undo puts the sheet back with it.
- Effective date: the target's date can't be later than the date of a single
  revision sheet that it absorbs. The rule does not apply to a sheet with
  revisions, because nothing moves. A dated period starts on or after its
  sheet's effective date, so after a move it also starts on or after the
  target's, as the history rules require. The opening period has no date.
- Target: any revision of any sheet with identical rates. An inactive target
  can absorb only inactive sheets. A moved period may then be on an inactive
  sheet, as a period keeps its sheet when the sheet is deactivated.
- Commit writes one `rate_sheet_name_changes` row for each absorbed name, with
  the target as its anchor and the committing user as `renamed_by`, so the
  target shows "Formerly ...". An absorbed sheet with no name adds no row. The
  name and code history rows (`rate_sheet_code_changes`) of the absorbed
  sheets move to the target on commit too. Nothing is written on absorb.
- A request that becomes a job while the action is pending moves its pending
  undo data to the job (`transferRequestMergeReferences`), as for clients.
  The job's opening period, written on the target with the job, is added to
  the undo data too, so undo moves both back.
- Moved jobs get `_imported = false` and a new `updated`, so the jobs
  writeback sends them with the target's id, and its rate sheet and entries
  queries include the target sheet. `job_rate_sheets` has no `_imported`
  column, so a moved period changes neither its own `updated` nor its job.
  A job whose only moved reference is an older period is not sent again: the
  writeback sends only `jobs.rate_sheet`, which did not change.
- `jobs.rate_sheet` is the sheet of the job's latest period. The job and the
  latest period move together, and undo puts both back, so they stay equal.
- While the action is pending, the history routes refuse (409) any change to
  a project whose job or periods the action moved
  (`checkPendingRateSheetAbsorb`), as a client change does during a client
  merge. Undo writes back the saved sheet of each moved row, so a change in
  between would be overwritten or leave `jobs.rate_sheet` unequal to the
  latest period. A period can't be pointed at an absorbed sheet, because the
  sheet is deleted.
- A move can leave two neighbouring periods on the target. The history rules
  then refuse only a pair that a change makes, so the pair does not block
  other changes, and deleting the later period removes it.

The UI: the Duplicates view (`/rate-sheets/duplicates`) lists each group of
identical sheets from `GET /api/rate_sheets/duplicates`, marks the target (the
earliest effective date) and shows why a sheet or group can't be absorbed.
The absorb button opens `/rate-sheets/[id]/absorb`, which offers only the
sheets from `GET /api/rate_sheets/{id}/absorb_candidates` that the server would
accept.
