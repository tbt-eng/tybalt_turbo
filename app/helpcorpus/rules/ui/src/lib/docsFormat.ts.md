<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"8faf74621d34a3ac9d036d494796b0c2605b6759586171ac7e0fc8ba19959b8d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages use the link style

Lines 33-43

Asked as: help link looks plain; link not styled in help; help links not blue; anchor looks like plain text

Every Markdown link in a help page renders as an anchor with the CSS class "link" and the given address. This applies to links in Markdown text and to links in flow step details and notes.

## What's new headings carry their last commit

Lines 45-116

Asked as: commit hash in heading; git hash in what's new; version number in title; link to commit on github

A level 2 heading that ends with "(to " followed by 8 to 40 hex characters and ")" sets the section's last commit. The commit text is removed from the section title. A commit links to the tbt-eng/tybalt_turbo repository on GitHub.

## Help references find a page and section by slug and id

Lines 51-57

Asked as: help link goes nowhere; broken help link; help link shows nothing; page link does not open

A help reference such as "project-requests#approve-a-request" names a page slug and a section id. If the page or the section does not exist, Help finds nothing. Section ids come from the section title in lowercase, with apostrophes removed and other characters turned into hyphens.

## Word count counts only words that start with a letter

Lines 59-61

Asked as: word count skips numbers; numbers not counted words; why word count is low; symbols not counted in count

Help counts a word as a run that starts with a letter (A–Z, a–z) and continues with letters, digits, underscores, apostrophes or hyphens. Numbers on their own and symbols are not counted.

## Flow lines need three parts or Help fails to load the page

Lines 63-80

Asked as: flow line error message; help page fails to load; actor step detail error; flow line missing part

A flow line in a doc must read "Actor | Step | Detail" with all three parts filled in. A line that starts with "Note:" is a note instead. Any other line fails with the error: A flow line needs "Actor | Step | Detail": followed by the line.

## A doc needs a level 1 title

Lines 102-121

Asked as: doc missing title error; page won't load no title; missing level one heading; doc needs heading one

A doc page takes its title from its first level 1 heading. A doc without one fails with the error: Doc <slug> needs a "# Title" line. Level 2 headings start the page's sections.

## Text before the first level 2 heading is the page intro

Lines 122-125

Asked as: text above first heading; intro paragraph before sections; content before first section; text before sections shows

A help page's intro (intro) holds the content before its first level 2 heading, apart from the level 1 title. Flow blocks in the intro render the same way as in sections.

## Section anchor ids come from the section title

Lines 126-129

Asked as: section link not jumping; heading link anchor broken; section id from heading; link to heading broken

Each help section gets an id (id) made from its title after the commit suffix is removed. The id is the slugified title. Links to a section use this id.

## Section and page word counts

Lines 131-133

Asked as: word count per section; how many words in section; page word total; section word count includes heading

A section's word count (words) covers only the content under its heading, not the heading itself. A page's word count covers the whole source file.

## A test fails when a page is not in the order lists

Lines 137-144

Asked as: test failing new page; page not in order list test; unit test fails help page; build fails new help page

The page order in Help follows set rules: Start here, then task pages in the order of the main menu (Time, Purchase Orders, Expenses, Business), then the profile page, then the glossary and What's new. A test fails when a page is in neither docPageOrder nor childDocPages.

## Order of pages in Help

Lines 145-155

Asked as: order of help menu; where pages sit in help; help sidebar order; help pages sequence

Help lists pages in this order: Start here, Time, Purchase Orders, Expenses, Proposals, Project Requests, Profile, Glossary, What's new. The list is docPageOrder.

## Child pages show under their parent page in Help

Lines 159-192

Asked as: sub pages under time; nested pages in help; indented menu item help; child page under parent

Six pages show indented under a parent page with a shorter label. Time has Review and track. Purchase Orders has Revisions and Find and print. Expenses has Vendor invoices and Review and commit. Profile has Admin profiles. Children appear in the order of the childDocPages list.

## Pages missing from the order list go before What's new

Lines 170-186

Asked as: page at bottom of help; new page missing from help menu; unlisted help page position; page sits above what's new

A page that is not in docPageOrder and is not a child page goes just before What's new. Several such pages are sorted by title.

## Links in help popovers open in a new tab

Lines 194-207

Asked as: popover link opens new tab; help link closes my form; lost form data clicking link; link in form help opens tab

In a help popover over a form, every link opens in a new tab (target _blank, rel noopener). This keeps unsaved work on the form. It applies to text blocks, flow step details and flow notes.
