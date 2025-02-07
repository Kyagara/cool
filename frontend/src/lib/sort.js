export function sortMediaArray(filters, data) {
  if (data && data.length === 0) {
    return data;
  }

  let temp = [...data];
  return sort(filters, temp);
}

function sort(filters, data) {
  if (!data || data.length === 0) return data;

  if (filters.search?.trim()) {
    const searchTerm = filters.search.toLowerCase();
    data = data.filter(user => user.u.toLowerCase().includes(searchTerm));
  }

  if (data.length === 0) return data;

  if (filters.type && filters.type !== "all") {
    const typeFilter = filters.type === "video" ? 1 : 0;
    data = data.filter(item => item.t === typeFilter);
  }

  if (data.length === 0) return data;

  if (filters.sortBy) {
    const order = filters.order === "desc" ? -1 : 1;

    data.sort((a, b) => {
      switch (filters.sortBy) {
        case "posts":
          return (a.p - b.p) * order;

        case "date":
          return (new Date(a.ca).getTime() - new Date(b.ca).getTime()) * order;

        case "likes":
          return (a.l - b.l) * order;

        case "name":
          return a.u.localeCompare(b.u, undefined, { sensitivity: "base" }) * order;

        default:
          return 0;
      }
    });
  }

  return data;
}
