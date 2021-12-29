import axios from "axios";

export default {
  namespaced: true,
  state: () => ({
    apiClient: null,
  }),
  mutations: {
    initApiClient(state, client) {
      state.apiClient = client;
    },
  },
  actions: {
    init({ commit }) {
      const client = axios.create({
        baseURL: import.meta.env.VITE_API_URL,
      });

      commit("initApiClient", client);
    },
    setAuthenticationToken({ commit }, { token }) {
      const client = axios.create({
        baseURL: import.meta.env.VITE_API_URL,
        headers: {
          Authorization: `bearer ${token}`,
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      });

      commit("initApiClient", client);
    },
    clearAuthenticationToken({ dispatch }) {
      dispatch("init");
    },
  },
  getters: {
    apiClient: (state) => state.apiClient,
  },
};
