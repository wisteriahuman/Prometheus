package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type sampleNote struct {
	path    string
	title   string
	tags    []string
	content string
}

var sampleNotes = []sampleNote{
	{
		path:  "welcome.md",
		title: "Prometheusへようこそ",
		tags:  []string{"getting-started", "prometheus"},
		content: `
# Prometheusへようこそ

**Prometheus** はMarkdownベースのノートアプリです。

## 基本操作

| キー | 機能 |
|------|------|
| ` + "`⌘K`" + ` | コマンドパレット |
| ` + "`⌘P`" + ` | ファイル切替 |
| ` + "`⌘S`" + ` | 保存 |
| ` + "`⌘D`" + ` | デイリーノート |
| ` + "`⌘G`" + ` | グラフビュー |
| ` + "`⌘.`" + ` | テーマ切替 |
| ` + "`⌘K`" + ` | ショートカット一覧（コマンドパレットから） |

## はじめに

- ` + "`[[wikilink]]`" + ` でノート同士をリンクできます
- サイドバーの ` + "`+`" + ` から新規ノート作成
- タグはノートヘッダーの ` + "`+`" + ` から追加
- 右クリックでノートのリネーム・削除

詳しくは [[setup]] をご覧ください。
[[ideas]] も参考にどうぞ。

## Mermaid図のサンプル

` + "```mermaid" + `
graph LR
  A[ノート作成] --> B[リンクを貼る]
  B --> C[グラフで可視化]
  C --> D[知識が繋がる]
` + "```" + `

## AI / Agent連携

vaultディレクトリでClaude Code や agent 対応ツールを起動すると、ノート操作をAIが支援します:

` + "```bash" + `
cd ~/my-notes && claude
cd ~/my-notes && codex
` + "```" + `

設定ファイル: ` + "`CLAUDE.md`" + ` ` + "`AGENTS.md`" + `  
補助ディレクトリ: ` + "`.claude/`" + ` ` + "`.agents/`" + `

> *ノートは .md ファイルとして保存されます。Neovimで直接編集可能です。*
`,
	},
	{
		path:  "setup.md",
		title: "セットアップ",
		tags:  []string{"getting-started", "guide"},
		content: `
# セットアップ

[[welcome]] の使い方ガイドです。

## ノートの保存先

起動時に指定したディレクトリに ` + "`.md`" + ` ファイルとして保存されます。

` + "```bash" + `
prm dev ~/my-notes        # ノートの保存先を指定して起動
prm dev ~/work -p 3001    # 別ポートで起動
` + "```" + `

- Neovimやお好みのエディタで直接編集できます
- 変更はPrometheusに自動反映されます
- Gitで管理すればバージョン管理も可能です

## ファイル名とタイトル

- **ファイル名**: ` + "`setup.md`" + ` のようにURLやリンクで使う名前
- **タイトル**: フロントマターの ` + "`title`" + ` で自由に設定。日本語OK
- ファイル名は英語、タイトルは日本語が推奨

## リンクの書き方

` + "```markdown" + `
[[ファイル名]]           → そのファイルへリンク
[[ファイル名|表示テキスト]] → 別名でリンク
` + "```" + `

## タグ

ノートヘッダーの ` + "`+`" + ` ボタンからタグを追加できます。
グラフビューではタグごとにノートがグループ化されます。

[[welcome]] | [[ideas]]
`,
	},
	{
		path:  "ideas.md",
		title: "アイデア",
		tags:  []string{"ideas"},
		content: `
# アイデア

ひらめきを逃さないためのノート。

## やりたいこと

- [ ] プロジェクトのアイデアを整理する
- [ ] 読書メモをまとめる
- [ ] 学習ログをつける

## フロー

` + "```mermaid" + `
graph TD
  A[ひらめき] --> B{重要？}
  B -->|Yes| C[ノートに書く]
  B -->|No| D[忘れてOK]
  C --> E[タグをつける]
  E --> F[リンクで繋げる]
` + "```" + `

## メモ

ここに自由にメモを書いてください。

[[welcome]] | [[setup]]
`,
	},
}

func InitVault(vault *Vault) bool {
	entries, err := os.ReadDir(vault.Path())
	if err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") || (!strings.HasPrefix(e.Name(), ".") && e.IsDir()) {
				return false // Already has content
			}
		}
	}

	// Create .prometheus config
	configDir := filepath.Join(vault.Path(), ".prometheus", "themes")
	os.MkdirAll(configDir, 0o755)

	configData, _ := json.MarshalIndent(map[string]interface{}{
		"name":              "My Notes",
		"theme":             "prometheus",
		"dailyNoteTemplate": "---\ntitle: \"{{date}}\"\ntags: [daily]\n---\n\n# {{dateJa}}\n\n## タスク\n\n- [ ] \n\n## メモ\n\n",
	}, "", "  ")
	os.WriteFile(filepath.Join(vault.Path(), ".prometheus", "config.json"), configData, 0o644)

	// Create sample notes
	for _, note := range sampleNotes {
		fm := NoteFrontmatter{
			ID:       newULID(),
			Title:    note.title,
			Created:  fmt.Sprintf("%s", timeNow()),
			Modified: fmt.Sprintf("%s", timeNow()),
			Tags:     note.tags,
		}
		content := serializeFrontmatter(fm) + note.content
		fullPath := filepath.Join(vault.Path(), note.path)
		os.WriteFile(fullPath, []byte(content), 0o644)
	}

	// Create AI workspace config
	ensureAgentWorkspace(vault.Path())

	// Create vault .gitignore
	gitignorePath := filepath.Join(vault.Path(), ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		os.WriteFile(gitignorePath, []byte(vaultGitignore), 0o644)
	}

	return true
}

func ensureAgentWorkspace(vaultPath string) {
	ensureWorkspaceAssets(vaultPath, rootWorkspaceGuides())
	ensureWorkspaceAssets(filepath.Join(vaultPath, ".claude"), claudeWorkspaceAssets())
	ensureWorkspaceAssets(filepath.Join(vaultPath, ".agents"), genericAgentWorkspaceAssets())
}

func EnsureAgentWorkspacePublic(vaultPath string) {
	ensureAgentWorkspace(vaultPath)
}

func writeIfMissing(path string, content []byte) {
	if _, err := os.Stat(path); err == nil {
		return
	}
	os.WriteFile(path, content, 0o644)
}

func ensureWorkspaceAssets(basePath string, assets []workspaceAsset) {
	for _, asset := range assets {
		fullPath := filepath.Join(basePath, asset.RelativePath)
		os.MkdirAll(filepath.Dir(fullPath), 0o755)
		writeIfMissing(fullPath, []byte(asset.Content))
	}
}

func timeNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}
