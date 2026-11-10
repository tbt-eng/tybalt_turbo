package helpcorpus

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	// The pure-Go SQLite driver that PocketBase also uses. It registers "sqlite".
	_ "modernc.org/sqlite"
)

// Limits on what one call returns, so a tool result stays a small part of a
// prompt.
const (
	defaultSearchLimit = 8
	maxSearchLimit     = 20
	maxReadChunks      = 8
	maxQueryTerms      = 24
)

// Index searches a corpus in memory. Every method that returns chunk text takes
// the user's claims and leaves out the chunks that those claims don't open.
// There is no way to get a chunk from an Index without claims.
type Index struct {
	corpus *Corpus
	byPath map[string][]int

	// One connection holds the in-memory database, so the lock also keeps the
	// database alive: a second connection would see an empty one.
	mu   sync.Mutex
	db   *sql.DB
	conn *sql.Conn
}

// Hit is one search result.
type Hit struct {
	Chunk
	// Score is higher for a better match. It only orders the hits of one search.
	Score float64
}

// SearchOptions narrows a search.
type SearchOptions struct {
	// Limit is the most hits to return. Zero means 8, and the maximum is 20.
	Limit int
	// Kinds keeps only chunks of these kinds. Empty means all kinds.
	Kinds []string
}

// NewIndex builds the search index for a corpus. It takes about 200 ms for the
// app's own corpus.
func NewIndex(corpus *Corpus) (*Index, error) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	index := &Index{corpus: corpus, byPath: map[string][]int{}, db: db, conn: conn}
	if err := index.fill(ctx); err != nil {
		index.Close()
		return nil, fmt.Errorf("build the Help chat index: %w", err)
	}
	return index, nil
}

