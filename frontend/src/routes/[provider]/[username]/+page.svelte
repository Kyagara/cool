<script>
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { PUBLIC_URL } from "$env/static/public";

  import Masonry from "$lib/components/Masonry.svelte";
  import Filters from "$lib/components/Filters.svelte";
  import Header from "$lib/components/Header.svelte";
  import PostModal from "$lib/components/PostModal.svelte";

  import { sortMediaArray } from "$lib/sort.js";

  let { provider, username } = $page.params;

  let filters = $state({
    type: "all",
    order: "desc",
    sortBy: "date",
  });

  let user = $state({
    displayName: "",
    bio: "",
    avatar: "",
    banner: "",
    totalImages: 0,
    totalVideos: 0,
    totalPosts: 0,
    links: [],
  });

  let posts = $state({});
  let media = $state([]);
  let sortedMedia = $state([]);

  let modal = $state({});

  let loading = $state(true);

  $effect(() => {
    sortedMedia = sortMediaArray(filters, media);
    loading = false;
  });

  onMount(() => {
    fetch(`${PUBLIC_URL}/api/users?provider=${provider}&username=${username}`)
      .then((res) => res.json())
      .then((newData) => {
        user = {
          displayName: newData.displayName,
          bio: newData.bio,
          avatar: newData.avatar,
          banner: newData.banner,
          totalImages: newData.totalImages,
          totalVideos: newData.totalVideos,
          totalPosts: newData.totalPosts,
          links: newData.links,
        };

        Object.entries(newData.posts).forEach(([slug, post]) => {
          let index = 0;
          post.m.forEach((postMedia) => {
            postMedia.u = username;
            postMedia.s = slug;
            postMedia.ca = post.ca;
            postMedia.l = post.l;
            postMedia.i = index++;
            media.push(postMedia);
          });
        });

        posts = newData.posts;
        sortedMedia = sortMediaArray(filters, media);
      });
  });
</script>

<svelte:head>
  <title>{`${provider[0].toUpperCase() + provider.slice(1)}/${username}`}</title
  >
</svelte:head>

<PostModal bind:this={modal} {posts} {provider} />

<Header page="" />

<div class="grid gap-4">
  <div class="w-full flex items-center justify-center">
    <div class="w-1/2 flex justify-center items-center">
      {#if user.banner}
        <img
          class="w-full"
          src={`${PUBLIC_URL}/provider/${provider}/${username}/${user.banner}`}
          alt="Banner"
        />
      {:else}
        <div class="w-full h-64 bg-black"></div>
      {/if}
    </div>
  </div>

  <div class="flex flex-col gap-4 items-center justify-center">
    <div class="flex flex-col gap-4 items-center justify-evenly">
      <div class="flex gap-4">
        {#if user.avatar}
          <img
            class="w-32 h-32 rounded-full"
            src={`${PUBLIC_URL}/provider/${provider}/${username}/${user.avatar}`}
            alt="Avatar"
          />
        {:else}
          <div class="w-32 h-32 rounded-full bg-black"></div>
        {/if}

        <div class="flex flex-col gap-2 items-center justify-center">
          <p class="text-3xl font-semibold">
            {user.displayName ? user.displayName : username}
            {#if user.displayName}
              <span class="text-sm opacity-80">{username}</span>
            {/if}
          </p>

          <div class="flex gap-2">
            <p class="text-sm">
              {user.totalPosts} Posts
            </p>

            <p class="text-sm">-</p>

            <span class="text-sm opacity-80">
              {`${user.totalImages} Images | ${user.totalVideos} Videos`}
            </span>
          </div>

          <div class="flex gap-2 items-center justify-center">
            {#each user.links as social}
              <a
                href={social.u}
                target="_blank"
                class="hover:underline text-semibold"
              >
                {social.w}
              </a>
            {/each}
          </div>
        </div>
      </div>

      {#if user.bio}
        <p class="text-xl m-4 w-1/2">
          {user.bio}
        </p>
      {/if}
    </div>

    <Filters bind:filters showType={true} />
  </div>

  {#key sortedMedia}
    {#if !loading}
      {#if sortedMedia}
        <Masonry media={sortedMedia} openModal={modal.openModal} {provider} />
      {/if}
    {:else if loading}
      <p class="text-center">Loading...</p>
    {/if}
  {/key}
</div>
