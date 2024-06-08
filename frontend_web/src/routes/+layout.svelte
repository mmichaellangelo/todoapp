<script lang="ts">
    import Header from "$lib/elements/Header/Header.svelte";
    import "./global.css"
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