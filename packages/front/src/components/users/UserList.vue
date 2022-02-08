<template>
  <v-row>
    <v-col cols="12">
      <v-table :fixed-header="true">
        <thead>
        <tr>
          <th>Nom</th>
          <th>Equipe</th>
          <th>Statut</th>
        </tr>
        </thead>
        <tbody>

        </tbody>
      </v-table>
    </v-col>
  </v-row>
</template>

<script setup>
import { defineProps } from "vue";

defineProps({
  users: {
    type: Array,
    required: true,
  },
})
</script>

<style scoped>

</style>