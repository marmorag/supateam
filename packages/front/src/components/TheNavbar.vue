<template>
  <v-app-bar app>
    <v-app-bar-title>Team Superflu</v-app-bar-title>
    <v-spacer />
    <template v-if="isAuthenticated && !isMobile">
      <v-btn v-if="isAdmin" prepend-icon="mdi-chess-king">
        Administration
        <v-menu activator="parent" anchor="bottom end">
          <v-sheet>
            <v-list>
              <v-list-item v-for="item in adminLink" :key="item.to" :to="{ name: item.to }">
                <v-list-item-title>{{ item.title }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-sheet>
        </v-menu>
      </v-btn>

      <v-btn v-if="canCreateEvent" class="mr-5" prepend-icon="mdi-calendar-plus" :to="{ name: 'create-event' }">évènement</v-btn>
      <v-btn class="mr-5" prepend-icon="mdi-calendar" :to="{ name: 'calendar' }">Calendrier</v-btn>
      <v-divider inset vertical></v-divider>
      <v-btn class="mx-5" prepend-icon="mdi-logout" @click="handleDisconnect">Déconnexion</v-btn>
    </template>
    <template v-else-if="isAuthenticated && isMobile">
      <v-btn icon="mdi-menu" @click="toggleDrawer"></v-btn>
    </template>
    <template v-else>
      <v-btn class="mx-5" prepend-icon="mdi-login" :to="{ name: 'login' }">Connexion</v-btn>
    </template>
  </v-app-bar>

  <v-navigation-drawer v-model="drawer" absolute :temporary="true">
    <v-list :nav="true" density="compact">
        <v-list-item v-if="canCreateEvent" :to="{ name: 'create-event' }" @click="toggleDrawer">
          <v-list-item-avatar>
            <v-icon>mdi-calendar-plus</v-icon>
          </v-list-item-avatar>
          <v-list-item-title>Évènements</v-list-item-title>
        </v-list-item>

        <v-list-item :to="{ name: 'calendar' }" @click="toggleDrawer">
          <v-list-item-avatar>
            <v-icon>mdi-calendar</v-icon>
          </v-list-item-avatar>
          <v-list-item-title>Calendrier</v-list-item-title>
        </v-list-item>

        <v-list-item  @click="handleDisconnect">
          <v-list-item-avatar>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-avatar>
          <v-list-item-title>Déconnexion</v-list-item-title>
        </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import { computed, ref } from "vue";
import { useStore } from "vuex";
import useAuthorization from "../services/authorization";
import { useRouter } from "vue-router";
import useAppDisplay from "../services/display";

const router = useRouter();
const store = useStore();
const { authorize, USERS_API_GROUP, EVENTS_API_GROUP, WRITE_API_ACTION, ALL_API_ACTION } = useAuthorization(store);
const { isMobile } = useAppDisplay();

const adminLink = [
  {
    title: "Utilisateurs",
    to: "admin-users",
    icon: "account",
  },
  // {
  //   title: "Évènements",
  //   to: "admin-events"
  // },
];

const drawer = ref(false);

const isAuthenticated = computed(() => {
  return store.getters["auth/isAuthenticated"];
})

const canCreateEvent = computed(() => authorize({ api: EVENTS_API_GROUP, action: WRITE_API_ACTION }) || false);
const isAdmin = computed(() => authorize({ api: USERS_API_GROUP, action: ALL_API_ACTION }) || false);

const toggleDrawer = () => {
  drawer.value = !drawer.value
}

const handleDisconnect = async () => {
  await store.dispatch("auth/disconnect");
  await store.dispatch("service/clearAuthenticationToken");
  await router.push({ name: "login" })
  toggleDrawer()
};
</script>