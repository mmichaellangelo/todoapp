<script lang="ts">
    import { goto } from '$app/navigation';
    import { Username } from '$lib/state/userstore';
    import { onMount } from 'svelte';

    onMount(async () => {
        Username.set(undefined)
        goto("/")
    })

</script>

<p>Logging out</p>