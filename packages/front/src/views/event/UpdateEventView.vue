<template class="update-event-view">
  <ThePageTitle offset="2" />
  <EventForm ref="eventForm" v-model="event" kind="update" :event-id="id">
    <template #action>
      <v-btn class="ml-4" append-icon="mdi-check" color="blue lighten-1" text-color="white" @click="handleUpdateEvent">
        Mise à jour
      </v-btn>
      <v-btn :to="{ name: 'calendar' }">
        annuler
      </v-btn>
    </template>
  </EventForm>
</template>

<script setup>
import { defineProps, ref } from "vue";
import { useStore } from "vuex";
import ThePageTitle from "../../components/ThePageTitle.vue";
import EventForm from "../../components/events/EventForm.vue";
import useEvents from "../../services/events";
import { notify } from "@kyvg/vue3-notification";

const store = useStore();
const { updateEvent } = useEvents(store, false);

defineProps({
  id: {
    type: String,
    required: true,
  }
})

const eventForm = ref(null);
const event = ref({
  Title: "",
  Description: "",
  Date: new Date(),
  Duration: 1,
  Kind: "",
  Players: [],
  Teams: [],
});

const handleUpdateEvent = async () => {
  if (!await eventForm.value.validate()) {
    return;
  }

  const eventToUpdate = { ...event.value };
  eventToUpdate.Players = eventToUpdate.Players.map((player) => player.id);
  eventToUpdate.Teams = eventToUpdate.Teams.map((team) => team.id);

  const { status, data } = await updateEvent(eventToUpdate);
  if (!status) {
    notify({
      title: "Impossible de mettre à jour l'évènement.",
      type: "error"
    });
    throw Error(data)
  }

  notify({
    title: "L'évènement a bien été mis à jour.",
    type: "success"
  });
};
</script>

<style>
.v-radio-group > .v-selection-control-group {
  margin-top: 0 !important;
}
</style>