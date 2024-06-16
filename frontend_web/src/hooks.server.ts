import { getUsernameFromAccessToken } from "$lib/util/tokenValidation.server";
import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    console.log("< RUN > Server handle hook")

    const accessCookie = event.cookies.get("accesstoken");
    const refreshCookie = event.cookies.get("refreshtoken");

    var access = accessCookie as string;
    var refresh = refreshCookie as string;


    var username;
    try {
        username = getUsernameFromAccessToken(access);
        event.locals.username = username;
    } catch (err) {
        try {
            const res = await fetch("http://api/login/refresh/", {
                method: "POST",
                headers: {
                    "refreshtoken": refresh,
                }
            })
            if (!res.ok) {
                console.log("oops")
                console.log(res)
            }
            
            var newAccess = res.headers.get("accesstoken")
            console.log(`NEW ACCESS: ${newAccess}`)
            if (newAccess && typeof newAccess == 'string') {
                access = newAccess;
                event.cookies.set("accesstoken", access, {path: "/"})
            }
            
        } catch (err) {
            console.log(err)
        }
    }
    
    
    var response = await resolve(event);
    response.headers.set("accesstoken", access);
    response.headers.set("refreshtoken", refresh);
    return response;
}