import { Username } from "$lib/state/userstore";
import { getUsernameFromAccessToken } from "$lib/util/tokenValidation.server";
import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {

    const accessCookie = event.cookies.get("accesstoken");
    const refreshCookie = event.cookies.get("refreshtoken");

    var username;
    try {
        username = getUsernameFromAccessToken(accessCookie as string);
        event.locals.username = username;
    } catch (err) {
        console.log(err);
    }
    
    var response = await resolve(event);
    response.headers.set("accesstoken", accessCookie as string);
    response.headers.set("refreshtoken", refreshCookie as string);
    return response;
}