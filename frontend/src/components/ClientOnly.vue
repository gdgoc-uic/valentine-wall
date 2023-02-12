<script lang="ts">
import {
  defineComponent,
  onMounted,
  ref,
  createCommentVNode,
  h
} from "@vue/runtime-core";

export default defineComponent({
  setup(_, { slots }) {
    const show = ref(false);
    onMounted(() => { show.value = true; });

    const defaultVnode = h(() => slots.default?.()) ?? createCommentVNode('Client only rendering with empty default slot')
    const placeholderVNode = h(() => slots.placeholder?.()) ?? createCommentVNode(`Client only rendering component placeholder`)

    return () => show.value ? defaultVnode : placeholderVNode
  }
});
</script>