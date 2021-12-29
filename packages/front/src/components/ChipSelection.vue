<template>
  <v-row class="chip-selection">
<!--    display selection-->
    <v-col cols="12" class="py-0">
      <v-row>
        <v-col v-for="(selection, i) in modelValue" :key="selection" class="shrink">
          <v-chip :closable="true" @click:close="deselectItem(i)">
            {{ selection[selectorKey] }}
          </v-chip>
        </v-col>
      </v-row>
    </v-col>
<!--    search fiedld-->
    <v-col cols="12">
      <v-text-field
        v-model="searchField"
        :label="label"
        :messages="messages"
        hide-details="true"
        @focusin="focusEnter"
        @focusout="focusExit"
      />
    </v-col>
<!--    display available results-->
    <v-col cols="12" class="py-0">
      <v-list v-if="focused" class="pa-0">
        <template v-for="item in selectable">
          <v-list-item v-if="filter(item)" :key="item.id" @click="selectItem(item)">
            <v-list-item-title>{{ item[selectorKey] }}</v-list-item-title>
          </v-list-item>
        </template>
      </v-list>
    </v-col>
  </v-row>
</template>

<script setup>
import { defineEmits, defineProps, ref } from "vue";

const searchField = ref("");
const focused = ref(false);
const messages = ref([]);

const props = defineProps({
  modelValue: {
    type: Array,
    required: false,
    default: () => [],
  },
  selectable: {
    type: Array,
    required: true,
  },
  selectorKey: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    default: "",
  }
});

const emit = defineEmits(['update:modelValue']);

const filter = (item) => {
  return !props.modelValue.includes(item) && (searchField.value.length > 0 ? item[props.selectorKey].includes(searchField.value) : true)
}

const selectItem = (item) => {
  const newValue = [...props.modelValue, item]
  emit('update:modelValue', newValue);
}

const deselectItem = (indice) => {
  const newValue = [...props.modelValue];
  newValue.splice(indice, 1);

  emit('update:modelValue', newValue);
}

const focusEnter = () => {
  focused.value = true
  messages.value = [`Liste des ${props.label}`]
}

const focusExit = () => {
  setTimeout(() => {
    focused.value = false
    messages.value = []
  }, 100)
}

</script>

<style>
.chip-selection .v-input__details {
  margin-bottom: 0 !important;
}
</style>