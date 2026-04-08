package service

import (
	"fmt"
	"strings"
)

type workspaceAsset struct {
	RelativePath string
	Content      string
}

func rootWorkspaceGuides() []workspaceAsset {
	return []workspaceAsset{
		{RelativePath: "CLAUDE.md", Content: buildClaudeWorkspaceGuide()},
		{RelativePath: "AGENTS.md", Content: buildAgentWorkspaceGuide()},
	}
}

func claudeWorkspaceAssets() []workspaceAsset {
	return []workspaceAsset{
		{RelativePath: "settings.json", Content: claudeSettings},
		{RelativePath: "rules/frontmatter.md", Content: frontmatterRule},
		{RelativePath: "rules/wikilink.md", Content: wikilinkRule},
		{RelativePath: "skills/create-note/SKILL.md", Content: buildClaudeSkill("create-note", "フロントマター付きの新規ノートを作成", "[title]", "Read Write Glob", createNoteSkillBody(true))},
		{RelativePath: "skills/daily/SKILL.md", Content: buildClaudeSkill("daily", "今日のデイリーノートを作成または開く", "", "Read Write Glob", dailySkillBody(true))},
		{RelativePath: "skills/link-check/SKILL.md", Content: buildClaudeSkill("link-check", "壊れたwikilinkを検出し修復提案する", "", "Read Glob Grep", linkCheckSkillBody(true))},
		{RelativePath: "agents/note-explorer/AGENT.md", Content: noteExplorerAgent},
	}
}

func genericAgentWorkspaceAssets() []workspaceAsset {
	return []workspaceAsset{
		{RelativePath: "rules/frontmatter.md", Content: frontmatterRule},
		{RelativePath: "rules/wikilink.md", Content: wikilinkRule},
		{RelativePath: "skills/create-note/SKILL.md", Content: buildAgentSkill("create-note", createNoteSkillBody(false))},
		{RelativePath: "skills/daily/SKILL.md", Content: buildAgentSkill("daily", dailySkillBody(false))},
		{RelativePath: "skills/link-check/SKILL.md", Content: buildAgentSkill("link-check", linkCheckSkillBody(false))},
		{RelativePath: "agents/note-explorer/AGENT.md", Content: noteExplorerAgent},
	}
}

func buildClaudeWorkspaceGuide() string {
	return buildWorkspaceGuide(".claude/rules", "CLAUDE.local.md", "操作のコツ", []string{
		fmt.Sprintf("%s を省略（アプリが自動付与）", bt("id")),
		"ファイル名はkebab-case英語、タイトルは日本語推奨",
		fmt.Sprintf("%s のテーマやテンプレートは変更可能", bt(".prometheus/config.json")),
		fmt.Sprintf("個人設定は %s を作成（gitignore済み）", bt("CLAUDE.local.md")),
	})
}

func buildAgentWorkspaceGuide() string {
	return buildWorkspaceGuide(".agents/rules", "AGENTS.local.md", "エージェント向けヒント", []string{
		fmt.Sprintf("追加ルールは %s に置く", bt(".agents/rules/")),
		fmt.Sprintf("再利用可能な手順は %s に置く", bt(".agents/skills/")),
		fmt.Sprintf("補助エージェント定義は %s に置く", bt(".agents/agents/")),
		fmt.Sprintf("個人設定は %s を作成（gitignore済み）", bt("AGENTS.local.md")),
	})
}

func buildWorkspaceGuide(ruleDir, localConfig, tipsHeading string, tips []string) string {
	return fmt.Sprintf(`# Prometheus Vault

このディレクトリ（CWD）がvaultのルート。全ノートのパスはここからの相対パス。
Go製のWebサーバー %s がこのディレクトリを監視し、Web UIを提供する。

## ディレクトリ構造

- %s — ノート本体。YAMLフロントマター必須
- %s — カテゴリ別ノート。自由にディレクトリ作成可
- %s — デイリーノート
- %s — 画像・PDFなどのアップロードファイル
- %s — vault設定。読み書き可能
- %s — SQLite FTS5インデックス。**変更禁止**

ノートのパスはvaultルートからの相対パス（例: %s）。
Web UIのURLも同じパスになる: %s

## フロントマター

詳細は %s を参照。例:

%s

- %s: ULID形式。**変更禁止**。省略すればアプリが自動生成
- %s: RFC3339 UTC
- %s: 配列。グラフビューでクラスタ化される
- %s: ノート個別テーマ（省略可）

## Wikilink

%s または %s
slugは小文字化、スペース→ハイフン変換。%s → %s
詳細は %s を参照。

## 対応記法

- **Mermaid**: %s コードブロックで図を描画
- **LaTeX数式**: %s / %s
- **GFM**: テーブル、タスクリスト、取り消し線、脚注

## 禁止事項

- %s を変更・削除しない
- フロントマターの %s を変更しない
- %s のファイルを削除しない

## %s

%s
`,
		bt("prm"),
		bt("*.md"),
		bt("サブディレクトリ/*.md"),
		bt("daily/YYYY-MM-DD.md"),
		bt("assets/"),
		bt(".prometheus/config.json"),
		bt(".prometheus/data.db"),
		bt("resources/design-md.md"),
		bt("/note/resources/design-md.md"),
		bt(ruleDir+"/frontmatter.md"),
		fenced("yaml", `---
id: "01KNE4B96D..."
title: "ノートタイトル"
created: "2026-04-06T09:00:00Z"
modified: "2026-04-06T09:00:00Z"
tags: [tag1, tag2]
---`),
		bt("id"),
		bt("created")+"/"+bt("modified"),
		bt("tags"),
		bt("theme"),
		bt("[[ファイル名]]"),
		bt("[[ファイル名|表示テキスト]]"),
		bt("[[Setup]]"),
		bt("setup.md"),
		bt(ruleDir+"/wikilink.md"),
		bt("```mermaid"),
		bt("$...$"),
		bt("$$...$$"),
		bt(".prometheus/data.db"),
		bt("id"),
		bt("assets/"),
		tipsHeading,
		bulletList(tips),
	)
}

