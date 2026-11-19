<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"0376d25557d1d735c4a82688f71eb42773e789fa263c61b47ff7824d95ec01fb","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links open as styled links

Lines 33-39

Asked as: link looks like plain text; help link not underlined; hyperlink not styled; link formatting in help

Every link in a help page renders as an anchor with the class "link" and the href that the Markdown gives. The link text keeps its inline formatting.

## Commit links point to the tbt-eng GitHub repository

Lines 45-49

Asked as: what's new commit link; commit link goes to wrong place; link to commit on github; commit link broken

A What's new commit links to https://github.com/tbt-eng/tybalt_turbo/commit followed by the commit hash (commitUrl).

## Section ids drop apostrophes and punctuation

Lines 51-57

Asked as: section link url looks odd; why dashes in anchor; apostrophe in heading link; section name in url

A section id is the title in lower case. Apostrophes (' and ’) are removed. Each run of other characters that are not a-z or 0-9 becomes one hyphen, and hyphens at the start or end are removed.

## Word count counts only words that start with a letter

Lines 59-61

Asked as: word count is wrong; numbers not counted as words; why is word count low; help page word count

Help counts a word only when it starts with a letter (A-Z). Numbers alone are not counted. Words can contain digits, underscores, apostrophes and hyphens.

## Flow block lines need three parts or a Note

Lines 63-80

Asked as: flow line refused error; help step line error; actor step detail error; flow line not accepted

A flow block in a help page takes one line for each step in the form "Actor | Step | Detail". A line that does not have exactly three non-empty parts is refused with: A flow line needs "Actor | Step | Detail": <line>. A line that starts with "Note:" becomes a note, and blank lines are skipped.

## Blank lines are skipped inside flow blocks

Lines 67-68

Asked as: empty lines in flow; blank line error in steps; extra gaps between steps; blank line in help block

Blank lines inside a flow block are ignored. They create neither a step nor an error.

## Flow step and note text supports inline Markdown

Lines 69-77

Asked as: bold text in flow step; formatting in step details; note prefix showing in help; markdown not working in steps

The detail of a flow step and the text of a Note go through inline Markdown. Note text loses its "Note:" prefix and surrounding spaces.

## Flow blocks render as steps and notes

Lines 82-100

Asked as: flow block shows as code; steps not displaying properly; code block looks wrong in help; help steps show as text

A code block marked "flow" in a help page becomes a flow block of steps and notes, apart from the surrounding text. Other Markdown around it is grouped into separate HTML blocks.

## Help page needs a level-1 title

Lines 102-121

Asked as: help page won't load; doc needs title error; page missing heading error; help page has no title

Turbo refuses to load a help page that has no level-1 heading (a "# Title" line). The error is: Doc <slug> needs a "# Title" line. Only the first level-1 heading becomes the page title.

## Level-2 headings split a help page into sections

Lines 106-135

Asked as: help page split into sections; new section starts here; subheading creates a section; sections in help page

A level-2 heading starts a new section of a help page. Content before the first level-2 heading is the page intro. Each section's id comes from its title (slugify).

## Only the first level-1 heading sets the page title

Lines 107-108

Asked as: two titles on help page; second heading shows as text; extra big heading in help; which heading is the title

The first level-1 heading sets the page title. Any later level-1 heading is not the title and is treated as ordinary content.

## What's new heading carries its last commit

Lines 109-116

Asked as: commit hash in heading; to hash in what's new title; suffix removed from heading; commit shown next to section

A level-2 heading that ends with "(to <8 to 40 hex characters>)" gives the section a commit and removes that suffix from its title. A heading that does not match keeps its full text and has no commit.

## Text before the first level-2 heading is the page intro

Lines 117-119

Asked as: text above first section; intro before the headings; content before first heading; help page opening text

Content that comes before the first level-2 heading goes into the page intro. Content after a level-2 heading belongs to that section.

## Section word counts differ from the page word count

Lines 131-133

Asked as: section word count mismatch; section count differs from page; why totals don't add up; word count per section

A section's word count covers the raw text of its own content, and the heading is not counted. The page word count covers the whole source, headings included.

## Help page order is fixed

Lines 145-192

Asked as: help pages in wrong order; why pages listed this way; move page up in help; help list order

Help lists pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, client-invoices, profile, glossary, whats-new. A page that is not in the list goes just before What's new, sorted by title.

## Child pages follow their parent in Help

Lines 159-198

Asked as: subpage under time; child page label shorter; review and track page; page nested under another

Some pages show under a parent page, indented and with a shorter label. For example, review-time shows under time as "Review and track", and rate-sheet-absorb shows under rate-sheets as "Duplicates". The children of one page show in the order of the child list. Other pages keep their own title as the label.

## Pages missing from the order list sort before What's new

Lines 176-185

Asked as: new help page location; page not in list; where does new page appear; help page missing from menu

A page that is not in docPageOrder and is not a child page goes just before What's new, sorted by title.

## Nav label uses the child label and indents child pages

Lines 195-198

Asked as: sidebar sub item indented; side nav shows short label; menu item indented under parent; side menu child page

In the side navigation a child page shows its shorter label and is indented. Other pages show their title without indent.

## Help popover links open in a new tab

Lines 200-213

Asked as: link in popover opens new tab; help link closes my form; lose form when clicking help; help link opens separate window

Links in a help popover open in a new tab (target _blank, rel noopener), so a person does not lose unsaved work in the form. This covers the text, the step details and the notes of flow blocks.

## Help reference needs a matching page and section

Lines 215-224

Asked as: help link finds nothing; section reference not working; hash link gives no result; help link to section empty

A help reference has the form "page-slug#section-id". Turbo finds the section only when both the page and the section exist. If either is missing, it finds nothing.
