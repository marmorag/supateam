import { createStore } from "vuex";
import VuexPersistence from "vuex-persist";
import authModule from "../store/auth";
import serviceModule from "../store/service";

const vuexLocal = new VuexPersistence({
  storage: window.localStorage,
  modules: ["auth"],
});

export default createStore({
  modules: {
    auth: authModule,
    service: serviceModule,
  },
  plugins: [vuexLocal.plugin],
});
