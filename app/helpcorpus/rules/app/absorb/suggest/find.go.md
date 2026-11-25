<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/suggest/find.go","sha256":"d8a19256afbdb7b57c8d5b723af4cf5f97a5c5e9ad35e50217175d93b9514e79","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/suggest/find.go

## Each member's score is measured against the suggested target

Lines 9-23

Asked as: how similar is each duplicate; score compared to main client; why score not 100; match percentage duplicate; how many references move

In an absorb suggestion, each member other than the suggested target carries a score (Score) that compares it with the target. The target itself always shows a score of 1. Each member's Refs is the number of references that an absorb would move.

## Suggested target is the first member of a bucket

Lines 25-50

Asked as: which record is kept; which duplicate is the main one; which client survives merge; who is the target record

In an absorb suggestion bucket, the first member (Members[0]) is the suggested target. The target's score is 1. The target's id (TargetID) is the first id in the list of member ids (IDs).

## How Turbo groups probable duplicates into absorb suggestions

Lines 52-78

Asked as: find duplicate clients; same vendor listed twice; possible duplicate vendors; similar client names; two records same company

Find compares every client or vendor name and alias and returns buckets of probable duplicates. A pair counts as a match when its score is at least the minimum score (MinScore). The result does not depend on the order of the records.

## Buckets are listed by rank, best first

Lines 71-76

Asked as: order of duplicate groups; which group shows first; sort duplicate suggestions; group ranking order

Turbo sorts buckets by rank, highest first. Buckets with equal rank sort by the target's id, lowest first.

## Records are compared in order of id

Lines 94-97

Asked as: does record order matter; why results change; records sorted by id; same duplicates different order

Turbo sorts the records by id before it compares them. The result of the suggestions does not depend on the order in which the records arrive.

## Only the name counts for how common a word is

Lines 108-125

Asked as: common words in names; alias ignored in word count; why words counted once; word frequency names

Turbo counts how many names contain a word (df) to decide how rare it is. It counts each word once for each record. It counts words from the name only, not from the alias.

## Rare name words weigh more than common ones

Lines 130-137

Asked as: rare name words matter more; why Ltd or Inc ignored; common words count less; word weight in matching

Each word of a name weighs ln(N/df) + 0.5. N is the number of records and df is the number of names that contain the word. A word that few names share counts for more when Turbo scores a pair.

## Adjacent words are also joined for comparison

Lines 138-148

Asked as: words split differently; name with missing space; run together name; space missing between words; joined name spelling

For each name or alias, Turbo also builds the whole text joined together and each pair of adjacent words joined together. It uses them to find records whose words are split or run together differently.

## Alias is compared only when it differs from the name

Lines 149-156

Asked as: alias not checked; alias empty so ignored; alias skipped for duplicates; alias not compared

Turbo compares a record's name. It also compares the alias, but only when the alias has words and its words differ from the name's words. A record with no name words and no alias words is not compared.

## Alias key is skipped when it equals the name

Lines 153-155

Asked as: alias same as name; alias repeats the name; duplicate alias entry; alias adds nothing

Turbo adds the alias as a second key only when its words differ from the name's words (same words in the same order). An alias that equals the name adds nothing.

## Alias with the same words in another order counts as different

Lines 161-171

Asked as: alias words reordered; same words different order; alias with extra words; word order differs alias

Turbo skips the alias only when its words are the same as the name's words, in the same order and in the same number (sameTokens). An alias with the same words in a different order, or with extra words, gets its own key and is compared as well.

## Only similar long words without digits can match as spelling variants

Lines 173-218

Asked as: typo in client name; misspelled vendor name; similar long words; digits must match exactly; spelling variant match

Two different words can match only when both have at least minEditLength letters and no digit. The longer word can be at most 1.5 times as long as the shorter. Words with digits match only when equal.

## Spelling tolerance and same-prefix shortcut for similar words

Lines 206-215

Asked as: how many typos allowed; same start letters match; spelling tolerance setting; prefix match words

The edits allowed for a word pair are the token tolerance (TokenTolerance) times the length of the longer word. When the Prefix setting is on and both words start with the same first letters, Turbo compares them even if the edit count is too high. With no allowed edits and no same prefix, Turbo skips the pair.

## Which record pairs Turbo scores

Lines 229-334

Asked as: why these not matched; clients never compared; no match found; names never scored

Turbo scores two records only when they share an equal or similar word, or have the same whole name joined together. It also scores them when a joined pair of adjacent words is within one letter of a word in the other record. Pairs that meet none of these tests are never scored, so they never match.

## Match threshold includes the minimum score

Lines 328-333

Asked as: exactly minimum score match; minimum score cutoff; equal score counts as match; threshold boundary duplicates

A pair is a match when its record score is at least MinScore, with a tiny tolerance. A pair that scores exactly MinScore matches.

## Buckets are stars and do not chain

Lines 341-382

Asked as: why A B C not grouped; duplicates chain together; group size two; matches of matches grouped

Turbo picks as centre the unassigned record with the most unassigned matches. Ties go to the record with more references. Turbo adds the centre's unassigned matches to its bucket. A match of a match does not join, so buckets do not chain. Records with no unassigned matches form no bucket, so every bucket has at least two members.

## Suggested target is the member with most references

Lines 384-415

Asked as: which client to keep; pick main record for merge; most references kept; choose target for merge

The suggested target of a bucket is the member with the most references. On a tie the centre wins, and the next tie goes to the earlier record by id. Turbo scores the other members against the target and lists them with the highest score first. Ties go to more references, then the lower id.

## Confidence, moved references and rank of a bucket

Lines 416-418

Asked as: what confidence means; how group rank computed; bucket rank score; meaning of moved references; average score of duplicates

Confidence is the mean score of the other members against the target. Moved references are the references of every member except the target. Rank is Confidence² × ((members − 1) + log2(1 + moved references)).
