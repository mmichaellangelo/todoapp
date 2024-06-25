import type { iList } from "$lib/types";
import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({fetch, locals }) => {
   console.log("Layout server load!")
      const res = await fetch(`http://api/accounts/${locals.userid}/lists/`)
      if (!res.ok) {
         console.log(res.statusText)
         throw error(res.status, res.statusText)
      }
      try {
         const resJ = await res.json()
         console.log(resJ)
         const lists: iList[] = resJ as iList[];
         return { success: true, lists: lists }
      } catch (error) {
         console.log(error)
         return { success: false }
   }
}