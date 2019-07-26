import Vue from "https://cdn.jsdelivr.net/npm/vue/dist/vue.esm.browser.min.js";
import {notificationMixin} from "./notifications.js";

const app = new Vue({
  el: "#app",
  data: {
    newShortLink: "",
    newLongLink: "",
    shortLinks: window.shortLinks,
  },
  computed: {
    newShortLinkWidth: function() {
      let length = this.newShortLink.length + 1;
      if(length < 9) length = 9;
      if(length > 30) length = 30;
      return length;
    },
    createLinkDisabled: function() {
      return this.newShortLink === "" || this.newLongLink === "";
    }
  },
  methods: {
    createLink: async function () {
      if(this.createLinkDisabled) return;

      const id = this.newShortLink;
      const url = this.newLongLink;

      const res = await fetch("/api/link", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({id, url})
      });

      if(res.ok) {
        this.newShortLink = "";
        this.newLongLink = "";

        const body = await res.json();
        this.shortLinks.unshift(body);
        this.showNotification("Link shortened!");
      } else if(res.status === 409) {
        this.showNotification("A link with this short URL already exists!")
      } else {
        this.showNotification("An unexpected error occurred!");
      }
    },
    deleteLink: async function(id) {
      const res = await fetch("/api/link", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({id})
      });

      if(res.ok) {
        this.shortLinks = this.shortLinks.filter(shortLink => shortLink.id !== id);
        this.showNotification("Link deleted!");
      } else if(res.status === 404) {
        this.showNotification("Couldn't find link to delete!");
      } else {
        this.showNotification("An unexpected error occurred!");
      }
    }
  },
  mixins: [notificationMixin]
});

const clipboard = new ClipboardJS(".box.link", {
  text: function(trigger) {
    return window.location.protocol + "//" + trigger.dataset.link;
  }
});
clipboard.on("success", e => {
  app.showNotification(e.text.substring((window.location.protocol + "//").length) + " copied to the clipboard!");
});
