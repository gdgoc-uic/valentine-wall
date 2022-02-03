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
        v-if="!$store.state.loggedIn && $route.params.recipientId" 
        class="max-w-7xl space-y-8 lg:space-y-0 mx-auto mt-3 p-8 w-full border bg-[#ffeded] border-rose-500 rounded-xl shadow-lg flex flex-col lg:flex-row text-center lg:text-left items-center justify-between">
        <div class="flex flex-col">
          <h3 class="text-2xl lg:text-3xl mb-1 text-rose-400 font-semibold">Are you {{ $route.params.recipientId }}?</h3>
          <p class="text-lg lg:text-xl">Join now and get access to exclusive features!</p>
        </div>

        <login-button />
      </section>
    </client-only>

    <paginated-response-handler :origin-endpoint="endpoint" :process-fn="processResults" :fail-fn="checkMessagesLength">
      <template #default="{ data: messages, links, goto }">
        <div class="max-w-7xl w-full mx-auto flex flex-col">
          <masonry class="message-results -mx-2">
            <!-- TODO: make card widths and --colors-- different -->
            <div :key="msg.id" v-for="msg in messages" class="w-1/2 md:w-1/3 lg:w-1/4 message-paper-wrapper">
              <router-link
                :to="{ name: 'message-page', params: { recipientId: msg.recipient_id, messageId: msg.id } }"
                class="message-paper"
                :class="[`paper-variant-${msg.paperColor}`]">
                  <div class="px-6 pt-6">
                    <p>{{ msg.content }}</p>
                  </div>
                  <div class="message-meta-info">
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
import LoginButton from '../components/LoginButton.vue';

dayjs.extend(relativeTime);

export default defineComponent({
  components: {
    IconReply,
    ClientOnly,
    Masonry,
    PaginatedResponseHandler,
    PaginationLoadMoreButton,
    LoginButton
  },
  created() {
    this.endpoint = this.getMessagesEndpoint({ hasGift: this.hasGift });
  },
  mounted() {
    if (this.stats.messages_count == 0) {
      this.loadData();
    }
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
    hasGift(newVal, oldVal) {
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
    processResults(data: any[]): any[] {
      const availablePaperColorId = [1,2,3,4];
      const quo = data.length / availablePaperColorId.length;
      let paperColorIds: number[] = [];
      let j = 0;
      let times = Math.floor(quo);
      if (times < 1) times = Math.ceil(quo);
      for (let i = 0; i < data.length; i++) {
        paperColorIds.push(availablePaperColorId[j]);
        if (i != 0 && i % times == 0) {
          j++;
        }
      }

      // add a check to avoid repetitions
      for (let i = 0; i < times; i++) {
        paperColorIds = this.shuffle(paperColorIds);
      }

      return data.map((d, i) => {
        const paperColor = paperColorIds[i];
        return {
          ...d, 
          paperColor
        }
      });
    },
    shuffle(array: any[]) {
        var i = array.length,
            j = 0,
            temp;
        while (i--) {
            j = Math.floor(Math.random() * (i+1));
            // swap randomly chosen element with current element
            temp = array[i];
            array[i] = array[j];
            array[j] = temp;
        }
        return array;
    },
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
      return `/messages/${recipientId}?order=created_at,desc&limit=10&${hasGift == null ? 'has_gift=2' : hasGift ? 'has_gift=1' : 'has_gift=0'}`
    },
    humanizeTime(date: Date | string): string {
      return dayjs(date).toNow(true);
    }
  }
})
</script>

<style lang="postcss">
.message-results > .message-paper-wrapper {
  @apply p-2 min-h-16;
}

.message-paper-wrapper > .message-paper {
  @apply rounded-lg flex flex-col justify-between shadow-lg h-64 hover:scale-110 transition-transform;
  background-image: linear-gradient(to bottom,rgb(254, 243, 199) 21px,#00b0d7 1px); 
  background-size: 100% 22px;
}

.message-paper .message-meta-info {
  @apply  px-6 pb-6 flex justify-between items-center mt-4 rounded-b-lg;
}

.message-paper.paper-variant-1 {
  background-image: linear-gradient(to bottom,rgb(254, 243, 199) 21px,#00b0d7 1px); 
}

.message-paper.paper-variant-1 .message-meta-info {
  background-color: rgb(254, 243, 199);
}

.message-paper.paper-variant-2 {
  background-image: linear-gradient(to bottom,rgb(152, 221, 255) 21px,#213381 1px); 
}

.message-paper.paper-variant-2 .message-meta-info {
  background-color: rgb(152, 221, 255);
}

.message-paper.paper-variant-3 {
  background-image: linear-gradient(to bottom,rgb(155, 255, 183) 21px,#00b0d7 1px); 
}

.message-paper.paper-variant-3 .message-meta-info {
  background-color: rgb(155, 255, 183);
}

.message-paper.paper-variant-4 {
  background-image: linear-gradient(to bottom,rgb(255, 194, 175) 21px,#213381 1px); 
}

.message-paper.paper-variant-4 .message-meta-info {
  background-color: rgb(255, 194, 175);
}
</style>