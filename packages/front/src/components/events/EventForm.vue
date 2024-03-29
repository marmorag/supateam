<template>
  <v-row>
    <v-col md="8" offset-md="2" sm="10" offset-sm="1" class="pa-5">
      <v-row :class="{'d-flex': isMobile, 'flex-column': isMobile }">
        <v-col md="4" sm="10" class="pa-0">
          <v-text-field
            v-model="event.Title"
            label="Nom de l'évènement"
            :error="v$.Title.$error"
            @blur="v$.Title.$touch()"
            @input="v$.Title.$touch()"
          />
        </v-col>
        <v-col md="8" sm="10" class="pa-0">
          <v-radio-group
            v-model="event.Kind"
            label="Type d'évènement"
            :inline="true"
            :error="v$.Kind.$error"
            @blur="v$.Kind.$touch()"
            @input="v$.Kind.$touch()"
          >
            <v-radio  v-for="(eventkind, index) in eventKindList" :key="`kind-${index}`" :label="eventkind" :value="eventkind" />
          </v-radio-group>
        </v-col>
        <v-col md="8" sm="10" :class="{'py-0': true, 'pl-0': !isMobile }">
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
        <v-col md="4" sm="10" :class="{'py-0': true, 'pr-0': !isMobile }">
          <v-text-field
            v-model="event.Duration"
            label="Durée (jours)"
            type="number"
            :error="v$.Duration.$error"
            @blur="v$.Duration.$touch()"
            @input="v$.Duration.$touch()"
          />
        </v-col>
        <v-col md="12" sm="12" class="pa-0">
          <v-textarea
            v-model="event.Description"
            label="Notes pour l'évènement"
            :error="v$.Description.$error"
            @blur="v$.Description.$touch()"
            @input="v$.Description.$touch()"
          />
        </v-col>
        <v-col md="6" sm="12" :class="{'py-0': true, 'pl-0': !isMobile }">
          <ChipSelection v-model="event.Players" selector-key="name" label="Participants" :selectable="users" />
        </v-col>
        <v-col md="6" sm="12" :class="{'py-0': true, 'pr-0': !isMobile }">
          <ChipSelection v-model="event.Teams" selector-key="name" label="Équipes" :selectable="teams" />
        </v-col>
      </v-row>
      <v-row class="d-flex flex-row-reverse">
        <v-col cols="4" class="d-flex flex-row-reverse pa-0">
          <slot name="action"></slot>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script setup>
import { computed, defineEmits, defineProps, onMounted } from "vue";
import { useStore } from "vuex";
import { DatePicker } from "v-calendar";
import useVuelidate from "@vuelidate/core";
import { integer, required } from "@vuelidate/validators";
import ChipSelection from "../../components/ChipSelection.vue";
import useEvents from "../../services/events";
import useUsers from "../../services/users";
import useTeams from "../../services/teams";
import useAppDisplay from "../../services/display";

const { isMobile } = useAppDisplay();
const store = useStore();
const { eventKindList } = useEvents(store, false);
const { users, fetchUsers } = useUsers(store, false);
const { teams, fetchTeams } = useTeams(store, false);
const { fetchEvent } = useEvents(store, false);

const eventKind = (value) => eventKindList.includes(value);

const emit = defineEmits(["update:modelValue"]);

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  eventId: {
    type: String,
    required: false,
    default: "",
  },
  kind: {
    type: String,
    required: true,
    validator: (value) => ["update", "create"].includes(value),
  },
});

const event = computed({
  get: () => props.modelValue,
  set: (value) => {
    emit('update:modelValue', value)
  },
});

onMounted(() => {
  Promise.all([
    fetchUsers(),
    fetchTeams(),
    props.kind === "update" ? fetchEvent(props.eventId) : async () => ({}),
  ]).then(([fetchedUsers, fetchedTeams, fetchedEvent]) => {
    users.value = fetchedUsers
    teams.value = fetchedTeams

    if (props.kind === "update") {
      event.value = {
        Id: fetchedEvent.id,
        Title: fetchedEvent.title,
        Description: fetchedEvent.description,
        Date: new Date(fetchedEvent.date),
        Duration: fetchedEvent.duration,
        Kind: fetchedEvent.kind,
        Players: fetchedEvent.players.map((p) => fetchedUsers.find((fp) => p === fp.id)),
        Teams: fetchedEvent.teams.map((t) => fetchedTeams.find((ft) => t === ft.id)),
      }
    }
  })
})

const rules = {
  Title: { required },
  Description: {},
  Date: { required },
  Duration: { required, integer },
  Kind: { required, eventKind }
};
const v$ = useVuelidate(rules, event);

const validate = async () => {
  return await v$.value.$validate();
};

const resetForm = () => {
  event.value.Title = ""
  event.value.Description = ""
  event.value.Date = new Date()
  event.value.Duration = 1
  event.value.Kind = ""
  event.value.Players = []
  event.value.Teams = []

  v$.value.$reset()
};

defineExpose({
  validate,
  resetForm,
})
</script>

<style>
.v-radio-group > .v-selection-control-group {
  margin-top: 0 !important;
}
</style>