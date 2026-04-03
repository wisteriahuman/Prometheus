package service

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

type WikilinkInfo struct {
	Slug        string `json:"slug"`
	DisplayText string `json:"displayText"`
}

type TaskInfo struct {
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	Line      int    `json:"lineNumber"`
}

type Markdown struct {
	md       goldmark.Markdown
	sanitize *bluemonday.Policy
}

func NewMarkdown() *Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").OnElements("a", "code", "span", "pre")
	policy.AllowAttrs("data-slug").OnElements("a")
	policy.AllowAttrs("type", "checked", "disabled").OnElements("input")

	return &Markdown{md: md, sanitize: policy}
}

func (m *Markdown) ToHTML(markdown string) string {
	// Remove frontmatter before rendering
	body := removeFrontmatter(markdown)

	// Convert wikilinks to HTML links before goldmark processing
	body = wikilinkRegex.ReplaceAllStringFunc(body, func(match string) string {
		sub := wikilinkRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		slug := strings.TrimSpace(sub[1])
		display := slug
		if len(sub) >= 3 && sub[2] != "" {
			display = strings.TrimSpace(sub[2])
		}
		href := "/note/" + slugToPath(slug)
		return `<a href="` + href + `" class="wikilink" data-slug="` + slug + `">` + display + `</a>`
	})

	var buf bytes.Buffer
	m.md.Convert([]byte(body), &buf)

	sanitized := m.sanitize.SanitizeBytes(buf.Bytes())
	return string(sanitized)
}

func removeFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---\n") {
		return content
	}
	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return content
	}
	return content[4+end+4:]
}

func slugToPath(slug string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(slug), " ", "-")) + ".md"
}

func ExtractWikilinks(content string) []WikilinkInfo {
	matches := wikilinkRegex.FindAllStringSubmatch(content, -1)
	var links []WikilinkInfo
	for _, m := range matches {
		slug := strings.TrimSpace(m[1])
		display := ""
		if len(m) >= 3 {
			display = strings.TrimSpace(m[2])
		}
		links = append(links, WikilinkInfo{Slug: slug, DisplayText: display})
	}
	return links
}

func ExtractTasks(content string) []TaskInfo {
	lines := strings.Split(content, "\n")
	var tasks []TaskInfo
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "- [ ] ") {
			tasks = append(tasks, TaskInfo{
				Content:   strings.TrimPrefix(trimmed, "- [ ] "),
				Completed: false,
				Line:      i + 1,
			})
		} else if strings.HasPrefix(trimmed, "- [x] ") || strings.HasPrefix(trimmed, "- [X] ") {
			text := trimmed[6:]
			tasks = append(tasks, TaskInfo{
				Content:   text,
				Completed: true,
				Line:      i + 1,
			})
		}
	}
	return tasks
}
