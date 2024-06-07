// example load func from docs

import type { Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ cookies }) => {
	// Check refresh token. If not expired, redirect to home
    const refresh = cookies.get("refreshtoken")
    console.log(refresh)
};

export const actions = {
    default: async ({cookies, request}) => {
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
        console.log("Success!")
        const responseJson = await response.json()
        const accesstoken = responseJson.accesstoken
        const refreshtoken = responseJson.refreshtoken
        console.log("Access: ", accesstoken)
        console.log("Refresh:", refreshtoken)

        cookies.set("refreshtoken", refreshtoken, {path: "/", secure: true, httpOnly: true})
        return {success: true, accesstoken: accesstoken}

    }
} satisfies Actions;