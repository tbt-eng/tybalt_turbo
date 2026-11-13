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
clients and vendors, beside each contact on the client's page, and beside each sheet that can absorb on the rate sheets Duplicates view. Options for
the icon are:

    - fluent-emoji-high-contrast:sponge
    - mdi:merge
    - material-symbols:merge

The buttons use the mdi:merge icon.

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

# Rate sheets

Duplicate rate sheets are absorbed like the other collections (issue #191).
The rule behind every check: an absorb never reprices time.

- Claims: `admin`, or both `absorb` and `rate_sheet_revise`. `absorb` alone is
  refused for absorb, undo and commit.
- References: `jobs.rate_sheet` and `job_requests.rate_sheet`. The entries
  (`rate_sheet_entries`) are owned rows.
- Rates: the target and the absorbed sheet must have the same set of roles,
  with the same rate and overtime rate in cents. A missing role is a
  difference. A sheet without entries is never a duplicate.
- Names: the test is on the name, not on the selected revision. A sheet that
  is the only revision of its name can be absorbed. A name with revisions can
  be absorbed only when no job and no project request uses any of its
  revisions; then every revision and its entries are absorbed in one action,
  and only the latest revision must match the target. Otherwise no revision
  of the name can be absorbed. A revision of the target's own name can't be
  absorbed.
- Effective date: the target's date can't be later than the date of a single
  revision sheet that it absorbs. The rule does not apply to a name with
  revisions, because nothing moves.
- Target: any revision of any sheet with identical rates. An inactive target
  can absorb only inactive sheets.
- Commit writes one `rate_sheet_name_changes` row for each absorbed name, with
  the target as its anchor and the committing user as `renamed_by`, so the
  target shows "Formerly ...". The history rows of the absorbed sheets move to
  the target on commit too. Nothing is written on absorb.
- A request that becomes a job while the action is pending moves its pending
  undo data to the job (`transferRequestMergeReferences`), as for clients.
- Moved jobs get `_imported = false` and a new `updated`, so the jobs
  writeback sends them with the target's id, and its rate sheet and entries
  queries include the target sheet.

The UI: the Duplicates view (`/rate-sheets/duplicates`) lists each group of
identical sheets from `GET /api/rate_sheets/duplicates`, marks the target (the
earliest effective date) and shows why a sheet or group can't be absorbed.
The absorb button opens `/rate-sheets/[id]/absorb`, which offers only the
sheets from `GET /api/rate_sheets/{id}/absorb_candidates` that the server would
accept.
