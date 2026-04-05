# Prometheus

軽量なMarkdownノートアプリ。Goシングルバイナリ ~30MB。

> `brew install` して `prm dev ~/notes` で起動。ノートは `.md` ファイル。Neovimでそのまま編集できる。

## なぜ Prometheus?

- **Obsidianは重い。** 300MB / 500MBメモリ。Prometheusは30MB / 23MBメモリ
- **Notionはオフラインで使えない。** Prometheusは完全ローカル。クラウド不要
- **ノートは .md ファイル。** `cat`, `grep`, `vim`, AI — UNIXツールチェインがそのまま使える
- **SQLiteでインデックス。** FTS5全文検索がサブミリ秒。1000ファイル超えても速い
- **1コマンドで起動。** `prm dev ~/notes` で終わり。Dockerもアカウント作成も不要
- **複数vault同時起動。** 1インスタンス23MB。3つ開いてもObsidian1つより軽い

## インストール

```bash
# Homebrew (macOS / Linux)
brew install wisteriahuman/tap/prm

# Go
go install github.com/wisteriahuman/prometheus/cmd/prometheus@latest

# GitHub Releases
# https://github.com/wisteriahuman/Prometheus/releases
```

## 使い方

```bash
prm dev ~/my-notes           # 起動
prm dev ~/work -p 3001       # 別vaultを別ポートで
prm init ~/new-vault         # vault作成
prm info                     # 設定確認
```

## 機能

- **Vimキーバインド** — `:w` で保存。エディタ内でVimモーション
- **[[wikilink]]** — 双方向リンク + タグクラスタ付きグラフビュー
- **全文検索** — SQLite FTS5。⌘P でファイル切替、⌘⇧F で検索
- **デイリーノート** — ⌘D で今日のノートを開く。テンプレートカスタマイズ可能
- **タスク管理** — `- [ ]` をチェック切替。フィルタリング付き一覧
- **数式** — `$E=mc^2$` がKaTeXでレンダリング
- **Mermaid** — コードブロックで図を描画
- **13テーマ + カスタム** — ダーク/ライト/ビジネス。ノート個別テーマも
- **エクスポート** — HTML（テーマ付き）/ Markdown / PDF
- **ファイルアップロード** — D&Dやペーストで画像挿入
- **コマンドパレット** — ⌘K で全機能にアクセス

## Neovimとの併用

vaultは普通のディレクトリ。中の `.md` ファイルをNeovimで直接編集できる。

```bash
nvim ~/my-notes/idea.md    # Neovimで編集
prm dev ~/my-notes          # 同じvaultをWeb UIで表示
```

AIとも相性がいい。`.md` ファイルだから `cat` でそのままLLMに渡せる。

## ショートカット

| キー | 機能 |
|------|------|
| `⌘K` | コマンドパレット |
| `⌘P` | ファイル切替 |
| `⌘⇧F` | 全文検索 |
| `⌘D` | デイリーノート |
| `⌘G` | グラフビュー |
| `⌘.` | テーマ切替 |
| `⌘\` | サイドバー切替 |
| `:w` | 保存 |
| `⌘↵` | エディタ/プレビュー切替 |

## 技術スタック

| | |
|---|---|
| バックエンド | Go (chi, modernc.org/sqlite, goldmark) |
| フロントエンド | SvelteKit (SPA) + Tailwind CSS v4 |
| エディタ | CodeMirror 6 + Vim |
| 検索 | SQLite FTS5 (unicode61) |
| グラフ | D3.js |
| 数式 | KaTeX |
| 図表 | Mermaid |
| 配信 | GoReleaser + Homebrew |

## 開発

```bash
git clone https://github.com/wisteriahuman/Prometheus.git
cd Prometheus
make build
./prm dev ~/my-notes
```
