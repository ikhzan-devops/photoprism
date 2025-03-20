<template>
  <v-snackbar id="p-notify" v-model="visible" :class="'p-notify--' + message.color">
    <v-icon v-if="message.icon" :icon="'mdi-' + message.icon" :color="message.color" start></v-icon>
    {{ message.text }}
    <template #actions>
      <v-btn icon="mdi-close" :color="'on-' + message.color" variant="text" @click="close"></v-btn>
    </template>
  </v-snackbar>
</template>
<script>
export default {
  name: "PNotify",
  data() {
    return {
      visible: false,
      messages: [],
      message: {
        icon: "",
        color: "transparent",
        text: "",
        delay: this.defaultDelay,
      },
      lastText: "",
      lastId: 1,
      subscriptionId: "",
      defaultColor: "info",
      defaultDelay: 2000,
      warningDelay: 3000,
      errorDelay: 8000,
    };
  },
  watch: {
    visible: function (show) {
      if (!show) {
        if (this.messages.length > 0) {
          this.show();
        } else {
          this.lastText = "";
        }
      }
    },
  },
  created() {
    this.subscriptionId = this.$event.subscribe("notify", this.onNotify);
  },
  beforeUnmount() {
    this.messages = [];
    this.visible = false;
    this.$event.unsubscribe(this.subscriptionId);
  },
  methods: {
    onNotify: function (ev, data) {
      const type = ev.split(".")[1];

      // Get the message.
      let m = data.message;

      // Skip empty messages.
      if (!m || !m.length) {
        console.warn("notify: empty message");
        return;
      }

      // Log notifications in test mode.
      if (this.$config.test) {
        console.log(type + ": " + m.toLowerCase());
        return;
      }

      // First letter of the message should be uppercase.
      m = m.replace(/^./, m[0].toUpperCase());

      switch (type) {
        case "warning":
          this.addWarningMessage(m);
          break;
        case "error":
          this.addErrorMessage(m);
          break;
        case "success":
          this.addSuccessMessage(m);
          break;
        case "info":
          this.addInfoMessage(m);
          break;
        default:
          alert(m);
      }
    },

    addSuccessMessage: function (message) {
      this.addMessage("success", "check-circle", message, this.defaultDelay);
    },

    addInfoMessage: function (message) {
      this.addMessage("info", "information-outline", message, this.defaultDelay);
    },

    addWarningMessage: function (message) {
      this.addMessage("warning", "alert", message, this.warningDelay);
    },

    addErrorMessage: function (message) {
      this.addMessage("error", "alert-circle-outline", message, this.errorDelay);
    },

    addMessage: function (color, icon, text, delay) {
      if (!text || text === this.lastText) {
        return;
      }

      this.lastId++;
      this.lastText = text;

      const m = {
        id: this.lastId,
        color,
        icon,
        text,
        delay,
      };

      this.messages.push(m);

      if (!this.visible) {
        this.show();
      }
    },
    close: function () {
      this.visible = false;
    },
    show: function () {
      const message = this.messages.shift();

      if (message) {
        this.message = message;

        if (!this.message.icon) {
          this.message.icon = "";
        }

        if (!this.message.color) {
          this.message.color = this.defaultColor;
        }

        if (!this.message.delay || this.message.delay <= 0) {
          this.message.delay = this.defaultDelay;
        }

        this.visible = true;

        setTimeout(() => {
          this.lastText = "";
          this.show();
        }, this.message.delay);
      } else {
        this.lastText = "";
        this.visible = false;
      }
    },
  },
};
</script>
