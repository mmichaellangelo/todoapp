<script lang="ts">

    import { enhance } from "$app/forms";
    import { goto } from "$app/navigation";
    import { Username } from "$lib/state/userstore";
    import { error, type SubmitFunction } from '@sveltejs/kit';

    let awaitingResponse = false;
    let errorMessage = "";
    let success = false;

    const handleEnhance: SubmitFunction = () => {
        awaitingResponse = true;
        errorMessage = "";
        return async ({ result, update }) => {
            await update();
            awaitingResponse = false;
            if (result.type == "success") {
                if (result.data?.success) {
                    success = true;
                    Username.set(result.data.username)
                    setTimeout(() => {
                        goto("/dashboard")
                    }, 1000)
                }
                if (!result.data?.success) {
                    errorMessage = "Error creating account";
                }
            }
        }
    } 

</script>

<div id="createaccount_container">
    <h2>Create Account</h2>
    <form method="POST" use:enhance={handleEnhance}>
        <label>
            Email:
            <input type="text" name="email">
        </label>
        <label>
            Username:
            <input type="text" name="username">
        </label>
        <label>
            Password:
            <input type="password" name="password">
        </label>
        <button type="submit">Create Account</button>
    </form>
</div>

{#if awaitingResponse}
    <p>Creating account...</p>
{/if}

{#if errorMessage != ""}
    <p style="color: red">{errorMessage}</p>
{/if}

{#if success}
    <p style="color: green">Account Created! Redirecting...</p>
{/if}

<style>
    #createaccount_container {
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
        width: fit-content;
        margin-left: auto;
    }
</style>