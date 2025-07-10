<template>
  <v-dialog
    v-if="formData"
    ref="dialog"
    :model-value="visible"
    :fullscreen="$vuetify.display.mdAndDown"
    scrim
    scrollable
    class="p-dialog p-photo-edit-batch v-dialog--sidepanel v-dialog--sidepanel-wide"
    @click.stop="onClick"
  >
    <v-card ref="content" class="edit-batch__card" :tile="$vuetify.display.mdAndDown" tabindex="1">
      <v-toolbar flat color="navigation" :density="$vuetify.display.mdAndDown ? 'compact' : 'comfortable'">
        <v-btn icon class="action-close" @click.stop="onClose">
          <v-icon>mdi-close</v-icon>
        </v-btn>

        <v-toolbar-title>
          {{ formTitle }}
        </v-toolbar-title>
      </v-toolbar>

      <v-row dense :class="!$vuetify.display.mdAndDown ? 'overflow-hidden' : ''">
        <!-- Desktop view -->
        <v-col v-if="!$vuetify.display.mdAndDown" cols="12" lg="4" class="scroll-col">
          <div v-if="model.models" class="edit-batch photo-results list-view v-table">
            <table>
              <tbody>
                <tr>
                  <td class="col-select" :class="{ 'is-selected': isAllSelected }">
                    <button class="input-select ma-auto" @click.stop.prevent="toggleAll">
                      <i class="mdi mdi-checkbox-marked select-on" />
                      <i class="mdi mdi-checkbox-blank-outline select-off" />
                    </button>
                  </td>
                  <td class="media result col-preview">Pictures</td>
                </tr>
                <tr v-for="(m, index) in model.models" :key="m.ID" ref="items" :data-index="index">
                  <td :data-id="m.ID" :data-uid="m.UID" class="col-select" :class="{ 'is-selected': isSelected(m) }">
                    <button
                      class="input-select ma-auto"
                      @touchstart.passive="onMouseDown($event, index)"
                      @touchend.stop="onSelectClick($event, index, true)"
                      @mousedown="onMouseDown($event, index)"
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
                    <!--                        @click.stop.prevent="onClick($event, index, false)"-->
                    <!--                    >-->
                    <div
                      :style="`background-image: url(${m.thumbnailUrl('tile_224')})`"
                      class="preview"
                      @touchstart.passive="onMouseDown($event, index)"
                      @touchend.stop="onSelectClick($event, index, false)"
                      @mousedown="onMouseDown($event, index)"
                      @click.stop.prevent="onSelectClick($event, index, false)"
                    >
                      <div class="preview__overlay"></div>
                      <button
                        v-if="m.Type === 'video' || m.Type === 'live' || m.Type === 'animated'"
                        class="input-open"
                        @click.stop.prevent="openPhoto(index)"
                      >
                        <i v-if="m.Type === 'live'" class="action-live" :title="$gettext('Live')"
                          ><icon-live-photo
                        /></i>
                        <i
                          v-else-if="m.Type === 'animated'"
                          class="mdi mdi-file-gif-box"
                          :title="$gettext('Animated')"
                        />
                        <i v-else-if="m.Type === 'video'" class="mdi mdi-play" :title="$gettext('Video')" />
                      </button>
                    </div>
                  </td>
                  <td
                    class="meta-data meta-title col-auto text-start clickable"
                    :title="m.FileName"
                    @click.exact="openPhoto(index)"
                  >
                    {{ m.getOriginalName() }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </v-col>

        <!-- Mobile view -->
        <v-col v-else cols="12">
          <div v-if="model.models" class="edit-batch photo-results list-view v-table">
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
                      <tr v-for="(m, index) in model.models" :key="m.ID" ref="items" :data-index="index">
                        <td
                          :data-id="m.ID"
                          :data-uid="m.UID"
                          class="col-select"
                          :class="{ 'is-selected': isSelected(m) }"
                        >
                          <button
                            class="input-select ma-auto"
                            @touchstart.passive="onMouseDown($event, index)"
                            @touchend.stop="onSelectClick($event, index, true)"
                            @mousedown="onMouseDown($event, index)"
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
                          <!--                        @click.stop.prevent="onClick($event, index, false)"-->
                          <!--                    >-->
                          <div
                            :style="`background-image: url(${m.thumbnailUrl('tile_224')})`"
                            class="preview"
                            @touchstart.passive="onMouseDown($event, index)"
                            @touchend.stop="onSelectClick($event, index, false)"
                            @mousedown="onMouseDown($event, index)"
                            @click.stop.prevent="onSelectClick($event, index, false)"
                          >
                            <div class="preview__overlay"></div>
                            <button
                              v-if="m.Type === 'video' || m.Type === 'live' || m.Type === 'animated'"
                              class="input-open"
                              @click.stop.prevent="openPhoto(index)"
                            >
                              <i v-if="m.Type === 'live'" class="action-live" :title="$gettext('Live')"
                                ><icon-live-photo
                              /></i>
                              <i
                                v-else-if="m.Type === 'animated'"
                                class="mdi mdi-file-gif-box"
                                :title="$gettext('Animated')"
                              />
                              <i v-else-if="m.Type === 'video'" class="mdi mdi-play" :title="$gettext('Video')" />
                            </button>
                          </div>
                        </td>
                        <td
                          class="meta-data meta-title col-auto text-start clickable edit-batch__file-name"
                          :title="m.FileName"
                          @click.exact="openPhoto(index)"
                        >
                          {{ m.getOriginalName() }}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </div>
        </v-col>

        <v-col v-if="model.values" cols="12" lg="8" class="scroll-col">
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
                        :model-value="formData.Title.value"
                        :placeholder="getFieldData('text-field', 'Title').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'Title').persistent"
                        :append-inner-icon="getIcon('text-field', 'Title')"
                        autocomplete="off"
                        density="comfortable"
                        class="input-title"
                        @click:append-inner="toggleField('Title', $event)"
                        @update:model-value="(val) => changeValue(val, 'text-field', 'Title')"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-textarea
                        hide-details
                        autocomplete="off"
                        auto-grow
                        :label="$gettext('Subject')"
                        :model-value="formData.DetailsSubject.value"
                        :placeholder="getFieldData('text-field', 'DetailsSubject').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'DetailsSubject').persistent"
                        :append-inner-icon="getIcon('text-field', 'DetailsSubject')"
                        :rows="1"
                        density="comfortable"
                        class="input-subject"
                        @click:append-inner="toggleField('DetailsSubject', $event)"
                        @update:model-value="(val) => changeValue(val, 'text-field', 'DetailsSubject')"
                      ></v-textarea>
                    </v-col>
                  </v-row>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-textarea
                      hide-details
                      autocomplete="off"
                      auto-grow
                      :label="$gettext('Caption')"
                      :model-value="formData.Caption.value"
                      :placeholder="getFieldData('text-field', 'Caption').placeholder"
                      :persistent-placeholder="getFieldData('text-field', 'Caption').persistent"
                      :append-inner-icon="getIcon('text-field', 'Caption')"
                      :rows="1"
                      density="comfortable"
                      class="input-caption"
                      @click:append-inner="toggleField('Caption', $event)"
                      @update:model-value="(val) => changeValue(val, 'text-field', 'Caption')"
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
                        v-model="formData.Day.value"
                        :label="$gettext('Day')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :items="getFieldData('select-field', 'Day').items"
                        :placeholder="getFieldData('select-field', 'Day').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Day').persistent"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-day"
                        @update:model-value="(val) => changeSelectValue(val, 'select-field', 'Day')"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="6" md="3">
                      <v-combobox
                        v-model="formData.Month.value"
                        :label="$gettext('Month')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :items="getFieldData('select-field', 'Month').items"
                        :placeholder="getFieldData('select-field', 'Month').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Month').persistent"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-month"
                        @update:model-value="(val) => changeSelectValue(val, 'select-field', 'Month')"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <v-combobox
                        v-model="formData.Year.value"
                        :label="$gettext('Year')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :items="getFieldData('select-field', 'Year').items"
                        :placeholder="getFieldData('select-field', 'Year').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Year').persistent"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-year"
                        @update:model-value="(val) => changeSelectValue(val, 'select-field', 'Year')"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-autocomplete
                        v-model="formData.TimeZone.value"
                        :label="$gettext('Time Zone')"
                        hide-no-data
                        :items="getFieldData('select-field', 'TimeZone').items"
                        :placeholder="getFieldData('select-field', 'TimeZone').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'TimeZone').persistent"
                        item-value="ID"
                        item-title="Name"
                        density="comfortable"
                        class="input-timezone"
                        @update:model-value="(val) => changeSelectValue(val, 'select-field', 'TimeZone')"
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
                        v-model="formData.Country.value"
                        :label="$gettext('Country')"
                        hide-details
                        hide-no-data
                        autocomplete="off"
                        item-value="Code"
                        item-title="Name"
                        :items="getFieldData('select-field', 'Country').items"
                        :placeholder="getFieldData('select-field', 'Country').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Country').persistent"
                        density="comfortable"
                        validate-on="input"
                        class="input-country"
                        @update:model-value="(val) => changeSelectValue(val, 'select-field', 'Country')"
                      >
                      </v-autocomplete>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <v-text-field
                        hide-details
                        flat
                        autocomplete="off"
                        autocorrect="off"
                        autocapitalize="none"
                        :label="$gettext('Altitude (m)')"
                        :model-value="formData.Altitude.value"
                        :placeholder="getFieldData('input-field', 'Altitude').placeholder"
                        :persistent-placeholder="getFieldData('input-field', 'Altitude').persistent"
                        :append-inner-icon="getIcon('input-field', 'Altitude')"
                        color="surface-variant"
                        density="comfortable"
                        validate-on="input"
                        class="input-altitude"
                        @click:append-inner="toggleField('Altitude', $event)"
                        @update:model-value="(val) => changeValue(val, 'input-field', 'Altitude')"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <p-location-input
                        :latlng="currentCoordinates"
                        :placeholder="isLocationMixed ? '' : locationPlaceholder"
                        hide-details
                        :label="isLocationMixed ? '<mixed>' : locationLabel"
                        density="comfortable"
                        validate-on="input"
                        :show-map-button="!placesDisabled"
                        :map-button-title="$gettext('Adjust Location')"
                        :map-button-disabled="placesDisabled"
                        class="input-coordinates"
                        @update:latlng="updateLatLng"
                        @changed="onLocationChanged"
                        @open-map="adjustLocation"
                      ></p-location-input>
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
                        :model-value="formData.DetailsArtist.value"
                        :placeholder="getFieldData('text-field', 'DetailsArtist').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'DetailsArtist').persistent"
                        :append-inner-icon="getIcon('text-field', 'DetailsArtist')"
                        density="comfortable"
                        class="input-artist"
                        @click:append-inner="toggleField('DetailsArtist', $event)"
                        @update:model-value="(val) => changeValue(val, 'text-field', 'DetailsArtist')"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        :label="$gettext('Copyright')"
                        :model-value="formData.DetailsCopyright.value"
                        :placeholder="getFieldData('text-field', 'DetailsCopyright').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'DetailsCopyright').persistent"
                        :append-inner-icon="getIcon('text-field', 'DetailsCopyright')"
                        density="comfortable"
                        class="input-copyright"
                        @click:append-inner="toggleField('DetailsCopyright', $event)"
                        @update:model-value="(val) => changeValue(val, 'text-field', 'DetailsCopyright')"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-textarea
                      hide-details
                      autocomplete="off"
                      auto-grow
                      :label="$gettext('License')"
                      :model-value="formData.DetailsLicense.value"
                      :placeholder="getFieldData('text-field', 'DetailsLicense').placeholder"
                      :persistent-placeholder="getFieldData('text-field', 'DetailsLicense').persistent"
                      :append-inner-icon="getIcon('text-field', 'DetailsLicense')"
                      :rows="1"
                      density="comfortable"
                      class="input-license"
                      @click:append-inner="toggleField('DetailsLicense', $event)"
                      @update:model-value="(val) => changeValue(val, 'text-field', 'DetailsLicense')"
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

                <div>
                  <v-col cols="12" md="6">
                    <p>File Type</p>
                  </v-col>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-combobox
                      v-model="formData.Type.value"
                      :label="$gettext('Type')"
                      autocomplete="off"
                      hide-details
                      hide-no-data
                      :items="getFieldData('select-field', 'Type').items"
                      :placeholder="getFieldData('select-field', 'Type').placeholder"
                      :persistent-placeholder="getFieldData('select-field', 'Type').persistent"
                      item-title="text"
                      item-value="value"
                      density="comfortable"
                      validate-on="input"
                      class="input-type"
                      @update:model-value="(val) => changeSelectValue(val, 'select-field', 'Type')"
                    >
                    </v-combobox>
                  </v-col>
                </div>

                <v-row>
                  <v-col v-for="fieldName in toggleFieldsArray" :key="fieldName" cols="12" sm="6" md="6">
                    <div class="d-flex flex-column">
                      <v-col cols="12">
                        <p>{{ fieldName }}</p>
                      </v-col>
                      <v-col cols="12" class="d-flex align-self-stretch flex-column">
                        <v-btn-toggle v-model="formData[fieldName].value" mandatory color="primary">
                          <v-btn
                            v-for="option in toggleOptions(fieldName)"
                            :key="option.value"
                            :value="option.value"
                            @click="changeToggleValue(option.value, fieldName)"
                          >
                            <span>{{ option.text }}</span>
                          </v-btn>
                        </v-btn-toggle>
                      </v-col>
                    </div>
                  </v-col>
                </v-row>
              </div>
            </div>

            <div class="form-actions form-actions--sticky">
              <div class="action-buttons">
                <v-btn color="button" variant="flat" class="action-close" @click.stop="close">
                  {{ $gettext(`Close`) }}
                </v-btn>
                <v-btn color="highlight" variant="flat" class="action-apply action-approve" @click.stop="save(false)">
                  <span>{{ $gettext(`Apply`) }}</span>
                </v-btn>
              </div>
            </div>
          </v-form>
        </v-col>
      </v-row>
    </v-card>
    <p-location-dialog
      :visible="locationDialog"
      :latlng="currentCoordinates"
      @close="locationDialog = false"
      @confirm="confirmLocation"
    ></p-location-dialog>
  </v-dialog>
</template>
<script>
import * as options from "options/options";
import countries from "options/countries.json";
import IconLivePhoto from "../../icon/live-photo.vue";
import { Batch } from "model/batch-edit";
import Thumb from "../../../model/thumb";
import PLocationDialog from "component/location/dialog.vue";
import PLocationInput from "component/location/input.vue";

export default {
  name: "PPhotoEditBatch",
  components: {
    IconLivePhoto,
    PLocationDialog,
    PLocationInput,
  },
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
      model: new Batch(),
      uid: "",
      loading: false,
      subscriptions: [],

      selectionsFullInfo: [],
      selectedPhotosLength: 0,
      expanded: [0],
      isBatchDialog: true,
      isAllSelected: true,
      allSelectedLength: 0,
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
      values: {},
      formData: null,
      previousFormData: {},
      deletedFields: {},
      toggleFieldsArray: ["Scan", "Favorite", "Private", "Panorama"],
      actions: { none: "none", update: "update", add: "add", remove: "remove" },
      locationDialog: false,
      placesDisabled: !this.$config.feature("places"),
      locationLabel: this.$gettext("Location"),
    };
  },
  computed: {
    formTitle() {
      return this.$gettext(`Batch Edit (${this.allSelectedLength})`);
    },
    currentCoordinates() {
      const latData = this.values?.Lat;
      const lngData = this.values?.Lng;

      // If no form data yet, return default
      if (!this.formData || !latData || !lngData) {
        return [0, 0];
      }

      const lat = this.formData.Lat.value;
      const lng = this.formData.Lng.value;

      // If form data has been updated, use the form data values
      if (this.formData.Lat.action === this.actions.update || this.formData.Lng.action === this.actions.update) {
        return [parseFloat(lat) || 0, parseFloat(lng) || 0];
      }

      // Handle mixed values or empty values
      if (
        latData.mixed ||
        lngData.mixed ||
        lat === "<mixed>" ||
        lng === "<mixed>" ||
        lat === "" ||
        lng === "" ||
        lat === null ||
        lng === null
      ) {
        return [0, 0];
      }

      return [parseFloat(lat) || 0, parseFloat(lng) || 0];
    },
    locationPlaceholder() {
      const latData = this.values?.Lat;
      const lngData = this.values?.Lng;

      // If either coordinate has mixed values
      if (latData?.mixed || lngData?.mixed) {
        return "<mixed>";
      }

      // If both coordinates are empty/null (no location info)
      const lat = this.formData?.Lat?.value;
      const lng = this.formData?.Lng?.value;

      if ((!lat || lat === 0) && (!lng || lng === 0)) {
        return ""; // Empty state
      }

      return ""; // Default placeholder
    },
    isLocationMixed() {
      const latData = this.values?.Lat;
      const lngData = this.values?.Lng;

      // If form data has been updated, not mixed anymore
      if (this.formData?.Lat?.action === this.actions.update || this.formData?.Lng?.action === this.actions.update) {
        return false;
      }

      return latData?.mixed || lngData?.mixed;
    },
  },
  watch: {
    visible: async function (show) {
      if (show) {
        this.expanded = [];

        await this.model.getData(this.selection);
        this.values = this.model.values;
        this.setFormData();
        this.allSelectedLength = this.model.getLengthOfAllSelected();
      } else {
        this.model = new Batch();
      }
    },
  },
  created() {
    // this.subscriptions.push(this.$event.subscribe("photos.updated", (ev, data) => this.onUpdate(ev, data)));
  },
  beforeUnmount() {
    for (let i = 0; i < this.subscriptions.length; i++) {
      this.$event.unsubscribe(this.subscriptions[i]);
    }
  },
  methods: {
    changeValue(newValue, fieldType, fieldName) {
      if (!fieldName) return;

      const previousValue = this.previousFormData[fieldName].value;
      this.formData[fieldName].action = this.actions.update;
      this.formData[fieldName].value = newValue;

      if (newValue === previousValue) {
        this.formData[fieldName].action = this.actions.none;
      }

      this.getIcon(fieldType, fieldName);
    },
    changeSelectValue(newValue, fieldType, fieldName) {
      if (!fieldName) return;

      const previousValue = this.previousFormData[fieldName].value;
      this.formData[fieldName].action = this.actions.update;

      if (fieldName === "Day" || fieldName === "Month" || fieldName === "Year" || fieldName === "Type") {
        this.formData[fieldName].value = newValue.text;

        if (newValue.text === previousValue) {
          this.formData[fieldName].action = this.actions.none;
        }
      } else {
        this.formData[fieldName].value = newValue;

        const newVal = newValue !== -2 ? newValue : "<mixed>";
        if (newVal === previousValue) {
          this.formData[fieldName].action = this.actions.none;
        }
      }
    },
    changeToggleValue(newValue, fieldName) {
      if (!fieldName) return;

      const previousValue = this.previousFormData[fieldName].value;
      this.formData[fieldName].action = this.actions.update;
      this.formData[fieldName].value = newValue;

      if (newValue === String(previousValue) || newValue === previousValue) {
        this.formData[fieldName].action = this.actions.none;
      }
    },
    setFormData() {
      this.formData = this.model.getDefaultFormData();

      const fieldConfigs = [
        { type: "text-field", name: "Title" },
        { type: "text-field", name: "DetailsSubject" },
        { type: "text-field", name: "Caption" },
        { type: "select-field", name: "Day" },
        { type: "select-field", name: "Month" },
        { type: "select-field", name: "Year" },
        { type: "select-field", name: "TimeZone" },
        { type: "select-field", name: "Country" },
        { type: "input-field", name: "Altitude" },
        { type: "input-field", name: "Lat" },
        { type: "input-field", name: "Lng" },
        { type: "text-field", name: "DetailsArtist" },
        { type: "text-field", name: "DetailsCopyright" },
        { type: "text-field", name: "DetailsLicense" },
        { type: "select-field", name: "Type" },
      ];

      fieldConfigs.forEach(({ type, name, key }) => {
        const formKey = key || name;
        const fieldData = this.values[formKey];

        const { value, placeholder } = this.getFieldData(type, name);
        this.formData[formKey] = {
          action: this.actions.none,
          mixed: fieldData.mixed,
          value: value ? value : "",
        };

        if (type === "text-field" || type === "input-field") {
          this.previousFormData[formKey] = { value, placeholder, action: fieldData.action, mixed: fieldData.mixed };
        } else {
          this.previousFormData[formKey] = { value, action: fieldData.action, mixed: fieldData.mixed };
        }
      });

      // Set values for toggle fields
      this.toggleFieldsArray.forEach((fieldName) => {
        const fieldData = this.values[fieldName];
        const toggleValue = this.getToggleValue(fieldName);

        this.formData[fieldName] = { action: this.actions.none, mixed: fieldData.mixed, value: toggleValue };
        this.previousFormData[fieldName] = { action: fieldData.action, mixed: fieldData.mixed, value: toggleValue };
      });
    },
    toggleOptions(fieldName) {
      const fieldData = this.values[fieldName];
      if (!fieldData) return [];

      const options = [
        { text: "Yes", value: "true" },
        { text: "No", value: "false" },
      ];

      if (fieldData.mixed) {
        options.push({ text: "<mixed>", value: "<mixed>" });
      }

      return options;
    },
    getToggleValue(fieldName) {
      const fieldData = this.values[fieldName];
      if (!fieldData) return null;

      if (fieldData.mixed) {
        return "<mixed>";
      } else {
        return fieldData.value ? "true" : "false";
      }
    },
    toggleField(fieldName, event) {
      const classList = event.target.classList;

      if (classList.contains("mdi-undo")) {
        this.deletedFields[fieldName] = false;
        this.formData[fieldName].action = this.actions.none;
        this.formData[fieldName].value = this.previousFormData[fieldName]?.value || "";

        // TODO: add this if it is necessary to change the mixed value
        // if (this.formData[fieldName].mixed !== this.previousFormData[fieldName].mixed) {
        //   this.formData[fieldName].mixed = true;
        // }
      } else if (classList.contains("mdi-delete")) {
        this.deletedFields[fieldName] = true;
        this.formData[fieldName].action = this.actions.remove;
        this.formData[fieldName].value = "";
      }
    },
    getIcon(fieldType, fieldName) {
      const fieldData = this.values[fieldName];
      const isDeleted = this.deletedFields?.[fieldName];

      if (!fieldData) return;
      const previousField = this.previousFormData[fieldName];

      if (this.formData[fieldName].value !== previousField?.value || isDeleted) {
        return "mdi-undo";
      } else if (fieldData.mixed) {
        return "mdi-delete";
      } else if (fieldType === "text-field" && fieldData.value !== null && fieldData.value !== "") {
        return "mdi-delete";
      } else if (
        fieldType === "input-field" &&
        fieldData.value !== 0 &&
        fieldData.value !== null &&
        fieldData.value !== ""
      ) {
        return "";
      }
    },
    getFieldData(fieldType, fieldName) {
      const fieldData = this.values[fieldName];
      const isDeleted = this.deletedFields?.[fieldName];

      if (!fieldData) return { value: "", placeholder: "", persistent: false };

      if (fieldType === "text-field") {
        if (isDeleted) {
          return { value: "", placeholder: "<deleted>", persistent: true };
        } else if (fieldData.mixed) {
          return { value: "", placeholder: "<mixed>", persistent: true };
        } else if (fieldData.value !== null && fieldData.value !== "") {
          return { value: fieldData.value, placeholder: "", persistent: false };
        } else {
          return { value: "", placeholder: "", persistent: false };
        }
      }

      if (fieldType === "input-field") {
        if (isDeleted) {
          return { value: 0, placeholder: "", persistent: false };
        } else if (fieldData.mixed) {
          return { value: "", placeholder: "<mixed>", persistent: true };
        } else if (fieldData.value !== 0 && fieldData.value !== null && fieldData.value !== "") {
          return { value: fieldData.value, placeholder: "", persistent: false };
        } else {
          return { value: fieldData.value || 0, placeholder: "", persistent: false };
        }
      }

      if (fieldType === "select-field") {
        if (fieldData.mixed) {
          const items = this.getItemsArray(fieldName, true);
          const value = this.getValue(fieldName, items);
          return { value: value, placeholder: "<mixed>", persistent: true, items: items };
        } else if (!fieldData.mixed) {
          const items = this.getItemsArray(fieldName, false);
          return { value: fieldData.value, placeholder: "", persistent: false, items: items };
        }
      }
    },
    getValue(fieldName, items) {
      if (fieldName === "Day" || fieldName === "Month" || fieldName === "Year") {
        return items.find((item) => item.value === -2).text;
      }
      if (fieldName === "Country") {
        return items.find((item) => item.Code === -2).Name;
      }
      if (fieldName === "TimeZone") {
        return items.find((item) => item.ID === -2).Name;
      }
      if (fieldName === "Type") {
        return items.find((item) => item.value === "mixed").text;
      }
    },
    getItemsArray(fieldName, isMixed) {
      if (fieldName === "Day") {
        return isMixed ? options.DaysBatchDialog() : options.Days();
      }
      if (fieldName === "Month") {
        return isMixed ? options.MonthsShortBatchDialog() : options.MonthsShort();
      }
      if (fieldName === "Year") {
        return isMixed ? options.YearsBatchDialog(1900) : options.Years(1900);
      }
      if (fieldName === "Country") {
        const newCountries = this.getCountriesArray(this.countries);
        return isMixed ? newCountries : this.countries;
      }
      if (fieldName === "TimeZone") {
        return isMixed ? options.TimeZonesBatchDialog() : options.TimeZones();
      }
      if (fieldName === "Type") {
        return isMixed ? options.PhotoTypesBatchDialog() : options.PhotoTypes();
      }
    },
    getCountriesArray(array) {
      const hasMixed = array.some((item) => item.Code === -2);
      if (!hasMixed) {
        array.push({ Code: -2, Name: "<mixed>" });
      }
      return array;
    },
    // onUpdate() {
    //   // console.log('Add event on update');
    // },
    onSelect(val, fieldName) {
      if (this.values[fieldName]) {
        console.log("onSelect");
      }
    },
    openPhoto(index) {
      this.$lightbox.openModels(Thumb.fromPhotos([this.model.models[index]]), 0, null, this.isBatchDialog);
    },
    isSelected(model) {
      return this.model.isSelected(model.UID);
    },
    onClick(ev) {
      // Closes dialog when user clicks on background and model data is unchanged.
      if (!ev || !ev?.target?.classList?.contains("v-overlay__scrim")) {
        return;
      }
      ev.preventDefault();
      this.onClose();
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
        this.toggle(this.model.models[index]);
      } else if (this.model.models[index]) {
        this.openPhoto(index);
      }
    },
    toggle(photo) {
      this.model.toggle(photo.UID);
      this.updateToggleAll();
      this.allSelectedLength = this.model.getLengthOfAllSelected();
    },
    updateToggleAll() {
      this.isAllSelected = this.model.selection.every((photo) => photo.selected);
    },
    toggleAll() {
      this.isAllSelected = !this.isAllSelected;
      this.model.toggleAll(this.isAllSelected);
      this.allSelectedLength = this.model.getLengthOfAllSelected();
    },
    onMouseDown(ev, index) {
      this.mouseDown.index = index;
      this.mouseDown.scrollY = window.scrollY;
      this.mouseDown.timeStamp = ev.timeStamp;
    },
    onClose() {
      // Closes the dialog only if model data is unchanged.
      // TODO: change the functionality if something was changed
      // if (this.model?.hasId() && this.model?.wasChanged()) {
      //   this.$refs?.dialog?.animateClick();
      // } else {
      this.close();
      // }
    },
    close() {
      // Close the dialog.
      this.$emit("close");
    },
    updateLatLng(latlng) {
      this.formData.Lat.value = latlng[0];
      this.formData.Lng.value = latlng[1];
      this.formData.Lat.action = this.actions.update;
      this.formData.Lng.action = this.actions.update;
    },
    onLocationChanged(data) {
      if (data && data.lat !== undefined && data.lng !== undefined) {
        this.formData.Lat.value = data.lat;
        this.formData.Lng.value = data.lng;
        this.formData.Lat.action = this.actions.update;
        this.formData.Lng.action = this.actions.update;
      }
    },
    adjustLocation() {
      this.locationDialog = true;
    },
    confirmLocation(data) {
      if (data && data.lat !== undefined && data.lng !== undefined) {
        this.formData.Lat.value = data.lat;
        this.formData.Lng.value = data.lng;
        this.formData.Lat.action = this.actions.update;
        this.formData.Lng.action = this.actions.update;
      }

      this.locationDialog = false;
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
      this.model.save(this.selection, this.formData);

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
