import type { PageLoad } from "./$types"

export const load: PageLoad = async ({ parent }) => {
    console.log("< RUN > page load")
    
    const data = await parent();
    return { username: data.username }
}