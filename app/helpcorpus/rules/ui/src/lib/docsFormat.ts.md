<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"9c593c10d9f1ecd2c5afb61e743672b32f17c9f7669be96167bf372fe88c044a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages get the link style

Lines 33-39

Asked as: link style in help; help links look wrong; anchor class link; link colour in docs

Turbo renders every Markdown link in a help page as an anchor with the class "link". Its text comes from the parsed link text and its address from the link's href.

## What's new heading carries its last commit

Lines 45-49

Asked as: commit shown on heading; what's new commit link; commit hash in heading; release commit in whats new

A level-2 heading that ends in '(to ' plus 8 to 40 hex characters and ')' sets the section's last commit (commit). The commit text is removed from the section title. Other headings have no commit. The commit link points to the tbt-eng/tybalt_turbo repository on GitHub.

## Slugs for section ids drop apostrophes and symbols

Lines 51-57

Asked as: section anchor id; heading id wrong; apostrophe in anchor; url fragment for heading

Turbo builds a section id (slugify) by lowercasing the heading text and removing apostrophes (' and ’). It then turns each run of characters other than a-z and 0-9 into one hyphen and trims a hyphen at the start or end. For example, "What's new" gives "whats-new".

## Help word count

Lines 59-61

Asked as: help word count; how many words in help; count words per section; words on help page

Turbo counts a word as a run that starts with a letter and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone are not words. Each section and the whole page have their own count.

## Flow block line format

Lines 63-80

Asked as: flow line error; flow format wrong; actor step detail format; flow line needs three parts

A code block marked 'flow' becomes a step flow. Each line must read 'Actor | Step | Detail' with exactly three non-empty parts, or the parser throws 'A flow line needs "Actor | Step | Detail": <line>'. A line that starts with 'Note:' becomes a note, and blank lines are skipped.

## A flow block turns its lines into steps and notes

Lines 63-99

Asked as: flow block shows as text; steps not appearing in help; code block flow render; help steps rendered wrong

A code block marked "flow" in a help page becomes a flow of steps and notes, not ordinary text. Turbo ends the HTML text before it and starts a new text block after it. Turbo also renders Markdown inline in each step's detail and in each note.

## Blank space between help blocks is dropped

Lines 82-100

Asked as: extra space between help blocks; empty gap in help page; blank lines in help; stray whitespace between blocks

When Turbo splits a help page into blocks, it drops the blank-space tokens. Other tokens between flow blocks are grouped into one HTML block.

## Level-2 headings split a help page into sections

Lines 102-135

Asked as: help sections split; second heading new section; section per heading in help; help page section breaks

Each level-2 heading starts a new section of a help page. Content before the first level-2 heading is the page intro. Each section id is its title in lowercase with apostrophes removed and other runs of characters turned into single hyphens.

## Help page needs a level-1 title

Lines 106-121

Asked as: help page missing title; doc needs title line; help page won't load; no heading in help doc

Turbo's docs parser refuses a help page that has no level-1 heading ('# Title'). The error is 'Doc <slug> needs a "# Title" line'. Only the first level-1 heading becomes the page title.

## Help page order and unlisted pages

Lines 137-186

Asked as: help menu order; help page sorting; where page appears in help; help list order wrong

Help lists pages in the fixed order: start-here, time, purchase-orders, expenses, proposals, project-requests, client-invoices, profile, glossary, whats-new. A page not on the list goes just before What's new, sorted by title. Pages with the same position are also sorted by title.

## Child pages sit under their parent in Help

Lines 158-192

Asked as: indented help page under parent; sub page in help menu; nested help topics; help child page label

Five pages show under a parent with a shorter label and an indent: review-time under time ('Review and track'), revise-a-purchase-order ('Revisions') and find-a-purchase-order ('Find and print') under purchase-orders, commit-expenses under expenses ('Review and commit'), and admin-profiles under profile ('Admin profiles'). Children follow the order of that list.

## Help popover links open in a new tab

Lines 194-207

Asked as: help link opens new tab; help link losing form; popover link new window; help link closes form

Links in a help popover open in a new tab (target _blank, rel noopener), so a person does not lose unsaved work on the form. This applies to links in text blocks, flow step details and flow notes.

## Help reference needs a page and a section

Lines 209-218

Asked as: help link to section not working; broken help anchor; help reference shows nothing; link to help heading fails

A help reference such as 'project-requests#approve-a-request' finds a section only when both the page slug and the section id exist. If either is missing, Turbo finds nothing.
