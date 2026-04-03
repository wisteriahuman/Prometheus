import type { Root, ListItem } from "mdast";
import { visit } from "unist-util-visit";

export interface TaskInfo {
  content: string;
  completed: boolean;
  lineNumber: number;
}

export function extractTasks(ast: Root): TaskInfo[] {
  const tasks: TaskInfo[] = [];

  visit(ast, "listItem", (node: ListItem) => {
    if (node.checked === null || node.checked === undefined) return;

    const content = getTextContent(node).trim();
    const lineNumber = node.position?.start.line ?? 0;

    tasks.push({
      content,
      completed: node.checked,
      lineNumber,
    });
  });

  return tasks;
}

function getTextContent(node: any): string {
  if (node.type === "text") return node.value;
  if (node.children) {
    return node.children.map(getTextContent).join("");
  }
  return "";
}
