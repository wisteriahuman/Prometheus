import type { Handle } from "@sveltejs/kit";
import { initVault } from "$lib/server/fs/init-vault.js";
import { indexAllNotes } from "$lib/server/indexer.js";
import { validateSession, getUserById } from "$lib/server/auth/session.js";

// Initialize vault and index on startup
(async () => {
  const initialized = await initVault();
  if (initialized) {
    console.log("Vault initialized with sample notes");
  }

  const count = await indexAllNotes();
  console.log(`Indexed ${count} notes`);
})().catch((e) => {
  console.error("Startup error:", e);
});

export const handle: Handle = async ({ event, resolve }) => {
  const sessionId = event.cookies.get("session");

  if (sessionId) {
    const session = await validateSession(sessionId);
    if (session) {
      const user = await getUserById(session.userId);
      if (user) {
        event.locals.user = {
          id: user.id,
          email: user.email,
          displayName: user.displayName,
        };
        event.locals.session = { id: sessionId, expiresAt: session.expiresAt };
      }
    } else {
      event.cookies.delete("session", { path: "/" });
    }
  }

  return resolve(event);
};
