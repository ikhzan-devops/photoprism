<template>
  <v-dialog
    :model-value="visible"
    :max-width="900"
    :fullscreen="$vuetify.display.xs"
    persistent
    scrim
    scrollable
    class="p-location-dialog"
    @keydown.esc="close"
    @after-enter="afterEnter"
    @after-leave="afterLeave"
  >
    <v-card :tile="$vuetify.display.xs">
      <v-toolbar v-if="$vuetify.display.xs" flat color="navigation" class="mb-4" density="compact">
        <v-btn icon @click.stop="close">
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <v-toolbar-title>
          {{ $gettext("Adjust Location") }}
        </v-toolbar-title>
      </v-toolbar>
      <v-card-title v-else class="d-flex justify-start align-center ga-3">
        <v-icon size="28" color="primary">mdi-map-marker</v-icon>
        <h6 class="text-h6">{{ $gettext("Adjust Location") }}</h6>
      </v-card-title>
      <v-card-text class="pb-3">
        <div class="d-flex flex-column flex-md-row ga-5">
          <div class="flex-grow-1 position-relative mb-4 mb-md-0">
            <div ref="map" class="p-map" style="height: 45vh; min-height: 250px; width: 100%; border-radius: 4px"></div>
          </div>

          <div
            class="map-sidebar d-flex flex-column"
            :class="$vuetify.display.xs ? `ga-3` : 'ga-5'"
            :style="{
              width: $vuetify.display.smAndDown ? '100%' : '300px',
              maxWidth: $vuetify.display.smAndDown ? '100%' : '300px',
              minWidth: 0,
            }"
          >
            <div>
              <v-autocomplete
                ref="search"
                v-model="selectedPlace"
                :items="searchResults"
                :loading="searchLoading"
                :search="searchQuery"
                prepend-inner-icon="mdi-magnify"
                density="compact"
                variant="outlined"
                :placeholder="$gettext(`Search`)"
                item-title="name"
                item-value="id"
                return-object
                auto-select-first
                clearable
                autocomplete="off"
                no-filter
                menu-icon=""
                :menu-props="{ maxHeight: 300 }"
                @update:search="onSearchQueryChange"
                @update:model-value="onPlaceSelected"
                @click:clear="clearSearch"
              >
                <template #item="{ props }">
                  <v-list-item v-bind="props" density="compact">
                    <template #prepend>
                      <v-icon>mdi-map-marker</v-icon>
                    </template>
                  </v-list-item>
                </template>
                <template #no-data>
                  <v-list-item
                    v-if="searchQuery && searchQuery.length >= 2 && !searchLoading && searchResults.length === 0"
                  >
                    <v-list-item-title>{{ $gettext("No results") }}</v-list-item-title>
                  </v-list-item>
                </template>
              </v-autocomplete>
            </div>
            <!-- div v-if="locationInfo">
              <div class="text-subtitle-2 mb-2">{{ $gettext("Location Details") }}</div>
              <div class="text-body-2">
                {{ simplifiedLocationDisplay }}
              </div>
            </div -->

            <div class="text-body-2 mt-3">
              {{ $gettext("You can search for a location or move the marker on the map to change the position:") }}
            </div>

            <div class="flex-grow-1">
              <p-location-input
                :lat="currentLat"
                :lng="currentLng"
                density="comfortable"
                :enable-undo="true"
                :auto-apply="true"
                :label="locationLabel"
                @update:lat="setLat"
                @update:lng="setLng"
                @changed="onLocationChanged"
                @cleared="onLocationCleared"
              ></p-location-input>
            </div>

            <div class="action-buttons">
              <v-btn variant="flat" color="button" class="action-cancel" min-width="120" @click.stop="close">
                {{ $gettext("Cancel") }}
              </v-btn>
              <v-btn
                color="highlight"
                min-width="120"
                :disabled="!(currentLat !== null && currentLng !== null) || locationLoading"
                :loading="locationLoading"
                @click="confirm"
              >
                {{ $gettext("Confirm") }}
              </v-btn>
            </div>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import PLocationInput from "component/location/input.vue";
import * as map from "common/map";

let maplibregl = null;

