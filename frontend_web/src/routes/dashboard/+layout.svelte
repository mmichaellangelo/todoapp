<script lang="ts">

    import ListExplorer from "$lib/elements/ListExplorer/ListExplorer.svelte";
    export let data;
</script>
<title>Todo Overview</title>
<h2>Dashboard</h2>
{#if data.lists}
    <ListExplorer data={data.lists} />
{/if}
<slot/>