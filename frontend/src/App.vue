<template>
  <div class="container">
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
  name: "app",
  data() {
    return {
      message: "Hello Vue!",
      showQRCode: true,
      qrcodeStyle: {
        backgroundImage: null,
      },
    };
  },
  created: function () {
    this.generateQRCode();
  },
  methods: {
    toggleConnectWay: function () {
      this.showQRCode = !this.showQRCode;
    },
    generateQRCode: function () {
      var $this = this;
      window.backend &&
        window.backend.main.App.GenerateQRCode().then((result) => {
          $this.qrcodeStyle.backgroundImage = "url(" + result + ")";
        });
    },
    refreshQRCode: function () {
      var $this = this;
      window.backend &&
        window.backend.main.App.GenerateQRCode().then((result) => {
          $this.qrcodeStyle.backgroundImage = "url(" + result + ")";
        });
    },
    connect: function () {},
  },
};
</script>

<style>
html {
  text-align: center;
  color: white;
}

body {
  color: white;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen",
    "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
  margin: 0;
}

button {
  -webkit-appearance: default-button;
  padding: 6px;
  margin: auto;
  margin-top: 8px;
  display: block;
}

input {
  border: solid 1px #0066cc;
  padding: 5px;
}

.container {
  padding: 20px;
}

.container a {
  font-size: 12px;
  color: #0066cc;
  position: fixed;
  bottom: 30px;
  margin: auto 0;
  left: 0;
  right: 0;
}

.memo {
  color: #666;
  margin-top: 20px;
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
  background-image: url("./assets/qrcode.png");
}

#logo {
  width: 140px;
  height: 140px;
  margin: auto;
  display: block;
  background-position: center;
  background-repeat: no-repeat;
  background-image: url("./assets/appicon.png");
}
</style>
