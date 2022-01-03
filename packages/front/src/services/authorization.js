export default function useAuthorization(store) {
  const EVENTS_API_GROUP = "events";
  const USERS_API_GROUP = "users";
  const PARTICIPATIONS_API_GROUP = "participations";
  const TEAMS_API_GROUP = "teams";

  const validApiGroup = [
    EVENTS_API_GROUP,
    USERS_API_GROUP,
    PARTICIPATIONS_API_GROUP,
    TEAMS_API_GROUP,
  ];

  const ALL_API_ACTION = "*";
  const READ_API_ACTION = "read";
  const READ_SELF_API_ACTION = "read:self";
  const WRITE_API_ACTION = "write";
  const WRITE_SELF_API_ACTION = "write:self";
  const UPDATE_API_ACTION = "update";
  const UPDATE_SELF_API_ACTION = "update:self";
  const DELETE_API_ACTION = "delete";

  const validApiAction = [
    ALL_API_ACTION,
    READ_API_ACTION,
    READ_SELF_API_ACTION,
    WRITE_API_ACTION,
    WRITE_SELF_API_ACTION,
    UPDATE_API_ACTION,
    UPDATE_SELF_API_ACTION,
    DELETE_API_ACTION,
  ];

  const unproxyfy = (proxy) => JSON.parse(JSON.stringify(proxy));

  const authorize = ({ api, action }) => {
    if (
      !store.getters["auth/isAuthenticated"] ||
      !validApiGroup.includes(api) ||
      !validApiAction.includes(action)
    ) {
      return false;
    }

    const userAuthorizations = unproxyfy(
      store.getters["auth/getAuthenticated"].userAuthorization
    );
    return (
      userAuthorizations[api].includes(action) ||
      userAuthorizations[api].includes(ALL_API_ACTION)
    );
  };

  return {
    authorize,
    EVENTS_API_GROUP,
    USERS_API_GROUP,
    PARTICIPATIONS_API_GROUP,
    TEAMS_API_GROUP,
    ALL_API_ACTION,
    READ_API_ACTION,
    READ_SELF_API_ACTION,
    WRITE_API_ACTION,
    WRITE_SELF_API_ACTION,
    UPDATE_API_ACTION,
    UPDATE_SELF_API_ACTION,
    DELETE_API_ACTION,
  };
}
