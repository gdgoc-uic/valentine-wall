<template>
  <!-- Notifications -->
  <client-only>
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
  </client-only>

  <!-- ID Modal -->
  <portal>
    <submit-id-modal v-if="!$store.state.user.associatedId && $store.getters.isLoggedIn" />
  </portal>

  <router-view v-slot="{ Component }">
    <suspense>
      <component :is="Component" />
    </suspense>
  </router-view>
</template>

<script lang="ts">
import { computed, defineComponent } from "@vue/runtime-core";
import { auth } from "./firebase";
import { HeadAttrs, useHead } from "@vueuse/head";

import BasicAlert from "./components/BasicAlert.vue";
import SubmitIDModal from "./components/SubmitIDModal.vue";
import { useRoute } from "vue-router";
import { getPageTitle } from "./router";
import ClientOnly from "./components/ClientOnly.vue";
import Portal from "./components/Portal.vue";

export default defineComponent({
  components: { 
    BasicAlert, 
    SubmitIdModal: SubmitIDModal, 
    ClientOnly,
    Portal 
  },
  setup() {
    const route = useRoute();
    useHead({
      htmlAttrs: {
        lang: 'en'
      },
      link: [
        { rel: 'icon', href: '/favicon.ico' }
      ],
      title: computed(() => getPageTitle(route)),
      meta: computed(() => {
        if (route.meta.metaTags && route.meta.metaTags instanceof Function) {
          return route.meta.metaTags(route);
        }
        return [
          { charset: 'UTF-8' },
          { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
          ...(route.meta.metaTags as HeadAttrs[] ?? [])
        ];
      }),
    });
  },
  mounted() {
    this.$store.dispatch('getGiftList');
    auth.onAuthStateChanged((user) => {
      this.$store.dispatch('onReceiveUser', user);
    });
  },
})
</script>

<style src="./assets/index.css"></style>

<style lang="postcss">
body {
  @apply bg-pink-200;
}
</style>