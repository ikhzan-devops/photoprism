<template>
  <div v-if="visible" ref="container" class="p-viewer" tabindex="-1" role="dialog">
    <div ref="lightbox" class="p-viewer__lightbox" :class="{ slideshow: slideshow.active, sidebar: sidebarVisible, 'is-favorite': model.Favorite }"></div>
    <div v-if="sidebarVisible" ref="sidebar" class="p-viewer__sidebar"></div>

    <!-- TODO: All previously available features and controls must be preserved in the new hybrid photo/video viewer:
    <div class="pswp__ui pswp__ui--hidden">
      <div class="pswp__top-bar">
        <div class="pswp__taken hidden-xs">
          {{ formatDate(slide.TakenAtLocal) }}
        </div>

        <div class="pswp__counter"></div>

        <button class="pswp__button pswp__button--close action-close" :title="$gettext('Close')"></button>

        <button v-if="canDownload" class="pswp__button action-download" style="background: none" :title="$gettext('Download')" @click.exact="onDownload">
          <v-icon size="16" color="white">mdi-download</v-icon>
        </button>

        <button v-if="canEdit" class="pswp__button action-edit hidden-shared-only" style="background: none" :title="$gettext('Edit')" @click.exact="onEdit">
          <v-icon size="16" color="white">mdi-pencil</v-icon>
        </button>

        <button class="pswp__button action-select" style="background: none" :title="$gettext('Select')" @click.exact="onSelect">
          <v-icon v-if="selection.length && $clipboard.has(item)" size="16" color="white">mdi-check-circle</v-icon>
          <v-icon v-else size="16" color="white">mdi-circle-outline</v-icon>
        </button>

        <button v-if="canLike" class="pswp__button action-like hidden-shared-only" style="background: none" :title="$gettext('Like')" @click.exact="onLike">
          <v-icon v-if="slide.Favorite" icon="mdi-star" size="19" color="white"></v-icon>
          <v-icon v-else icon="mdi-star-outline" size="19" color="white"></v-icon>
        </button>

        <button class="pswp__button pswp__button--fs action-toggle-fullscreen" :title="$gettext('Fullscreen')"></button>

        <button class="pswp__button pswp__button--zoom action-zoom" :title="$gettext('Zoom in/out')"></button>

        <button class="pswp__button action-slideshow" style="background: none" :title="$gettext('Start/Stop Slideshow')" @click.exact="onSlideshow">
          <v-icon v-show="!interval" size="18" color="white">mdi-play</v-icon>
          <v-icon v-show="interval" size="16" color="white">mdi-pause</v-icon>
        </button>
      </div>

      <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
        <div class="pswp__share-tooltip"></div>
      </div>
    </div>
     -->
  </div>
</template>

<script>
import PhotoSwipe from "photoswipe";
import Lightbox from "photoswipe/lightbox";
import PhotoSwipeDynamicCaption from "photoswipe-dynamic-caption-plugin";
import Util from "common/util";
import Api from "common/api";
import Thumb from "model/thumb";
import { Photo } from "model/photo";

