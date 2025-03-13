<template>
  <v-form
    ref="form"
    validate-on="invalid-input"
    autocomplete="off"
    class="p-photo-toolbar p-album-toolbar"
    accept-charset="UTF-8"
    @submit.prevent="updateQuery()"
  >
    <v-toolbar
      flat
      :density="$vuetify.display.smAndDown ? 'compact' : 'default'"
      class="page-toolbar"
      color="secondary"
    >
      <v-toolbar-title :title="album.Title" class="flex-grow-1">
        <span class="hidden-xs">
          <router-link :to="{ name: collectionRoute }">
            {{ T(collectionTitle) }}
          </router-link>
          <v-icon>{{ navIcon }}</v-icon>
        </span>
        {{ album.Title }}
      </v-toolbar-title>

      <v-btn-toggle
        :model-value="settings.view"
        :title="$gettext('Toggle View')"
        :density="$vuetify.display.smAndDown ? 'comfortable' : 'default'"
        base-color="secondary"
        variant="flat"
        rounded="pill"
        mandatory
        border
        group
        class="ms-1"
      >
        <v-btn value="cards" icon="mdi-view-column" class="ps-1" @click="setView('cards')"></v-btn>
        <v-btn v-if="listView" value="list" icon="mdi-view-list" @click="setView('list')"></v-btn>
        <v-btn value="mosaic" icon="mdi-view-comfy" class="pe-1" @click="setView('mosaic')"></v-btn>
      </v-btn-toggle>

      <v-menu transition="slide-y-transition" open-on-click open-on-hover>
        <template #activator="{ props }">
          <v-btn density="comfortable" icon="mdi-dots-vertical" v-bind="props" class="action-menu ms-1"></v-btn>
        </template>

        <v-list min-width="128" density="comfortable" bg-color="navigation" slim class="ra-8 opacity-95">
          <v-list-item
            prepend-icon="mdi-refresh"
            :subtitle="$gettext('Refresh')"
            class="action-reload action-refresh"
            @click="refresh"
          ></v-list-item>
          <v-list-item
            v-if="canManage"
            :subtitle="$gettext('Edit')"
            prepend-icon="mdi-pencil"
            class="action-edit"
            @click="dialog.edit = true"
          ></v-list-item>
          <v-list-item
            v-if="canShare"
            :subtitle="$gettext('Share')"
            prepend-icon="mdi-share-variant"
            class="action-share"
            @click="dialog.share = true"
          ></v-list-item>
          <v-list-item
            v-if="canDownload"
            :subtitle="$gettext('Download')"
            prepend-icon="mdi-download"
            class="action-download"
            @click="download"
          ></v-list-item>
          <v-list-item
            v-if="canUpload"
            :subtitle="$gettext('Upload')"
            prepend-icon="mdi-cloud-upload"
            class="action-upload"
            @click="showUpload"
          ></v-list-item>
        </v-list>
      </v-menu>
    </v-toolbar>

    <div v-if="album.Description" class="toolbar-details-panel">
      {{ album.Description }}
    </div>

    <p-share-dialog
      :visible="dialog.share"
      :model="album"
      @upload="webdavUpload"
      @close="dialog.share = false"
    ></p-share-dialog>
    <p-service-upload
      :visible="dialog.upload"
      :items="{ albums: album.getId() }"
      :model="album"
      @close="dialog.upload = false"
      @confirm="dialog.upload = false"
    ></p-service-upload>
    <p-album-edit-dialog :visible="dialog.edit" :album="album" @close="dialog.edit = false"></p-album-edit-dialog>
  </v-form>
</template>
<script>
import $notify from "common/notify";
import download from "common/download";
import { T } from "common/gettext";

export default {
  name: "PAlbumToolbar",
  props: {
    album: {
      type: Object,
      default: () => {},
    },
    filter: {
      type: Object,
      default: () => {},
    },
    updateFilter: {
      type: Function,
      default: () => {},
    },
    updateQuery: {
      type: Function,
      default: () => {},
    },
    settings: {
      type: Object,
      default: () => {},
    },
    refresh: {
      type: Function,
      default: () => {},
    },
  },
  data() {
    const cameras = [
      {
        ID: 0,
        Name: this.$gettext("All Cameras"),
      },
    ].concat(this.$config.get("cameras"));
    const countries = [
      {
        ID: "",
        Name: this.$gettext("All Countries"),
      },
    ].concat(this.$config.get("countries"));
    const features = this.$config.getSettings().features;
    return {
      expanded: false,
      canUpload: this.$config.allow("files", "upload") && features.upload,
      canDownload: this.$config.allow("albums", "download") && features.download,
      canShare: this.$config.allow("albums", "share") && features.share,
      canManage: this.$config.allow("albums", "manage"),
      experimental: this.$config.get("experimental"),
      isFullScreen: !!document.fullscreenElement,
      categories: this.$config.albumCategories(),
      collectionTitle: this.$route.meta?.collectionTitle ? this.$route.meta.collectionTitle : this.$gettext("Albums"),
      collectionRoute: this.$route.meta?.collectionRoute ? this.$route.meta.collectionRoute : "albums",
      navIcon: this.$isRtl ? "mdi-chevron-left" : "mdi-chevron-right",
      listView: this.$config.getSettings()?.search?.listView,
      dialog: {
        share: false,
        upload: false,
        edit: false,
      },
      titleRule: (v) => v.length <= this.$config.get("clip") || this.$gettext("Name too long"),
    };
  },
  methods: {
    hideExpansionPanel() {
      if (this.expanded) {
        this.expanded = false;
      }
    },
    T() {
      return T.apply(this, arguments);
    },
    webdavUpload() {
      this.dialog.share = false;
      this.dialog.upload = true;
    },
    showUpload() {
      // Pre-select manually managed albums in upload dialog.
      if (this.album.Type === "album") {
        this.$event.publish("dialog.upload", { albums: [this.album] });
      } else {
        this.$event.publish("dialog.upload", { albums: [] });
      }
    },
    onUpdate(v) {
      this.updateQuery(v);
    },
    setView(name) {
      if (name) {
        if (name === "list" && !this.listView) {
          name = "mosaic";
        }

        this.refresh({ view: name });
      }
    },
    download() {
      this.onDownload(`${this.$config.apiUri}/albums/${this.album.UID}/dl?t=${this.$config.downloadToken}`);
    },
    onDownload(path) {
      $notify.success(this.$gettext("Downloading…"));

      download(path, "album.zip");
    },
  },
};
</script>
