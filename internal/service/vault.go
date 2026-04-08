package service

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/yaml.v3"
)

type NoteFrontmatter struct {
	ID       string   `yaml:"id" json:"id"`
	Title    string   `yaml:"title" json:"title"`
	Created  string   `yaml:"created" json:"created"`
	Modified string   `yaml:"modified" json:"modified"`
	Tags     []string `yaml:"tags" json:"tags"`
	Theme    *string  `yaml:"theme,omitempty" json:"theme"`
}

type VaultNote struct {
	Path        string          `json:"path"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	RawContent  string          `json:"-"`
	Frontmatter NoteFrontmatter `json:"frontmatter"`
	Checksum    string          `json:"checksum"`
}

type FileEntry struct {
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	IsDirectory bool         `json:"isDirectory"`
	Children    []*FileEntry `json:"children,omitempty"`
}

type Vault struct {
	path string
}

func NewVault(vaultPath string) *Vault {
	os.MkdirAll(vaultPath, 0o755)
	return &Vault{path: vaultPath}
}

func (v *Vault) Path() string {
	return v.path
}

func computeChecksum(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

func newULID() string {
	return ulid.Make().String()
}

// ParseFrontmatterPublic is a public wrapper around parseFrontmatter.
func ParseFrontmatterPublic(raw []byte) (NoteFrontmatter, string) {
	return parseFrontmatter(raw)
}

func parseFrontmatter(raw []byte) (NoteFrontmatter, string) {
	content := string(raw)

	if !strings.HasPrefix(content, "---\n") {
		return NoteFrontmatter{}, content
	}

	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return NoteFrontmatter{}, content
	}

	fmRaw := content[4 : 4+end]
	body := content[4+end+4:] // skip \n---\n

	var fm NoteFrontmatter
	yaml.Unmarshal([]byte(fmRaw), &fm)

	if fm.Tags == nil {
		fm.Tags = []string{}
	}

	return fm, body
}

func serializeFrontmatter(fm NoteFrontmatter) string {
	// Build YAML manually to control output
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("id: \"%s\"\n", fm.ID))
	b.WriteString(fmt.Sprintf("title: \"%s\"\n", fm.Title))
	b.WriteString(fmt.Sprintf("created: \"%s\"\n", fm.Created))
	b.WriteString(fmt.Sprintf("modified: \"%s\"\n", fm.Modified))

	if len(fm.Tags) > 0 {
		tags := make([]string, len(fm.Tags))
		for i, t := range fm.Tags {
			tags[i] = t
		}
		b.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(tags, ", ")))
	} else {
		b.WriteString("tags: []\n")
	}

	if fm.Theme != nil && *fm.Theme != "" {
		b.WriteString(fmt.Sprintf("theme: %s\n", *fm.Theme))
	}

	b.WriteString("---\n")
	return b.String()
}

func (v *Vault) ReadNote(notePath string) (*VaultNote, error) {
	fullPath := filepath.Join(v.path, notePath)
	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	checksum := computeChecksum(raw)
	fm, body := parseFrontmatter(raw)

	if fm.ID == "" {
		fm.ID = newULID()
	}
	if fm.Title == "" {
		fm.Title = strings.TrimSuffix(filepath.Base(notePath), ".md")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if fm.Created == "" {
		fm.Created = now
	}
	if fm.Modified == "" {
		fm.Modified = now
	}

	return &VaultNote{
		Path:        notePath,
		Title:       fm.Title,
		Content:     body,
		RawContent:  string(raw),
		Frontmatter: fm,
		Checksum:    checksum,
	}, nil
}

func (v *Vault) WriteNote(notePath, body string, fm NoteFrontmatter) (*VaultNote, error) {
	fullPath := filepath.Join(v.path, notePath)
	os.MkdirAll(filepath.Dir(fullPath), 0o755)

	fm.Modified = time.Now().UTC().Format(time.RFC3339)

	rawContent := serializeFrontmatter(fm) + body
	if err := os.WriteFile(fullPath, []byte(rawContent), 0o644); err != nil {
		return nil, err
	}

	checksum := computeChecksum([]byte(rawContent))
	return &VaultNote{
		Path:        notePath,
		Title:       fm.Title,
		Content:     body,
		RawContent:  rawContent,
		Frontmatter: fm,
		Checksum:    checksum,
	}, nil
}

func (v *Vault) CreateNote(notePath string, title string) (*VaultNote, error) {
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(notePath), ".md")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	fm := NoteFrontmatter{
		ID:       newULID(),
		Title:    title,
		Created:  now,
		Modified: now,
		Tags:     []string{},
	}
	body := "\n# " + title + "\n\n"
	return v.WriteNote(notePath, body, fm)
}

func (v *Vault) DeleteNote(notePath string) error {
	return os.Remove(filepath.Join(v.path, notePath))
}

func (v *Vault) NoteExists(notePath string) bool {
	_, err := os.Stat(filepath.Join(v.path, notePath))
	return err == nil
}

// isHiddenNote returns true for config files that should not appear in the UI.
func isHiddenNote(name string) bool {
	upper := strings.ToUpper(name)
	return upper == "CLAUDE.MD" || upper == "CLAUDE.LOCAL.MD" || upper == "AGENTS.MD" || upper == "AGENTS.LOCAL.MD"
}

func (v *Vault) ListNotes() ([]string, error) {
	var notes []string
	err := filepath.WalkDir(v.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") && !strings.HasPrefix(d.Name(), ".") && !isHiddenNote(d.Name()) {
			rel, _ := filepath.Rel(v.path, path)
			notes = append(notes, rel)
		}
		return nil
	})
	sort.Strings(notes)
	return notes, err
}

func (v *Vault) GetFileTree(dirPath string) ([]*FileEntry, error) {
	targetPath := filepath.Join(v.path, dirPath)
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, err
	}

	var result []*FileEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryRelPath := entry.Name()
		if dirPath != "" {
			entryRelPath = filepath.Join(dirPath, entry.Name())
		}

		if entry.IsDir() {
			children, _ := v.GetFileTree(entryRelPath)
			result = append(result, &FileEntry{
				Name:        entry.Name(),
				Path:        entryRelPath,
				IsDirectory: true,
				Children:    children,
			})
		} else if strings.HasSuffix(entry.Name(), ".md") && !isHiddenNote(entry.Name()) {
			result = append(result, &FileEntry{
				Name:        entry.Name(),
				Path:        entryRelPath,
				IsDirectory: false,
			})
		}
	}

	// Sort: directories first, then alphabetically
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDirectory && !result[j].IsDirectory {
			return true
		}
		if !result[i].IsDirectory && result[j].IsDirectory {
			return false
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}
