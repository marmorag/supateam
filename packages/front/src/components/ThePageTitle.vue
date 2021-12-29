<template>
  <v-row>
    <v-col :offset="offset">
      <h2>{{ currentRoute.meta.title }}</h2>
    </v-col>
  </v-row>
</template>

<script setup>
import { useRouter } from "vue-router";
import { defineProps } from "vue";

const { currentRoute } = useRouter();
defineProps({
  offset: {
    type: String,
    default: "1",
  }
});

</script>

<style scoped>

</style>