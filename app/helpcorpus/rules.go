package helpcorpus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// A rule file says in plain words what one source file does: who may do an
// action, when Turbo refuses it, the limits, and the cases that give a result a
// person does not expect. A model writes it from the code (go run
// ./cmd/helprules), and the file is in the repo, so a build needs no model.
//
// The code often has none of the words of a question. The line that approves an
// expense at submit when its owner is also its approver has no "self", "own" or
// "automatic" in it. A rule says it in those words, so the search finds it, and
// it names the lines of the code, so the model can check it.
//
// A rule file records the hash of the source that it was written from. The
// corpus takes a rule file only when its source has that hash today. A source
// that changed has no rules until someone writes them again, so a rule can be
// missing and can't be out of date.

// rulesDir is the folder of the rule files, inside the app folder. A rule file
// has the path of its source and ".md": helpcorpus/rules/app/hooks/expenses.go.md.
const rulesDir = "helpcorpus/rules"

// ruleRoots are the folders whose files get rules: where the app decides what
// a person may do. A folder of plumbing or of screens has none.
var ruleRoots = []string{
	"app/absorb/",
	"app/constants/",
	"app/cron/",
	"app/hooks/",
	"app/notifications/",
	"app/reports/",
	"app/routes/",
	"app/utilities/",
}

// ruleFiles are the folders where only the files directly inside get rules.
var ruleFiles = []string{"ui/src/lib"}

// HasRules reports if a source file of the corpus gets a rule file.
func HasRules(sourcePath string) bool {
	for _, root := range ruleRoots {
		if strings.HasPrefix(sourcePath, root) {
			return true
		}
	}
	for _, dir := range ruleFiles {
		if path.Dir(sourcePath) == dir && path.Ext(sourcePath) == ".ts" {
			return true
		}
	}
	return false
}

// Rule is one statement about what a source file does.
type Rule struct {
	// Title names the case in a few plain words. The search weighs it most.
	Title string `json:"title"`
	Text  string `json:"text"`
	// Start and End are the lines of the source that make the rule true.
	Start int `json:"start"`
	End   int `json:"end"`
}

// RuleFile is the rules of one source file.
type RuleFile struct {
	// Source is the repo path of the source file.
	Source string `json:"source"`
	// SHA256 is the hash of the source that the rules were written from.
	SHA256 string `json:"sha256"`
	// Model is the model that wrote the rules.
	Model string `json:"model"`
	Rules []Rule `json:"-"`
}

// SourceHash is the hash that a rule file records for its source.
func SourceHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

// RuleFilePath is where the rule file of a source file is on disk. appDir is
// the app folder.
func RuleFilePath(appDir, sourcePath string) string {
	return filepath.Join(appDir, filepath.FromSlash(rulesDir), filepath.FromSlash(sourcePath)+".md")
}

// RulesDir is the folder that holds every rule file.
func RulesDir(appDir string) string {
	return filepath.Join(appDir, filepath.FromSlash(rulesDir))
}

const ruleFileMark = "<!-- helprules "

