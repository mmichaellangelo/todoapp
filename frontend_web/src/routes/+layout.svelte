<script lang="ts">
    import { setContext, writable } from "svelte";
    import Header from "$lib/elements/Header/Header.svelte";
    import "./global.css"
        
    export let data;

    const username = writable<string | undefined>(data.username)

    setContext("username", username)

    

</script>
<Header/>
<div>
    <slot/>
</div>


<style>
    div {
        padding: 1rem;
    }
</style>