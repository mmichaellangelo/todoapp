<script>
    import { Username } from "$lib/state/userstore";

</script>
<div id="header_container">
    <a href="/" id="title_link">
        <h1>ToDo</h1>
    </a>
    
    <div id="login_area_container">
        {#if $Username == ("" || undefined)}
            <a href="/login">Log in</a>
            <span>or</span>
            <a href="/create-account">Create an Account</a>
        {:else}
        <p>Logged in as {$Username}</p>
        <a href="/logout" data-sveltekit-preload-data="off">Log Out</a>
        {/if}
        
    </div>
</div>

<style>
    #header_container {
        display: flex;
        flex-direction: row;
        align-items: center;
        background-color: var(--col-grayblue);
        padding: 1rem;
    }

    #login_area_container {
        margin-left: auto;
        margin-right: 2rem;
        text-align: right;
    }

    a {
        color: white;
        text-decoration: underline;
    }

    #title_link {
        text-decoration: none;
    }

    #title_link:hover {
        color: white;
    }

</style>
