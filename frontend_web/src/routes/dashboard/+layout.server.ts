import type { iList } from "$lib/types";
import { error, redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";
import { goto } from "$app/navigation";

export const load: LayoutServerLoad = async ({ fetch, locals }) => {
   console.log("Layout server load!")
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
         console.log(resJ)
         const lists: iList[] = resJ as iList[];
         return { success: true, lists: lists }
      } catch (error) {
         console.log(error)
         return { success: false }
   }
}