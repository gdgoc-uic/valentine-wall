<template>
  <div>
    <div class="hidden lg:block bg-[#f3b3ae] px-10 py-1">
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
    <header ref="mainNavbar" :class="[isHome ? 'pt-2 lg:pt-8' : 'py-2 lg:py-8']" class="navbar mb-2 px-2 lg:px-10 items-center bg-gradient-to-b from-[#FFEFEF] via-[#ffefef7e] to-transparent']">
      <!-- TODO: mobile view -->
      <div class="flex-none flex lg:hidden">
        <button @click="menuOpen = !menuOpen" class="btn btn-square btn-ghost">
          <icon-menu class="h-6 w-6 inline-block" />
        </button>
      </div>
      <router-link :class="[isHome ? 'md:hidden' : 'lg:flex']" :to="{ name: 'home-page' }" class="flex-none flex-nowrap mr-2 lg:mr-8 space-x-2">
        <img src="../assets/images/icon.png" class="h-full w-14 md:w-20" alt="Icon" />
        <span class="flex-1 text-lg hidden md:hidden lg:block"> <span class="font-bold">Valentine</span>Wall </span>
      </router-link>
      <search-form :class="{'md:hidden': isHome}" class="flex-1 lg:flex-none md:mr-2 lg:mr-0 lg:w-1/3">
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
          <div v-if="authState.isLoggedIn" class="dropdown dropdown-end dropdown-hover -mb-2">
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
          <!-- TODO: add 'add coins' modal -->
          <button 
            :class="[!shouldSendButtonHide ? 'rounded-none' : 'rounded-l-none']" 
            v-if="authState.isLoggedIn" 
            class="btn shadow-md normal-case text-black bg-white border-0 hover:bg-gray-100">
            <icon-coin class="mr-2" />
            <span>ღ{{ authState.user!.expand.wallet?.balance ?? 'unknown' }}</span>
          </button>
          <button
            v-if="!shouldSendButtonHide && authState.isLoggedIn"
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

    <div 
      @click.self="menuOpen = false" 
      :class="[menuOpen ? 'block' : 'hidden']" 
      class="lg:hidden bg-[#FFEFEF] bg-opacity-50 h-screen fixed inset-x-0 bottom-0">
      <div class="bg-[#FFEFEF] flex h-full lg:hidden p-8 flex-col w-[85vw] drop-shadow-xl">
        <button
          @click="menuOpen = false"
          style="right: 20px; top: 20px;"
          class="bg-rose-500 hover:bg-rose-600 transition-colors text-white p-2 rounded-full absolute">
          <icon-close />
        </button>
        <client-only>
          <div class="py-8 flex flex-col space-y-8">
            <router-link
              :key="'links_' + i"
              v-for="(link, i) in navbarLinks"
              :to="link.to"
              @click="menuOpen = false"
              class="text-3xl text-gray-600 hover:text-gray-800">
              {{ link.label }}
            </router-link>
          </div>
          <div v-if="authState.isLoggedIn" class="bg-white bg-opacity-60 p-4 rounded-xl mt-auto">
            <p>Signing in as</p>
            <h3 class="text-2xl text-ellipsis overflow-hidden font-bold">
              {{ authState.user.username }}
            </h3>
            <ul class="space-y-4 py-4">
              <li class="text-xl"
                  @click="menuOpen = false"
                  :class="al.liClass"
                  :key="'account_link_' + i"
                  v-for="(al, i) in accountLinks">
                <component :is="al.tag" v-bind="al.props">
                  {{ al.children }}
                </component>
              </li>
            </ul>
            <button 
              v-if="authState.isLoggedIn" 
              class="btn shadow-md normal-case text-black bg-white border-0 hover:bg-gray-100 w-full mb-2">
              <icon-coin class="mr-2" />
              <span>ღ{{ authState.user!.expand.wallet?.balance ?? 'unknown' }}</span>
            </button>
            <button
              v-if="!shouldSendButtonHide"
              @click="store.state.isSendMessageModalOpen = true; menuOpen = false"
              class="shadow-md btn border-none w-full bg-rose-500 hover:bg-rose-600 px-8 space-x-2">
              <icon-send />
              <span>Send a Message</span>
            </button>
          </div>
          <login-button v-else @click="menuOpen = false" class="mt-auto" />
        </client-only>
      </div>
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
import { useRoute, useRouter } from 'vue-router';
import { ref, computed } from 'vue';
import { pb } from '../client';
import { useAuth, useStore } from '../store_new';

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

const menuOpen = ref(false);
const navbarLinks = [
  {
    label: 'Rankings',
    to: { name: 'rankings-page' }
  },
  {
    label: 'About',
    to: { name: 'about-page' }
  }
];

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
          // TODO: bring back later
          // catchAndNotifyError(this, e);
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