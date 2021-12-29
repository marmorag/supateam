import { onMounted, ref } from "vue";
import EventsService, {
  KIND_ENTRAINEMENT,
  KIND_EQUIPE,
  KIND_GRAND_PRIX,
} from "./api/events";

export default function useEvents(store, triggerMounted = true) {
  const events = ref([]);

  const kindColorMapping = {
    [KIND_GRAND_PRIX]: "#F44336",
    [KIND_ENTRAINEMENT]: "#2196F3",
    [KIND_EQUIPE]: "#4CAF50",
  };

  const eventKindList = [KIND_ENTRAINEMENT, KIND_EQUIPE, KIND_GRAND_PRIX];

  const fetchEvents = async () => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new EventsService(client);

    const { status, data } = await eventsApi.get();
    if (status) {
      events.value = data;
    }
  };

  const fetchEvent = async (id) => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new EventsService(client);

    const { status, data } = await eventsApi.get(id);
    if (status) {
      return data;
    }
    return null;
  };

  const fetchEventParticipations = async (id) => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new EventsService(client);

    const { status, data } = await eventsApi.getParticipation(id);
    if (status) {
      return data;
    }
    return [];
  };

  const createEvent = async (event) => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new EventsService(client);

    const { status, data } = await eventsApi.post(event);
    return { status, data };
  };

  const styleEventByKind = (kind) => ({
    backgroundColor: kindColorMapping[kind],
    width: "5px",
    height: "5px",
    borderRadius: "50%",
    display: "inline-block",
  });

  if (triggerMounted) {
    onMounted(fetchEvents);
  }

  return {
    events,
    eventKindList,
    kindColorMapping,
    fetchEvents,
    fetchEvent,
    fetchEventParticipations,
    createEvent,
    styleEventByKind,
  };
}
