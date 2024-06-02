<script lang="ts">
    import type { ToDoItemData } from "../ToDoItem/ToDoItem";
    import ToDoItem from "../ToDoItem/ToDoItem.svelte";
    var toDoItems: ToDoItemData[] = [
        {
            name: "Wash dishes",
            completed: false,
        },
        {
            name: "Do laundry",
            completed: false,
        },
        {
            name: "Eat food",
            completed: false,
        },
        {
            name: "Shit bricks",
            completed: false,
        },
        {
            name: "Mow lawn",
            completed: false,
        },
    ]
</script>

{#each toDoItems as item}
    <ToDoItem data={item} />
{/each}