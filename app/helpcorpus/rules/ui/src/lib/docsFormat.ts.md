<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"4f615762b82e5e9eec0a40df63264cbedc5cb964c69c73d899c7097ca117de03","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages use the link style

Lines 33-39

Asked as: link style in help; help page links look plain; anchor class link; link looks wrong in help

Turbo renders each Markdown link in a help page as an anchor with the class "link". The anchor keeps the link's address (href) and its text.

## Commit link goes to GitHub

Lines 45-49

Asked as: commit link goes where; link to commit on github; commit link wrong repo; open commit in github

A commit link points to https://github.com/tbt-eng/tybalt_turbo/commit/<commit>.

## Section ids come from the title

Lines 51-57

Asked as: section anchor id; heading id from title; why section link broken; section id hyphens

A section's id is its title in lowercase, with apostrophes removed and each run of other characters turned into a hyphen. Hyphens at the start and end are trimmed.

## Word counts ignore numbers and symbols

Lines 59-61

Asked as: word count wrong; numbers not counted in words; word count skips digits; symbols counted as words

Turbo counts a word only if it starts with a letter. Text that starts with a digit does not count.

## Flow block line format

Lines 63-80

Asked as: flow line error; flow step format wrong; actor step detail missing part; flow block line fails

Each line of a flow block must be "Actor | Step | Detail" with exactly three non-empty parts, or "Note: text". Any other line fails with: A flow line needs "Actor | Step | Detail": <line>. Blank lines are skipped.

## Flow blocks in help pages become step lists

Lines 82-100

Asked as: flow code block shows as code; steps shown as list in help; flow block not rendering; help steps look like code

In a help page, a code block marked as "flow" is shown as a list of steps and notes, not as plain code. The text before and after the block stays as ordinary page content, in the same order.

## Blank space between help content is dropped

Lines 90-97

Asked as: extra space between help blocks; blank lines in help page; empty gaps in help content; whitespace tokens in page

When Turbo builds a help page, it ignores blank-space tokens between blocks. Only content tokens go into the page.

## A doc page needs a level-1 title

Lines 102-121

Asked as: doc page won't load; missing title in help page; doc page error no title; help page fails to load

Turbo refuses to load a doc page that has no "# Title" line. The error says: Doc <slug> needs a "# Title" line. Only the first level-1 heading is the title.

## Help page text before the first section is the intro

Lines 102-120

Asked as: text above first heading; help intro before section; content before first section; where intro text goes

Content that comes before the first level-2 heading goes into the page intro. Content after a level-2 heading belongs to that section until the next level-2 heading.

## Only the first level-1 heading is the page title

Lines 107-119

Asked as: second h1 heading; two page titles in help; extra level one heading; wrong page title in help

The first level-1 heading becomes the page title. A later level-1 heading is not a title, so it goes into the intro or the current section as ordinary content.

## What's new heading carries its last commit

Lines 109-116

Asked as: whats new commit hash; heading shows commit id; to commit at end of heading; commit shown in whats new

A level-2 heading that ends with "(to <8 to 40 hex characters>)" sets the section's last commit and is removed from the section title. A heading without that ending has no commit and keeps its full title.

## Section word count uses the section's own content

Lines 126-133

Asked as: section word count; word count per section; page word count total; section count doesn't match

Each section's word count comes from the raw text of the content under its heading. The page word count comes from the whole source, headings included.

## Help page order in the side navigation

Lines 137-187

Asked as: help menu order; side nav page order; why help pages out of order; where new help page appears

Help lists pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new (docPageOrder). A page that this list lacks goes before What's new. Unlisted pages sort by title.

## Child pages show under their parent with a short label

Lines 158-193

Asked as: subpage under help page; indented help link; child page label; review and track label

Some Help pages show under a parent page, indented, with a shorter label (childDocPages). Examples: review-time under time is "Review and track", child-purchase-orders under purchase-orders is "Child POs", and admin-profiles under profile is "Admin profiles". The children of one page show in the order of that list.

## Help links in popovers open in a new tab

Lines 195-208

Asked as: help popover link opens; link replaces form in popover; help link new tab; popover link closes form

In a help popover over a form, every link opens in a new tab (target=_blank, rel=noopener). This keeps unsaved work. It applies to links in text blocks, flow step details and flow notes.

## Help references need a page and a section

Lines 210-219

Asked as: help link shows nothing; broken help reference; help anchor not found; link to help section empty

A help reference such as "project-requests#approve-a-request" finds a page by slug and a section by id. If the page or the section is missing, Turbo finds nothing.
