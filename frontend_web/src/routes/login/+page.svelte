<script lang="ts">
    import { enhance } from '$app/forms';
    import { goto } from '$app/navigation';
    import { Session } from '$lib/state/userstore';
    import type { SubmitFunction } from '@sveltejs/kit';
    
    export let form;

    export let data

    let awaitingResponse = false;
    
    let loginError = {
        error: false,
        type: "",
    }

    const handleEnhance: SubmitFunction = () => {
        loginError.error = false;
        awaitingResponse = true;
        return async ({ result, update }) => {
            await update();
            awaitingResponse = false;
            if (result.type == 'success') {
                if (result.data?.success) {
                    const username = result.data?.username as string;
                    const userid = result.data?.userid as number;
                    Session.set({ username, userid });
                    setTimeout(() => {
                        goto("/dashboard")
                    }, 500)
                } else {
                    loginError.error = true;
                    loginError.type = "Incorrect username/password"
                }
            }
        }
    }

</script>

<div id="login_container">
    <h2>Login</h2>
    <form method="POST" use:enhance={handleEnhance}>
        <label>
            Username/Email:
            <input type="text" name="emailorusername" required>
        </label>
        <label>
            Password:
            <input type="password" name="password" required>
        </label>
        <div id="login_create_container">
            <a href="/create-account">Create an Account</a>
            <button type="submit">Log In</button>
        </div>
    </form>
</div>

{#if awaitingResponse}
    <p>Loading...</p>
{/if}

{#if form?.success}
    <p>Login success! Redirecting...</p>
{/if}

{#if loginError.error}
    <p>Error: {loginError.type}</p>
{/if}

<style>
    #login_container {
        display: flex;
        flex-direction: column;
    }

    #login_create_container a {
        margin-right: 1rem;
    }

    form {
        display: flex;
        flex-direction: column;
        width: fit-content;
        text-align: right;
    }

    label {
        margin-bottom: 0.5rem;
    }

    input {
        padding: 0.25rem;
    }

    button {
        padding: 0.25rem;
    }
</style>