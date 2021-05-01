<template>
  <div class="container" v-show="showpage == 'index'">
    <div id="qrcode">
      <div
        id="qrcode_img"
        v-bind:style="qrcodeStyle"
        v-on:click="refreshQRCode"
      ></div>
    </div>
    <div class="memo">请使用app扫码连接</div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      message: "Hello Vue!",
      showQRCode: true,
      qrcodeStyle: {
        backgroundImage: null,
      },
      showpage: "index",
      rawData: null,
    };
  },
  created: function () {
    this.generateQRCode();
    this.registerEvent();
  },
  methods: {
    toggleConnectWay: function () {
      this.showQRCode = !this.showQRCode;
    },
    generateQRCode: function () {
      var $this = this;
      window.backend &&
        window.backend.main.App.GenerateQRCode().then((result) => {
          $this.rawData = result.data;
          $this.qrcodeStyle.backgroundImage = "url(" + result.base64data + ")";
        });
    },
    refreshQRCode: function () {
      var $this = this;
      window.backend &&
        window.backend.main.App.GenerateQRCode().then((result) => {
          $this.rawData = result.data;
          $this.qrcodeStyle.backgroundImage = "url(" + result.base64data + ")";
        });
    },
    connect: function () {},
    registerEvent: function () {
      var $this = this;
      if (!window.wails) return;

      window.wails.Events.On("app.scan.connected", function (code) {
        if (code.indexOf($this.rawData.token) >= 0) {
          window.wails.Events.Emit("app.window.close");
        }
      });

      window.wails.Events.On("app.qrcode.window.show", function () {
        console.log("app.qrcode.window.show");
        $this.refreshQRCode();
      });
    },
  },
};
</script>

<style scoped>
.memo {
  color: #666;
  margin-top: 20px;
  text-align: center;
}

#name {
  border-radius: 3px;
  outline: none;
  -webkit-font-smoothing: antialiased;
}

#qrcode {
  width: 200px;
  height: 200px;
  background-color: #fff;
  padding: 5px;
  margin: auto;
  display: block;
}
#qrcode_img {
  width: 100%;
  height: 100%;
  margin: auto;
  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-image: url(./assets/qrcode.png);
}

#logo {
  width: 140px;
  height: 140px;
  margin: auto;
  display: block;
  background-position: center;
  background-repeat: no-repeat;
  background-image: url(./assets/appicon.png);
}
</style>
