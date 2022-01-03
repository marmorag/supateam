import { createApp } from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router.js";
import store from "./plugins/store.js";
import "./plugins/webfontloader";
import Notifications from "@kyvg/vue3-notification";

router.store = store;

const app = createApp(App);

app.use(router).use(vuetify).use(store).use(Notifications);
app.mount("#app");
