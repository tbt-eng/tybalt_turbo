<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsSearch.ts","sha256":"610e0e9f21c38747611dab17452e403082a25575e901a731e1d65a758e465efa","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsSearch.ts

## Help search shows at most 8 results

Lines 27-27

Asked as: too many search results; only eight results; results cut off; why fewer hits; search list limit

A docs search returns at most 8 hits (maxHits). Hits whose title holds the whole query come first, and the rest follow in relevance order.

## Search results show a 160-character excerpt around the first match

Lines 27-29

Asked as: long preview text; search snippet too long; preview cut short; how much text shown; help result preview length

Each hit shows an excerpt of up to 160 characters (excerptLength). It starts about 50 characters before the first matched word (excerptLead) and begins at a whole word. Turbo adds "…" at the start or end when text is cut. Matched words, including their endings, are marked. When nothing matches, the excerpt starts at the beginning of the text.

## Search ignores filler words in a question

Lines 31-36

Asked as: question words ignored; search skips how and what; why do my question words; stop words in search; search fails with a question

Docs search drops filler words such as "how", "do", "i", "the", "what" and "who" from the query (fillerWords). If every word of the query is a filler word, Turbo keeps them all and searches for them.

## Search indexes each page intro and each section separately

Lines 38-83

Asked as: search doesn't find my step; page missing from help search; steps not searchable; help search misses section; intro not in search

Each docs page gives one searchable entry for its intro and one for each section. An intro entry has no section name. The text of a step list includes each step's actor, step and detail, plus its notes. HTML tags are removed and common entities such as &amp; are turned back into characters.

## Search strips formatting from doc text

Lines 38-62

Asked as: weird symbols in results; html showing in search; &amp, in search results; stray code in help text; search shows tags

Search reads doc text without HTML tags. It decodes &lt; &gt; &quot; &#39; and &amp;, and collapses extra spaces. A step in a flow is searched as "actor: step. detail" followed by its notes.

## Search turns discouraged words into the glossary word

Lines 85-109

Asked as: timesheet search finds nothing; old term in search; employee search not working; search uses different word; PA search gives no results

Docs search reads the glossary table's "Don't use" column. It replaces each listed word, such as "timesheet", "employee", "unbundle", "PA" or "billing profile", with the word the docs use ("time sheet", "staff", "recall", "project authorization", "invoicing profile"). The match ignores case and allows a plural "s". Replacement works on whole words only.

## Search changes discouraged words only as whole words, with optional s

Lines 93-108

Asked as: timesheets search swapped word; search changes my word; old words in a list; plural old term search; search rewrites my query

A glossary "Don't use" word in the query is replaced only as a whole word, with or without a final "s" (such as "timesheets"). Several discouraged words in one row, separated by commas, each map to the same glossary word. Rows with an empty "Don't use" cell give no replacement.

## Search treats a glossary word in parentheses as the plain word

Lines 96-99

Asked as: PO search finds purchase order; searching abbreviation in brackets; search PO gives purchase order; brackets in glossary term; acronym search results

When Turbo builds the glossary synonyms, it removes a trailing parenthesised part from the glossary word. For example, "Purchase order (PO)" becomes "purchase order". A search for a discouraged word in that row finds "purchase order".

## Search treats plural and singular as the same word

Lines 111-114

Asked as: order and orders same; plural search no match; expense versus expenses; singular plural search; search ends with s

Docs search cuts a trailing "s" from a word when the word has more than 3 letters and does not end in "ss". "Orders" and "order" match each other, and "expenses" matches "expense". Words of 3 letters or fewer, such as "pos", keep their "s".

## Search returns nothing when the query has no words

Lines 125-129

Asked as: search with only symbols; blank search gives nothing; empty query no results; just punctuation search; no results for dashes

If the query has no letters or digits after cleaning, Turbo returns no hits. Only a-z and 0-9 count as word characters.

## Search words do not include punctuation

Lines 125-128

