<template>
  <div class="container">
    <div class="settings">
      <p>
        <label class="p-form-checkbox-cont">
          <input
            type="checkbox"
            name="example"
            v-model="startAtLogin.value"
            :disabled="startAtLogin.disabled"
            v-on:click="toggleClipboardSync"
          />
          <span> </span>
          {{ startAtLogin.text }}
        </label>
      </p>
      <p>
        <label class="p-form-checkbox-cont">
          <input
            type="checkbox"
            name="example"
            v-model="config.ClipboardSync"
            v-on:click="toggleClipboardSync"
          />
          <span></span>
          同步剪贴板到手机
        </label>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      config: {
        StartAtLogin: false,
        HideOnStartup: false,
        ConnectOnStartup: false,
        ClipboardSync: false,
      },
      startAtLogin: {
        text: "开机自启动",
        value: false,
        disabled: false,
      },
    };
  },
  created: function () {
    this.loadConfig();
  },
  methods: {
    loadConfig: function () {
      var $this = this;
      if (!window.backend) return;

      window.backend.main.App.LoadConfig().then((result) => {
        console.log(result);
        $this.config.ClipboardSync = result.ClipboardSync;
      });

      window.backend.main.App.GetStartAtLoginState()
        .then((result) => {
          console.log(result);
          $this.startAtLogin.value = result;
        })
        .catch((e) => {
          console.log(e);
          $this.startAtLogin.text = "开机自启动(不支持)";
          $this.startAtLogin.value = false;
          $this.startAtLogin.disabled = true;
        });
    },
    toggleStartAtLogin: function (ev) {
      var $this = this;
      window.backend.main.App.toggleStartAtLogin(ev.target.checked)
        .then((result) => {
          console.log(result);
        })
        .catch((e) => {
          console.log(e);
        });
    },
    toggleClipboardSync: function (ev) {
      var $this = this;
      console.log(ev.target.checked);
      window.backend.main.App.ToggleClipboardSync(ev.target.checked)
        .then((result) => {
          console.log(result);
        })
        .catch((e) => {
          console.log(e);
        });
    },
  },
};
</script>

<style scoped>
.settings {
  text-align: left;
  margin-left: 30px;
}
</style>
