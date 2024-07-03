<script lang="ts">
    import ListExplorer from "$lib/elements/ListExplorer/ListExplorer.svelte";
    export let data;
    import List from "$lib/elements/List/List.svelte";
</script>
<title>Todo Overview</title>
<h2>Dashboard</h2>

<ListExplorer data={data.lists} />


<h3>Lists</h3>
{#if data}
{#if data.lists}
    {#each data.lists as list}
        <List data={list}/>
    {/each}
{/if}
{/if}