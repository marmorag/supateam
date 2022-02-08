<template>
  <ThePageTitle offset="2" />
  <UserList :users="users" />
</template>

<script setup>
import ThePageTitle from "../../components/ThePageTitle.vue";
import UserList from "../../components/users/UserList.vue";
import useUsers from "../../services/users";
import { useStore } from "vuex";

const store = useStore()
const { users } = useUsers(store)
</script>
