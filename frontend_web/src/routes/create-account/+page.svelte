
<div id="createaccount_container">
    <h2>Create Account</h2>
    <form method="POST">
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