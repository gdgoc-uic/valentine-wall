<template>
  <div class="background"></div>

  <div style="z-index: 999" class="fixed inset-x-0 bottom-0 py-4 bg-rose-200 text-center hidden md:block">
    <div class="max-w-7xl mx-auto w-full flex flex-col md:flex-row space-y-2 md:space-y-0 md:space-x-4 items-center justify-center">
      <p class="text-lg"><b>Let's not make your Valentine's ruined by bugs.</b> </p>
      
      <feedback-form v-slot="{ openDialog }">
        <div class="tooltip" data-tip="Problems, suggestions? Post it here!">
          <button @click="openDialog" class="btn btn-sm btn-primary space-x-2">
            <icon-comment-add />
            <span>Add your Feedback</span>
          </button>                                                                                                                                                                                   
        </div>
      </feedback-form>
    </div>
  </div>

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
    <setup-dialog v-if="!user?.expand.details?.student_id && isLoggedIn" />
    <submit-message-modal
      :key="$route.fullPath" 
      :open="store.state.isSendMessageModalOpen" 
      @update:open="store.state.isSendMessageModalOpen = $event" />
  </portal>

  <div v-if="isAuthLoading" class="h-screen flex items-center justify-center">
    <loading />
  </div>

  <template v-else>
    <portal>
      <welcome-modal />
    </portal>

    <navbar v-if="!$route.meta.disableAppHeader" :is-home="$route.name === 'home-page'" class="sticky top-0 z-50" />
    <main class="min-h-[80vh]">
      <router-view v-slot="{ Component }">
        <suspense>
          <component :is="Component" class="relative" />
        </suspense>
      </router-view>
    </main>

    <app-footer />
  </template>
</template>

<script lang="ts" setup>
import { computed } from "@vue/runtime-core";
import { HeadAttrs, useHead } from "@vueuse/head";

import BasicAlert from "./components/BasicAlert.vue";
import SetupDialog from "./components/SetupDialog.vue";
import { useRoute } from "vue-router";
import { getPageTitle } from "./router";
import ClientOnly from "./components/ClientOnly.vue";
import Portal from "./components/Portal.vue";
import Navbar from "./components/Navbar.vue";
import SubmitMessageModal from './components/SendMessageModal.vue';
import AppFooter from "./components/Footer.vue";
import Loading from "./components/Loading.vue";
import WelcomeModal from "./components/WelcomeModal.vue";
import IconCommentAdd from "~icons/uil/comment-add";
import FeedbackForm from "./components/FeedbackForm.vue";
import { useAuth, useStore } from "./store_new";
import { onMounted, onUnmounted, toRefs } from "vue";
import { pb } from "./client";
import { User } from "./types";
import fallbackImage from './assets/images/fallback.png';
import { useQueryClient } from "@tanstack/vue-query";

const queryClient = useQueryClient();
const route = useRoute();
const { state, methods: { onReceiveUser } } = useAuth();
const { isLoggedIn, isAuthLoading, user } = toRefs(state);
const store = useStore();

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
    let metaTags: HeadAttrs[] = [];
    if (route.meta.metaTags && route.meta.metaTags instanceof Function) {
      metaTags = route.meta.metaTags(route);
    } else if (Array.isArray(route.meta.metaTags)) {
      metaTags = route.meta.metaTags;
    }

    if (metaTags.findIndex(m => m.name === 'og:image') === -1) {
      metaTags.push({ name: 'og:image', content: fallbackImage });
    }

    return [
      { charset: 'UTF-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
      { name: 'og:image:width', content: '1200' },
      { name: 'og:image:height', content: '630'},
      ...metaTags
    ];
  }),
});

if (!import.meta.env.SSR) {
  onMounted(() => {
    store.methods.loadGiftsAndDepartments();

    pb.authStore.onChange((_, user) => {
      store.methods.checkFirstTimeVisitor();

      onReceiveUser(user as User | null, store.state)
        .then(() => {
          return queryClient.refetchQueries();
        });
    }, true);
  });

  onUnmounted(() => {
    pb.collection('virtual_wallets').unsubscribe();
  });
}
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

.badge-primary, .tabs-boxed .tab-active {
  @apply bg-rose-500 border-rose-500 text-white hover:text-white;
}

.btn.btn-primary {
  @apply bg-rose-500 border-rose-500 hover:bg-rose-600 hover:border-rose-600 text-white hover:text-white;
}
</style>