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
| Runtime | Bun | 1.x |
| CSS Framework | Tailwind CSS | v4 |
| Markdown Editor | CodeMirror 6 | - |
| Markdown Parser | unified + remark + rehype | - |
| Graph Visualization | D3.js (d3-force) | - |
| ORM | Drizzle ORM | - |
| Validation | Zod | - |

### 理由 (Reasons)

#### SvelteKit

- **バンドルサイズ最小**: Svelteはコンパイル時にVanilla JSへ変換される。Reactのような仮想DOMランタイム（~80KB gzip）が不要
- **カスタムアプリに最適**: Prometheusのコアコンポーネント（エディタ、グラフ、コマンドパレット）は全てカスタム実装。既存UIコンポーネントライブラリへの依存が少ないため、Reactエコシステムの優位性が活きない
- **トランジション・アニメーション組み込み**: ノート切替、サイドバー開閉、テーマ切替のアニメーションがフレームワークレベルでサポートされる
- **リアクティビティがシンプル**: Svelte 5のRunes（`$state`, `$derived`, `$effect`）はReactのuseState/useEffect/useMemoより直感的で、バグが入り込む余地が少ない
- **コード量削減**: 同等の機能をReactより少ないコードで実装でき、メンテナンス性が高い

#### Bun

- **SQLite内蔵** (`bun:sqlite`): better-sqlite3のインストールやネイティブバインディングの問題が不要。ノートアプリのインデックスDBとして即座に利用可能
- **ファイルI/O高速**: .mdファイルの読み書きが頻繁なアプリケーションにおいて、Node.jsの数倍の速度
- **起動速度**: 開発時のHMRサイクルが高速
- **TypeScript直接実行**: マイグレーションスクリプト等をトランスパイルなしで実行可能

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
- **Bunの安定性**: Node.jsほど枯れていない。一部のnpmパッケージとの互換性に問題が生じる可能性がある
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

### Node.js (Runtime)
- **メリット**: 最も安定、全パッケージとの互換性保証
- **デメリット**: SQLiteは別途better-sqlite3が必要、起動・I/O速度がBunに劣る
- **不採用理由**: Bunの速度優位が.mdファイル中心のアプリで顕著に効く

## 参考 (References)
- [SvelteKit Documentation](https://svelte.dev/docs/kit)
- [Bun Documentation](https://bun.sh/docs)
- [Tailwind CSS v4](https://tailwindcss.com/docs)
- [CodeMirror 6](https://codemirror.net/)
- [D3 Force-Directed Graph](https://d3js.org/d3-force)
- [Drizzle ORM](https://orm.drizzle.team/)

## 議論 (Discussion)
- フロントエンドフレームワークの選定において、当初はNext.jsを「エコシステムの総合力」で推奨する方向だったが、Prometheusのコンポーネントがほぼ全てカスタム実装である点を考慮し、SvelteKitの「フレームワーク自体の軽さ」がそのままユーザー体験に反映されると判断
- Bunの採用はSQLite内蔵が決定的だった。bun:sqliteはbetter-sqlite3と異なりネイティブバインディングのコンパイルが不要で、Dockerイメージのサイズ削減とCI/CDの簡素化にも寄与する
- CodeMirror 6のVimモードはNeovimの完全な再現ではないが、外部エディタとの併用（chokidarファイル監視）により、重い編集はNative Neovimで、軽い編集とプレビューはWeb UIでという使い分けが可能
