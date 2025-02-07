<script>
  let { page } = $props();
</script>

<div class="flex w-full items-center justify-center gap-2 my-2">
  {#if page === "Users"}
    <p>Users</p>
  {:else}
    <a href="/" class="underline" data-sveltekit-preload-data> Users </a>
  {/if}

  {#if page === "Gallery"}
    <p>Gallery</p>
  {:else}
    <a href="/gallery" class="underline" data-sveltekit-preload-data>
      Gallery
    </a>
  {/if}
</div>
