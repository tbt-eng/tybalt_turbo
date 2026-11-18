<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/printWhenReady.ts","sha256":"472b7892f82cbb4c1832aa1ea4bb77db91731bb5f97ca8536b5d87b5894d92a7","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/printWhenReady.ts

## Print dialog opens only after page, fonts and images are ready

Lines 3-14

Asked as: print preview blank; print dialog opens too early; printout half drawn; print waits for page to load; fonts not loaded when printing

Turbo opens the browser print dialog only after the page has rendered, its fonts have loaded and its images have decoded. The dialog therefore does not show a half-drawn page.

## Image that fails to decode does not block printing

Lines 12-12

Asked as: image missing from printout; picture not printing; broken image when printing; print cut off image; image blank on printed page

When an image fails to decode, Turbo ignores the error and still opens the print dialog once the other images are ready. That image may print incomplete.

## Leaving the page before it is ready cancels the print

Lines 13-17

Asked as: print dialog did not open; closed page then print gone; cancel printing after navigating away; print vanished after leaving; printing stopped when I left

When a person leaves the page before the page, fonts and images are ready, Turbo does not open the print dialog. The cancel function that printWhenReady returns sets this.
