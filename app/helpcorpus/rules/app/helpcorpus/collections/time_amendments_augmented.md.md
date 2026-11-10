<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_amendments_augmented.md","sha256":"746e79e05e379a14aa3faf6fb72e976e68ed93f9a987fa901d105a36c6c43e87","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_amendments_augmented.md

## Time amendments view is read-only through the API

Lines 5-5

Asked as: cannot edit amendment through api; read only time amendments; delete amendment api not allowed; time amendment view changes

The time_amendments_augmented view has no create, update or delete rules. Nobody can create, change or delete its records through the API. Turbo's own server code is not limited by these rules, and the view reads from the time_amendments collection.

## Who can list time amendment view records

Lines 7-13

Asked as: cannot see time amendments list; no time amendments showing; missing tame claim access; who can list amendments

Only a person with the time amendment claim (the 'tame' claim) or the report claim (the 'report' claim) can list records in the time_amendments_augmented view. Everyone else gets no records.

## Who can view one time amendment view record

Lines 15-21

Asked as: cannot open one amendment; time amendment record access denied; view single amendment permission; amendment details not visible

Only a person with the time amendment claim (the 'tame' claim) or the report claim (the 'report' claim) can view a single record in the time_amendments_augmented view. The rule is copied from the time_amendments collection.

## View adds names, codes and job details to amendments

Lines 26-63

Asked as: amendment missing names; show full name on amendment; amendment blank job or category; time amendment extra details fields

The view shows each time amendment with the full names of the person (uid_name), the creator (creator_name) and the committer (committer_name), built as given name plus surname from profiles. It also adds the time type code and name, the job number and description, the category name, and the division code and name. The joins are left joins, so an amendment with no job, category, division or committer still appears, with those values empty.

## Time amendment hours limits

Lines 73-76

Asked as: how many hours on amendment; negative hours amendment; meals hours maximum; amendment hours too big

Hours (hours) on a time amendment are required and must be between -18 and 18. Negative hours are allowed. Meals hours (meals_hours) can be at most 3.

## Work record number pattern

Lines 78-78

Asked as: work record number format; work record rejected; job reference format wrong; record code must start with F K Q

A work record (work_record) on a time amendment must start with F, K or Q, then two digits, a hyphen and at least three digits, with an optional hyphen and digits at the end (for example K24-123 or Q24-123-1).

## Amendment date formats

Lines 80-87

Asked as: amendment date format; week ending date wrong; date must be YYYY-MM-DD; committed week ending format

The date, week ending (week_ending) and committed week ending (committed_week_ending) are text in the form YYYY-MM-DD. The date and week ending are required.
