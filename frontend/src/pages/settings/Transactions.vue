<template>
  <paginated-response-handler origin-endpoint="/user/transactions">
    <template #default="{ data: transactions, links, goto }">
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
          <tr :key="t.id" v-for="(t, i) in transactions" :class="{'border-b-2': i < transactions.length - 1}">
            <td class="text-md font-bold text-gray-700">{{t.id}}</td>
            <td class="text-md font-semibold text-gray-500">{{ t.description }}</td>
            <td class="text-md text-gray-500 text-center">{{ t.amount }}</td>
            <td class="text-md text-gray-500 text-center">{{ t.created_at }}</td>
          </tr>
        </tbody>
      </table>

      <pagination-load-more-button :link="links.next" @click="goto(links.next, true)" />
    </template>

    <template #error="{ error }">
      <p>{{ error.message || error }}</p>
    </template>
  </paginated-response-handler>
</template>

<script>
import PaginatedResponseHandler from '../../components/PaginatedResponseHandler.vue'
import PaginationLoadMoreButton from '../../components/PaginationLoadMoreButton.vue'
export default {
  components: { PaginatedResponseHandler, PaginationLoadMoreButton },
}
</script>