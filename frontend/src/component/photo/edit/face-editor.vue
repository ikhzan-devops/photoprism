<template>
  <div class="photo-face-editor">
    <div v-if="primaryFile" class="photo-preview-container mb-4">
      <div
        ref="photoPreviewWrapper"
        class="photo-preview-wrapper"
        :style="{ cursor: wrapperCursor }"
        @mousedown="handleWrapperMouseDown"
      >
        <img ref="photoPreview" :src="photoUrl" class="photo-preview" @load="onPhotoLoaded" />
        <div
          v-for="marker in markers"
          :key="marker.UID"
          :ref="(el) => (markerRefs[marker.UID] = el)"
          class="face-marker"
          :class="{
            'face-marker-selected': selectedMarker?.UID === marker.UID,
            'face-marker-editing': isEditingMarker && selectedMarker?.UID === marker.UID,
          }"
          :style="getMarkerStyle(marker)"
          @mousedown.stop="handleMarkerMouseDown($event, marker)"
          @mouseenter="hoveredMarkerUID = marker.UID"
          @mouseleave="hoveredMarkerUID = null"
        >
          <div v-if="selectedMarker?.UID === marker.UID || hoveredMarkerUID === marker.UID" class="face-marker-name">
            {{ marker.Name || $gettext("Unnamed") }}
            <div class="face-marker-actions">
              <v-btn
                v-if="!marker.SubjUID && !marker.Invalid"
                size="x-small"
                variant="text"
                icon
                color="error"
                density="compact"
                @click.stop="onReject(marker)"
              >
                <v-icon size="x-small">mdi-close</v-icon>
              </v-btn>
            </div>
          </div>
          <template v-if="selectedMarker?.UID === marker.UID && isEditingMarker && !interaction.active">
            <div
              class="marker-handle handle-nw"
              @mousedown.stop="handleResizeHandleMouseDown($event, marker, 'nw')"
            ></div>
            <div
              class="marker-handle handle-ne"
              @mousedown.stop="handleResizeHandleMouseDown($event, marker, 'ne')"
            ></div>
            <div
              class="marker-handle handle-sw"
              @mousedown.stop="handleResizeHandleMouseDown($event, marker, 'sw')"
            ></div>
            <div
              class="marker-handle handle-se"
              @mousedown.stop="handleResizeHandleMouseDown($event, marker, 'se')"
            ></div>
          </template>
        </div>
        <div
          v-if="interaction.drawingPreview"
          class="face-marker face-marker-new-preview"
          :style="getPreviewMarkerStyle()"
        ></div>
      </div>
      <div class="photo-actions mt-2 d-flex justify-space-between">
        <div>
          <v-btn
            :color="isDrawingMode ? 'primary' : 'default'"
            variant="outlined"
            prepend-icon="mdi-plus"
            class="mr-2"
            @click="toggleDrawingMode"
          >
            {{ isDrawingMode ? $gettext("Cancel") : $gettext("Add Face Marker") }}
          </v-btn>
          <v-btn v-if="isEditingMarker" color="primary" variant="outlined" @click="saveMarkerChanges">
            {{ $gettext("Save Changes") }}
          </v-btn>
        </div>
        <div>
          <v-btn color="success" variant="outlined" prepend-icon="mdi-check" @click="$emit('close')">
            {{ $gettext("Done") }}
          </v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Marker from "model/marker";
import $api from "common/api";
import "../../../css/face-markers.css";

