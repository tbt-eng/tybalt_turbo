// Package helpcorpus builds the corpus that Help chat searches: the user docs,
// the design notes and the app's code, cut into chunks. Each chunk says who may
// read it.
//
// A deploy builds the corpus from the repo with `go run ./cmd/helpcorpus`. The
// same input always gives the same bytes, so a build identifies its corpus.
package helpcorpus

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// FormatVersion changes when the shape of the corpus file changes.
const FormatVersion = 1

// The kinds of source that a chunk comes from.
const (
	KindDocs  = "docs"  // user docs in ui/src/docs
	KindNotes = "notes" // design notes in descriptions/
	KindCode  = "code"  // app and UI source
	KindRules = "rules" // what a source file does, in plain words (see rules.go)
)

// Chunk is one piece of one file.
type Chunk struct {
	// ID is the path and the first line, such as "app/hooks/expenses.go:120".
	ID   string `json:"id"`
	Path string `json:"path"`
	Kind string `json:"kind"`
	// Start and End are line numbers in the file. Both count from 1, and End is
	// the last line of the chunk.
	Start int `json:"start"`
	End   int `json:"end"`
	// Heading says where the chunk is in its file. For Markdown it is the
	// headings that cover the chunk, such as "Time > Submit a week". For a chunk
	// that starts inside a Go or TypeScript declaration, it is that
	// declaration's first line.
	Heading string `json:"heading,omitempty"`
	// Requires says who may read the chunk. See Allows.
	Requires [][]string `json:"requires"`
	Text     string     `json:"text"`
	// Asks is more text for the search to match, which a reader of the chunk
	// does not get. A rule has it: see Rule.Asks.
	Asks string `json:"asks,omitempty"`
}

// Claims are the names of the claims that a user holds. A nil value is a user
// with no claims.
type Claims map[string]bool

// Allows reports if a user with these claims may read the chunk. The user must
// hold at least one claim from every group in Requires. A chunk with no groups
// is for everyone.
func (c Chunk) Allows(claims Claims) bool {
	for _, group := range c.Requires {
		held := false
		for _, claim := range group {
			if claims[claim] {
				held = true
				break
			}
		}
		if !held {
			return false
		}
	}
	return true
}

// ForEveryone reports if any signed-in user may read the chunk.
func (c Chunk) ForEveryone() bool {
	return len(c.Requires) == 0
}

// Corpus is everything that Help chat can read.
type Corpus struct {
	Format int `json:"format"`
	// Hash identifies the handbook and the chunks. It is a SHA-256 in hex.
	Hash string `json:"hash"`
	// Handbook is the text that goes into every prompt. It holds only text that
	// every signed-in user may read.
	Handbook string  `json:"handbook"`
	Chunks   []Chunk `json:"chunks"`
}

// Encode writes the corpus as JSON with one chunk on each line, so a diff of two
// corpus files is readable.
func (c *Corpus) Encode() ([]byte, error) {
	var out bytes.Buffer
	handbook, err := marshal(c.Handbook)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(&out, "{\"format\":%d,\"hash\":%q,\"handbook\":%s,\"chunks\":[", c.Format, c.Hash, handbook)
	for i, chunk := range c.Chunks {
		line, err := marshal(chunk)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			out.WriteByte(',')
		}
		out.WriteByte('\n')
		out.Write(line)
	}
	out.WriteString("\n]}\n")
	return out.Bytes(), nil
}

// Decode reads a corpus file and checks its format and its hash.
func Decode(data []byte) (*Corpus, error) {
	var corpus Corpus
	if err := json.Unmarshal(data, &corpus); err != nil {
		return nil, fmt.Errorf("read corpus: %w", err)
	}
	if corpus.Format != FormatVersion {
		return nil, fmt.Errorf("corpus format is %d, this build reads %d", corpus.Format, FormatVersion)
	}
	hash, err := hashOf(corpus.Handbook, corpus.Chunks)
	if err != nil {
		return nil, err
	}
	if hash != corpus.Hash {
		return nil, fmt.Errorf("corpus hash does not match its content")
	}
	return &corpus, nil
}

func hashOf(handbook string, chunks []Chunk) (string, error) {
	sum := sha256.New()
	sum.Write([]byte(handbook))
	for _, chunk := range chunks {
		line, err := marshal(chunk)
		if err != nil {
			return "", err
		}
		sum.Write([]byte{'\n'})
		sum.Write(line)
	}
	return hex.EncodeToString(sum.Sum(nil)), nil
}

// marshal writes JSON without the HTML escapes that json.Marshal adds, because
// the code in a chunk is full of < and >.
func marshal(value any) ([]byte, error) {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return bytes.TrimRight(out.Bytes(), "\n"), nil
}
