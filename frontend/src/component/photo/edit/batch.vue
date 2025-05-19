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
          <div v-if="model.models" class="edit-batch photo-results list-view v-table">
            <table>
              <tbody>
                <tr>
                  <td class="col-select" :class="{ 'is-selected': isAllSelected }">
                    <button
                      class="input-select ma-auto"
                      @click.stop.prevent="toggleAll"
                    >
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
                        <i v-if="m.Type === 'live'" class="action-live" :title="$gettext('Live')"><icon-live-photo /></i>
                        <i v-else-if="m.Type === 'animated'" class="mdi mdi-file-gif-box" :title="$gettext('Animated')" />
                        <i v-else-if="m.Type === 'video'" class="mdi mdi-play" :title="$gettext('Video')" />
                      </button>
                    </div>
                  </td>
                  <td
                    class="meta-data meta-title col-auto text-start clickable"
                    :title="m.FileName"
                    @click.exact="openPhoto(index)"
                  >
                    {{ m.FileName }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </v-col>

        <!-- Phone view -->
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
                              <i v-if="m.Type === 'live'" class="action-live" :title="$gettext('Live')"><icon-live-photo /></i>
                              <i v-else-if="m.Type === 'animated'" class="mdi mdi-file-gif-box" :title="$gettext('Animated')" />
                              <i v-else-if="m.Type === 'video'" class="mdi mdi-play" :title="$gettext('Video')" />
                            </button>
                          </div>
                        </td>
                        <td
                          class="meta-data meta-title col-auto text-start clickable"
                          :title="m.FileName"
                          @click.exact="openPhoto(index)"
                        >
                          {{ m.FileName }}
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
                        :model-value="getFieldData('text-field', 'Title').value"
                        :placeholder="getFieldData('text-field', 'Title').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'Title').persistent"
                        autocomplete="off"
                        density="comfortable"
                        class="input-title"
                        @update:modelValue="onInput($event, 'Title')"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-textarea
                        hide-details
                        autocomplete="off"
                        auto-grow
                        :label="$gettext('Subject')"
                        :model-value="getFieldData('text-field', 'DetailsSubject').value"
                        :placeholder="getFieldData('text-field', 'DetailsSubject').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'DetailsSubject').persistent"
                        :rows="1"
                        density="comfortable"
                        class="input-subject"
                        @update:modelValue="onInput($event, 'DetailsSubject')"
                      ></v-textarea>
                    </v-col>
                  </v-row>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-textarea
                      hide-details
                      autocomplete="off"
                      auto-grow
                      :label="$gettext('Caption')"
                      :model-value="getFieldData('text-field', 'Caption').value"
                      :placeholder="getFieldData('text-field', 'Caption').placeholder"
                      :persistent-placeholder="getFieldData('text-field', 'Caption').persistent"
                      :rows="1"
                      density="comfortable"
                      class="input-caption"
                      @update:modelValue="onInput($event, 'Caption')"
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
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :model-value="getFieldData('select-field', 'Day').value"
                        :items="getFieldData('select-field', 'Day').items"
                        :placeholder="getFieldData('select-field', 'Day').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Day').persistent"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-day"
                        @update:modelValue="onSelect($event, 'Day')"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="6" md="3">
                      <v-combobox
                        :label="$gettext('Month')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :model-value="getFieldData('select-field', 'Month').value"
                        :items="getFieldData('select-field', 'Month').items"
                        :placeholder="getFieldData('select-field', 'Month').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Month').persistent"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-month"
                        @update:modelValue="onSelect($event, 'Month')"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <v-combobox
                        :label="$gettext('Year')"
                        autocomplete="off"
                        hide-details
                        hide-no-data
                        :model-value="getFieldData('select-field', 'Year').value"
                        :items="getFieldData('select-field', 'Year').items"
                        :placeholder="getFieldData('select-field', 'Year').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Year').persistent"
                        item-title="text"
                        item-value="value"
                        density="comfortable"
                        validate-on="input"
                        class="input-year"
                        @update:modelValue="onSelect($event, 'Month')"
                      >
                      </v-combobox>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-autocomplete
                        :label="$gettext('Time Zone')"
                        hide-no-data
                        :model-value="getFieldData('select-field', 'TimeZone').value"
                        :items="getFieldData('select-field', 'TimeZone').items"
                        :placeholder="getFieldData('select-field', 'TimeZone').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'TimeZone').persistent"
                        item-value="ID"
                        item-title="Name"
                        density="comfortable"
                        class="input-timezone"
                        @update:modelValue="onSelect($event, 'TimeZone')"
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
                        :label="$gettext('Country')"
                        hide-details
                        hide-no-data
                        autocomplete="off"
                        item-value="Code"
                        item-title="Name"
                        :model-value="getFieldData('select-field', 'Country').value"
                        :items="getFieldData('select-field', 'Country').items"
                        :placeholder="getFieldData('select-field', 'Country').placeholder"
                        :persistent-placeholder="getFieldData('select-field', 'Country').persistent"
                        density="comfortable"
                        validate-on="input"
                        class="input-country"
                        @update:modelValue="onSelect($event, 'Country')"
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
                        :placeholder="getFieldData('input-field', 'Lat').placeholder"
                        :persistent-placeholder="getFieldData('input-field', 'Lat').persistent"
                        :model-value="getFieldData('input-field', 'Lat').value"
                        density="comfortable"
                        validate-on="input"
                        class="input-latitude"
                        @update:modelValue="onInput($event, 'Lat')"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="3">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        autocorrect="off"
                        autocapitalize="none"
                        :label="$gettext('Longitude')"
                        :placeholder="getFieldData('input-field', 'Lng').placeholder"
                        :persistent-placeholder="getFieldData('input-field', 'Lng').persistent"
                        :model-value="getFieldData('input-field', 'Lng').value"
                        density="comfortable"
                        validate-on="input"
                        class="input-longitude"
                        @update:modelValue="onInput($event, 'Lng')"
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
                        :placeholder="getFieldData('input-field', 'Altitude').placeholder"
                        :persistent-placeholder="getFieldData('input-field', 'Altitude').persistent"
                        :model-value="getFieldData('input-field', 'Altitude').value"
                        color="surface-variant"
                        density="comfortable"
                        validate-on="input"
                        class="input-altitude"
                        @update:modelValue="onInput($event, 'Altitude')"
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
                        :model-value="getFieldData('text-field', 'DetailsArtist').value"
                        :placeholder="getFieldData('text-field', 'DetailsArtist').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'DetailsArtist').persistent"
                        density="comfortable"
                        class="input-artist"
                        @update:modelValue="onInput($event, 'DetailsArtist')"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field
                        hide-details
                        autocomplete="off"
                        :label="$gettext('Copyright')"
                        :model-value="getFieldData('text-field', 'DetailsCopyright').value"
                        :placeholder="getFieldData('text-field', 'DetailsCopyright').placeholder"
                        :persistent-placeholder="getFieldData('text-field', 'DetailsCopyright').persistent"
                        density="comfortable"
                        class="input-copyright"
                        @update:modelValue="onInput($event, 'DetailsCopyright')"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <v-col cols="12" class="d-flex align-self-stretch flex-column">
                    <v-textarea
                      hide-details
                      autocomplete="off"
                      auto-grow
                      :label="$gettext('License')"
                      :model-value="getFieldData('text-field', 'DetailsLicense').value"
                      :placeholder="getFieldData('text-field', 'DetailsLicense').placeholder"
                      :persistent-placeholder="getFieldData('text-field', 'DetailsLicense').persistent"
                      :rows="1"
                      density="comfortable"
                      class="input-license"
                      @update:modelValue="onInput($event, 'DetailsLicense')"
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
                  :disabled="this.selectionsFullInfo.length < 1"
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
import { Batch } from "model/batch-edit";
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
      // deletedFields: {},
    };
  },
  computed: {
    title() {
      return this.$gettext(`Batch Edit (${this.allSelectedLength})`);
    },
  },
  watch: {
    visible: async function (show) {
      if (show) {
        this.expanded = [];

        await this.model.getData(this.selection);
        this.values = this.model.values;
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
    getFieldData(fieldType, fieldName) {
      const fieldData = this.values[fieldName];

      if (!fieldData) return { placeholder: "", persistent: false };

      if (fieldType === "text-field") {
        if (fieldData.mixed) {
          return { value: "", placeholder: "<mixed>", persistent: true };
        } else if (fieldData.value !== null && fieldData.value !== "") {
          return { value: fieldData.value, placeholder: "", persistent: false };
        } else {
          return { value: "", placeholder: "", persistent: false };
        }
      }

      if (fieldType === "input-field") {
        if (fieldData.mixed) {
          return { value: "", placeholder: "<mixed>", persistent: true };
        } else if (fieldData.value === 0 || (fieldData.value !== null && fieldData.value !== "")) {
          return { value: fieldData.value, placeholder: "", persistent: false };
        } else {
          return { value: "", placeholder: "", persistent: false };
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
        return items.find(item => item.value === -2).text;
      }
      if (fieldName === "Country") {
        return items.find(item => item.Code === -2).Name;
      }
      if (fieldName === "TimeZone") {
        return items.find(item => item.ID === -2).Name;
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
    },

    getCountriesArray(array) {
      const hasMixed = array.some(item => item.Code === -2);
      if (!hasMixed) {
        array.push({ Code: -2, Name: "<mixed>" });
      }
      return array;
    },

    onInput(val, fieldName) {
      if (this.values[fieldName]) {
        console.log("onInput");
      }
    },

    onSelect(val, fieldName) {
      if (this.values[fieldName]) {
        console.log("onSelect");
      }
    },

    // deleteValue(fieldName) {
    //   const fieldData = this.values[fieldName];
    //   if (fieldData.mixed || fieldData.value !== null) {
    //     this.deletedFields[fieldName] = {
    //       originalValue: fieldData.value,
    //       mixed: fieldData.mixed,
    //     };
    //     this.values[fieldName].value = '<deleted>';
    //   }
    // },
    //
    // undoDeleteValue(fieldName) {
    //   if (this.deletedFields[fieldName]) {
    //     const { originalValue, mixed } = this.deletedFields[fieldName];
    //     const fieldData = this.values[fieldName];
    //
    //     if (mixed) {
    //       fieldData.value = '<mixed>';
    //     } else {
    //       fieldData.value = originalValue;
    //     }
    //
    //     this.deletedFields[fieldName] = null;
    //   }
    // },
    openPhoto(index) {
      this.$lightbox.openModels(Thumb.fromPhotos([this.model.models[index]]), 0, null , this.isBatchDialog);
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
    // onUpdate() {
    //   // console.log('Add event on update');
    // },
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
      this.isAllSelected = this.model.selection.every(photo => photo.selected);
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
