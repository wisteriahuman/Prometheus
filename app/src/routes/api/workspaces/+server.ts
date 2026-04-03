import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq, or } from "drizzle-orm";
import { ulid } from "ulid";

// GET /api/workspaces — list user's workspaces
export const GET: RequestHandler = async ({ locals }) => {
  if (!locals.user) {
    throw error(401, "Authentication required");
  }

  const workspaces = await db
    .select({
      id: schema.workspaces.id,
      name: schema.workspaces.name,
      engine: schema.workspaces.engine,
      theme: schema.workspaces.theme,
      createdAt: schema.workspaces.createdAt,
      ownerId: schema.workspaces.ownerId,
    })
    .from(schema.workspaces)
    .leftJoin(
      schema.workspaceMembers,
      eq(schema.workspaces.id, schema.workspaceMembers.workspaceId),
    )
    .where(
      or(
        eq(schema.workspaces.ownerId, locals.user.id),
        eq(schema.workspaceMembers.userId, locals.user.id),
      ),
    );

  return json(workspaces);
};

// POST /api/workspaces — create a workspace
export const POST: RequestHandler = async ({ request, locals }) => {
  if (!locals.user) {
    throw error(401, "Authentication required");
  }

  const { name, engine, dbUrl } = await request.json();

  if (!name || typeof name !== "string") {
    throw error(400, "name is required");
  }

  const validEngine = engine === "postgresql" ? "postgresql" : "sqlite";

  const id = ulid();
  const vaultPath = `./vaults/${id}`;

  await db.insert(schema.workspaces).values({
    id,
    name,
    vaultPath,
    engine: validEngine,
    dbUrl: validEngine === "postgresql" ? dbUrl : null,
    createdAt: new Date(),
    ownerId: locals.user.id,
  });

  // Add owner as member
  await db.insert(schema.workspaceMembers).values({
    workspaceId: id,
    userId: locals.user.id,
    role: "owner",
  });

  return json({ id, name, engine: validEngine }, { status: 201 });
};
