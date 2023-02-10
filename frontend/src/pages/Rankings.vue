<template>
  <main class="flex">
    <!-- Polish UI -->
    <!-- TODO: add empty state -->
    <div class="bg-white max-w-7xl shadow-lg w-full flex flex-col mx-auto self-start mt-4 p-6 lg:p-12 rounded-lg">
      <div class="flex flex-col md:flex-row justify-between items-center mb-8">
        <h1 class="text-center text-3xl font-bold">Valentine Ranking Board</h1>

        <div class="tabs tabs-boxed">
          <button 
            v-for="sex in availableSexes"
            @click="rankingsSex = sex.id" 
            :class="{ 'tab-active': rankingsSex == sex.id }" 
            class="tab tab-lg">{{ sex.label }}</button>
        </div>
      </div>

      <response-handler :query="query">
        <template #default>
          <table class="table w-full">
            <thead>
              <tr>
                <th class="w-1/4 text-lg normal-case text-red-400">
                  <span class="hidden lg:block">Ranking</span>
                </th>
                <th class="w-2/4 text-lg normal-case text-red-400">
                  <span class="block lg:hidden">Dept.</span>
                  <span class="hidden lg:block">Department</span>
                </th>
                <th class="w-1/4 text-lg text-center normal-case space-x-2 text-red-400">
                  <div class="space-x-2 flex justify-center items-center">
                    <span>Accumulated Coins</span>
                    <icon-coin />
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <template :key="`rankings_`+page" v-for="(ranking, page) in rankings?.pages">
                <tr :key="r.recipient_id" v-for="(r, i) in ranking.items" :class="{'border-b-2': i < ranking.items.length - 1}">
                  <td class="text-xl font-bold text-gray-700">#{{ i + 1 }}</td>
                  <td class="text-xl font-semibold text-gray-500">{{ r.department }}</td>
                  <td class="text-xl text-gray-500 text-center">
                    <div>
                      <span class="text-rose-500 font-bold">{{ r.total_coins }}</span>
                      <span class="text-2xl">ღ</span>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>

          <pagination-load-more-button 
            :should-go-next="hasNextPage" @click="fetchNextPage" />
        </template>

        <template #error="{ error }">
          <p>{{ error ? (error as Error).message : error }}</p>
        </template>
      </response-handler>
    </div>
  </main>
</template>

<script lang="ts" setup>
import ResponseHandler from "../components/ResponseHandler2.vue";
import PaginationLoadMoreButton from '../components/PaginationLoadMoreButton.vue';
import IconCoin from '~icons/twemoji/coin';
import { ref, watch } from "vue";
import { useInfiniteQuery, useQuery, useQueryClient } from "@tanstack/vue-query";
import { pb } from "../client";

const availableSexes = [
  {
    id: 'male',
    label: 'Male'
  },
  {
    id: 'female',
    label: 'Female'
  }
];

const queryClient = useQueryClient();
const rankingsSex = ref('male');

// TODO: integrate notiwind into tanstack query
const { fetchNextPage, hasNextPage, ...query } = useInfiniteQuery(
  ['rankings', rankingsSex], 
  ({ pageParam = 1 }) => pb.collection('rankings')
          .getList(pageParam, 10, { filter: `sex="${rankingsSex.value}"` }), {
    getNextPageParam: (result) => result.page + 1 <= result.totalPages ? result.page + 1 : undefined
  });

const rankings = query.data;
</script>