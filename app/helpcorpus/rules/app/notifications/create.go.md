<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/notifications/create.go","sha256":"5f9b07c9e52cf6b501ccef5ec1bf00cd69ea0744227df005df02fcc31c66dae8","model":"claude-sonnet-5-5"} -->

# Rules in app/notifications/create.go

## Notification delivery mode must be deferred or immediate

Lines 28-30

Turbo refuses to dispatch a notification when the delivery mode is not deferred (DeliveryDeferred) or immediate (DeliveryImmediate). The error is 'invalid delivery mode' followed by the mode. No notification record is created.

## Deferred notifications are created as pending and not sent

Lines 37-49

In deferred mode (DeliveryDeferred), Turbo creates the notification record with status 'pending' and does not send it at that moment. Turbo returns the new notification ID.

## Immediate notifications are sent right after creation

Lines 37-39

In immediate mode (DeliveryImmediate), Turbo creates the notification and then tries to send it (SendNotificationByID). Turbo sends only when a notification was actually created.

## Immediate send failure does not block the action

Lines 37-48

In immediate mode, when the send fails, Turbo logs the error and does not return it. The caller still gets the notification ID and no error, so the business action goes on. The notification record stays as created.

## Disabled notification feature skips creation silently

Lines 52-68

When the notifications feature flag for the template code is disabled, Turbo creates no notification. It returns an empty ID and no error, and logs that creation was skipped. In immediate mode nothing is sent.

## Feature config read failure skips notification (fail-closed)

Lines 52-60

When Turbo cannot read the notifications feature configuration, it creates no notification. It treats this as disabled (fail-closed), returns an empty ID with no error, and logs the error.

## Notification collection must exist

Lines 70-77

When Turbo cannot find the notifications collection, it fails with the error 'error finding notifications collection' and creates nothing.

## Missing notification template refuses creation

Lines 79-89

When no notification template has a code equal to the template code, Turbo does not create the notification. It returns the error 'error finding notification template' with the code and the cause.

## New notification copies template content and starts pending

Lines 91-97

Turbo sets the recipient to the recipient user, and copies the template reference, subject and email text (text_email) from the template. It sets the status to 'pending' and records whether it is a system notification (system_notification).

## Actor is recorded only when given

Lines 98-100

Turbo sets the notification's user field to the actor only when an actor ID is given. Otherwise it leaves the field unset.

## Notification data is stored as JSON only when not empty

Lines 102-112

When the notification has data, Turbo stores it as a JSON string in the data field. When the data is empty, Turbo leaves the field unset. A data value that cannot be converted to JSON fails with 'error marshaling notification data'.

## Notification save failure returns an error

Lines 114-123

When Turbo cannot save the notification record, it returns the error 'error saving notification' with the cause and no ID.
