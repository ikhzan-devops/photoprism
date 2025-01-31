<template>
  <footer class="p-about-footer text-ltr">
    <p class="flex-fill text-sm-start">
      <strong>
        <router-link to="/about" class="text-link text-selectable">{{ about }} {{ getMembership() }}</router-link>
      </strong>
      <span class="body-link text-selectable">
        <span class="cursor-text" @click.stop.prevent="$util.copyText(about, version)">Build</span>
        <a href="https://docs.photoprism.app/release-notes/" target="_blank" :title="version" class="body-link">{{
          build
        }}</a>
      </span>
    </p>
    <p class="hidden-xs text-sm-end">
      <a href="https://raw.githubusercontent.com/photoprism/photoprism/develop/NOTICE" target="_blank" class="text-link"
        >3rd-party software packages</a
      >
      <a href="https://www.photoprism.app/about/team/" target="_blank" class="body-link">© 2018-2025 PhotoPrism UG</a>
    </p>
  </footer>
</template>

<script>
export default {
  name: "PAboutFooter",
  data() {
    const ver = this.$config.getVersion().split("-");
    const build = ver.slice(0, 2).join("-");
    const about = this.$config.getAbout();
    const membership = this.$config.getMembership();
    const customer = this.$config.getCustomer();

    return {
      rtl: this.$rtl,
      build: build,
      about: about,
      membership: membership,
      customer: customer,
      version: this.$config.getVersion(),
      isDemo: this.$config.isDemo(),
    };
  },
  methods: {
    getMembership() {
      if (this.isDemo) {
        return "Demo";
      }

      const tier = this.$config.getTier();
      const edition = this.$config.getEdition();

      if (edition === "plus" && tier > 7) {
        return "Plus";
      } else if (edition === "plus" && tier > 5) {
        return "Essentials+";
      } else if (tier > 3) {
        return "Essentials";
      }

      return "CE";
    },
  },
};
</script>
