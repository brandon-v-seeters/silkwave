import { POST } from "$lib/api/Api";
import { redirect } from "@sveltejs/kit";
import { fail, setError, superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { z } from "zod";

const schema = z.object({
    name: z.string().min(2, { message: 'Artist name must be at least 2 characters.' }),
});

export const load = async ({ locals }) => {
    const form = await superValidate(zod4(schema));

    const { user } = locals;

    if (!user) {
        throw redirect(307, '/login');
    }

    return {
        form,
        user
    };
};


export const actions = {
    default: async ({ request, fetch }) => {
        const form = await superValidate(request, zod4(schema));
        const { name } = form.data;

        if (!form.valid) {
            return fail(400, { form });
        }

        try {
            await POST('/register/artist-name', fetch, { name });
        } catch (error) {
            const message =
                (error as { body?: { message?: string } })?.body?.message ??
                'Could not set your artist name. Please try again.';
            return setError(form, 'name', message);
        }

        return redirect(303, '/upload');
    }
}