export function load() {

}

export const actions = {
    default: async ({request, fetch, cookies}) => {
        let formData;
        try {
            formData = await request.formData();
        } catch(err) {
            console.log("error parsing form data: ", err)
            return {success: false}
        }


        const response = await fetch("http://api:8080/accounts/", {
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