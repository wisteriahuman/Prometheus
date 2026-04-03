import type { Actions } from "./$types.js";
import { fail, redirect } from "@sveltejs/kit";
import { createUser, emailExists, createSession } from "$lib/server/auth/session.js";

export const actions: Actions = {
  default: async ({ request, cookies }) => {
    const data = await request.formData();
    const email = data.get("email")?.toString()?.trim();
    const displayName = data.get("displayName")?.toString()?.trim();
    const password = data.get("password")?.toString();

    if (!email || !displayName || !password) {
      return fail(400, { error: "All fields are required" });
    }

    if (password.length < 8) {
      return fail(400, { error: "Password must be at least 8 characters" });
    }

    if (await emailExists(email)) {
      return fail(400, { error: "Email already registered" });
    }

    const userId = await createUser(email, displayName, password);
    const sessionId = await createSession(userId);

    cookies.set("session", sessionId, {
      path: "/",
      httpOnly: true,
      sameSite: "lax",
      secure: false, // Set to true in production with HTTPS
      maxAge: 30 * 24 * 60 * 60, // 30 days
    });

    throw redirect(303, "/");
  },
};
