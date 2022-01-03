<template>
  <v-chip :color="kindColorMapping[kind]" text-color="white" >{{ kind }}</v-chip>
</template>

<script setup>
import { defineProps } from "vue";
import useEvents from "../../services/events";

const { kindColorMapping } = useEvents(null, false);

defineProps({
  kind: {
    type: String,
    required: true,
  }
})
</script>