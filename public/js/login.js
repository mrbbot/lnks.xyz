import Vue from "https://cdn.jsdelivr.net/npm/vue/dist/vue.esm.browser.min.js";
import { notificationMixin } from "./notifications.js";

new Vue({
  el: "#app",
  mixins: [notificationMixin]
});