export default {
  name: "PPhotoFaceEditor",
  props: {
    uid: {
      type: String,
      default: "",
    },
    primaryFile: {
      type: Object,
      default: null,
    },
    initialMarkers: {
      type: Array,
      default: () => [],
    },
  },
  emits: ["close", "markers-updated"],
  data() {
    return {
      markers: [],
      markerRefs: {},
      busy: false,
      disabled: !this.$config.feature("edit"),
      readonly: this.$config.get("readonly"),

      mode: null,
      selectedMarker: null,
      hoveredMarkerUID: null,

      interaction: {
        active: false,
        type: null,
        startX: 0,
        startY: 0,
        wrapperRect: null,
        initialMarkerRect: null,
        resizeHandle: null,
        drawingPreview: null,
      },

      photoAspectRatio: 1,
      photoLoaded: false,
    };
  },
  computed: {
    photoUrl() {
      if (!this.primaryFile) return "";
      return `${this.$config.contentUri}/t/${this.primaryFile.Hash}/${this.$config.previewToken}/fit_1280`;
    },

    isDrawingMode() {
      return this.mode === "drawing";
    },
    isEditingMarker() {
      return this.selectedMarker && this.mode === "editing";
    },
    canResize() {
      return this.isEditingMarker && !this.interaction.active;
    },
    wrapperCursor() {
      if (this.isDrawingMode) return "crosshair";
      if (this.interaction.active && this.interaction.type === "moving") return "grabbing";
      if (this.interaction.active && this.interaction.type === "resizing") {
        const handle = this.interaction.resizeHandle;
        if (handle === "nw" || handle === "se") return "nwse-resize";
        if (handle === "ne" || handle === "sw") return "nesw-resize";
      }
      return "default";
    },
  },
  watch: {
    initialMarkers: {
      handler(newMarkers) {
        this.markers = newMarkers.map((markerData) => new Marker(markerData));
      },
      immediate: true,
    },
  },
  mounted() {
    window.addEventListener("resize", this.updateWrapperRect);
    if (this.$refs.photoPreviewWrapper) {
      this.updateWrapperRect();
    }
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.updateWrapperRect);
    this.removeDocumentListeners();
  },
  methods: {
    updateWrapperRect() {
      if (this.$refs.photoPreviewWrapper) {
        const rect = this.$refs.photoPreviewWrapper.getBoundingClientRect();
        if (this.interaction.active) {
          this.interaction.wrapperRect = rect;
        }
      }
    },
    onPhotoLoaded(event) {
      const img = event.target;
      this.photoAspectRatio = img.naturalWidth / img.naturalHeight;
      this.photoLoaded = true;
      this.updateWrapperRect();
    },
    getMarkerStyle(marker) {
      return {
        left: marker.X * 100 + "%",
        top: marker.Y * 100 + "%",
        width: marker.W * 100 + "%",
        height: marker.H * 100 + "%",
      };
    },
    getPreviewMarkerStyle() {
      if (!this.interaction.drawingPreview) return {};
      const preview = this.interaction.drawingPreview;
      return {
        left: preview.x + "%",
        top: preview.y + "%",
        width: preview.w + "%",
        height: preview.h + "%",
      };
    },
    selectAndEditMarker(marker) {
      this.selectedMarker = marker;
      this.mode = "editing";
    },
    addDocumentListeners() {
      document.addEventListener("mousemove", this.handleDocumentMouseMove);
      document.addEventListener("mouseup", this.handleDocumentMouseUp);
    },
    removeDocumentListeners() {
      document.removeEventListener("mousemove", this.handleDocumentMouseMove);
      document.removeEventListener("mouseup", this.handleDocumentMouseUp);
    },
    handleWrapperMouseDown(event) {
      if (this.isDrawingMode && event.button === 0) {
        this.startInteraction("drawing");
        this.setDrawingStart(event);
      }
    },
    handleMarkerMouseDown(event, marker) {
      if (event.button === 0) {
        this.selectAndEditMarker(marker);
        this.startInteraction("moving", marker);
        this.setInteractionStart(event);
      }
    },
    handleResizeHandleMouseDown(event, marker, handle) {
      if (event.button === 0) {
        this.selectAndEditMarker(marker);
        this.startInteraction("resizing", marker, handle);
        this.setInteractionStart(event);
      }
    },
    handleDocumentMouseMove(event) {
      if (!this.interaction.active) return;
      window.requestAnimationFrame(() => {
        const { type } = this.interaction;
        if (type === "drawing") this.handleDrawing(event);
        else if (type === "moving") this.handleMoving(event);
        else if (type === "resizing") this.handleResizing(event);
      });
    },
    handleDocumentMouseUp() {
      if (!this.interaction.active) return;
      const { type, drawingPreview } = this.interaction;
      if (type === "drawing" && drawingPreview) {
        if (drawingPreview.w > 0.01 && drawingPreview.h > 0.01) {
          this.createMarker(
            drawingPreview.x / 100,
            drawingPreview.y / 100,
            drawingPreview.w / 100,
            drawingPreview.h / 100
          );
        }
      }
      this.interaction.active = false;
      this.removeDocumentListeners();
    },
    setDrawingStart(event) {
      const rect = this.interaction.wrapperRect;
      if (!rect) return;
      this.interaction.startX = (event.clientX - rect.left) / rect.width;
      this.interaction.startY = (event.clientY - rect.top) / rect.height;
    },
    setInteractionStart(event) {
      this.interaction.startX = event.clientX;
      this.interaction.startY = event.clientY;
    },
    handleDrawing(event) {
      const rect = this.interaction.wrapperRect;
      if (!rect) return;
      const currentX = Math.max(0, Math.min(1, (event.clientX - rect.left) / rect.width));
      const currentY = Math.max(0, Math.min(1, (event.clientY - rect.top) / rect.height));
      const startX = this.interaction.startX;
      const startY = this.interaction.startY;
      this.interaction.drawingPreview = {
        x: Math.min(startX, currentX) * 100,
        y: Math.min(startY, currentY) * 100,
        w: Math.abs(currentX - startX) * 100,
        h: Math.abs(currentY - startY) * 100,
      };
    },
    handleMoving(event) {
      if (!this.selectedMarker || !this.interaction.wrapperRect) return;
      const rect = this.interaction.wrapperRect;
      const deltaX = (event.clientX - this.interaction.startX) / rect.width;
      const deltaY = (event.clientY - this.interaction.startY) / rect.height;
      const initial = this.interaction.initialMarkerRect;
      let newX = initial.X + deltaX;
      let newY = initial.Y + deltaY;
      newX = Math.max(0, Math.min(1 - initial.W, newX));
      newY = Math.max(0, Math.min(1 - initial.H, newY));
      this.selectedMarker.X = newX;
      this.selectedMarker.Y = newY;
      this.updateMarkerElement(this.selectedMarker);
    },
    handleResizing(event) {
      if (!this.selectedMarker || !this.interaction.wrapperRect) return;
      const rect = this.interaction.wrapperRect;
      const mouseX = Math.max(0, Math.min(rect.width, event.clientX - rect.left));
      const mouseY = Math.max(0, Math.min(rect.height, event.clientY - rect.top));
      const initial = this.interaction.initialMarkerRect;
      const handle = this.interaction.resizeHandle;
      let newX = initial.X,
        newY = initial.Y,
        newW = initial.W,
        newH = initial.H;
      if (handle.includes("w")) {
        newX = mouseX / rect.width;
        newW = initial.X + initial.W - newX;
      }
      if (handle.includes("e")) {
        newW = (mouseX - initial.X * rect.width) / rect.width;
      }
      if (handle.includes("n")) {
        newY = mouseY / rect.height;
        newH = initial.Y + initial.H - newY;
      }
      if (handle.includes("s")) {
        newH = (mouseY - initial.Y * rect.height) / rect.height;
      }
      newW = Math.max(0.01, newW);
      newH = Math.max(0.01, newH);
      newX = Math.max(0, Math.min(1 - newW, newX));
      newY = Math.max(0, Math.min(1 - newH, newY));
      this.selectedMarker.X = newX;
      this.selectedMarker.Y = newY;
      this.selectedMarker.W = newW;
      this.selectedMarker.H = newH;
      this.updateMarkerElement(this.selectedMarker);
    },
    updateMarkerElement(marker) {
      const markerRef = this.markerRefs[marker.UID];
      if (markerRef) {
        markerRef.style.left = `${marker.X * 100}%`;
        markerRef.style.top = `${marker.Y * 100}%`;
        markerRef.style.width = `${marker.W * 100}%`;
        markerRef.style.height = `${marker.H * 100}%`;
      }
    },
    createMarker(x, y, w, h) {
      if (!this.primaryFile) return;
      this.busy = true;
      this.$notify.blockUI("busy");
      const markerData = {
        FileUID: this.primaryFile.UID,
        Type: "face",
        Src: "manual",
        X: x,
        Y: y,
        W: w,
        H: h,
        MarkerName: "",
        MarkerReview: false,
        MarkerInvalid: false,
      };
      $api
        .post("markers", markerData)
        .then((response) => {
          const newMarker = new Marker(response.data);
          this.markers.push(newMarker);
          this.$nextTick(() => {
            this.selectAndEditMarker(newMarker);
          });
          this.$notify.success(this.$gettext("Face marker added"));
          this.$emit("markers-updated", this.markers);
        })
        .catch((error) => {
          console.error("Error creating marker:", error);
          this.$notify.error(this.$gettext("Failed to add face marker"));
        })
        .finally(() => {
          this.busy = false;
          this.$notify.unblockUI();
          this.mode = null;
        });
    },
    updateMarkerPosition(marker) {
      if (!marker || this.busy) return;
      const original = this.interaction.initialMarkerRect;
      if (
        original &&
        marker.X === original.X &&
        marker.Y === original.Y &&
        marker.W === original.W &&
        marker.H === original.H
      ) {
        this.$notify.info(this.$gettext("No changes to save."));
        this.resetState();
        return;
      }
      this.busy = true;
      this.$notify.blockUI("busy");
      const markerDataToUpdate = {
        FileUID: marker.FileUID,
        Type: marker.Type,
        Src: marker.Src,
        SubjSrc: marker.SubjSrc,
        Name: marker.Name,
        MarkerReview: marker.Review,
        Invalid: marker.Invalid,
        X: marker.X,
        Y: marker.Y,
        W: marker.W,
        H: marker.H,
        UpdatePosition: true,
      };
      $api
        .put(`markers/${marker.UID}`, markerDataToUpdate)
        .then((response) => {
          const serverMarkerData = response.data;
          const index = this.markers.findIndex((m) => m.UID === marker.UID);
          if (index !== -1) {
            const localMarker = this.markers[index];
            Object.assign(localMarker, serverMarkerData);
            if (localMarker.Thumb) {
              const timestamp = new Date().getTime();
              const baseUrl = localMarker.thumbnailUrl("tile_320");
              localMarker.thumbWithTimestamp = `${baseUrl}?ts=${timestamp}`;
            }
            this.updateMarkerElement(localMarker);
          }
          this.$notify.success(this.$gettext("Face marker updated"));
          this.$emit("markers-updated", this.markers);
        })
        .catch((error) => {
          console.error("Error updating marker:", error);
          this.$notify.error(this.$gettext("Failed to update face marker"));
          const index = this.markers.findIndex((m) => m.UID === marker.UID);
          if (index !== -1 && this.interaction.initialMarkerRect) {
            Object.assign(this.markers[index], this.interaction.initialMarkerRect);
            this.updateMarkerElement(this.markers[index]);
          }
        })
        .finally(() => {
          this.busy = false;
          this.$notify.unblockUI();
          this.resetState();
        });
    },
    toggleDrawingMode() {
      if (this.isDrawingMode) {
        this.resetState();
      } else {
        this.mode = "drawing";
        this.selectedMarker = null;
      }
    },
    onReject(model) {
      if (this.busy || !model) return;
      this.busy = true;
      this.$notify.blockUI("busy");
      model.reject().finally(() => {
        this.$notify.unblockUI();
        this.busy = false;
        this.$emit("markers-updated", this.markers);
      });
    },
    resetState() {
      this.mode = null;
      this.selectedMarker = null;
      this.interaction = {
        active: false,
        type: null,
        startX: 0,
        startY: 0,
        wrapperRect: null,
        initialMarkerRect: null,
        resizeHandle: null,
        drawingPreview: null,
      };
      this.removeDocumentListeners();
    },
    startInteraction(type, marker = null, resizeHandle = null) {
      this.interaction = {
        active: true,
        type,
        startX: 0,
        startY: 0,
        wrapperRect: this.$refs.photoPreviewWrapper?.getBoundingClientRect() || null,
        initialMarkerRect: marker ? { X: marker.X, Y: marker.Y, W: marker.W, H: marker.H } : null,
        resizeHandle,
        drawingPreview: type === "drawing" ? { x: 0, y: 0, w: 0, h: 0 } : null,
      };
      this.addDocumentListeners();
    },
    saveMarkerChanges() {
      if (this.selectedMarker) {
        this.updateMarkerPosition(this.selectedMarker);
      }
      this.resetState();
    },
  },
};
</script>
