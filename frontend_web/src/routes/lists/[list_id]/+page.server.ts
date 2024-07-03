import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import type { iList } from "$lib/types";

export const load: PageServerLoad = async ({fetch, params}) => {
    const res = await fetch(`http://api/lists/${params.list_id}`)
    if (!res.ok) {
        const errmsg = await res.text()
        throw error(res.status, errmsg)
    }
    const resJ = await res.json()
    const list = resJ as iList;

    return { list: list }
    
}

export const actions = {
    delete: async ({params, fetch}) => {
        console.log("DELETE")
        const res = await fetch(`http://api/accounts/${params.account_id}/lists/${params.list_id}`,
            {
                method: "DELETE"  
            }
        )
    }
}