/*
  TODO: All previously available features and controls must be preserved in the new hybrid photo/video viewer:
    1. Some of the controls that the old viewer had (e.g. (a) select, (b) play slideshow, (c) fullscreen,
       (d) edit, (e) date info,...) are still missing.
    2. The already added controls may need some improvements (e.g. (a) the sidebar toggle button (info icon) shows
       the sidebar, but the functionality there is not implemented yet, (b) the zoom doesn't load a larger version
       of the image yet).
    3. Finally, after the refactoring/upgrade, (a) the old/unused code (e.g. for the separate video player) needs
       to be removed and (b) everything needs to be thoroughly tested on all major browsers and mobile devices.
*/
export default {
  name: "PViewer",
  data() {
    return {
      visible: false,
      sidebarVisible: false,
      lightbox: null, // Current PhotoSwipe lightbox instance.
      captionPlugin: null, // Current PhotoSwipe caption plugin instance.
      captionTimer: false,
      hasTouch: false,
      idleTime: 6000, // Automatically hide viewer controls after 6 seconds until user settings are implemented.
      controlsShown: -1, // -1 or a positive Date.now() timestamp indicates that the PhotoSwipe controls are shown.
      canEdit: this.$config.allow("photos", "update") && this.$config.feature("edit"),
      canLike: this.$config.allow("photos", "manage") && this.$config.feature("favorites"),
      canDownload: this.$config.allow("photos", "download") && this.$config.feature("download"),
      experimental: this.$config.get("experimental"), // Experimental features flag.
      selection: this.$clipboard.selection,
      config: this.$config.values,
      model: new Thumb(), // Current slide.
      models: [],
      index: 0,
      subscriptions: [], // Event subscriptions.
      interval: false,
      slideshow: {
        active: false,
        next: 0,
      },
    };
  },
  created() {
    // this.subscriptions["viewer.change"] = this.$event.subscribe("viewer.change", this.onChange);
    this.subscriptions["viewer.pause"] = this.$event.subscribe("viewer.pause", this.onPause);
    // this.subscriptions["viewer.show"] = this.$event.subscribe("viewer.show", this.onShow);
    this.subscriptions["viewer.close"] = this.$event.subscribe("viewer.close", this.onClose);
  },
  beforeUnmount() {
    this.onPause();
    this.destroyLightbox();

    for (let i = 0; i < this.subscriptions.length; i++) {
      this.$event.unsubscribe(this.subscriptions[i]);
    }
  },
  methods: {
    // Returns the PhotoSwipe container HTML element, if visible.
    getLightbox() {
      return this.$refs?.lightbox;
    },
    // Returns the PhotoSwipe config options, see https://photoswipe.com/options/.
    getLightboxOptions() {
      return {
        appendToEl: this.getLightbox(),
        pswpModule: PhotoSwipe,
        dataSource: this.models,
        index: this.index,
        mouseMovePan: true,
        arrowPrev: true,
        arrowNext: true,
        loop: true,
        zoom: true,
        close: true,
        counter: false,
        trapFocus: false,
        returnFocus: false,
        initialZoomLevel: "fit",
        secondaryZoomLevel: "fill",
        maxZoomLevel: 3,
        bgOpacity: 1,
        preload: [1, 1],
        showHideAnimationType: "none",
        tapAction: (point, e) => this.toggleControls(e),
        imageClickAction: "zoom",
        mainClass: "media-viewer-lightbox",
        bgClickAction: (point, e) => this.onBgClick(e),
        paddingFn: (s) => this.getLightboxPadding(s),
        getViewportSizeFn: () => this.getLightboxViewport(),
        closeTitle: this.$gettext("Close"),
        zoomTitle: this.$gettext("Zoom"),
        arrowPrevTitle: this.$gettext("Previous"),
        arrowNextTitle: this.$gettext("Next"),
        errorMsg: this.$gettext("Error"),
      };
    },
    // Displays the thumbnail images and/or videos that belong to the specified models in the lightbox.
    showThumbs(models, index = 0) {
      // Check if at least one model was passed, as otherwise no content can be displayed.
      if (!Array.isArray(models) || models.length === 0 || index >= models.length) {
        console.log("model list passed to viewer is empty:", models);
        return;
      }

      this.onShow();

      this.$nextTick(() => {
        this.renderLightbox(models, index);
      });
    },
    // Loads the pictures that belong to the page context and displays them in the lightbox.
    showContext(ctx, index) {
      if (ctx.loading || !ctx.listen || ctx.viewer.loading || !ctx.results[index]) {
        return false;
      }

      const selected = ctx.results[index];

      if (!ctx.viewer.dirty && ctx.viewer.results && ctx.viewer.results.length > index) {
        // Reuse existing viewer result if possible.
        let i = -1;

        if (ctx.viewer.results[index] && ctx.viewer.results[index].UID === selected.UID) {
          i = index;
        } else {
          i = ctx.viewer.results.findIndex((p) => p.UID === selected.UID);
        }

        if (i > -1 && (((ctx.viewer.complete || ctx.complete) && ctx.viewer.results.length >= ctx.results.length) || i + ctx.viewer.batchSize <= ctx.viewer.results.length)) {
          this.showThumbs(ctx.viewer.results, i);
          return;
        }
      }

      // Fetch photos from server API.
      ctx.viewer.loading = true;

      const params = ctx.searchParams();
      params.count = params.offset + ctx.viewer.batchSize;
      params.offset = 0;

      // Fetch viewer results from API.
      return Api.get("photos/view", { params })
        .then((response) => {
          const count = response && response.data ? response.data.length : 0;
          if (count === 0) {
            ctx.$notify.warn(ctx.$gettext("No pictures found"));
            ctx.viewer.dirty = true;
            ctx.viewer.complete = false;
            return;
          }

          // Process response.
          if (response.headers && response.headers["x-count"]) {
            const c = parseInt(response.headers["x-count"]);
            const l = parseInt(response.headers["x-limit"]);
            ctx.viewer.complete = c < l;
          } else {
            ctx.viewer.complete = ctx.complete;
          }

          let i;

          if (response.data[index] && response.data[index].UID === selected.UID) {
            i = index;
          } else {
            i = response.data.findIndex((p) => p.UID === selected.UID);
          }

          ctx.viewer.results = Thumb.wrap(response.data);

          // Show pictures.
          this.showThumbs(ctx.viewer.results, i);
          ctx.viewer.dirty = false;
        })
        .catch(() => {
          ctx.viewer.dirty = true;
          ctx.viewer.complete = false;
        })
        .finally(() => {
          // Unblock.
          ctx.viewer.loading = false;
        });
    },
    // Initializes and opens the PhotoSwipe lightbox with the
    // images and/or videos that belong to the specified models.
    renderLightbox(models, index = 0) {
      // Check if at least one model was passed, as otherwise no content can be displayed.
      if (!Array.isArray(models) || models.length === 0 || index >= models.length) {
        console.log("model list passed to viewer is empty:", models);
        return;
      }

      // Set the initial model list and start index.
      // TODO: In the future, additional models should be dynamically loaded when the index reaches the end of the list.
      this.models = models;
      this.index = index;

      // Focus lightbox element.
      this.getLightbox().focus();

      // Get PhotoSwipe lightbox config options, see https://photoswipe.com/options/.
      const options = this.getLightboxOptions();

      // Create PhotoSwipe instance.
      let lightbox = new Lightbox(options);
      let firstPicture = true;

      // Keep reference to PhotoSwipe instance.
      this.lightbox = lightbox;
      this.captionTimer = false;
      this.hasTouch = false;

      // Use dynamic caption plugin,
      // see https://github.com/dimsemenov/photoswipe-dynamic-caption-plugin.
      this.captionPlugin = new PhotoSwipeDynamicCaption(lightbox, {
        type: "auto",
        captionContent: (slide) => {
          if (!slide || !this.models || slide?.index < 0) {
            return "";
          }

          const model = this.models[slide.index];

          if (model) {
            return this.formatCaption(model);
          }

          return "";
        },
      });

      // Add a close event handler to destroy the viewer after use,
      // see https://photoswipe.com/events/#closing-events.
      lightbox.on("close", () => {
        this.$event.publish("viewer.pause");
        this.$event.publish("viewer.close");
      });

      // Add viewer controls, see https://photoswipe.com/adding-ui-elements/.
      //
      // TODO: The same controls as with PhotoSwipe 4 should be usable/available!
      lightbox.on("uiRegister", () => {
        // Add a sidebar toggle button only if the window is large enough.
        // TODO: Proof-of-concept only, the sidebar needs to be fully implemented before this can be released.
        // TODO: Once this is fully implemented, remove the "this.experimental" flag check below.
        // IDEA: We can later try to add styles that display the sidebar at the bottom
        //       instead of on the side, to allow use on mobile devices.
        if (this.experimental && this.canEdit && window.innerWidth > 600) {
          lightbox.pswp.ui.registerElement({
            name: "sidebar-button",
            className: "pswp__button--sidebar-button pswp__button--mdi", // Sets the icon style/size in viewer.css.
            order: 9,
            isButton: true,
            html: {
              isCustomSVG: true,
              inner: '<path d="M11 7V9H13V7H11M14 17V15H13V11H10V13H11V15H10V17H14M22 12C22 17.5 17.5 22 12 22C6.5 22 2 17.5 2 12C2 6.5 6.5 2 12 2C17.5 2 22 6.5 22 12M20 12C20 7.58 16.42 4 12 4C7.58 4 4 7.58 4 12C4 16.42 7.58 20 12 20C16.42 20 20 16.42 20 12Z" id="pswp__icn-sidebar"/>',
              outlineID: "pswp__icn-sidebar", // Add this to the <path> in the inner property.
              size: 24, // Depends on the original SVG viewBox, e.g. use 24 for viewBox="0 0 24 24".
            },
            onClick: (e) => {
              return this.toggleSidebar(e);
            },
          });
        }

        // Add download button if user has permission to download pictures,
        // see https://photoswipe.com/adding-ui-elements/.
        if (this.canDownload) {
          lightbox.pswp.ui.registerElement({
            name: "download-button",
            className: "pswp__button--download-button pswp__button--mdi", // Sets the icon style/size in viewer.css.
            order: 10,
            isButton: true,
            html: {
              isCustomSVG: true,
              inner: `<path d="M5,20H19V18H5M19,9H15V3H9V9H5L12,16L19,9Z" id="pswp__icn-download" />`,
              outlineID: "pswp__icn-download", // Add this to the <path> in the inner property.
              size: 24, // Depends on the original SVG viewBox, e.g. use 24 for viewBox="0 0 24 24".
            },
            onClick: (e) => {
              return this.onDownload(e);
            },
          });
        }

        // Add edit button if user has permission to edit pictures,
        // see https://photoswipe.com/adding-ui-elements/.
        if (this.canEdit) {
          lightbox.pswp.ui.registerElement({
            name: "edit-button",
            className: "pswp__button--edit-button pswp__button--mdi", // Sets the icon style/size in viewer.css.
            order: 10,
            isButton: true,
            html: {
              isCustomSVG: true,
              inner: `<path d="M20.71,7.04C21.1,6.65 21.1,6 20.71,5.63L18.37,3.29C18,2.9 17.35,2.9 16.96,3.29L15.12,5.12L18.87,8.87M3,17.25V21H6.75L17.81,9.93L14.06,6.18L3,17.25Z" id="pswp__icn-edit" />`,
              outlineID: "pswp__icn-edit", // Add this to the <path> in the inner property.
              size: 26, // Depends on the original SVG viewBox, e.g. use 24 for viewBox="0 0 24 24".
            },
            onClick: () => {
              return this.onEdit();
            },
          });
        }

        if (this.canLike) {
          lightbox.pswp.ui.registerElement({
            name: "toggle-like-button",
            className: "pswp__button--mdi pswp__button--toggle-like-button", // Sets the icon style/size in viewer.css.
            order: 10,
            isButton: true,
            html: {
              isCustomSVG: true,
              inner: `<path d="M12,17.27L18.18,21L16.54,13.97L22,9.24L14.81,8.62L12,2L9.19,8.62L2,9.24L7.45,13.97L5.82,21L12,17.27Z" id="pswp__icn-like" /><path d="M12,15.39L8.24,17.66L9.23,13.38L5.91,10.5L10.29,10.13L12,6.09L13.71,10.13L18.09,10.5L14.77,13.38L15.76,17.66M22,9.24L14.81,8.63L12,2L9.19,8.63L2,9.24L7.45,13.97L5.82,21L12,17.27L18.18,21L16.54,13.97L22,9.24Z" id="pswp__icn-like-outline" />`,
              size: 24, // Depends on the original SVG viewBox, e.g. use 24 for viewBox="0 0 24 24".
            },
            onClick: () => {
              return this.onLike();
            },
          });
        }
      });

      // Trigger onChange() event handler when slide is changed and on initialization,
      // see https://photoswipe.com/events/#initialization-events.
      this.lightbox.on("change", () => {
        this.onChange();
      });

      // Processes model data for rendering slides with PhotoSwipe, see https://photoswipe.com/filters/#itemdata.
      lightbox.addFilter("itemData", (el, i) => {
        /*
         TODO: Rendering of slides needs to be improved to allow dynamic zooming (loading higher resolution thumbs
               depending on zoom level) and playing videos in their native format whenever possible (see below).
        */

        // Get the current slide model data.
        const model = this.models[i];

        // Get the screen (window) resolution in real pixels,
        // depending on the width/height and pixel density.
        const pixels = this.getWindowPixels();

        // Get the right thumbnail size based on the screen resolution in pixels.
        const s = Util.thumbSize(pixels.width, pixels.height);

        // Get thumbnail image URL.
        const imgSrc = model.Thumbs[s].src;

        // Render videos and animations as custom HTML.
        if (model.Playable) {
          const videoSrc = Util.videoUrl(model.Hash);
          /*
            TODO: (a) Check if there is a more convenient and/or secure way to render the video slide, then perform
                      security tests to ensure that no code can be injected, e.g. create an HTMLVideoElement object,
                      set the properties based on the media type/video duration, and then return it instead of the
                      plain HTML as implemented in the proof-of-concept.
                  (b) Live Photos and Animations (e.g. GIFs) must be looped and played automatically (autoplay attribute).
                  (c) If the browser can naively handle the video file format, don't default to the AVC video URL, as
                      this may require transcoding, which is slow and resource-intensive. For this, the Util.videoUrl()
                      function has a second argument for the codec (might need to be added to the server response,
                      which is something we can help with).

                  Once this is released, the following enhancements can be worked on and shipped in a future release:

                  (d) We should consider using the .m3u8 file format for specifying the stream URL(s), so that the
                      browser can choose the best format/codec (first develop a simple/static proof-of-concept to see
                      if/how it works).
                  (e) The server should (additionally) provide a video/animation still from time index 0 that can be
                      used as a poster (the current thumbnail is taken later for longer videos, since the first frame is
                      often black).
          */
          if (firstPicture) {
            firstPicture = false;
            return {
              html: `<video class="pswp__video" autoplay controls playsinline poster="${imgSrc}" preload="auto"><source src="${videoSrc}" /></video>`,
            };
          } else {
            return {
              html: `<video class="pswp__video" controls playsinline poster="${imgSrc}" preload="metadata"><source src="${videoSrc}" /></video>`,
            };
          }
        }

        if (firstPicture) {
          firstPicture = false;
        }

        // Return the data that PhotoSwipe needs to show the image,
        // see https://photoswipe.com/data-sources/#dynamically-generated-data.
        return {
          src: imgSrc, // Thumbnail image URL.
          width: model.Thumbs[s].w, // Actual thumbnail image width (x).
          height: model.Thumbs[s].h, // Actual thumbnail image height (y).
        };
      });

      // Init PhotoSwipe.
      lightbox.init();

      // Show first image.
      lightbox.loadAndOpen(index);

      // Publish event to be consumed by other components.
      this.$event.publish("viewer.opened");
    },
    // Destroys the PhotoSwipe lightbox instance after use, see onClose().
    destroyLightbox() {
      if (this.lightbox) {
        this.lightbox?.destroy();
        this.lightbox = null;
      }
    },
    // Returns the picture (model) caption as sanitized HTML, if any.
    formatCaption(model) {
      if (!model) {
        return "";
      }

      let caption = "";

      if (model.Title) {
        caption += `<h4>${Util.encodeHTML(model.Title)}</h4>`;
      }

      /* TODO: Find a good position for the date information that
               works for all screen sizes and image dimensions. */
      /* if (model.TakenAtLocal) {
         caption += `<div>${Util.formatDate(model.TakenAtLocal)}</div>`;
      } */

      if (model.Caption) {
        caption += `<p>${Util.encodeHTML(model.Caption)}</p>`;
      } else if (model.Description) {
        caption += `<p>${Util.encodeHTML(model.Description)}</p>`;
      }

      // TODO: Perform security tests to see if unwanted code can be injected.
      return Util.sanitizeHtml(caption);
    },
    onShow() {
      // Hide the browser scrollbar as it is not wanted in the viewer.
      this.$scrollbar.hide();

      // Render the component template.
      this.visible = true;

      // Publish event to be consumed by other components.
      this.$event.publish("viewer.show");
    },
    // Destroys the PhotoSwipe lightbox, resets the component state, and unhides the browser scrollbar.
    onClose() {
      // Pause slideshow and any videos that are playing.
      this.onPause();

      // Destroy PhotoSwipe lightbox.
      this.destroyLightbox();

      // Reset component state.
      this.onReset();

      // Hide lightbox and sidebar.
      this.hideViewer();

      // Publish event to be consumed by other components.
      this.$event.publish("viewer.closed");
    },
    // Pauses the lightbox slideshow and any videos that are playing.
    onPause() {
      this.pauseVideos();
      this.pauseSlideshow();
    },
    // Resets the component state after closing the lightbox.
    onReset() {
      this.resetTimer();
      this.resetControls();
      this.resetModels();
    },
    // Resets the timer for hiding the viewer controls.
    resetTimer() {
      if (this.captionTimer) {
        window.clearTimeout(this.captionTimer);
        this.captionTimer = false;
      }
    },
    // Resets the state of the viewer controls.
    resetControls() {
      this.hasTouch = false;
      this.controlsShown = -1;
    },
    // Reset the viewer models and index.
    resetModels() {
      this.model = new Thumb();
      this.models = [];
      this.index = 0;
    },
    // Hides the viewer and restores the scrollbar state.
    hideViewer() {
      // Hide sidebar.
      this.hideSidebar();

      // Remove lightbox focus and hide viewer.
      if (this.visible) {
        this.$refs?.lightbox?.blur();
        this.visible = false;
      }

      // Restore browser scrollbar state.
      this.$scrollbar.show();
    },
    // Returns the active PhotoSwipe instance, if any.
    // Be sure to check the result before using it!
    pswp() {
      return this.lightbox?.pswp;
    },
    // Called when the slide is changed and on initialization,
    // see https://photoswipe.com/events/#initialization-events.
    onChange() {
      // Get active PhotoSwipe instance.
      const pswp = this.pswp();

      if (!pswp) {
        return;
      }

      // Find and pause videos that are currently playing.
      this.pauseVideos();

      // Attach touch and mouse event handlers to automatically hide controls.
      document.addEventListener(
        "touchstart",
        () => {
          this.resetTimer();
          this.hasTouch = true;
        },
        { once: true }
      );
      document.addEventListener(
        "mousemove",
        () => {
          this.startTimer();
        },
        { once: true }
      );

      // Set current slide (model) list index.
      if (typeof pswp.currIndex === "number") {
        this.index = pswp.currIndex;
      }

      // Set current slide model.
      if (this.index >= 0 && this.models.length > 0 && this.index < this.models.length) {
        this.model = this.models[this.index];
      }

      // Pause the slideshow if the index of the next slide does not match.
      if (this.slideshow.next !== this.index) {
        this.pauseSlideshow();
      }
    },
    // Called when the user clicks on the PhotoSwipe lightbox background,
    // see https://photoswipe.com/click-and-tap-actions.
    onBgClick(e) {
      if (this.controlsVisible()) {
        this.onClose();
      } else {
        this.showControls();
      }

      if (e && typeof e.stopPropagation === "function") {
        e.stopPropagation();
      }
    },
    onLike() {
      this.model.toggleLike();
    },
    // TODO: Toggles the selection of the current picture in the global photo clipboard.
    onSelect() {
      this.$clipboard.toggle(this.model);
    },
    // Returns the <video> elements in the lightbox container as an HTMLCollection.
    getVideos() {
      const el = this.getLightbox();

      // Call https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementsByTagName to find videos.
      if (el) {
        return el.getElementsByTagName("video");
      }

      return [];
    },
    // Finds and pauses all videos currently playing in the lightbox.
    pauseVideos() {
      const videos = this.getVideos();

      if (!videos || !videos.length) {
        return false;
      }

      for (let video of videos) {
        if (typeof video.pause === "function") {
          try {
            // Calling pause() before a play promise has been resolved may result in an error,
            // see https://github.com/flutter/flutter/issues/136309 (we'll ignore this for now).
            if (!video.paused) {
              video.pause();
            }
          } catch (e) {
            console.log(e);
          }
        }
      }
    },
    // Pauses the lightbox slideshow, if currently active.
    pauseSlideshow() {
      this.slideshow.active = false;

      if (this.interval) {
        clearInterval(this.interval);
        this.interval = false;
      }
    },
    // Toggles the lightbox slideshow.
    // TODO: Does not work yet, needs to be reimplemented for the new viewer.
    onSlideshow() {
      if (this.interval) {
        this.pauseSlideshow();
        return;
      }

      this.slideshow.active = true;

      const pswp = this.pswp();

      self.interval = setInterval(() => {
        if (pswp && typeof pswp.next === "function") {
          pswp.next();
          this.slideshow.next = pswp.currIndex;
        } else {
          this.pauseSlideshow();
        }
      }, 5000);
    },
    // Downloads the original files of the current picture.
    onDownload(e) {
      if (e && typeof e.stopPropagation === "function") {
        e.stopPropagation();
      }

      this.pauseSlideshow();

      /* TODO: Once all the viewer's core functionality has been restored, add a file size/type
               selection dialog so the user can choose which format and quality to download. */

      if (!this.model || !this.model.DownloadUrl) {
        console.warn("photo viewer: no download url");
        return;
      }

      this.$notify.success(this.$gettext("Downloading…"));

      new Photo().find(this.model.UID).then((p) => p.downloadAll());
    },
    onEdit() {
      this.onPause();

      const pswp = this.pswp();
      let index = 0;

      // remove duplicates
      let filtered = this.models?.filter(function (p, i, s) {
        return !(i > 0 && p.UID === s[i - 1].UID);
      });

      let selection = filtered.map((p, i) => {
        if (this.model.UID === p.UID) {
          index = i;
        }

        return p.UID;
      });

      let album = null;

      pswp.close(); // Close Gallery

      this.$event.publish("dialog.edit", { selection, album, index }); // Open Edit Dialog
    },
    toggleSidebar(e) {
      this.sidebarVisible = !this.sidebarVisible;

      this.$nextTick(() => {
        const pswp = this.pswp();
        if (pswp) {
          pswp.updateSize(true);
        }
      });

      if (e && typeof e.stopPropagation === "function") {
        e.stopPropagation();
      }
    },
    // Hides the viewer sidebar, if visible.
    hideSidebar() {
      if (this.sidebarVisible) {
        this.$refs?.sidebar?.blur();
        this.sidebarVisible = false;
      }
    },
    toggleControls(e) {
      if (this.pswp() && this.pswp().element) {
        const el = this.pswp().element;
        if (el.classList.contains("pswp--ui-visible")) {
          this.controlsShown = 0;
          el.classList.remove("pswp--ui-visible");
        } else {
          this.controlsShown = Date.now();
          el.classList.add("pswp--ui-visible");
        }
      }

      if (e && typeof e.stopPropagation === "function") {
        e.stopPropagation();
      }
    },
    showControls() {
      if (this.pswp() && this.pswp().element) {
        this.controlsShown = Date.now();
        this.pswp().element.classList.add("pswp--ui-visible");
      }
    },
    hideControls() {
      if (this.pswp() && this.pswp().element) {
        this.controlsShown = 0;
        this.pswp().element.classList.remove("pswp--ui-visible");
      }
    },
    controlsVisible() {
      if (this.controlsShown === 0) {
        return false;
      } else if (this.controlsShown < 0) {
        return true;
      }

      if (this.pswp() && this.pswp().element) {
        const el = this.pswp().element;
        if (el.classList.contains("pswp--ui-visible") && Date.now() - this.controlsShown > 120) {
          return true;
        }
      }

      return false;
    },
    mouseMove() {
      this.resetTimer();
      if (this.lightbox) {
        this.showControls();
        this.startTimer();
      }
    },
    startTimer() {
      if (this.hasTouch) {
        return;
      }

      this.resetTimer();
      this.captionTimer = window.setTimeout(() => {
        this.hideControls();
      }, this.idleTime);
      document.addEventListener(
        "mousemove",
        () => {
          this.mouseMove();
        },
        { once: true }
      );
    },
    getWindowPixels() {
      return {
        width: window.innerWidth * window.devicePixelRatio,
        height: window.innerHeight * window.devicePixelRatio,
      };
    },
    getLightboxViewport() {
      const el = this.getLightbox();

      if (el) {
        return {
          x: el.clientWidth,
          y: el.clientHeight,
        };
      } else {
        return {
          x: window.innerWidth,
          y: window.innerHeight,
        };
      }
    },
    getLightboxPadding(s) {
      if (!s || (s.x <= 600 && s.x < s.y)) {
        // Vertical padding on mobile screens to avoid obscuring controls (except when zooming into pictures).
        return {
          top: 56,
          bottom: 8,
          left: 0,
          right: 0,
        };
      } else if (s.x === 720 || s.x === 1080 || s.x === 1280 || s.x === 1920 || s.x === 2560 || s.x === 3840 || s.x === 4096 || s.x === 4096 || s.x === 7680) {
        // Viewport has a standardized size, e.g. on a TV or a browser in full-screen mode.
        return {
          top: 0,
          bottom: 0,
          left: 0,
          right: 0,
        };
      } else {
        // Default.
        return {
          top: 4,
          bottom: 4,
          left: 4,
          right: 4,
        };
      }
    },
  },
};
</script>
