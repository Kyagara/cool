<script>
  import { PUBLIC_URL } from "$env/static/public";

  let { username, slug, filename, provider, classNames } = $props();

  let f = $state(filename);
  let gif = $state(false);
  if (filename.endsWith(".mp4")) {
    f = filename.replace(".mp4", ".gif");
    gif = true;
  }
</script>

<!-- svelte-ignore a11y_missing_attribute -->
<div
  class={`${gif ? "border-2 border-red-500 " : ""} flex justify-center items-center w-full h-full relative ${classNames}`}
>
  <img
    src={`${PUBLIC_URL}/provider/${provider}/${username}/${slug}/${f}`}
    class="max-w-full max-h-full h-auto object-contain"
  />
</div>
