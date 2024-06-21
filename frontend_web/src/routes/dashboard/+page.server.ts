import type { iList } from "$lib/types";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({fetch}) => {
    const res = await fetch("http://api/accounts/1/lists/")
    const resJ = await res.json()

    const lists: iList[] = resJ as iList[];
    
    console.log("TYPEOF RESJ", typeof resJ)
    console.log(resJ)
    
     return {
        lists: lists,
     }
}