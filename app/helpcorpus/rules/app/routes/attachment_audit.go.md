<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/attachment_audit.go","sha256":"e057bc2c8f2cff00bf1a484d1132c114d0ba9b56fd23fd2d6ef9a775744c5f55","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/attachment_audit.go

## Audit run status values

Lines 24-29

An attachment audit run has one of three statuses: "running", "completed" or "failed".

## Delete-orphans result counts each file by outcome

Lines 62-77

The orphan cleanup response counts deleted files (deleted_files), files skipped because a record now references them (skipped_referenced_files), files already gone from storage (already_missing_files), rows skipped as invalid (skipped_invalid_files) and failures (failed_files). It also lists each failure with its storage path and error.

## Attachment audit covers two targets only

Lines 112-118

Turbo can audit two attachment targets. They are Expense Documents (expense_documents, field attachment) and Purchase Orders (purchase_orders, field attachment). Any other target key gets the message "attachment audit target not found".

## Only admins can use attachment audits

Lines 120-129

Turbo requires the admin claim for every attachment audit action. These are listing targets, viewing a run's status, refreshing, downloading reports and deleting orphans. Anyone else gets "you do not have permission to" followed by the action, for example "view attachment audits" or "delete orphaned attachment files".

## Listing audit targets shows latest run

Lines 131-152

An admin who lists the audit targets sees each target with its latest audit run. A target that has never been audited shows no latest run.

## Status of never-audited target is blank

Lines 154-175

When an admin asks for the status of a target that has never been audited, Turbo returns the target key with an empty status. It does not return an error. An unknown target returns "attachment audit target not found".

## Refresh starts audit in background

Lines 177-216

Refreshing an audit sets the run to "running", records the requesting admin and the start time, and clears the finish time and error. The audit then runs in the background and the admin receives the running record straight away.

## Refresh while target already running

Lines 188-197

If an audit for the same target is already running, a refresh does not start a second one. Turbo returns the current run record, or a running placeholder if no record exists yet.

## Failed audit is recorded as failed

Lines 206-213

If a background audit fails, Turbo sets the run status to "failed", stores the error message and sets the finish time.

## Report download rules

Lines 219-260

An admin can download the missing or orphaned report as a CSV file named target key plus _missing.csv or _orphaned.csv. If the target has no run, the report is empty, or the file is not in storage, Turbo says "attachment audit report not found".

## Delete orphans needs idle target and report

Lines 262-288

Deleting orphaned files is refused with "attachment audit target is already running" while an audit or cleanup runs for that target. It returns "attachment audit orphaned report not found" if there is no run or no orphaned report.

## Unknown audit target is refused

Lines 290-297

When the target in a status, refresh, download or delete-orphans request is not one of the known target keys, Turbo refuses with "attachment audit target not found".

## Audit run record is created on first use for a target

Lines 307-327

When a target has no saved audit run, Turbo creates one with the target's key, label, collection name and field name, and sets its status to "running" (getOrCreateAttachmentAuditRun). Later runs reuse the same single record per target.

## Starting an audit resets the run record

Lines 329-347

When an admin starts an audit (refresh), Turbo sets the target's run record to status "running". It records the admin as the requester (requested_by) and sets the start time (started_at). It clears the finish time (finished_at) and any error. Each target keeps one run record, which later runs reuse.

## Failed audit keeps error text and finish time

Lines 349-362

When an audit fails, Turbo sets the run's status to "failed", stores the error message (error) and sets the finish time (finished_at). The earlier counts and reports stay on the record.

## Completed audit saves counts and reports

Lines 364-406

When an audit finishes, Turbo sets the status to "completed" and saves the finish time. It also saves the total, referenced, matching, missing and orphaned counts, and clears any error. It attaches two CSV reports, one for missing files and one for orphaned files.

## Audit reports are CSV files with fixed columns

Lines 370-385

