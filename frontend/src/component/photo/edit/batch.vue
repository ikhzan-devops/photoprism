<template>
  <v-dialog
    ref="dialog"
    :model-value="visible"
    :fullscreen="$vuetify.display.mdAndDown"
    scrim
    scrollable
    class="p-dialog p-photo-edit-batch v-dialog--sidepanel v-dialog--sidepanel-wide"
    @click.stop="onClick"
  >
    <v-card class="edit-batch__card" ref="content" :tile="$vuetify.display.mdAndDown" tabindex="1">
      <v-toolbar flat color="navigation" :density="$vuetify.display.mdAndDown ? 'compact' : 'comfortable'">
        <v-btn icon class="action-close" @click.stop="onClose">
          <v-icon>mdi-close</v-icon>
        </v-btn>

        <v-toolbar-title>
          {{ title }}
        </v-toolbar-title>
      </v-toolbar>

      <v-row dense :class="!$vuetify.display.mdAndDown ? 'overflow-hidden' : ''">
        <!-- Desktop view -->
        <v-col v-if="!$vuetify.display.mdAndDown" cols="12" lg="4" class="scroll-col">
          <div class="edit-batch photo-results list-view v-table">
            <table>
              <tbody>
                <tr>
                  <td class="col-select" :class="{ 'is-selected': isAllSelected }">
                    <button
                      class="input-select ma-auto"
                      @click.stop.prevent="onSelectAllToggle"
                    >
                      <i class="mdi mdi-checkbox-marked select-on" />
                      <i class="mdi mdi-checkbox-blank-outline select-off" />
                    </button>
                  </td>
                  <td class="media result col-preview">Pictures</td>
                </tr>
                <tr v-for="(m, index) in selectionsFullInfo" :key="m.ID" ref="items" :data-index="index">
                  <td :data-id="m.ID" :data-uid="m.UID" class="col-select" :class="{ 'is-selected': isSelected(m) }">
                    <button
                      class="input-select ma-auto"
                      @touchstart.passive="onMouseDown($event, index)"
                      @touchend.stop="onSelectClick($event, index, true)"
                      @mousedown="onMouseDown($event, index)"
                      @contextmenu.stop="onContextMenu($event, index)"
                      @click.stop.prevent="onSelectClick($event, index, true)"
                    >
                      <i class="mdi mdi-checkbox-marked select-on" />
                      <i class="mdi mdi-checkbox-blank-outline select-off" />
                    </button>
                  </td>
                  <td :data-id="m.ID" :data-uid="m.UID" class="media result col-preview">
