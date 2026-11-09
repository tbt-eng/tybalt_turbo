package helpcorpus

import (
	"regexp"
	"strings"
)

// Chunk sizes, in bytes of text. A chunk ends at a natural break once it is big
// enough, so a function or a docs section usually stays in one piece.
const (
	// chunkMin is the least text before a chunk may end at the start of a
	// top-level declaration.
	chunkMin = 400
	// chunkMinMarkdown is the least text before a chunk may end at a heading. It
	// is small, so each section of the user docs is its own chunk.
	chunkMinMarkdown = 120
	// chunkTarget is the least text before a chunk may end at any blank line.
	chunkTarget = 1500
	// chunkMax ends a chunk at any line. Only one very long line goes over it.
	chunkMax = 3000
)

type span struct {
	start, end int // line numbers from 1; end is the last line
	heading    string
	text       string
}

// splitChunks cuts a file into spans of whole lines. The spans are in order and
// don't overlap. Blank lines between spans belong to no span.
func splitChunks(path, content string) []span {
	markdown := strings.HasSuffix(path, ".md")
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	blank := func(i int) bool { return strings.TrimSpace(lines[i]) == "" }

	minSize := chunkMin
	if markdown {
		minSize = chunkMinMarkdown
	}
	names := declarationNames(path)

	var spans []span
	var trail [6]string
	declaration := ""
	start, size, heading := -1, 0, ""
	flush := func(end int) {
		if start < 0 {
			return
		}
		for end > start && blank(end-1) {
			end--
		}
		if end > start {
			spans = append(spans, span{
				start:   start + 1,
				end:     end,
				heading: heading,
				text:    strings.Join(lines[start:end], "\n"),
			})
		}
		start, size = -1, 0
	}

	inFence := false
	for i, line := range lines {
		level := 0
		if markdown && !inFence {
			level = headingLevel(line)
		}
		afterBlank := i > 0 && blank(i-1) && !blank(i)
		strong := level > 0 || (!markdown && afterBlank && startsTopLevel(line))

		if start >= 0 {
			switch {
			case strong && size >= minSize,
				afterBlank && size >= chunkTarget,
				size+len(line)+1 > chunkMax:
				flush(i)
			}
		}
		if level > 0 {
			trail[level-1] = strings.TrimSpace(line[level:])
			for deeper := level; deeper < len(trail); deeper++ {
				trail[deeper] = ""
			}
		}
		if markdown && strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
		}
		if start < 0 {
			if blank(i) {
				continue
			}
			start = i
			switch {
			case markdown:
				heading = headingTrail(trail)
			case strong || i == 0:
				heading = ""
			default:
				// The chunk starts inside a declaration, so it names it.
				heading = declaration
			}
		}
		if names != nil && names.MatchString(line) {
			declaration = strings.TrimRight(strings.TrimSpace(line), " {(")
		}
		size += len(line) + 1
	}
	flush(len(lines))

	// A short tail joins the span before it, so a file doesn't end in a scrap.
	if n := len(spans); n >= 2 && len(spans[n-1].text) < minSize/2 {
		last, before := spans[n-1], &spans[n-2]
		before.end = last.end
		before.text = strings.Join(lines[before.start-1:last.end], "\n")
		spans = spans[:n-1]
	}
	return spans
}

// headingLevel returns 1 to 6 for a Markdown heading line, and 0 for other lines.
func headingLevel(line string) int {
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level == 0 || level > 6 || level >= len(line) || line[level] != ' ' {
		return 0
	}
	return level
}

func headingTrail(trail [6]string) string {
	var parts []string
	for _, title := range trail {
		if title != "" {
			parts = append(parts, title)
		}
	}
	return strings.Join(parts, " > ")
}

// startsTopLevel reports if a line of code starts at the left edge with
// something other than a closing bracket. After a blank line, that is the start
// of a declaration or of its comment.
func startsTopLevel(line string) bool {
	switch line[0] {
	case ' ', '\t', '}', ')', ']':
		return false
	}
	return true
}

var (
	goDeclaration = regexp.MustCompile(`^(func|type|var|const) `)
	tsDeclaration = regexp.MustCompile(`^(export |async |function |class |interface |type |const |let )`)
)

// declarationNames matches the first line of a top-level declaration in a Go or
// TypeScript file. Other files have no pattern.
func declarationNames(path string) *regexp.Regexp {
	switch {
	case strings.HasSuffix(path, ".go"):
		return goDeclaration
	case strings.HasSuffix(path, ".ts"):
		return tsDeclaration
	}
	return nil
}
