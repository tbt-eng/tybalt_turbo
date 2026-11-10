<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docs.ts","sha256":"b7ab75a6491899d2093a428b3f67aa1fa6c2676807e1c1e482ae4c6be95e4173","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docs.ts

## README.md is not a Help page

Lines 14-18

Asked as: readme showing in help; readme file in help search; readme appears as help page; hide readme from help; help page list shows readme

Turbo builds its Help pages from the Markdown files in src/docs, but leaves out src/docs/README.md. The README does not appear in the Help pages or in Help search.

## Help popover shows links that open in a new tab

Lines 36-40

Asked as: help links open same tab; help popover link opens new window; popover links not opening in new tab; help popover empty; help popover shows nothing

A help popover shows the section that matches its reference, and runs its content through withNewTabLinks so its links open in a new tab. When no section matches the reference, the popover has no section to show (null).
