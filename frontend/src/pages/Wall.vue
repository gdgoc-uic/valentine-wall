<template>
  <main class="min-h-[60vh] flex flex-col px-4">
    <div class="text-center flex flex-col items-center justify-center mb-4">
      <div class="pt-8 pb-6">
        <div v-if="$route.params.recipientId">
          <p class="text-3xl mb-2 text-gray-500">Messages for </p>
          <div class="indicator">
            <div v-if="totalCount != 0" class="indicator-item badge">{{ totalCount }}</div>
            <h2 class="text-5xl lg:text-6xl font-bold text-rose-600">{{ $route.params.recipientId }}</h2>
          </div>
        </div>
        <h2 v-else class="text-5xl lg:text-6xl font-bold text-rose-600">Recent Wall</h2>
      </div>
    </div>
    <client-only>
      <div class="max-w-7xl mx-auto w-full bg-white rounded-xl shadow-lg flex justify-center">
        <div 
          v-if="$route.params.recipientId && $store.getters.isLoggedIn && $store.state.user.associatedId === $route.params.recipientId" 
          class="tabs">
          <button @click="hasGift = null" :class="{'tab-active': hasGift == null}" class="tab tab-lg px-12 space-x-2 tab-bordered">
            <span>All</span>
          </button>
          <button @click="hasGift = false" :class="{'tab-active': hasGift == false}" class="tab tab-lg px-6 space-x-2 tab-bordered">
            <span>Messages</span>
            <span class="badge">{{ stats.messages_count }}</span>
          </button>
          <button @click="hasGift = true" :class="{'tab-active': hasGift == true}" class="tab tab-lg px-10 space-x-2 tab-bordered">
            <span>Gifts</span>
            <span class="badge">{{ stats.gift_messages_count }}</span>
          </button>
        </div>
      </div>

      <section 
        v-if="!$store.getters.isLoggedIn && $route.params.recipientId" 
        class="max-w-7xl space-y-8 lg:space-y-0 mx-auto mt-3 p-8 w-full border bg-[#ffeded] border-rose-500 rounded-xl shadow-lg flex flex-col lg:flex-row text-center lg:text-left items-center justify-between">
        <div class="flex flex-col">
          <h3 class="text-2xl lg:text-3xl mb-1 text-rose-400 font-semibold">Are you {{ $route.params.recipientId }}?</h3>
          <p class="text-lg lg:text-xl">Join now and get access to exclusive features!</p>
        </div>

        <login-button />
      </section>
    </client-only>

    <paginated-response-handler :origin-endpoint="endpoint" :fail-fn="checkMessagesLength">
      <template #default="{ data: messages, links, goto }">
        <div class="max-w-7xl w-full mx-auto flex flex-col">
          <message-tiles :messages="messages" replace />
          <pagination-load-more-button :link="links.next" @click="goto(links.next, true)" />
        </div>
      </template>

      <template #error="{ error, isResponseError }">
        <div 
          class="text-center w-full py-16 h-full flex flex-col items-center">
          <template 
            v-if="
              (isResponseError && error.rawResponse.status == 404) || 
              (error.message && error.message == 'No messages found.')
            ">
            <p class="text-4xl">Nothing to see here!</p>
          </template> 
          <template v-else-if="error.message">
            <p class="text-4xl">Something went wrong.</p>
          </template>
        </div>
      </template>
    </paginated-response-handler>
  </main>
</template>

<script lang="ts">
import { catchAndNotifyError } from '../notify';
import { defineComponent } from '@vue/runtime-core';
import ClientOnly from '../components/ClientOnly.vue';
import { APIResponse } from '../client';
import PaginatedResponseHandler from '../components/PaginatedResponseHandler.vue';
import PaginationLoadMoreButton from '../components/PaginationLoadMoreButton.vue';
import LoginButton from '../components/LoginButton.vue';
import MessageTiles from '../components/MessageTiles.vue';

export default defineComponent({
  components: {
    ClientOnly,
    PaginatedResponseHandler,
    PaginationLoadMoreButton,
    LoginButton,
    MessageTiles
  },
  created() {
    if (this.$route.params.recipientId && this.$store.getters.isLoggedIn) {
      this.hasGift = null;
    }
    this.endpoint = this.getMessagesEndpoint({ hasGift: this.hasGift });
  },
  mounted() {
    if (this.stats.messages_count == 0 && !import.meta.env.SSR) {
      this.loadData();
    }
  },
  data() {
    return {
      stats: { messages_count: 0, gift_messages_count: 0 },
      // TODO: disable_restricted_access_to_gift_messages
      hasGift: false as boolean | null,
      // hasGift: null as boolean | null,
      endpoint: ''
    };
  },
  computed: {
    totalCount(): number {
      if (import.meta.env.SSR) {
        return 0;
      } else if (this.$store.getters.isLoggedIn 
        && this.$store.state.user.associatedId === this.$route.params.recipientId) {
        return this.stats.gift_messages_count + this.stats.messages_count;
      } else {
        return this.stats.messages_count;
      }
    }
  },
  watch: {
    hasGift(newVal, oldVal) {
      if (newVal === oldVal) return;
      this.endpoint = this.getMessagesEndpoint({ hasGift: newVal });
    },
    '$route'() {
      this.stats = { messages_count: 0, gift_messages_count: 0 };
      this.endpoint = this.getMessagesEndpoint({ hasGift: this.hasGift });
      this.loadData();
    },
  },
  methods: {
    loadData(): Promise<any> {
      return this.loadStats();
    },
    checkMessagesLength({ data }: APIResponse) {
      if (!Array.isArray(data['data']) || data['data'].length == 0) {
        throw new Error('No messages found.');
      }
    },
    async loadStats(): Promise<void> {
      try {
        if (!this.$route.params.recipientId) return;
        const recipientId = this.$route.params.recipientId ?? '';
        const { data: json } = await this.$client.get(`/messages/${recipientId}/stats`);
        this.stats = json;
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    },
    getMessagesEndpoint({ hasGift = false }: { hasGift?: boolean | null }): string {
      const recipientId = this.$route.params.recipientId ?? '';
      return `/messages/${recipientId}?order=created_at,desc&limit=10&has_gift=${hasGift == null ? '2' : hasGift ? '1' : '0'}`
    }
  }
})
</script>