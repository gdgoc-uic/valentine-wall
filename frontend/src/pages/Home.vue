<template>
  <main class="max-w-7xl mx-auto flex flex-col px-4">
    <section class="flex flex-col lg:flex-row">
      <section class="lg:w-1/3 pb-8 flex flex-col text-center space-y-8">
        <div>
          <img src="../assets/images/logo.png" class="w-4/5 md:w-2/3 lg:w-full pb-8 mx-auto" alt="Valentine Wall">
          <p class="text-gray-500 text-lg font-bold pb-4">Send, confess, and share your feelings anonymously!</p>
          <login-button v-if="!authState.isLoggedIn" class="btn-lg" />
        </div>
        <aside class="min-h-[40rem] bg-white shadow-md rounded-2xl flex flex-col">
          <div class="flex items-center space-x-4 rounded-t-2xl py-4 px-8 font-bold bg-rose-400 text-white">
            <img src="../assets/images/home/leaderboard.png">
            <p>Valentine Ranking Board</p>
          </div>
          <div class="tabs">
            <button
              v-for="sex in store.state.sexList"
              :key="'ranking_btn_' + sex.value"
              @click="rankingsSex = sex.value"
              :class="{ 'tab-active': rankingsSex == sex.value }"
              class="tab tab-lg flex-1 tab-bordered">{{ sex.label }}</button>
          </div>
          <div class="flex-1 flex flex-col ranking-board">
            <response-handler :query="rankingsQuery">
              <template #default>
                <!-- TODO: add empty state -->
                <div class="min-h-12 flex my-4 shadow ranking-info" :key="i" v-for="(r, i) in rankingsQuery.data.value?.items">
                  <div class="w-2/12 bg-black text-white inline-flex items-center justify-center font-bold ranking-placement">
                    {{ ordinalSuffixOf(i + 1) }}
                  </div>
                  <div class="flex-1 py-2 pl-2 inline-flex items-center">
                    <img
                      :src="r.sex == 'female' ? queenImg : kingImg"
                      class="w-2/12 mx-4" :alt="r.sex" />
                    <span class="font-bold">{{ (r.expand.college_department as PbRecord)?.uid ?? 'Unknown' }}</span>
                  </div>
                  <div class="flex w-3/12 px-2 bg-white">
                    <div class="pl-1 inline-flex items-center space-x-1 py-6">
                      <icon-coin />
                      <span>{{ r.total_coins }}ღ</span>
                    </div>
                  </div>
                </div>
              </template>
              <template #error="{ error }">
                <p>{{ errorMessage(error) }}</p>
              </template>
            </response-handler>
          </div>
          <router-link
            :to="{ name: 'rankings-page' }"
            class="btn w-full normal-case rounded-b-2xl rounded-t-none bg-rose-400 hover:bg-rose-500 border-none">
            Show all
          </router-link>
        </aside>
      </section>
      <section class="lg:w-2/3 flex flex-col space-y-2 md:space-y-8 lg:pl-8">
        <div class="bg-white p-12 space-y-2 rounded-2xl shadow-md">
          <div v-if="authState.isLoggedIn">
            <h2 class="text-3xl font-bold">Start writing your message!</h2>
            <send-message-form />
          </div>
          <div v-else>
            <div class="flex flex-col">
              <h2 class="text-2xl lg:text-4xl text-center font-bold mb-4">Welcome to UIC Valentine Wall!</h2>
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
              <p class="w-2/3 text-xl text-gray-500">Search your messages or even other's messages for free through the school ID.</p>
            </div>
            <search-form>
              <div class="form-control space-y-4">
                <div class="flex space-x-2 items-stretch">
                  <input type="text" class="flex-1 input input-lg input-bordered" name="recipient_id" placeholder="6 to 12-digit Student ID (e.g. 200xxxxxxxxx)">
                  <button class="btn bg-rose-600 hover:bg-rose-700 border-0 px-16 h-16">Search</button>
                </div>
              </div>
            </search-form>
          </div>
        </div>
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
import kingImg from '../assets/images/home/king.png';
import queenImg from '../assets/images/home/queen.png';

import SendMessageForm from '../components/SendMessageForm.vue';
import IconCoin from '~icons/twemoji/coin';
import SearchForm from '../components/SearchForm.vue';
import LoginButton from '../components/LoginButton.vue';
import ResponseHandler from '../components/ResponseHandler2.vue';
import MessageTiles from '../components/MessageTiles.vue';
import { pb } from '../client';
import { Record as PbRecord, UnsubscribeFunc } from 'pocketbase';
import { onMounted, onUnmounted, ref, reactive } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import { useAuth, useStore } from '../store_new';

// welcome modal but in home page
import IconWelcomeFeature1 from '~icons/home-icons/welcome_feature_1';
import IconWelcomeFeature2 from '~icons/home-icons/welcome_feature_2';
import IconWelcomeFeature3 from '~icons/home-icons/welcome_feature_3';

function ordinalSuffixOf(i: number): string {
  var j = i % 10,
    k = i % 100;
  if (j == 1 && k != 11) {
    return i + "st";
  }
  if (j == 2 && k != 12) {
    return i + "nd";
  }
  if (j == 3 && k != 13) {
    return i + "rd";
  }
  return i + "th";
}

const store = useStore();
const { state: authState } = useAuth();
const rankingsSex = ref('male');
const recentMessages = reactive<PbRecord[]>([]);
const recentsSSE = ref<UnsubscribeFunc>();
const errorMessage = (e: unknown) => e instanceof Error ? e.message : 'Unknown error';

const rankingsQuery = useQuery(
  ['rankings_summary', rankingsSex],
  () => pb.collection('rankings').getList(1, 3, { 
    filter: `sex = "${rankingsSex.value}"`,
    expand: 'college_department'
  }),
  {
    refetchOnWindowFocus: () => false,
  }
);

onMounted(() => {
  if (!import.meta.env.SSR) {
    pb.collection('messages').getList(1, 10, {
      sort: '-created'
    }).then(records => {
      recentMessages.push(...records.items);

      return pb.collection('messages').subscribe('*', (e) => {
        recentMessages.splice(0, recentMessages.length);
        recentMessages.push(e.record);
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

<style lang="postcss" scoped>
.ranking-info:nth-child(1) .ranking-placement {
  background: linear-gradient(42.38deg, rgba(139, 90, 0, 0.2) 0%, rgba(255, 255, 255, 0.13) 100%), #FFA500;
}

.ranking-info:nth-child(2) .ranking-placement {
  background: linear-gradient(42.38deg, rgba(119, 119, 119, 0.2) 0%, rgba(255, 255, 255, 0.13) 100%), #BABABA;
}

.ranking-info:nth-child(3) .ranking-placement {
  background: linear-gradient(0deg, rgba(156, 101, 0, 0.4), rgba(156, 101, 0, 0.4)), linear-gradient(49.39deg, #9C6500 0%, rgba(156, 101, 0, 0) 100%);
}
</style>