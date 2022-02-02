<script lang="ts">
import { defineComponent, h, TeleportProps } from "@vue/runtime-core";
import { Teleport } from 'vue';
import ClientOnly from "./ClientOnly.vue";

export default defineComponent({
    render() {
        return h(ClientOnly, [
            // @ts-ignore
            h(Teleport, { to: 'body' } as TeleportProps, this.$slots)
        ]);
    }
})
</script>