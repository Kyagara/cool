<script>
  import { PUBLIC_URL } from "$env/static/public";

  let { username, slug, filename, provider } = $props();
</script>

<!-- svelte-ignore a11y_media_has_caption -->
<video
  autoplay
  loop
  playsinline
  controls
  class="h-full object-contain will-change-transform translate-z-0"
>
  <source
    src={`${PUBLIC_URL}/provider/${provider}/${username}/${slug}/${filename}`}
    type="video/mp4"
  />
</video>
