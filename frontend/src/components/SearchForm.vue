<template>
  <form @submit.prevent="searchMessageForm">
    <slot></slot>
  </form>
</template>

<script lang="ts">
import { defineComponent } from "@vue/runtime-core";
import { notify } from "../notify";

const idRegex = /^[0-9]{6,12}$/;

export default defineComponent({
  methods: {
    async searchMessageForm(e: SubmitEvent) {
      if (e.target && e.target instanceof HTMLFormElement) {
        const formData = new FormData(e.target);
        const recipientId = formData.get("recipient_id");
        if (
          !recipientId ||
          typeof recipientId !== "string" ||
          recipientId.length == 0
        ) {
          notify({ type: "error", text: "Invalid search query input." });
          return;
        }

        if (!idRegex.test(recipientId)) {
          notify({ type: "error", text: "ID must be numeric and should be 6-12 characters long." });
          return;
        }

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