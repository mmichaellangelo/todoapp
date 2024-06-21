<script lang="ts">
    import type { iTodo } from "$lib/types";
    export let data: iTodo;
    var isEditing = false;
</script>

<div class="todoitem_container">
    <input type="checkbox" bind:checked={data.completed} tabindex="0">
    {#if isEditing}
        <!-- svelte-ignore a11y-autofocus 
         it just works so why change it? -->
        <!-- <input type="text" bind:value={data.body} class="item_edit_input" on:blur={() => {isEditing = false}} autofocus>  -->
         <div contenteditable="true" bind:innerText={data.body} class="item_edit_input" on:blur={() => {isEditing = false}} autofocus>{data.body}</div>
    {:else}
        <!-- svelte-ignore a11y-no-noninteractive-tabindex 
         again -- just works. why change it? -->
        <p class={data.completed ? "completed" : ""} tabindex="0" on:focus={() => {isEditing = true}}>{data.body}</p>
    {/if}
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

    .item_edit_input {
        width: fit-content;
        margin-top: 0.5rem;
        padding: 0.5rem;
        margin-bottom: 0.5rem;
        margin-left: -0.5rem;
    }
</style>