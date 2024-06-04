// example load func from docs

import type { Actions } from "@sveltejs/kit";

// export const load: PageServerLoad = async ({ cookies }) => {
// 	const user = await db.getUserFromSession(cookies.get('sessionid'));
// 	return { user };
// };

export const actions = {
    default: async ({cookies, request}) => {
        const data = await request.formData();

    }
} satisfies Actions;