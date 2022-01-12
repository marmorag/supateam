<template class="create-event-view">
  <ThePageTitle offset="2" />
  <EventForm ref="eventForm" v-model="event" kind="create">
    <template #action>
      <v-btn class="ml-4" append-icon="mdi-check" color="blue lighten-1" text-color="white" @click="handleCreateEvent">
        créer
      </v-btn>
      <v-btn @click="router.back()">
        annuler
      </v-btn>
    </template>
  </EventForm>
</template>

<script setup>
import { ref } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import ThePageTitle from "../../components/ThePageTitle.vue";
import EventForm from "../../components/events/EventForm.vue";
import useEvents from "../../services/events";
import { notify } from "@kyvg/vue3-notification";

const router = useRouter();
const store = useStore();
const { createEvent } = useEvents(store, false);

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

const handleCreateEvent = async () => {
  console.log(eventForm);
  if (!await eventForm.value.validate()) {
    return;
  }

  const eventToCreate = { ...event.value };
  eventToCreate.Players = eventToCreate.Players.map((player) => player.id);
  eventToCreate.Teams = eventToCreate.Teams.map((team) => team.id);

  const { status, data } = await createEvent(eventToCreate);
  if (!status) {
    notify({
      title: "Impossible de créer l'évènement.",
      type: "error"
    });
    throw Error(data)
  }

  notify({
    title: "L'évènement a bien été créé.",
    type: "success"
  });
  eventForm.value.resetForm();
};
</script>

<style>
.v-radio-group > .v-selection-control-group {
  margin-top: 0 !important;
}
</style>