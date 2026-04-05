# ADR-005: エクスポート戦略

## 背景・課題 (Background/Problem)
- ノートを他者に共有したり、外部ツールで利用するためにエクスポート機能が必要
- PDF、HTML、Markdownの3形式が候補
- テーマ（ダーク/ライト/カスタム）がエクスポートに反映される必要がある
- Goシングルバイナリの軽量さを維持したい（外部依存の追加を最小限に）

## 決定事項 (Decision)

3つのエクスポート形式を提供:

| 形式 | 方式 | 用途 |
|---|---|---|
| **HTML** | Go APIがテーマ付きスタンドアロンHTMLを生成 → ダウンロード | 共有、アーカイブ |
| **Markdown** | フロントマターなしのプレーンMD → ダウンロード | 他ツールへの移行 |
| **PDF** | テーマ付きHTMLを新タブで開く → ブラウザ印刷ダイアログ | 印刷、フォーマルな共有 |

### テーマ解決

エクスポート時のテーマは以下の優先順位で決定:
1. URLクエリパラメータ `?theme=ocean`
2. ノートのfrontmatter `theme` フィールド
3. アプリ全体テーマ（フロントエンドから渡される）
4. フォールバック: `light`

### エクスポートHTML構成

```html
<style>
  :root { --bg: #...; --text: #...; ... }  /* テーマ色 */
  * { print-color-adjust: exact !important; }  /* 背景色印刷 */
  @page { margin: 15mm; }  /* 印刷余白 */
  @media print { ... }  /* 印刷最適化 */
</style>
<body>
  <!-- goldmark変換済みHTML -->
  <footer>Exported from Prometheus</footer>
</body>
```

### 理由 (Reasons)
- **HTMLが最も汎用的**: ブラウザで開ける、スタイル付き、テーマ反映、軽量
- **PDFはブラウザ印刷に委譲**: Go側でPDF生成ライブラリを入れると、CJKフォント埋め込みでバイナリが+5MB以上。ブラウザの印刷エンジンが最も高品質（ベクター描画、テキスト選択可能）
- **Markdownプレーンは移行用**: フロントマターなしで他ツールにそのまま貼れる

### 受け入れるトレードオフ (Accepted Trade-offs)
- **PDFのヘッダー/フッター**: ブラウザの印刷設定に依存。ユーザーが「詳細設定」でヘッダーとフッターをオフにする必要がある
- **PDFの背景色**: `print-color-adjust: exact`をCSSに設定しているが、一部ブラウザでは「背景のグラフィック」を手動でオンにする必要がある

## 検討した別の選択肢 (Alternatives Considered)

### html2pdf.js（クライアントサイドPDF生成）
- **メリット**: ヘッダー/フッターなし、完全制御
- **デメリット**: html2canvasがテキストをラスタライズ → PDFからテキストコピー不可。+200KBのバンドル増加
- **不採用理由**: テキスト選択不可は致命的

### Go側PDF生成（gpdf、gofpdf等）
- **メリット**: サーバー側で完結、ブラウザ不要
- **デメリット**: CJKフォント埋め込みでバイナリ+5MB以上。Markdown→PDFのスタイリング手動実装
- **不採用理由**: シングルバイナリの軽量さを犠牲にする

## 参考 (References)
- [print-color-adjust CSS](https://developer.mozilla.org/en-US/docs/Web/CSS/print-color-adjust)
- [@page CSS](https://developer.mozilla.org/en-US/docs/Web/CSS/@page)

## 議論 (Discussion)
- PDFエクスポートは一度削除したが、ユーザーからの「ブラウザの詳細設定で十分対応できる」というフィードバックにより復活
- Go側にビルトインテーマ色定義（`internal/service/themes.go`）を持ち、フロントエンドと同期。エクスポートHTMLにCSS変数として埋め込む
- wikilinkのHTML変換がコードブロック内でも走る問題があった。goldmark処理後のHTMLに対してwikilink変換を行い、`<code>`/`<pre>`タグ内はスキップするように修正
