# ADR-001: 技術スタックの選定

## 背景・課題 (Background/Problem)
- Prometheusは「第二の脳」を目指すMarkdownベースのノートアプリである
- Webアプリとして構築し、デスクトップ・モバイル両方で快適に動作する必要がある
- パフォーマンス（体感速度）とUIの完成度を最優先とする
- .mdファイルを実データとして保持しつつ、リッチなWeb UIを提供するというハイブリッドな要件がある
- 外部エディタ（Neovim等）との併用を想定し、キーボードファーストのUXを実現する
- 将来的にマルチユーザー対応を見据えるが、個人利用の軽さを犠牲にしない

## 決定事項 (Decision)

| レイヤー | 技術 | バージョン |
|---|---|---|
| Frontend Framework | SvelteKit | 2.x (Svelte 5) |
| Runtime | **Node.js + pnpm** | Node.js 22+, pnpm 10+ |
| CSS Framework | Tailwind CSS | v4 |
| Markdown Editor | CodeMirror 6 | - |
| Markdown Parser | unified + remark + rehype | - |
| Graph Visualization | D3.js (d3-force) | - |
| ORM | Drizzle ORM | - |
| SQLite Driver | better-sqlite3 | - |
| Validation | Zod | - |
| 環境変数 | dotenv | - |

### 理由 (Reasons)

#### SvelteKit

- **バンドルサイズ最小**: Svelteはコンパイル時にVanilla JSへ変換される。Reactのような仮想DOMランタイム（~80KB gzip）が不要
- **カスタムアプリに最適**: Prometheusのコアコンポーネント（エディタ、グラフ、コマンドパレット）は全てカスタム実装。既存UIコンポーネントライブラリへの依存が少ないため、Reactエコシステムの優位性が活きない
- **トランジション・アニメーション組み込み**: ノート切替、サイドバー開閉、テーマ切替のアニメーションがフレームワークレベルでサポートされる
- **リアクティビティがシンプル**: Svelte 5のRunes（`$state`, `$derived`, `$effect`）はReactのuseState/useEffect/useMemoより直感的で、バグが入り込む余地が少ない
- **コード量削減**: 同等の機能をReactより少ないコードで実装でき、メンテナンス性が高い

#### Node.js + pnpm

- **Vite SSR互換**: SvelteKitの開発サーバー（Vite）はNode.js上で動作する。SSRコンテキストではNode.jsのAPIが使用されるため、Node.js互換のライブラリが必要
- **better-sqlite3**: Node.js向けの高速SQLiteドライバ。bun:sqliteと同等のパフォーマンスでネイティブバインディングを提供
- **エコシステムの安定性**: 全npmパッケージとの完全な互換性。ネイティブモジュール（@node-rs/argon2、better-sqlite3等）が問題なく動作する
- **pnpmのワークスペース管理**: monorepo構成で`pnpm --filter`による効率的なパッケージ管理が可能
- **pnpm onlyBuiltDependencies**: ネイティブモジュールのビルドをpnpmが適切に管理

#### Tailwind CSS v4

- **CSS変数ベースの設計**: `@theme`ディレクティブでCSS変数を定義し、ランタイムでの値変更が可能。ノート/ワークスペースごとのテーマ切替に最適
- **JIT**: 未使用クラスを含まない最小CSSを生成

#### CodeMirror 6

- **Vim keybinding対応**: `@codemirror/vim`で本格的なVimモーションを提供。Neovimユーザーにとって自然な編集体験
- **軽量**: Monaco Editor（~2MB）に対して~150KB（必要な拡張のみ）
- **モバイル対応**: IME（日本語入力）との互換性が高い
- **増分パーシング**: 大きなMarkdownファイルでも入力レスポンスが劣化しない
- **カスタム拡張容易**: wikilink構文、タスクチェックボックス等のカスタム構文ハイライトを追加可能

#### D3.js (d3-force)

- **最小依存**: vis-networkやCytoscape.jsのような大きなライブラリに依存せず、d3-forceモジュールのみを利用
- **フル制御**: ノードの描画、インタラクション、スタイリングを完全にカスタマイズ可能
- **SVG/Canvas切替**: ノード数に応じてSVG（少数時、高品質）とCanvas（大量時、高性能）を使い分け可能

