<template>
  <div>
    <modal v-model:open="isCoinsModalOpen" with-closing-button title="Your coins">
    <div class="flex flex-col space-y-4">
      <div class="bg-gray-200 text-3xl flex items-center justify-center p-3 rounded-md">
        <icon-coin class="mr-2" />
        <span>₱{{ authState.user?.expand.wallet?.balance.toFixed(2) ?? 'unknown' }}</span>
      </div>

      <div>
        <h3 class="font-bold mb-2">How to earn coins?</h3>
        <div class="flex flex-col space-y-2">
          <div :key="'how_to_earn_' + hi" 
            v-for="(ht, hi) in howToEarn" 
            class="flex flex-row bg-gray-200 p-3 rounded-md">
            <div class="w-3/4">
              <p class="font-bold text-rose-500">{{ ht.description }}</p>
            </div>

            <div class="w-1/4 text-xl flex items-center">
              <icon-coin class="mr-2" />
              <span>₱{{ ht.amount }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-end">
        <router-link class="btn btn-primary" :to="{ name: 'settings-transactions-section' }">
          Transactions
        </router-link>
      </div>
    </div>
  </modal>

  <div>
    <div class="hidden lg:block bg-[#EF9B95] px-10 py-1">
      <div :class="[isHome ? 'mx-auto' : 'ml-auto mr-8']" class="flex flex-none space-x-8">
        <router-link
          :key="'links_' + i"
          v-for="(link, i) in navbarLinks"
          :to="link.to"
          class="text-white hover:text-gray-200">
          {{ link.label }}
        </router-link>
      </div>
    </div>
    <header ref="mainNavbar" :class="[isHome ? 'pt-2 lg:pt-8' : 'py-2 lg:py-8']" class="navbar mb-2 px-2 lg:px-10 items-center bg-gradient-to-b from-[#FFEFEF] via-[#ffefef7e] to-transparent">
      <!-- TODO: mobile view -->
      <div class="flex-none flex lg:hidden">
        <button @click="menuOpen = !menuOpen" class="btn btn-square btn-ghost btn-sm">
          <icon-menu class="h-5 w-5 inline-block" />
        </button>
      </div>
      <router-link :class="[isHome ? 'md:hidden' : 'lg:flex']" :to="{ name: 'home-page' }" class="flex-none flex items-center flex-nowrap lg:mr-8 space-x-2">
        <img src="../assets/images/icon.png" class="w-8 h-8 sm:w-14 sm:h-14 md:w-20 md:h-20" alt="Icon" />
        <span class="flex-1 text-lg hidden md:hidden lg:block"> <span class="font-bold">Valentine</span>Wall </span>
      </router-link>
      <search-form :class="{'md:hidden': isHome}" class="flex-1 lg:flex-none lg:w-1/3 min-w-0 ml-1 sm:ml-2">
        <div class="form-control p-1 lg:p-2 bg-white shadow-md rounded-xl w-full">
          <div class="flex space-x-2 items-center">
            <input
              type="text"
              placeholder="Search by ID"
              name="recipient_id"
              class="flex-1 input input-sm input-ghost"
              :value="$route.params.recipientId || ''"
            />
            <button
              class="btn btn-sm px-4 py-3 lg:py-2 h-full border-0 bg-rose-500 hover:bg-rose-600 rounded-xl">
              <icon-search />
            </button>
          </div>
        </div>
      </search-form>
      <client-only>
        <div class="flex-none hidden md:ml-auto md:flex">
          <div v-if="!isReadOnly() && authState.isLoggedIn" class="dropdown dropdown-end dropdown-hover -mb-2">
            <div tabindex="0" class="rounded-r-none mb-2 shadow-md btn normal-case text-black bg-white border-0 hover:bg-gray-100">
              <span class="overflow-hidden text-ellipsis">
                {{ authState.user.username }}
              </span>
              <icon-dropdown />
            </div>
            <ul tabindex="0" class="p-2 shadow-md menu dropdown-content bg-base-100 rounded-box w-52">
              <li :class="al.liClass"
                  :key="'account_link_' + i"
                  v-for="(al, i) in accountLinks">
                <component :is="al.tag" v-bind="al.props">
                  {{ al.children }}
                </component>
              </li>
            </ul>
          </div>
          
          <button
            v-if="!isReadOnly() && authState.isLoggedIn" 
            @click="isCoinsModalOpen = true"
            :class="[!shouldSendButtonHide ? 'rounded-none' : 'rounded-l-none']" 
            class="btn shadow-md normal-case text-black bg-white border-0 hover:bg-gray-100">
            <icon-coin class="mr-2" />
            <span>₱{{ authState.user!.expand.wallet?.balance.toFixed(2) ?? 'unknown' }}</span>
          </button>
          <button
            v-if="!isReadOnly() && !shouldSendButtonHide && authState.isLoggedIn"
            @click="store.state.isSendMessageModalOpen = true"
            class="shadow-md btn border-none rounded-l-none bg-rose-500 hover:bg-rose-600 lg:px-8 space-x-2">
            <icon-send />
            <span class="hidden lg:block">Send a Message</span>
            <span class="lg:hidden">Send</span>
          </button>
          <login-button v-if="!isHome && !authState.isLoggedIn" />
        </div>
      </client-only>
    </header>

    <!-- Mobile Sidebar Overlay -->
    <Transition name="sidebar-overlay">
      <div 
        v-if="menuOpen"
        @click.self="menuOpen = false" 
        style="z-index: 99999;"
        class="lg:hidden fixed inset-0 bg-black/30 backdrop-blur-sm">
        <!-- Sidebar Panel -->
        <Transition name="sidebar-slide">
          <div v-if="menuOpen" class="bg-white h-full w-[75vw] max-w-[320px] flex flex-col shadow-2xl">
            <!-- Header -->
            <div class="flex items-center justify-between px-5 py-4 bg-gradient-to-r from-rose-500 to-rose-400">
              <router-link :to="{ name: 'home-page' }" @click="menuOpen = false" class="flex items-center space-x-2">
                <img src="../assets/images/icon.png" class="w-8 h-8" alt="Icon" />
                <span class="text-white font-bold text-lg">Valentine Wall</span>
              </router-link>
              <button @click="menuOpen = false" class="text-white/80 hover:text-white p-1">
                <icon-close class="w-5 h-5" />
              </button>
            </div>

            <client-only>
              <!-- User Card -->
              <div v-if="!isReadOnly() && authState.isLoggedIn" class="px-5 py-4 bg-rose-50 border-b border-rose-100">
                <div class="flex items-center space-x-3">
                  <div class="w-10 h-10 rounded-full bg-rose-400 flex items-center justify-center text-white font-bold text-lg">
                    {{ authState.user.username?.charAt(0)?.toUpperCase() || '?' }}
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="font-semibold text-gray-800 truncate">{{ authState.user.username }}</p>
                    <button @click="isCoinsModalOpen = true; menuOpen = false" class="flex items-center space-x-1 text-sm text-gray-500 hover:text-rose-500">
                      <icon-coin class="w-4 h-4" />
                      <span>₱{{ authState.user!.expand.wallet?.balance.toFixed(2) ?? '0.00' }}</span>
                    </button>
                  </div>
                </div>
              </div>

              <!-- Nav Links -->
              <nav class="flex-1 overflow-y-auto py-2">
                <div class="px-3">
                  <router-link
                    :key="'links_' + i"
                    v-for="(link, i) in navbarLinks"
                    :to="link.to"
                    @click="menuOpen = false"
                    class="flex items-center px-3 py-3 rounded-lg text-gray-700 hover:bg-rose-50 hover:text-rose-600 transition-colors font-medium">
                    {{ link.label }}
                  </router-link>
                </div>

                <!-- Divider -->
                <div v-if="!isReadOnly() && authState.isLoggedIn" class="my-2 border-t border-gray-100"></div>

                <!-- Account Links -->
                <div v-if="!isReadOnly() && authState.isLoggedIn" class="px-3">
                  <div
                    :key="'account_link_' + i"
                    v-for="(al, i) in accountLinks"
                    @click="menuOpen = false">
                    <component 
                      :is="al.tag" 
                      v-bind="al.props"
                      class="flex items-center px-3 py-3 rounded-lg hover:bg-rose-50 transition-colors font-medium"
                      :class="al.liClass">
                      {{ al.children }}
                    </component>
                  </div>
                </div>
              </nav>

              <!-- Bottom Action -->
              <div v-if="!isReadOnly() && authState.isLoggedIn" class="p-4 border-t border-gray-100">
                <button
                  v-if="!shouldSendButtonHide"
                  @click="store.state.isSendMessageModalOpen = true; menuOpen = false"
                  class="btn border-none w-full bg-rose-500 hover:bg-rose-600 text-white rounded-xl space-x-2">
                  <icon-send />
                  <span>Send a Message</span>
                </button>
              </div>

              <div v-if="!authState.isLoggedIn" class="p-4 mt-auto border-t border-gray-100">
                <login-button @click="menuOpen = false" class="w-full" />
              </div>
            </client-only>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
  </div>
</template>

<script lang="ts" setup>
import IconMenu from '~icons/uil/align-center-alt';
import IconDropdown from '~icons/uil/angle-down';
import IconSend from '~icons/uil/message';
import IconSearch from '~icons/uil/search';
import IconClose from '~icons/uil/multiply';
import IconCoin from '~icons/twemoji/coin';
import { catchAndNotifyError } from '../notify';
import ClientOnly from './ClientOnly.vue';
import LoginButton from './LoginButton.vue';
import SearchForm from './SearchForm.vue';
import Modal from './Modal.vue';
import { useRoute, useRouter } from 'vue-router';
import { ref, computed } from 'vue';
import { useAuth, useStore } from '../store_new';
import { isReadOnly } from '../utils';

const props = defineProps({
  isHome: {
    type: Boolean,
    default: false
  }
});

const router = useRouter();
const route = useRoute();
const { state: authState, methods: {logout} } = useAuth();
const store = useStore();
const isCoinsModalOpen = ref(false);
const howToEarn = computed(() => [
  {
    description: 'Earn from idle time (per second)',
    amount: 0.05
  },
  {
    description: 'Share posts',
    amount: 300,
  },
  {
    description: 'Receive money virtual gift',
    amount: store.state.giftList.find(g => g.uid === 'money')?.price ?? 1000,
  },
  {
    description: 'Ask the admins?',
    amount: '???'
  }
]);

const menuOpen = ref(false);
const navbarLinks = computed(() => [
  ...(!isReadOnly() && authState.isLoggedIn ? [{
    label: 'Your Wall',
    to: { name: 'message-wall-page', params: { recipientId: authState.user.expand.details.student_id } }
  }] : []),
  {
    label: 'Recent',
    to: { name: 'recent-wall-page' }
  },
  {
    label: 'Rankings',
    to: { name: 'rankings-page' }
  }
]);

const accountLinks = [
  {
    liClass: 'text-black',
    tag: 'router-link',
    props: {
      to: { name: 'settings-page' }
    },
    children: 'Settings'
  },
  {
    liClass: 'text-red-500',
    tag: 'a',
    props: {
      class: 'cursor-pointer',
      onClick: function(e: Event) {
        e.preventDefault();
        try {
          router.replace({ name: 'home-page' });
          logout();
        } catch(e) {
          catchAndNotifyError(e);
        }
      }
    },
    children: 'Logout' 
  }
];

const shouldSendButtonHide = computed(() => {
  return props.isHome || 
    route.path.startsWith('/settings') || 
    route.path.startsWith('/about');
});
</script>

<style scoped>
.sidebar-overlay-enter-active,
.sidebar-overlay-leave-active {
  transition: opacity 0.25s ease;
}
.sidebar-overlay-enter-from,
.sidebar-overlay-leave-to {
  opacity: 0;
}
.sidebar-slide-enter-active {
  transition: transform 0.25s ease-out;
}
.sidebar-slide-leave-active {
  transition: transform 0.2s ease-in;
}
.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
  transform: translateX(-100%);
}
</style>