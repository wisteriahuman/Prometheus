import type { Actions } from "./$types.js";
import { fail, redirect } from "@sveltejs/kit";
import { verifyPassword, createSession } from "$lib/server/auth/session.js";

export const actions: Actions = {
  default: async ({ request, cookies }) => {
    const data = await request.formData();
    const email = data.get("email")?.toString()?.trim();
    const password = data.get("password")?.toString();

    if (!email || !password) {
      return fail(400, { error: "Email and password are required" });
    }

    const userId = await verifyPassword(email, password);
    if (!userId) {
      return fail(400, { error: "Invalid email or password" });
    }

    const sessionId = await createSession(userId);

    cookies.set("session", sessionId, {
      path: "/",
      httpOnly: true,
      sameSite: "lax",
      secure: false,
      maxAge: 30 * 24 * 60 * 60,
    });

    throw redirect(303, "/");
  },
};
