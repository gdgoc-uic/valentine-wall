<template>
  <aside class="min-h-[30rem] bg-white shadow-md rounded-2xl flex flex-col">
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
</template>

<script lang="ts" setup>
import { useQuery } from '@tanstack/vue-query';
import { ref } from 'vue';
import { pb } from '../client';
import { useStore } from '../store_new';
import ResponseHandler from './ResponseHandler2.vue';
import IconCoin from '~icons/twemoji/coin';
import kingImg from '../assets/images/home/king.png';
import queenImg from '../assets/images/home/queen.png';
import { Record as PbRecord } from 'pocketbase';

const store = useStore();
const rankingsSex = ref('male');
const errorMessage = (e: unknown) => e instanceof Error ? e.message : 'Unknown error';

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

const rankingsQuery = useQuery(
  ['rankings_summary', rankingsSex],
  () => pb.collection('rankings').getList(1, 3, { 
    filter: `sex = "${rankingsSex.value}"`,
    expand: 'college_department'
  }),
  {
    refetchOnWindowFocus: () => false,
    refetchInterval: 10000 // refetch every 10 seconds
  }
);
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