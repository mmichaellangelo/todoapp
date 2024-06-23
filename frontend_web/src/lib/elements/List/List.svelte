<script lang="ts">
    import type { iList } from "$lib/types";
    export let data: iList;
    import ToDoItem from "../ToDoItem/ToDoItem.svelte";

    var titleCache = data.title || "";

    function handleUpdateTitle(e: FocusEvent) {
        console.log("blur")
        e.relatedTarget
        const element = e.target as HTMLElement;
        const newTitle = element.innerText;
        if (!(newTitle == titleCache)) {
            
        }
    }

</script>

{#if data}
    <div class="list_container">
        <input type="text" name="title" class="title_input" bind:value={data.title} on:blur={handleUpdateTitle} placeholder="Title">
        <p>{data.description || "description"}</p>
        {#if data.todos}
            {#each data.todos as todo}
                <ToDoItem data={todo}/>
            {/each}
        {/if}
        <input type="text" placeholder="Add Todo"/>
        <button>Add Todo</button>
    </div>
{/if}

<style>
    .list_container {
        border: 2px solid white;
        padding: 0.5rem;
        margin-bottom: 0.5rem;
        background-color: var(--col-grayblue);
    }

    .title_input {
        font-size: medium;
    }

    p {
        opacity: 75%;
        font-size: smaller;
        margin-top: 0.5rem;
        text-indent: 0.5rem;
    }

    

</style>