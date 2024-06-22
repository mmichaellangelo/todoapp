import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = ({ cookies, locals }) => {
    cookies.delete("accesstoken", {path: "/"})
    cookies.delete("refreshtoken", {path: "/"})
    locals.userid = undefined;
    locals.username = undefined;
    return { success: true }
}