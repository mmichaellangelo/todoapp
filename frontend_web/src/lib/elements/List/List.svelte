<script lang="ts">
    import type { iList } from "$lib/types";
    export let data: iList
</script>

<ul>
    {#each data.items as item}
        <li>{item.title}</li>
    {/each}
</ul>