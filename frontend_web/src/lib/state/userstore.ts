import { writable, type Writable } from "svelte/store";

export const Username = writable<string | undefined>(undefined)
export const AccessToken = writable("")