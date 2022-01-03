<template>
  <h3>Liste des évènements</h3>
  <v-expansion-panels v-if="monthEvent.length > 0">
    <v-expansion-panel v-for="event in monthEvent" :key="event.id">
      <v-expansion-panel-title>
        {{ event.title }}
      </v-expansion-panel-title>
      <v-expansion-panel-text>
        <v-row>
          <v-col cols="12">
            <v-row>
              <v-col class="d-flex align-center">
                <EventKindChip :kind="event.kind" />
              </v-col>
              <v-col class="d-flex flex-row-reverse">
                <v-btn v-if="canEditEvent" icon="mdi-calendar-edit" size="small" class="mr-2" :to="{ name: 'update-event', params: { id: event.id } }" />
                <v-btn icon="mdi-calendar-search" size="small" :to="{ name: 'event', params: { id: event.id } }" />
              </v-col>
            </v-row>
          </v-col>
          <v-col cols="12">
            Date : {{ formatDate(event) }} - {{ event.duration }} jour(s)
          </v-col>
          <v-col>
            Notes : {{ event.description.length === 0 ? 'Aucune.' : event.description }}
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <EventParticipation
              :participation="getParticipationForEvent(event.id)"
              @participation:create="(payload) => handleParticipationCreate(payload, event.id)"
              @participation:update="handleParticipationUpdate"
            />
          </v-col>
        </v-row>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
  <v-row v-else>
    <v-col>
      <h4>Aucuns évènements programmé ce mois.</h4>
    </v-col>
    <v-col v-if="canCreateEvent">
      <v-btn :to="{ name: 'create-event' }" prepend-icon="mdi-calendar-plus">Ajouter</v-btn>
    </v-col>
  </v-row>
</template>

<script setup>
import { computed, defineProps, defineEmits } from "vue";
import { isSameMonth, isSameYear, setMonth, setYear, format } from "date-fns";
import { fr } from "date-fns/locale";
import { useStore } from "vuex";
import useAuthorization from "../../services/authorization";
import EventKindChip from "./EventKindChip.vue";
import EventParticipation from "./EventParticipation.vue";
import useParticipations from "../../services/participations";

const store = useStore();
const { authorize, EVENTS_API_GROUP, WRITE_API_ACTION } = useAuthorization(store);
const { createParticipation } = useParticipations(store);

const emit = defineEmits(['participation:refresh']);

const props = defineProps({
  events: {
    type: Array,
    required: true,
  },
  participations: {
    type: Array,
    required: false,
    default: () => [],
  },
  month: {
    type: Number,
    required: true,
  },
  year: {
    type: Number,
    required: true,
  },
});

const monthEvent = computed(() => {
  const events = [...props.events];
  const currentDate = setYear(setMonth(new Date(), props.month - 1), props.year)

  return events.filter((event) => {
    const eventDate = new Date(event.date);
    return isSameYear(currentDate, eventDate) && isSameMonth(currentDate, eventDate);
  });
})
const canCreateEvent = computed(() => authorize({ api: EVENTS_API_GROUP, action: WRITE_API_ACTION }) || false)
const canEditEvent = computed(() => authorize({ api: EVENTS_API_GROUP, action: WRITE_API_ACTION }) || false)

const formatDate = (event) => {
  return format(new Date(event.date), "EEEE d LLLL", { locale: fr })
}

const getParticipationForEvent = (eventId) => {
  return props.participations.find((participation) => participation.event === eventId) ?? {};
}

const handleParticipationCreate = async ({ status }, eventId) => {
  const participation = {
    Event: eventId,
    Player: store.getters["auth/getAuthenticated"].userId,
    Team: null,
    Status: status,
  };

  const { status: responseStatus } = await createParticipation(participation);
  if (responseStatus) {
    emit('participation:refresh');
  }
}

const handleParticipationUpdate = () => {
  emit('participation:refresh');
}
</script>