### 受け入れるトレードオフ (Accepted Trade-offs)
- **Svelteのエコシステムの小ささ**: React向けのUIコンポーネントライブラリ（shadcn/ui等）が使えない。ただしPrometheusはほぼ全コンポーネントをカスタム実装するため影響は限定的
- **better-sqlite3のネイティブコンパイル**: bun:sqliteと異なりネイティブバインディングのコンパイルが必要。ただしpnpmのonlyBuiltDependenciesで管理され、実運用上の問題はない
- **CodeMirror VimモードとNative Neovimの差**: 完全なNeovim互換ではない。しかし外部エディタとの併用（chokidarによるファイル監視）でカバーする

## 検討した別の選択肢 (Alternatives Considered)

### Next.js 15 (App Router)
- **メリット**: Reactエコシステム最大、Server Components、情報・チュートリアル豊富
- **デメリット**: Reactランタイム（~80KB）のオーバーヘッド、App Routerの複雑性、Prometheusのカスタムコンポーネント中心の設計ではエコシステムの優位が活きない
- **不採用理由**: パフォーマンスとコードのシンプルさでSvelteKitに劣る

### SolidStart
- **メリット**: JSXベースで最高パフォーマンス、細粒度リアクティビティ
- **デメリット**: メタフレームワーク（SolidStart）がまだ成熟途上、エコシステムが最も小さい
- **不採用理由**: SvelteKit同等のパフォーマンスだが、フレームワークの安定性でリスクが高い

### Monaco Editor
- **メリット**: VS Code同等の編集機能
- **デメリット**: バンドルサイズ~2MB、モバイル対応が弱い、Web Worker必須
- **不採用理由**: ノートアプリにはオーバースペック。CodeMirror 6の方が軽量で要件を満たす

### Bun (Runtime)
- **メリット**: SQLite内蔵（bun:sqlite）でネイティブバインディング不要、ファイルI/O高速、起動速度が高速
- **デメリット**: SvelteKitの開発サーバー（Vite）はNode.js上でSSRを実行するため、bun:sqliteが使用不可。Vite SSRコンテキストではBun固有のAPIにアクセスできない
- **不採用理由**: ViteのSSRがNode.jsコンテキストで動作するため、bun:sqliteを利用できない。結局better-sqlite3等のNode.js互換ドライバが必要となり、Bunの最大の優位性（SQLite内蔵）が活きない

## 参考 (References)
- [SvelteKit Documentation](https://svelte.dev/docs/kit)
- [Node.js Documentation](https://nodejs.org/docs/latest-v22.x/api/)
- [pnpm Documentation](https://pnpm.io/)
- [better-sqlite3](https://github.com/WiseLibs/better-sqlite3)
- [Tailwind CSS v4](https://tailwindcss.com/docs)
- [CodeMirror 6](https://codemirror.net/)
- [D3 Force-Directed Graph](https://d3js.org/d3-force)
- [Drizzle ORM](https://orm.drizzle.team/)

## 議論 (Discussion)
- フロントエンドフレームワークの選定において、当初はNext.jsを「エコシステムの総合力」で推奨する方向だったが、Prometheusのコンポーネントがほぼ全てカスタム実装である点を考慮し、SvelteKitの「フレームワーク自体の軽さ」がそのままユーザー体験に反映されると判断
- 当初はBunをランタイムとして採用し、bun:sqliteを活用する計画だった。しかしSvelteKitの開発サーバー（Vite）がNode.js上でSSRを実行するため、bun:sqliteはSSRコンテキストで利用できないことが判明。結果としてNode.js + better-sqlite3の構成に変更した。better-sqlite3はbun:sqliteと同等のSQLiteパフォーマンスを提供し、Node.js環境での安定性も高い
- パッケージマネージャはpnpmを採用。monorepoのワークスペース管理に優れ、onlyBuiltDependenciesによるネイティブモジュール（better-sqlite3、@node-rs/argon2）のビルド管理も適切に行われる
- CodeMirror 6のVimモードはNeovimの完全な再現ではないが、外部エディタとの併用（chokidarファイル監視）により、重い編集はNative Neovimで、軽い編集とプレビューはWeb UIでという使い分けが可能
