# ADR-002: ハイブリッドストレージモデルの採用

## 背景・課題 (Background/Problem)
- Prometheusの根幹要件として「Markdownファイル（.md）でノートを管理する」がある
- 一方で、全文検索、双方向リンク解決、タグフィルタリング、タスク集約といった機能にはインデックスが必要
- ファイルシステムだけでは検索が遅く（毎回全ファイルをパースする必要がある）、DBだけではNeovim等の外部エディタとの併用ができない
- .mdファイルのポータビリティ（他のMarkdownツールへの移行容易性）を維持したい
- 外部エディタ（Neovim等）でファイルを直接編集した場合にも、Web UIに反映される必要がある

## 決定事項 (Decision)
- **.mdファイルを唯一の真実（Single Source of Truth）** とする
- **SQLiteはインデックス（使い捨て可能なキャッシュ）** として利用。DBを削除しても.mdファイルから再構築できる
- サーバー起動時にIndexerが全.mdファイルをパースし、メタデータをSQLiteに書き込む
- ファイル変更時（Web UI経由の保存）にIndexerが該当ファイルを再インデックス

### データフロー

```
Vault (.mdファイル) ← 唯一の真実 (Source of Truth)
      │
      ├── Go API (CRUD) ← Web UIからの編集
      │     │
      │     └── Indexer (goldmark解析) → SQLite
      │
      └── Neovim等 ← 外部エディタからの直接編集
            │
            └── サーバー再起動時にIndexerが再インデックス
```

### .mdファイルフォーマット

```yaml
---
id: "01HX4..."           # ULID (ソート可能・一意)
title: "ノートタイトル"
created: 2026-04-02T09:00:00Z
modified: 2026-04-02T14:30:00Z
tags: [tag1, tag2]
theme: ocean              # オプショナル: ノート個別テーマ
---
```

- フロントマターにメタデータを記述。YAML形式で`gopkg.in/yaml.v3`で解析
- IDにはULIDを採用。UUIDと異なりソート可能で、作成順の並び替えが自然にできる
- `[[wikilink]]` 構文でノート間リンクを表現。slugベースで解決（ファイル名の拡張子なし、大文字小文字非区別）
- `- [ ]` / `- [x]` 構文でタスクを表現。Indexerがline_numberと共にtasksテーブルへ抽出

### SQLiteスキーマ

| テーブル | 用途 |
|---|---|
| `notes` | .mdファイルのインデックス（id, path, title, content, checksum） |
| `links` | 双方向リンク（source_id, target_slug） |
| `tags` / `note_tags` | タグの多対多リレーション |
| `tasks` | タスク（content, completed, line_number） |
| `notes_fts` | FTS5仮想テーブル（全文検索用） |

### Indexerの処理フロー

1. ファイルを読み込み
2. SHA-256チェックサムを計算
3. DBの既存チェックサムと比較 → 同じならスキップ
4. goldmarkでMarkdownをパース
5. wikilink抽出（正規表現）
6. タスク抽出（`- [ ]` / `- [x]`）
7. フロントマターからタグ抽出
8. SQLiteにupsert（notes, links, tags, note_tags, tasks）
9. 起動時にvaultに存在しないDBレコードを削除（stale entry cleanup）

### 理由 (Reasons)
- **.mdファイルの普遍性**: Obsidian、Typora、VS Code、Neovim — あらゆるツールで読み書きできる。ベンダーロックインがゼロ
- **外部エディタとの互換性**: Neovimで編集 → サーバー再起動で反映
- **DBの再構築可能性**: DBが破損してもゼロコストで復旧。サーバー再起動で全ファイルから再インデックス
- **Git管理の容易さ**: .mdファイルをGitで管理すれば、バージョン管理・差分・マージが標準ツールで可能

### 受け入れるトレードオフ (Accepted Trade-offs)
- **外部エディタの変更はリアルタイム反映されない**: 現在はサーバー再起動時にインデックスが更新される。ファイル監視（fsnotify）は将来の改善候補
- **フロントマターの記述コスト**: 新規ノート作成時にフロントマター（id, title, created等）を自動生成する仕組みが必須

## 検討した別の選択肢 (Alternatives Considered)

### DB Only（.mdファイルなし）
- **メリット**: 二重管理不要、検索・クエリが常に高速
- **デメリット**: 外部エディタとの併用不可、ベンダーロックイン、Git管理不可
- **不採用理由**: Prometheusの根幹要件「.mdファイルで管理」に反する

### ファイルシステム Only（DBなし）
- **メリット**: シンプル、二重管理不要
- **デメリット**: 全文検索に毎回全ファイルパースが必要（遅い）
- **不採用理由**: ノート数が増えるにつれてパフォーマンスが劣化

## 参考 (References)
- [SQLite FTS5](https://www.sqlite.org/fts5.html)
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
- [goldmark](https://github.com/yuin/goldmark)
- [ULID Spec](https://github.com/ulid/spec)

## 議論 (Discussion)
- チェックサムベースの変更検知を採用。SHA-256で実際に内容が変わったかを判定し、不要なインデックス更新を回避
- 当初はchokidar（Node.js）でファイル監視を行う設計だったが、Go移行後はfsnotifyに置き換え予定。現時点ではサーバー起動時の全インデックスで対応
- Markdown処理はgoldmark + カスタムwikilink変換。Node.js時代はunified/remark/rehypeを使用していたが、Go移行でgoldmarkに統一
- FTS5（unicode61トークナイザー）で全文検索を実装済み。日本語にも対応。FTS5が利用できない場合のみLIKEにフォールバック
