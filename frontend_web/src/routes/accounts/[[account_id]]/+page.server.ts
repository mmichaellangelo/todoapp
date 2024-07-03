import { error, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({params, fetch}) => {
    const res = await fetch(`http://api/accounts/${params.account_id}`)
    if (!res.ok) {
        const errmsg = await res.text()
        throw error(res.status, errmsg)
    }
    const resJ = await res.json()
    return { account: {
        username: resJ.username
    } }
}

export const actions: Actions = {
    create: async ({request, fetch}) => {
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

        return { success: true }
    }
}