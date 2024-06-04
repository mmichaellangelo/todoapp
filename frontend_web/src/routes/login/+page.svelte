
<div id="login_container">
    <h2>Login</h2>
    <form method="POST">
        <label>
            Username/Email:
            <input type="text">
        </label>
        <label>
            Password:
            <input type="password">
        </label>
        <div id="login_create_container">
            <button type="submit">Log In</button>
            <span>or <a href="#">Create an Account</a></span>
        </div>
    </form>
</div>

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