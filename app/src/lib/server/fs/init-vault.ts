import { mkdir, writeFile, readdir } from "node:fs/promises";
import { resolve, join } from "node:path";
import { ulid } from "ulid";
import { getVaultPath } from "./vault.js";

const SAMPLE_NOTES: { path: string; title: string; tags: string[]; content: string }[] = [
  {
    path: "welcome.md",
    title: "Prometheusへようこそ",
    tags: ["getting-started", "prometheus"],
    content: `
# Prometheusへようこそ

**Prometheus** はMarkdownベースのノートアプリです。

## 基本操作

| キー | 機能 |
|------|------|
| \`⌘K\` | コマンドパレット |
| \`⌘P\` | ファイル切替 |
| \`⌘S\` | 保存 |
| \`⌘D\` | デイリーノート |
| \`⌘G\` | グラフビュー |
| \`⌘.\` | テーマ切替 |
| \`⌘K\` | ショートカット一覧（コマンドパレットから） |

## はじめに

- \`[[wikilink]]\` でノート同士をリンクできます
- サイドバーの \`+\` から新規ノート作成
- タグはノートヘッダーの \`+\` から追加
- 右クリックでノートのリネーム・削除

詳しくは [[setup]] をご覧ください。
[[ideas]] も参考にどうぞ。

## Mermaid図のサンプル

\`\`\`mermaid
graph LR
  A[ノート作成] --> B[リンクを貼る]
  B --> C[グラフで可視化]
  C --> D[知識が繋がる]
\`\`\`

> *ノートは .md ファイルとして保存されます。Neovimで直接編集可能です。*
`,
  },
  {
    path: "setup.md",
    title: "セットアップ",
    tags: ["getting-started", "guide"],
    content: `
# セットアップ

[[welcome]] の使い方ガイドです。

## ノートの保存先

ノートは \`.env\` で設定したディレクトリに \`.md\` ファイルとして保存されます。

\`\`\`bash
# app/.env
PROMETHEUS_VAULT_PATH=~/my-notes  # 好きなパスに変更可能
\`\`\`

- Neovimやお好みのエディタで直接編集できます
- 変更はPrometheusに自動反映されます
- Gitで管理すればバージョン管理も可能です

## ファイル名とタイトル

- **ファイル名**: \`setup.md\` のようにURLやリンクで使う名前
- **タイトル**: フロントマターの \`title\` で自由に設定。日本語OK
- ファイル名は英語、タイトルは日本語が推奨

## リンクの書き方

\`\`\`markdown
[[ファイル名]]           → そのファイルへリンク
[[ファイル名|表示テキスト]] → 別名でリンク
\`\`\`

例: \`[[setup]]\` → このノートへのリンク

## タグ

ノートヘッダーの \`+\` ボタンからタグを追加できます。
グラフビューではタグごとにノートがグループ化されます。

[[welcome]] | [[ideas]]
`,
  },
  {
    path: "ideas.md",
    title: "アイデア",
    tags: ["ideas"],
    content: `
# アイデア

ひらめきを逃さないためのノート。

## やりたいこと

- [ ] プロジェクトのアイデアを整理する
- [ ] 読書メモをまとめる
- [ ] 学習ログをつける

## フロー

\`\`\`mermaid
graph TD
  A[ひらめき] --> B{重要？}
  B -->|Yes| C[ノートに書く]
  B -->|No| D[忘れてOK]
  C --> E[タグをつける]
  E --> F[リンクで繋げる]
\`\`\`

## メモ

ここに自由にメモを書いてください。

[[welcome]] | [[setup]]
`,
  },
];

function createFrontmatter(title: string, tags: string[]): string {
  const now = new Date().toISOString();
  return `---
id: "${ulid()}"
title: "${title}"
created: "${now}"
modified: "${now}"
tags: [${tags.join(", ")}]
---
`;
}

/**
 * Initialize a vault with sample notes if empty.
 */
export async function initVault(): Promise<boolean> {
  const vaultPath = await getVaultPath();

  // Check if vault already has .md files
  try {
    const entries = await readdir(vaultPath);
    const hasNotes = entries.some(
      (e) => e.endsWith(".md") || (e !== ".prometheus" && !e.startsWith(".")),
    );
    if (hasNotes) return false; // Already initialized
  } catch {
    // Directory doesn't exist yet, will be created
  }

  // Create .prometheus config
  const configDir = join(vaultPath, ".prometheus", "themes");
  await mkdir(configDir, { recursive: true });
  await writeFile(
    join(vaultPath, ".prometheus", "config.json"),
    JSON.stringify(
      {
        name: "My Notes",
        engine: "sqlite",
        theme: "prometheus",
        dailyNoteTemplate:
          '---\ntitle: "{{date}}"\ntags: [daily]\n---\n\n# {{dateJa}}\n\n## タスク\n\n- [ ] \n\n## メモ\n\n',
      },
      null,
      2,
    ),
  );

  // Create sample notes
  for (const note of SAMPLE_NOTES) {
    const fullPath = resolve(vaultPath, note.path);
    const content = createFrontmatter(note.title, note.tags) + note.content;
    await writeFile(fullPath, content, "utf-8");
  }

  return true;
}
