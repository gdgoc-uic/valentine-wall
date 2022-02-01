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

      <button v-if="links.next" @click="loadRankings({ url: links.next, merge: true })" class="mt-8 btn text-gray-900 px-12 self-center bg-white hover:bg-gray-100 border-gray-300 hover:border-gray-500">Load More</button>
    </div>
  </main>
</template>

<script lang="ts">
import { catchAndNotifyError } from '../notify';
export default {
  mounted() {
    this.loadRankings({});
  },
  data() {
    return {
      links: { first: null, last: null, next: null, previous: null },
      page: 1,
      perPage: 10,
      pageCount: 1,
      rankings: [],
      rankingsGender: 'male',
    }
  },
  watch: {
    rankingsGender(newVal, oldVal) {
      if (newVal == oldVal) return;
      this.loadRankings({});
    }
  },
  methods: {
    async loadRankings({ url, merge = false }: { url?: string|null, merge?: boolean }): Promise<void> {
      try {
        // const rankingsGender = this.rankingsGender;
        const rankingsGender = 'unknown';
        const { data: json } = await this.$client.get(url ?? '/rankings?limit=2');
        this.links = json['links'];
        this.page = json['page'];
        this.perPage = json['per_page'];
        this.pageCount = json['page_count'];
        this.rankings = merge ? this.rankings.concat(...json['data']) : json['data'];
      } catch (e) {
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>