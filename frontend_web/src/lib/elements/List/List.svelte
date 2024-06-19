<script lang="ts">
    import type { iList } from "$lib/types";
    export let data: iList
    import ToDoItem from "../ToDoItem/ToDoItem.svelte";
</script>

<details>
    <summary>
        <h3>{data.title}</h3>
        <p>{data.description}</p>
    </summary>
    {#each data.items as item}
        <ToDoItem data={item}/>
    {/each}
</details>

<style>
    details {
        border: 2px solid white;
        padding: 0.5rem;
        border-radius: 0.2rem;
        margin-bottom: 0.5rem;
        background-color: var(--col-grayblue);
    }
    h3 {
        display: inline;
    }
</style>