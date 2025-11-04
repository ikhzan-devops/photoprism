<template>
  <v-dialog
    ref="dialog"
    :model-value="visible"
    persistent
    max-width="360"
    class="p-dialog p-confirm-dialog"
    retain-focus
    @keydown.esc.exact.stop.prevent="close"
    @keydown.enter.exact.stop.prevent="confirm"
    @after-enter="afterEnter"
  >
    <v-card ref="content" tabindex="1">
      <v-card-title class="d-flex justify-start align-center ga-3">
        <v-icon :icon="icon" :size="iconSize" color="primary"></v-icon>
        <div class="text-subtitle-1">{{ text ? text : $gettext(`Are you sure?`) }}</div>
      </v-card-title>
      <v-card-actions class="action-buttons">
        <v-btn variant="flat" color="button" class="action-cancel action-close" @click.stop="close">
          {{ $gettext(`Cancel`) }}
        </v-btn>
        <v-btn color="highlight" variant="flat" class="action-confirm" @click.stop="confirm">
          {{ action ? action : $gettext(`Yes`) }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
export default {
  name: "PConfirmDialog",
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    icon: {
      type: String,
      default: "mdi-delete-outline",
    },
    iconSize: {
      type: Number,
      default: 54,
    },
    text: {
      type: String,
      default: "",
    },
    action: {
      type: String,
      default: "",
    },
  },
  emits: ["close", "confirm"],
  data() {
    return {};
  },
  watch: {
    visible(show) {
      if (show) {
        this.$nextTick(() => this.$view.enter(this, this.$refs?.content, ".action-confirm"));
      } else {
        this.$view.leave(this);
      }
    },
  },
  methods: {
    afterEnter() {
      this.$nextTick(() => this.$view.enter(this, this.$refs?.content, ".action-confirm"));
    },
    close() {
      this.$emit("close");
    },
    confirm() {
      this.$emit("confirm");
    },
  },
};
</script>