<!--                    <div v-if="index < firstVisibleElementIndex || index > lastVisibleElementIndex" class="preview"></div>-->
<!--                    <div-->
<!--                        v-else-->
<!--                        :style="`background-image: url(${m.thumbnailUrl('tile_224')})`"-->
<!--                        class="preview"-->
<!--                        @touchstart.passive="onMouseDown($event, index)"-->
<!--                        @touchend.stop="onClick($event, index, false)"-->
<!--                        @mousedown="onMouseDown($event, index)"-->
<!--                        @contextmenu.stop="onContextMenu($event, index)"-->
<!--                        @click.stop.prevent="onClick($event, index, false)"-->
<!--                    >-->
                    <div
                      :style="`background-image: url(${m.thumbnailUrl('tile_224')})`"
                      class="preview"
                      @touchstart.passive="onMouseDown($event, index)"
                      @touchend.stop="onSelectClick($event, index, false)"
                      @mousedown="onMouseDown($event, index)"
                      @contextmenu.stop="onContextMenu($event, index)"
                      @click.stop.prevent="onSelectClick($event, index, false)"
                    >
                      <div class="preview__overlay"></div>
                      <button
                        v-if="m.Type === 'video' || m.Type === 'live' || m.Type === 'animated'"
                        class="input-open"
                        @click.stop.prevent="openPhoto(index)"
                      >
                        <i v-if="m.Type === 'live'" class="action-live" :title="$gettext('Live')"><icon-live-photo /></i>
                        <i v-else-if="m.Type === 'animated'" class="mdi mdi-file-gif-box" :title="$gettext('Animated')" />
                        <i v-else-if="m.Type === 'video'" class="mdi mdi-play" :title="$gettext('Video')" />
                      </button>
                    </div>
                  </td>
                  <td
                    class="meta-data meta-title col-auto text-start clickable"
                    :title="m.Title"
                    @click.exact="openPhoto(index)"
                  >
                    {{ m.Title }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </v-col>

        <!-- Phone view -->
        <v-col v-else cols="12">
          <div class="edit-batch photo-results list-view v-table">
            <v-expansion-panels
              v-model="expanded"
              variant="accordion"
              density="compact"
              rounded="6"
              tabindex="1"
              class="elevation-0"
            >
              <v-expansion-panel title="Pictures" color="secondary" class="pa-0 elevation-0">
                <v-expansion-panel-text>
                  <table class="w-100">
                    <tbody>
                      <tr v-for="(m, index) in selectionsFullInfo" :key="m.ID" ref="items" :data-index="index">
                        <td :data-id="m.ID" :data-uid="m.UID" class="col-select" :class="{ 'is-selected': isSelected(m) }">
                          <button
                            class="input-select ma-auto"
                            @touchstart.passive="onMouseDown($event, index)"
                            @touchend.stop="onSelectClick($event, index, true)"
                            @mousedown="onMouseDown($event, index)"
                            @contextmenu.stop="onContextMenu($event, index)"
                            @click.stop.prevent="onSelectClick($event, index, true)"
                          >
                            <i class="mdi mdi-checkbox-marked select-on" />
                            <i class="mdi mdi-checkbox-blank-outline select-off" />
                          </button>
                        </td>
                        <td :data-id="m.ID" :data-uid="m.UID" class="media result col-preview">
                          <!--                    <div v-if="index < firstVisibleElementIndex || index > lastVisibleElementIndex" class="preview"></div>-->
                          <!--                    <div-->
                          <!--                        v-else-->
                          <!--                        :style="`background-image: url(${m.thumbnailUrl('tile_224')})`"-->
                          <!--                        class="preview"-->
                          <!--                        @touchstart.passive="onMouseDown($event, index)"-->
                          <!--                        @touchend.stop="onClick($event, index, false)"-->
                          <!--                        @mousedown="onMouseDown($event, index)"-->
                          <!--                        @contextmenu.stop="onContextMenu($event, index)"-->
                          <!--                        @click.stop.prevent="onClick($event, index, false)"-->
                          <!--                    >-->
                          <div
                            :style="`background-image: url(${m.thumbnailUrl('tile_224')})`"
                            class="preview"
                            @touchstart.passive="onMouseDown($event, index)"
                            @touchend.stop="onSelectClick($event, index, false)"
                            @mousedown="onMouseDown($event, index)"
                            @contextmenu.stop="onContextMenu($event, index)"
                            @click.stop.prevent="onSelectClick($event, index, false)"
                          >
                            <div class="preview__overlay"></div>
                            <button
                              v-if="m.Type === 'video' || m.Type === 'live' || m.Type === 'animated'"
                              class="input-open"
                              @click.stop.prevent="openPhoto(index)"
                            >
                              <i v-if="m.Type === 'live'" class="action-live" :title="$gettext('Live')"><icon-live-photo /></i>
                              <i v-else-if="m.Type === 'animated'" class="mdi mdi-file-gif-box" :title="$gettext('Animated')" />
                              <i v-else-if="m.Type === 'video'" class="mdi mdi-play" :title="$gettext('Video')" />
                            </button>
                          </div>
                        </td>
                        <td
                          class="meta-data meta-title col-auto text-start clickable"
                          :title="m.Title"
                          @click.exact="openPhoto(index)"
                        >
                          {{ m.Title }}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </div>
        </v-col>

        <v-col cols="12" lg="8">
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
                        :label="$gettext('Title')"
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
                    <v-col cols="6" md="2">
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
                    <v-col cols="6" md="3">
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
                    <v-col cols="12" sm="6" md="3">
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
                    <v-col cols="12" sm="6" md="4">
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
                    <v-col cols="12" sm="6" md="3">
                      <v-autocomplete
                        :placeholder="$gettext('Country')"
                        hide-details
                        hide-no-data
                        autocomplete="off"
                        item-value="Code"
                        item-title="Name"
                        :items="countries"
                        density="comfortable"
                        validate-on="input"
                        class="input-country"
                      >
                      </v-autocomplete>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
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
                    <v-col cols="12" sm="6" md="3">
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
                    <v-col cols="12" sm="6" md="3">
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
                  :disabled="this.selection.length < 1"
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
import * as options from "options/options";
import countries from "options/countries.json";
import IconLivePhoto from "../../icon/live-photo.vue";
import { PhotoClipboard } from "common/clipboard";
import { Photo } from "model/photo";
import Thumb from "../../../model/thumb";

export default {
  name: "PPhotoEditBatch",
  components: {IconLivePhoto},
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    selection: {
      type: Array,
      default: () => [],
    },
    openDate: {
      type: Function,
      default: () => {},
    },
    openLocation: {
      type: Function,
      default: () => {},
    },
    editPhoto: {
      type: Function,
      default: () => {},
    },
  },
  emits: ["close"],
  data() {
    return {
      model: new Photo(),
      uid: "",
      loading: false,
      subscriptions: [],

      selectionsFullInfo: [],
      expanded: [0],
      isBatchDialog: true,
      isAllSelected: true,
      options,
      countries,
      albums: [],
      firstVisibleElementIndex: 0,
      lastVisibleElementIndex: 0,
      mouseDown: {
        index: -1,
        scrollY: window.scrollY,
        timeStamp: -1,
      },
    };
  },
  computed: {
    title() {
      if (this.selection.length > 0) {
        return this.$gettext(`Batch Edit (${this.selection.length})`);
      }

      return this.$gettext(`Batch Edit (${this.selection.length})`);
    },
  },
  watch: {
    visible: function (show) {
      if (show) {
        this.expanded = [];
      }
    },
    selection: function () {
      this.getSelectionsFullInfo();
    },
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
    getSelectionsFullInfo() {
      const photosUids = this.selection.join('|');
      const photosLength = photosUids.length;

      const params = {
        count: photosLength,
        uid: photosUids,
      };

      Photo.search(params)
        .then((response) => {
          this.selectionsFullInfo = response.models;
        });
    },
    openPhoto(index) {
      this.$lightbox.openModels(Thumb.fromPhotos([this.selectionsFullInfo[index]]), 0, null , this.isBatchDialog);
    },
    isSelected(m) {
      return PhotoClipboard.has(m);
    },
    onClick(ev) {
      // Closes dialog when user clicks on background and model data is unchanged.
      if (!ev || !ev?.target?.classList?.contains("v-overlay__scrim")) {
        return;
      }
      ev.preventDefault();
      this.onClose();
    },
    onSelectAllToggle() {
      PhotoClipboard.toggleAllIds(this.selectionsFullInfo);
      this.isAllSelected = !this.isAllSelected;
    },
    onUpdate() {
      console.log('Add event on update');
    },
    onSelectClick(ev, index, select) {
      const longClick = this.mouseDown.index === index && ev.timeStamp - this.mouseDown.timeStamp > 400;
      const scrolled = this.mouseDown.scrollY - window.scrollY !== 0;

      if (!select && scrolled) {
        return;
      }

      ev.preventDefault();
      ev.stopPropagation();

      if (select !== false && (select || longClick || this.selectMode)) {
        if (longClick || ev.shiftKey) {
          this.selectRange(index);
        } else {
          this.toggle(this.selectionsFullInfo[index]);
        }
      } else if (this.selectionsFullInfo[index]) {
        this.openPhoto(index);
      }
    },
    toggle(photo) {
      this.$clipboard.toggle(photo);
    },
    onMouseDown(ev, index) {
      this.mouseDown.index = index;
      this.mouseDown.scrollY = window.scrollY;
      this.mouseDown.timeStamp = ev.timeStamp;
    },
    onContextMenu(ev, index) {
      if (this.$isMobile) {
        ev.preventDefault();
        ev.stopPropagation();
        this.selectRange(index);
      }
    },
    onClose() {
      // Closes the dialog only if model data is unchanged.
      if (this.model?.hasId() && this.model?.wasChanged()) {
        this.$refs?.dialog?.animateClick();
      } else {
        this.close();
      }
    },
    selectRange(index) {
      this.$clipboard.addRange(index, this.selectionsFullInfo);
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
  },
};
</script>
