<script lang="ts">
    import Header from "$lib/elements/Header/Header.svelte";
    import "./global.css"
    import { Username } from "$lib/state/userstore";
    import { onMount } from "svelte";
    export let data;

    onMount(() => { 
        Username.set(data.username)
    })

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