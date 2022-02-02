<template>
  <main class="flex">
    <!-- Polish UI -->
    <!-- TODO: add empty state -->
    <div class="bg-white max-w-7xl shadow-lg w-full flex flex-col mx-auto self-start mt-4 p-12 rounded-lg">
      <div class="flex flex-row justify-between items-center mb-8">
        <h1 class="text-center text-3xl font-bold">Valentine Ranking Board</h1>

        <div class="tabs tabs-boxed">
          <button @click="rankingsGender = 'male'" :class="{ 'tab-active': rankingsGender == 'male' }" class="tab tab-lg">Male</button>
          <button @click="rankingsGender = 'female'" :class="{ 'tab-active': rankingsGender == 'female' }" class="tab tab-lg">Female</button>
        </div>
      </div>

      <paginated-response-handler :origin-endpoint="endpoint">
        <template #default="{ data: rankings, links, goto }">
          <table class="table w-full">
            <thead>
              <tr>
                <th class="w-1/4 text-lg normal-case text-red-400">Ranking</th>
                <th class="w-2/4 text-lg normal-case text-red-400">Department</th>
                <th class="w-1/4 text-lg text-center normal-case text-red-400">Messages</th>
                <th class="w-1/4 text-lg text-center normal-case text-red-400">Gift Messages</th>
              </tr>
            </thead>
            <tbody>
              <tr :key="r.recipient_id" v-for="(r, i) in rankings" :class="{'border-b-2': i < rankings.length - 1}">
                <td class="text-xl font-bold text-gray-700">#{{ i + 1 }}</td>
                <td class="text-xl font-semibold text-gray-500">{{ r.department }}</td>
                <td class="text-xl text-gray-500 text-center">{{ r.messages_count }}</td>
                <td class="text-xl text-gray-500 text-center">{{ r.gift_messages_count }}</td>
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

export default {
  components: { 
    PaginatedResponseHandler,
    PaginationLoadMoreButton 
  },
  created() {
    this.endpoint = this.getRankingsEndpoint();
  },
  data() {
    return {
      endpoint: '',
      rankingsGender: 'male',
    }
  },
  watch: {
    rankingsGender(newVal, oldVal) {
      if (newVal == oldVal) return;
      this.endpoint = this.getRankingsEndpoint();
    }
  },
  methods: {
    getRankingsEndpoint(): string {
      // const rankingsGender = this.rankingsGender;
      const rankingsGender = 'unknown';
      return `/rankings?limit=2&gender=${rankingsGender}`;
    },
  }
}
</script>