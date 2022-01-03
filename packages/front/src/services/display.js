import { computed } from "vue";
import { useDisplay } from "vuetify";

export default function useAppDisplay() {
  const display = useDisplay();
  const isMobile = computed(() => display.mobile.value);

  return {
    isMobile,
  };
}
