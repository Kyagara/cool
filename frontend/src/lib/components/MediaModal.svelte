<script>
  import Image from "$lib/components/Image.svelte";
  import Video from "$lib/components/Video.svelte";

  let { provider } = $props();

  let media = $state({ media: {} });
  let { status } = $state({ status: false });

  export function openModal(newMedia) {
    media = newMedia;
    status = true;
  }

  const handleKeyDown = (event) => {
    if (event.keyCode == 27) {
      status = false;
    }
  };
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if status}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-90"
    aria-modal="true"
    role="dialog"
  >
    <div
      class="flex flex-col items-center gap-2 w-full h-full max-h-screen overflow-hidden"
    >
      <a
        href={`/${provider}/${media.u}`}
        class="text-lg underline hover:text-blue-500 mt-2"
      >
        {media.u}
      </a>

      <div
        class="flex flex-grow justify-center items-center w-full max-h-full overflow-hidden"
      >
        {#if media.t === 0}
          <Image
            username={media.u}
            slug={media.s}
            filename={media.f}
            {provider}
          />
        {:else}
          <Video
            username={media.u}
            slug={media.s}
            filename={media.f}
            {provider}
          />
        {/if}
      </div>
    </div>
  </div>
{/if}
