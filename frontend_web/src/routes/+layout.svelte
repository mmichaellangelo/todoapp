<script lang="ts">
    import { setContext } from "svelte";
    import Header from "$lib/elements/Header/Header.svelte";
    import "./global.css"
    import { Session } from "$lib/state/userstore";
        
    export let data;
    
    $: Session.set(data.username ? { username: data.username, userid: data.userid } : undefined);
    

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