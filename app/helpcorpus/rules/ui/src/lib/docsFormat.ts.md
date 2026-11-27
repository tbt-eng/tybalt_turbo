<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"92d574aabe34a399032c40fc576a535001ea3b4ba45147ad5ebb3c996f168de6","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Markdown links in help pages get the link style

Lines 33-39

Asked as: link colour in help; help links look wrong; markdown link style; link class anchor

Turbo renders each Markdown link in a help page as an anchor with the class "link". The link text is parsed as inline Markdown.

## Commit links point to the GitHub repository

Lines 45-49

Asked as: commit link goes nowhere; whats new commit url; link to github commit; broken commit link

A commit link in What's new goes to https://github.com/tbt-eng/tybalt_turbo/commit/<commit>.

## Doc sections come from level-2 headings; commit from a '(to hash)' suffix

Lines 51-57

Asked as: section split by heading; week hash in heading; what's new week id; heading suffix hash

Each level-2 heading starts a new section of a doc page. A heading ending in '(to' plus a hex commit id of 8 to 40 characters')' is a What's new week: the id is stored as the last commit it covers, and the title drops that suffix. Content before the first level-2 heading is the page intro. The section id is the title in lowercase with apostrophes removed and other characters turned into hyphens.

## Doc word counts count only words that start with a letter

Lines 59-61

Asked as: doc word count wrong; numbers counted as words; word count of doc; why words miscounted

Turbo counts a doc's words as runs that start with a letter (letters, digits, apostrophes and hyphens may follow). Numbers on their own do not count.

## Flow blocks need three parts per line

Lines 63-80

Asked as: flow line error; flow step missing field; actor step detail error; flow block parse failure

A code block marked 'flow' in a doc shows as a step-by-step flow. Each line must be 'Actor | Step | Detail', with all three parts filled in. Otherwise parsing fails with 'A flow line needs "Actor | Step | Detail": <line>'. A line starting with 'Note:' becomes a note, and blank lines are skipped.

## Flow blocks and notes in help pages

Lines 82-100

Asked as: flow block not showing; steps in help page; help page code block; extra spaces in help

A code block marked "flow" in a help page becomes a flow block with steps and notes. Turbo renders all other content as HTML. Blank space between tokens is dropped.

## A doc page needs a title

Lines 102-121

Asked as: doc page parse error; missing title error; page has no heading; doc needs # title

A doc page must have a level-1 heading ('# Title'). Without one, parsing fails with the error 'Doc <slug> needs a "# Title" line'. Only the first level-1 heading is the title.

## Help text before the first section is the page intro

Lines 106-120

Asked as: text above first heading; intro text on help; content before section; help page introduction

Turbo puts content that comes before the first level-2 heading into the page intro. Content after a level-2 heading goes into that section.

## Only the first level-1 heading sets the page title

Lines 107-108

Asked as: wrong page title; second h1 heading; two titles on page; title not updating

Turbo takes the first level-1 heading as the page title. A later level-1 heading is treated as ordinary content.

## Help page and section word counts

Lines 126-134

Asked as: word count per section; help page length; how many words page; section word total

Turbo counts words for the whole page source and for each section separately. A section count uses only the raw text of that section's content, not its heading.

## Order of pages in the Help side navigation

Lines 137-160

Asked as: help menu order; sidebar page ordering; why pages in wrong order; proposal before project request

Help lists pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, client-invoices, clients, absorb-proposals, computers, profile, glossary, whats-new (docPageOrder). A proposal page comes before the project request page.

## Child pages appear under their parent with a shorter label

Lines 162-180

Asked as: subpage under parent; nested help page label; short name under time; child page indent

Some Help pages are listed indented under a parent page with a shorter label (childDocPages). Examples: review-time shows as 'Review and track' under time, and close-a-purchase-order shows as 'Close or cancel' under purchase-orders. Children show in the order of the childDocPages list, right after their parent.

## Help pages missing from the order list go before What's new

Lines 182-200

Asked as: page missing from help menu; new help page not listed; unlisted help page position; page appears above whats new

A Help page that is not in docPageOrder and is not a child page goes just before What's new. Such pages are sorted by title among themselves.

## Help page label and indent in navigation

Lines 202-206

Asked as: sidebar label differs; help nav indent; shorter menu name; page title in navigation

A child page shows its short label and is indented in the side navigation. Any other page shows its own title with no indent (docNavEntry).

## Help popover links open in a new tab

Lines 208-221

Asked as: help link opens new tab; popover link closes form; help link same window; unsaved work lost clicking link

Links inside a help popover open in a new tab (target _blank), so a person does not lose unsaved work on the form under the popover. This applies to links in text, flow steps and flow notes.

## Help reference needs both page and section to match

Lines 223-232

Asked as: help link opens nothing; section link not working; hash link empty help; help anchor not found

A help reference such as 'project-requests#approve-a-request' opens a section only when both the page and the section id exist. Otherwise Turbo finds nothing (findSection returns null).
