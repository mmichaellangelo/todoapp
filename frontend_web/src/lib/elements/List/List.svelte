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
        <form method="POST">
            <input type="text" name="title" class="title_input" bind:value={data.title} on:blur={handleUpdateTitle} placeholder="Title">
            <button class="delete_button" formaction={`/accounts/${data.account_id}/lists/${data.id}?/delete`} type="submit">Delete</button>
            <p>{data.description || "description"}</p>
            {#if data.todos}
                {#each data.todos as todo}
                    <ToDoItem data={todo}/>
                {/each}
            {/if}
            <input type="text" placeholder="Add Todo"/>
            <button type="submit" formaction="?/add">Add Todo</button>
        </form>
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
        margin-right: 2rem;
    }

    p {
        opacity: 75%;
        font-size: smaller;
        margin-top: 0.5rem;
        text-indent: 0.5rem;
    }

    .delete_button {
        
    }
    

</style>