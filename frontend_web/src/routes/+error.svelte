<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";

    
    let imgsrc: string;
    let errormessage: string;
    
    onMount(() => {
        const images = ["fail1.gif", "fail2.gif", "fail3.gif"]
        const errormessages = ["oops", "aw shucks", "aaaaaaa", "oh no", "scheisse"]
        imgsrc = images[Math.floor(Math.random()*images.length)]
        errormessage = errormessages[Math.floor(Math.random()*errormessages.length)]
    })

    function isServerError(code: number): boolean {
        if (code >= 500 && code < 600) {
            return true
        } else {
            return false
        }
    }
    
 

</script>
<div id="error_container">
    <div id="text_container">
        <h1>{errormessage ? errormessage : "..."}</h1>
        <h2>Error {$page.status}</h2>
        <p>{$page.error?.message}</p>
        {#if (isServerError($page.status)) }
            <p>server says sorry</p>
        {/if}
    </div>

    <img src={imgsrc} alt="funny cat fail">
</div>



<style>
    #error_container {
        display: flex;
        align-items: center;
        justify-content: center;
        max-width: fit-content;
        padding: 1rem;
    }

    img {
        margin-left: 2rem;
        max-width: 200px;
        animation: fadeIn 5s;
    }

    h1 {
        animation: fadeIn 5s;
    }

    @keyframes fadeIn {
        0% { opacity: 0; }
        100% { opacity: 1; }
}
</style>