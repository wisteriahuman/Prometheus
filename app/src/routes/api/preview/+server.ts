import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { markdownToHtml } from "@prometheus/core/markdown";

export const POST: RequestHandler = async ({ request }) => {
  const { content } = await request.json();

  if (typeof content !== "string") {
    return json({ error: "content is required" }, { status: 400 });
  }

  const html = await markdownToHtml(content);
  return json({ html });
};
