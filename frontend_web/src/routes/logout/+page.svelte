<script lang="ts">
    import { goto } from '$app/navigation';
    import { Session } from '$lib/state/userstore';
    import { onMount } from 'svelte';

    onMount(async () => {
        Session.set(undefined)
        goto("/")
    })

</script>

<p>Logging out</p>