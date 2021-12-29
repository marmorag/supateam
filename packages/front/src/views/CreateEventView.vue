<template class="create-event-view">
  <ThePageTitle offset="2" />
  <v-row>
    <v-col cols="8" offset="2" class="pa-5">
      <v-row>
        <v-col cols="4" class="pa-0">
          <v-text-field
            v-model="event.Title"
            label="Nom de l'évènement"
            :error="v$.Title.$error"
            @blur="v$.Title.$touch()"
            @input="v$.Title.$touch()"
          />
        </v-col>
        <v-col cols="8" class="pa-0">
          <v-radio-group
            v-model="event.Kind"
            label="Type d'évènement"
            :inline="true"
            :error="v$.Kind.$error"
            @blur="v$.Kind.$touch()"
            @input="v$.Kind.$touch()"
          >
            <v-radio  v-for="(kind, index) in eventKindList" :key="`kind-${index}`" :label="kind" :value="kind" />
          </v-radio-group>
        </v-col>
        <v-col cols="8" class="py-0 pl-0">
          <DatePicker v-model="event.Date" mode="date" >
            <template #default="{ inputValue, inputEvents }">
              <v-text-field
                :model-value="inputValue"
                label="Date de l'évènement"
                :error="v$.Date.$error"
                v-on="inputEvents"
                @blur="v$.Date.$touch()"
                @input="v$.Date.$touch()"
              />
            </template>
          </DatePicker>
        </v-col>
        <v-col cols="4" class="py-0 pr-0">
          <v-text-field
            v-model="event.Duration"
            label="Durée (jours)"
            type="number"
            :error="v$.Duration.$error"
            @blur="v$.Duration.$touch()"
            @input="v$.Duration.$touch()"
          />
        </v-col>
        <v-col cols="12" class="pa-0">
          <v-textarea
            v-model="event.Description"
            label="Notes pour l'évènement"
            :error="v$.Description.$error"
            @blur="v$.Description.$touch()"
            @input="v$.Description.$touch()"
          />
        </v-col>
        <v-col cols="6" class="py-0 pl-0">
          <ChipSelection v-model="event.Players" selector-key="name" label="Participants" :selectable="users" />
        </v-col>
        <v-col cols="6" class="py-0 pr-0">
          <ChipSelection v-model="event.Teams" selector-key="name" label="Équipes" :selectable="teams" />
        </v-col>
      </v-row>
      <v-row class="d-flex flex-row-reverse">
        <v-col cols="4" class="d-flex flex-row-reverse pa-0">
          <v-btn class="ml-4" append-icon="mdi-check" color="primary" @click="handleCreateEvent">
            créer
          </v-btn>
          <v-btn @click="router.back()">
            annuler
          </v-btn>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script setup>
import { ref } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import { DatePicker } from "v-calendar";
import useVuelidate from "@vuelidate/core";
import { integer, required } from "@vuelidate/validators";
import ThePageTitle from "../components/ThePageTitle.vue";
import ChipSelection from "../components/ChipSelection.vue";
import useEvents from "../services/events";
import useUsers from "../services/users";
import useTeams from "../services/teams";

const router = useRouter();
const store = useStore();
const { eventKindList, createEvent } = useEvents(store, false);
const { users } = useUsers(store);
const { teams } = useTeams(store);

const eventKind = (value) => eventKindList.includes(value);

const event = ref({
  Title: "",
  Description: "",
  Date: new Date(),
  Duration: 1,
  Kind: "",
  Players: [],
  Teams: [],
});

const rules = {
  Title: { required },
  Description: {},
  Date: { required },
  Duration: { required, integer },
  Kind: { required, eventKind }
};
const v$ = useVuelidate(rules, event);

const handleCreateEvent = async () => {
  if (!await v$.value.$validate()) {
    return;
  }

  const eventToCreate = { ...event.value };
  eventToCreate.Players = eventToCreate.Players.map((player) => player.id);
  eventToCreate.Teams = eventToCreate.Teams.map((team) => team.id);

  const { status, data } = await createEvent(eventToCreate);
  if (!status) {
    console.log(data);
  }
};
</script>

<style>
.v-radio-group > .v-selection-control-group {
  margin-top: 0 !important;
}
</style>