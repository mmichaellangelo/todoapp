<script lang="ts">
    import { setContext } from "svelte";
    import Header from "$lib/elements/Header/Header.svelte";
    import "./global.css"
    import { Username } from "$lib/state/userstore";
        
    export let data;

    Username.set(data.username);

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