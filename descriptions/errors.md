# ErrorSystem

## MessageBar component

The message bar is the app's message area. It displays the error messages and,
after them, any flash message (see below). Next to each message is a button that
dismisses it. If there are no messages, nothing is shown. It is placed in the
layout so it is visible on every page, and errors persist across navigation. It
is at ui/src/lib/components/MessageBar.svelte, has no props of its own, and is
included in ui/src/routes/+layout.svelte as the first element of the main
content area.

MessageBar only decides which messages to show. Each one is drawn by
ui/src/lib/components/DSMessageBar.svelte: a full-width strip at the top of the
page with no border radius, a bottom border, a tone icon, and a tone colour
(red for errors, green for success, blue for info). The split follows the UI
conventions in ui/CLAUDE.md: DS components own the styling and hold no app
state.

## Flash messages

A flash message is a one-line message for the page navigated to, such as
"Approved X" after approving an expense moves the user on to the next one. A
page passes it with the navigation: `goto(url, { state: { flash } })`, where
`flash` is a `Flash` from ui/src/lib/flash.ts (a tone, the message, and an
optional link). It lives in that history entry, so Back and Forward keep it and
a reload or fresh visit drops it. Dismissing it clears it from page state.

## errorMessages in global store

The global store is at ui/src/lib/stores/global.ts. The errorMessages array is
used to store the error messages that are displayed in the MessageBar component.
It is an array of objects with the following properties:

- `message`: The error message to display.
- `id`: The id of the error message.

The id is used to identify the error message when it is dismissed.

The errorMessages array is a writable store, so it can be updated by any
component. The global store is the only component that should update the
errorMessages array.

### Adding errors

The `addError` function is a function of the global store and is used to add an
error to the errorMessages array. It takes a single string argument representing
the error message. A unique id is generated for the error message and added to
the array.

### Dismissing errors

The `dismissError` function is a function of the global store and is used to
dismiss an error from the errorMessages array. It takes a single string argument
representing the id of the error to dismiss. The error is removed from the
array.
