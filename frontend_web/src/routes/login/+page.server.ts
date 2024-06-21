import { redirect, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

import jwt, { type JwtPayload } from "jsonwebtoken";

interface MyPayload extends JwtPayload {
    username: string;
}

export const load: PageServerLoad = async ({ cookies, locals }) => {
	// Check refresh token. If not expired, redirect to home
    const access = cookies.get("accesstoken") as string;
    if (access == "") {
        return;
    }
    jwt.verify(access, "secret key", (err, decoded) => {
        if (err) {
            switch (err.name) {
                case "TokenExpiredError":
                    console.log("Token expired. Deleting.");
                    cookies.delete("refresh", {path: "/"});
                    return { username: undefined, userid: undefined }
            }   
            return { username: undefined, user_id: undefined}
        }
        if (typeof decoded == 'object') {
            console.log("DECODED USERNAME:", decoded.username)
            return { username: decoded.username, userid: decoded.userid }
        } else {
            return { username: undefined, user_id: undefined }
        }
        
    })   
}
    

export const actions = {
    default: async ({ fetch, request, cookies, locals }) => {
        let formData;
        try {
            formData = await request.formData();
        } catch (err) {
            console.log("Error parsing form data:", err);
            return {success: false}
        }

        const response = await fetch("http://api/login/", {
            method: "POST",
            body: formData,
        })


        if (!response.ok) {
            console.log("Error logging in:", response)
            return { success: false }
        }
        
        const accesstoken = response.headers.get("accesstoken")
        const refreshtoken = response.headers.get("refreshtoken")

        if (!accesstoken || !refreshtoken) {
            return { success: false }
        }
        if (accesstoken == "" || refreshtoken == "") {
            return { success: false }
        }
        
        let username;
        jwt.verify(accesstoken, "secret key", (err, decoded) => {
            if (err) {
                return { success: false };
            }
            let data = decoded as MyPayload;
            console.log(data.username)
            if (data.username == "") {
                return { success: false }
            }
            username = data.username;
        })   

        cookies.set("refreshtoken", refreshtoken, {path: "/", secure: true, httpOnly: true, sameSite: "lax"})
        cookies.set("accesstoken", accesstoken, {path: "/", secure: true, httpOnly: true, sameSite: "lax"})
        locals.username = username;
        return {success: true, username: username}

    }
} satisfies Actions;