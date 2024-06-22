<script lang="ts">
    import List from "$lib/elements/List/List.svelte";
    import {Session} from "$lib/state/userstore"
    import type { iList } from "$lib/types"
    var listData: iList = {
        account_id: $Session?.userid || -1,


    };
</script>
<h2>Create List</h2>

<List data={listData} />
