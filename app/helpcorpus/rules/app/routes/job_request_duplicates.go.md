<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_request_duplicates.go","sha256":"2ec5d8a03c4b85891146d9672b701530e384a64ebe10bea3db13fd6b332ff05c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_request_duplicates.go

## Which jobs the duplicate-work review lists for a project request

Lines 19-29

Asked as: duplicate work review jobs; which jobs show as duplicates; possible duplicate projects; jobs listed for project request

The review lists a job in the exact tier when it matches the request on proposal, client PO, client reference, or parent, whatever its age or status. It also lists the request's parent project and any sibling subjob with the same parent. It lists the client's other jobs only when they are open and recent. A job of another client is listed only if it is an exact match.

## Duplicate review lists exact matches first, then similar, then client

Lines 35-43

Asked as: duplicate list order; why exact matches on top; order of duplicate rows; tier order in review

Each listed row in the duplicate-work review has a tier (exact, similar or client). The list orders tiers as exact first, similar second and client third.

## Recent window is measured from today's date in UTC

Lines 45-52

Asked as: recent window timezone; duplicate check time zone; UTC date duplicate review; wrong day in duplicate check

The duplicate-work review measures its recent window from the current time in UTC. The server's local time zone does not change it.

## Fields compared in the duplicate-work review

Lines 65-89

Asked as: what fields are compared; duplicate check fields; how requests are matched; which details find duplicates

Turbo compares a project request with other jobs and requests on seven fields: client, proposal, parent, client PO (client_po), client reference (client_reference_number), location and description. It reads these from the request under review.

## Each listed job shows reasons, tier and last activity

Lines 91-112

Asked as: last activity date blank; duplicate row reasons; what is last activity; why job listed as duplicate

Each listed job carries its reasons (reasons), its tier (tier: exact, similar or client) and its last activity date (last_activity). The last activity is the later of the project's award date and its latest time entry date. It is empty when the project has neither.

## Blank values never match as duplicates

Lines 118-127

Asked as: blank fields matching; empty client PO duplicate; empty location same match; missing proposal not duplicate

Two blank IDs (client, proposal, parent) or two blank locations never count as a match. Two blank client POs or client references never match either.

## Client PO and reference compare ignoring case and spaces

Lines 118-123

Asked as: PO match ignoring case; client PO spaces matching; reference case sensitive; same PO different capitals

Client PO and client reference match when they are equal after trimming surrounding spaces and ignoring case. Location is compared with surrounding spaces trimmed, and case counts.

## Reasons Turbo gives for a duplicate match

Lines 134-167

Asked as: duplicate reason labels; why flagged as duplicate; what does same proposal mean; meaning of similar description

Each listed row carries reasons: Same client, Same proposal, Same parent, Same location, Same client PO, Same client reference, Similar description, and Parent project. Same proposal, Same parent, Same client PO, Same client reference and Parent project make the row an exact match. Same client, Same location and Similar description do not make it exact.

## Same client adds a reason but is not exact by itself

Lines 136-139

Asked as: same client not listed; client jobs only shown; why same client job missing; same customer no duplicate flag

A job with the same client as the request gets the reason "Same client". This alone does not make an exact match. The job also needs to be open and recent to be listed.

## Same proposal, client PO or client reference is an exact match

Lines 140-158

Asked as: same PO shows as duplicate; matching proposal number; same client reference listed; old job same PO

A job with the same proposal gets "Same proposal". A matching client PO gets "Same client PO". A matching client reference gets "Same client reference". Each one is an exact match, so Turbo lists the job whatever its age or status.

## Same parent or the request's parent is an exact match

Lines 144-147

Asked as: parent project duplicate; same parent job listed; parent shown in duplicates; child project of same parent

A job with the same parent as the request gets the reason "Same parent" and is an exact match. A job that is the request's own parent gets the reason "Parent project" and is an exact match. Both are listed whatever their age or status.

## Location match uses trimmed text and is only a supporting reason

Lines 148-150

Asked as: same location duplicate; location alone flagged; matching site not duplicate; address match duplicate

Turbo compares locations after trimming spaces and adds the reason "Same location" when they match. A blank location never matches. A matching location alone never makes a project exact and never lists it.

## Description similarity is compared only within the same client

Lines 159-161

Asked as: similar description other client; description match different client; similar wording flagged; same words duplicate warning

The Similar description reason applies only when the other job or request has the same client. Location alone never lists a job.

## How description words are compared

Lines 169-198

Asked as: description word filtering; which words ignored description; stopwords in description; short words duplicate check

Turbo lowercases a description and splits it on anything that is not a letter or digit. It drops stopwords such as and, the, for, project, phase, work, service, and services. It drops words shorter than 2 characters, and drops 2-character words that contain no digit. A short word with a digit, such as 17, is kept.

## Threshold for similar descriptions

Lines 200-221

Asked as: how similar is similar; how many words must match; description similarity threshold; wording match rules

Two descriptions are similar when they share at least 2 meaningful words and at least half of the meaningful words of the shorter one. If the shorter has exactly one meaningful word, that word must be shared. If the shorter has none, they are not similar.

## Last activity ignores dates that are not valid YYYY-MM-DD

Lines 223-243

Asked as: invalid date activity; bad award date; date format last activity; no activity date shown

