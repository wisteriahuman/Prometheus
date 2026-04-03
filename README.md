# Prometheus

Markdownベースの「第二の脳」ノートアプリ。

## 特徴

- **Markdownファイル管理** — ノートは `.md` ファイルとして保存。Neovimや好きなエディタで直接編集可能
- **双方向リンク** — `[[wikilink]]` でノート同士を自然に接続
- **グラフビュー** — ノートの関係をタグクラスタ付きで可視化
- **全文検索** — 高速なインクリメンタル検索
- **デイリーノート** — テンプレート付きの日記機能
- **タスク管理** — `- [ ]` 構文でタスクを自動抽出
- **6つのテーマ** — Prometheus / Ocean / Forest / Sakura / Nord / Solarized
- **Vimキーバインド** — エディタ内でVimモーション
- **Mermaid図** — コードブロックで図を描画

## セットアップ

```bash
git clone https://github.com/your/prometheus.git
cd prometheus
pnpm install
```

### ノートの保存先を設定

```bash
cd app
cp .env.example .env
```

`.env` を編集して `PROMETHEUS_VAULT_PATH` を設定:

```env
# 例: ホームディレクトリに保存
PROMETHEUS_VAULT_PATH=~/my-notes

# 例: デフォルト（リポジトリ内の vault/）
PROMETHEUS_VAULT_PATH=./vault
```

### 起動

```bash
# データベース初期化（初回のみ）
pnpm --filter app run db:migrate

# 開発サーバー起動
pnpm --filter app run dev
```

http://localhost:5173 でアクセス。初回起動時にサンプルノートが自動生成されます。

## 本番デプロイ（Mac mini等）

```bash
# ビルド
pnpm --filter app run build

# 起動
cd app
PORT=3000 PROMETHEUS_VAULT_PATH=/path/to/vault node build/index.js
```

### Docker

```bash
docker compose up -d
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
| `⌘K` → `?` | ショートカット一覧（コマンドパレットから） |

## Neovimとの併用

`PROMETHEUS_VAULT_PATH` で指定したディレクトリ内の `.md` ファイルを直接編集できます。
Prometheusはファイル変更を自動検知してインデックスを更新します。

```bash
# 例
nvim ~/my-notes/new-idea.md
```

## 技術スタック

- SvelteKit + Tailwind CSS v4
- CodeMirror 6 (Vimモード)
- SQLite (FTS5全文検索)
- D3.js (グラフビュー)
- Mermaid (図表描画)
