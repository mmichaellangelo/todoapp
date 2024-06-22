<script lang="ts">
    import List from "$lib/elements/List/List.svelte";

    export let data;
</script>
<button on:click={() => history.back()}>{"< Back"}</button>
<List data={data.list} />

<style>
    button {
        margin-bottom: 1rem;
    }
</style>