import type { LayoutServerLoad } from "./$types.js";
import { VAULT_PATH } from "$lib/server/env.js";

export const load: LayoutServerLoad = async ({ locals }) => {
  return {
    user: locals.user ?? null,
    vaultPath: VAULT_PATH,
  };
};
