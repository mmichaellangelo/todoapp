import type { iList } from "$lib/types";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = () => {
     return {
        lists: [
            {
                title: "a list",
                description: "baba",
                items: [
                    {title: "item 1"},
                    {title: "item 2"},
                ]
            },
            {
                title: "another list",
                description: "baba baba",
                items: [
                    {title: "item uno"},
                    {title: "item dos"},
                ]
            },
            {
                title: "another another list",
                description: "baba baba baba",
                items: [
                    {title: "item juan"},
                    {title: "item too"},
                ]
            },
        ] satisfies iList[]
     }
}