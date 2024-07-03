import { goto } from "$app/navigation";
import { getSessionDataFromToken } from "$lib/util/tokenValidation.server";
import { error, redirect, type Handle, type HandleFetch, type HandleServerError } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const accessCookie = event.cookies.get("accesstoken");
    const refreshCookie = event.cookies.get("refreshtoken");

    var access = accessCookie as string;
    var refresh = refreshCookie as string;
    var sessionData;

    if (!access || access == "") {
        event.locals.userid = undefined;
        event.locals.username = undefined;
    } else {
        try {
            sessionData = getSessionDataFromToken(access);
            event.locals.username = sessionData.username;
            event.locals.userid = sessionData.userid;
        } catch (err) {
            try {
                const res = await fetch("http://api/login/refresh/", {
                    method: "POST",
                    headers: {
                        "refreshtoken": refresh,
                    }
                })
                if (!res.ok) {
                    event.cookies.delete("refreshtoken", {path: "/"})
                    event.cookies.delete("accesstoken", {path: "/"})
                    throw Error("unable to refresh")
                    
                }
                
                var newAccess = res.headers.get("accesstoken")
                if (newAccess && typeof newAccess == 'string') {
                    access = newAccess;
                    event.cookies.set("accesstoken", access, {path: "/"})
                    sessionData = getSessionDataFromToken(newAccess);
                    event.locals.username = sessionData.username;
                    event.locals.userid = sessionData.userid;
                }
                
            } catch (err) {
                console.log(err)
                throw redirect(302, "/login")
            }
        }
    }
    var response = await resolve(event);
    response.headers.set("X-AccessToken", access);
    response.headers.set("X-RefreshToken", refresh);
    return response;
}

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
    const accessCookie = event.cookies.get("accesstoken") as string;

    console.log(`ACCESS: ${accessCookie}`)

    if (request.url.startsWith("http://api/")) {
        request.headers.set("accesstoken", accessCookie)
    }
    return fetch(request)
}

// export const handleError: HandleServerError = async ({event}) => {
//     event.request.text()

//     return {
//         message: "oopsie"
//     }
// }