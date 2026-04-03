import type { Actions } from "./$types.js";
import { redirect } from "@sveltejs/kit";
import { invalidateSession } from "$lib/server/auth/session.js";

export const actions: Actions = {
  default: async ({ cookies }) => {
    const sessionId = cookies.get("session");
    if (sessionId) {
      await invalidateSession(sessionId);
      cookies.delete("session", { path: "/" });
    }

    throw redirect(303, "/auth/login");
  },
};
