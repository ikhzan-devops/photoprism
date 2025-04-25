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
                <tr class="pa-3">
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
                <tr v-for="(m, index) in selectedPhotos" :key="m.ID" ref="items" :data-index="index">
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
                      <tr v-for="(m, index) in selectedPhotos" :key="m.ID" ref="items" :data-index="index">
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

        <v-col cols="12" lg="8" class="scroll-col">
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
                    <v-col cols="6" md="3">
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
                    <v-col cols="6" md="4">
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

      selections: [],
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
      selectedPhotos: [
        new Photo({
          "ID": 105,
          "TakenAt": "2022-09-14T21:32:38Z",
          "TakenAtLocal": "2022-09-14T21:32:38Z",
          "TakenSrc": "",
          "UID": "psik3xv214ri31l8",
          "Type": "image",
          "TypeSrc": "",
          "Title": "",
          "TitleSrc": "",
          "Caption": "",
          "CaptionSrc": "",
          "Path": "demo/demo",
          "Name": "Fox",
          "OriginalName": "",
          "Stack": 0,
          "Favorite": true,
          "Private": false,
          "Scan": false,
          "Panorama": false,
          "TimeZone": "UTC",
          "PlaceID": "zz",
          "PlaceSrc": "",
          "CellID": "zz",
          "CellAccuracy": 0,
          "Altitude": 0,
          "Lat": 0,
          "Lng": 0,
          "Country": "zz",
          "Year": 2022,
          "Month": 9,
          "Day": 14,
          "Iso": 0,
          "Exposure": "",
          "FNumber": 0,
          "FocalLength": 0,
          "Quality": 4,
          "Resolution": 1,
          "Color": 6,
          "CameraID": 1,
          "CameraSerial": "",
          "CameraSrc": "",
          "LensID": 1,
          "Details": {
            "PhotoID": 105,
            "Keywords": "blue, demo, fox, outdoor",
            "KeywordsSrc": "",
            "Notes": "",
            "NotesSrc": "",
            "Subject": "",
            "SubjectSrc": "",
            "Artist": "",
            "ArtistSrc": "",
            "Copyright": "",
            "CopyrightSrc": "",
            "License": "",
            "LicenseSrc": "",
            "Software": "",
            "SoftwareSrc": "",
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2025-02-07T10:05:49Z"
          },
          "Camera": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown"
          },
          "Lens": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown",
            "Type": ""
          },
          "Cell": {
            "ID": "zz",
            "Name": "",
            "Street": "",
            "Postcode": "",
            "Category": "",
            "Place": {
              "PlaceID": "zz",
              "Label": "Unknown",
              "District": "Unknown",
              "City": "Unknown",
              "State": "Unknown",
              "Country": "zz",
              "Keywords": "",
              "Favorite": false,
              "PhotoCount": 65,
              "CreatedAt": "2024-12-16T17:00:50Z",
              "UpdatedAt": "2024-12-16T17:00:50Z"
            },
            "CreatedAt": "2024-12-16T17:00:50Z",
            "UpdatedAt": "2024-12-16T17:00:50Z"
          },
          "Place": {
            "PlaceID": "zz",
            "Label": "Unknown",
            "District": "Unknown",
            "City": "Unknown",
            "State": "Unknown",
            "Country": "zz",
            "Keywords": "",
            "Favorite": false,
            "PhotoCount": 65,
            "CreatedAt": "2024-12-16T17:00:50Z",
            "UpdatedAt": "2024-12-16T17:00:50Z"
          },
          "Albums": [],
          "Files": [
            {
              "UID": "fsoljb5rl78xzafx",
              "PhotoUID": "psik3xv214ri31l8",
              "Name": "demo/demo/Fox.avif.jpg",
              "Root": "sidecar",
              "Hash": "147a9632feb0c6ad84ba347eefd61e874372eeab",
              "Size": 153277,
              "Primary": true,
              "TimeIndex": "79779085786762-9999999895-0-fsoljb5rl78xzafx",
              "MediaID": "9999999895-0-fsoljb5rl78xzafx",
              "MediaUTC": 1663191158000,
              "FileType": "jpg",
              "MediaType": "image",
              "Mime": "image/jpeg",
              "Width": 1204,
              "Height": 800,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1.51,
              "MainColor": "blue",
              "Colors": "611661160",
              "Diff": 1023,
              "Chroma": 8,
              "ModTime": 1724224243,
              "CreatedAt": "2024-12-16T17:01:53Z",
              "CreatedIn": 593690959,
              "UpdatedAt": "2025-01-14T16:40:43Z",
              "UpdatedIn": 608092958,
              "Markers": []
            },
            {
              "UID": "fsoljb47oqu1lyp3",
              "PhotoUID": "psik3xv214ri31l8",
              "Name": "demo/demo/Fox.avif",
              "Root": "/",
              "Hash": "0867c29e938e03d7130fbdffde2c9ce31cd536fb",
              "Size": 80743,
              "Primary": false,
              "TimeIndex": "79779085786762-9999999895-1-fsoljb47oqu1lyp3",
              "MediaID": "9999999895-1-fsoljb47oqu1lyp3",
              "Codec": "avif",
              "FileType": "avif",
              "MediaType": "image",
              "Mime": "image/avif",
              "Width": 1204,
              "Height": 800,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1.51,
              "Diff": -1,
              "Chroma": -1,
              "ModTime": 1663191158,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "CreatedIn": 13990667,
              "UpdatedAt": "2025-01-14T16:40:43Z",
              "UpdatedIn": 19678625,
              "Markers": []
            }
          ],
          "Labels": [
            {
              "PhotoID": 105,
              "LabelID": 13,
              "LabelSrc": "image",
              "Uncertainty": 52,
              "Photo": null,
              "Label": {
                "ID": 13,
                "UID": "lsoljaq4fh7vszjs",
                "Slug": "outdoor",
                "CustomSlug": "outdoor",
                "Name": "Outdoor",
                "Priority": -1,
                "Favorite": false,
                "Description": "",
                "Notes": "",
                "PhotoCount": 3,
                "Thumb": "147a9632feb0c6ad84ba347eefd61e874372eeab",
                "CreatedAt": "2024-12-16T17:01:38Z",
                "UpdatedAt": "2025-02-07T10:05:49Z"
              }
            }
          ],
          "CreatedAt": "2024-08-21T07:10:43Z",
          "UpdatedAt": "2025-02-07T10:05:49Z",
          "EditedAt": null,
          "CheckedAt": "2025-04-07T18:02:13Z",
          "DeletedAt": null
        }),
        new Photo({
          "ID": 102,
          "TakenAt": "2022-05-15T03:03:32Z",
          "TakenAtLocal": "2022-05-15T03:03:32Z",
          "TakenSrc": "manual",
          "UID": "psik3xvfp8pq05nm",
          "Type": "image",
          "TypeSrc": "manual",
          "Title": "Animated Earth / 2022",
          "TitleSrc": "",
          "Caption": "",
          "CaptionSrc": "",
          "Path": "demo/demo/Animated",
          "Name": "animated-earth",
          "OriginalName": "",
          "Stack": -1,
          "Favorite": true,
          "Private": false,
          "Scan": false,
          "Panorama": false,
          "TimeZone": "UTC",
          "PlaceID": "zz",
          "PlaceSrc": "",
          "CellID": "zz",
          "CellAccuracy": 0,
          "Altitude": 0,
          "Lat": 0,
          "Lng": 0,
          "Country": "zz",
          "Year": 2022,
          "Month": 5,
          "Day": 15,
          "Iso": 0,
          "Exposure": "",
          "FNumber": 0,
          "FocalLength": 0,
          "Quality": 5,
          "Resolution": 0,
          "Color": 0,
          "CameraID": 1,
          "CameraSerial": "",
          "CameraSrc": "",
          "LensID": 1,
          "Details": {
            "PhotoID": 102,
            "Keywords": "animated, black, demo, earth",
            "KeywordsSrc": "manual",
            "Notes": "Notes",
            "NotesSrc": "manual",
            "Subject": "",
            "SubjectSrc": "",
            "Artist": "",
            "ArtistSrc": "",
            "Copyright": "",
            "CopyrightSrc": "manual",
            "License": "License",
            "LicenseSrc": "manual",
            "Software": "",
            "SoftwareSrc": "",
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2025-02-04T22:23:43Z"
          },
          "Camera": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown"
          },
          "Lens": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown",
            "Type": ""
          },
          "Cell": {
            "ID": "zz",
            "Name": "",
            "Street": "",
            "Postcode": "",
            "Category": "",
            "Place": {
              "PlaceID": "zz",
              "Label": "Unknown",
              "District": "Unknown",
              "City": "Unknown",
              "State": "Unknown",
              "Country": "zz",
              "Keywords": "",
              "Favorite": false,
              "PhotoCount": 65,
              "CreatedAt": "2024-12-16T17:00:50Z",
              "UpdatedAt": "2024-12-16T17:00:50Z"
            },
            "CreatedAt": "2024-12-16T17:00:50Z",
            "UpdatedAt": "2024-12-16T17:00:50Z"
          },
          "Place": {
            "PlaceID": "zz",
            "Label": "Unknown",
            "District": "Unknown",
            "City": "Unknown",
            "State": "Unknown",
            "Country": "zz",
            "Keywords": "",
            "Favorite": false,
            "PhotoCount": 65,
            "CreatedAt": "2024-12-16T17:00:50Z",
            "UpdatedAt": "2024-12-16T17:00:50Z"
          },
          "Albums": [],
          "Files": [
            {
              "UID": "fsoljb4wh92twxcp",
              "PhotoUID": "psik3xvfp8pq05nm",
              "Name": "demo/demo/Animated/animated-earth.gif.jpg",
              "Root": "sidecar",
              "Hash": "a5534f864efdf58f9f867f2d6c3d4f1fbe88d9e4",
              "Size": 12271,
              "Primary": true,
              "TimeIndex": "79779484969668-9999999898-0-fsoljb4wh92twxcp",
              "MediaID": "9999999898-0-fsoljb4wh92twxcp",
              "MediaUTC": 1649991812000,
              "FileType": "jpg",
              "MediaType": "image",
              "Mime": "image/jpeg",
              "Width": 300,
              "Height": 300,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1,
              "MainColor": "black",
              "Colors": "000100100",
              "Diff": 703,
              "Chroma": 4,
              "ModTime": 1738707439,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "CreatedIn": 172989792,
              "UpdatedAt": "2025-02-04T22:17:19Z",
              "UpdatedIn": 130020958,
              "Markers": []
            }
          ],
          "Labels": [],
          "CreatedAt": "2024-08-21T07:10:43Z",
          "UpdatedAt": "2025-02-04T22:23:43Z",
          "EditedAt": "2025-02-04T22:23:43Z",
          "CheckedAt": "2025-04-07T18:02:13Z",
          "DeletedAt": null
        }),
        new Photo({
          "ID": 100,
          "TakenAt": "2020-08-31T16:00:54Z",
          "TakenAtLocal": "2020-08-31T18:00:54Z",
          "TakenSrc": "meta",
          "UID": "pqowera1e1siaw6z",
          "Type": "image",
          "TypeSrc": "",
          "Title": "Bench / Berlin / 2020",
          "TitleSrc": "",
          "Caption": "",
          "CaptionSrc": "",
          "Path": "demo/demo",
          "Name": "20200831-180054-Bench-Berlin-2020-2wy",
          "OriginalName": "",
          "Stack": 0,
          "Favorite": false,
          "Private": false,
          "Scan": false,
          "Panorama": false,
          "TimeZone": "Europe/Berlin",
          "PlaceID": "de:z0vP8a5RZU2e",
          "PlaceSrc": "meta",
          "CellID": "s2:47a85a6214b4",
          "CellAccuracy": 0,
          "Altitude": 67,
          "Lat": 52.45218658444445,
          "Lng": 13.309621810833335,
          "Country": "de",
          "Year": 2020,
          "Month": 8,
          "Day": 31,
          "Iso": 50,
          "Exposure": "1/316",
          "FNumber": 1.8,
          "FocalLength": 27,
          "Quality": 4,
          "Resolution": 10,
          "Color": 9,
          "CameraID": 10,
          "CameraSerial": "",
          "CameraSrc": "meta",
          "LensID": 1,
          "Details": {
            "PhotoID": 100,
            "Keywords": "attraction, bauern, bench, berlin, botanical, botanischer, demo, garden, garten, germany, green, hortensienplatz, lichterfelde, nutzpflanzen, tourist",
            "KeywordsSrc": "",
            "Notes": "",
            "NotesSrc": "",
            "Subject": "",
            "SubjectSrc": "",
            "Artist": "",
            "ArtistSrc": "",
            "Copyright": "",
            "CopyrightSrc": "",
            "License": "",
            "LicenseSrc": "",
            "Software": "ELE-L29 10.1.0.140(C431E22R2P5)",
            "SoftwareSrc": "meta",
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2025-01-22T23:10:00Z"
          },
          "Camera": {
            "ID": 10,
            "Slug": "huawei-p30",
            "Name": "HUAWEI P30",
            "Make": "HUAWEI",
            "Model": "P30",
            "Type": "phone"
          },
          "Lens": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown",
            "Type": ""
          },
          "Cell": {
            "ID": "s2:47a85a6214b4",
            "Name": "Botanischer Garten",
            "Street": "Hortensienplatz",
            "Postcode": "12203",
            "Category": "tourist attraction",
            "Place": {
              "PlaceID": "de:z0vP8a5RZU2e",
              "Label": "Lichterfelde, Berlin, Germany",
              "District": "Lichterfelde",
              "City": "Berlin",
              "State": "Berlin",
              "Country": "de",
              "Keywords": "",
              "Favorite": false,
              "PhotoCount": 6,
              "CreatedAt": "2024-12-16T17:01:40Z",
              "UpdatedAt": "2024-12-16T17:01:52Z"
            },
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2024-12-16T17:01:52Z"
          },
          "Place": {
            "PlaceID": "de:z0vP8a5RZU2e",
            "Label": "Lichterfelde, Berlin, Germany",
            "District": "Lichterfelde",
            "City": "Berlin",
            "State": "Berlin",
            "Country": "de",
            "Keywords": "",
            "Favorite": false,
            "PhotoCount": 6,
            "CreatedAt": "2024-12-16T17:01:40Z",
            "UpdatedAt": "2024-12-16T17:01:52Z"
          },
          "Albums": [],
          "Files": [
            {
              "UID": "fsoljb48qx2t8xgy",
              "PhotoUID": "pqowera1e1siaw6z",
              "Name": "demo/demo/20200831-180054-Bench-Berlin-2020-2wy.yml",
              "Root": "/",
              "Hash": "a8c891a22823718dbb915100e6b1d1430d212a7e",
              "Size": 489,
              "Primary": false,
              "TimeIndex": "79799168819946-9999999900-2-fsoljb48qx2t8xgy",
              "MediaID": "9999999900-2-fsoljb48qx2t8xgy",
              "MediaUTC": 1598896854000,
              "FileType": "yml",
              "MediaType": "sidecar",
              "Mime": "text/plain",
              "Sidecar": true,
              "Diff": -1,
              "Chroma": -1,
              "ModTime": 1613943478,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "CreatedIn": 9475042,
              "UpdatedAt": "2025-01-14T16:40:43Z",
              "UpdatedIn": 16244292,
              "Markers": []
            },
            {
              "UID": "fsoljb4lzjhnqmr2",
              "PhotoUID": "pqowera1e1siaw6z",
              "Name": "demo/demo/20200831-180054-Bench-Berlin-2020-2wy.jpg",
              "Root": "/",
              "Hash": "da7ba4f7c21e5902976cdf0eafb7d744183e4852",
              "Size": 6085448,
              "Primary": true,
              "TimeIndex": "79799168819946-9999999900-0-fsoljb4lzjhnqmr2",
              "MediaID": "9999999900-0-fsoljb4lzjhnqmr2",
              "MediaUTC": 1598889654000,
              "Codec": "jpeg",
              "FileType": "jpg",
              "MediaType": "image",
              "Mime": "image/jpeg",
              "Width": 3648,
              "Height": 2736,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1.33,
              "MainColor": "green",
              "Colors": "999921911",
              "Diff": 1007,
              "Chroma": 19,
              "Software": "ELE-L29 10.1.0.140(C431E22R2P5)",
              "ModTime": 1598943367,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "CreatedIn": 656223000,
              "UpdatedAt": "2025-01-14T16:40:43Z",
              "UpdatedIn": 642219042,
              "Markers": []
            }
          ],
          "Labels": [
            {
              "PhotoID": 100,
              "LabelID": 24,
              "LabelSrc": "location",
              "Uncertainty": 0,
              "Photo": null,
              "Label": {
                "ID": 24,
                "UID": "lsoljasrdnnyslnb",
                "Slug": "tourist-attraction",
                "CustomSlug": "tourist-attraction",
                "Name": "Tourist Attraction",
                "Priority": -1,
                "Favorite": false,
                "Description": "",
                "Notes": "",
                "PhotoCount": 9,
                "Thumb": "da7ba4f7c21e5902976cdf0eafb7d744183e4852",
                "CreatedAt": "2024-12-16T17:01:40Z",
                "UpdatedAt": "2025-01-22T23:10:00Z"
              }
            },
            {
              "PhotoID": 100,
              "LabelID": 59,
              "LabelSrc": "image",
              "Uncertainty": 17,
              "Photo": null,
              "Label": {
                "ID": 59,
                "UID": "lsoljb487i972ybl",
                "Slug": "bench",
                "CustomSlug": "bench",
                "Name": "Bench",
                "Priority": 0,
                "Favorite": false,
                "Description": "",
                "Notes": "",
                "PhotoCount": 1,
                "Thumb": "da7ba4f7c21e5902976cdf0eafb7d744183e4852",
                "CreatedAt": "2024-12-16T17:01:52Z",
                "UpdatedAt": "2025-01-22T23:10:00Z"
              }
            }
          ],
          "CreatedAt": "2021-02-21T21:37:58Z",
          "UpdatedAt": "2025-01-22T23:10:00Z",
          "EditedAt": "2025-01-22T23:10:00Z",
          "CheckedAt": "2025-04-07T18:02:13Z",
          "DeletedAt": null
        }),
        new Photo({
          "ID": 99,
          "TakenAt": "2020-07-27T11:55:07Z",
          "TakenAtLocal": "2020-07-27T13:55:07Z",
          "TakenSrc": "meta",
          "UID": "pqower816p8ztqdu",
          "Type": "image",
          "TypeSrc": "",
          "Title": "Peacock / Tübingen / 2020",
          "TitleSrc": "",
          "Caption": "",
          "CaptionSrc": "",
          "Path": "demo/demo",
          "Name": "20200727-135507-Peacock-Tubingen-2020-2gu",
          "OriginalName": "",
          "Stack": 0,
          "Favorite": true,
          "Private": false,
          "Scan": false,
          "Panorama": false,
          "TimeZone": "Europe/Berlin",
          "PlaceID": "de:Nm0DsLKL5cyn",
          "PlaceSrc": "meta",
          "CellID": "s2:4799fb37ea3c",
          "CellAccuracy": 0,
          "Altitude": 417,
          "Lat": 48.51729965194444,
          "Lng": 9.023340225,
          "Country": "de",
          "Year": 2020,
          "Month": 7,
          "Day": 27,
          "Iso": 50,
          "Exposure": "1/50",
          "FNumber": 2.4,
          "FocalLength": 81,
          "Quality": 7,
          "Resolution": 10,
          "Color": 1,
          "CameraID": 10,
          "CameraSerial": "",
          "CameraSrc": "meta",
          "LensID": 1,
          "Details": {
            "PhotoID": 99,
            "Keywords": "baden-württemberg, demo, germany, green, grey, peacock, schwärzloch, tubingen, tübingen, weststadt",
            "KeywordsSrc": "manual",
            "Notes": "",
            "NotesSrc": "",
            "Subject": "",
            "SubjectSrc": "",
            "Artist": "",
            "ArtistSrc": "",
            "Copyright": "Copyright",
            "CopyrightSrc": "manual",
            "License": "",
            "LicenseSrc": "",
            "Software": "ELE-L29 10.1.0.133(C431E22R2P5)",
            "SoftwareSrc": "meta",
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2025-02-04T22:24:54Z"
          },
          "Camera": {
            "ID": 10,
            "Slug": "huawei-p30",
            "Name": "HUAWEI P30",
            "Make": "HUAWEI",
            "Model": "P30",
            "Type": "phone"
          },
          "Lens": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown",
            "Type": ""
          },
          "Cell": {
            "ID": "s2:4799fb37ea3c",
            "Name": "",
            "Street": "Schwärzloch",
            "Postcode": "72070",
            "Category": "",
            "Place": {
              "PlaceID": "de:Nm0DsLKL5cyn",
              "Label": "Tübingen, Baden-Württemberg, Germany",
              "District": "Weststadt",
              "City": "Tübingen",
              "State": "Baden-Württemberg",
              "Country": "de",
              "Keywords": "",
              "Favorite": false,
              "PhotoCount": 1,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "UpdatedAt": "2024-12-16T17:01:52Z"
            },
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2024-12-16T17:01:52Z"
          },
          "Place": {
            "PlaceID": "de:Nm0DsLKL5cyn",
            "Label": "Tübingen, Baden-Württemberg, Germany",
            "District": "Weststadt",
            "City": "Tübingen",
            "State": "Baden-Württemberg",
            "Country": "de",
            "Keywords": "",
            "Favorite": false,
            "PhotoCount": 1,
            "CreatedAt": "2024-12-16T17:01:52Z",
            "UpdatedAt": "2024-12-16T17:01:52Z"
          },
          "Albums": [],
          "Files": [
            {
              "UID": "fsoljb4g9jt4dx5i",
              "PhotoUID": "pqower816p8ztqdu",
              "Name": "demo/demo/20200727-135507-Peacock-Tubingen-2020-2gu.yml",
              "Root": "/",
              "Hash": "1d7bfeb03a4516f0ee2841fe3fa3c677e89d7490",
              "Size": 474,
              "Primary": false,
              "TimeIndex": "79799272864493-9999999901-2-fsoljb4g9jt4dx5i",
              "MediaID": "9999999901-2-fsoljb4g9jt4dx5i",
              "MediaUTC": 1595858107000,
              "FileType": "yml",
              "MediaType": "sidecar",
              "Mime": "text/plain",
              "Sidecar": true,
              "Diff": -1,
              "Chroma": -1,
              "ModTime": 1613943476,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "CreatedIn": 14463042,
              "UpdatedAt": "2025-01-14T16:40:42Z",
              "UpdatedIn": 14955833,
              "Markers": []
            },
            {
              "UID": "fsoljb42f1qyx91q",
              "PhotoUID": "pqower816p8ztqdu",
              "Name": "demo/demo/20200727-135507-Peacock-Tubingen-2020-2gu.jpg",
              "Root": "/",
              "Hash": "e296654221cb9badc0de201874cc4c816fdd8b24",
              "Size": 4508193,
              "Primary": true,
              "TimeIndex": "79799272864493-9999999901-0-fsoljb42f1qyx91q",
              "MediaID": "9999999901-0-fsoljb42f1qyx91q",
              "MediaUTC": 1595850907000,
              "Codec": "jpeg",
              "FileType": "jpg",
              "MediaType": "image",
              "Mime": "image/jpeg",
              "Width": 3648,
              "Height": 2736,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1.33,
              "MainColor": "grey",
              "Colors": "111116126",
              "Diff": 767,
              "Chroma": 8,
              "Software": "ELE-L29 10.1.0.133(C431E22R2P5)",
              "ModTime": 1597414766,
              "CreatedAt": "2024-12-16T17:01:52Z",
              "CreatedIn": 643280876,
              "UpdatedAt": "2025-01-14T16:40:42Z",
              "UpdatedIn": 655232042,
              "Markers": []
            }
          ],
          "Labels": [
            {
              "PhotoID": 99,
              "LabelID": 58,
              "LabelSrc": "image",
              "Uncertainty": 13,
              "Photo": null,
              "Label": {
                "ID": 58,
                "UID": "lsoljb42ygdy532y",
                "Slug": "peacock",
                "CustomSlug": "peacock",
                "Name": "Peacock",
                "Priority": 0,
                "Favorite": false,
                "Description": "",
                "Notes": "",
                "PhotoCount": 1,
                "Thumb": "e296654221cb9badc0de201874cc4c816fdd8b24",
                "CreatedAt": "2024-12-16T17:01:52Z",
                "UpdatedAt": "2025-02-04T22:24:54Z"
              }
            }
          ],
          "CreatedAt": "2021-02-21T21:37:56Z",
          "UpdatedAt": "2025-02-04T22:24:54Z",
          "EditedAt": "2025-02-04T22:24:54Z",
          "CheckedAt": "2025-04-07T18:02:13Z",
          "DeletedAt": null
        }),
        new Photo({
          "ID": 164,
          "TakenAt": "2020-06-09T14:42:50Z",
          "TakenAtLocal": "2020-06-09T16:42:50Z",
          "TakenSrc": "meta",
          "UID": "pqowerj3d8kfr56s",
          "Type": "video",
          "TypeSrc": "",
          "Title": "Plant / Meckenheim / 2020",
          "TitleSrc": "",
          "Caption": "",
          "CaptionSrc": "",
          "Path": "demo/demo/Videos",
          "Name": "20201216_141910_6E9CCAEC",
          "OriginalName": "",
          "Stack": 0,
          "Favorite": true,
          "Private": false,
          "Scan": false,
          "Panorama": false,
          "TimeZone": "Europe/Berlin",
          "PlaceID": "de:bM4IIOpnFIDL",
          "PlaceSrc": "meta",
          "CellID": "s2:479635f947fc",
          "CellAccuracy": 0,
          "Altitude": 0,
          "Lat": 49.4006,
          "Lng": 8.2519,
          "Country": "de",
          "Year": 2020,
          "Month": 6,
          "Day": 9,
          "Iso": 0,
          "Exposure": "",
          "FNumber": 0,
          "FocalLength": 0,
          "Quality": 7,
          "Resolution": 2,
          "Duration": 3200000000,
          "Color": 9,
          "CameraID": 1,
          "CameraSerial": "",
          "CameraSrc": "",
          "LensID": 1,
          "Details": {
            "PhotoID": 164,
            "Keywords": "böhler, demo, fruit, germany, green, meckenheim, plant, rheinland-pfalz, straße, videos",
            "KeywordsSrc": "",
            "Notes": "",
            "NotesSrc": "",
            "Subject": "",
            "SubjectSrc": "",
            "Artist": "",
            "ArtistSrc": "",
            "Copyright": "",
            "CopyrightSrc": "",
            "License": "",
            "LicenseSrc": "",
            "Software": "",
            "SoftwareSrc": "",
            "CreatedAt": "2024-12-16T17:02:03Z",
            "UpdatedAt": "2025-01-14T16:40:52Z"
          },
          "Camera": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown"
          },
          "Lens": {
            "ID": 1,
            "Slug": "zz",
            "Name": "Unknown",
            "Make": "",
            "Model": "Unknown",
            "Type": ""
          },
          "Cell": {
            "ID": "s2:479635f947fc",
            "Name": "",
            "Street": "Böhler Straße",
            "Postcode": "67149",
            "Category": "",
            "Place": {
              "PlaceID": "de:bM4IIOpnFIDL",
              "Label": "Meckenheim, Rheinland-Pfalz, Germany",
              "District": "",
              "City": "Meckenheim",
              "State": "Rheinland-Pfalz",
              "Country": "de",
              "Keywords": "",
              "Favorite": false,
              "PhotoCount": 3,
              "CreatedAt": "2024-12-16T17:02:03Z",
              "UpdatedAt": "2024-12-16T17:02:03Z"
            },
            "CreatedAt": "2024-12-16T17:02:03Z",
            "UpdatedAt": "2024-12-16T17:02:03Z"
          },
          "Place": {
            "PlaceID": "de:bM4IIOpnFIDL",
            "Label": "Meckenheim, Rheinland-Pfalz, Germany",
            "District": "",
            "City": "Meckenheim",
            "State": "Rheinland-Pfalz",
            "Country": "de",
            "Keywords": "",
            "Favorite": false,
            "PhotoCount": 3,
            "CreatedAt": "2024-12-16T17:02:03Z",
            "UpdatedAt": "2024-12-16T17:02:03Z"
          },
          "Albums": [],
          "Files": [
            {
              "UID": "fsoljbf02fbmge8t",
              "PhotoUID": "pqowerj3d8kfr56s",
              "Name": "demo/demo/Videos/20201216_141910_6E9CCAEC.yml",
              "Root": "/",
              "Hash": "d91c5cd0ea6d5d38791b937b48c45767fa8b48b1",
              "Size": 402,
              "Primary": false,
              "TimeIndex": "79799390835750-9999999836-2-fsoljbf02fbmge8t",
              "MediaID": "9999999836-2-fsoljbf02fbmge8t",
              "MediaUTC": 1591720970000,
              "FileType": "yml",
              "MediaType": "sidecar",
              "Mime": "text/plain",
              "Sidecar": true,
              "Diff": -1,
              "Chroma": -1,
              "ModTime": 1613943487,
              "CreatedAt": "2024-12-16T17:02:03Z",
              "CreatedIn": 19157250,
              "UpdatedAt": "2025-01-14T16:40:51Z",
              "UpdatedIn": 9623166,
              "Markers": []
            },
            {
              "UID": "fsoljbf6i2i4ikz3",
              "PhotoUID": "pqowerj3d8kfr56s",
              "Name": "demo/demo/Videos/20201216_141910_6E9CCAEC.mp4.jpg",
              "Root": "sidecar",
              "Hash": "9df3283ccb3aab06cfc7ba15004a32aefe040188",
              "Size": 109748,
              "Primary": true,
              "TimeIndex": "79799390835750-9999999836-0-fsoljbf6i2i4ikz3",
              "MediaID": "9999999836-0-fsoljbf6i2i4ikz3",
              "MediaUTC": 1591720970000,
              "FileType": "jpg",
              "MediaType": "image",
              "Mime": "image/jpeg",
              "Width": 1920,
              "Height": 1080,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1.78,
              "MainColor": "green",
              "Colors": "999939999",
              "Diff": 767,
              "Chroma": 35,
              "ModTime": 1724224263,
              "CreatedAt": "2024-12-16T17:02:03Z",
              "CreatedIn": 614768375,
              "UpdatedAt": "2025-01-14T16:40:52Z",
              "UpdatedIn": 515091209,
              "Markers": []
            },
            {
              "UID": "fsoljbfbv80i2t62",
              "PhotoUID": "pqowerj3d8kfr56s",
              "Name": "demo/demo/Videos/20201216_141910_6E9CCAEC.mp4",
              "Root": "/",
              "Hash": "48eeda078916bfc4b75116d14f2b193b3c389fd9",
              "Size": 7112164,
              "Primary": false,
              "TimeIndex": "79799390835750-9999999836-1-fsoljbfbv80i2t62",
              "MediaID": "9999999836-1-fsoljbfbv80i2t62",
              "MediaUTC": 1591713770000,
              "Codec": "avc1",
              "FileType": "mp4",
              "MediaType": "video",
              "Mime": "video/mp4",
              "Video": true,
              "Duration": 3200000000,
              "FPS": 28.8113241030388,
              "Frames": 92,
              "Width": 1920,
              "Height": 1080,
              "Orientation": 1,
              "OrientationSrc": "meta",
              "AspectRatio": 1.78,
              "MainColor": "green",
              "Colors": "999939999",
              "Diff": 767,
              "Chroma": 35,
              "ModTime": 1608130191,
              "CreatedAt": "2024-12-16T17:02:03Z",
              "CreatedIn": 26222417,
              "UpdatedAt": "2025-01-14T16:40:52Z",
              "UpdatedIn": 16329125,
              "Markers": []
            }
          ],
          "Labels": [
            {
              "PhotoID": 164,
              "LabelID": 3,
              "LabelSrc": "image",
              "Uncertainty": 7,
              "Photo": null,
              "Label": {
                "ID": 3,
                "UID": "lsoljaqor2dlaylg",
                "Slug": "plant",
                "CustomSlug": "plant",
                "Name": "Plant",
                "Priority": 0,
                "Favorite": false,
                "Description": "",
                "Notes": "",
                "PhotoCount": 10,
                "Thumb": "a7be9d1a20d0bab9bc3e3b2053ba0210eebfc700",
                "CreatedAt": "2024-12-16T17:01:38Z",
                "UpdatedAt": "2025-01-14T16:40:52Z"
              }
            }
          ],
          "CreatedAt": "2021-02-21T21:38:07Z",
          "UpdatedAt": "2025-01-14T16:40:52Z",
          "EditedAt": null,
          "CheckedAt": "2025-04-07T18:02:14Z",
          "DeletedAt": null
        })
      ],
      toggledPhotos: this.selectedPhotos,
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
    openPhoto(index) {
      this.$lightbox.openModels(Thumb.fromFiles([this.selectedPhotos[index]]), 0, null , this.isBatchDialog);
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
      PhotoClipboard.toggleAllIds(this.selectedPhotos);
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
          this.toggle(this.selectedPhotos[index]);
        }
      } else if (this.selectedPhotos[index]) {
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
      this.$clipboard.addRange(index, this.selectedPhotos);
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
