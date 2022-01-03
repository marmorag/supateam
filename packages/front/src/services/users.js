import { onMounted, ref } from "vue";
import UsersService from "./api/users";

export default function useUsers(store, triggerMounted = true) {
  const users = ref([]);
  const fetchUsersLoading = ref(false);
  const fetchUserParticipationsLoading = ref(false);

  const fetchUsers = async () => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new UsersService(client);

    fetchUsersLoading.value = true;
    const { status, data } = await eventsApi.get();
    fetchUsersLoading.value = false;
    if (status) {
      return data;
    }
    return null;
  };

  const fetchUserParticipations = async (id) => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new UsersService(client);

    fetchUserParticipationsLoading.value = true;
    const { status, data } = await eventsApi.getParticipations(id);
    fetchUserParticipationsLoading.value = false;
    if (status) {
      return data;
    }
    return [];
  };

  if (triggerMounted) {
    onMounted(async () => (users.value = await fetchUsers()));
  }

  return {
    users,
    fetchUsersLoading,
    fetchUserParticipationsLoading,
    fetchUsers,
    fetchUserParticipations,
  };
}
