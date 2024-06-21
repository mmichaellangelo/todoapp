<script lang="ts">
    import type { iList } from "$lib/types";

    export let data: iList[];
</script>

<h3>Lists</h3>
<ul>
    {#each data as list}
        <li>{list.title}</li>
    {/each}
</ul>