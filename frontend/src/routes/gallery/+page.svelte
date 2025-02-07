<script>
  import { onMount } from "svelte";
  import { PUBLIC_URL } from "$env/static/public";

  import Filters from "$lib/components/Filters.svelte";
  import PageHeader from "$lib/components/Header.svelte";
  import Masonry from "$lib/components/Masonry.svelte";
  import MediaModal from "$lib/components/MediaModal.svelte";

  import { sortMediaArray } from "$lib/sort.js";

  let sortedMedia = $state([]);
  let gallery = $state([]);

  let filters = $state({
    provider: "umate",
    type: "all",
    order: "desc",
    sortBy: "date",
  });

  let modal = $state({});
  let loading = $state(true);

  $effect(() => {
    sortedMedia = sortMediaArray(filters, gallery);
    loading = false;
  });

  onMount(() => {
    fetch(`${PUBLIC_URL}/api/gallery?provider=${filters.provider}`)
      .then((res) => res.json())
      .then((data) => {
        gallery = data;
        sortedMedia = sortMediaArray(filters, gallery);
      });
  });
</script>

<svelte:head>
  <title>Gallery</title>
</svelte:head>

<MediaModal bind:this={modal} provider={filters.provider} />

<div class="fixed flex flex-col bg-black w-full">
  <PageHeader page="Gallery" />

  <div class="flex gap-4 items-center justify-center my-2">
    <p class="text-xl">
      {#if !loading && sortedMedia && sortedMedia.length > 0}
        {sortedMedia.length}
        {filters.type === "all"
          ? "Media"
          : filters.type === "video"
            ? "Videos"
            : "Images"}
      {:else}
        Loading
      {/if}
    </p>

    <Filters bind:filters showType={true} showProvider={true} />
  </div>
</div>

<div class="pt-20"></div>

{#key sortedMedia}
  {#if !loading}
    {#if sortedMedia}
      <Masonry
        media={sortedMedia}
        openModal={modal.openModal}
        provider={filters.provider}
      />
    {/if}
  {:else if loading}
    <p class="text-center">Loading...</p>
  {/if}
{/key}
