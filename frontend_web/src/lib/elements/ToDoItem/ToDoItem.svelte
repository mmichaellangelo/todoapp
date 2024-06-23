<script lang="ts">
    import type { iTodo } from "$lib/types";
    export let data: iTodo;
    var isEditing = false;
</script>

<div class="todoitem_container">
    <input type="checkbox" bind:checked={data.completed}>
    <input type="text" bind:value={data.body} class="item_edit_input" on:blur={() => {isEditing = false}}>
</div>

<style>

    .todoitem_container {
        display: flex;
        flex-direction: row;
        max-width: 300px;
        align-items: center;
        margin-bottom: 0.5rem;
        background-color: var(--col-grayblue);
    }
    
</style>