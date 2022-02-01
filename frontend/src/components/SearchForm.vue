<template>
  <form @submit.prevent="searchMessageForm">
    <slot></slot>
  </form>
</template>

<script lang="ts">
import { defineComponent } from "@vue/runtime-core";

export default defineComponent({
  methods: {
    async searchMessageForm(e: SubmitEvent) {
      console.log(e.target);
      if (e.target && e.target instanceof HTMLFormElement) {
        const formData = new FormData(e.target);
        const recipientId = formData.get("recipient_id");
        if (
          !recipientId ||
          typeof recipientId !== "string" ||
          recipientId.length == 0
        ) {
          this.$notify({ type: "error", text: "Invalid search query input." });
          return;
        }

        // FIXME: vue won't navigate to wall page due to issues when unmounting the component
        if (!import.meta.env.SSR) {
          this.$router.push({
            name: "message-wall-page",
            params: { recipientId: recipientId },
          });
          e.target.reset();
        }
      }
    },
  },
});
</script>