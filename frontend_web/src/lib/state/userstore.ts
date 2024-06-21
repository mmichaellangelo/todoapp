import { writable, type Writable } from "svelte/store";

import type { SessionData } from "$lib/types";


export const Session = writable<SessionData | undefined>(undefined)
