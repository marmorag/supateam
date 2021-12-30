import { createApp } from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router.js";
import store from "./plugins/store.js";
import "./plugins/webfontloader";

router.store = store;

const app = createApp(App);

app.use(router).use(vuetify).use(store);
app.mount("#app");
