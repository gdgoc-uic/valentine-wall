<template>
  <div class="background"></div>

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
    <submit-message-modal
      :key="$route.fullPath" 
      :open="$store.state.isSendMessageModalOpen" 
      @update:open="$store.commit('SET_SEND_MESSAGE_MODAL_OPEN', $event)" />
  </portal>

  <div v-if="$store.state.isAuthLoading" class="h-screen flex items-center justify-center">
    <loading />
  </div>

  <template v-else>
    <navbar :is-home="$route.name === 'home-page'" class="sticky top-0 z-50" />
    <router-view v-slot="{ Component }">
      <suspense>
        <component :is="Component" class="relative" />
      </suspense>
    </router-view>

    <app-footer />
  </template>
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
import Navbar from "./components/Navbar.vue";
import SubmitMessageModal from './components/SendMessageModal.vue';
import Footer from "./components/Footer.vue";
import { catchAndNotifyError } from "./notify";
import Loading from "./components/Loading.vue";

export default defineComponent({
  components: { 
    BasicAlert, 
    SubmitMessageModal,
    SubmitIdModal: SubmitIDModal, 
    ClientOnly,
    Portal,
    Navbar,
    AppFooter: Footer,
    Loading
  },
  setup() {
    const route = useRoute();
    useHead({
      htmlAttrs: {
        lang: 'en'
      },
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700&display=swap' }
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
    if (!import.meta.env.SSR) {
      auth.onAuthStateChanged((user) => {
        this.$store.dispatch('onReceiveUser', user);
      });

      Promise.all([
        this.$store.dispatch('getGiftList'),
        this.$store.dispatch('getDepartmentList')
      ]).catch((err) => {
        catchAndNotifyError(this, err);
      });
    }
  },
})
</script>

<style src="./assets/index.css"></style>

<style lang="postcss">
/* TODO: Customize Daisy CSS based on color palette. Remove "patched" stylings to inline classes as much as possible. */

html {
  height: 100%;
}

body {
  min-height: 100%;
}

.background {
  background-image: url(./assets/images/background.png);
  background-size: 250% 100%;
  background-repeat: no-repeat;
  background-position: top center;
  height: 110vh;
  widows: 100vw;
  z-index: -1;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

@media screen(md) {
  .background {
    background-size: 150% 100%;
  }
}

@media screen(lg) {
  .background {
    background-size: 100% 100%;
  }
}
</style>