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
        <h4 contenteditable="true" bind:innerText={data.title} on:blur={handleUpdateTitle}>Title</h4>
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
        border-radius: 0.2rem;
        margin-bottom: 0.5rem;
        background-color: var(--col-grayblue);
    }

    h4 {
        margin-top: 0.5rem;
        margin-bottom: 0.5rem;
        text-decoration: underline;
    }

    p {
        opacity: 75%;
        font-size: smaller;
        margin-top: 0.5rem;
        text-indent: 0.5rem;
    }
</style>