<template>
  <main class="min-h-[60vh] flex flex-col px-4">
    <div class="text-center flex flex-col items-center justify-center mb-4">
      <div class="pt-8 pb-6">
        <div v-if="recipient">
          <p class="text-2xl mb-2 text-gray-500">Messages for </p>
          <div class="indicator">
            <div v-if="totalCount != 0" class="indicator-item badge">{{ totalCount }}</div>
            <h2 class="text-5xl lg:text-6xl font-bold text-rose-600">{{ recipient }}</h2>
          </div>
        </div>
        <h2 v-else class="text-5xl lg:text-6xl font-bold text-rose-600">Recent Wall</h2>
      </div>
    </div>
    <client-only>
      <div 
        v-if="authState.isLoggedIn && recipient" 
        class="max-w-3xl bg-white p-8 mb-8 mx-auto w-full space-y-2 rounded-2xl shadow-md h-full">
        <h2 class="text-3xl font-bold text-center">Send a message</h2>

        <send-message-form :existing-recipient="recipient" />
      </div>
      
      <div class="max-w-7xl mx-auto w-full bg-white rounded-xl shadow-lg flex justify-center">
        <div 
          v-if="recipient && authState.isLoggedIn && authState.user!.expand.details.student_id === recipient" 
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
        v-if="!authState.isLoggedIn && recipient" 
        class="max-w-7xl space-y-8 lg:space-y-0 mx-auto mt-3 p-8 w-full border bg-[#ffeded] border-rose-500 rounded-xl shadow-lg flex flex-col lg:flex-row text-center lg:text-left items-center justify-between">
        <div class="flex flex-col">
          <h3 class="text-2xl lg:text-3xl mb-1 text-rose-400 font-semibold">Are you {{ recipient }}?</h3>
          <p class="text-lg lg:text-xl">Join now and get access to exclusive features!</p>
        </div>

        <login-button />
      </section>
    </client-only>

    <response-handler :query="query">
      <template #default>
        <div class="max-w-7xl w-full mx-auto flex flex-col">
          <message-tiles :messages="currentMessages" replace />
          <pagination-load-more-button 
            :should-go-next="hasNextPage" 
            @click="fetchNextPage" />
        </div>
      </template>

      <template #error="{ error }">
        <div 
          class="text-center w-full py-16 h-full flex flex-col items-center">
          <template 
            v-if="isEmptyError(error)">
            <p class="text-4xl">Nothing to see here!</p>
          </template> 
          <template v-else-if="hasError(error)">
            <p class="text-4xl">Something went wrong.</p>
          </template>
        </div>
      </template>
    </response-handler>
  </main>
</template>

<script lang="ts" setup>
import SendMessageForm from '../components/SendMessageForm.vue';
import ClientOnly from '../components/ClientOnly.vue';
import ResponseHandler from '../components/ResponseHandler2.vue';
import PaginationLoadMoreButton from '../components/PaginationLoadMoreButton.vue';
import LoginButton from '../components/LoginButton.vue';
import MessageTiles from '../components/MessageTiles.vue';
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { onBeforeRouteUpdate, useRoute } from 'vue-router';
import { pb } from '../client';
import { useInfiniteQuery } from '@tanstack/vue-query';
import { ClientResponseError, Record as PbRecord, UnsubscribeFunc } from 'pocketbase';
import { useAuth } from '../store_new';

function isEmptyError(error: unknown) {
  if (error instanceof ClientResponseError && error.status === 404) {
    return true;
  } else if (error instanceof Error && error.message === 'No messages found.') {
    return true;
  } 
  return false;
}

function hasError(error: unknown) {
  return error instanceof Error && error.message;
}

const recipient = computed(() => route.params.recipientId ?? '');

async function loadStats(): Promise<void> {
  try {
    if (!recipient.value) return;

    stats.messages_count = 0;
    stats.gift_messages_count = 0;

    const msgs = pb.collection('messages');
    const [giftsResp, wGiftsResp] = await Promise.all([
      msgs.getList(1, 1, { filter: `recipient="${recipient.value}"` }),
      msgs.getList(1, 1, { filter: `recipient="${recipient.value}" && gifts != "[]"` })
    ]);

    stats.messages_count = giftsResp.totalItems;
    stats.gift_messages_count = wGiftsResp.totalItems;
  } catch(e) {
    console.error(e);
  }
}

const route = useRoute();
const { state: authState } = useAuth();
const hasGift = ref<boolean | null>(false);

if (route.params.receipientId && authState.isLoggedIn) {
  hasGift.value = null;
}

const stats = reactive({
  messages_count: 0,
  gift_messages_count: 0
});

const totalCount = computed(() => {
  if (import.meta.env.SSR) {
    return 0;
  } else if (authState.isLoggedIn
    && authState.user!.expand.details?.student_id === recipient.value) {
    return stats.gift_messages_count + stats.messages_count;
  } else {
    return stats.messages_count;
  }
});

onBeforeRouteUpdate((to, from) => {
  if (to.params === from.params) return;
  loadStats();
});

const { hasNextPage, fetchNextPage, ...query } = useInfiniteQuery(
  ['wall', recipient, hasGift],
  async ({ pageParam = 1 }) => {
    let hasGiftsFilter = '';
    if (hasGift.value !== null) {
      hasGiftsFilter = '&& ' + (hasGift.value ? `gifts != "[]"` : `gifts = "[]"`);
    }

    const resp = await pb.collection('messages').getList(pageParam, 10, { 
      sort: '-created',
      filter: `(recipient = "${recipient.value}" ${hasGiftsFilter})`
    });

    if (resp.items.length === 0) {
      throw new Error('No messages found.');
    }

    return resp;
  },
  {
    keepPreviousData: true,
    retry: 1,
    refetchOnWindowFocus: () => false,
    getNextPageParam: (result) => 
      result.page + 1 <= result.totalPages ? result.page + 1 : undefined
  }
);

const rawCurrentMessages = reactive<PbRecord[]>([]);

const currentMessages = computed(() => {
  if (!query.data.value) {
    return [];
  }

  rawCurrentMessages.splice(0, rawCurrentMessages.length - 1);
  rawCurrentMessages.push(...query.data.value.pages.map(p => p.items).flat());
  return rawCurrentMessages;
});

const unsubscribeFunc = ref<UnsubscribeFunc | null>(null);

onMounted(() => {
  if (!import.meta.env.SSR) {
    if (stats.messages_count == 0) {
      loadStats();
    }

    pb.collection('messages').subscribe('*', (data) => {
      if (data.action === 'create' && data.record.recipient === recipient.value) {
        if (data.record.gifts.length === 0) {
          stats.messages_count++;
        } else {
          stats.gift_messages_count++;
        }
       
        if (hasGift.value === true && data.record.gifts.length === 0) {
          return;
        }
        
        rawCurrentMessages.unshift(data.record);
      }
    }).then((unsub) => {
      unsubscribeFunc.value = unsub;
    });
  }
});

onUnmounted(() => {
  unsubscribeFunc.value?.();
})
</script>