import type { LayoutServerLoad } from "./$types.js";
import { VAULT_PATH } from "$lib/server/env.js";

export const load: LayoutServerLoad = async () => {
  return {
    vaultPath: VAULT_PATH,
  };
};
