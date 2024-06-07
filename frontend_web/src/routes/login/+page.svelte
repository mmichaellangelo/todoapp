<script lang="ts">
    import { enhance } from '$app/forms';
    import { goto } from '$app/navigation';
    import { AccessToken } from '$lib/state/userstore';
    import type { SubmitFunction } from '@sveltejs/kit';
    export let form;

    let awaitingResponse = false;

    const handleEnhance: SubmitFunction = () => {
        console.log("enhancement")
        awaitingResponse = true;
        return async ({ result, update }) => {
            await update()
            awaitingResponse = false;
            if (result.type == 'success') {
                const accesstoken = result.data?.accesstoken as string;
                console.log("ACCESS TOKEN:", accesstoken)
                AccessToken.set(accesstoken)
                setTimeout(() => {
                    goto("/")
                }, 500)
            } 
        }
    }

</script>

<div id="login_container">
    <h2>Login</h2>
    <form method="POST" use:enhance={handleEnhance}>
        <label>
            Username/Email:
            <input type="text" name="emailorusername">
        </label>
        <label>
            Password:
            <input type="password" name="password">
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