func buildClaudeSkill(name, description, argumentHint, allowedTools, body string) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("name: %s\n", name))
	b.WriteString(fmt.Sprintf("description: %s\n", description))
	if argumentHint != "" {
		b.WriteString(fmt.Sprintf("argument-hint: %q\n", argumentHint))
	}
	b.WriteString(fmt.Sprintf("allowed-tools: %s\n", allowedTools))
	b.WriteString("---\n\n")
	b.WriteString(body)
	return b.String()
}

func buildAgentSkill(name, body string) string {
	return fmt.Sprintf("# %s\n\n%s", name, body)
}

func createNoteSkillBody(includeClaudeDetails bool) string {
	body := `タイトルからフロントマター付きの新規ノートを作成する。

## 手順

1. タイトルからkebab-caseのファイル名を生成
2. 同名ファイルが存在しないか確認
3. 以下のテンプレートで作成

` + fenced("markdown", `---
title: "$ARGUMENTS"
created: "<現在のRFC3339 UTC>"
modified: "<現在のRFC3339 UTC>"
tags: []
---

# $ARGUMENTS
`) + `

4. ` + bt("id") + ` は省略する`

	if includeClaudeDetails {
		return "引数で渡されたタイトルから新規ノートを作成する。\n\n## 手順\n\n1. タイトルからファイル名を生成（kebab-case英語）\n2. 同名ファイルが存在しないか確認\n3. 以下のテンプレートでファイルを作成:\n\n" +
			fenced("markdown", `---
title: "$ARGUMENTS"
created: "<現在のRFC3339 UTC>"
modified: "<現在のRFC3339 UTC>"
tags: []
---

# $ARGUMENTS
`) + "\n\n4. " + bt("id") + " は省略する（Prometheusが自動付与）\n5. ファイルパスを報告"
	}

	return body + `
5. 作成した相対パスを返す`
}

func dailySkillBody(includeClaudeDetails bool) string {
	if includeClaudeDetails {
		return "今日の日付のデイリーノートを作成する（既に存在する場合は内容を表示）。\n\n## 手順\n\n1. " + bt("YYYY-MM-DD") + " 形式で今日の日付を取得\n2. " + bt("daily/YYYY-MM-DD.md") + " が存在するか確認\n3. 存在する場合: 内容を表示\n4. 存在しない場合: " + bt("daily/") + " ディレクトリを作成し、以下のテンプレートで作成:\n\n" +
			fenced("markdown", `---
title: "YYYY-MM-DD"
created: "<現在のRFC3339 UTC>"
modified: "<現在のRFC3339 UTC>"
tags: [daily]
---

# <日本語の日付（例: 2026年4月6日日曜日）>

## タスク

- [ ]

## メモ
`) + "\n\n5. " + bt("id") + " は省略する"
	}

	return `今日のデイリーノートを作成または開く。

## 手順

1. ` + bt("YYYY-MM-DD") + ` 形式で今日の日付を取得
2. ` + bt("daily/YYYY-MM-DD.md") + ` があるか確認
3. なければ作成
4. ` + bt("id") + ` は省略する`
}

func linkCheckSkillBody(includeClaudeDetails bool) string {
	if includeClaudeDetails {
		return "vault内の全ノートからwikilinkを抽出し、リンク先が存在するか確認する。\n\n## 手順\n\n1. " + bt("**/*.md") + " ファイルを全取得\n2. 各ファイルから " + bt("[[slug]]") + " と " + bt("[[slug|text]]") + " パターンを抽出\n3. slugのファイル解決（大文字小文字非区別、ハイフン↔スペース変換）:\n   - " + bt("[[setup]]") + " → setup.md を探す\n   - " + bt("[[Graph View]]") + " → graph-view.md を探す\n4. 解決できないリンクを一覧表示\n5. 各壊れたリンクに対して:\n   - 類似ファイル名の候補を提示\n   - スタブノートの作成を提案"
	}

	return `vault内のwikilinkを抽出し、壊れたリンクを検出する。

## 手順

1. ` + bt("**/*.md") + ` を列挙
2. ` + bt("[[slug]]") + ` と ` + bt("[[slug|text]]") + ` を抽出
3. slugをvault内のファイルに解決
4. 解決できないリンクと修復候補をまとめる`
}

