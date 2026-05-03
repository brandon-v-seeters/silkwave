import { zod } from 'sveltekit-superforms/adapters';
import { setError, superValidate } from 'sveltekit-superforms';
import { z } from 'zod';
import { redirect, setFlash } from 'sveltekit-flash-message/server';
import { POST } from '$lib/api/Api';
import { fail, isHttpError, type Actions } from '@sveltejs/kit';

export const ssr = false;

const changeEmailSchema = z.object({
    email: z.string().email().default(''),
    password: z.string().min(8).default('')
});

const changePasswordSchema = z.object({
    oldPassword: z.string().min(8).default(''),
    newPassword: z.string().min(8).default('')
});

const deleteAccountSchema = z.object({
    password: z.string().min(8).default('')
});

export async function load({ locals }) {
    const form = await superValidate(zod(changeEmailSchema));
    const passwordForm = await superValidate(zod(changePasswordSchema));
    const deleteForm = await superValidate(zod(deleteAccountSchema));

    if (locals.user) {
        form.data.email = locals.user.email;
    }

    return {
        form,
        passwordForm,
        deleteForm
    };
}

export const actions: Actions = {
    changeEmail: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, zod(changeEmailSchema));
        const { email, password } = form.data;

        if (!password) {
            setError(form, 'password', 'Password is required');
            return fail(400, { form });
        }

        try {
            await POST('/settings/email', fetch, { email, password });
            setFlash({ type: 'success', message: 'Your email has been changed.' }, cookies);
            return { form };
        } catch (error) {
            const msg = `Failed to change email${isHttpError(error) ? `: ${error.body.message}` : ''}`;
            setError(form, '', msg);
            return fail(400, { form });
        }
    },

    changePassword: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, zod(changePasswordSchema));
        const { oldPassword, newPassword } = form.data;

        if (!oldPassword) {
            setError(form, 'oldPassword', 'Old password is required');
            return fail(400, { form });
        }

        if (!newPassword) {
            setError(form, 'newPassword', 'New password is required');
            return fail(400, { form });
        }

        try {
            await POST('/settings/password', fetch, { oldPassword, newPassword });
            setFlash({ type: 'success', message: 'Your password has been changed.' }, cookies);
            return { form };
        } catch (error) {
            const msg = `Failed to change password${isHttpError(error) ? `: ${error.body.message}` : ''}`;
            setError(form, '', msg);
            return fail(400, { form });
        }
    },

    deleteAccount: async ({ request, fetch, cookies }) => {
        const deleteForm = await superValidate(request, zod(deleteAccountSchema));
        const { password } = deleteForm.data;

        if (!password) {
            setError(deleteForm, 'password', 'Password is required');
            return fail(400, { deleteForm });
        }

        try {
            await POST('/settings/delete-account', fetch, { password });

            cookies.delete('token', { path: '/' });

            redirect(303, '/', { type: 'success', message: 'Your account has been deleted.' }, cookies);  // Todo: Redirect the user to the discover page with all the songs that are 'sad'
        } catch {
            setError(deleteForm, 'password', 'Failed to delete account');
            return fail(400, { deleteForm });
        }
    }
};
