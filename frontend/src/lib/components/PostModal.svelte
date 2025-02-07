<script>
  import Image from "$lib/components/Image.svelte";
  import Video from "$lib/components/Video.svelte";

  let { posts, provider } = $props();

  let { post } = $state({ post: {} });
  let { status } = $state({ status: false });

  export function openModal(newPost) {
    post = posts[newPost.s];
    index = newPost.i;
    status = true;
  }

  let { index } = $state({ index: 0 });

  const nextItem = () => {
    const newValue = (index + 1) % post.m.length;
    index = newValue;
  };

  const prevItem = () => {
    const newValue = (index - 1 + post.m.length) % post.m.length;
    index = newValue;
  };

  function handleKeyDown(event) {
    if (event.key === "Escape") {
      closeModal();
    }
  }

  function closeModal() {
    status = false;
    index = 0;
  }
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if status}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-90"
    aria-modal="true"
    role="dialog"
  >
    <div class="max-w-screen flex h-full w-full">
      {#if post.m.length > 1}
        <button
          onclick={prevItem}
          class="w-16 bg-gray-800 hover:text-blue-500"
          aria-label="Previous item"
        >
          ←
        </button>
      {/if}

      <div class="relative flex flex-1 items-center justify-center">
        {#key post.m[index].f}
          {#if post.m[index].t === 0}
            <Image
              username={post.m[index].u}
              slug={post.m[index].s}
              filename={post.m[index].f}
              {provider}
            />
          {:else if post.m[index].t === 1}
            <Video
              username={post.m[index].u}
              slug={post.m[index].s}
              filename={post.m[index].f}
              {provider}
            />
          {/if}
        {/key}
      </div>

      {#if post.m.length > 1}
        <button
          onclick={nextItem}
          class="w-16 bg-gray-800 hover:text-blue-500"
          aria-label="Next item"
        >
          →
        </button>
      {/if}

      <div class="flex h-full w-1/5 flex-col bg-white text-black">
        <p class="mt-4 text-center">
          {index + 1} / {post.m.length}
        </p>

        <div class="flex-1 overflow-y-auto p-4">
          {#if post.c}
            <p
              class="text-lg font-medium text-gray-900 text-pretty break-words"
            >
              {post.c}
            </p>
          {/if}

          <div class="mt-4 flex flex-col gap-1 text-gray-500">
            <span>{new Date(post.ca).toDateString()}</span>
            <span>{post.l} likes</span>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
