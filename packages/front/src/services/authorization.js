export default function useAuthorization(store) {
  const EVENTS_API_GROUP = "events";
  const USERS_API_GROUP = "users";

  const validApiGroup = [EVENTS_API_GROUP, USERS_API_GROUP];

  const ALL_API_ACTION = "*";
  const READ_API_ACTION = "*";
  const WRITE_API_ACTION = "*";

  const validApiAction = [ALL_API_ACTION, READ_API_ACTION, WRITE_API_ACTION];

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
      userAuthorizations[api].includes("*")
    );
  };

  return {
    authorize,
    EVENTS_API_GROUP,
    USERS_API_GROUP,
    ALL_API_ACTION,
    READ_API_ACTION,
    WRITE_API_ACTION,
  };
}
