<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/jobs_latest.sql","sha256":"e989a7f5fd181f9af3df600d5494a6fca48803fca894c2ab1a7fdd4d2d13061e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/jobs_latest.sql

## Missing outstanding balance shows as 0 or blank date

Lines 13-14

Asked as: outstanding balance shows 0; blank outstanding date; job balance zero; no outstanding date

A job with no outstanding balance shows 0. A job with no outstanding balance date shows an empty date.

## Latest jobs list splits proposals and projects

Lines 18-57

Asked as: proposals and projects split; job number starts with P; proposal vs project list; why is job in proposals

The latest jobs list has two groups: 'Proposals' and 'Projects'. A job is a proposal when its number starts with P (number LIKE 'P%'). Every other job is a project.

## Latest jobs list caps each group at the limit

Lines 18-43

Asked as: latest jobs limit; too many jobs listed; limit per group; list longer than limit

Turbo takes up to the requested limit (:limit) of the newest jobs from each group separately. The limit is not shared between the groups. The list can therefore hold up to twice the limit.

## Missing client alias, branch or manager shows blank

Lines 52-63

Asked as: job missing client alias; blank branch on job; manager name empty; job hidden missing client

A job with no client alias, branch or manager shows an empty value for that field. The branch shows as its code and the manager as given name and surname. The job still appears in the list.

## Latest jobs list order: proposals first, newest first

Lines 57-64

Asked as: latest jobs order; newest jobs first; proposals listed first; sort latest jobs

The list shows proposals before projects. Within each group the jobs are ordered by creation date, newest first (created DESC).
