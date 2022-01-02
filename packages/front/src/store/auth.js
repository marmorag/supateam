const getDefault = () => ({
  user: null,
  authenticated: false,
  authToken: null,
});

export default {
  namespaced: true,
  state: () => getDefault(),
  mutations: {
    authenticate(state, user) {
      state.user = user;
      state.authenticated = true;
    },
    disconnect(state) {
      const { user, authenticated, authToken } = getDefault();
      state.user = user;
      state.authenticated = authenticated;
      state.authToken = authToken;
    },
    setAuthToken(state, token) {
      state.authToken = token;
    },
  },
  actions: {
    authenticate({ commit }, { user, token }) {
      commit("authenticate", user);
      commit("setAuthToken", token);
    },
    disconnect({ commit }) {
      commit("disconnect");
    },
  },
  getters: {
    isAuthenticated: (state) => state.authenticated,
    getAuthenticated: (state) => state.user,
    getAuthToken: (state) => state.authToken,
  },
};
