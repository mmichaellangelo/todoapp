import type { iList } from "$lib/types";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({fetch, locals, }) => {
   
   try {
      const res = await fetch(`http://api/accounts/${locals.userid}/lists/`)
      const resJ = await res.json()
      const lists: iList[] = resJ as iList[];
      return { success: true, lists: lists }
   } catch (error) {
      return { success: false }
   }
   
}