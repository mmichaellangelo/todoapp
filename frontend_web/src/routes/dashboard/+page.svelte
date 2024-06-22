<script lang="ts">
    import List from "$lib/elements/List/List.svelte";
    export let data;
</script>


<h3>Lists</h3>
{#if data}
{#if data.lists}
    {#each data.lists as list}
        <List data={list}/>
    {/each}
{/if}
{/if}