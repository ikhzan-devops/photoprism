<template>
  <v-dialog
    ref="dialog"
    :model-value="visible"
    :close-delay="0"
    :open-delay="0"
    persistent
    scrim
    max-width="360"
    class="p-dialog p-confirm-dialog"
    @keyup.esc.exact="close"
    @keyup.enter.exact="confirm"
    @after-enter="afterEnter"
    @after-leave="afterLeave"
    @focusout="onFocusOut"
  >
    <v-card ref="content" tabindex="0">
      <v-card-title class="d-flex justify-start align-center ga-3">
        <v-icon :icon="icon" :size="iconSize" color="primary"></v-icon>
        <div class="text-subtitle-1">{{ text ? text : $gettext(`Are you sure?`) }}</div>
      </v-card-title>
      <v-card-actions class="action-buttons">
        <v-btn variant="flat" tabindex="0" color="button" class="action-cancel action-close" @click.stop="close">
          {{ $gettext(`Cancel`) }}
        </v-btn>
        <v-btn color="highlight" tabindex="0" variant="flat" class="action-confirm" @click.stop="confirm">
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
  methods: {
    afterEnter() {
      this.$view.enter(this);
    },
    afterLeave() {
      this.$view.leave(this);
    },
    onFocusOut(ev) {
      if (!this.$view.isActive(this)) {
        return;
      }

      const el = this.$refs.content?.$el;

      if (!ev || !ev.target || !(ev.target instanceof HTMLElement) || !(el instanceof HTMLElement)) {
        return;
      }

      const next = ev.relatedTarget;
      const leavingDialog = !next || !(next instanceof Node) || !el.contains(next);

      if (leavingDialog) {
        el.focus();
        ev.preventDefault();
      }
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
