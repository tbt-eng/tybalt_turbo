<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"3563310351f2291ba2a7a0e44eb9f6ef2fdce212c81ab2634bde8195e8e05294","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Doc links use the link style

Lines 33-39

Asked as: link looks like plain text; doc link not blue; docs links not clickable look; anchor link style

Markdown links in the docs show as an anchor with the "link" class.

## Commit links point to the tybalt_turbo repository

Lines 45-49

Asked as: commit link wrong repo; link to commit on github; commit link broken; where commit link goes

A commit link goes to https://github.com/tbt-eng/tybalt_turbo/commit/ followed by the commit id.

## Section ids come from heading text

Lines 51-57

Asked as: section anchor name; heading link id; why section id changed; anchor from heading title

A section id is its title in lower case. Apostrophes are removed, and each run of other characters becomes one hyphen, with no hyphen at the start or end.

## Word counts for help pages and sections

Lines 59-61

Asked as: word count wrong; how words are counted; numbers not counted as words; help page word total

Turbo counts a word as a run that starts with a letter and may contain letters, digits, underscores, apostrophes and hyphens. Numbers alone do not count. A section's count uses only its own content, and the page's count uses the whole source.

## Flow lines need three parts

Lines 63-80

Asked as: flow line error; flow step missing part; actor step detail error; flow block won't save

A flow block line must have the form "Actor | Step | Detail", with all three parts filled in. Otherwise Turbo fails with the error 'A flow line needs "Actor | Step | Detail": <line>'. A line that starts with "Note:" is a note, and blank lines are skipped.

## Only flow code blocks become flow diagrams

Lines 82-100

Asked as: code block not showing as diagram; flow diagram not rendering; code block shows as html; flow code not working

A code block with the language "flow" becomes a flow block of steps and notes. Other content becomes HTML. Blank space between tokens is dropped.

## A help page needs a title

Lines 103-121

Asked as: help page won't load; missing title error doc; page fails to load heading; doc needs title line

A help page needs a level-1 heading. Only the first one is the title. Without it, Turbo fails to load the page with the error "Doc <slug> needs a "# Title" line".

## What's new heading carries its last commit

Lines 109-116

Asked as: what's new commit in heading; to commit suffix heading; what's new week last commit; commit id at end of title

A level-2 heading that ends with "(to <8 to 40 hex characters>)" sets the last commit that the What's new week covers. Turbo removes that suffix from the section title. A heading without the suffix keeps its full title and has no commit.

## Text before the first section heading is the page intro

Lines 125-132

Asked as: page intro text; text before sections; where intro comes from; text above first heading

A help page's intro (intro) is the content that comes after the title and before the first level-2 heading. Everything after a level-2 heading belongs to a section. Turbo shows the sections in the order of their headings.

## Section word count leaves out the heading

Lines 131-131

Asked as: section word count low; heading not counted in words; section counts don't add up; words under heading only

A help section's word count (words) counts only the body text under its heading. The heading is not counted, so the section counts add up to less than the page count.

## Page word count covers the whole document

Lines 133-133

Asked as: page word total; total words help page; word count includes title; whole page word count

A help page's word count (words) counts every word in its source text. That includes the title, the section headings and the body of every section. Turbo counts words with countWords.

## Help page order in the side navigation

Lines 137-158

Asked as: side navigation order; help menu order wrong; page order in help; help pages sequence

Help lists pages in a fixed order (docPageOrder): Start here, Time, Purchase Orders, Expenses, Proposals, Project Requests, Rate Sheets, Absorb Proposals, Computers, Profile, Glossary, What's new. A proposal comes before the project request that it can lead to.

## Child pages appear under their parent with short labels

Lines 160-175

Asked as: child page label; indented help page; child POs label; review and track name

Some Help pages show under a parent page, indented, with a shorter label. Examples: "review-time" shows as "Review and track" under Time, and "child-purchase-orders" shows as "Child POs" under Purchase Orders. A child page that has no entry here shows with its own title and no indent.

## Help pages not in the order list go before What's new

Lines 177-195

Asked as: help page missing from menu; new page where in nav; unlisted help page position; page sits before what's new

A page that is not in the order list goes just before What's new, sorted by title. Pages in the list keep the list order. A child page follows its parent, in the order of the child list.

## Help popover links open in a new tab

Lines 203-216

Asked as: popover link opens new tab; help link closes my form; link in help popup; help link loses unsaved work

Links inside a help popover open in a new tab (target="_blank", rel="noopener"). This keeps unsaved work on the form that the popover sits over. It applies to text blocks, flow step details and flow notes.

## Help section reference lookup

Lines 218-227

Asked as: help section link broken; section anchor not found; hash link to help section; help reference shows nothing

A help reference such as "project-requests#approve-a-request" names a page and a section id. If the page or the section does not exist, Turbo finds nothing and shows no section.
