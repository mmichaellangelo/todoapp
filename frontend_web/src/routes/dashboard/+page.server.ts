import type { iList } from "$lib/types";
import { error, redirect, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ fetch, locals }) => {
      const res = await fetch(`http://api/accounts/${locals.userid}/lists/`)
      if (!res.ok) {
         if (res.status == 303) {
            throw redirect(303, "/login")
         }
         console.log(res.statusText)
         throw error(res.status, res.statusText)
      }
      try {
         const resJ = await res.json()
         const lists: iList[] = resJ as iList[];
         return { success: true, lists: lists }
      } catch (error) {
         console.log(error)
         return { success: false }
   }
}

      