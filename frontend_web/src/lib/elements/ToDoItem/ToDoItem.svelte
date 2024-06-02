<script lang="ts">
    import type { ToDoItemData } from "./ToDoItem";
    export let data: ToDoItemData;
</script>

<div class="todoitem_container">
    <input type="checkbox" bind:checked={data.completed}>
    <h3 class={data.completed ? "completed" : ""}>{data.name}</h3>
</div>

<style>

    .todoitem_container {
        display: flex;
        flex-direction: row;
    }

    .completed {
        text-decoration: line-through;
    }
</style>