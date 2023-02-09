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
              <tr :key="r.recipient_id" v-for="(r, i) in rankings?.items" :class="{'border-b-2': i < (rankings?.items ?? []).length - 1}">
                <td class="text-xl font-bold text-gray-700">#{{ i + 1 }}</td>
                <td class="text-xl font-semibold text-gray-500">{{ r.department }}</td>
                <td class="text-xl text-gray-500 text-center">
                  <div>
                    <span class="text-rose-500 font-bold">{{ r.total_coins }}</span>
                    <span class="text-2xl">ღ</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <pagination-load-more-button 
            :should-go-next="(rankings?.page ?? 0) < (rankings?.totalPages ?? 0)" 
            @click="page++" />
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
import { useQuery, useQueryClient } from "@tanstack/vue-query";
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
const page = ref(1);
const rankingsSex = ref('male');

// TODO: integrate notiwind into tanstack query
// TODO: retrofit tanstack query to ResponseHandler
const query = useQuery(
  ['rankings', rankingsSex, page], 
  () => pb.collection('rankings')
          .getList(page.value, 10, { filter: `sex="${rankingsSex.value}"` }), {
    keepPreviousData: true
  });

const rankings = query.data;
</script>