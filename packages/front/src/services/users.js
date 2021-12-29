import { onMounted, ref } from "vue";
import UsersService from "./api/users";

export default function useUsers(store, triggerMounted = true) {
  const users = ref([]);

  const fetchUsers = async () => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new UsersService(client);

    const { status, data } = await eventsApi.get();
    if (status) {
      users.value = data;
    }
  };

  const fetchUserParticipations = async (id) => {
    const client = store.getters["service/apiClient"];
    const eventsApi = new UsersService(client);

    const { status, data } = await eventsApi.getParticipations(id);
    if (status) {
      return data;
    }
    return [];
  };

  if (triggerMounted) {
    onMounted(fetchUsers);
  }

  return {
    users,
    fetchUsers,
    fetchUserParticipations,
  };
}
