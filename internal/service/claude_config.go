package service

const claudeMd = `# Prometheus Vault

このディレクトリ（CWD）がvaultのルート。全ノートのパスはここからの相対パス。
Go製のWebサーバー ` + "`prm`" + ` がこのディレクトリを監視し、Web UIを提供する。

## ディレクトリ構造

- ` + "`*.md`" + ` — ノート本体。YAMLフロントマター必須
- ` + "`サブディレクトリ/*.md`" + ` — カテゴリ別ノート。自由にディレクトリ作成可
- ` + "`daily/YYYY-MM-DD.md`" + ` — デイリーノート
- ` + "`assets/`" + ` — 画像・PDFなどのアップロードファイル
- ` + "`.prometheus/config.json`" + ` — vault設定。読み書き可能
- ` + "`.prometheus/data.db`" + ` — SQLite FTS5インデックス。**変更禁止**

ノートのパスはvaultルートからの相対パス（例: ` + "`resources/design-md.md`" + `）。
Web UIのURLも同じパスになる: ` + "`/note/resources/design-md.md`" + `

## フロントマター

詳細は ` + "`.claude/rules/frontmatter.md`" + ` を参照。例:

` + "```yaml" + `
---
id: "01KNE4B96D..."
title: "ノートタイトル"
created: "2026-04-06T09:00:00Z"
modified: "2026-04-06T09:00:00Z"
tags: [tag1, tag2]
---
` + "```" + `

- ` + "`id`" + `: ULID形式。**変更禁止**。省略すればアプリが自動生成
- ` + "`created`/`modified`" + `: RFC3339 UTC
- ` + "`tags`" + `: 配列。グラフビューでクラスタ化される
- ` + "`theme`" + `: ノート個別テーマ（省略可）

## Wikilink

` + "`[[ファイル名]]`" + ` または ` + "`[[ファイル名|表示テキスト]]`" + `
slugは小文字化、スペース→ハイフン変換。` + "`[[Setup]]`" + ` → ` + "`setup.md`" + `
詳細は ` + "`.claude/rules/wikilink.md`" + ` を参照。

## 対応記法

- **Mermaid**: ` + "````mermaid`" + ` コードブロックで図を描画（フローチャート、シーケンス図、ER図等）
- **LaTeX数式**: ` + "`$...$`" + ` でインライン数式、` + "`$$...$$`" + ` でブロック数式（KaTeX）
- **GFM**: テーブル、タスクリスト（` + "`- [ ]`" + `）、取り消し線、脚注

## 禁止事項

- ` + "`.prometheus/data.db`" + ` を変更・削除しない
- フロントマターの ` + "`id`" + ` を変更しない
- ` + "`assets/`" + ` のファイルを削除しない

## 操作のコツ

- ノート作成時は ` + "`id`" + ` を省略（アプリが自動付与）
- ファイル名はkebab-case英語、タイトルは日本語推奨
- ` + "`.prometheus/config.json`" + ` のテーマやテンプレートは変更可能
- 個人設定は ` + "`CLAUDE.local.md`" + ` を作成（gitignore済み）
`

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

const ruleFrontmatter = `---
paths:
  - "**/*.md"
---

# フロントマター規約

全ノートに以下のYAMLフロントマターが必要:

` + "```yaml" + `
---
id: "<ULID — 省略可。アプリが自動生成>"
title: "<ノートタイトル（日本語OK）>"
created: "<RFC3339 UTC — 例: 2026-04-06T09:00:00Z>"
modified: "<RFC3339 UTC>"
tags: [tag1, tag2]
theme: <オプション: prometheus, ocean, forest, sakura, nord 等>
---
` + "```" + `

- ` + "`id`" + ` はULID形式。手動作成時は省略してよい（Prometheusが自動付与）
- ` + "`id`" + ` を変更してはならない
- ` + "`modified`" + ` はノート更新時に現在時刻に更新する
- ` + "`tags`" + ` は空配列 ` + "`[]`" + ` でもよい
`

const ruleWikilink = `---
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

const skillCreateNote = `---
name: create-note
description: フロントマター付きの新規ノートを作成
argument-hint: "[title]"
allowed-tools: Read Write Glob
---

引数で渡されたタイトルから新規ノートを作成する。

## 手順

1. タイトルからファイル名を生成（kebab-case英語）
2. 同名ファイルが存在しないか確認
3. 以下のテンプレートでファイルを作成:

` + "```markdown" + `
---
title: "$ARGUMENTS"
created: "<現在のRFC3339 UTC>"
modified: "<現在のRFC3339 UTC>"
tags: []
---

# $ARGUMENTS

` + "```" + `

4. ` + "`id`" + ` は省略する（Prometheusが自動付与）
5. ファイルパスを報告
`

const skillDaily = `---
name: daily
description: 今日のデイリーノートを作成または開く
allowed-tools: Read Write Glob
---

今日の日付のデイリーノートを作成する（既に存在する場合は内容を表示）。

## 手順

1. 今日の日付を ` + "`YYYY-MM-DD`" + ` 形式で取得
2. ` + "`daily/YYYY-MM-DD.md`" + ` が存在するか確認
3. 存在する場合: 内容を表示
4. 存在しない場合: ` + "`daily/`" + ` ディレクトリを作成し、以下のテンプレートで作成:

` + "```markdown" + `
---
title: "YYYY-MM-DD"
created: "<現在のRFC3339 UTC>"
modified: "<現在のRFC3339 UTC>"
tags: [daily]
---

# <日本語の日付（例: 2026年4月6日日曜日）>

## タスク

- [ ]

## メモ

` + "```" + `

5. ` + "`id`" + ` は省略する
`

const skillLinkCheck = `---
name: link-check
description: 壊れたwikilinkを検出し修復提案する
allowed-tools: Read Glob Grep
---

vault内の全ノートからwikilinkを抽出し、リンク先が存在するか確認する。

## 手順

1. ` + "`**/*.md`" + ` ファイルを全取得
2. 各ファイルから ` + "`[[slug]]`" + ` と ` + "`[[slug|text]]`" + ` パターンを抽出
3. slugのファイル解決（大文字小文字非区別、ハイフン↔スペース変換）:
   - ` + "`[[setup]]`" + ` → setup.md を探す
   - ` + "`[[Graph View]]`" + ` → graph-view.md を探す
4. 解決できないリンクを一覧表示
5. 各壊れたリンクに対して:
   - 類似ファイル名の候補を提示
   - スタブノートの作成を提案
`

const agentNoteExplorer = `---
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

# OS
.DS_Store

# Large binary assets (uncomment if you don't want to track images)
# assets/
`
