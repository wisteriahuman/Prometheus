# ADR-003: アプリケーションアーキテクチャ

## 背景・課題 (Background/Problem)
- PrometheusはGoバックエンド + SvelteフロントエンドのWebアプリケーション
- シングルバイナリで配信するため、フロントエンドをGoバイナリに埋め込む必要がある
- 個人利用に特化し、認証やマルチテナントは不要
- アーキテクチャは過剰に複雑にせず、ノートアプリの規模に適したシンプルな構成にする

## 決定事項 (Decision)

### 全体構成

```
┌─────────────────────────────────────────┐
│           Go シングルバイナリ              │
│                                         │
│  ┌──────────┐   ┌────────────────────┐  │
│  │ cobra CLI │   │  embed (SPA)       │  │
│  └──────────┘   │  ┌──────────────┐  │  │
│                  │  │ Svelte SPA   │  │  │
│  ┌──────────┐   │  │ (静的HTML/JS) │  │  │
│  │ chi HTTP │◄──│  └──────────────┘  │  │
│  │ Router   │   └────────────────────┘  │
│  └────┬─────┘                           │
│       │                                 │
│  ┌────┴─────┐                           │
│  │ Handler  │  ← 18 APIエンドポイント     │
│  └────┬─────┘                           │
│       │                                 │
│  ┌────┴──────┐                          │
│  │ Service   │  ← Vault, Markdown,      │
│  │           │    Indexer, Daily         │
│  └────┬──────┘                          │
│       │                                 │
│  ┌────┴───┐  ┌──────────┐              │
│  │ SQLite │  │ Vault    │              │
│  │ (DB)   │  │ (.mdファイル)│              │
│  └────────┘  └──────────┘              │
└─────────────────────────────────────────┘
```

### レイヤー構成

2層+インフラの軽量アーキテクチャを採用。

| レイヤー | ディレクトリ | 責務 |
|---|---|---|
| **Handler** | `internal/handler/` | HTTPリクエスト処理、JSON入出力、バリデーション |
| **Service** | `internal/service/` | ビジネスロジック（Vault操作、Markdown変換、インデックス） |
| **DB** | `internal/db/` | SQLite接続、マイグレーション |
| **Config** | `internal/config/` | 環境変数/CLIフラグ読み込み |
| **Server** | `internal/server/` | DI、ルーティング、SPA配信 |

### ディレクトリ構成

```
prometheus/
├── cmd/prometheus/main.go       # エントリポイント（cobra CLI）
├── internal/
│   ├── server/
│   │   ├── server.go            # DI + HTTPサーバー起動
│   │   ├── router.go            # chiルーティング（18 API）
│   │   ├── spa.go               # embed SPA配信
│   │   └── web/                 # ビルド済みSPA（go:embed対象）
│   ├── handler/                 # 12ファイル（notes, search, graph等）
│   ├── service/
│   │   ├── vault.go             # .mdファイルCRUD、フロントマター
│   │   ├── markdown.go          # goldmark MD→HTML、wikilink変換
│   │   ├── indexer.go           # ノート→SQLite同期
│   │   ├── daily.go             # デイリーノート生成
│   │   └── init.go              # サンプルノート生成
│   ├── db/
│   │   ├── db.go                # SQLite接続
│   │   └── migration.go         # CREATE TABLE
│   └── config/config.go         # 設定
├── web/                          # Svelteフロントエンド（ソース）
│   ├── src/
│   │   ├── lib/components/      # UI コンポーネント
│   │   ├── lib/stores/          # テーマ、サイドバー
│   │   └── routes/              # SPA ルート
│   ├── svelte.config.js         # adapter-static
│   └── package.json
├── docs/adr/                     # ADR
├── go.mod
├── Makefile
└── .goreleaser.yaml              # リリース設定
```

### DI（依存性注入）

`internal/server/server.go`で手動DI:

```go
func NewServer(cfg *config.Config) *Server {
    database := db.New(cfg.DBPath)     // SQLite接続
    vault := service.NewVault(...)     // Vault操作
    md := service.NewMarkdown()        // Markdown変換
    indexer := service.NewIndexer(database, vault, md)  // インデックス
    daily := service.NewDaily(vault)   // デイリーノート

    handlers := &Handlers{
        Notes:   handler.NewNotesHandler(vault, indexer),
        Search:  handler.NewSearchHandler(database),
        // ... 全12ハンドラ
    }

    router := NewRouter(handlers)      // chiルーティング
}
```

go-college-09-fujiiと同じ手動DI（コンストラクタ注入）パターン。DIコンテナは不使用。

### 理由 (Reasons)
- **クリーンアーキテクチャは不採用**: ノートアプリのドメインロジックは薄い（ファイルCRUD + インデックス）。usecase/domain/infra の4層分離は過剰
- **Handler → Service の2層**: Handlerがリクエストを受け取り、Serviceに委譲。シンプルで見通しが良い
- **SQLは直接書く**: ORMを使わず`database/sql`で直接SQLを書く。クエリが単純なのでORMのオーバーヘッドが不要

### 受け入れるトレードオフ (Accepted Trade-offs)
- **Serviceが肥大化する可能性**: vault.goが大きくなったらファイル分割で対応
- **テストのDB依存**: Service層がSQLiteに直接依存。モック化はしていないが、個人プロジェクトでは実DB使用のテストで十分

## 検討した別の選択肢 (Alternatives Considered)

### クリーンアーキテクチャ（go-college-09-fujii方式）
- **メリット**: レイヤー分離が明確、テスタビリティ高い
- **デメリット**: ノートアプリには過剰。handler → usecase → domain → infra の4層でボイラープレートが増える
- **不採用理由**: ドメインロジックが薄い（CRUDが中心）ため、2層で十分

### Echo HTTPフレームワーク
- **メリット**: go-college-09-fujiiと同じ、バリデーション組み込み
- **デメリット**: chiで十分。Echoの機能（バインディング、バリデーション等）を使う場面がない
- **不採用理由**: net/http + chi の方が軽量で標準ライブラリに近い

## 参考 (References)
- [chi router](https://github.com/go-chi/chi)
- [Go embed](https://pkg.go.dev/embed)
- [cobra CLI](https://github.com/spf13/cobra)

## 議論 (Discussion)
- フロントエンドとバックエンドの接続はREST API（19エンドポイント）。GraphQLやgRPCは不要（クライアントが1つしかない）
- SPA配信は`//go:embed`で静的ファイルをバイナリに埋め込み、SPAフォールバック（存在しないパスはindex.htmlを返す）で実現
- 起動時の初期化フロー: vault作成(初回) → DBマイグレーション → サンプルノート生成(初回) → 全ノートインデックス → HTTPサーバー起動
- 全文検索はSQLite FTS5（unicode61トークナイザー）を使用。日本語対応。検索レスポンスはサブミリ秒（163μs）。FTS5が利用できない場合はLIKEにフォールバック
- テーマ設定（アプリ全体テーマ + カスタムテーマ）はvaultの`.prometheus/config.json`に保存。localStorageは使用しない。理由: 別ブラウザ・別端末からアクセスしても同じテーマが維持される。vaultが全設定の中心
- グラフ描画はSVGを採用。Canvasは描画速度で優位だがテキストがぼやけ、ヒットテストが複雑。個人ノートアプリで現実的なノート数（100-500）ではSVGで十分な性能
