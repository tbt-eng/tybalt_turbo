<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"6d5a3cdb4a0cee1a03be4b674951cb0adb397a8c06bfb82be18824d597fd76ea","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Doc links render with the link style

Lines 33-39

Asked as: doc link style; links in user docs look wrong; help link not styled; anchor class link

Turbo renders every Markdown link in the user docs as an anchor with the class "link". The link keeps its href and its text.

## Commit links point to the GitHub repository

Lines 45-49

Asked as: commit link goes to wrong place; link to commit on github; commit hash link; where commit links point

Turbo links a commit to https://github.com/tbt-eng/tybalt_turbo/commit/ followed by the commit hash.

## Section anchor ids come from the section title

Lines 51-57

Asked as: section anchor id; heading link slug; why section link broken; how section ids are made

Turbo builds a section id (slugify) from the title. It lowercases the title and removes apostrophes. It then turns each run of characters other than a-z and 0-9 into one hyphen and trims hyphens at the start and end.

## Word count counts only words that start with a letter

Lines 59-61

Asked as: word count wrong; numbers not counted as words; how words are counted; count of words in section

Turbo counts a word (countWords) when it starts with a letter and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone are not counted. Text with no such word counts as 0.

## Flow lines need exactly three non-empty parts

Lines 63-80

Asked as: flow line error; flow line needs three parts; flow step missing detail; empty part in flow

In a flow block, each line is "Actor | Step | Detail" or "Note: text". A line with other than 3 parts, or with an empty part, fails with: A flow line needs "Actor | Step | Detail": followed by the line. Empty lines are skipped. The Detail and Note text accept Markdown inline formatting.

## Only code blocks marked flow become flow blocks

Lines 82-100

Asked as: flow code block not working; code block shows as text; which code blocks become flows; blank tokens dropped

A code block with the language "flow" becomes a flow block of steps and notes. All other content becomes HTML blocks. Blank-space tokens are dropped.

## A doc page needs one level-1 title

Lines 102-121

Asked as: doc page error no title; missing page title; page has no heading; second level one heading

The first level-1 heading (# Title) is the page title. A page with no such heading fails with: Doc <slug> needs a "# Title" line. A later level-1 heading is treated as ordinary content.

## Level-2 headings split a page into sections

Lines 106-135

Asked as: section split on h2; page intro before section; why headings split page; section word count

Each level-2 heading starts a new section. Content before the first one is the page intro. Each section counts its own words.

## What's new heading carries its last commit

Lines 109-116

Asked as: whats new commit suffix; commit shown in heading; release heading shows commit; commit at end of heading

A level-2 heading that ends with "(to <commit hash>)", where the hash has 8 to 40 hexadecimal characters, sets the section's last commit. The heading text shows without that suffix. A heading without that suffix has no commit.

## Commit hash must be 8 to 40 hex characters

Lines 110-115

Asked as: commit hash too short; heading commit not recognised; invalid commit in whats new; hex hash length rule

A What's new heading takes a commit only when it ends with "(to " plus 8 to 40 lowercase hex characters plus ")". Any other heading keeps its full text as the title and has no commit.

## Order of the pages in Help

Lines 145-156

Asked as: help page order; order of help pages; why pages listed differently; sequence of help topics

Help lists the pages in a fixed order (docPageOrder). The order is: Start here, Time, Purchase Orders, Expenses, Proposals, Project Requests, Rate Sheets, Profile, Glossary, What's new.

## Child pages appear under their parent with a shorter label

Lines 160-198

Asked as: subpage indented in help; child page under parent; shorter label help menu; nested help topic order

Help lists some pages under a parent page, indented, with a shorter label. Examples are "Review and track" under Time, "Revisions", "Find and print" and "Child POs" under Purchase Orders, and "Periods", "Duplicates", "Codes and names" and "Delete" under Rate Sheets. A child page follows its parent. The children of one page show in the order of the childDocPages list.

## Pages missing from the order list go before What's new

Lines 176-185

Asked as: missing page in help list; new doc page not showing; where new page appears; page not in order list

A page that is not in docPageOrder and is not a child page is placed just before What's new. These pages are sorted by title.

## Links in help popovers open in a new tab

Lines 200-213

Asked as: help link opens new tab; popover link new window; link closes my form; help link keeps form safe

Links inside a help popover open in a new browser tab (target="_blank"), so a person does not lose unsaved work on the form. This applies to links in text blocks and in the steps and notes of a flow.

## A help reference that matches no page or section gives nothing

Lines 215-224

Asked as: help link returns nothing; broken help reference; help anchor not working; empty help for link

A help reference such as "project-requests#approve-a-request" needs both a page slug and a section id that exist. If either one is missing, Turbo returns no help content for it.
