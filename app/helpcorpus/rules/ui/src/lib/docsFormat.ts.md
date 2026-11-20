<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"15af5da77b9d6592a915b5c518e6d93ef69a29d5ecc9d6c6ae68fd49e7320143","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Docs links open as Turbo styled links

Lines 33-39

Asked as: docs link style; links look plain in help; link opens same tab in docs; docs anchor class

The docs Markdown parser renders every link as an anchor with the class "link" and the given href. It adds no target attribute at this step. Only the help popover step adds new-tab behavior.

## What's new heading carries its last commit

Lines 45-49

Asked as: what's new commit link; weekly update commit hash; which commit does week cover; github commit in heading

A level-2 heading that ends with "(to <8 to 40 hex characters>)" gives the section a last commit (the commit that the week covers) and removes the suffix from the title. Other level-2 headings have no commit. A commit links to https://github.com/tbt-eng/tybalt_turbo/commit/<commit>.

## Section ids and word counts

Lines 51-61

Asked as: section anchor id wrong; word count per section; heading to id slug; page word count wrong

A section id comes from its title in lowercase. Apostrophes are removed and other runs of non-alphanumeric characters become a hyphen. A word is a run that starts with a letter. Section word counts use only the section's content, and page word counts use the whole source.

## Flow block line format

Lines 63-80

Asked as: flow line error; flow step format wrong; actor step detail missing; flow line needs three parts

In a flow block each line is "Actor | Step | Detail" with three non-empty parts, or "Note: text". Any other line fails with the error "A flow line needs "Actor | Step | Detail": <line>". Blank lines are ignored.

## Flow blocks in docs become step lists

Lines 82-100

Asked as: flow code block not rendering; steps shown as code; markdown flow fence; flow block gaps between

In a doc page, a fenced code block marked "flow" becomes a flow block of steps and notes (parseFlow). Other Markdown around it becomes HTML blocks, in order. Blank space between blocks is dropped.

## Doc page needs a title line

Lines 102-121

Asked as: doc page won't load; missing title in doc; doc needs heading error; page fails no title

A doc page must have a level-1 heading ("# Title"). The first level-1 heading is the page title. Without one, the page fails to load with the error "Doc <slug> needs a "# Title" line".

## Doc pages split into intro and sections by headings

Lines 102-120

Asked as: doc sections not splitting; second heading new section; intro text before sections; doc page layout headings

The first level-1 heading is the page title. Each level-2 heading starts a new section. Content before the first level-2 heading goes into the page intro.

## Later level-1 headings stay as content

Lines 106-119

Asked as: second h1 heading ignored; extra top heading shows; multiple titles in page; later title treated as text

Only the first level-1 heading becomes the page title. Any later level-1 heading is treated as ordinary content in the intro or the current section.

## Section word count and page word count

Lines 122-135

Asked as: how many words counted; word count off; section length count; words in raw text

Each section carries a word count of its own content. The page word count covers the whole source text. Words are counted from the raw text.

## Help page order and unlisted pages

Lines 137-186

Asked as: help menu order; help page sort order; unlisted help page position; where help page appears

Help lists top-level pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, computers, profile, glossary, whats-new. A page that is not in that list and not a child page goes just before What's new. Unlisted pages are sorted by title.

## Child pages sit under their parent with short labels

Lines 159-199

Asked as: sub page under parent; indented help page label; shorter name in help menu; child page nested help

Some Help pages show under a parent page, in the order of the childDocPages list, with a shorter label and an indent. For example, review-time shows as "Review and track" under time, and rate-sheet-absorb shows as "Duplicates" under rate-sheets. Other pages show their own title with no indent.

## Help popover links open in a new tab

Lines 201-214

Asked as: help link opens new tab; popover link leaves form; help link unsaved work lost; link in help window

Links in a help popover over a form open in a new tab (target _blank, rel noopener). This keeps unsaved work on the form.

## Help section reference lookup

Lines 216-225

Asked as: help link broken section; anchor not found in help; missing help section link; help reference doesn't work

A help reference such as "project-requests#approve-a-request" names a page slug and a section id. Turbo finds nothing if the page or the section does not exist.
