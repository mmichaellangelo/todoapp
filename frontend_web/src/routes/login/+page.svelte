<script lang="ts">
    import { enhance } from '$app/forms';
    import { goto } from '$app/navigation';
    import { AccessToken, Username } from '$lib/state/userstore';
    import type { SubmitFunction } from '@sveltejs/kit';
    
    export let form;

    let awaitingResponse = false;
    let loginError = {
        error: false,
        type: "",
    }

    const handleEnhance: SubmitFunction = () => {
        loginError.error = false;
        awaitingResponse = true;
        return async ({ result, update }) => {
            await update()
            awaitingResponse = false;
            if (result.type == 'success') {
                if (result.data?.success) {
                    const accesstoken = result.data?.accesstoken as string;
                    AccessToken.set(accesstoken);
                    const username = result.data?.username as string;
                    Username.set(username);
                    setTimeout(() => {
                        goto("/")
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
            <button type="submit">Log In</button>
            <span>or <a href="/create-account">Create an Account</a></span>
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