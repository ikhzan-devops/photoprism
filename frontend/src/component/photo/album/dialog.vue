<template>
  <v-dialog
    :model-value="visible"
    persistent
    max-width="500"
    class="p-dialog p-photo-album-dialog"
    @keydown.esc.exact="close"
    @after-enter="afterEnter"
    @after-leave="afterLeave"
  >
    <v-form ref="form" validate-on="invalid-input" accept-charset="UTF-8" tabindex="1" @submit.prevent="confirm">
      <v-card>
        <v-card-title class="d-flex justify-start align-center ga-3">
          <v-icon icon="mdi-bookmark" size="28" color="primary"></v-icon>
          <h6 class="text-h6">{{ $gettext(`Add to album`) }}</h6>
        </v-card-title>
        <v-card-text>
          <v-combobox
            ref="input"
            v-model="selectedAlbums"
            :disabled="loading"
            :loading="loading"
            hide-details
            chips
            closable-chips
            multiple
            class="input-albums"
            :items="items"
            item-title="Title"
            item-value="UID"
            :placeholder="$gettext('Select or create albums')"
            return-object
          >
            <template #no-data>
              <v-list-item>
                <v-list-item-title>
                  {{ $gettext(`Press enter to create a new album.`) }}
                </v-list-item-title>
              </v-list-item>
            </template>
            <template #chip="chip">
              <v-chip
                :model-value="chip.selected"
                :disabled="chip.disabled"
                prepend-icon="mdi-bookmark"
                class="text-truncate"
                @click:close="removeSelection(chip.index)"
              >
                {{ chip.item.title ? chip.item.title : chip.item }}
              </v-chip>
            </template>
          </v-combobox>
        </v-card-text>
        <v-card-actions class="action-buttons">
          <v-btn variant="flat" color="button" class="action-cancel action-close" @click.stop="close">
            {{ $gettext(`Cancel`) }}
          </v-btn>
          <v-btn
            :disabled="selectedAlbums.length === 0"
            variant="flat"
            color="highlight"
            class="action-confirm text-white"
            @click.stop="confirm"
          >
            {{ $gettext(`Confirm`) }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
<script>
import Album from "model/album";

// TODO: Handle cases where users have more than 10000 albums.
const MaxResults = 10000;

export default {
  name: "PPhotoAlbumDialog",
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["close", "confirm"],
  data() {
    return {
      loading: false,
      selectedAlbums: [],
      albums: [],
      items: [],
      labels: {
        addToAlbum: this.$gettext("Add to album"),
        createAlbum: this.$gettext("Create album"),
      },
    };
  },
  watch: {
    visible: function (show) {
      if (show) {
        this.reset();
        this.load("");
      }
    },
    selectedAlbums: {
      handler(newVal) {
        if (!Array.isArray(newVal)) return;

        let changed = false;
        const processed = [];
        const seenUids = new Set();

        newVal.forEach((item) => {
          // If it's a string, try to match it with existing albums
          if (typeof item === "string" && item.trim().length > 0) {
            const matchedAlbum = this.items.find(
              (album) => album.Title && album.Title.toLowerCase() === item.trim().toLowerCase()
            );

            if (matchedAlbum && !seenUids.has(matchedAlbum.UID)) {
              // Replace string with actual album object
              processed.push(matchedAlbum);
              seenUids.add(matchedAlbum.UID);
              changed = true;
            } else if (!matchedAlbum) {
              // Keep as string for new album creation
              processed.push(item.trim());
            }
          } else if (typeof item === "object" && item?.UID && !seenUids.has(item.UID)) {
            // Keep existing album objects, but prevent duplicates
            processed.push(item);
            seenUids.add(item.UID);
          } else if (typeof item === "object" && item?.UID && seenUids.has(item.UID)) {
            // Skip duplicate album objects
            changed = true;
          }
        });

        // Update selectedAlbums if changes were made
        if (changed || processed.length !== newVal.length) {
          this.$nextTick(() => {
            this.selectedAlbums = processed;
          });
        }
      },
    },
  },
  methods: {
    afterEnter() {
      this.$view.enter(this);
    },
    afterLeave() {
      this.$view.leave(this);
    },
    close() {
      this.$emit("close");
    },
    confirm() {
      if (this.loading) {
        return;
      }

      const existingUids = [];
      const namesToCreate = [];

      (this.selectedAlbums || []).forEach((a) => {
        if (typeof a === "object" && a?.UID) {
          existingUids.push(a.UID);
        } else if (typeof a === "string" && a.length > 0) {
          namesToCreate.push(a);
        }
      });

      // Deduplicate existing UIDs
      const uniqueExistingUids = [...new Set(existingUids)];

      this.loading = true;

      if (namesToCreate.length === 0) {
        this.$emit("confirm", uniqueExistingUids);
        this.loading = false;
        return;
      }

      Promise.all(
        namesToCreate.map((title) => {
          const newAlbum = new Album({ Title: title, UID: "", Favorite: false });
          return newAlbum
            .save()
            .then((a) => a?.UID)
            .catch(() => null);
        })
      )
        .then((created) => {
          const createdUids = created.filter((u) => typeof u === "string" && u.length > 0);
          this.$emit("confirm", [...uniqueExistingUids, ...createdUids]);
        })
        .catch((error) => {
          console.error("Failed to create some albums:", error);
          // Still emit successful ones if any exist
          if (uniqueExistingUids.length > 0) {
            this.$emit("confirm", uniqueExistingUids);
          }
        })
        .finally(() => {
          this.loading = false;
        });
    },
    onLoad() {
      this.loading = true;
      this.$nextTick(() => {
        if (document.activeElement !== this.$refs.input) {
          this.$refs.input.focus();
        }
      });
    },
    onLoaded() {
      this.loading = false;
      this.$nextTick(() => {
        if (document.activeElement !== this.$refs.input) {
          this.$refs.input.focus();
        }
      });
    },
    reset() {
      this.loading = false;
      this.selectedAlbums = [];
      this.albums = [];
      this.items = [];
    },
    removeSelection(index) {
      this.selectedAlbums.splice(index, 1);
    },
    load(q) {
      if (this.loading) {
        return;
      }

      this.onLoad();

      const params = {
        q: q,
        count: MaxResults,
        offset: 0,
        type: "album",
      };

      Album.search(params)
        .then((response) => {
          this.albums = response.models;
          this.items = [...this.albums];
        })
        .finally(() => {
          this.onLoaded();
        });
    },
  },
};
</script>
