import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = ({ cookies }) => {
    cookies.delete("accesstoken", {path: "/"})
    cookies.delete("refreshtoken", {path: "/"})
    return { success: true }
}