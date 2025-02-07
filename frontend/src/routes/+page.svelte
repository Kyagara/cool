<script>
  import { onMount } from "svelte";
  import { PUBLIC_URL } from "$env/static/public";

  import VirtualList from "svelte-tiny-virtual-list";

  import Filters from "$lib/components/Filters.svelte";
  import Header from "$lib/components/Header.svelte";
  import User from "$lib/components/User.svelte";

  import { sortMediaArray } from "$lib/sort.js";

  let users = $state([]);
  let sortedUsers = $state([]);
  let totalPosts = $state(0);
  let loading = $state(true);

  let filters = $state({
    search: "",
    sortBy: "posts",
    order: "desc",
    provider: "umate",
  });

  $effect(() => {
    sortedUsers = sortMediaArray(filters, users);
    loading = false;
  });

  onMount(() => {
    fetch(`${PUBLIC_URL}/api/users`)
      .then((res) => res.json())
      .then((data) => {
        users = data.users;
        totalPosts = data.totalPosts;
      });
  });
</script>

<svelte:head>
  <title>Users</title>
</svelte:head>

<Header page="Users" />

<div class="flex gap-4 items-center justify-center my-2">
  <div class="flex gap-2 text-xl">
    <p>
      {users.length} Users
    </p>
    <p>
      {totalPosts} Posts
    </p>
  </div>

  <Filters
    bind:filters
    showSearch={true}
    showProvider={true}
    showType={false}
  />
</div>

{#key sortedUsers}
  {#if !loading}
    {#if sortedUsers.length > 0}
      <VirtualList
        width="100%"
        height={670}
        itemCount={sortedUsers.length}
        itemSize={135}
      >
        <div slot="item" let:index let:style {style} id={index}>
          {#if sortedUsers[index]}
            <User user={sortedUsers[index]} provider={filters.provider} />
          {/if}
        </div>
      </VirtualList>
    {/if}
  {:else if loading}
    <p class="text-center">Loading...</p>
  {/if}
{/key}
