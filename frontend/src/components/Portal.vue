<script lang="ts">
import { defineComponent, h, TeleportProps } from "@vue/runtime-core";
import { Teleport } from 'vue';
import ClientOnly from "./ClientOnly.vue";

export default defineComponent({
    setup(_, { slots }) {
        return () => h(ClientOnly, () => [
            h(Teleport, { to: 'body' } as TeleportProps, slots)
        ]);
    }
})
</script>