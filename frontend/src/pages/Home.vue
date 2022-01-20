<template>
  <main class="bg-pink-200 h-screen w-screen flex">
    <div style="min-height: 70vh" class="relative bg-white shadow-2xl shadow-pink-400 max-w-5xl w-full mx-auto lg:mt-20 self-start p-12 lg:p-14 rounded-lg flex flex-col">
      <div v-if="$store.state.isAuthLoading" class="rounded-lg z-30 absolute inset-0 bg-white bg-opacity-75 flex flex-col justify-center items-center h-full">
        <p>Loading...</p>
      </div>
      
      <div class="flex-1 flex flex-col h-full">
        <div class="flex flex-col justify-center text-center items-center mb-8">
          <h1 class="text-5xl lg:text-7xl mb-4 lg:mb-8 font-extrabold bg-gradient-to-tr from-pink-400 to-rose-500 bg-clip-text text-transparent">Valentine Wall</h1>
          <p class="lg:w-5/6 text-2xl text-gray-500">Lorem ipsum dolor sit amet consectetur adipisicing elit. Tempore, dignissimos? Accusantium, mollitia.</p>
        </div>

        <span class="text-center uppercase tracking-widest text-gray-500 mt-auto mb-6 font-semibold">I want to...</span>
        <div class="flex flex-col lg:flex-row divide-y-2 lg:divide-y-0 lg:divide-x-2">
          <div class="w-full lg:w-1/2 flex flex-col justify-between py-8 lg:pr-8">
            <div class="text-center mb-8">
              <h3 class="text-rose-500 text-2xl font-bold">Search my messages</h3>
              <p class="text-gray-500">Lorem ipsum dolor sit amet consectetur adipisicing elit.</p>
            </div>

            <form @submit.prevent="searchMessageForm" class="flex flex-col justify-center items-center">
              <input class="input input-bordered w-full mb-2" type="text" name="student_id" pattern="[0-9]{12}" placeholder="12-digit student ID">
              <button class="btn w-full border-none rounded-full bg-rose-100 text-rose-500 hover:bg-rose-200">Search</button>
            </form>
          </div>
  
          <div class="w-full lg:w-1/2 flex flex-col justify-between items-center py-8 lg:pl-8">
            <div class="text-center mb-8">
              <h3 class="text-rose-500 text-2xl font-bold">Post a message</h3>
              <p class="text-gray-500">Lorem ipsum dolor sit amet consectetur adipisicing elit.</p>
            </div>
            <div v-if="!$store.getters.isLoggedIn" class="lg:w-2/3 flex flex-col items-center">
              <button @click="login" class="btn btn-lg w-full border-none rounded-full bg-rose-500 hover:bg-rose-700">Login</button>
              <span class="mt-4 text-gray-500 text-sm">Using UIC Google Account</span>
            </div>
            <div v-else class="flex flex-col items-center w-full mt-auto">
              <button @click="isFormOpen = true" class="btn btn-lg w-full border-none rounded-full bg-rose-500 hover:bg-rose-700">Write Message</button>
              <span class="mt-4 text-gray-500 text-sm">Maximum of 240 characters. 10 minutes per post</span>
              <teleport to="body">
                <submit-message-modal v-model:open="isFormOpen" />
              </teleport>
            </div>
          </div>
        </div>

        <div v-if="$store.getters.isLoggedIn" class="bg-gray-200 px-4 py-2 rounded-lg flex flex-wrap justify-center self-center mt-8 text-center">
          <span>Signed in as {{ $store.state.user.email }}</span>
          <span
            class="cursor-pointer font-bold hover:underline lg:ml-2"
            @click="$store.dispatch('logout')">Log out</span>
        </div>
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { logEvent } from '@firebase/analytics';
import SubmitMessageModal from '../components/SendMessageModal.vue';
import { analytics } from '../firebase';
import { catchAndNotifyError } from '../notify';

export default {
  components: { SubmitMessageModal },
  data() {
    return {
      isFormOpen: false
    }
  },
  methods: {
    async login() {
      try {
        await this.$store.dispatch('login');
        if (this.$store.state.isSetupModalOpen) {
          logEvent(analytics, 'sign_up');
        } else {
          logEvent(analytics, 'login');
        }
      } catch (e) {
        catchAndNotifyError(this, e);
      } finally {
        this.$store.commit('SET_AUTH_LOADING', false);
      }
    },

    async searchMessageForm(e: SubmitEvent) {
      let target = e.target;
      if (target && target instanceof HTMLFormElement) {
        const studentId = target.children.namedItem('student_id');
        if (!studentId || !(studentId instanceof HTMLInputElement) || studentId.value.length == 0) return;
        this.$router.push({ name: 'message-wall-page', params: { recipientId: studentId.value } });
        target.reset();
      }
    },
  },
}
</script>