import { createApp } from 'vue/dist/vue.esm-bundler.js';
import { createRouter, createMemoryHistory } from "vue-router";
import './puppertino.css'
import './style.css'


import App from "./App.vue"
import Home from "./Home.vue";
import Device from "./Device.vue";
import Settings from "./Settings.vue";

const routes = [
  { path: "/", name:"home", component: Home },
  { path: "/device", name:"device", component: Device },
  { path: "/settings", name:"settings", component: Settings },
];

const router = createRouter({
  history: createMemoryHistory(),
  routes,
});

createApp(App).use(router).mount("#app");