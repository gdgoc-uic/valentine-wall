<template>
  <main class="max-w-7xl mx-auto flex flex-col px-4">
    <section class="py-8 self-center flex flex-col items-center text-center">
      <img src="../assets/images/logo.png" class="w-1/2 pb-8" alt="Valentine Wall">
      <p class="text-gray-500 text-lg font-bold pb-4">Send, confess, and share your feelings to the person you admire the most anonymously!</p>

      <button v-if="!$store.getters.isLoggedIn" @click="login" class="btn btn-lg text-gray-500 bg-white hover:bg-gray-100 hover:text-gray-500 normal-case border-0 shadow-md space-x-3">
        <google-icon />
        <span>Sign in with UIC Google</span>
      </button>
      <button 
        v-else 
        @click="$store.commit('SET_SEND_MESSAGE_MODAL_OPEN', true)" 
        class="btn btn-lg bg-rose-500 hover:bg-rose-600 normal-case border-0 shadow-md space-x-3">
        <span>Start Writing</span>
        <icon-send />
      </button>
    </section>

    <section class="flex mt-8 -mx-8">
      <div class="px-8 w-2/3">
        <div class="bg-white p-16 space-y-8 rounded-2xl shadow-md h-full">
          <div>
            <h2 class="text-4xl font-bold mb-4">Search Message</h2>
            <p class="text-xl text-gray-500">Lorem ipsum dolor sit amet consectetur adipisicing elit.</p>
          </div>

          <search-form>
            <div class="form-control space-y-4">
              <input type="text" class="input input-bordered" name="recipient_id" placeholder="12-digit Student ID">
              <button class="btn self-end px-16">Search</button>
            </div>
          </search-form>
        </div>
      </div>

      <div class="px-8 w-1/3">
        <aside class="bg-white h-full shadow-md rounded-2xl flex flex-col">
          <div class="flex items-center space-x-4 rounded-t-2xl py-4 px-8 font-bold bg-rose-400 text-white">
            <img src="../assets/images/home/leaderboard.png">
            <p>Valentine Ranking Board</p>
          </div>

          <div class="tabs">
            <button @click="rankingGender = 'male'" :class="{ 'tab-active': rankingGender == 'male' }" class="tab tab-lg flex-1 tab-bordered">Male</button>
            <button @click="rankingGender = 'female'" :class="{ 'tab-active': rankingGender == 'female' }" class="tab tab-lg flex-1 tab-bordered">Female</button>
          </div>

          <div class="flex-1 flex flex-col ranking-board">
            <!-- TODO: add empty state -->
            <div class="min-h-12 flex my-4 shadow ranking-info" :key="i" v-for="(r, i) in rankings">
              <div class="w-2/12 bg-black text-white inline-flex items-center justify-center font-bold ranking-placement">{{ ordinalSuffixOf(i + 1) }}</div>
              <div class="flex-1 py-2 pl-2 inline-flex items-center">
                <img v-if="r.gender == 'female'" src="../assets/images/home/queen.png" class="w-2/12 mx-4" :alt="r.gender" />
                <img v-else src="../assets/images/home/king.png" class="w-2/12 mx-4" :alt="r.gender" />
                <!-- TODO: use shorthand dept name -->
                <span class="font-bold">{{ r.department }}</span>
              </div>
              <div class="flex w-3/12 px-2 bg-white">
                <div class="w-1/2 pr-1 inline-flex items-center space-x-1">
                  <icon-envelope />
                  <span>{{ r.messages_count }}</span>
                </div>

                <div class="w-1/2 pl-1 inline-flex items-center space-x-1">
                  <icon-gift />
                  <span>{{ r.gift_messages_count }}</span>
                </div>
              </div>
            </div>
          </div>

          <router-link 
            :to="{ name: 'rankings-page' }"
            class="btn w-full normal-case rounded-b-2xl rounded-t-none bg-rose-400 hover:bg-rose-500 border-none">
            Show all
          </router-link>
        </aside>
      </div>
    </section>
  </main>
</template>

<script lang="ts">
import { catchAndNotifyError } from '../notify';
import ClientOnly from '../components/ClientOnly.vue';
import Portal from '../components/Portal.vue';
import GoogleIcon from '~icons/logos/google-icon';
import IconGift from '~icons/uil/gift';
import IconEnvelope from '~icons/uil/envelope';
import SearchForm from '../components/SearchForm.vue';
import IconSend from '~icons/uil/message';

export default {
  components: { 
    ClientOnly, 
    Portal, 
    GoogleIcon,
    IconGift,
    IconEnvelope,
    IconSend,
    SearchForm
  },
  mounted() {
    if (!this.rankings.length) {
      this.loadRankings();
    }
  },
  data() {
    return {
      isFormOpen: false,
      rankingGender: 'male',
      rankings: []
    }
  },
  watch: {
    rankingGender(newVal, oldVal) {
      if (newVal == oldVal) return;
      this.rankings = [];
      this.loadRankings();
    }
  },
  methods: {
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
    async login() {
      try {
        await this.$store.dispatch('login');
      } catch (e) {
        catchAndNotifyError(this, e);
      }
    },
    async loadRankings(): Promise<void> {
      try {
        // const rankingGender = this.rankingGender;
        const rankingGender = 'unknown';
        const { data: json } = await this.$client.get(`/rankings?limit=3&gender=${rankingGender}`);
        this.rankings = json['data'];
      } catch (e) {
        catchAndNotifyError(this, e);
      }
    }
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