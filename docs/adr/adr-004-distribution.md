# ADR-004: 配信戦略

## 背景・課題 (Background/Problem)
- Prometheusをユーザーが簡単にインストールできるようにしたい
- Go製のシングルバイナリなので複数の配信方法が可能
- Mac (Apple Silicon / Intel)、Linux、Windowsに対応する必要がある
- フロントエンド（Svelte SPA）をバイナリに同梱する必要がある

## 決定事項 (Decision)

3つの配信方法を提供:

| 方法 | コマンド | 対象 |
|---|---|---|
| **Homebrew** | `brew install wisteriahuman/tap/prometheus` | macOS / Linux |
| **go install** | `go install github.com/wisteriahuman/prometheus/cmd/prometheus@latest` | Go開発者 |
| **GitHub Releases** | バイナリ直接ダウンロード | 全プラットフォーム |

### ビルド・リリースフロー

```
git tag v0.x.x → GitHub Actions → GoReleaser
                                    ├── pnpm build (Svelte SPA)
                                    ├── go build (5プラットフォーム)
                                    ├── GitHub Releases (バイナリ + チェックサム)
                                    └── homebrew-tap (Formula自動更新)
```

### 対応プラットフォーム

| OS | Arch | 形式 |
|---|---|---|
| macOS | arm64 (Apple Silicon) | tar.gz |
| macOS | x86_64 (Intel) | tar.gz |
| Linux | arm64 | tar.gz |
| Linux | x86_64 | tar.gz |
| Windows | x86_64 | zip |

### GoReleaser設定

- `CGO_ENABLED=0`: pure Goビルド（modernc.org/sqliteはcgo不要）
- before hooks: `sh -c "cd web && pnpm install && pnpm run build"` でSPAをビルド
- `//go:embed`でSPAビルド成果物をバイナリに同梱
- `brews`セクションで`wisteriahuman/homebrew-tap`に自動push

### go install対応

- SPAビルド済みファイル（`internal/server/web/`）をgitにコミット
- `go install`はソースから直接ビルドするため、embedするファイルがリポジトリに必要
- `.gitignore`で除外しない

### 理由 (Reasons)
- **Homebrew**: macOSユーザーの最も自然なインストール方法。1コマンドで完了
- **go install**: Go開発者はこれが一番速い。ソースからビルドするので常に最新
- **GitHub Releases**: Goが入っていないユーザーでもバイナリをダウンロードして即利用可能
- **GoReleaser**: タグpushだけで全プラットフォームのビルド+リリース+brew更新が自動完了

### 受け入れるトレードオフ (Accepted Trade-offs)
- **SPAビルドファイルのgitコミット**: リポジトリサイズが増える（~2MB）。`go install`対応に必要
- **GitHub Actionsの実行時間**: SPAビルド + Goクロスコンパイルで~3分。許容範囲
- **HOMEBREW_TAP_TOKEN**: GitHub PATが必要。90日で期限切れ、更新が必要

## 検討した別の選択肢 (Alternatives Considered)

### npmパッケージ化
- **メリット**: `npx prometheus-notes dev`で起動
- **デメリット**: better-sqlite3のネイティブビルド問題。Node.jsランタイムが必要
- **不採用理由**: Goに移行したことで不要になった

### Docker配信
- **メリット**: 環境依存問題ゼロ
- **デメリット**: Docker Desktopのインストールが必要（重い）
- **不採用理由**: 個人用ノートアプリにDockerは過剰

### Goバイナリ + 別途SPAダウンロード
- **メリット**: リポジトリサイズが小さい
- **デメリット**: インストール手順が2段階になる。UXが悪い
- **不採用理由**: シングルバイナリの利点が失われる

## 参考 (References)
- [GoReleaser](https://goreleaser.com/)
- [GoReleaser Homebrew](https://goreleaser.com/customization/homebrew/)
- [GitHub Actions](https://docs.github.com/en/actions)
- [wisteriahuman/homebrew-tap](https://github.com/wisteriahuman/homebrew-tap)

## 議論 (Discussion)
- GoReleaserのbefore hookで`cd && pnpm build`を使うには`sh -c`ラップが必要（GoReleaserはシェルではなくexecで実行するため）
- `pnpm/action-setup@v4`でバージョン指定すると`package.json`の`packageManager`と競合する。バージョン指定を省略して自動検出させる
- `.gitignore`の`prometheus`パターンが`cmd/prometheus/`にマッチする問題があった。`/prometheus`（ルートのみ）に修正
- Windowsのarm64ビルドは需要が少ないため除外。x86_64のみ
