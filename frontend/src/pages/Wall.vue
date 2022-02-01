<template>
  <main class="flex flex-col px-4">
    <!-- TODO: loader -->
    <div class="text-center flex flex-col items-center justify-center mb-4">
      <div class="pt-8 pb-6">
        <div v-if="$route.params.recipientId">
          <p class="text-3xl mb-2 text-gray-500">Messages for </p>
          <h2 class="text-6xl font-bold text-rose-600">{{ $route.params.recipientId }}</h2>
        </div>
        <h2 v-else class="text-6xl font-bold text-rose-600">Recent Wall</h2>
      </div>
    </div>
    <div class="bg-white rounded-xl shadow-lg flex justify-center">
      <client-only>
          <!-- v-if="$store.getters.isLoggedIn && $store.state.user.associatedId === $route.params.recipientId"  -->
        <div v-if="$route.params.recipientId" class="tabs">
          <button @click="hasGift = null" :class="{'tab-active': hasGift == null}" class="tab tab-lg px-12 space-x-2 tab-bordered">
            <span>All</span>
            <span class="badge">{{ totalCount }}</span>
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
      </client-only>
    </div>

    <div v-if="messages.length != 0" class="flex flex-col -mx-2">
      <masonry>
        <!-- TODO: make card widths and colors different -->
        <div :key="msg.id" v-for="msg in messages" class="w-1/5 p-2 min-h-16 block">
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
      </masonry>

      <button v-if="links.next" @click="loadMessages({ url: links.next, merge: true, hasGift })" class="mt-8 btn text-gray-900 px-12 self-center bg-white hover:bg-gray-100 border-gray-300 hover:border-gray-500">Load More</button>
    </div>

    <div v-else class="text-center">
      <p>No messages found</p>
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
import { defineComponent } from '@vue/runtime-core';
import ClientOnly from '../components/ClientOnly.vue';
import Masonry from '../components/Masonry.vue';


dayjs.extend(relativeTime);

export default defineComponent({
  components: {
    IconReply,
    ClientOnly,
    Masonry
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
      totalCount: 0,
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
        const endpoint = url ?? `/messages/${recipientId}?order=created_at,desc&limit=10&${hasGift == null ? 'has_gift=2' : hasGift ? 'has_gift=1' : 'has_gift=0'}`;
        const { data: json, rawResponse: resp } = await this.$client.get(endpoint);
        this.links = json['links'];
        this.page = json['page'];
        this.perPage = json['per_page'];
        this.pageCount = json['page_count'];
        this.totalCount = json['total_count'];
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