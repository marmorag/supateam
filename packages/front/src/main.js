import { createApp } from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router.js";
import store from "./plugins/store.js";
import "./plugins/webfontloader";

router.store = store;

const app = createApp(App);

// broken
// if (import.meta.env.PROD) {
//   Sentry.init({
//     app,
//     environment: "production",
//     dsn: "https://ed98515495174d489e951f8f1f1246cc@o473284.ingest.sentry.io/6129020",
//     integrations: [
//       new Integrations.BrowserTracing({
//         routingInstrumentation: Sentry.vueRouterInstrumentation(router),
//         tracingOrigins: ["supateam.marmog.cloud", /^\//],
//       }),
//     ],
//     // Set tracesSampleRate to 1.0 to capture 100%
//     // of transactions for performance monitoring.
//     // We recommend adjusting this value in production
//     tracesSampleRate: 1.0,
//   });
// }

app.use(router).use(vuetify).use(store);
app.mount("#app");
