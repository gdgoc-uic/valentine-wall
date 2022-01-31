<template>
  <main class="min-h-screen flex">
    <div class="bg-white max-w-4xl w-full mx-auto self-start mt-4 p-12 rounded-lg">
      <div class="text-center pb-12 flex flex-col items-center justify-center">
        <template v-if="$route.params.recipientId">
          <p class="text-2xl mb-2 text-gray-500">Messages for </p>
          <h2 class="text-5xl font-bold">{{ $route.params.recipientId }}</h2>
          <p>{{ stats.messages_count }} Messages, {{ stats.gift_messages_count }} Gift Messages</p>

          <!-- TODO: -->
          <client-only>
            <div v-if="$store.getters.isLoggedIn && $store.state.user.associatedId === $route.params.recipientId" class="btn-group mt-8">
              <button @click="hasGift = false" :class="toggleBtnStyling(hasGift == false)" class="btn btn-md px-8">Messages</button> 
              <button @click="hasGift = true" :class="toggleBtnStyling(hasGift == true)" class="btn btn-md px-8">Gifts</button>
              <button @click="hasGift = null" :class="toggleBtnStyling(hasGift == null)" class="btn btn-md px-8">All</button>
            </div>
          </client-only>
        </template>
        <template v-else>
          <h2 class="text-5xl font-bold">Recent Wall</h2>
        </template>
      </div>

      <div v-if="messages.length != 0" class="flex flex-col">
        <div class="flex flex-wrap items-stretch justify-center">
          <div :key="msg.id" v-for="msg in messages" class="w-1/3 p-2  min-h-16 block">
            <router-link 
              :to="{ name: 'message-page', params: { recipientId: msg.recipient_id, messageId: msg.id } }" 
              style="background: linear-gradient(to bottom,rgb(254, 243, 199) 21px,#00b0d7 1px); background-size: 100% 22px;"
              class="rounded-lg flex flex-col justify-between border shadow-lg h-64 hover:scale-110 transition-transform">
                <div class="px-6 pt-6">
                  <p>{{ msg.content }}</p>
                </div>
                <div class="bg-amber-100  px-6 pb-6 flex justify-between items-center mt-4">
                  <p class="text-gray-500 text-sm">{{ humanizeTime(msg.created_at) }} ago</p>
                  <icon-reply v-if="msg.has_replied" class="text-pink-500" />
                </div>
            </router-link>
          </div>
        </div>

        <button v-if="links.next" @click="loadMessages({ url: links.next, merge: true, hasGift })" class="mt-8 btn text-gray-900 px-12 self-center bg-white hover:bg-gray-100 border-gray-300 hover:border-gray-500">Load More</button>
      </div>

      <div v-else class="text-center">
        <p>No messages found</p>
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import IconReply from '~icons/uil/comment-heart';
import { analytics } from '../firebase';
import { logEvent } from '@firebase/analytics';
import { catchAndNotifyError } from '../notify';
import ClientOnly from '../components/ClientOnly.vue';
import { defineComponent } from '@vue/runtime-core';

dayjs.extend(relativeTime);

export default defineComponent({
  components: {
    IconReply,
    ClientOnly
  },
  async serverPrefetch(): Promise<void> {
    await this.loadData();
  },
  mounted() {
    if (!this.messages.length) {
      this.loadData();
    }
  },
  data() {
    return {
      stats: { messages: 0, gift_messages: 0 },
      links: { 
        first: null as string | null, 
        last: null as string | null, 
        next: null as string | null, 
        previous: null as string | null
      },
      page: 1,
      perPage: 10,
      pageCount: 1,
      messages: [] as any[],
      // TODO: disable_restricted_access_to_gift_messages
      // hasGift: false
      hasGift: null as boolean | null
    };
  },
  watch: {
    hasGift(newVal: boolean, oldVal: boolean) {
      if (newVal === oldVal) return;
      this.loadMessages({ hasGift: newVal });
    },
    '$route'() {
      this.messages = [];
      this.stats = { messages: 0, gift_messages: 0 };
      this.loadData();
    }
  },
  methods: {
    toggleBtnStyling(hasGift: boolean): string[] {
      return [
        hasGift
          ? 'text-white border-rose-700 bg-rose-500 hover:bg-rose-600 hover:border-rose-800'
          : 'text-gray-900 border-gray-300 bg-white hover:bg-rose-50 hover:border-rose-400'
      ]
    },
    loadData(): Promise<[any, any]> {
      return Promise.all([
        this.loadMessages({ hasGift: this.hasGift }),
        this.loadStats()
      ]);
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
    async loadMessages({ hasGift = false, url, merge = false }: { hasGift?: boolean | null, url?: string | null, merge?: boolean }): Promise<void> {
      try {
        const recipientId = this.$route.params.recipientId ?? '';
        const endpoint = url ?? `/messages/${recipientId}?order=created_at,desc&limit=6&${hasGift == null ? 'has_gift=2' : hasGift ? 'has_gift=1' : 'has_gift=0'}`;
        const { data: json, rawResponse: resp } = await this.$client.get(endpoint);
        this.links = json['links'];
        this.page = json['page'];
        this.perPage = json['per_page'];
        this.pageCount = json['page_count'];
        this.messages = merge ? this.messages.concat(...json['data']) : json['data'];
        logEvent(analytics!, 'search_messages', { recipient_id: recipientId, status_code: resp.status });
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    },
    humanizeTime(date: Date | string): string {
      return dayjs(date).toNow(true);
    }
  }
})
</script>