func (ix *Index) fill(ctx context.Context) error {
	// The porter tokenizer makes "approves" and "approval" one word.
	if _, err := ix.conn.ExecContext(ctx,
		`CREATE VIRTUAL TABLE chunks USING fts5(path, heading, body, parts, tokenize = 'porter unicode61')`,
	); err != nil {
		return err
	}
	tx, err := ix.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	insert, err := tx.PrepareContext(ctx, `INSERT INTO chunks(rowid, path, heading, body, parts) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer insert.Close()
	for i, chunk := range ix.corpus.Chunks {
		parts := identifierParts(chunk.Path + " " + chunk.Heading + " " + chunk.Text)
		body := chunk.Text
		if chunk.Asks != "" {
			body += "\n" + chunk.Asks
		}
		if _, err := insert.ExecContext(ctx, i+1, chunk.Path, chunk.Heading, body, parts); err != nil {
			return err
		}
		ix.byPath[chunk.Path] = append(ix.byPath[chunk.Path], i)
	}
	return tx.Commit()
}

// Close frees the index.
func (ix *Index) Close() error {
	ix.mu.Lock()
	defer ix.mu.Unlock()
	ix.conn.Close()
	return ix.db.Close()
}

// Handbook returns the text that goes into every prompt. It holds only text for
// everyone, so it takes no claims.
func (ix *Index) Handbook() string { return ix.corpus.Handbook }

// Hash identifies the corpus that the index holds.
func (ix *Index) Hash() string { return ix.corpus.Hash }

// Search returns the chunks that best match the query, best first, and only
// chunks that the claims open. The query is plain words: a chunk matches if it
// holds any of them, and holds more of them to rank higher.
func (ix *Index) Search(claims Claims, query string, options SearchOptions) ([]Hit, error) {
	match := matchExpression(query)
	if match == "" {
		return nil, nil
	}
	limit := options.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	limit = min(limit, maxSearchLimit)

	ix.mu.Lock()
	defer ix.mu.Unlock()
	// bm25 weighs a match in the heading most, then the path, then the text.
	// A lower rank is a better match.
	rows, err := ix.conn.QueryContext(context.Background(),
		`SELECT rowid, bm25(chunks, 2.0, 5.0, 1.0, 0.5) AS rank FROM chunks WHERE chunks MATCH ? ORDER BY rank`,
		match,
	)
	if err != nil {
		return nil, fmt.Errorf("search the Help chat index: %w", err)
	}
	defer rows.Close()

	var hits []Hit
	for rows.Next() && len(hits) < limit {
		var rowid int
		var rank float64
		if err := rows.Scan(&rowid, &rank); err != nil {
			return nil, err
		}
		chunk := ix.corpus.Chunks[rowid-1]
		// The filter runs before the limit, so a chunk that the user can't read
		// never takes the place of one that they can.
		if !chunk.Allows(claims) {
			continue
		}
		if len(options.Kinds) > 0 && !slices.Contains(options.Kinds, chunk.Kind) {
			continue
		}
		hits = append(hits, Hit{Chunk: chunk, Score: -rank})
	}
	return hits, rows.Err()
}

// ChunkByID returns one chunk, if the claims open it.
func (ix *Index) ChunkByID(claims Claims, id string) (Chunk, bool) {
	for _, i := range ix.byPath[pathOfID(id)] {
		if chunk := ix.corpus.Chunks[i]; chunk.ID == id && chunk.Allows(claims) {
			return chunk, true
		}
	}
	return Chunk{}, false
}

// pathOfID returns the path in a chunk id: "path:start", or
// "rule:path:number:start" for a rule.
func pathOfID(id string) string {
	parts := strings.Split(id, ":")
	if parts[0] == "rule" && len(parts) >= 4 {
		return strings.Join(parts[1:len(parts)-2], ":")
	}
	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], ":")
	}
	return id
}

// Read returns the chunks of one file that hold the lines from..to, in order,
// and only chunks that the claims open. It returns at most 8 chunks. A line
// number below 1 means the start or the end of the file.
func (ix *Index) Read(claims Claims, path string, from, to int) []Chunk {
	var chunks []Chunk
	for _, i := range ix.byPath[path] {
		chunk := ix.corpus.Chunks[i]
		// A rule has the path and the lines of its code. A read is for the code.
		if (to >= 1 && chunk.Start > to) || chunk.End < from || !chunk.Allows(claims) || chunk.Kind == KindRules {
			continue
		}
		chunks = append(chunks, chunk)
		if len(chunks) == maxReadChunks {
			break
		}
	}
	return chunks
}

var (
	wordPattern = regexp.MustCompile(`[\pL\pN_]+`)
	camelBreak  = regexp.MustCompile(`([\p{Ll}\pN])(\p{Lu})`)
	camelWord   = regexp.MustCompile(`\pL*\p{Ll}\p{Lu}[\pL\pN]*`)
)

// identifierParts returns the words inside camelCase names, such as "clean Time
// Entry" for cleanTimeEntry. The index holds them beside the text, so a search
// for "time entry" finds the function. Names with underscores need no help,
// because the tokenizer splits them.
func identifierParts(text string) string {
	names := camelWord.FindAllString(text, -1)
	for i, name := range names {
		names[i] = camelBreak.ReplaceAllString(name, "$1 $2")
	}
	return strings.Join(names, " ")
}

// matchExpression turns any text into a safe FTS5 query: quoted words joined by
// OR. Text from the model or a user never reaches FTS5 as query syntax.
func matchExpression(query string) string {
	var terms []string
	add := func(word string) {
		term := `"` + strings.ToLower(word) + `"`
		if len(terms) < maxQueryTerms && !slices.Contains(terms, term) {
			terms = append(terms, term)
		}
	}
	for _, word := range wordPattern.FindAllString(query, -1) {
		add(word)
		if split := camelBreak.ReplaceAllString(word, "$1 $2"); split != word {
			for _, part := range strings.Fields(split) {
				add(part)
			}
		}
	}
	return strings.Join(terms, " OR ")
}
