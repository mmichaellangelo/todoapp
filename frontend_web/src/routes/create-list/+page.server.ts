import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ fetch, locals }) => {
    console.log(`USERID IS ${locals.userid}`)
    const res = await fetch(`http://api/accounts/${locals.userid}/lists/`,
        {
            method: "POST"
        }
    )
    if (!res.ok) {
        const msg = await res.text()
        console.log(msg)
    } else {
        const resJ = await res.json()
        const list_id = resJ.id
        throw redirect(303, `/accounts/${locals.userid}/lists/${list_id}/`)
    }
}