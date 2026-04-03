import { initVault } from "$lib/server/fs/init-vault.js";
import { indexAllNotes } from "$lib/server/indexer.js";

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
