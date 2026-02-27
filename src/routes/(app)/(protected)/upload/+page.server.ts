import { POST } from "$lib/api/Api";
import { redirect } from "@sveltejs/kit";
import { fail, setError, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { z } from "zod";

const schema = z.object({
    name: z.string().min(2, { message: 'Name is required' }),
});

export const load = async ({ locals }) => {
    const form = await superValidate(zod(schema));

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
        const form = await superValidate(request, zod(schema));
        const { name } = form.data;

        if (!form.valid) {
            return fail(400, { form });
        }

        try {
            await POST('/register/artist-name', fetch, { name })
        } catch (error) {
            return fail(400, { form });
        }

        return redirect(303, '/upload');
    }
}