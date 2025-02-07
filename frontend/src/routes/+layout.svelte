<script>
  import "../app.css";

  let { children } = $props();

  if (import.meta.env.MODE === "production") {
    if ("serviceWorker" in navigator) {
      navigator.serviceWorker
        .register("/service-worker.js")
        .then(() => {
          console.log("Service Worker registered successfully");
        })
        .catch((error) => {
          console.error("Service Worker registration failed:", error);
        });
    }
  }
</script>

{@render children()}
