import { redirect, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

import jwt, { type JwtPayload } from "jsonwebtoken";

interface MyPayload extends JwtPayload {
    username: string;
}

export const load: PageServerLoad = async ({ cookies }) => {
	// Check refresh token. If not expired, redirect to home
    const refresh = cookies.get("refreshtoken") as string;
    if (refresh == "") {
        return;
    }
    jwt.verify(refresh, "secret key", (err, decoded) => {
        if (err) {
            switch (err.name) {
                case "TokenExpiredError":
                    console.log("Token expired. Deleting.");
                    cookies.delete("refresh", {path: "/"});
                    break;
            }   
            return;
        }
    })   
}
    

export const actions = {
    default: async ({ fetch, cookies, request }) => {
        let formData;
        try {
            formData = await request.formData();
        } catch (err) {
            console.log("Error parsing form data:", err);
            return {success: false}
        }

        const response = await fetch("http://api:8080/login/", {
            method: "POST",
            body: formData
        })
        if (!response.ok) {
            console.log("Error logging in:", response)
            return { success: false }
        }

        const responseJson = await response.json();
        console.log(responseJson)
        const accesstoken = responseJson.accesstoken as string;
        const refreshtoken = responseJson.refreshtoken as string;

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

        cookies.set("refreshtoken", refreshtoken, {path: "/", secure: true, httpOnly: true})
        return {success: true, accesstoken: accesstoken, username: username}

    }
} satisfies Actions;