// oneLine makes text safe for the file: a line of a rule can't start a heading.
func oneLine(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// Format writes the rule file. It is Markdown, so a person can read it, and it
// has one fixed shape, so ParseRuleFile reads it back.
func (f RuleFile) Format() string {
	meta, _ := json.Marshal(f)
	var out strings.Builder
	out.WriteString("<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->\n")
	fmt.Fprintf(&out, "%s%s -->\n\n# Rules in %s\n", ruleFileMark, meta, f.Source)
	if len(f.Rules) == 0 {
		out.WriteString("\nThis file has no rule that a person who uses Turbo can observe.\n")
	}
	for _, rule := range f.Rules {
		fmt.Fprintf(&out, "\n## %s\n\nLines %d-%d\n\n%s\n", oneLine(rule.Title), rule.Start, rule.End, oneLine(rule.Text))
	}
	return out.String()
}

var ruleLines = regexp.MustCompile(`^Lines (\d+)-(\d+)$`)

// ParseRuleFile reads what Format wrote.
func ParseRuleFile(text string) (RuleFile, error) {
	var file RuleFile
	var rule *Rule
	var body []string
	flush := func() error {
		if rule == nil {
			return nil
		}
		rule.Text = strings.TrimSpace(strings.Join(body, " "))
		if rule.Title == "" || rule.Text == "" || rule.Start < 1 || rule.End < rule.Start {
			return fmt.Errorf("the rule %q has no text or no lines", rule.Title)
		}
		file.Rules = append(file.Rules, *rule)
		rule, body = nil, nil
		return nil
	}
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, ruleFileMark):
			meta := strings.TrimSuffix(strings.TrimPrefix(line, ruleFileMark), " -->")
			if err := json.Unmarshal([]byte(meta), &file); err != nil {
				return file, fmt.Errorf("the helprules line is not valid: %w", err)
			}
		case strings.HasPrefix(line, "## "):
			if err := flush(); err != nil {
				return file, err
			}
			rule = &Rule{Title: strings.TrimSpace(strings.TrimPrefix(line, "## "))}
		case rule != nil && rule.Start == 0 && ruleLines.MatchString(line):
			match := ruleLines.FindStringSubmatch(line)
			rule.Start, _ = strconv.Atoi(match[1])
			rule.End, _ = strconv.Atoi(match[2])
		case rule != nil && strings.TrimSpace(line) != "":
			body = append(body, strings.TrimSpace(line))
		}
	}
	if err := flush(); err != nil {
		return file, err
	}
	if file.Source == "" || file.SHA256 == "" {
		return file, fmt.Errorf("the file has no helprules line")
	}
	return file, nil
}

// ReadRuleFile reads the rule file of a source file. ok is false when there is
// none.
func ReadRuleFile(appDir, sourcePath string) (file RuleFile, ok bool, err error) {
	data, err := os.ReadFile(RuleFilePath(appDir, sourcePath))
	if os.IsNotExist(err) {
		return file, false, nil
	}
	if err != nil {
		return file, false, err
	}
	file, err = ParseRuleFile(string(data))
	if err != nil {
		return file, false, fmt.Errorf("%s: %w", RuleFilePath(appDir, sourcePath), err)
	}
	return file, true, nil
}

// SourceFile is one source file that gets rules.
type SourceFile struct {
	Path    string
	Content string
}

// RuleSources returns the source files that get rules, in the order of their
// paths. They are files of the corpus, so what the corpus leaves out has no
// rules.
func RuleSources(options Options) ([]SourceFile, error) {
	var out []SourceFile
	for _, src := range sources {
		if src.kind != KindCode {
			continue
		}
		files, err := readSource(options, src, &Stats{})
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if HasRules(f.path) {
				out = append(out, SourceFile{Path: f.path, Content: f.content})
			}
		}
	}
	return out, nil
}

// appDirOf returns the app folder for the options.
func appDirOf(options Options) string {
	if options.AppDir != "" {
		return options.AppDir
	}
	return filepath.Join(options.RepoDir, "app")
}

// ruleChunks returns the chunks for the rules of one source file: one chunk for
// each rule, with the path and the lines of the code that it describes. So a
// source link for a rule opens the code. It returns nothing when the file has
// no rules, and stale is true when the rules are for an older source.
func ruleChunks(options Options, f file) (chunks []Chunk, stale bool, err error) {
	rules, ok, err := ReadRuleFile(appDirOf(options), f.path)
	if err != nil || !ok {
		return nil, false, err
	}
	if rules.SHA256 != SourceHash(f.content) {
		return nil, true, nil
	}
	// A file with no rules still has its rule file, so it is not reported as
	// a file that needs one.
	chunks = []Chunk{}
	for i, rule := range rules.Rules {
		// The chat treats a commit as final. See handbookBans.
		if handbookBans.MatchString(rule.Title + " " + rule.Text) {
			continue
		}
		chunks = append(chunks, Chunk{
			ID:       fmt.Sprintf("rule:%s:%d:%d", f.path, i+1, rule.Start),
			Path:     f.path,
			Kind:     KindRules,
			Start:    rule.Start,
			End:      rule.End,
			Heading:  rule.Title,
			Requires: requiresFor(f.path, rule.Title+" "+rule.Text),
			Text:     rule.Text,
		})
	}
	return chunks, false, nil
}