Each completed audit saves two CSV reports. The missing report (<target key>_missing.csv) has columns collection, field, record_id, filename, storage_path. The orphaned report (<target key>_orphaned.csv) has columns collection, field, storage_path, record_id_from_path, filename. A report with no rows holds only the header line.

## Orphan deletion continues after failures

Lines 457-482

A storage error on one file is listed in the failures with its path and does not stop the other deletions. The response counts deleted, referenced-skipped, already-missing, invalid-skipped and failed files.

## Audit re-runs after orphan cleanup

Lines 484-495

After deleting orphans, Turbo runs a fresh audit so the counts and reports match storage. If that refresh fails, the response still succeeds and carries a refresh error with no latest run.

## Delete orphans uses cached report as candidates

Lines 498-546

Deleting orphans works from the latest cached orphaned report. A report row without 5 columns counts as invalid. Turbo rechecks each row before it deletes the file.

## Orphan report rows with fewer than 5 columns are invalid

Lines 516-545

When Turbo reads the cached orphaned report, it skips the header row. A row with fewer than 5 columns becomes an empty row. That row is later counted as skipped invalid and is not deleted. A completely empty report file gives no rows.

## Orphan rows revalidated before deletion

Lines 548-586

A row is skipped as invalid if its collection or field differs from the target, its path is blank, the path is outside the collection's folder, or no record id or filename can be read from it. A file is skipped as referenced if its record now holds the same filename. A file already gone from storage counts as already missing.

## Orphan row is invalid if it does not match the target

Lines 549-565

Before deleting a file from the orphaned report, Turbo checks the row. The row's collection and field must equal the audit target's, and its storage path must not be empty. The path must also start with the target collection's file folder. Otherwise Turbo counts the row as skipped invalid and does not delete the file. A path that gives no record id or no filename is also invalid.

## Orphan still referenced by its record is not deleted

Lines 567-577

Before deleting an orphaned file, Turbo looks up the record whose id is in the file's path. If that record's attachment field (trimmed) still holds the same filename, Turbo skips the delete and counts the file as skipped referenced. If the record does not exist, Turbo goes on to delete the file.

## Orphan already gone from storage counts as already missing

Lines 579-585

When Turbo deletes an orphaned file and storage says it is not found, Turbo counts it as already missing. Any other delete error counts as a failed file and is reported with the file's path.

## Audit fails if field missing

Lines 593-595

The audit fails with "field collection.field not found" if the collection does not have the audited field.

## Missing and matching file counting

Lines 597-655

The audit checks every record in the target collection. A record counts as referenced when its attachment field holds a filename after trimming spaces. Blank fields are skipped. A referenced record is matching if the file exists in storage at the collection path, record id and filename. Otherwise it is missing.

## Audit scans every record and ignores thumbnails

Lines 597-626

An audit loads all records of the target collection and lists the stored files under the collection's file folder. It ignores any stored file whose path contains "/thumbs_", so thumbnails never count as orphans. The total records count (total_records) includes records without an attachment.

## Orphaned files definition

Lines 617-678

A stored file under the collection's storage folder that no record references is orphaned. Thumbnail files (paths containing "/thumbs_") are ignored. Both reports are sorted by storage path.

## Records with blank attachment are not referenced

Lines 628-632

The audit trims the attachment field of each record. A record with an empty value is skipped. It counts as neither referenced, matching nor missing.

## Audit report rows are sorted by storage path

Lines 673-678

Both the missing and the orphaned rows are sorted by storage path in ascending order, so the CSV reports have a stable order.

## Orphan path parsing fallback

Lines 683-690

The record id and filename of an orphaned file come from its storage path. If the path has no record folder, the record id is blank and the filename is the last part of the path.

## Date fields in audit responses are UTC or blank

Lines 728-766

Audit start and finish times are shown in UTC in RFC3339 format. An unset time shows as an empty string. The response flags that a missing report and an orphaned report exist (has_missing_report, has_orphaned_report) when a report file is stored.
