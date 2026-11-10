<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/functions.go","sha256":"426dd20a2233b1de77f0147b6787a4a114a0085448d8452b39b92c28a2aaf481","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/functions.go

## Empty report produces an empty CSV

Lines 22-24

Asked as: empty report csv; report with no rows export; blank csv output; csv missing header row

When a report has no rows, the CSV export (convertToCSV) returns an empty string. It writes no header row.

## CSV columns default to sorted union of all row keys

Lines 26-39

Asked as: csv column order; csv headers sorted; columns in csv export; csv column names wrong order

When no headers are given, the CSV export collects every column name that appears in any row. It sorts the names alphabetically and uses them as the header row, so the column order is the same each time.

## Empty or missing CSV values are blank

Lines 50-58

Asked as: null values in csv; blank cells csv; empty cell export; missing column csv

In a CSV export, a NULL value or a column that a row lacks is written as an empty cell. Only values that are not NULL are written with their text.

## Zip export refuses when there are no attachments

Lines 82-84

Asked as: zip no attachments error; download zip nothing attached; no attachments to zip; zip export fails empty

When a report has no attachments, the attachment zip (zipAttachments) fails with the error "no attachments to zip". Turbo creates no zip.

## Attachment source collection falls back to the report's collection

Lines 98-102

Asked as: attachment file not found; attachment collection missing; zip uses report collection; attachment source path wrong

When an attachment names no collection of its own, Turbo reads its file from the collection of the report. The file path is the collection id, a slash, and the attachment's source path. This applies to the zip and to its manifest.

## Each attachment is stored in the zip under its zip filename

Lines 103-114

Asked as: file names inside zip; zip file name wrong; attachment name in zip; zip entries named oddly

Turbo names each file in the zip after the attachment's zip filename (ZipFilename), not its stored source path.

## One missing or unreadable file fails the whole zip

Lines 105-136

Asked as: zip failed file error; attachment unreadable zip; partial zip download; zip broke on one file

If Turbo cannot read, add or copy any attachment file, or close it after copying, the zip fails with an error that names the file. Turbo produces no partial zip and saves no cache record.

## Zip file is named from class and key

Lines 150-153

Asked as: zip file name format; name of downloaded zip; class key zip name; what is zip called

The finished zip file is named class_key.zip, built from the class and the key of the cache.

## Zip cache record is created or updated by key and class

Lines 155-188

Asked as: zip cache not saved; cached zip overwritten; zip_cache record duplicate; zip cache entry update

After zipping, Turbo looks in the zip_cache collection for a record with the same key and class. If none exists, it creates one with the key, class, manifest and zip. If one exists, it replaces that record's zip and manifest. Either way it returns the record.

## Zip manifest lists each file in a fixed order

Lines 198-225

Asked as: manifest order wrong; zip manifest contents; manifest file listing; sha256 in manifest

The manifest lists the collection id, source path, zip filename and sha256 of each attachment. A missing collection id is replaced by the default collection. Entries are sorted by collection id, then source path, then zip filename, then sha256.

## Cache with no stored manifest counts as a miss

Lines 227-231

Asked as: cache always rebuilds zip; empty zip manifest cache; zip rebuilt every time; cached zip ignored

If the cached zip record has no manifest, or the manifest is empty or null, the cache lookup (zipCacheLookup) returns no record. Turbo then rebuilds the zip.

## Unreadable stored manifest causes an error

Lines 233-237

Asked as: manifest unmarshal error; corrupt zip cache manifest; failed to unmarshal manifest; broken cached manifest

If the manifest of a cached zip cannot be parsed, the cache lookup fails with the error "failed to unmarshal manifest". It does not count as a cache miss.

## Cache hit needs exact manifest match

Lines 250-273

Asked as: zip cache not reused; zip cache miss attachments changed; stale cached zip; zip cache mismatch

The cache lookup returns the cached zip only when a record with the same key and class exists and its manifest equals the manifest of the current attachments. Any difference gives a cache miss. A failed lookup of the record is also a cache miss, and not an error.
