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
          <v-expansion-panels v-model="expanded" class="pa-0 elevation-0" variant="accordion" multiple>
            <v-expansion-panel tabindex="1" style="margin-top: 1px" class="pa-0 elevation-0">
              <v-expansion-panel-title>
                <div class="text-caption font-weight-bold filename">Date & Location</div>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <p>second</p>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-col>
      </v-row>
    </v-card>
  </v-dialog>
</template>
<script>
import Photo from "model/photo";

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

      expanded: [0],
      selections: [],
      view: this.$view.data(),
    };
  },
  computed: {
    title() {
      if (this.selection.length > 0) {
        return this.$gettext(`Edit ${this.selection.length} selected photos`);
      }

      return this.title;
    },
  },
  watch: {
    visible: function (show) {
      if (show) {
        this.find(this.index);
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
