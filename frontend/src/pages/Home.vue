<template>
  <main class="max-w-7xl mx-auto flex flex-col lg:flex-row px-4">
    <section class="lg:w-1/3 pb-8 flex flex-col text-center space-y-8">
      <div>
        <img src="../assets/images/logo.png" class="w-4/5 md:w-2/3 lg:w-full pb-8 mx-auto" alt="Valentine Wall">
        <p class="text-gray-500 text-lg font-bold pb-4">Send, confess, and share your feelings anonymously!</p>
        <login-button v-if="!$store.getters.isLoggedIn" class="btn-lg" />
        <button
          v-else
          @click="$store.commit('SET_SEND_MESSAGE_MODAL_OPEN', true)"
          class="btn btn-lg bg-rose-500 hover:bg-rose-600 normal-case border-0 shadow-md rounded-2xl w-2/3 lg:w-full space-x-3">
          <span>Start Writing</span>
          <icon-send />
        </button>
      </div>

      <aside class="bg-white min-h-[30rem] shadow-md rounded-2xl flex flex-col">
        <div class="flex items-center space-x-4 rounded-t-2xl py-4 px-8 font-bold bg-rose-400 text-white">
          <img src="../assets/images/home/leaderboard.png">
          <p>Valentine Ranking Board</p>
        </div>

        <div class="tabs">
          <button @click="rankingSex = 'male'" :class="{ 'tab-active': rankingSex == 'male' }" class="tab tab-lg flex-1 tab-bordered">Male</button>
          <button @click="rankingSex = 'female'" :class="{ 'tab-active': rankingSex == 'female' }" class="tab tab-lg flex-1 tab-bordered">Female</button>
        </div>

        <div class="flex-1 flex flex-col ranking-board">
          <paginated-response-handler :origin-endpoint="rankingsEndpoint">
            <template #default="{ data: rankings }">
              <!-- TODO: add empty state -->
              <div class="min-h-12 flex my-4 shadow ranking-info" :key="i" v-for="(r, i) in rankings">
                <div class="w-2/12 bg-black text-white inline-flex items-center justify-center font-bold ranking-placement">{{ ordinalSuffixOf(i + 1) }}</div>
                <div class="flex-1 py-2 pl-2 inline-flex items-center">
                  <img v-if="r.sex == 'female'" src="../assets/images/home/queen.png" class="w-2/12 mx-4" :alt="r.sex" />
                  <img v-else src="../assets/images/home/king.png" class="w-2/12 mx-4" :alt="r.sex" />
                  <!-- TODO: use shorthand dept name -->
                  <span class="font-bold">{{ r.department }}</span>
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
              <p>{{ error.message }}</p>
            </template>
          </paginated-response-handler>
        </div>

        <router-link 
          :to="{ name: 'rankings-page' }"
          class="btn w-full normal-case rounded-b-2xl rounded-t-none bg-rose-400 hover:bg-rose-500 border-none">
          Show all
        </router-link>
      </aside>
    </section>

    <section class="lg:w-2/3 flex flex-col space-y-2 md:space-y-8 lg:pl-8">
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

      <div>
        <h2 class="text-xl font-bold text-rose-600 mb-3">Recent Messages</h2>
        <message-tiles 
          :limit="20" 
          box-class="w-1/2 md:w-1/3"
          :messages="recentMessages" prepend />
      </div>
    </section>
  </main>
</template>

<script lang="ts">
import ClientOnly from '../components/ClientOnly.vue';
import Portal from '../components/Portal.vue';
import IconGift from '~icons/uil/gift';
import IconCoin from '~icons/twemoji/coin';
import IconEnvelope from '~icons/uil/envelope';
import SearchForm from '../components/SearchForm.vue';
import IconSend from '~icons/uil/message';
import LoginButton from '../components/LoginButton.vue';
import PaginatedResponseHandler from '../components/PaginatedResponseHandler.vue';
import MessageTiles from '../components/MessageTiles.vue';
import { expandAPIEndpoint } from '../client';

export default {
  components: { 
    ClientOnly, 
    Portal,
    IconGift,
    IconEnvelope,
    IconSend,
    IconCoin,
    SearchForm,
    LoginButton,
    PaginatedResponseHandler,
    MessageTiles
  },
  created() {
    if (!this.rankingsEndpoint.length) {
      this.rankingsEndpoint = this.getRankingsEndpoint();
    }
  },
  mounted() {
    if (!import.meta.env.SSR) {
      this.recentsSSE = new EventSource(expandAPIEndpoint('/recent-messages'));
      this.recentsSSE.onmessage = (event) => {
        if (event.data === 'null') return;
        // console.log(event);
        this.recentMessages = [JSON.parse(event.data)];
      }
    }
  },
  unmounted() {
    if (!import.meta.env.SSR) {
      this.recentsSSE.close();
    }
  },
  data() {
    return {
      isFormOpen: false,
      rankingSex: 'male',
      rankingsEndpoint: '',
      recentMessages: [] as any[],
      recentsSSE: null as unknown as EventSource
    }
  },
  watch: {
    rankingSex(newVal, oldVal) {
      if (newVal == oldVal) return;
      this.rankingsEndpoint = this.getRankingsEndpoint();
    }
  },
  methods: {
    getRankingsEndpoint(): string {
      const rankingSex = this.rankingSex;
      // const rankingSex = 'unknown';
      return `/rankings?limit=3&sex=${rankingSex}`;
    },
    ordinalSuffixOf(i: number): string {
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
    },
  },
}
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