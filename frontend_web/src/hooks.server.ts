import { getSessionDataFromToken } from "$lib/util/tokenValidation.server";
import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const accessCookie = event.cookies.get("accesstoken");
    const refreshCookie = event.cookies.get("refreshtoken");

    var access = accessCookie as string;
    var refresh = refreshCookie as string;
    var sessionData;
    
    try {
        sessionData = getSessionDataFromToken(access);
        event.locals.username = sessionData.username;
        event.locals.userid = sessionData.userid;
    } catch (err) {
        try {
            console.log("REFRESH")
            const res = await fetch("http://api/login/refresh/", {
                method: "POST",
                headers: {
                    "refreshtoken": refresh,
                }
            })
            if (!res.ok) {
                throw Error("unable to refresh")
            }
            
            var newAccess = res.headers.get("accesstoken")
            console.log(`NEW ACCESS: ${newAccess}`)
            if (newAccess && typeof newAccess == 'string') {
                access = newAccess;
                event.cookies.set("accesstoken", access, {path: "/"})
                sessionData = getSessionDataFromToken(newAccess);
                event.locals.username = sessionData.username;
                event.locals.userid = sessionData.userid;
            }
            
        } catch (err) {
            console.log(typeof err)
        }
    }
    
    
    var response = await resolve(event);
    response.headers.set("accesstoken", access);
    response.headers.set("refreshtoken", refresh);
    return response;
}