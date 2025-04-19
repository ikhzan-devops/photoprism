<template>
  <div ref="map" class="p-map"></div>
</template>

<script>
import * as map from "common/map";

let maplibregl = null;

export default {
  name: "PMap",
  props: {
    lat: {
      type: Number,
      required: true,
    },
    lng: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      map: null,
      marker: null,
      position: [0.0, 0.0],
      options: {
        container: null,
        // To test new styles, put the style file in /assets/static/geo
        // and include it from there e.g. "/static/geo/embedded.json".
        // Styles can be edited/created with https://maplibre.org/maputnik/.
        style: "https://cdn.photoprism.app/maps/embedded.json",
        glyphs: `https://cdn.photoprism.app/maps/font/{fontstack}/{range}.pbf`,
        zoom: 9,
        interactive: true,
        attributionControl: false,
      },
      loaded: false,
    };
  },
  watch: {
    lat() {
      this.updatePosition();
    },
    lng() {
      this.updatePosition();
    },
  },
  mounted() {
    map.load().then((m) => {
      maplibregl = m;
      this.initMap();
    });
  },
  beforeUnmount() {
    if (this.map) {
      this.map.remove();
    }
  },
  methods: {
    initMap() {
      if (this.map || !this.$refs.map || !maplibregl) {
        return;
      }

      try {
        this.options.container = this.$refs.map;
        this.map = new maplibregl.Map(this.options);

        // Add controls.
        /* this.map.addControl(
          new maplibregl.NavigationControl({
            showCompass: false,
            showZoom: true,
            visualizePitch: false,
          }),
          "top-right"
        );

        this.map.addControl(new maplibregl.ScaleControl({ maxWidth: 80, unit: "metric" }), "bottom-left");
        */

        this.map.on("error", (e) => {
          console.error("map:", e);
        });

        this.map.on("load", () => {
          this.loaded = true;
          this.updatePosition();
        });
      } catch (error) {
        console.error("map: initialization failed", error);
        this.loaded = false;
      }
    },
    updatePosition() {
      if (this.map && this.loaded) {
        if (this.position[0] === this.lng && this.position[1] === this.lat) {
          return;
        }

        this.position = [this.lng, this.lat];
        this.map.setCenter(this.position);

        if (this.marker) {
          this.marker.setLngLat(this.position);
        } else {
          this.marker = new maplibregl.Marker({
            color: "#3fb4df",
            draggable: false,
          })
            .setLngLat(this.position)
            .addTo(this.map);
        }
      }
    },
  },
};
</script>
