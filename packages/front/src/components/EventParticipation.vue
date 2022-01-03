<template>
  <v-btn-group :border="false" class="justify-center">
    <v-btn
      v-for="status in participationStatusList"
      :key="`${participation.id}-${status}`"
      :value="status"
      :color="btnColor(status)"
      :disabled="disabled"
      @click="handleUpdateParticipation(status)"
    >
      <v-icon :color="iconColor(status)">
        mdi-{{ participationStatusStyleMapping[status].icon }}
      </v-icon>
    </v-btn>
  </v-btn-group>
</template>

<script setup>
import { defineProps, defineEmits, computed } from "vue";
import useParticipations from "../services/participations";
import { useStore } from "vuex";

const store = useStore();
const { participationStatusStyleMapping, participationStatusList, updateParticipation } = useParticipations(store);

const props = defineProps({
  participation: {
    type: Object,
    required: true
  },
  disabled: {
    type: Boolean,
    required: false,
    default: false,
  }
});

const emit = defineEmits([
  'participation:create',
  'participation:update',
]);

const hasParticipation = computed(() => {
  return props.participation.id !== undefined;
})

const iconColor = (status) => {
  if (hasParticipation.value && props.participation.status === status) {
    return participationStatusStyleMapping[status].selectedIconColor;
  }
  return participationStatusStyleMapping[status].iconColor;
};

const btnColor = (status) => {
  if (hasParticipation.value && props.participation.status === status) {
    return participationStatusStyleMapping[status].selectedBgColor;
  }
  return undefined;
};

const handleUpdateParticipation = async (newStatus) => {
  if (!hasParticipation.value) {
    emit('participation:create', { status: newStatus })
    return;
  }
  console.log(newStatus)

  const updatedParticipation = {...props.participation};
  updatedParticipation.status = newStatus;

  const { status } = await updateParticipation(updatedParticipation);
  if (status) {
    emit('participation:update')
  }
};
</script>