<template>
  <v-dialog
    ref="dialog"
    :model-value="visible"
    :fullscreen="$vuetify.display.smAndDown"
    scrim
    scrollable
    class="p-dialog p-photo-edit-batch v-dialog--sidepanel"
    @click.stop="onClick"
  >
    <v-card ref="content" :tile="$vuetify.display.smAndDown" tabindex="1">
      <v-toolbar flat color="navigation" :density="$vuetify.display.smAndDown ? 'compact' : 'comfortable'">
        <v-btn icon class="action-close" @click.stop="onClose">
          <v-icon>mdi-close</v-icon>
        </v-btn>

        <v-toolbar-title>
          {{ title }}
        </v-toolbar-title>
      </v-toolbar>

      <v-row dense>
        <v-col cols="12 d-none d-md-flex" md="4">
          <p>First</p>
        </v-col>

        <v-col cols="12" md="8">
          <v-form
            ref="form"
            validate-on="invalid-input"
            class="p-form p-form-photo-details-meta"
            accept-charset="UTF-8"
            tabindex="1"
            @submit.prevent="save"
          >
            <div class="form-body">
              <div class="form-controls">
                <div>
                  <v-col cols="12" md="6">
                    <p>Description</p>
                  </v-col>
                  <v-row dense>
                    <v-col cols="12" md="6">
                      <v-text-field
                        hide-details
                        :label="$gettext('Photo', 'Title')"
                        placeholder=""
                        autocomplete="off"
                        density="comfortable"
                        class="input-title"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-textarea
                        hide-details
                        autocomplete="off"
                        auto-grow
                        :label="$gettext('Subject')"
                        placeholder=""
                        :rows="1"
                        density="comfortable"
                        class="input-subject"
                      ></v-textarea>
                    </v-col>
                  </v-row>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-textarea
                      hide-details
                      autocomplete="off"
                      auto-grow
                      :label="$gettext('Caption')"
                      placeholder=""
                      :rows="1"
                      density="comfortable"
                      class="input-caption"
                    ></v-textarea>
                  </v-col>
                </div>

                <div>
                  <v-col cols="12" md="6">
                    <p>Date</p>
                  </v-col>
                  <v-row dense>
                    <v-col cols="3" lg="2">
                      <v-combobox
                        :label="$gettext('Day')"
                        :placeholder="$gettext('Unknown')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :items="options.Days()"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-day"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="3" lg="3">
                      <v-combobox
                        :label="$gettext('Month')"
                        :placeholder="$gettext('Unknown')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :items="options.MonthsShort()"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-month"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="3" lg="3">
                      <v-combobox
                        :label="$gettext('Year')"
                        :placeholder="$gettext('Unknown')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :items="options.Years(1900)"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-year"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="3" lg="4">
                      <v-autocomplete
                        :label="$gettext('Time Zone')"
                        hide-no-data
                        :items="options.TimeZones()"
                        item-value="ID"
                        item-title="Name"
                        density="comfortable"
                        class="input-timezone"
                      ></v-autocomplete>
                    </v-col>
                  </v-row>
                </div>

                <div>
                  <v-col cols="12" md="6">
                    <p>Location</p>
                  </v-col>
                  <v-row dense>
                    <v-col cols="3" lg="3">
                      <v-autocomplete
                        :placeholder="$gettext('Country')"
                        hide-details
                        hide-no-data
                        autocomplete="off"
                        item-value="Code"
                        item-title="Name"
                        :items="countries"
                        prepend-inner-icon="mdi-map-marker"
                        density="comfortable"
                        validate-on="input"
                        class="input-country"
                      >
                      </v-autocomplete>
                    </v-col>
                    <v-col cols="3" lg="3">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        autocorrect="off"
                        autocapitalize="none"
                        :label="$gettext('Latitude')"
                        placeholder=""
                        density="comfortable"
                        validate-on="input"
                        class="input-latitude"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="3" lg="3">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        autocorrect="off"
                        autocapitalize="none"
                        :label="$gettext('Longitude')"
                        placeholder=""
                        density="comfortable"
                        validate-on="input"
                        class="input-longitude"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="3" lg="3">
                      <v-text-field
                        hide-details
                        flat
                        autocomplete="off"
                        autocorrect="off"
                        autocapitalize="none"
                        :label="$gettext('Altitude (m)')"
                        placeholder=""
                        color="surface-variant"
                        density="comfortable"
                        validate-on="input"
                        class="input-altitude"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                </div>

                <div>
                  <v-col cols="12" md="6">
                    <p>Copyright</p>
                  </v-col>
                  <v-row dense>
                    <v-col cols="12" md="6">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        :label="$gettext('Artist')"
                        placeholder=""
                        density="comfortable"
                        class="input-artist"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        :label="$gettext('Copyright')"
                        placeholder=""
                        density="comfortable"
                        class="input-copyright"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-textarea
                      hide-details
                      autocomplete="off"
                      auto-grow
                      :label="$gettext('License')"
                      placeholder=""
                      :rows="1"
                      density="comfortable"
                      class="input-license"
                    ></v-textarea>
                  </v-col>
                </div>

                <div>
                  <v-col cols="12" md="6">
                    <p>Albums</p>
                  </v-col>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-combobox
                      hide-details
                      chips
                      closable-chips
                      multiple
                      class="input-albums"
                      :items="albums"
                      item-title="Title"
                      item-value="UID"
                      :placeholder="$gettext('Select or create an album')"
                      return-object
                    >
                      <template #no-data>
                        <v-list-item>
                          <v-list-item-title>
                            {{ $gettext(`Press enter to create a new album.`) }}
                          </v-list-item-title>
                        </v-list-item>
                      </template>
