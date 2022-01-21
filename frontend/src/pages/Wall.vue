<template>
  <main class="min-h-screen flex">
    <div class="bg-white max-w-4xl w-full mx-auto self-start mt-4 p-12 rounded-lg">
      <div class="text-center pb-12 flex flex-col items-center justify-center">
        <template v-if="$route.params.recipientId">
          <p class="text-2xl mb-2 text-gray-500">Messages for </p>
          <h2 class="text-5xl font-bold">{{ $route.params.recipientId }}</h2>

          <div v-if="$store.getters.isLoggedIn && $store.state.user.associatedId === $route.params.recipientId" class="btn-group mt-8">
            <button @click="hasGift = false" :class="toggleBtnStyling(hasGift == false)" class="btn btn-md px-8">Notes</button> 
            <button @click="hasGift = true" :class="toggleBtnStyling(hasGift == true)" class="btn btn-md px-8">Gifts</button>
            <button @click="hasGift = null" :class="toggleBtnStyling(hasGift == null)" class="btn btn-md px-8">All</button>
          </div>
        </template>
        <template v-else>
          <h2 class="text-5xl font-bold">Recent Wall</h2>
        </template>
      </div>

      <div v-if="messages.length != 0" class="flex flex-wrap items-stretch justify-center">
        <div :key="msg.id" v-for="msg in messages" class="w-1/3 p-2  min-h-16 block">
          <router-link 
            :to="{ name: 'message-page', params: { recipientId: $route.params.recipientId, messageId: msg.id } }" 
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
          <!-- TODO: add infinite scrolling -->
        </div>
      </div>
      <div>
        <!-- TODO: Empty State -->
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
import client from '../client';

dayjs.extend(relativeTime);

export default {
  components: {
    IconReply
  },
  mounted() {
    this.loadMessages();
  },
  data() {
    return {
      links: { first: null, last: null, next: null, previous: null },
      page: 1,
      perPage: 10,
      pageCount: 1,
      messages: [],
      hasGift: false
    }
  },
  watch: {
    hasGift(newVal: boolean, oldVal: boolean) {
      if (newVal === oldVal) return;
      this.loadMessages(newVal);
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
    async loadMessages(hasGift: boolean = false): Promise<void> {
      const recipientId = this.$route.params.recipientId ?? '';
      const resp = await client.get(`/messages/${recipientId}?${hasGift == null ? 'has_gift=2' : hasGift ? 'has_gift=1' : 'has_gift=0'}`);
      const json = await resp.json();

      if (resp.status == 200) {
        this.links = json['links'];
        this.page = json['page'];
        this.perPage = json['per_page'];
        this.pageCount = json['page_count'];
        this.messages = json['data'];
      }

      logEvent(analytics, 'search_messages', { recipient_id: recipientId, status_code: resp.status });
    },
    humanizeTime(date: Date | string): string {
      return dayjs(date).toNow(true);
    }
  }
}
</script>