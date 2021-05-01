<template>
  <div class="container">
    <div class="message">
      <div class="textarea-wapper">
        <textarea
          name="textarea"
          class="p-form-text p-form-no-validate"
          v-model="message"
        ></textarea>
      </div>
      <div class="toolbar">
        <div class="left">
          <button class="p-btn" v-on:click="clear">清空</button>
        </div>
        <div class="right">
          <button class="p-btn" v-on:click="cancel">取消</button>
          <button class="p-btn p-prim-col" v-on:click="submit">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script scoped>
export default {
  data() {
    return {
      message: "",
    };
  },
  created: function () {
    this.loadClipboardText();
  },
  methods: {
    loadClipboardText: function () {
      var $this = this;
      window.backend &&
        window.backend.main.App.GetClipboardText().then((result) => {
          $this.message = result;
        });
    },
    clear: function () {
      this.message = "";
    },
    cancel: function () {
      window.wails.Events.Emit("app.window.close");
    },
    submit: function () {
      var $this = this;
      window.backend &&
        window.backend.main.App.SendText($this.message).then((result) => {
          console.log(result);
          window.wails.Events.Emit("app.window.close");
        });
    },
    connect: function () {},
  },
};
</script>

<style scoped>
.message {
  margin: 0 20px;
}

textarea {
  width: 100%;
  -webkit-box-sizing: border-box;
  -moz-box-sizing: border-box;
  box-sizing: border-box;
  padding: 5px;
  margin: 0px;
  height: 180px;
  font-size: 14px;
}
.toolbar {
  display: flex;
  flex-direction: row;
}
.toolbar > .left {
  text-align: left;
  flex: auto;
}
.toolbar > .right {
  flex: auto;
  text-align: right;
}
</style>
