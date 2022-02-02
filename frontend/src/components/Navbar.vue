<template>
  <header :class="[isHome ? 'pt-8' : 'py-8']" class="navbar mb-2 px-10 items-center bg-gradient-to-b from-[#FFEFEF] via-[#ffefef7e] to-transparent']">
    <template v-if="!isHome">
      <!-- TODO: mobile view -->
      <!-- <div class="flex-none flex lg:hidden">
        <button class="btn btn-square btn-ghost">
          <icon-menu class="h-6 w-6 inline-block" />
        </button>
      </div> -->
      <router-link :to="{ name: 'home-page' }" class="flex-none flex-nowrap mr-8 lg:flex space-x-2">
        <img src="../assets/images/icon.png" class="h-full w-20" alt="Icon" />
        <span class="flex-1 text-lg hidden md:hidden lg:block"> Valentine<span class="font-bold">Wall</span> </span>
      </router-link>
      <search-form class="flex-none w-1/3">
        <div class="form-control p-2 bg-white shadow-md rounded-xl w-full">
          <div class="flex space-x-2">
            <input
              type="text"
              placeholder="Search by ID"
              name="recipient_id"
              class="flex-1 input input-sm input-ghost"
              :value="$route.params.recipientId || ''"
            />
            <button
              class="
                btn btn-sm
                px-4
                py-2
                h-full
                border-0
                bg-rose-500
                hover:bg-rose-600
                rounded-xl
              "
            >
              <icon-search />
            </button>
          </div>
        </div>
      </search-form>
    </template>
    <div :class="[isHome ? 'mx-auto' : 'ml-auto mr-8']" class="flex-none flex space-x-8">
      <router-link :to="{ name: 'message-wall-page' }" class="text-gray-600 hover:text-gray-800">Recent</router-link>
      <router-link :to="{ name: 'rankings-page' }" class="text-gray-600 hover:text-gray-800">Rankings</router-link>
      <router-link :to="{ name: 'about-page' }" class="text-gray-600 hover:text-gray-800">About</router-link>
    </div>
    <client-only>
      <div class="flex-none">
        <div v-if="$store.getters.isLoggedIn" class="dropdown dropdown-end dropdown-hover -mb-2">
          <div :class="{'rounded-r-none': !shouldSendButtonHide}" tabindex="0" class="mb-2 shadow-md btn px-8 normal-case text-black bg-white border-0 hover:bg-gray-100">
            <span>{{ $store.state.user.email }}</span>
            <icon-dropdown />
          </div>
          <ul tabindex="0" class="p-2 shadow-md menu dropdown-content bg-base-100 rounded-box w-52">
            <li class="text-black"><router-link :to="{ name: 'settings-page' }">Settings</router-link></li>
            <li class="text-red-500">
              <a class="cursor-pointer" @click="logout">Logout</a>
            </li>
          </ul>
        </div>

        <button
          v-if="!shouldSendButtonHide && $store.getters.isLoggedIn"
          @click="$store.commit('SET_SEND_MESSAGE_MODAL_OPEN', true)"
          class="shadow-md btn border-none rounded-l-none bg-rose-500 hover:bg-rose-600 px-8 space-x-2">
          <icon-send />
          <span>Send a Message</span>
        </button>

        <login-button v-if="!isHome && !$store.getters.isLoggedIn" />
      </div>
    </client-only>
  </header>
</template>

<script lang="ts">
import IconMenu from '~icons/uil/align-center-alt';
import IconDropdown from '~icons/uil/angle-down';
import IconSend from '~icons/uil/message';
import IconSearch from '~icons/uil/search';
import { catchAndNotifyError } from '../notify';
import ClientOnly from './ClientOnly.vue';
import LoginButton from './LoginButton.vue';
import SearchForm from './SearchForm.vue';

export default {
  props: {
    isHome: {
      type: Boolean,
      default: false
    }
  },
  components: {
    IconMenu,
    IconDropdown,
    IconSend,
    IconSearch,
    SearchForm,
    ClientOnly,
    LoginButton
  },
  computed: {
    shouldSendButtonHide(): boolean {
      return this.isHome || this.$route.path.startsWith('/settings') || this.$route.path.startsWith('/about');
    }
  },
  methods: {
    async logout() {
      try {
        await this.$store.dispatch('logout');
        this.$router.replace('/');
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    }
  }
};
</script>

<style>
</style>