Asked as: hyphen breaks search; dash in search words; special characters in query; punctuation in search terms; apostrophe search problem

Search splits the query into words made of the letters a-z and the digits 0-9 only, in lower case. Other characters such as hyphens and punctuation separate words and are dropped.

## Excerpt starts about 50 characters before the first match

Lines 132-136

Asked as: preview starts too early; excerpt begins in wrong place; match not at start of preview; preview starts before word; snippet offset wrong

Each search result excerpt begins about 50 characters (excerptLead) before the first matching word. The start moves forward to the next space so no word is cut. When the match is near the start of the text, the excerpt starts at the beginning.

## Excerpt starts at the beginning when nothing matches

Lines 133-135

Asked as: result has no highlight; no bold word in preview; preview not highlighted; search term missing from snippet; result shows start of page

When no query word appears in the text of a result, the excerpt starts from the beginning of the text and nothing is highlighted. An empty word list also highlights nothing.

## Excerpt highlights every word that starts with a query word

Lines 133-150

Asked as: approves not highlighted; partial word highlighted; highlight wrong words; why is that word bold; words with same start marked

The excerpt marks each word that begins with a query word (match is true), whatever the case. A query word "approve" also marks "approves" and "approver". Only whole-word starts count, so a query word in the middle of a word is not marked.

## Excerpt ends at a word boundary and shows ellipses

Lines 137-139

Asked as: three dots in result; ellipsis at end of preview; dots around snippet; preview cut mid sentence; trailing dots in search

A search result excerpt ends at the last space before the 160-character limit, so it does not cut a word. Turbo adds "… " before the excerpt when it does not start at the beginning of the text. Turbo adds " …" after it when text remains.

## Search needs every word, then falls back to any word

Lines 153-178

Asked as: results only match one word; search returns partial matches; why not every word; results too loose; search finds loose matches

Docs search first looks for entries that hold every word of the query (AND). If no entry does, it returns entries that hold any of the words. Words match as prefixes, so "approv" finds "approve". Words longer than 4 letters also allow small typos (fuzzy 0.2).

## Search ranks titles above body text

Lines 157-173

Asked as: right page not on top; best match ranked low; title result lower than body; page name not first; why is my result buried

A match in a section title counts 3 times and a match in the page title 1.5 times, compared with a match in the body text. A page intro has no section title, so its page title is used as its title. A search for a page's name puts that page first.

## Search ignores letter case and plural when indexing

Lines 159-163

Asked as: capital letters in search; case sensitive search; Orders finds order; uppercase query no results; search ignores capitals

Search lowercases every indexed word and turns it to the singular, in the same way as the query. A search for "Orders" finds "order".

## Section hits use the section title as the result title

Lines 159-166

Asked as: result title is heading; which name shows on result; section name as title; result label wrong; page title in result

For the title ranking, a section's title counts as its title. A page intro uses the page title instead. Page title matches also get a smaller boost (1.5) in all of the page's entries.

## Search matches word beginnings and allows typos in long words

Lines 171-175

Asked as: misspelled word search; spelling mistake still finds; typo in search; fuzzy matching search; near match words

Search matches words that begin with a query word (prefix search). For a query word longer than 4 characters, it also accepts a small spelling difference (fuzzy 0.2). Shorter words get no fuzzy match.

## Search puts a title that contains the whole query first

Lines 179-185

Asked as: exact phrase at top; question phrase ranked first; whole question in heading; who approves first result; best title match first

Turbo normalizes the query (lower case, glossary replacements) and moves the hits whose section title, or page title for an intro, contains that whole phrase to the top. "Who approves" finds "Who approves what" before sections that only say "approve". The other hits keep their order after those.

## Title-first ordering compares the query with the section or page title

Lines 181-184

Asked as: title match goes first; heading matches order; why heading results first; title order in search; sort by title match

Turbo checks whether the title contains the whole normalized query, ignoring case. It uses the section title, or the page title for an intro. Matching results come first, in their score order, and the others follow in their score order.
