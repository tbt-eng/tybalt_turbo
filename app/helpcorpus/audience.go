package helpcorpus

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

// This file is the audience rules. It decides who may read each chunk.
//
// The repo is public, so its docs, design notes and code are open to every
// signed-in user. Nothing here is a secret. The rules exist so that a folder or
// a word can be restricted later with one line.
//
// A rule names claims. A user passes a rule if they hold any one of them. A
// chunk can match a path rule and several term rules, and the user must pass
// every one of them, so the strictest rule wins.

// everyone is a rule that every signed-in user passes.
var everyone = []string{}

// defaultRequires applies to a file that no path rule matches. A new source is
// then closed until a person adds a rule for it.
var defaultRequires = []string{"admin"}

type pathRule struct {
	// Prefix is a repo path, such as "ui/src/docs/". The longest prefix that
	// matches a file wins.
	Prefix   string
	Requires []string
}

var pathRules = []pathRule{
	{Prefix: "ui/src/docs/", Requires: everyone},
	{Prefix: "descriptions/", Requires: everyone},
	{Prefix: "app/", Requires: everyone},
	{Prefix: "ui/src/", Requires: everyone},
}

type termRule struct {
	// Term is text to find in a chunk or in its path. Upper and lower case are
	// the same.
	Term     string
	Requires []string
}

// No word is restricted. Uncommit is not a secret: the chat's instructions tell
// the model not to raise it, so that people don't come to rely on it. That is
// guidance for the answer, not a limit on what the model may read.
var termRules = []termRule{}

// handbookBans matches text that the user docs must not hold, because the
// handbook goes into every prompt. The user docs treat a commit as final (see
// ui/src/docs/README.md). "Uncommitted", for a record that nobody has committed
// yet, is fine.
var handbookBans = regexp.MustCompile(`(?i)\buncommit(s|ting)?\b`)

// requiresFor returns the groups of claims for a chunk of a file. The result is
// sorted, and is empty when the chunk is for everyone.
func requiresFor(path, text string) [][]string {
	groups := [][]string{pathRequires(path)}
	for _, rule := range matchingTerms(path, text) {
		groups = append(groups, rule.Requires)
	}
	return simplify(groups)
}

func pathRequires(path string) []string {
	best := pathRule{Requires: defaultRequires}
	matched := false
	for _, rule := range pathRules {
		if strings.HasPrefix(path, rule.Prefix) && (!matched || len(rule.Prefix) > len(best.Prefix)) {
			best, matched = rule, true
		}
	}
	return best.Requires
}

func matchingTerms(path, text string) []termRule {
	haystack := strings.ToLower(path + "\n" + text)
	var matches []termRule
	for _, rule := range termRules {
		if strings.Contains(haystack, strings.ToLower(rule.Term)) {
			matches = append(matches, rule)
		}
	}
	return matches
}

// simplify sorts the groups and drops the ones that add nothing: a group for
// everyone, a copy of another group, and a group that a smaller group implies.
// A user who holds "admin" also passes "admin or commit".
func simplify(groups [][]string) [][]string {
	var clean [][]string
	for _, group := range groups {
		if len(group) == 0 {
			continue
		}
		sorted := slices.Clone(group)
		sort.Strings(sorted)
		clean = append(clean, slices.Compact(sorted))
	}
	sort.Slice(clean, func(a, b int) bool {
		if len(clean[a]) != len(clean[b]) {
			return len(clean[a]) < len(clean[b])
		}
		return strings.Join(clean[a], ",") < strings.Join(clean[b], ",")
	})
	kept := [][]string{}
	for _, group := range clean {
		implied := false
		for _, smaller := range kept {
			if containsAll(group, smaller) {
				implied = true
				break
			}
		}
		if !implied {
			kept = append(kept, group)
		}
	}
	return kept
}

func containsAll(group, subset []string) bool {
	for _, claim := range subset {
		if !slices.Contains(group, claim) {
			return false
		}
	}
	return true
}
