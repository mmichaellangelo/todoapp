import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
    return {
        username: locals.username,
        userid: locals.userid
    }
}