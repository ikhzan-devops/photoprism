<template>
  <v-dialog :model-value="show" persistent max-width="500" class="p-dialog p-service-add" @keydown.esc="close">
    <v-card>
      <v-card-title class="d-flex justify-start align-center ga-3">
        <v-icon size="28" color="primary">mdi-swap-horizontal</v-icon>
        <h6 class="text-h6">
          {{ $gettext(`Add Account`) }}
        </h6>
      </v-card-title>
      <v-card-text class="dense">
        <v-row align="center" dense>
          <v-col cols="12">
            <v-text-field
              v-model="model.AccURL"
              hide-details
              autofocus
              :label="$gettext('Service URL')"
              placeholder="https://www.example.com/"
              autocorrect="off"
              autocapitalize="none"
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="model.AccUser"
              hide-details
              :label="$gettext('Username')"
              :placeholder="$gettext('optional')"
              autocorrect="off"
              autocapitalize="none"
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="model.AccPass"
              hide-details
              autocomplete="new-password"
              autocapitalize="none"
              :label="$gettext('Password')"
              :placeholder="$gettext('optional')"
              :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
              :type="showPassword ? 'text' : 'password'"
              @click:append-inner="showPassword = !showPassword"
            ></v-text-field>
          </v-col>
          <v-col cols="12" class="text-start text-caption">
            {{
              $gettext(
                `Note: Only WebDAV servers, like Nextcloud or PhotoPrism, can be configured as remote service for backup and file upload.`
              )
            }}
            {{ $gettext(`Support for additional services, like Google Drive, will be added over time.`) }}
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="action-buttons">
        <v-btn variant="flat" color="button" class="action-cancel action-close" @click.stop="close">
          <span>{{ label.cancel }}</span>
        </v-btn>
        <v-btn variant="flat" color="highlight" class="action-confirm" @click.stop="confirm">
          <span>{{ label.confirm }}</span>
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import Service from "model/service";
import * as options from "options/options";

export default {
  name: "PServiceAdd",
  props: {
    show: Boolean,
  },
  data() {
    return {
      options: options,
      showPassword: false,
      loading: false,
      search: null,
      model: new Service(false),
      label: {
        cancel: this.$gettext("Cancel"),
        confirm: this.$gettext("Connect"),
      },
    };
  },
  watch: {
    show: function () {},
  },
  methods: {
    close() {
      this.$emit("close");
    },
    confirm() {
      this.loading = true;

      this.model.save().then((a) => {
        this.loading = false;
        this.$notify.success(this.$gettext("Account created"));
        this.$emit("confirm", a.UID);
      });
    },
  },
};
</script>
