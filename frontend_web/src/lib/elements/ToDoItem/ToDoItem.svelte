<script lang="ts">
    import type { ToDoItemData } from "./ToDoItem";
    export let data: ToDoItemData;
</script>

<div class="todoitem_container">
    <input type="checkbox" bind:checked={data.completed}>
    <p class={data.completed ? "completed" : ""}>{data.name}</p>
</div>

<style>

    .todoitem_container {
        display: flex;
        flex-direction: row;
        border: 1px solid darkgray;
        border-radius: 0.5rem;
        max-width: 300px;
        align-items: center;
        margin-bottom: 0.5rem;
        background-color: var(--col-grayblue);
    }

    .completed {
        text-decoration: line-through;
    }

    input {
        margin: 1rem;
        width: 1rem;
        height: 1rem;
    }
</style>