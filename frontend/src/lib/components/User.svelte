<script>
  import { PUBLIC_URL } from "$env/static/public";

  let { provider, user } = $props();

  let {
    u: username,
    d: displayName,
    a: avatar,
    b: banner,
    p: posts,
    l: links,
  } = user;
</script>

<div class="flex flex-col items-center justify-center">
  <div class="relative w-1/3 h-32 rounded-2xl shadow-lg overflow-hidden">
    <div class="absolute inset-0 bg-black bg-opacity-30"></div>

    <img
      src={`${PUBLIC_URL}/provider/${provider}/${username}/${banner}`}
      alt="User Banner"
      class="absolute w-full h-full object-cover opacity-80"
    />

    <div class="absolute inset-0 flex items-center p-4 space-x-4">
      <img
        src={`${PUBLIC_URL}/provider/${provider}/${username}/${avatar}`}
        alt={username}
        class="w-20 h-20 rounded-full border-2 border-white shadow-md object-cover"
      />

      <div class="flex flex-col gap-1">
        <div class="flex items-center gap-2">
          <a
            data-sveltekit-preload-data
            href="{provider}/{username}"
            class="text-xl font-semibold hover:text-blue-400"
          >
            {displayName ? displayName : username}
          </a>

          {#if displayName}
            <p class="text-sm opacity-80">{username}</p>
          {/if}
        </div>

        <p class="text-sm">
          Posts: <span class="font-medium">{posts}</span>
        </p>

        <div class="flex space-x-3">
          {#each links as link}
            <a
              href={link.u}
              target="_blank"
              class="text-sm hover:text-blue-400"
            >
              {link.w}
            </a>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>
