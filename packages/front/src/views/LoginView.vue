<template>
  <v-row class="mt-5">
    <v-col md="4" offset-md="4" sm="8" offset-sm="2">
      <v-card class="elevation-12" title="Connexion">
        <v-card-text>
          <v-form>
            <v-text-field
              v-model="phoneNumber"
              prepend-inner-icon="mdi-phone"
              name="identity"
              label="N° Téléphone"
              type="text"
              :rules="rules"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="handleConnect">Se Connecter</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup>
import { ref } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import AuthenticationService from "../services/api/authentication";
import jwtDecode from "jwt-decode";

const store = useStore();
const router = useRouter();

const phoneNumber = ref(null);
const rules = ref([
  value => !!value || "Requis.",
  value => /^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$/im.test(value) || "Numéro de téléphone valide requis."
]);

const handleConnect = async () => {
  const client = store.getters["service/apiClient"];
  const authService = new AuthenticationService(client);

  const { status, data } = await authService.authenticate({ identity: phoneNumber.value })
  if (status) {
    store.dispatch("auth/authenticate", { user: jwtDecode(data) })
    store.dispatch("service/setAuthenticationToken", { token: data})

    router.push({ name: "calendar" });
  }
}
</script>
