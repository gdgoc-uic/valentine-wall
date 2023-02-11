<template>
  <response-handler :query="query">
    <template #default>
      <table class="table w-full">
        <thead>
          <tr>
            <th class="w-1/4 text-md normal-case text-red-400">
              <span class="hidden lg:block">ID</span>
            </th>
            <th class="w-2/4 text-md normal-case text-red-400">
              <span class="block lg:hidden">Desc.</span>
              <span class="hidden lg:block">Description</span>
            </th>
            <th class="w-1/4 text-md text-center normal-case text-red-400">
              <span>Amount</span>
            </th>
            <th class="w-1/4 text-md text-center normal-case text-red-400">
              <span>Date</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <template :key="'transactions_' + page" v-for="(records, page) in transactions?.pages">
            <tr :key="t.id" v-for="(t, i) in records.items" :class="{'border-b-2': i < (records.items.length ?? 0) - 1}">
              <td class="text-md font-bold text-gray-700">
                <span class="tooltip z-10" :data-tip="t.id">
                  {{t.id.substring(0, 6)}}
                </span>
              </td>
              <td class="text-md font-semibold text-gray-500">{{ t.description }}</td>
              <td class="text-md text-gray-500 text-center">{{ t.amount }}</td>
              <td class="text-md text-gray-500 text-right">{{ prettifyDateTime(t.created) }}</td>
            </tr>
          </template>
        </tbody>
      </table>

      <div class="flex justify-center">
        <pagination-load-more-button 
          :should-go-next="hasNextPage"
          @click="fetchNextPage" />
      </div>
    </template>

    <template #error="{ error }">
      <p>{{ (error as any).message || error }}</p>
    </template>
  </response-handler>
</template>

<script lang="ts" setup>
import ResponseHandler from '../../components/ResponseHandler2.vue'
import PaginationLoadMoreButton from '../../components/PaginationLoadMoreButton.vue'
import { prettifyDateTime } from '../../time_utils';
import { pb } from '../../client';
import { useInfiniteQuery } from '@tanstack/vue-query';

const { fetchNextPage, hasNextPage, ...query } = useInfiniteQuery(['transactions'], ({ pageParam = 1 }) => {
  console.log(pageParam);
  return pb.collection('virtual_transactions').getList(pageParam, 10, {
    sort: '-created'
  })
}, {
  keepPreviousData: true,
  getNextPageParam: (result) => {
    return result.page + 1 <= result.totalPages ? 
      result.page + 1 : undefined;
  }
});

const transactions = query.data;
</script>