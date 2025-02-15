<script>
  let { filters, showSearch, showType, showProvider } = $props();
</script>

<div class="flex gap-2">
  {#if showSearch}
    <input
      bind:value={filters.search}
      name="search"
      type="text"
      placeholder="Search"
      class="border-2 p-1 rounded-md bg-black"
    />
  {/if}

  {#if showProvider}
    <select
      bind:value={filters.provider}
      class="border-2 rounded-md bg-black"
      name="provider"
    >
      <option value="umate">Umate</option>
    </select>
  {/if}

  {#if showType}
    <select
      bind:value={filters.type}
      name="type"
      class="border-2 rounded-md bg-black"
    >
      <option value="all">Show all media</option>
      <option value="video">Videos</option>
      <option value="image">Images</option>
    </select>
  {/if}

  {#if showType}
    <select
      bind:value={filters.sortBy}
      name="sortBy"
      class="border-2 rounded-md bg-black"
    >
      <option value="date">Date</option>
      <option value="likes">Likes</option>
    </select>
  {:else}
    <select
      bind:value={filters.sortBy}
      name="sortBy"
      class="border-2 rounded-md bg-black"
    >
      <option value="posts">Posts</option>
      <option value="name">Name</option>
    </select>
  {/if}

  <select
    bind:value={filters.order}
    name="order"
    class="border-2 rounded-md bg-black"
  >
    <option value="desc">Descending</option>
    <option value="asc">Ascending</option>
  </select>
</div>
