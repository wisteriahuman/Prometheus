import type { PageServerLoad } from "./$types.js";
import { VAULT_PATH } from "$lib/server/env.js";

export const load: PageServerLoad = async () => {
  return {
    vaultPath: VAULT_PATH,
  };
};
