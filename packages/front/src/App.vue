<template>
  <v-layout>
    <v-app>
      <TheNavbar />
      <v-main>
        <v-container>
          <router-view />
        </v-container>
        <notifications :max="3"/>
      </v-main>
    </v-app>
  </v-layout>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useStore } from "vuex";
import TheNavbar from "./components/TheNavbar.vue";

onBeforeMount(() => {
  const store = useStore();

  if (store.getters["auth/isAuthenticated"]) {
    const token = store.getters["auth/getAuthToken"]
    store.dispatch("service/setAuthenticationToken", { token })
  } else {
    store.dispatch("service/init");
  }
})

</script>
