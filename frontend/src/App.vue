<template>
  <!-- Notifications -->
  <notification-group>
    <div style="z-index: 9999;" class="fixed inset-0 flex flex-col items-center px-4 py-6 pointer-events-none">
      <div class="w-full max-w-2xl">
        <notification 
          v-slot="{ notifications }"
          enter="transform ease-out duration-300 transition"
          enter-from="translate-y-3 opacity-0 sm:translate-y-0 sm:translate-x-4"
          enter-to="translate-y-0 opacity-100 sm:translate-x-0"
          leave="transition ease-in duration-500"
          leave-from="opacity-100"
          leave-to="opacity-0"
          move="transition duration-500"
          move-delay="delay-300">
          <basic-alert 
            class="shadow my-3"
            v-for="notification in notifications"
            :key="notification.id"
            :type="notification.type"
            :message="notification.text" />
        </notification>
      </div>
    </div>
  </notification-group>

  <!-- ID Modal -->
  <teleport to="body">
    <submit-id-modal v-if="$store.getters.isLoggedIn" />
  </teleport>

  <router-view></router-view>
</template>

<script lang="ts">
import { defineComponent } from "@vue/runtime-core";
import { auth } from "./firebase";

import BasicAlert from "./components/BasicAlert.vue";
import SubmitIDModal from "./components/SubmitIDModal.vue";

export default defineComponent({
  components: { BasicAlert, SubmitIdModal: SubmitIDModal },
  mounted() {
    this.$store.dispatch('getGiftList');
    auth.onAuthStateChanged((user) => {
      this.$store.dispatch('onReceiveUser', user);
    });
  },
})
</script>