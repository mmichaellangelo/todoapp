import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
    console.log("< RUN > page server load")
    return {
        username: locals.username
    }
}