Turbo reads only the first 10 characters of a date as YYYY-MM-DD. A value shorter than 10 characters, or one that is not a valid date, counts as no date. A project with no valid award date or time entry date has no last activity, so it is never recent.

## Recent project window for the duplicate review

Lines 235-255

Asked as: recent project days; how old is recent; project considered recent; future award date recent

A project is recent when its last activity, the later of its award date and its latest time entry date, is on or after today minus JOB_REQUEST_DUPLICATE_RECENT_DAYS. A future award date counts as recent. A project with no valid award date or time entry is not recent.

## Recent window uses JOB_REQUEST_DUPLICATE_RECENT_DAYS and includes the cutoff day

Lines 245-255

Asked as: recent days setting; cutoff day recent; JOB_REQUEST_DUPLICATE_RECENT_DAYS; how many days recent window

A project is recent when its last activity date is on or after the day JOB_REQUEST_DUPLICATE_RECENT_DAYS days before today (UTC). A project on the cutoff day counts as recent. An award date in the future also counts as recent.

## Open projects in the duplicate review

Lines 257-260

Asked as: closed projects duplicates; cancelled job still shown; open projects only; status hides duplicate

A project is open unless its status is Closed or Cancelled. Only open and recent projects of the same client are listed when they are not exact matches.

## Proposals are never listed as duplicate jobs

Lines 267-270

Asked as: proposal shown as duplicate; P numbers in duplicate list; proposal numbers missing; why proposals not listed

The duplicate-work review skips any job whose number starts with P (a proposal), even when it matches the request exactly.

## Tiers and order of duplicate jobs

Lines 272-299

Asked as: duplicate list ordering; which duplicate shows first; similar versus client tier; newest job first duplicates

A non-exact job of the same client that is open and recent is listed in the similar tier if its description is similar, otherwise in the client tier. Exact matches come first, then similar, then client. Within a tier, the most recent last activity comes first, then the higher job number.

## Client jobs outside the rules are counted, not listed

Lines 272-285

Asked as: hidden client jobs count; jobs not shown same client; how many jobs hidden; old client projects missing

A job of the same client that is not an exact match and is closed, cancelled or not recent is left out of the list. Turbo adds it to the hidden count (hidden_client_jobs). Jobs of other clients that are not exact matches are neither listed nor counted.

## Similar tier needs a similar description

Lines 275-279

Asked as: similar tier meaning; similar but client tier; why not similar tier; description not similar listed

An open, recent project of the same client is in the similar tier when it has the reason "Similar description". Otherwise it is in the client tier.

## Count of hidden client jobs

Lines 280-284

Asked as: hidden jobs number; count of hidden projects; jobs left out of review; how many hidden

The review counts the same client's jobs that are not exact matches and are closed, cancelled or not recent, as hidden_client_jobs. It does not list them. Jobs of other clients are not counted.

## Listed jobs sort by tier, last activity, then number

Lines 290-299

Asked as: job sort order duplicates; latest activity first list; same date job number order; duplicate list sorting

Turbo sorts listed jobs by tier (exact, similar, client). Inside a tier, the job with the latest activity comes first. When the activity dates are equal, the higher job number comes first.

## Which other pending requests are listed and in what order

Lines 303-316

Asked as: other pending requests; requests waiting same client; pending request duplicates; other requests listed

Another pending request is listed if it is an exact match or has the same client, with no date rule. Exact matches come first, then same client with similar description, then other same-client requests. Within a tier, the most recently edited comes first.

## Duplicate review opens only for a visible request

Lines 318-324

Asked as: cannot open duplicate review; duplicate review access denied; request not visible duplicates; permission to see duplicates

The duplicate review loads the request with loadVisibleRequest. Whether a person may see the request depends on that function, which is in another file.

## Duplicate review response has jobs, hidden count and pending requests

Lines 318-369

Asked as: duplicate response fields; jobs hidden requests response; client PO lowercase query; duplicate review data

The duplicate-work review for a project request returns three things: the listed jobs (jobs), the number of the client's projects left out (hidden_client_jobs), and other pending project requests (requests). The jobs come from a database query that Turbo binds to the request's client, proposal, parent, client PO and client reference. Turbo lowercases and trims the client PO and client reference before the query.

## Pending project requests appear only for procurement

Lines 337-353

Asked as: who sees pending requests; procurement pending requests; requests list empty; why no other requests

Only a person with the job claim (procurement) sees other pending project requests in the review. Other people get an empty requests list. The request under review is left out.

## Only users with the job claim see pending requests

Lines 337-342

Asked as: job claim requests; only procurement sees requests; no pending requests shown; requests missing for me

Turbo adds other pending project requests to the duplicate review only for a person who has the job claim (procurement). For anyone else the requests list is empty.

## Pending requests are sorted by tier then newest edit

Lines 347-367

Asked as: pending requests order; newest edited request first; request sort duplicate; pending list tier order

Turbo sorts the pending requests by last edit, newest first. It then sorts them by tier, keeping the edit order inside each tier. Exact matches come first, then similar, then client.

## Pending requests other than the one under review

Lines 351-356

Asked as: own request in duplicates; request lists itself; current request appears duplicate; review shows same request

The review never lists the request under review among the pending requests. Turbo compares every other pending project request with it using the same reasons as for jobs.
