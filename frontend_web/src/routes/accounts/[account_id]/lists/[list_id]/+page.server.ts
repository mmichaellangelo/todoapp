import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import type { iList } from "$lib/types";

export const load: PageServerLoad = async ({fetch, params}) => {
    const listID = params.list_id;
    const accountID = params.account_id;
    const res = await fetch(`http://api/accounts/${accountID}/lists/${listID}/`)
    if (!res.ok) {
        throw error(res.status, res.statusText)
    }
    const resJ = await res.json()
    const list = resJ as iList;

    return { list: list }
    
}