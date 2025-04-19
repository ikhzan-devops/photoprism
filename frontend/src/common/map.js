let loading = false;
let maplibregl = null;

// Loads the maplibregl library.
export async function load() {
  if (maplibregl !== null || loading) {
    return Promise.resolve(maplibregl);
  }

  loading = true;

  try {
    const module = await import("./maplibregl.js");
    maplibregl = module.default;
    loading = false;
  } catch (e) {
    loading = false;
    console.error("maps: failed to load maplibregl", e);
    return Promise.reject(e);
  }

  return Promise.resolve(maplibregl);
}
