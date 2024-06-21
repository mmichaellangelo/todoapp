import { getSessionDataFromToken } from '$lib/util/tokenValidation.server.js';

export const actions = {
    default: async ({request, fetch, cookies}) => {
        let formData;
        try {
            formData = await request.formData();
        } catch(err) {
            console.log("error parsing form data: ", err)
            return {success: false}
        }


        const response = await fetch("http://api/accounts/", {
            method: "POST",
            body: formData
        })
        if (!response.ok) {
            console.log(response)
            return {success: false}
        }
        const responseJson = await response.json()
        console.log(responseJson)

        const accesstoken = response.headers.get("accesstoken")
        const refreshtoken = response.headers.get("refreshtoken")

        if ((!accesstoken || !refreshtoken) || (accesstoken == "" || refreshtoken == "")) {
            return { success: false }
        }
        cookies.set("accesstoken", accesstoken, {path: "/", httpOnly: true})
        cookies.set("refreshtoken", refreshtoken, {path: "/", httpOnly: true})

        const sessionData = getSessionDataFromToken(accesstoken)
        const username = sessionData.username
        const userid = sessionData.userid

        return { success: true, username: username, userid: userid}
    }

}