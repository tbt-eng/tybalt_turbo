package helpcorpus

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// The Docker build writes data/corpus.json before it compiles the app, so the
// binary holds the corpus of its own commit. The folder is not in pb_public,
// and no route serves it. Only .gitkeep is in git.
//
//go:embed all:data
var embedded embed.FS

// ErrNotBuilt means that this binary holds no corpus.
var ErrNotBuilt = errors.New("the Help chat corpus is not built: run `go run ./cmd/helpcorpus` in app/")

// Embedded returns the corpus that the build put in this binary.
func Embedded() (*Corpus, error) {
	data, err := embedded.ReadFile("data/corpus.json")
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrNotBuilt
	}
	if err != nil {
		return nil, err
	}
	return Decode(data)
}

// Open returns the corpus for a server to use.
//
// In development, pass the repo folder. If the sources are there, Open builds
// the corpus from them, so the answers match the files on disk. In production
// there is no repo, so Open returns the corpus in the binary.
func Open(repoDir string) (*Corpus, error) {
	if repoDir != "" {
		if _, err := os.Stat(filepath.Join(repoDir, filepath.FromSlash(docsRoot))); err == nil {
			corpus, _, err := Generate(Options{RepoDir: repoDir})
			return corpus, err
		}
	}
	return Embedded()
}
