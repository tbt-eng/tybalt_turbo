<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/work_records.go","sha256":"a8e61ebc54acaaac3ccea3b3e1952e89e7a21db050a1c7f57205c17add8d2e12","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/work_records.go

## Who can view work records

Lines 41-68

Staff with the report claim or the work_record claim can view work records. Everyone else is refused with the message "you are not authorized to view work records". This applies to both the work record list and the work record details.

## Work record list requires viewer access

Lines 70-87

The work record list is returned only after the viewer check passes. Each row gets a search text (search_text) that Turbo builds from the work record number.

## Work record details need a work record number

Lines 89-109

The details of a work record are returned for the work record number given in the request. Turbo trims spaces from the number first. If the number is empty, Turbo refuses with "workRecord is required". The viewer check runs before this check.

## Work record search text variants

Lines 111-137

Turbo builds the search text from the work record number in upper case, with spaces trimmed. It adds the number as is and the number without dashes. It then treats the first character as a prefix and adds the number without that prefix, with and without dashes. An empty number gives an empty search text.

## Search text splits on first dash

Lines 118-147

After the prefix is removed, Turbo splits the rest at the first dash. If there are two parts, it adds each part and the two parts joined. The search text is these tokens separated by spaces. Turbo leaves out duplicate tokens.

## Single-character work record has no prefix variants

Lines 131-136

A work record number of only one character gets no prefix-stripped variants in its search text. Turbo adds only the upper-case number and the version without dashes.
