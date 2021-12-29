import EventsService from "./api/participations";
import ParticipationsService from "./api/participations";

export default function useParticipations(store) {
  const PARTICIPATION_STATUS_UNKNOWN = "unknown";
  const PARTICIPATION_STATUS_ACCEPTED = "accepted";
  const PARTICIPATION_STATUS_REJECTED = "rejected";

  const participationStatusList = [
    PARTICIPATION_STATUS_UNKNOWN,
    PARTICIPATION_STATUS_REJECTED,
    PARTICIPATION_STATUS_ACCEPTED,
  ];

  const participationStatusStyleMapping = {
    [PARTICIPATION_STATUS_UNKNOWN]: {
      icon: "help",
      iconColor: "grey darken-1",
      selectedBgColor: "grey darken-1",
      selectedIconColor: "white",
    },
    [PARTICIPATION_STATUS_ACCEPTED]: {
      icon: "check",
      iconColor: "green darken-1",
      selectedBgColor: "green darken-1",
      selectedIconColor: "white",
    },
    [PARTICIPATION_STATUS_REJECTED]: {
      icon: "cancel",
      iconColor: "red darken-1",
      selectedBgColor: "red darken-1",
      selectedIconColor: "white",
    },
  };

  const createParticipation = async (participation) => {
    const client = store.getters["service/apiClient"];
    const participationsApi = new ParticipationsService(client);

    const { status, data } = await participationsApi.post(participation);
    return { status, data };
  };

  const updateParticipation = async (participation) => {
    const client = store.getters["service/apiClient"];
    const participationsApi = new ParticipationsService(client);

    const { status, data } = await participationsApi.put(participation);
    return { status, data };
  };

  return {
    PARTICIPATION_STATUS_UNKNOWN,
    PARTICIPATION_STATUS_ACCEPTED,
    PARTICIPATION_STATUS_REJECTED,
    participationStatusStyleMapping,
    participationStatusList,
    createParticipation,
    updateParticipation,
  };
}
