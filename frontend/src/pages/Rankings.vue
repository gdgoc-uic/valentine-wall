<template>
  <main class="min-h-screen flex">
    <!-- Polish UI -->
    <div class="bg-white max-w-4xl w-full mx-auto self-start mt-4 p-12 rounded-lg">
      <h1 class="text-center text-7xl font-bold mb-8">Rankings</h1>

      <table class="table w-full ">
        <thead>
          <tr>
            <th></th>
            <th>Student ID</th>
            <th>Messages</th>
            <th>Gift Messages</th>
          </tr>
        </thead>
        <tbody>
          <tr :key="r.recipient_id" v-for="(r, i) in rankings">
            <td>{{ i + 1 }}</td>
            <td>{{ r.recipient_id }}</td>
            <td>{{ r.messages_count }}</td>
            <td>{{ r.gift_messages_count }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<script>
import client from '../client'
import { catchAndNotifyError } from '../notify';
export default {
  mounted() {
    this.loadRankings();
  },
  data() {
    return {
      rankings: []
    }
  },
  methods: {
    async loadRankings() {
      try {
        const { data: json } = await client.get('/rankings');
        this.rankings = json;
      } catch (e) {
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>