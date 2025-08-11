<template>
  <div class="chip-selector">
    <div v-if="shouldRenderChips" class="chip-selector__chips">
      <v-tooltip
        v-for="item in processedItems"
        :key="item.value || item.title"
        :text="getChipTooltip(item)"
        location="top"
      >
        <template #activator="{ props }">
          <div
            v-bind="props"
            :class="getChipClasses(item)"
            :aria-pressed="item.selected"
            :tabindex="0"
            role="button"
            @click="handleChipClick(item)"
            @keydown.enter="handleChipClick(item)"
            @keydown.space.prevent="handleChipClick(item)"
          >
            <div class="chip__content">
              <v-icon v-if="getChipIcon(item)" class="chip__icon">
                {{ getChipIcon(item) }}
              </v-icon>
              <span class="chip__text">{{ item.title }}</span>
            </div>
          </div>
        </template>
      </v-tooltip>

      <div v-if="processedItems.length === 0 && !showInput" class="chip-selector__empty">
        {{ emptyText }}
      </div>
    </div>

    <div v-if="allowCreate" class="chip-selector__input-container">
      <v-combobox
        ref="inputField"
        v-model="newItemTitle"
        :placeholder="computedInputPlaceholder"
        :persistent-placeholder="true"
        :items="availableItems"
        item-title="title"
        item-value="value"
        density="comfortable"
        hide-details
        hide-no-data
        return-object
        class="chip-selector__input"
        @keydown.enter="addNewItem"
        @update:model-value="onComboboxChange"
      >
        <template #no-data>
          <v-list-item>
            <v-list-item-title>
              {{ $gettext("Press enter to create new item") }}
            </v-list-item-title>
          </v-list-item>
        </template>
      </v-combobox>
    </div>
  </div>
</template>

<script>
export default {
  name: "ChipSelector",
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    availableItems: {
      type: Array,
      default: () => [],
    },
    allowCreate: {
      type: Boolean,
      default: true,
    },
    emptyText: {
      type: String,
      default: "",
    },
    inputPlaceholder: {
      type: String,
      default: "",
    },
    loading: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:items"],
  data() {
    return {
      newItemTitle: null,
      originalStates: new Map(),
    };
  },
  computed: {
    processedItems() {
      return this.items.map((item) => ({
        ...item,
        // Ensure action is always a string, never null/undefined
        action: item.action || "none",
        selected: item.action === "add" || item.action === "remove",
      }));
    },
    computedInputPlaceholder() {
      return this.inputPlaceholder || this.$gettext("Enter item name...");
    },
    showInput() {
      return this.allowCreate;
    },
    shouldRenderChips() {
      // Render chips container only when there are chips
      return this.processedItems.length > 0 || !this.showInput;
    },
  },
  watch: {
    items: {
      handler(newItems) {
        newItems.forEach((item) => {
          const itemKey = item.value || item.title;
          if (!item.isNew && !this.originalStates.has(itemKey)) {
            this.originalStates.set(itemKey, {
              mixed: item.mixed,
              action: item.action || "none",
            });
          }
        });
      },
      immediate: true,
      deep: true,
    },
  },
  methods: {
    getChipClasses(item) {
      const baseClass = "chip";
      const classes = [baseClass];

      if (this.loading || this.disabled) {
        classes.push(`${baseClass}--loading`);
      }

      if (item.action === "add") {
        classes.push(item.mixed ? `${baseClass}--green-light` : `${baseClass}--green`);
      } else if (item.action === "remove") {
        classes.push(item.mixed ? `${baseClass}--red-light` : `${baseClass}--red`);
      } else if (item.mixed) {
        classes.push(`${baseClass}--gray-light`);
      } else {
        classes.push(`${baseClass}--gray`);
      }

      return classes;
    },

    getChipIcon(item) {
      if (item.action === "add") return "mdi-plus";
      if (item.action === "remove") return "mdi-minus";
      if (item.mixed) return "mdi-circle-half-full";
      return null;
    },

    getChipTooltip(item) {
      if (item.action === "add") {
        return item.mixed ? this.$gettext("Add to all selected photos") : this.$gettext("Add to all");
      } else if (item.action === "remove") {
        return item.mixed ? this.$gettext("Remove from all selected photos") : this.$gettext("Remove from all");
      } else if (item.mixed) {
        return this.$gettext("Part of some selected photos");
      }
      return this.$gettext("Part of all selected photos");
    },

    handleChipClick(item) {
      if (this.loading || this.disabled) return;

      let newAction;

      if (item.mixed) {
        // Handle mixed state cycling
        switch (item.action) {
          case "none":
            newAction = "add";
            break;
          case "add":
            newAction = "remove";
            break;
          case "remove":
            newAction = "none";
            break;
        }
      } else {
        // Handle normal state cycling
        if (item.isNew) {
          newAction = item.action === "add" ? "remove" : "add";
        } else {
          newAction = item.action === "remove" ? "none" : "remove";
        }
      }

      this.updateItemAction(item, newAction);
    },

    updateItemAction(itemToUpdate, action) {
      // Special case: remove new items completely
      if (itemToUpdate.isNew && action === "remove") {
        const updatedItems = this.items.filter(
          (item) => (item.value || item.title) !== (itemToUpdate.value || itemToUpdate.title)
        );
        this.$emit("update:items", updatedItems);
        return;
      }

      // Update action for existing item
      const updatedItems = this.items.map((item) =>
        (item.value || item.title) === (itemToUpdate.value || itemToUpdate.title) ? { ...item, action } : item
      );

      this.$emit("update:items", updatedItems);
    },

    onComboboxChange(value) {
      if (value && typeof value === "object" && value.title) {
        this.newItemTitle = value;
        this.addNewItem();
        // Immediately clear the input, remove focus and restore placeholder
        this.$nextTick(() => {
          this.newItemTitle = "";
          if (this.$refs.inputField) {
            this.$refs.inputField.blur();
            // Force the combobox to reset completely
            setTimeout(() => {
              this.newItemTitle = null;
            }, 10);
          }
        });
      } else {
        this.newItemTitle = value;
      }
    },

    addNewItem() {
      // Extract title and value from input
      let title, value;

      if (typeof this.newItemTitle === "string") {
        title = this.newItemTitle.trim();
        value = "";
      } else if (this.newItemTitle && typeof this.newItemTitle === "object") {
        title = this.newItemTitle.title;
        value = this.newItemTitle.value;
      } else {
        return;
      }

      if (!title) return;

      const existingItem = this.items.find(
        (item) => item.title.toLowerCase() === title.toLowerCase() || (item.value && item.value === value)
      );

      if (existingItem) {
        console.warn(`Item "${title}" already exists`);
        return;
      }

      const newItem = {
        value: value || "",
        title,
        mixed: false,
        action: "add",
        isNew: true,
      };

      this.$emit("update:items", [...this.items, newItem]);
      this.newItemTitle = null;

      // Refocus input field
      this.$nextTick(() => {
        if (this.$refs.inputField) {
          this.$refs.inputField.focus();
        }
      });
    },
  },
};
</script>

<style src="../../css/chip-selector.css"></style>
