package helpcorpus

import (
	"fmt"
	"strings"
	"testing"
)

// checkSpans confirms that the spans are whole lines of the content, in order,
// and that they leave out only blank lines.
func checkSpans(t *testing.T, content string, spans []span) {
	t.Helper()
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	covered := make([]bool, len(lines))
	previousEnd := 0
	for _, piece := range spans {
		if piece.start <= previousEnd || piece.end < piece.start || piece.end > len(lines) {
			t.Fatalf("span %d-%d is out of order or out of range", piece.start, piece.end)
		}
		if want := strings.Join(lines[piece.start-1:piece.end], "\n"); piece.text != want {
			t.Fatalf("span %d-%d does not hold its lines", piece.start, piece.end)
		}
		if strings.TrimSpace(lines[piece.start-1]) == "" || strings.TrimSpace(lines[piece.end-1]) == "" {
			t.Fatalf("span %d-%d starts or ends with a blank line", piece.start, piece.end)
		}
		for line := piece.start; line <= piece.end; line++ {
			covered[line-1] = true
		}
		previousEnd = piece.end
	}
	for i, line := range lines {
		if !covered[i] && strings.TrimSpace(line) != "" {
			t.Fatalf("line %d is in no span: %q", i+1, line)
		}
	}
}

func TestMarkdownSectionsAreChunksWithTheirHeadings(t *testing.T) {
	section := func(title string) string {
		return "## " + title + "\n\n" + strings.Repeat("Select **Save** and wait for the approver. ", 4) + "\n\n"
	}
	content := "# Time\n\n" + strings.Repeat("A time sheet holds your entries. ", 5) + "\n\n" +
		section("Enter time") + section("Submit a week") + "### A detail\n\n" + strings.Repeat("More detail here. ", 8) + "\n"
	spans := splitChunks("ui/src/docs/time.md", content)
	checkSpans(t, content, spans)

	var headings []string
	for _, piece := range spans {
		headings = append(headings, piece.heading)
	}
	want := []string{"Time", "Time > Enter time", "Time > Submit a week", "Time > Submit a week > A detail"}
	if strings.Join(headings, "|") != strings.Join(want, "|") {
		t.Fatalf("headings = %q, want %q", headings, want)
	}
	if !strings.HasPrefix(spans[1].text, "## Enter time") {
		t.Fatalf("a section's chunk must start with its heading, got %q", spans[1].text[:20])
	}
}

func TestHeadingInACodeFenceIsNotAHeading(t *testing.T) {
	content := "# Notes\n\n" + strings.Repeat("Some text about the build. ", 8) + "\n\n```bash\n# not a heading\nmake all\n```\n\n" +
		strings.Repeat("More text. ", 20) + "\n"
	spans := splitChunks("descriptions/build.md", content)
	checkSpans(t, content, spans)
	for _, piece := range spans {
		if strings.Contains(piece.heading, "not a heading") {
			t.Fatalf("a comment in a code fence became the heading %q", piece.heading)
		}
	}
}

func TestCodeBreaksAtDeclarationsAndStaysInBounds(t *testing.T) {
	var content strings.Builder
	content.WriteString("package hooks\n\n")
	for fn := 0; fn < 12; fn++ {
		fmt.Fprintf(&content, "// step%d does one thing.\nfunc step%d(app core.App) error {\n", fn, fn)
		// Function 5 is long, so it must be cut inside.
		for line := 0; line < 12+fn%6*40; line++ {
			fmt.Fprintf(&content, "\tvalue%d := compute(app, %d)\n", line, line)
			if line%10 == 9 {
				content.WriteString("\n")
			}
		}
		content.WriteString("\treturn nil\n}\n\n")
	}
	text := content.String()
	spans := splitChunks("app/hooks/steps.go", text)
	checkSpans(t, text, spans)

	if len(spans) < 12 {
		t.Fatalf("got %d spans for 12 functions", len(spans))
	}
	insideDeclaration := 0
	for i, piece := range spans {
		if len(piece.text) > chunkMax+chunkMin {
			t.Fatalf("span %d has %d bytes", i, len(piece.text))
		}
		startsAtTop := strings.HasPrefix(piece.text, "//") || strings.HasPrefix(piece.text, "func") || strings.HasPrefix(piece.text, "package")
		switch {
		case startsAtTop && piece.heading != "":
			t.Fatalf("span %d starts a declaration, so it needs no heading, got %q", i, piece.heading)
		case !startsAtTop:
			insideDeclaration++
			if !strings.HasPrefix(piece.heading, "func step") || strings.HasSuffix(piece.heading, "{") {
				t.Fatalf("span %d starts inside a function, so it must name it, got %q", i, piece.heading)
			}
		}
	}
	if insideDeclaration == 0 {
		t.Fatal("the long functions were not cut")
	}
}

func TestOneLongBlockIsCutAtTheMaximum(t *testing.T) {
	content := strings.Repeat("SELECT a_long_column_name FROM a_table WHERE x = 1\n", 200)
	spans := splitChunks("app/routes/big.sql", content)
	checkSpans(t, content, spans)
	if len(spans) < 3 {
		t.Fatalf("got %d spans", len(spans))
	}
	for _, piece := range spans {
		if len(piece.text) > chunkMax {
			t.Fatalf("a span has %d bytes, the maximum is %d", len(piece.text), chunkMax)
		}
	}
}

func TestSmallAndOddFiles(t *testing.T) {
	if spans := splitChunks("a.go", ""); len(spans) != 0 {
		t.Fatalf("an empty file gave %d spans", len(spans))
	}
	if spans := splitChunks("a.go", "\n\n  \n"); len(spans) != 0 {
		t.Fatalf("a blank file gave %d spans", len(spans))
	}
	spans := splitChunks("a.ts", "export const a = 1;\r\n\r\nexport const b = 2;\r\n")
	if len(spans) != 1 || spans[0].start != 1 || spans[0].end != 3 || strings.Contains(spans[0].text, "\r") {
		t.Fatalf("a small Windows file gave %+v", spans)
	}
}

func TestAShortTailJoinsTheChunkBeforeIt(t *testing.T) {
	content := "package a\n\n" + strings.Repeat("var value = compute()\n", 30) + "\nfunc tail() {}\n"
	spans := splitChunks("app/a.go", content)
	checkSpans(t, content, spans)
	if len(spans) != 1 || !strings.HasSuffix(spans[0].text, "func tail() {}") {
		t.Fatalf("got %d spans", len(spans))
	}
}
