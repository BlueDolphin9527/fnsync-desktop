<template>
  <router-view></router-view>
</template>

<script>
export default {
  data() {
    return {};
  },
  created: function () {
    this.$router.push("/");
    this.registerEvent();
  },
  methods: {
    registerEvent: function () {
      var $this = this;
      if (!window.wails) return;

      window.wails.Events.On("tray.click.qrcode", function (ev) {
        $this.$router.push("/");
      });
      window.wails.Events.On("tray.click.settings", function (ev) {
        $this.$router.push("/settings");
      });
      window.wails.Events.On("tray.click.device", function (ev) {
        $this.$router.push("/device");
      });
    },
  },
};
</script>

<style>
</style>
