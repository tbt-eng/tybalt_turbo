<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"758372f31a1dbe11d0b22f6f8cd7a051801a41f34f918d4d96d228467a7408b3","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Doc links open the page's link style and render as links

Lines 33-39

Asked as: link looks wrong in help; help links not blue; link text formatting in docs; clickable text not styled

Turbo renders each Markdown link in a help page as an anchor with the class "link" and the given address. The link text is parsed as inline Markdown.

## Commit links point to the GitHub commit page

Lines 45-49

Asked as: commit link goes nowhere; link to github commit; what's new commit url; github commit button broken

A What's new commit links to https://github.com/tbt-eng/tybalt_turbo/commit followed by the commit hash (commitUrl).

## Section links use a slug of the title

Lines 51-57

Asked as: section anchor link broken; jump link to section; heading link id; section link not working

The id of a section is its title in lower case, with apostrophes removed and other runs of non-letters and non-digits replaced by one hyphen. Hyphens at the start and the end are trimmed.

## Word count of a doc section

Lines 59-61

Asked as: how many words counted; word count looks wrong; word count skips numbers; why word count is low

Turbo counts words as runs that start with a letter and may contain letters, digits, underscores, apostrophes or hyphens. Numbers alone do not count as words.

## Flow lines need three parts

Lines 63-80

Asked as: flow line error message; actor step detail error; flow line invalid format; flow step error text

Each line of a flow block must be 'Actor | Step | Detail' with all three parts filled, or a line that starts with 'Note:'. A line with another number of parts or an empty part fails with the error 'A flow line needs "Actor | Step | Detail"'.

## Flow note lines become notes, not steps

Lines 67-72

Asked as: note in flow diagram; note: line not showing; flow note formatting; flow note shows as step

In a flow block, a line that starts with "Note:" becomes a note, and Turbo shows the text after "Note:" as inline Markdown. Blank lines are skipped.

## Flow step line with a missing part is rejected

Lines 73-77

Asked as: docs won't load; help page fails to load; missing flow field error; flow step missing part

Every flow step line must have exactly three non-empty parts separated by "|": Actor, Step and Detail. Otherwise the docs fail to load with the error 'A flow line needs "Actor | Step | Detail": ' followed by the line. Only the Detail part is parsed as inline Markdown.

## Only code blocks marked flow become flow diagrams

Lines 82-100

Asked as: code block not a diagram; flow block not rendering; fenced code flow language; code block shows as text

A fenced code block whose language is "flow" becomes a flow block with steps and notes. All other content is rendered as HTML. Blank space tokens are dropped.

## A doc page needs a title

Lines 106-121

Asked as: help page missing title; doc failed to load title; no title error in doc; page won't load no heading

A doc page must have a first-level '# Title' line. A doc without one fails to load with the error 'Doc <slug> needs a "# Title" line'. Only the first level-1 heading is the title.

## Only the first level-1 heading is the page title

Lines 106-108

Asked as: second h1 heading; multiple titles in doc; extra heading shows as text; wrong page title shown

The first "#" heading is the page title. Later level-1 headings stay in the content as normal text.

## What's new heading carries its last commit

Lines 109-116

Asked as: what's new commit removed from title; to hash in heading; week commit in heading; heading shows commit id

A What's new heading that ends with '(to' and a hex commit of 8 to 40 characters sets the last commit that the week covers. The commit is removed from the shown title. A heading without it has no commit.

## Level-2 headings start sections; earlier content is the intro

Lines 109-119

Asked as: intro text before heading; where does a section start; content above first heading; text under which section

Each "##" heading starts a new section. Content before the first "##" heading is the page intro, and content after a heading belongs to that section.

## What's new commit needs 8 to 40 hex characters

Lines 110-116

Asked as: invalid commit hash heading; hash not recognised in what's new; commit missing from section; short hash in heading

A section heading ending in "(to <hash>)" sets the section's commit only when the hash is 8 to 40 lowercase hex characters. With any other ending, the heading keeps its full text and has no commit.

## Section IDs come from the heading without the commit

Lines 112-127

Asked as: section id has hash; anchor name in what's new; slug includes commit; section link has commit

A section's ID is the slug of its title. For a What's new heading, the title leaves out the "(to hash)" part.

## Section words and page words are counted differently

Lines 126-133

Asked as: section word count wrong; page word count differs; section count includes heading; word count per section

A section's word count covers only that section's content, not its heading. The page's word count covers the whole source.

## Help page order in the side navigation

Lines 137-200

Asked as: side navigation order; help pages in wrong order; help menu order; why pages are sorted this way

Help lists pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, clients, absorb-proposals, computers, profile, glossary, whats-new (docPageOrder). A page that is not in the list goes just before What's new. Unlisted pages sort by title.

## Child pages show under their parent page

Lines 161-206

Asked as: sub page under time; child page indent; review and track page location; child page not under parent

Help lists some pages under a parent page with a shorter label, for example 'review-time' under Time as 'Review and track' and 'child-purchase-orders' under Purchase Orders as 'Child POs'. A child page follows its parent and is indented in the side navigation. Children of one parent show in the order of the childDocPages list.

## Pages missing from the order list go before What's new

Lines 184-193

Asked as: new help page missing from menu; page appears before what's new; page not in navigation; help page placed oddly

A top-level page that is in neither docPageOrder nor the child list sits just before What's new, sorted by title.

## Help popover links open in a new tab

Lines 208-221

Asked as: popover link opens new tab; help link closes my form; help link loses unsaved work; popover link opens in tab

Links in a help popover open in a new tab (target _blank). This keeps unsaved work on the form that the popover sits over. It applies to the text blocks, and to the details and notes of flow steps.

## Help reference needs a page and a section

Lines 223-232

Asked as: help link does nothing; anchor link finds nothing; deep link not found; help reference shows no result

A help reference such as 'project-requests#approve-a-request' gives a result only when both the page and the section exist. If either one is missing, Turbo finds nothing.