func bulletList(items []string) string {
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = "- " + item
	}
	return strings.Join(lines, "\n")
}

func bt(s string) string {
	return "`" + s + "`"
}

func fenced(lang, body string) string {
	if lang == "" {
		return "```\n" + body + "\n```"
	}
	return "```" + lang + "\n" + body + "\n```"
}

const claudeSettings = `{
  "permissions": {
    "allow": [
      "Read",
      "Write",
      "Glob",
      "Grep",
      "Bash(ls *)",
      "Bash(cat *)",
      "Bash(head *)",
      "Bash(tail *)",
      "Bash(wc *)",
      "Bash(find *)",
      "Bash(date *)",
      "Bash(git add *)",
      "Bash(git commit *)",
      "Bash(git status)",
      "Bash(git log *)",
      "Bash(git diff *)"
    ],
    "deny": [
      "Bash(rm -rf *)",
      "Bash(rm -r *)",
      "Bash(mv .prometheus/data*)",
      "Bash(rm .prometheus/data*)",
      "Bash(curl *)",
      "Bash(wget *)",
      "Bash(git push *)",
      "Bash(git reset --hard *)"
    ]
  }
}
`

var frontmatterRule = buildFrontmatterRule()

const wikilinkRule = `---
paths:
  - "**/*.md"
---

# リンク規約

## 内部リンク（wikilink）— vault内のノートへのリンク
- ` + "`[[ファイル名]]`" + ` → そのファイルへのリンク
- ` + "`[[ファイル名|表示テキスト]]`" + ` → 別名リンク
- **vault内の.mdファイルへのリンクにのみ使う**

### Slug解決
- ファイル名（拡張子なし）で照合
- 大文字小文字非区別
- スペースとハイフンは相互変換（` + "`Graph View`" + ` = ` + "`graph-view`" + `）

### 例
- ` + "`[[setup]]`" + ` → setup.md
- ` + "`[[setup|セットアップガイド]]`" + ` → setup.md を「セットアップガイド」と表示
- ` + "`[[resources/design-md]]`" + ` → resources/design-md.md（サブディレクトリ内のノート）

## 外部リンク — URLへのリンク
- 通常のMarkdown記法を使う: ` + "`[表示テキスト](https://...)`" + `
- **` + "`[[https://...]]`" + ` は禁止**。wikilinkにURLを入れてはならない
- 外部リンクは新タブで開かれる

### 例
- ` + "`[GitHub](https://github.com/example/repo)`" + `
- ` + "`[公式ドキュメント](https://docs.example.com)`" + `
`

const noteExplorerAgent = `---
name: note-explorer
description: vault内のノートを探索・分析する
model: haiku
tools: Read Glob Grep
---

このvaultはPrometheusノートアプリのデータディレクトリ。
全ノートはYAMLフロントマター付きの.mdファイル。

探索時の注意:
- フロントマターのtagsフィールドでトピック分類を把握する
- [[wikilink]]でノート間の関連を辿る
- daily/YYYY-MM-DD.mdはデイリーノート
- .prometheus/は触らない
`

const vaultGitignore = `# Prometheus app data (regenerated from .md files)
.prometheus/data.db
.prometheus/data.db-journal
.prometheus/data.db-wal

# Claude Code personal config
CLAUDE.local.md
.claude/settings.local.json

# Generic agent personal config
AGENTS.local.md

# OS
.DS_Store

# Large binary assets (uncomment if you don't want to track images)
# assets/
`

func buildFrontmatterRule() string {
	return `---
paths:
  - "**/*.md"
---

# フロントマター規約

全ノートに以下のYAMLフロントマターが必要:

` + fenced("yaml", `---
id: "<ULID — 省略可。アプリが自動生成>"
title: "<ノートタイトル（日本語OK）>"
created: "<RFC3339 UTC — 例: 2026-04-06T09:00:00Z>"
modified: "<RFC3339 UTC>"
tags: [tag1, tag2]
theme: <オプション: prometheus, ocean, forest, sakura, nord 等>
---`) + `

- ` + "`id`" + ` はULID形式。手動作成時は省略してよい（Prometheusが自動付与）
- ` + "`id`" + ` を変更してはならない
- ` + "`modified`" + ` はノート更新時に現在時刻に更新する
- ` + "`tags`" + ` は空配列 ` + "`[]`" + ` でもよい
`
}
