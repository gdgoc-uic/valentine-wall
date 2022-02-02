<template>
  <main class="flex flex-col px-4">
    <div class="text-center flex flex-col items-center justify-center mb-4">
      <div class="pt-8 pb-6">
        <div v-if="$route.params.recipientId">
          <p class="text-3xl mb-2 text-gray-500">Messages for </p>
          <div class="indicator">
            <div class="indicator-item badge">{{ totalCount }}</div>
            <h2 class="text-6xl font-bold text-rose-600">{{ $route.params.recipientId }}</h2>
          </div>
        </div>
        <h2 v-else class="text-6xl font-bold text-rose-600">Recent Wall</h2>
      </div>
    </div>
    <client-only>
      <div class="bg-white rounded-xl shadow-lg flex justify-center">
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
    </client-only>

    <paginated-response-handler :origin-endpoint="endpoint">
      <template #default="{ data: messages, links, goto }">
        <!-- v-if="messages.length != 0"  -->
        <div class="flex flex-col -mx-2">
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

          <pagination-load-more-button :link="links.next" @click="goto(links.next, true)" />
        </div>
      </template>

      <template #error="{ error, isResponseError }">
        <div 
          v-if="(isResponseError && error.rawResponse.status == 404) || error.message == 'No messages found'" 
          class="text-center">
          <p>No messages found</p>
        </div>
      </template>
    </paginated-response-handler>
  </main>
</template>

<script lang="ts">
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import IconReply from '~icons/uil/comment-heart';
import { catchAndNotifyError } from '../notify';
import { defineComponent } from '@vue/runtime-core';
import ClientOnly from '../components/ClientOnly.vue';
import Masonry from '../components/Masonry.vue';
import { APIResponse } from '../client';
import PaginatedResponseHandler from '../components/PaginatedResponseHandler.vue';
import PaginationLoadMoreButton from '../components/PaginationLoadMoreButton.vue';

dayjs.extend(relativeTime);

export default defineComponent({
  components: {
    IconReply,
    ClientOnly,
    Masonry,
    PaginatedResponseHandler,
    PaginationLoadMoreButton
  },
  created() {
    this.endpoint = this.getMessagesEndpoint({ hasGift: this.hasGift });
  },
  data() {
    return {
      stats: { messages_count: 0, gift_messages_count: 0 },
      // TODO: disable_restricted_access_to_gift_messages
      // hasGift: false
      hasGift: null as boolean | null,
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
    hasGift(newVal: boolean, oldVal: boolean) {
      if (newVal === oldVal) return;
      this.endpoint = this.getMessagesEndpoint({ hasGift: newVal });
    },
    '$route'() {
      this.stats = { messages_count: 0, gift_messages_count: 0 };
      this.endpoint = this.getMessagesEndpoint({ hasGift: this.hasGift });
      this.loadData();
    }
  },
  methods: {
    loadData(): Promise<any> {
      return this.loadStats();
    },
    checkMessagesLength({ data }: APIResponse) {
        if (Array.isArray(data['data']) && data['data'].length == 0) {
            throw new Error('No messages found.');
        } else {
            throw new Error('Data not an array.');
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
      return `/messages/${recipientId}?order=created_at,desc&limit=10&${hasGift == null ? 'has_gift=2' : hasGift ? 'has_gift=1' : 'has_gift=0'}`
    },
    humanizeTime(date: Date | string): string {
      return dayjs(date).toNow(true);
    }
  }
})
</script>