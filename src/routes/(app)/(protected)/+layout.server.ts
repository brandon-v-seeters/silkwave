import { redirect, type ServerLoad } from "@sveltejs/kit";

export const load: ServerLoad = async ({ locals, url }) => {
    if (!locals.user) {
        const redirectTo = url.pathname + url.search;
        throw redirect(307, `/login?redirectTo=${encodeURIComponent(redirectTo)}`);
    }

    return {
        user: locals.user
    };
};