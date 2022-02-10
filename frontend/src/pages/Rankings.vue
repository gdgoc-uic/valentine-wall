<template>
  <main class="flex">
    <!-- Polish UI -->
    <!-- TODO: add empty state -->
    <div class="bg-white max-w-7xl shadow-lg w-full flex flex-col mx-auto self-start mt-4 p-6 lg:p-12 rounded-lg">
      <div class="flex flex-col md:flex-row justify-between items-center mb-8">
        <h1 class="text-center text-3xl font-bold">Valentine Ranking Board</h1>

        <div class="tabs tabs-boxed">
          <button @click="rankingsSex = 'male'" :class="{ 'tab-active': rankingsSex == 'male' }" class="tab tab-lg">Male</button>
          <button @click="rankingsSex = 'female'" :class="{ 'tab-active': rankingsSex == 'female' }" class="tab tab-lg">Female</button>
        </div>
      </div>

      <paginated-response-handler :origin-endpoint="endpoint">
        <template #default="{ data: rankings, links, goto }">
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
              <tr :key="r.recipient_id" v-for="(r, i) in rankings" :class="{'border-b-2': i < rankings.length - 1}">
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

          <pagination-load-more-button :link="links.next" @click="goto(links.next, true)" />
        </template>

        <template #error="{ error }">
          <p>{{ error.message || error }}</p>
        </template>
      </paginated-response-handler>
    </div>
  </main>
</template>

<script lang="ts">
import PaginatedResponseHandler from "../components/PaginatedResponseHandler.vue";
import PaginationLoadMoreButton from '../components/PaginationLoadMoreButton.vue';
import IconCoin from '~icons/twemoji/coin';

export default {
  components: { 
    PaginatedResponseHandler,
    PaginationLoadMoreButton,
    IconCoin
  },
  created() {
    this.endpoint = this.getRankingsEndpoint();
  },
  data() {
    return {
      endpoint: '',
      rankingsSex: 'male',
    }
  },
  watch: {
    rankingsSex(newVal, oldVal) {
      if (newVal == oldVal) return;
      this.endpoint = this.getRankingsEndpoint();
    }
  },
  methods: {
    getRankingsEndpoint(): string {
      const rankingsSex = this.rankingsSex;
      // const rankingsSex = 'unknown';
      return `/rankings?limit=10&sex=${rankingsSex}`;
    },
  }
}
</script>