<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/send.go","sha256":"b23ed0dba57c0aef57c2aaa2c9ad75aa83d6137821a57f45ef8a58b9f113294a","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/send.go

## Delivery runs in the background by default

Lines 23-36

Turbo delivers the email on a background process by default (sendNotificationAsync is true unless a test overrides it). The send call returns before delivery ends. A delivery failure does not make the send call return an error. It only sets the notification to 'error'.

## Delivery result sets status to sent or error

Lines 38-61

After the mail send, Turbo sets the notification status to 'sent' and clears the error text. If the mail send fails, it sets the status to 'error' and stores the error message. In both cases it records the time of the status change.

## Only pending notifications are sent; others are skipped silently

Lines 63-96

Sending a notification by its ID acts only on a notification whose status is 'pending'. If the notification is missing or already in another status, Turbo does nothing and reports no error.

## Invalid notification data blocks sending

Lines 98-109

If the notification's data is not valid JSON, Turbo does not send the notification and returns an error. The status stays 'pending'.

## Notification email body is rendered from a template with record data

Lines 111-139

Turbo renders the email body from the notification's template. The template can use the notification fields (recipient name and email, user name, subject, status, error and so on) and any values in the notification's data. A data value with the same name as a notification field replaces it. A placeholder with no value makes the render fail.

## Unresolved legacy placeholder blocks sending

Lines 140-142

If the rendered body still contains a legacy placeholder that was not resolved, Turbo refuses to send the notification and returns an error naming the placeholder. The status stays 'pending'.

## Subject line is rendered with the same data as the body

Lines 144-153

Turbo renders the subject as a template with the same data as the body, so a subject can name the record. A subject without placeholders is sent unchanged. A subject placeholder with no value makes the send fail.

## Email is sent from the app's sender settings to the recipient

Lines 155-160

Turbo sends the email from the sender name and address in the app settings. It sends it to the recipient's profile name and the email on the recipient's user account. The body is plain text.

## Notification moves from pending to inflight before delivery

Lines 162-177

Before it delivers the email, Turbo sets the notification status to 'inflight' and records the time. This happens in the same transaction as the fetch and render. If any step fails, the transaction rolls back and the notification stays 'pending'.

## A crash during email delivery is caught and logged

Lines 180-206

If delivery panics, Turbo recovers and logs the panic. The panic does not crash the app. In that case the status is not updated, and the notification stays 'inflight'.

## Queue send takes one pending notification and reports the remainder

Lines 211-246

Sending the next pending notification counts the notifications with status 'pending'. If there are none, it returns 0 and does nothing. Otherwise it sends one pending notification and returns the count minus one. If sending fails, it returns the error and the count from before the send.

## Draining the queue sends until none are pending

Lines 248-265

Sending all notifications repeats the single send until no pending notifications remain. It stops at the first error and returns that error. The count it returns is the number of loop passes that finished without error. It adds one for the final pass that finds an empty queue.
