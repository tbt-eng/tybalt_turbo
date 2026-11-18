<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/clipboard.ts","sha256":"5aada6bfb56b22aa49fb16595045ab0036461f10d5d6eb2372e8ef08bd15d64b","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/clipboard.ts

## Copy button reports failure when browser has no clipboard

Lines 9-9

Asked as: copy button does nothing; copy not working in browser; clipboard not available; copy fails on http page; old browser copy button

Turbo's copy function (copyText, used by DSCopyButton) returns false when the browser has no clipboard writeText function. This happens on an insecure page or in an old browser. Turbo does not copy anything in that case.

## Copy button reports failure when browser refuses the copy

Lines 10-15

Asked as: copy button says failed; browser blocked copying; copy refused by browser; copy returns false; clipboard permission denied

When the browser refuses to write the text to the clipboard, Turbo's copy function (copyText) returns false and does not raise an error. It returns true only after the browser has accepted the text.
