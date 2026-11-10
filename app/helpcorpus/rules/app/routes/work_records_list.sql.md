<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/work_records_list.sql","sha256":"f87070f7a94eb6fc13a47784f3df93a673b77afe9995c5f912df1772102de85e","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/work_records_list.sql

## Work records list counts time entries per work record

Lines 1-7

Asked as: count time entries per record; entries per work record; how many entries each record; time entries grouped by record

The work records list groups time entries (time_entries) by work record (work_record). It shows each work record with its number of entries (entry_count).

## Work records list shows first character as prefix

Lines 3-3

Asked as: what is the prefix column; first letter of record; record prefix letter; prefix shown in list

Each row of the work records list shows a prefix (prefix). The prefix is the first character of the work record.

## Work records list skips entries with an empty work record

Lines 6-6

Asked as: missing time entries in list; blank work record entries; entries with no record; empty record hidden from list

The work records list leaves out time entries whose work record is an empty string (work_record != ''). Turbo does not count these entries in any row.

## Work records list sorts by work record ascending

Lines 8-8

Asked as: work records sort order; order of records list; alphabetical work record list; records listed ascending

The work records list is sorted by work record in ascending order (ORDER BY work_record ASC).
