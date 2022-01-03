<template>
  <ThePageTitle offset="2" />
  <v-row>
    <v-col md="8" offset-md="2" sm="10" offset-sm="1">
      <EventDisplay v-if="loaded" :event="event" :participations="participations" @participation:update="handleParticipationUpdate"/>
    </v-col>
  </v-row>
</template>

<script setup>
import { defineProps, onMounted, ref } from "vue";
import ThePageTitle from "../components/ThePageTitle.vue";
import { useStore } from "vuex";
import useEvents from "../services/events";
import EventDisplay from "../components/EventDisplay.vue";

const store = useStore();
const { fetchEvent, fetchEventParticipations } = useEvents(store, false);

const event = ref(null);
const participations = ref(null);
const loaded = ref(false);

const props = defineProps({
  id: {
    type: String,
    required: true,
  },
});

onMounted(async () => {
  Promise.all([
    fetchEvent(props.id),
    fetchEventParticipations(props.id),
  ])
    .then(([eventResponse, participationsResponse ]) => {
      event.value = eventResponse;
      participations.value = participationsResponse
    })
    .finally(() => {
      loaded.value = true;
    })
})

const handleParticipationUpdate = async () => {
  participations.value = await fetchEventParticipations(props.id)
}

</script>
