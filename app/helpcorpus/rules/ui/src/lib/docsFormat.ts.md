<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"d2466d4fa4a54edb4bce3a938c7b9af9345d9fa4bd21c4ed7b5b6d9f3521fd99","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links render as styled links

Lines 33-43

Asked as: help link looks like plain text; link not clickable in help; help doc links style

Links in help docs render as an anchor with the 'link' class. Links in the details of flow steps and in notes use the same inline formatting.

## What's new headings link to their last commit

Lines 45-49

Asked as: what's new heading commit link; commit hash in what's new; week's last commit

A What's new section heading that ends with "(to <commit hash>)" of 8 to 40 hex characters shows the title without that suffix. Turbo keeps the hash as the last commit that the week covers, and builds a GitHub commit link from it.

## Help section links need a page and a section

Lines 51-57

Asked as: help link to section not opening; section link does nothing; broken help section anchor

A help reference such as "project-requests#approve-a-request" opens only when both the page and the section id exist. Turbo finds nothing if either is missing (findSection). A section id is its title in lower case, with apostrophes removed and other characters turned into hyphens.

## Help word counts count only words that start with a letter

Lines 59-61

Asked as: help word count wrong; how words are counted in help; numbers not counted as words

Turbo counts a word as a run that starts with a letter and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone do not count. A section's count uses its raw text. A page's count uses the whole source.

## Flow blocks need three parts per line

Lines 64-78

Asked as: flow block error; flow line needs three parts; help flow step format error

In a help doc, a code block marked 'flow' becomes a step-by-step flow. Each line must read 'Actor | Step | Detail' with all three parts filled in. A line that does not fails with the error 'A flow line needs "Actor | Step | Detail": <line>'. Blank lines are skipped.

## Flow lines starting with Note: become notes

Lines 69-72

Asked as: note in help flow; flow note text; note: line in help steps

A flow line that starts with 'Note:' is a note, not a step. Turbo trims the 'Note:' prefix and shows the rest as text under the flow.

## Help doc page needs a level-1 title

Lines 103-121

Asked as: help doc fails to load; missing title in help doc; doc needs title line

When Turbo parses a help doc, it requires a '# Title' line (a level-1 heading). Only the first level-1 heading becomes the page title. A doc without one fails with the error 'Doc <slug> needs a "# Title" line'.

## What's new heading carries its last commit hash

Lines 109-116

Asked as: what's new heading hash; commit suffix missing from title; what's new title shows hash

A level-2 heading that ends with '(to <hash>)' sets the section's last commit (commit) and drops that suffix from the title. The hash must be 8 to 40 lowercase hex characters. A heading in any other form has no commit and keeps its full text.

## Help page order in the side navigation

Lines 137-197

Asked as: help page order wrong; side navigation help order; help pages sorted incorrectly

Help lists pages in this order: Start here, Time, Purchase Orders, Expenses, Proposals, Project Requests, Rate Sheets, Clients, Absorb Proposals, Computers, Profile, Glossary, What's new (docPageOrder). A page that is not in this list appears just before What's new, sorted by title.

## Child pages appear indented under a parent with a short label

Lines 161-203

Asked as: help subpage indented; nested help page under parent; child help page not showing under

Some Help pages appear indented under a parent page with a shorter label. Examples: Review and track under Time; Revisions, Find and print, Child POs and Close or cancel under Purchase Orders; Review and commit under Expenses; Proposed rates under Project Requests; Admin profiles under Profile; Periods, Duplicates, Codes and names and Delete under Rate Sheets; Suggestions under Absorb Proposals. A child page follows its parent, in the order of the list (childDocPages).

## Pages missing from the order list go before What's new

Lines 179-197

Asked as: help page missing from menu; new help page in wrong place; unlisted help page position

A page that is neither in the page order list nor a child page is placed just before What's new. Such pages sort by title among themselves.

## Links in help popovers open in a new tab

Lines 205-218

Asked as: help link opens over form; help popover link new window; links closing my unsaved form

A help popover that sits over a form opens its links in a new tab (target=_blank), so that a person does not lose unsaved work. This applies to links in text blocks, in flow steps and in flow notes (withNewTabLinks).

## Help section lookup returns nothing for unknown pages or sections

Lines 220-229

Asked as: help section not found; help link goes nowhere; help reference returns nothing

A help reference such as 'project-requests#approve-a-request' needs both a page slug and a section id that exist. If the page or the section is missing, or the reference has no '#', the lookup finds nothing.
