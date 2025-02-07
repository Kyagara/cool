self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open("cool-cache").then(() => {
      console.log("Cache opened");
    })
  );
});

self.addEventListener("fetch", (event) => {
  const url = new URL(event.request.url);

  if (url.pathname.startsWith("/api")) {
    event.respondWith(
      caches.open("cool-cache").then(async (cache) => {
        const cachedResponse = await cache.match(event.request);
        if (cachedResponse) {
          console.log("Serving from cache:", event.request.url);
          return cachedResponse;
        }

        const networkResponse = await fetch(event.request);
        cache.put(event.request, networkResponse.clone());
        console.log("Fetched & Cached:", event.request.url);
        return networkResponse;
      })
    );
  }
});
