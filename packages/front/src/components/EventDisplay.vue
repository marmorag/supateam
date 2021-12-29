<template>
  <v-row>
    <v-col cols="12">
      <v-row class="d-flex justify-space-between">
        <v-col cols="6">
          <h3>Nom de l'évènement : {{ event.title }}</h3>
        </v-col>
        <v-col cols="2" align-self="end">
          <v-chip :color="kindColorMapping[event.kind]" text-color="white" >{{ event.kind }}</v-chip>
        </v-col>
      </v-row>
    </v-col>
    <v-col cols="12">
      {{ formatDate(event.date) }} - {{ event.duration }} jour(s)
    </v-col>
    <v-col cols="12">
      Description :<br>
      {{ event.description.length > 0 ? event.description : "Aucune notes." }}
    </v-col>
    <v-col cols="12">
      <v-table :fixed-header="true">
        <thead>
          <tr>
            <th>Joueur</th>
            <th>Equipe</th>
            <th>Statut</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="cParticipations.length > 0">
            <tr v-for="participation in cParticipations" :key="participation.id">
              <td>{{ participation.player[0].name }}</td>
              <td>{{ participation.team[0]?.name }}</td>
              <td>
                <v-icon>mdi-{{ statusIcon(participation.status) }}</v-icon>
              </td>
            </tr>
          </template>
          <template v-else>
            <tr>
              <td colspan="3">Aucuns joueur.</td>
            </tr>
          </template>
        </tbody>
      </v-table>
    </v-col>
  </v-row>
</template>

<script setup>
import { computed, defineProps } from "vue";
import { format } from "date-fns";
import { fr } from "date-fns/locale";
import useParticipations from "../services/participations";
import useEvents from "../services/events";

const { kindColorMapping } = useEvents(null, false)
const { participationStatusStyleMapping } = useParticipations();

const props = defineProps({
  event: {
    type: Object,
    required: true,
  },
  participations: {
    type: Array,
    required: true,
  },
})

const cParticipations = computed(() => {
  return props.participations ?? [];
})

const statusIcon = (status) => {
  return participationStatusStyleMapping[status].icon;
};

const formatDate = (date) => {
  return format(Date.parse(date), "EEEE d LLLL", { locale: fr });
}
</script>