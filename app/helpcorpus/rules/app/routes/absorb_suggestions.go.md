<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_suggestions.go","sha256":"1269131ca9c2129ed95c2b344e014c975d21d2f071ad1357bd4e95956ede8d4c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_suggestions.go

## Suggestion buckets are listed best first with a rank

Lines 19-58

Asked as: best duplicate first; ranking of duplicate suggestions; which suggestion to merge first; order of suggested duplicates

Absorb suggestions are buckets of probable duplicate clients or vendors, listed best first. Each bucket carries a confidence, a rank and the number of references that an absorb would move (moved_references). Each member shows its name, alias, reference count and score.

## A dismissal hides the bucket for every user

Lines 24-25

Asked as: hide suggestion for everyone; not duplicates for all users; dismissed by a colleague still showing; stop suggesting these clients

Turbo hides a bucket when a saved "Not duplicates" dismissal for that collection covers it. A dismissal by one user applies to all users of that collection.

## Suggested target is the record with most references

Lines 31-37

Asked as: which record to keep; which client survives merge; pick the main record; most used client kept

Turbo counts how many times each client or vendor is used. A client counts as the client or the owner (job_owner) of a job. A vendor counts on purchase orders and on expenses. These uses move to the kept record on absorb. The bucket's target is the member with the most (suggestionRefs).

## What counts as a reference when picking the target

Lines 34-37

Asked as: how references are counted; what counts as usage; jobs linked to a client; purchase orders per vendor count

For clients, Turbo counts jobs that use the client (client) or name it as job owner (job_owner). For vendors, it counts purchase orders and expenses that use the vendor (vendor). Blank references are not counted, and a record with no uses has 0 references.

## Each record appears in one suggestion bucket only

Lines 50-52

Asked as: same record in two groups; duplicate appears twice in suggestions; record in multiple buckets; group key of suggestion

A record is in one bucket only. The bucket key is the suggested target's id (key, target_id).

## Pending absorb action does not show who absorbed

Lines 74-77

Asked as: who merged these records; who absorbed the clients; absorb history user; audit of absorb action

The pending absorb action shown with suggestions has no field for the person who absorbed. Turbo does not record who absorbed in its absorb actions (absorb_actions). The action shows only its id, the target, the absorbed records and when it was created.

## Absorb suggestions cover clients and vendors only

Lines 86-89

Asked as: duplicates for jobs or expenses; suggestions not working for projects; dismiss suggestion for other records; merge suggestions unavailable for this list

Turbo refuses a request for absorb suggestions or a dismissal on any collection other than clients or vendors. The message is "absorb suggestions are for clients and vendors only".

## Only staff who may absorb the collection see or dismiss suggestions

Lines 90-96

Asked as: cannot see duplicate suggestions; no permission to merge clients; suggestions missing for my account; can't dismiss duplicate group

Turbo checks that the person holds the claims needed to absorb the collection (absorb.HasRequiredClaims). Without them Turbo refuses to list suggestions or to dismiss a bucket. The message is "you may not absorb these records".

## Suggestion level must be strict, normal or loose

Lines 103-109

Asked as: invalid suggestion strength; level option error; strictness setting wrong; matching sensitivity value rejected

The list of absorb suggestions takes a level of strict, normal or loose. Any other level gets the message "level must be strict, normal or loose". The access check runs before the level check.

## Suggestions are computed on each request

Lines 120-170

Asked as: suggestions not refreshing; duplicate list updates; why list changed; empty duplicate groups

Turbo finds the buckets of probable duplicates again on every request, from all records of the collection. Each bucket shows its target, confidence, moved references and rank. Each member shows name, alias, reference count and score. Empty buckets return an empty list.

## Dismissed buckets are hidden and counted

Lines 129-148

Asked as: dismissed count went up; where did the group go; hidden duplicate groups; suggestions disappeared after dismissing

When a bucket's records are covered by a "Not duplicates" dismissal for the collection, Turbo hides the bucket. It adds the bucket to the dismissed count in the response. Dismissals apply to every user.

## Reference counts cover every record in the collection

Lines 172-200

Asked as: usage count for each client; reference numbers blank alias; vendor usage totals; counts for all records

Turbo reads every client or vendor with its reference count. A blank alias is treated as an empty alias.

## Members show their open absorb proposal

Lines 204-225

Asked as: blank pending merge column; which merge proposal includes this client; proposal id on suggestion; record already queued for absorb

Each member shows the id of an open absorb proposal that names it, either as a record to absorb (open records) or as the target of a pending proposal. The member shows a blank when there is none. If several proposals name the record, the first one found is shown.

## Dismissals only apply to their own collection

Lines 228-246

Asked as: dismissal hides vendor duplicates; client dismissal affects vendors; vendor not dismissed in client list; dismissed records wrong list

Turbo loads only the "Not duplicates" dismissals saved for the collection being listed (collection_name). A dismissal of client records does not hide vendor suggestions, and a dismissal of vendor records does not hide client suggestions. Each dismissal stores the set of record ids that was dismissed (record_ids).

## Pending absorb action shown with suggestions

Lines 250-281

Asked as: undo merge banner; pending merge not showing; commit or undo merge; last absorb waiting

The response includes the collection's absorb action that waits for a commit or an undo. It gives the target id and name, the created time, and the absorbed records as name with alias in parentheses when the alias is not blank. The response is null when no action exists. Turbo does not record who absorbed.

## Dismiss request with an invalid body is refused

Lines 289-294

Asked as: dismiss request fails; bad request when dismissing; cannot submit not duplicates; request body error dismiss

When the dismiss request body cannot be read, Turbo refuses it with "the request body is not valid". This check runs after the access check and before the record checks.

## Dismissal needs at least two different records

Lines 295-304

Asked as: same client selected twice; need two records to dismiss; duplicate id in selection; only one record chosen

A dismissal fails with "a record is listed twice" when an id repeats. It fails with "choose at least two records" when fewer than two ids remain. The duplicate check runs first.

## Dismissal needs every record to exist

Lines 305-313

Asked as: deleted record in dismissal; client was removed while selecting; record missing when dismissing; stale selection error

Turbo checks that every listed id exists in the collection. If one is missing, it refuses with "a record no longer exists".

## Dismissal saves the record ids sorted and who dismissed

Lines 315-326

Asked as: who marked not duplicates; undo a not duplicates choice; saved dismissal record; reverse dismissal of duplicates

A valid dismissal saves the collection, the sorted record ids and the person who dismissed (dismissed_by). Turbo returns the id of the dismissal. Turbo does not offer a way to undo a dismissal here.
