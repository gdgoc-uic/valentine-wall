<script lang="ts">
import {
  createElementBlock,
  defineComponent,
  onMounted,
  ref,
} from "@vue/runtime-core";

export default defineComponent({
  name: "ClientOnly",
  props: ["fallback", "placeholder", "placeholderTag", "fallbackTag"],
  setup(_, { slots }) {
    const mounted = ref(false);
    if (!import.meta.env.SSR) {
      onMounted(() => {
        mounted.value = true;
      });
    }
    return (props: any) => {
      if (mounted.value) {
        return slots.default?.();
      }
      const slot = slots.fallback || slots.placeholder;
      if (slot) {
        return slot();
      }
      const fallbackStr = props.fallback || props.placeholder || "";
      const fallbackTag = props.fallbackTag || props.placeholderTag || "span";
      return createElementBlock(fallbackTag, null, fallbackStr);
    };
  },
});
</script>