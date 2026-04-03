import { db, schema } from "../db/index.js";
import { eq } from "drizzle-orm";
import { hash, verify } from "@node-rs/argon2";
import { ulid } from "ulid";
import { randomBytes } from "node:crypto";

const SESSION_DURATION = 30 * 24 * 60 * 60 * 1000; // 30 days

function generateSessionId(): string {
  return randomBytes(32).toString("hex");
}

export async function createUser(
  email: string,
  displayName: string,
  password: string,
): Promise<string> {
  const id = ulid();
  const passwordHash = await hash(password);

  await db.insert(schema.users).values({
    id,
    email: email.toLowerCase(),
    displayName,
    passwordHash,
    createdAt: new Date(),
  });

  return id;
}

export async function verifyPassword(
  email: string,
  password: string,
): Promise<string | null> {
  const [user] = await db
    .select()
    .from(schema.users)
    .where(eq(schema.users.email, email.toLowerCase()))
    .limit(1);

  if (!user) return null;

  const valid = await verify(user.passwordHash, password);
  if (!valid) return null;

  return user.id;
}

export async function createSession(userId: string): Promise<string> {
  const id = generateSessionId();
  const expiresAt = new Date(Date.now() + SESSION_DURATION);

  await db.insert(schema.sessions).values({
    id,
    userId,
    expiresAt,
  });

  return id;
}

export async function validateSession(
  sessionId: string,
): Promise<{ userId: string; expiresAt: Date } | null> {
  const [session] = await db
    .select()
    .from(schema.sessions)
    .where(eq(schema.sessions.id, sessionId))
    .limit(1);

  if (!session) return null;

  if (session.expiresAt < new Date()) {
    await db.delete(schema.sessions).where(eq(schema.sessions.id, sessionId));
    return null;
  }

  // Extend session if more than half expired
  const halfLife = SESSION_DURATION / 2;
  if (session.expiresAt.getTime() - Date.now() < halfLife) {
    const newExpiry = new Date(Date.now() + SESSION_DURATION);
    await db
      .update(schema.sessions)
      .set({ expiresAt: newExpiry })
      .where(eq(schema.sessions.id, sessionId));
    return { userId: session.userId, expiresAt: newExpiry };
  }

  return { userId: session.userId, expiresAt: session.expiresAt };
}

export async function invalidateSession(sessionId: string): Promise<void> {
  await db.delete(schema.sessions).where(eq(schema.sessions.id, sessionId));
}

export async function getUserById(userId: string) {
  const [user] = await db
    .select({
      id: schema.users.id,
      email: schema.users.email,
      displayName: schema.users.displayName,
      createdAt: schema.users.createdAt,
    })
    .from(schema.users)
    .where(eq(schema.users.id, userId))
    .limit(1);

  return user ?? null;
}

export async function emailExists(email: string): Promise<boolean> {
  const [user] = await db
    .select({ id: schema.users.id })
    .from(schema.users)
    .where(eq(schema.users.email, email.toLowerCase()))
    .limit(1);

  return !!user;
}
