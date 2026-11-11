<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"ae552abf87c378a285d6ff88b03a0a4cb6eb5c71503108455e667c908d4ba6af","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links open in the same tab by default

Lines 33-39

Asked as: link opens same tab; help link new window; links open where; clicking link leaves page

The Markdown renderer turns each link in a help page into a link with the class "link" and the given address. It adds no target, so the link opens in the same tab. The popover rule changes this for popovers.

## Commit links point to the tybalt_turbo GitHub repository

Lines 45-49

Asked as: what's new commit link; github commit link broken; link to commit on github; commit hash link

The link for a What's new commit is the repository address https://github.com/tbt-eng/tybalt_turbo/commit followed by a slash and the commit hash (commitUrl).

## Section ids and word counts in docs

Lines 51-61

Asked as: section anchor id; heading id slug; section link anchor; how ids are made

Each doc section gets an id made from its title: lowercase, apostrophes removed, other runs of non-alphanumeric characters turned into one hyphen. Words are counted as sequences that start with a letter. Content before the first level-2 heading is the page intro.

## Flow line format and errors in docs

Lines 63-80

Asked as: flow line error; flow step wrong format; flow line needs three parts; doc line fails parsing

A flow line in a doc must read "Actor | Step | Detail", with exactly three non-empty parts. A line that starts with "Note:" is a note instead. Any other line fails with the error: A flow line needs "Actor | Step | Detail": followed by the line. Blank lines are skipped.

## Flow blocks become step lists and other blocks stay as HTML

Lines 82-100

Asked as: flow block not showing; steps shown as text; code block flow rendering; help steps out of order

In a help page, a code block marked "flow" becomes a flow block of steps and notes (parseFlow). All other content becomes HTML blocks. Blank lines (space tokens) are dropped. A flow block ends the HTML block before it, so content keeps its order.

## Doc page needs a title

Lines 102-121

Asked as: doc page no title; missing page heading error; doc fails without title; which heading is title

A doc page takes its title from its first level-1 heading. A page without one fails with the error: Doc <slug> needs a "# Title" line. Later level-1 headings are not used as the title.

## Intro text and section content are split at level-2 headings

Lines 106-120

Asked as: intro text before section; where page intro goes; content split by heading; second title heading shows

Content before the first level-2 heading (##) is the page intro. Content after a level-2 heading belongs to that section. Only the first level-1 heading (#) is the title, and it does not appear in the intro. A later level-1 heading stays in the content.

## What's new heading ends with its last commit

Lines 109-116

Asked as: what's new week last commit; heading shows commit hash; to commit suffix in title; week covers commits

A level-2 heading that ends with "(to <8 to 40 hex characters>)", such as "(to af436ec9)", sets the last commit that the What's new week covers. The section title drops that suffix. Other level-2 headings keep their full text and have no commit.

## Section words and page words are counted differently

Lines 122-134

Asked as: word count wrong; section word count; page word count includes headings; why counts differ

A section's word count (words) comes from the raw text of its own content only, and its heading is not counted. The page's word count comes from the whole source, so it includes the title and all headings.

## Help page order and unlisted pages

Lines 137-178

Asked as: help pages order; page missing from help list; where new page appears; sort order of help pages

Help lists the pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, profile, glossary, whats-new. A page that the list lacks goes just before What's new. Unlisted pages sort by title among themselves.

## Child pages follow their parent in Help

Lines 157-185

Asked as: child page placement; sub page under parent; review time page order; subpages in help list

A child page (review-time, revise-a-purchase-order, find-a-purchase-order, commit-expenses, admin-profiles) shows right after its parent page. Children of one parent show in the order of the childDocPages list. A child page is not sorted as a top-level page.

## Side navigation label and indent for child pages

Lines 187-191

Asked as: side menu shorter name; nav label indented child; sidebar page name indent; child page nav title

In the Help side navigation, a child page shows its shorter label, such as "Review and track", "Revisions", "Find and print", "Review and commit" or "Admin profiles", and is indented. Other pages show their title and have no indent.

## Help popover links open in a new tab

Lines 193-206

Asked as: popover link new tab; help popup link opens tab; link in popup closes form; unsaved work lost popover

Links in a help popover open in a new tab (target="_blank", rel="noopener"), so people do not lose unsaved work on the form under the popover. This covers links in text blocks, in flow step details and in flow notes.

## Help section reference lookup

Lines 208-217

Asked as: help link to section not found; section reference returns nothing; help anchor lookup fails; hash link to help section

A help reference such as "project-requests#approve-a-request" names a page slug and a section id after the #. Turbo finds the section only when both the page and the section exist. Otherwise it finds nothing (null).
