<template>
  <ThePageTitle />
  <v-row>
    <v-col md="6" offset-md="1" sm="8" offset-sm="2">
      <v-row>
        <v-col cols="12">
          <Calendar ref="calendar" :attributes="calendarAttributes" locale="fr" is-expanded @update:to-page="handleToPage">
            <template #day-popover="{ day, format, attributes }">
              <v-row class="mb-2">
                <v-col class="pb-0">
                  {{ format(day.date, "WWWW D MMMM") }}
                </v-col>
              </v-row>
              <v-row v-for="attribute in attributes" :key="attribute.key">
                <v-col>
                  <i class="pa-1" :style="styleEventByKind(attribute.customData.kind)" /> {{ attribute.popover.label }} ({{ attribute.customData.kind }})
                </v-col>
              </v-row>
            </template>
          </Calendar>
        </v-col>
      </v-row>
    </v-col>
    <v-col md="3" sm="8">
      <EventList :events="events" :participations="participations" :month="currentMonth" :year="currentYear" @participation:refresh="syncParticipation" />
    </v-col>
  </v-row>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useStore } from "vuex";
import { Calendar } from "v-calendar";
import useEvents from "../services/events";
import EventList from "../components/events/EventList.vue";
import ThePageTitle from "../components/ThePageTitle.vue";
import useUsers from "../services/users";

const STORAGE_KEY = "calendar:current";

const store = useStore();
const { events, kindColorMapping, styleEventByKind } = useEvents(store, true);
const { fetchUserParticipations } = useUsers(store, false);

const currentMonth = ref((new Date()).getMonth() + 1);
const currentYear = ref((new Date()).getFullYear());
const participations = ref([]);
const calendar = ref(null);

const calendarAttributes = computed(() => {
  return events.value.map((event) => ({
    key: event.id,
    bar: {
      style: {
        backgroundColor: kindColorMapping[event.kind]
      }
    },
    dates: {
      start: new Date(event.date),
      span: event.duration,
    },
    popover: {
      label: event.title,
      visibility: "click"
    },
    customData: {
      ...event
    },
  }))
});

const syncParticipation = async () => {
  const userId = store.getters["auth/getAuthenticated"].userId
  participations.value = await fetchUserParticipations(userId)
}

const handleToPage = ({ month, year }) => {
  currentMonth.value = month;
  currentYear.value = year;
}

onMounted(syncParticipation)

onMounted(() => {
  const resolvedStorage = window.localStorage.getItem(STORAGE_KEY)
  if (resolvedStorage !== null) {
    const currentCalendar = JSON.parse(resolvedStorage)
    currentYear.value = currentCalendar.currentYear
    currentMonth.value = currentCalendar.currentMonth
    calendar.value.move({ month: currentMonth.value, year: currentYear.value })
  }
})

watch([currentMonth, currentYear], ([newCurrentMonth, newCurrentYear]) => {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify({ currentMonth: newCurrentMonth, currentYear: newCurrentYear }))
})
</script>