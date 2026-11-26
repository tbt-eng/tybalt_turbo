<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"8ccc04b2fcf4f9284f5ee533f70f0e461129a6917cecb74196c3d10523cf165e","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages carry the link style

Lines 33-39

Asked as: links look like plain text; help link not styled; link colour wrong in help; help links look odd

Every Markdown link in a help page is drawn as an anchor with the CSS class "link". It keeps the address (href) that the author wrote.

## Commit links point to GitHub

Lines 45-49

Asked as: commit link goes wrong place; link to a commit; github commit link broken; where commit link points

A commit link goes to https://github.com/tbt-eng/tybalt_turbo/commit/<commit>.

## Section ids come from titles

Lines 51-57

Asked as: section anchor name wrong; heading id has apostrophe; section url slug; link to a heading

A section id is its title in lower case, with apostrophes removed and other runs of non-letters and non-digits turned into one hyphen. Hyphens at the start and end are removed. Section word counts use the section's raw text.

## Word count rule

Lines 59-61

Asked as: word count is wrong; why numbers not counted; how docs count words; word count skips numbers

Docs count a word as a run that starts with a letter and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone do not count.

## Flow lines need three parts

Lines 63-80

Asked as: flow line rejected; help line error missing part; actor step detail error; note line in flow

A flow line in a doc must be "Actor | Step | Detail" with all three parts filled in, or it is rejected with: A flow line needs "Actor | Step | Detail": <line>. A line that starts with "Note:" is a note instead. Blank lines are skipped.

## A flow code block becomes a step list

Lines 82-100

Asked as: flow block shows as code; steps show as code; code block not turning into steps; flow code not rendering

A fenced code block whose language is "flow" becomes a flow block with steps and notes. It does not show as code. Other content around it shows as normal HTML. Blank spacing between blocks is dropped.

## Content before the first section heading is the page intro

Lines 102-120

Asked as: text above first heading; intro text in help page; content before first section; help page intro missing

A help page takes its title from the first level-1 heading. Each level-2 heading starts a new section. Content that comes before the first level-2 heading goes into the page intro. Content after it goes into the latest section.

## A doc page needs a title

Lines 107-121

Asked as: doc missing title error; no title line error; help page won't build; doc page error on save

A doc source with no "# Title" line (first level 1 heading) is rejected with the error: Doc <slug> needs a "# Title" line. Only the first level 1 heading is the title.

## What's new heading carries its last commit

Lines 109-116

Asked as: commit in what's new heading; version hash in heading; what's new commit hash; heading shows commit text

A level 2 heading that ends with "(to <8 to 40 hex characters>)" sets the section's last commit and drops that text from the title. Other level 2 headings have no commit.

## Section and page word counts

Lines 126-133

Asked as: section word total; page word total; word count ignores heading; why word count differs

A section's word count covers the text of its own content, and does not cover its heading. The page's word count covers the whole source, including the title and headings.

## Principle for the order of help pages

Lines 137-144

Asked as: why help pages in this order; reason for page ordering; order of help menu logic; how help pages are sorted

Help lists pages in this order. First comes Start here. Then the task pages, in the order of Turbo's main menu: Time, Purchase Orders, Expenses, Business. In Business, a proposal comes before the project request that it can lead to. Then the page about your profile. Last come the reference pages: the glossary, then What's new.

## Order of pages in Help

Lines 145-159

Asked as: help page sequence; pages listed in wrong order; help menu order list; proposal page comes before project

Help lists pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, clients, absorb-proposals, computers, profile, glossary, whats-new. A proposal page comes before the project request page.

## Child pages show under their parent in Help

Lines 163-179

Asked as: subpage under time; review and track missing; child PO indent; nested help page under parent

Some pages (childDocPages) appear indented under a parent page with a shorter label. Examples: review-time appears under time as "Review and track", and child-purchase-orders appears under purchase-orders as "Child POs". Children show in the order of the childDocPages list, right after their parent.

## Unlisted help pages go before What's new

Lines 183-199

Asked as: help page missing from list; new page not showing in help; page appears before what's new; unlisted help page placement

A page that is not in the page order list and is not a child page goes just before What's new. Unlisted pages are sorted by title.

## Side navigation label and indent for child pages

Lines 202-205

Asked as: side menu label wrong; sidebar shows full title; child page indent sidebar; help sidebar short name

The Help side navigation shows a child page with its short label and an indent. Any other page shows its full title with no indent.

## Links in help popovers open in a new tab

Lines 207-220

Asked as: help link opens same tab; popover link loses form; help popover link new window; link in help popup opens

Links in a help popover over a form open in a new tab (target=_blank, rel=noopener). This keeps unsaved work on the form safe. It applies to links in text, flow step details and flow notes.

## Help reference needs both page and section

Lines 223-231

Asked as: help link does nothing; section link not opening; help reference not found; help hash link empty

A help reference such as "project-requests#approve-a-request" opens only when both the page and the section exist. If either is missing, Turbo finds nothing.
