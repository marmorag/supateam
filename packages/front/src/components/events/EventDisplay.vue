<template>
  <v-row>
    <v-col cols="12">
      <v-row class="d-flex justify-space-between">
        <v-col md="6" sm="8">
          <h3>Nom de l'évènement : {{ event.title }}</h3>
        </v-col>
        <v-col md="2" sm="6" align-self="end">
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
                <EventParticipation :disabled="!canUpdateParticipation" :participation="simplify(participation)" @participation:update="handleParticipationUpdated"/>
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
import { computed, defineEmits, defineProps } from "vue";
import { format } from "date-fns";
import { fr } from "date-fns/locale";
import EventParticipation from "./EventParticipation.vue";
import useEvents from "../../services/events";
import useAuthorization from "../../services/authorization";
import store from "../../plugins/store";

const { authorize, PARTICIPATIONS_API_GROUP, WRITE_API_ACTION } = useAuthorization(store)
const { kindColorMapping } = useEvents(null, false)

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

const emit = defineEmits(["participation:update"])

const canUpdateParticipation = computed(() => authorize({ api: PARTICIPATIONS_API_GROUP, action: WRITE_API_ACTION }) || false)

const cParticipations = computed(() => {
  return props.participations ?? [];
})

const formatDate = (date) => {
  return format(Date.parse(date), "EEEE d LLLL", { locale: fr });
}

const handleParticipationUpdated = () => {
  emit("participation:update")
}

const simplify = (participation) => ({
  ...participation,
  player: participation.player[0].id,
  team: participation.team[0]?.id ?? null,
})
</script>

<style scoped>
td, th {
  text-align: center !important;
}

</style>