export default {
  name: "PLocationDialog",
  components: {
    PLocationInput,
  },
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    latlng: {
      type: Array,
      default: () => [0, 0],
    },
    style: {
      type: String,
      default: "embedded",
    },
  },
  emits: ["update:lat", "update:lng", "close", "confirm"],
  data() {
    return {
      map: null,
      marker: null,
      position: [0.0, 0.0],
      options: {
        container: null,
        style: `https://cdn.photoprism.app/maps/${this.style}.json`,
        glyphs: `https://cdn.photoprism.app/maps/font/{fontstack}/{range}.pbf`,
        zoom: 12,
        interactive: true,
        attributionControl: false,
      },
      loaded: false,
      currentLat: this.lat,
      currentLng: this.lng,
      location: null,
      locationLoading: false,
      searchQuery: "",
      searchResults: [],
      searchLoading: false,
      searchTimeout: null,
      selectedPlace: null,
    };
  },
  computed: {
    locationLabel() {
      if (!this.location || !this.location?.place?.label) {
        return "";
      }

      return this.location.place.label;
    },
  },
  watch: {
    visible(show) {
      if (show) {
        this.currentLat = this.latlng[0];
        this.currentLng = this.latlng[1];
        this.setPosition(this.currentLat, this.currentLng);
      }
    },
    latlng(val) {
      this.currentLat = val[0];
      this.currentLng = val[1];
      this.setPosition(this.currentLat, this.currentLng);
    },
  },
  beforeUnmount() {
    this.afterLeave();
  },
  methods: {
    close() {
      this.$emit("close");
    },
    confirm() {
      if (this.currentLat !== null && this.currentLng !== null) {
        this.$emit("update:lat", this.currentLat);
        this.$emit("update:lng", this.currentLng);
        this.$emit("confirm", {
          lat: this.currentLat,
          lng: this.currentLng,
          location: this.location,
        });
      }
    },
    removeMap() {
      if (this.map) {
        this.map.remove();
        this.map = null;
        this.marker = null;
        this.loaded = false;
      }
    },
    afterEnter() {
      map.load().then((m) => {
        maplibregl = m;
        this.initMap();
      });
    },
    afterLeave() {
      this.removeMap();
      this.location = null;
      this.locationLoading = false;
      this.resetSearchState();
    },
    initMap() {
      if (this.map || !this.$refs.map) {
        return;
      }

      try {
        this.options.container = this.$refs.map;

        if (!this.currentLat || !this.currentLng || (this.currentLat === 0 && this.currentLng === 0)) {
          this.options.zoom = 2;
          this.options.center = [0, 20];
        } else {
          this.options.zoom = 12;
          this.options.center = [this.currentLng, this.currentLat];
        }

        this.map = new maplibregl.Map(this.options);

        this.map.on("styleimagemissing", (e) => {
          const emptyImage = new ImageData(1, 1);
          if (e && e.id) {
            this.map.addImage(e.id, emptyImage);
          }
        });

        this.map.addControl(
          new maplibregl.NavigationControl({
            showCompass: true,
            showZoom: true,
            visualizePitch: false,
          }),
          "top-right"
        );
        this.map.addControl(new maplibregl.ScaleControl({ maxWidth: 80, unit: "metric" }), "bottom-left");
        this.map.addControl(
          new maplibregl.GeolocateControl({
            positionOptions: {
              enableHighAccuracy: true,
            },
            trackUserLocation: true,
          }),
          "top-right"
        );

        this.map.on("error", (e) => {
          console.error("map:", e);
        });

        this.map.on("load", () => {
          this.loaded = true;
          if (this.currentLat && this.currentLng && !(this.currentLat === 0 && this.currentLng === 0)) {
            this.setPosition(this.currentLat, this.currentLng);
            this.fetchLocationInfo(this.currentLat, this.currentLng);
          }
          this.map.resize();
        });

        this.map.on("click", (e) => {
          this.setPositionAndFetchInfo(e.lngLat.lat, e.lngLat.lng);
        });
      } catch (error) {
        console.error("map: initialization failed", error);
        this.loaded = false;
      }
    },
    setPosition(lat, lng) {
      if (!this.map || !this.loaded) {
        return;
      }

      if (this.position[0] === lng && this.position[1] === lat && this.marker) {
        return;
      }

      // Skip invalid or empty coordinates (0,0)
      if ((lat === 0 && lng === 0) || !lat || !lng) {
        if (this.marker) {
          this.marker.remove();
          this.marker = null;
        }
        return;
      }
      this.position = [lng, lat];

      // Always center map when position changes
      this.map.setCenter(this.position, {
        zoom: 12,
        animate: false,
      });

      if (this.marker) {
        this.marker.setLngLat(this.position);
      } else {
        this.marker = new maplibregl.Marker({
          color: "#3fb4df",
          draggable: true,
        })
          .setLngLat(this.position)
          .addTo(this.map);

        // Update coordinates when marker is dragged
        this.marker.on("dragend", () => {
          const lngLat = this.marker.getLngLat();
          this.setPositionAndFetchInfo(lngLat.lat, lngLat.lng);
        });
      }
    },
    setLat(lat) {
      this.currentLat = lat;
      this.setPosition(lat, this.currentLng);
    },
    setLng(lng) {
      this.currentLng = lng;
      this.setPosition(this.currentLat, lng);
    },
    onLocationChanged(data) {
      if (data.lat !== 0 || data.lng !== 0) {
        this.fetchLocationInfo(data.lat, data.lng);
      }
    },
    onLocationCleared() {
      this.location = null;
      this.locationLoading = false;

      if (this.marker) {
        this.marker.remove();
        this.marker = null;
      }

      if (this.map) {
        this.map.flyTo({
          center: [0, 20],
          zoom: 2,
          essential: true,
        });
      }
    },
    clearSearchTimeout() {
      if (this.searchTimeout) {
        clearTimeout(this.searchTimeout);
        this.searchTimeout = null;
      }
    },
    resetSearchState() {
      this.searchQuery = "";
      this.searchResults = [];
      this.selectedPlace = null;
      this.searchLoading = false;
      this.clearSearchTimeout();
    },
    setPositionAndFetchInfo(lat, lng) {
      this.currentLat = lat;
      this.currentLng = lng;
      this.setPosition(lat, lng);
      this.fetchLocationInfo(lat, lng);
    },
    fetchLocationInfo(lat, lng) {
      this.locationLoading = true;
      this.$api
        .get(`places/reverse?lat=${lat}&lng=${lng}`)
        .then((response) => {
          if (response.data && response.data?.place?.label) {
            this.location = response.data;
          } else {
            this.location = null;
          }
        })
        .catch((error) => {
          console.error("Reverse geocoding error:", error);
          this.location = null;
        })
        .finally(() => {
          this.locationLoading = false;
        });
    },
    onSearchQueryChange(query) {
      this.searchQuery = query;
      this.clearSearchTimeout();

      if (!query || query.length < 2) {
        this.searchResults = [];
        this.searchLoading = false;
        return;
      }

      this.searchLoading = true;
      this.searchTimeout = setTimeout(() => {
        this.performPlaceSearch(query);
      }, 300); // 300ms delay after user stops typing
    },
    async performPlaceSearch(query) {
      if (!query || query.length < 2) {
        this.searchLoading = false;
        return;
      }

      try {
        const response = await this.$api.get("places/search", {
          params: {
            q: query,
            count: 10,
            locale: this.$config.getLanguageLocale() || "en",
          },
        });

        if (this.searchQuery === query) {
          if (response.data && Array.isArray(response.data)) {
            this.searchResults = response.data;
          } else {
            this.searchResults = [];
          }
        }
      } catch (error) {
        console.error("Place search error:", error);
        if (this.searchQuery === query) {
          this.searchResults = [];
        }
      } finally {
        if (this.searchQuery === query) {
          this.searchLoading = false;
        }
      }
    },
    onPlaceSelected(place) {
      if (place && place.lat && place.lng) {
        this.setPositionAndFetchInfo(place.lat, place.lng);

        this.$nextTick(() => {
          this.resetSearchState();
        });
      }
    },
    clearSearch() {
      this.resetSearchState();
    },
  },
};
</script>
