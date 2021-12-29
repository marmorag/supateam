const getDefault = () => ({
  user: null,
  authenticated: false,
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
      const { user, authenticated } = getDefault();
      state.user = user;
      state.authenticated = authenticated;
    },
  },
  actions: {
    authenticate({ commit }, { user }) {
      commit("authenticate", user);
    },
    disconnect({ commit }) {
      commit("disconnect");
    },
  },
  getters: {
    isAuthenticated: (state) => state.authenticated,
    getAuthenticated: (state) => state.user,
  },
};
