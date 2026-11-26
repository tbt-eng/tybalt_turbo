<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/absorb_suggestions.go","sha256":"eb35702d1d2913b9e87475414fd026bb18563da28de735c881d40cbb5203ebc2","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/absorb_suggestions.go

## Reference counts decide each record's usage and the suggested target

Lines 31-37

Asked as: how many uses; usage count; which duplicate is main; most used record; reference count

For each record Turbo counts its uses, and the suggested target of a bucket is the member with the most. For clients it counts jobs that use the client as client and as job owner (job_owner). For vendors it counts purchase orders and expenses that use the vendor. A record with no uses counts 0. Empty references are not counted.

## Absorb suggestions cover clients and vendors only

Lines 34-93

Asked as: suggestions for jobs; merge employees; absorb not available; duplicates for purchase orders; only clients and vendors

Turbo refuses absorb suggestions and dismissals for any collection other than clients or vendors. The message is "absorb suggestions are for clients and vendors only".

## Only users who may absorb the collection can use suggestions

Lines 94-100

Asked as: can't see suggestions; no access to merge; not allowed to absorb; permission for duplicates; you may not absorb

Turbo checks that the user holds the claims needed to absorb the collection (absorb.HasRequiredClaims). A user without them cannot list or dismiss suggestions and sees "you may not absorb these records". The code does not show which claims are required.

## Suggestion level must be strict, normal or loose

Lines 110-113

Asked as: strictness setting; how strict matching; loose matching level; level not accepted; matching sensitivity

The list of suggested buckets takes a level of strict, normal or loose. Any other level is refused with "level must be strict, normal or loose". Turbo uses the level to choose the matching settings (suggest.LevelParams); this file does not show what each level does.

## Suggestions are computed fresh on every request

Lines 122-148

Asked as: suggestions not updating; stale duplicate list; refresh duplicate list; why list changed; duplicates not saved

Turbo computes the buckets of probable duplicate clients or vendors each time the list is requested. It reads every record of the collection and does not store the buckets.

## Dismissed buckets are hidden and counted

Lines 148-152

Asked as: hide duplicate pair; stop suggesting this; not duplicates count; dismissed suggestions hidden; dismissal applies everyone

Turbo hides a bucket that a dismissal covers ("Not duplicates"). It adds one to the dismissed count for each hidden bucket. Dismissals apply to every user of that collection.

## Each bucket reports its target, confidence, moved references and rank

Lines 153-172

Asked as: confidence score meaning; how many would move; suggested target field; rank of suggestion; score for member

A bucket lists its suggested target (key and target_id), a confidence, the number of references an absorb would move (moved_references) and a rank. Each member shows its name, alias, reference count and score. Turbo takes these values from the suggest package, and this file does not show how it calculates them.

## Reference counts include only non-blank uses of a record

Lines 178-204

Asked as: blank client not counted; empty job client; missing alias blank; zero uses shown; empty vendor field

Turbo counts the uses of each client or vendor from jobs (client and job_owner) or from purchase orders and expenses (vendor). A blank reference is not counted. A record with no uses has 0 references, and an alias that is missing shows as blank.

## Members show an open absorb proposal that names them

Lines 206-230

Asked as: proposal already open; why shows pending merge; merge already started; proposal id on member; record already proposed

A member of a bucket shows the id of an open absorb proposal that names it, either as a record to absorb or as the target of a pending proposal. Otherwise it shows blank. When several proposals name one record, the first one found is used.

## An open absorb proposal is a pending one or an open proposal record

Lines 208-229

Asked as: pending merge id; already being absorbed; which proposal is open; open proposal record; proposal target pending

A record shows an open absorb proposal when it is listed in an open proposal record, or when it is the target of a pending proposal. If one record is named by several proposals, Turbo shows the first one it finds.

## A dismissal hides a bucket that it covers

Lines 232-250

Asked as: duplicates hidden after dismiss; dismissal breaks suggestions; suggestions request fails; saved dismissals load; covered duplicates removed

Turbo loads every dismissal saved for the collection, each one a list of record IDs. A bucket is hidden when the dismissals cover its records (Dismissals.Covers). A dismissal that cannot be read makes the suggestion request fail.

## Suggestions show the collection's pending absorb action

Lines 252-287

Asked as: last merge shown; pending merge in response; absorbed records list; when merge was created; pending action missing

The response includes the collection's pending absorb action, if there is one: its target and the target's name, the absorbed records with their aliases in parentheses, when it was created and who absorbed. The absorbed-by user and name are blank for an action made before that field existed. When there is no pending action, the value is empty.

## Pending action shows who absorbed and the names absorbed

Lines 254-286

Asked as: who merged records; absorbed by name; target name blank; merged by user; who absorbed clients

A pending absorb action shows its target name, the user who absorbed and that user's name. Turbo lists each absorbed record by name, with its alias in parentheses when the alias is not blank. For an action made before the field existed, the absorbing user is blank. The target name stays blank if the target record no longer exists.

## A collection with no pending absorb action shows none

Lines 255-258

Asked as: no merge shown; pending action null; nothing pending to absorb; empty merge status; no absorb action

When the collection has no absorb action, the suggestions show no pending action (null).

## A dismissal with a bad request body is refused

Lines 295-300

Asked as: dismiss request invalid; malformed dismiss request; bad dismiss body; dismissal body error; request body not valid

When the body of a dismiss request is not valid, Turbo refuses it with the message "the request body is not valid".

## Dismiss needs at least two distinct records

Lines 301-310

Asked as: dismiss only one record; select two clients; need two records; same record twice; choose two records

A "Not duplicates" dismissal must list at least two records. A record listed twice is refused with "a record is listed twice". Fewer than two records is refused with "choose at least two records". Turbo checks for duplicates before it counts the records.

## A duplicate record in a dismissal is refused

Lines 301-307

Asked as: same client twice; pick same record; duplicate entry in dismiss; record listed twice error; repeat record in dismissal

When a dismiss request lists the same record twice, Turbo refuses it with the message "a record is listed twice". Turbo checks this before it checks the number of records.

## Fewer than two records in a dismissal is refused

Lines 308-310

Asked as: dismiss single record; only one client selected; need more records; not enough records; dismiss with one

When a dismiss request has fewer than two records, Turbo refuses it with the message "choose at least two records".

## Dismiss refuses records that no longer exist

Lines 311-319

Asked as: deleted client in dismissal; record gone dismiss; dismiss missing record; record no longer there; removed client dismiss error

Before it saves a dismissal, Turbo checks that every listed record still exists in the collection. If any record is missing it refuses with "a record no longer exists".

## A missing record in a dismissal is refused

Lines 317-319

Asked as: record was deleted; no longer exists error; dismiss deleted record; client already removed; record missing from dismissal

When any record of a dismiss request no longer exists in the collection, Turbo refuses it with the message "a record no longer exists".

## A dismissal is saved for the group and the user

Lines 321-332

Asked as: where dismissal is saved; dismissal id returned; order of records dismissed; who dismissed duplicates; save not duplicates

Turbo saves a dismissal with the collection, the sorted list of record ids and the user who dismissed. It returns the dismissal's id. Because the ids are sorted, the order in the request does not matter.
