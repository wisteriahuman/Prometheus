import { writable } from "svelte/store";

export const sidebarOpen = writable(true);

export function toggleSidebar() {
  sidebarOpen.update((v) => !v);
}
