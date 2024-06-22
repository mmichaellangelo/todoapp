<script lang="ts">
    import type { iList } from "$lib/types";

    export let data: iList[];
</script>

<div id="explorer_container">
    <h3>Lists</h3>
    <ul>
        {#each data as list}
            <li>
                <a href={`/accounts/${list.account_id}/lists/${list.id}`}>{list.title}</a>
            </li>
            
        {/each}
    </ul>
</div>

<style>
    #explorer_container {
        border: 2px solid white;
        padding: 1rem;
    }

    ul {
        list-style-type: none;
        padding: 0px;
    }

    li {
        margin-top: 0.5rem;
        text-indent: 0.5rem;
    }
</style>