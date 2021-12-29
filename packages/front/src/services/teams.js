import { onMounted, ref } from "vue";
import TeamsService from "./api/teams";

export default function useTeams(store, triggerMounted = true) {
  const teams = ref([]);

  const fetchTeams = async () => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new TeamsService(client);

    const { status, data } = await eventsApi.get();
    if (status) {
      teams.value = data;
    }
  };

  if (triggerMounted) {
    onMounted(fetchTeams);
  }

  return {
    teams,
    fetchTeams,
  };
}
