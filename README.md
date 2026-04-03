# Prometheus

Markdownベースの「第二の脳」ノートアプリ。

## 特徴

- **Markdownファイル管理** — ノートは `.md` ファイルとして保存。Neovimや好きなエディタで直接編集可能
- **双方向リンク** — `[[wikilink]]` でノート同士を自然に接続
- **グラフビュー** — ノートの関係をタグクラスタ付きで可視化
- **全文検索** — FTS5による高速インクリメンタル検索
- **デイリーノート** — テンプレート付きの日記機能
- **タスク管理** — `- [ ]` 構文でタスクを自動抽出
- **13テーマ + カスタムテーマ** — ダーク/ライト/ビジネス + 自作可能
- **Vimキーバインド** — エディタ内でVimモーション、`:w` で保存
- **Mermaid図** — コードブロックで図を描画
- **シングルバイナリ** — Go製。`brew install` で即利用可能

## インストール

```bash
# Homebrew (macOS / Linux)
brew install wisteriahuman/tap/prm

# Go
go install github.com/wisteriahuman/prometheus/cmd/prometheus@latest

# GitHub Releases からバイナリを直接ダウンロード
# https://github.com/wisteriahuman/Prometheus/releases
```

## 使い方

```bash
# ノートを ~/my-notes に保存して起動
prm dev ~/my-notes

# 別のvaultを別ポートで起動
prm dev ~/work -p 3001

# 新しいvaultを作成
prm init ~/new-vault

# 設定を確認
prm info
```

初回起動時にサンプルノートが自動生成されます。

## 開発

```bash
git clone https://github.com/wisteriahuman/Prometheus.git
cd Prometheus

# フロントエンドビルド + Goバイナリビルド
make build

# 起動
./prm dev ~/my-notes
```

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
| `:w` / `⌘S` | 保存 |
| `⌘↵` | エディタ/プレビュー切替 |

## Neovimとの併用

vaultディレクトリ内の `.md` ファイルを直接編集できます。
変更は自動検知されインデックスが更新されます。

## 技術スタック

- **バックエンド**: Go (chi, modernc.org/sqlite, goldmark)
- **フロントエンド**: SvelteKit (SPA) + Tailwind CSS v4
- **エディタ**: CodeMirror 6 (Vimモード)
- **グラフ**: D3.js
- **図表**: Mermaid
