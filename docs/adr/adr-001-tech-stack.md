# ADR-001: 技術スタックの選定

## 背景・課題 (Background/Problem)
- Prometheusは「第二の脳」を目指すMarkdownベースのノートアプリである
- シングルバイナリで配信し、`brew install`や`go install`で即利用可能にする
- パフォーマンス（体感速度）とUIの完成度を最優先とする
- .mdファイルを実データとして保持しつつ、リッチなWeb UIを提供する
- 外部エディタ（Neovim等）との併用を想定し、キーボードファーストのUXを実現する
- 個人利用に全振りし、チーム機能は持たない

## 決定事項 (Decision)

### バックエンド

| レイヤー | 技術 | 理由 |
|---|---|---|
| 言語 | **Go** | シングルバイナリ、クロスコンパイル、高速 |
| HTTP | **net/http + chi** | 標準ライブラリベース、軽量ルーター |
| SQLite | **modernc.org/sqlite** | pure Go（cgo不要）、クロスコンパイル可能 |
| Markdown | **goldmark** | 高速、GFM対応、拡張可能 |
| YAML | **gopkg.in/yaml.v3** | フロントマター解析 |
| CLI | **cobra** | サブコマンド管理 |
| 静的ファイル | **embed (Go標準)** | SPAをバイナリに同梱 |

### フロントエンド

| レイヤー | 技術 | 理由 |
|---|---|---|
| Framework | **SvelteKit (SPA)** | バンドル最小、adapter-staticでSPA出力 |
| CSS | **Tailwind CSS v4** | CSS変数ベースでランタイムテーマ切替 |
| エディタ | **CodeMirror 6** | 軽量、Vimモード、モバイル対応 |
| グラフ | **D3.js (d3-force)** | フル制御、タグクラスタ可視化 |
| アイコン | **Lucide** | 軽量SVGアイコン |
| 図表 | **Mermaid** | コードブロックで図を描画 |

### 理由 (Reasons)

#### Go バックエンド

- **シングルバイナリ配信**: `go build`で1ファイルに全て含まれる。Node.jsのようにランタイムやnode_modulesが不要
- **`//go:embed`でSPA同梱**: ビルド済みのSvelteアプリをGoバイナリに埋め込み、1ファイルで完結
- **クロスコンパイル**: `GOOS=darwin/linux/windows`で3プラットフォーム対応。CGO_ENABLED=0でpure Goビルド
- **brew/go install配信**: GoReleaserでGitHub Releases + homebrew-tapに自動配信
- **modernc.org/sqlite**: pure GoのSQLite実装。cgo不要なのでクロスコンパイルが問題なく動く

#### SvelteKit SPA

- **SSRなし**: `adapter-static`でクライアントサイドのみ。GoがAPIサーバーと静的ファイル配信を担当
- **バンドル最小**: 仮想DOMなし、コンパイル時にVanilla JSに変換
- **全データ取得はfetch**: `+page.server.ts`なし。全てクライアントサイドでAPIを呼ぶ

### 受け入れるトレードオフ (Accepted Trade-offs)
- **SSRなし**: 初回表示が数十ms遅くなるが、個人ノートアプリではSEO不要、体感差なし
- **SvelteのエコシステムがReactより小さい**: カスタムコンポーネント中心なので影響なし
- **Goのテンプレートエンジン不使用**: フロントエンドは完全にSvelte。Goはpure APIサーバー

## 検討した別の選択肢 (Alternatives Considered)

### Node.js + SvelteKit SSR（以前の構成）
- **メリット**: SSRで初回表示が速い、サーバーとフロントが一体
- **デメリット**: シングルバイナリ配信不可。Node.js + node_modules + better-sqlite3のネイティブビルドが必要
- **不採用理由**: 配信の容易さがGoに大きく劣る。`npm install`でネイティブビルドが失敗するリスク

### Rust (Axum/Actix)
- **メリット**: 最高パフォーマンス
- **デメリット**: ビルド時間が長い、開発速度が遅い
- **不採用理由**: ノートアプリの規模ではGoで十分な性能。開発速度を優先

### Electron / Tauri
- **メリット**: ネイティブデスクトップアプリ
- **デメリット**: 配信サイズが大きい（Electron: 100MB+）、Web UIと二重管理
- **不採用理由**: ブラウザで十分。Goシングルバイナリの方が軽い

## 参考 (References)
- [Go embed](https://pkg.go.dev/embed)
- [chi router](https://github.com/go-chi/chi)
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
- [goldmark](https://github.com/yuin/goldmark)
- [GoReleaser](https://goreleaser.com/)
- [SvelteKit adapter-static](https://svelte.dev/docs/kit/adapter-static)
- [cobra CLI](https://github.com/spf13/cobra)

## 議論 (Discussion)
- 当初はNode.js + SvelteKit SSRで構築。better-sqlite3のネイティブビルド問題と配信の困難さからGoに移行
- Go移行時、フロントエンドのSvelteコンポーネントは一切変更せず、サーバーサイドのみ書き換え。APIインターフェース（18エンドポイント）を維持したため、フロントの修正はSSR→SPA化（`+page.server.ts`削除、クライアントサイドfetch化）のみ
- SQLiteドライバは`modernc.org/sqlite`（pure Go）を採用。`mattn/go-sqlite3`はcgo必須でクロスコンパイルが複雑になるため不採用
- CLIはcobraを採用。`prometheus dev ~/my-notes`の1コマンドでvault作成→DB初期化→サーバー起動→ブラウザアクセスまで完結