<!--                      <template #chip="chip">-->
<!--                        <v-chip-->
<!--                          prepend-icon="mdi-bookmark"-->
<!--                          class="text-truncate"-->
<!--                        >-->
<!--                          {{ chip.item.title ? chip.item.title : chip.item }}-->
<!--                        </v-chip>-->
<!--                      </template>-->
                    </v-combobox>
                  </v-col>
                </div>

                <div>
                  <v-col cols="12" md="6">
                    <p>Labels</p>
                  </v-col>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-combobox
                      rows="2"
                      hide-details
                      chips
                      closable-chips
                      multiple
                      class="input-albums"
                      :items="albums"
                      item-title="Title"
                      item-value="UID"
                      :placeholder="$gettext('Select or create a label')"
                      return-object
                    >
                      <template #no-data>
                        <v-list-item>
                          <v-list-item-title>
                            {{ $gettext(`Press enter to create a new label.`) }}
                          </v-list-item-title>
                        </v-list-item>
                      </template>
<!--                      <template #chip="chip">-->
<!--                        <v-chip-->
<!--                          prepend-icon="mdi-bookmark"-->
<!--                          class="text-truncate"-->
<!--                        >-->
<!--                          {{ chip.item.title ? chip.item.title : chip.item }}-->
<!--                        </v-chip>-->
<!--                      </template>-->
                    </v-combobox>
                  </v-col>
                </div>
              </div>
            </div>

            <div class="form-actions form-actions--sticky">
              <div class="action-buttons">
                <v-btn color="button" variant="flat" class="action-close" @click.stop="close">
                  {{ $gettext(`Close`) }}
                </v-btn>
                <v-btn
                  color="highlight"
                  variant="flat"
                  class="action-apply action-approve"
                  @click.stop="save(false)"
                >
                  <span>{{ $gettext(`Apply`) }}</span>
                </v-btn>
              </div>
            </div>
          </v-form>
        </v-col>
      </v-row>
    </v-card>
  </v-dialog>
</template>
<script>
import Photo from "model/photo";
import * as options from "options/options";
import countries from "options/countries.json";

export default {
  name: "PPhotoEditBatch",
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    selection: {
      type: Array,
      default: () => [],
    },
  },
  emits: ["close"],
  data() {
    return {
      model: new Photo(),
      uid: "",
      loading: false,
      subscriptions: [],

      selections: [],
      view: this.$view.data(),
      options,
      countries,
      albums: [],
      selectedAlbums: [],
    };
  },
  computed: {
    title() {
      if (this.selection.length > 0) {
        return this.$gettext(`Batch Edit (${this.selection.length})`);
      }

      return this.title;
    },
  },
  watch: {
    visible: function (show) {
      if (show) {
        this.find(this.index);

        // Set currently selected albums.
        // if (this.data && Array.isArray(this.data.albums)) {
        //   this.selectedAlbums = this.data.albums;
        // } else {
        //   this.selectedAlbums = [];
        // }
      }
    },
    // selection: function () {
    //   this.getSelection();
    //   console.log('this.$view.data()', this.$view.data());
    // },
  },
  created() {
    this.subscriptions.push(this.$event.subscribe("photos.updated", (ev, data) => this.onUpdate(ev, data)));
  },
  beforeUnmount() {
    for (let i = 0; i < this.subscriptions.length; i++) {
      this.$event.unsubscribe(this.subscriptions[i]);
    }
  },
  methods: {
    onClick(ev) {
      // Closes dialog when user clicks on background and model data is unchanged.
      if (!ev || !ev?.target?.classList?.contains("v-overlay__scrim")) {
        return;
      }
      ev.preventDefault();
      this.onClose();
    },
    onClose() {
      // Closes the dialog only if model data is unchanged.
      if (this.model?.hasId() && this.model?.wasChanged()) {
        this.$refs?.dialog?.animateClick();
      } else {
        this.close();
      }
    },
    close() {
      // Close the dialog.
      this.$emit("close");
    },
    save(close) {
      // if (this.invalidDate) {
      //   this.$notify.error(this.$gettext("Invalid date"));
      //   return;
      // }
      //
      // this.updateModel();
      //
      // this.view.model.update().then(() => {
      if (close) {
        this.$emit("close");
      }
      //
      //   this.syncTime();
      // });
    },
    find(index) {
      if (this.loading) {
        return Promise.reject();
      }

      if (!this.selection) {
        this.$notify.error(this.$gettext("Invalid photos selected"));
        return Promise.reject();
      }

      this.loading = true;
      this.selected = index;
      this.selectedId = this.selection[index];
      this.loading = false;

      // const a = this.model
      //   .find(this.selectedId)
      //   .then((model) => {
      //     model.refreshFileAttr();
      //     this.model = model;
      //     this.loading = false;
      //     this.uid = this.selectedId;
      //   })
      //   .catch(() => (this.loading = false));
      // console.log('a', a);
      // return a;
    },
    // getImageData() {
    //   this.model.find(0).then((model) => {
    //     model.refreshFileAttr();
    //     this.model = model;
    //   });
    // },
    // getSelection() {
    //   if (!this.selection) {
    //     return Promise.reject();
    //   }
    //
    //   console.log("getImageData", this.getImageData());
    //   // this.selections = this.selection.find((val, index) => this.find(index));
    //   // console.log('this.selections', this.selections);
    // },
  },
};
</script>
