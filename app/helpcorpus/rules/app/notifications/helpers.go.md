<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/helpers.go","sha256":"69327da37d9cd596faa801c7fa7b2315c028c095584160376544d58a417c455b","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/helpers.go

## Action links use the app URL as their base

Lines 16-30

Turbo builds a notification action link (BuildActionURL) from the app URL setting (Meta.AppURL). It trims spaces and trailing slashes from that URL. When the app URL is empty, the link is an empty string.

## Empty action path gives the bare app URL

Lines 31-34

When the path for a notification action link is empty or only spaces, Turbo returns the base app URL alone (BuildActionURL).

## Action link path gets a leading slash

Lines 35-38

When the path for a notification action link has no leading slash, Turbo adds one before it joins the path to the app URL (BuildActionURL).

## Old placeholders are detected in notification text

Lines 41-49

Turbo detects two unresolved legacy placeholders in notification text: {APP_URL} and {:RECORD_ID}. It checks {APP_URL} first and reports only one placeholder. It reports an empty value when neither is present (unresolvedLegacyPlaceholder).

## Notification status time is set on creation or status change

Lines 57-62

Turbo sets the notification's status_updated field to the current time when the notification is new or when its status differs from the saved status. A save that leaves the status unchanged keeps the old status_updated (WriteStatusUpdated).

## Queue-only callers do not send notifications

Lines 64-78

When the send flag is false, Turbo only queues the notifications and sends nothing (sendQueuedIfRequested). When the flag is true, it sends the queued notifications. If sending fails, it returns the error "error sending notifications: ...".

## Profile display name joins given name and surname

Lines 80-89

Turbo finds a profile by its user (uid) and builds the display name as given name, a space, and surname, with outer spaces trimmed. When a profile has no surname or no given name, the name has no extra space. When no profile exists for the user, the lookup returns the error.

## Immediate notifications go to each recipient separately

Lines 91-128

Turbo sends an immediate notification to each recipient in the list, one at a time (createAndSendToRecipients). When one recipient fails, Turbo logs "error creating immediate notification" and goes on to the next. A recipient whose dispatch returns an empty notification ID is skipped and not counted. The function returns the number of notifications created.

## Missing or null row values read as empty text

Lines 130-136

When a row has no value for a key, or the value is null, Turbo reads it as an empty string (rowStringValue).
