import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq, and } from "drizzle-orm";

// GET /api/workspaces/:id
export const GET: RequestHandler = async ({ params, locals }) => {
  if (!locals.user) throw error(401, "Authentication required");

  const [workspace] = await db
    .select()
    .from(schema.workspaces)
    .where(eq(schema.workspaces.id, params.id))
    .limit(1);

  if (!workspace) throw error(404, "Workspace not found");

  // Check membership
  const [member] = await db
    .select()
    .from(schema.workspaceMembers)
    .where(
      and(
        eq(schema.workspaceMembers.workspaceId, params.id),
        eq(schema.workspaceMembers.userId, locals.user.id),
      ),
    )
    .limit(1);

  if (!member && workspace.ownerId !== locals.user.id) {
    throw error(403, "Access denied");
  }

  // Get members
  const members = await db
    .select({
      userId: schema.workspaceMembers.userId,
      role: schema.workspaceMembers.role,
      email: schema.users.email,
      displayName: schema.users.displayName,
    })
    .from(schema.workspaceMembers)
    .innerJoin(schema.users, eq(schema.workspaceMembers.userId, schema.users.id))
    .where(eq(schema.workspaceMembers.workspaceId, params.id));

  return json({ ...workspace, members });
};

// DELETE /api/workspaces/:id
export const DELETE: RequestHandler = async ({ params, locals }) => {
  if (!locals.user) throw error(401, "Authentication required");

  const [workspace] = await db
    .select()
    .from(schema.workspaces)
    .where(eq(schema.workspaces.id, params.id))
    .limit(1);

  if (!workspace) throw error(404, "Workspace not found");
  if (workspace.ownerId !== locals.user.id) throw error(403, "Only the owner can delete");

  await db.delete(schema.workspaces).where(eq(schema.workspaces.id, params.id));

  return json({ success: true });
};
