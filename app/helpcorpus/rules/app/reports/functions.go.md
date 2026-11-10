<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/reports/functions.go","sha256":"426dd20a2233b1de77f0147b6787a4a114a0085448d8452b39b92c28a2aaf481","model":"claude-sonnet-5-5"} -->

# Rules in app/reports/functions.go

## Empty report produces an empty CSV

Lines 22-24

When a report has no rows, the CSV export (convertToCSV) returns an empty string. It writes no header row.

## CSV columns default to sorted union of all row keys

Lines 26-39

When no headers are given, the CSV export collects every column name that appears in any row. It sorts the names alphabetically and uses them as the header row, so the column order is the same each time.

## Empty or missing CSV values are blank

Lines 50-58

In a CSV export, a NULL value or a column that a row lacks is written as an empty cell. Only values that are not NULL are written with their text.

## Zip export refuses when there are no attachments

Lines 82-84

When a report has no attachments, the attachment zip (zipAttachments) fails with the error "no attachments to zip". Turbo creates no zip.

## Attachment source collection falls back to the report's collection

Lines 98-102

When an attachment names no collection of its own, Turbo reads its file from the collection of the report. The file path is the collection id, a slash, and the attachment's source path. This applies to the zip and to its manifest.

## Each attachment is stored in the zip under its zip filename

Lines 103-114

Turbo names each file in the zip after the attachment's zip filename (ZipFilename), not its stored source path.

## One missing or unreadable file fails the whole zip

Lines 105-136

If Turbo cannot read, add or copy any attachment file, or close it after copying, the zip fails with an error that names the file. Turbo produces no partial zip and saves no cache record.

## Zip file is named from class and key

Lines 150-153

The finished zip file is named class_key.zip, built from the class and the key of the cache.

## Zip cache record is created or updated by key and class

Lines 155-188

After zipping, Turbo looks in the zip_cache collection for a record with the same key and class. If none exists, it creates one with the key, class, manifest and zip. If one exists, it replaces that record's zip and manifest. Either way it returns the record.

## Zip manifest lists each file in a fixed order

Lines 198-225

The manifest lists the collection id, source path, zip filename and sha256 of each attachment. A missing collection id is replaced by the default collection. Entries are sorted by collection id, then source path, then zip filename, then sha256.

## Cache with no stored manifest counts as a miss

Lines 227-231

If the cached zip record has no manifest, or the manifest is empty or null, the cache lookup (zipCacheLookup) returns no record. Turbo then rebuilds the zip.

## Unreadable stored manifest causes an error

Lines 233-237

If the manifest of a cached zip cannot be parsed, the cache lookup fails with the error "failed to unmarshal manifest". It does not count as a cache miss.

## Cache hit needs exact manifest match

Lines 250-273

The cache lookup returns the cached zip only when a record with the same key and class exists and its manifest equals the manifest of the current attachments. Any difference gives a cache miss. A failed lookup of the record is also a cache miss, and not an error.
