import { setContext } from "svelte"
import type { PageLoad } from "./$types"

export const load: PageLoad = async ({ data }) => {
    console.log("< RUN > page load")
    console.log(`username is ${data.username}`)

    return { username: data.username }
}