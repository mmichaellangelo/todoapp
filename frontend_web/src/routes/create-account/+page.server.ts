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

        // const email = data.get("email");
        // const username = data.get("username");
        // const password = data.get("password");
        

        const response = await fetch("http://api:8080/accounts/", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: formData
        })
        if (!response.ok) {
            console.log(response)
            return {success: false}
        }

    }

}