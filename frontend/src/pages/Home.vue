<template>
  <main class="max-w-7xl mx-auto flex flex-col px-4">
    <section class="flex flex-col lg:flex-row">
      <section class="lg:w-1/3 flex flex-col text-center space-y-8">
        <div>
          <img src="../assets/images/logo.png" class="w-4/5 md:w-2/3 lg:w-full pb-8 mx-auto" alt="Valentine Wall">
          <p class="text-gray-500 text-lg font-bold pb-4">Send, confess, and share your feelings anonymously!</p>
          <login-button v-if="!authState.isLoggedIn" class="btn-lg mb-8" />
        </div>

        <rankings-box class="hidden md:flex" />
      </section>
      <section class="lg:w-2/3 flex flex-col space-y-2 md:space-y-8 lg:pl-8">
        <div class="bg-white p-6 md:p-12 space-y-2 rounded-2xl shadow-md">
          <div v-if="authState.isLoggedIn">
            <h2 class="text-center md:text-left text-xl md:text-3xl font-bold">Start writing your message!</h2>
            <send-message-form />
          </div>
          <div v-else>
            <div class="flex flex-col">
              <h2 class="text-xl md:text-2xl lg:text-4xl text-center font-bold mb-4">Welcome to UIC Valentine Wall!</h2>
              <div class="flex flex-col text-lg md:text-2xl lg:text-lg">
                <div class="p-3 w-full flex flex-row items-center justify-center text-left">
                  <div class="w-1/3 md:w-1/4 bg-rose-400 rounded-full p-4 mr-8">
                    <icon-welcome-feature-1 class="w-full h-full text-white" />
                  </div>
                  <p>Post, confess, or share your thoughts to your fellow Ignacian Marian anonymously!</p>
                </div>
                <div class="p-3 w-full flex flex-row items-center justify-center text-left">
                  <div class="w-1/3 md:w-1/4 bg-blue-400 rounded-full p-4 mr-8">
                    <icon-welcome-feature-2 class="w-full h-full text-rose-100" />
                  </div>
                  <p>No money? No problem! You can also send virtual gifts alongside your message!</p>
                </div>
                <div class="p-3 w-full flex flex-row items-center justify-center text-left">
                  <div class="w-1/3 md:w-1/4 bg-orange-300 rounded-full p-4 mr-8">
                    <icon-welcome-feature-3 class="w-full h-full text-white" />
                  </div>
                  <p>Discover many more messages sent by others through our public wall!</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="hidden md:block w-full">
          <div class="bg-white p-12 space-y-8 rounded-2xl shadow-md h-full">
            <div>
              <h2 class="text-3xl font-bold mb-4">Search Messages</h2>
              <p class="w-2/3 text-xl text-gray-500">Find your messages or even other's messages.</p>
            </div>
            <search-form>
              <div class="form-control space-y-4">
                <div class="flex space-x-2 items-stretch">
                  <input type="text" class="flex-1 input input-lg input-bordered" name="recipient_id" placeholder="Type 'everyone' or 6 to 12-digit Student ID">
                  <button class="btn bg-rose-600 hover:bg-rose-700 border-0 px-16 h-16">Search</button>
                </div>
              </div>
            </search-form>
          </div>
        </div>
        <rankings-box class="md:hidden" />
      </section>
    </section>

    <div class="w-full mt-8">
      <h2 class="text-xl font-bold text-rose-600 mb-3">Recent Messages</h2>
      <message-tiles 
        :limit="12"
        :messages="recentMessages"
        prepend />
    </div>
  </main>
</template>

<script lang="ts" setup>
import SendMessageForm from '../components/SendMessageForm.vue';
import SearchForm from '../components/SearchForm.vue';
import LoginButton from '../components/LoginButton.vue';
import MessageTiles from '../components/MessageTiles.vue';
import { pb } from '../client';
import { Record as PbRecord, UnsubscribeFunc } from 'pocketbase';
import { onMounted, onUnmounted, ref, reactive } from 'vue';
import { useAuth } from '../store_new';

// welcome modal but in home page
import IconWelcomeFeature1 from '~icons/home-icons/welcome_feature_1';
import IconWelcomeFeature2 from '~icons/home-icons/welcome_feature_2';
import IconWelcomeFeature3 from '~icons/home-icons/welcome_feature_3';
import RankingsBox from '../components/RankingsBox.vue';

const { state: authState } = useAuth();
const recentMessages = reactive<PbRecord[]>([]);
const recentsSSE = ref<UnsubscribeFunc>();

onMounted(() => {
  if (!import.meta.env.SSR) {
    pb.collection('messages').getList(1, 10, {
      filter: 'recipient = "everyone"',
      sort: '-created',
      expand: 'gifts'
    }).then(records => {
      recentMessages.push(...records.items);

      return pb.collection('messages').subscribe('*', (e) => {
        if (e.action !== 'create' || e.record.recipient !== 'everyone') return;

        recentMessages.unshift(e.record);
      });
    }).then(unsub => {
      if (!unsub) return;
      recentsSSE.value = unsub;
    });
  }
});

onUnmounted(() => {
  if (!import.meta.env.SSR) {
    recentsSSE.value?.();
  }
});
</script>