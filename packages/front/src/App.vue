<template>
  <v-layout>
    <v-app>
      <TheNavbar />
      <v-main>
        <v-container>
          <router-view />
        </v-container>
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
  store.dispatch("service/init");
})

</script>
