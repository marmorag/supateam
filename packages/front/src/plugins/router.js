import { createRouter, createWebHashHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import CalendarView from "../views/CalendarView.vue";
import EventView from "../views/EventView.vue";
import CreateEventView from "../views/CreateEventView.vue";
import AdminUsersView from "../views/admin/UsersView.vue";
import store from "./store";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      name: "login",
      path: "/login",
      component: LoginView,
      meta: {
        title: "Connexion",
        requireAuthenticated: false,
        requireAdmin: false,
      },
      beforeEnter: (to, from, next) => {
        if (store.getters["auth/isAuthenticated"]) {
          next({ name: "calendar" });
        }
        next();
      },
    },
    {
      name: "calendar",
      path: "/calendar",
      component: CalendarView,
      meta: {
        title: "Calendrier",
        requireAuthenticated: true,
        requireAdmin: false,
      },
    },
    {
      name: "event",
      path: "/event/:id",
      component: EventView,
      props: true,
      meta: {
        title: "Evenement",
        requireAuthenticated: true,
        requireAdmin: false,
      },
    },
    {
      name: "create-event",
      path: "/create-event",
      component: CreateEventView,
      meta: {
        title: "Créer un évènement",
        requireAuthenticated: true,
        requireAdmin: false,
      },
    },
    {
      name: "admin-users",
      path: "/admin/users",
      component: AdminUsersView,
      meta: {
        title: "Administration - Utilisateurs",
        requireAuthenticated: true,
        requireAdmin: true,
      },
    },
  ],
});

router.beforeEach((to, from, next) => {
  if (to.meta.requireAuthenticated && !store.getters["auth/isAuthenticated"]) {
    next({ name: "login" });
  }

  document.title = `SupaTeam - ${to.meta.title}`;
  next();
});

export default router;
