<script>
  import { onMount } from "svelte";

  import Image from "./Image.svelte";
  import Video from "./Video.svelte";

  let { media, openModal, provider } = $props();

  let columns = $state([]);
  let columnCount = 5;
  let visibleItems = 0;
  let loadedItemsCount = 20;

  const arrangeItems = () => {
    columns = Array.from({ length: columnCount }, () => []);

    let index = 0;
    visibleItems.forEach((item) => {
      columns[index % columnCount].push(item);
      index++;
    });
  };

  const handleScroll = () => {
    const scrollPosition = window.scrollY + window.innerHeight;
    const documentHeight = document.documentElement.scrollHeight;
    const nearBottom = scrollPosition >= documentHeight * 0.8;

    if (nearBottom) {
      const newItems = media.slice(
        visibleItems.length,
        visibleItems.length + 20,
      );
      visibleItems = [...visibleItems, ...newItems];
      arrangeItems();
    }
  };

  onMount(() => {
    visibleItems = media.slice(0, loadedItemsCount);
    arrangeItems();

    window.addEventListener("scroll", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  });
</script>

<div class="w-full overflow-y-auto h-full">
  <div class="flex gap-1.5">
    {#each columns as column}
      <div class="flex flex-1 flex-col gap-1.5">
        {#each column as item}
          <button
            class="w-full overflow-hidden rounded-md shadow-md"
            onclick={() => openModal(item)}
          >
            <Image
              classNames="transition duration-100 ease-in-out hover:-translate-y-1 hover:scale-150 overflow-hidden"
              username={item.u}
              slug={item.s}
              filename={item.f}
              {provider}
            />
          </button>
        {/each}
      </div>
    {/each}
  </div>
</div>
