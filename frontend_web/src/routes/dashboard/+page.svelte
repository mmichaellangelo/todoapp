<script lang="ts">
    import List from "$lib/elements/List/List.svelte";
    import type { iList } from "$lib/types";
    export let data: {
        lists: iList[]
    }
</script>

<h2>Dashboard</h2>

<h3>Lists</h3>
{#each data.lists as list}
    <List data={list}/>
{/each}