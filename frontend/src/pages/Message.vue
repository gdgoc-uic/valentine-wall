<template>
  <main class="min-h-[60vh]">
    <div class="max-w-3xl w-full mx-auto pt-4 flex flex-col space-y-4 self-start">
      <response-handler 
        disappear-on-loading
        :query="query">
        <template #default>
          <div class="w-full bg-white rounded-xl divide-y-2 shadow-lg">
            <div class="p-6 lg:p-12">
              <div class="flex flex-col items-center text-center" v-if="message?.expand.gifts">
                <div class="flex flex-row space-x-2 items-center justify-center">
                  <div :key="gift.uid" v-for="gift in message!.expand.gifts" class="p-8 bg-white border-gray-200 border rounded-full shadow-md">
                    <gift-icon :uid="gift.uid" class="text-4xl" />
                  </div>
                </div>
                <p class="mt-4 text-gray-500 text-xl mb-2">Someone gifted {{ getDisplayName(message!.recipient) }}</p>
                <p class="text-3xl font-bold">{{ displayedGiftLabels }}</p>
              </div>
              <p v-else class="text-gray-500 text-xl mb-2">For {{ getDisplayName(message!.recipient) }}</p>
              <div class="mb-8 mt-8 bg-amber-100 rounded-lg text-center" :class="{ 'px-8 py-16': revealContent || !hasGifts }">
                <button v-if="hasGifts && !revealContent" class="w-full p-4 hover:bg-amber-200 rounded-lg" @click="revealContent = true">Reveal note</button>
                <div v-if="!hasGifts || revealContent" class="font-bold text-4xl">
                  <p :key="'content_' + ci" v-for="(c, ci) in message!.content.split('\r\n')">{{ c }}</p>
                </div>
              </div>
              <p class="text-gray-500" :class="{ 'text-center': hasGifts }">
                Posted {{ fromNow(message!.created) }} ({{ prettifyDateTime(message!.created) }})
              </p>
            </div>

            <div class="flex space-x-2 px-8 py-4">
              <share-dialog 
                :image-url="imageUrl"
                :image-file-name="`${message!.recipient}-${message!.id}.png`"
                :permalink="permalink"
                :title="'test'"
                :hashtags="['UICValentineWall']"
                @success="onShareSuccess"
                >
                <template #default="{ openDialog }">
                  <button 
                    @click="openDialog" 
                    class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-gray-900">
                    <icon-share class="text-gray-500" />
                    <span>Share</span>
                  </button>
                </template>
              </share-dialog>

              <delete-dialog @confirm="onMessageDeletion">
                <template #default="{ openDialog }">
                  <button
                    v-if="message?.user === authState.user?.details"
                    @click="openDialog"
                    class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-red-500">
                    <icon-trash class="text-red-500" />
                    <span>Delete</span>
                  </button>
                </template>
              </delete-dialog>

              <report-dialog 
                :link="permalink"
                :email="authState.user?.email ?? ''">
                <template #default="{ openDialog }">
                  <button 
                    @click="openDialog" 
                    class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-gray-900">
                    <icon-report class="text-gray-500" />
                    <span>Report</span>
                  </button>
                </template>
              </report-dialog>
            </div>
          </div>
          
          <reply-thread />
        </template>

        <template #error="{ error }">
          <div class="w-full bg-white p-14 rounded-lg flex flex-col items-center text-center">
            <icon-confused class="text-gray-500 text-9xl mb-2" />
            <h2 class="text-4xl font-bold mb-4">{{ isClientError(error) && error.status == 404 ? 'Message not found.' : 'Something went wrong.' }}</h2>
            <p class="text-xl text-gray-500">{{ isClientError(error) && error.status == 404 ? 'Double-check if your link is correct and try again.' : 'Might be an error on our side. Please try again.' }}</p>
          </div>
        </template>
      </response-handler>
    </div>
  </main>
</template>

<script lang="ts" setup>
import IconConfused from '~icons/uil/confused';
import IconTrash from '~icons/uil/trash-alt';
import IconShare from '~icons/uil/share-alt';
import GiftIcon from '../components/GiftIcon.vue';

import ReplyThread from '../components/ReplyThread.vue';
import ShareDialog from '../components/ShareDialog.vue';

import { logEvent } from '../analytics';
import { notify } from '../notify';
import ResponseHandler from '../components/ResponseHandler2.vue';
import IconReport from '~icons/uil/exclamation-circle';
import { fromNow, prettifyDateTime } from '../time_utils';
import DeleteDialog from '../components/DeleteDialog.vue';
import ReportDialog from '../components/ReportDialog.vue';
import { ref, computed, provide } from 'vue';
import { pb } from '../client';
import { useMutation, useQuery } from '@tanstack/vue-query';
import { useRoute, useRouter } from 'vue-router';
import { ClientResponseError } from 'pocketbase';
import { useAuth } from '../store_new';
import { Gift } from '../types';

const { state: authState, methods } = useAuth();
const router = useRouter();
const route = useRoute();
const revealContent = ref(false);

function isClientError(err: unknown): err is ClientResponseError {
  return err instanceof ClientResponseError;
}

function onShareSuccess(provider: string) {
  // share success
  if (provider === 'clipboard') {
    logEvent('share', { method: 'copy-url', item_id: message.value!.id });
  } else if (authState.isLoggedIn) {
    methods.reward(300, 'Social share');
  }
}

async function onMessageDeletion(confirmed: boolean) {
  if (confirmed) {
    await deleteMessage();
  }
}

function generateDisplayGiftLabelString(g: Gift, i: number, arr: Gift[]) {
  let displayStr = g.label;
  if (arr.length > 1 && i === arr.length - 1) {
    displayStr = 'and ' + displayStr;
  }
  return displayStr;
}

function getDisplayName(recipientId: string): string {
  if (recipientId === authState.user?.expand.details?.student_id) {
    return "you";
  } else {
    return recipientId;
  }
}

const { mutateAsync: deleteMessage } = useMutation(
  () => pb.collection('messages').delete(message.value!.id),
  {
    onSuccess() {
      notify({ type: 'success', text: 'Message was deleted successfully.' });
      router.replace({ name: 'message-wall-page', params: { recipientId: recipient.value } });
    }
  }
);

const hasGifts = computed(() => message.value?.gifts.length !== 0);
const recipient = computed(() => route.params.recipientId ?? ''); 
const messageId = computed(() => {
  const msgId = route.params.messageId;
  if (Array.isArray(msgId)) {
    return msgId[0];
  }
  return msgId;
});

const query = useQuery(
  ['message', recipient, messageId],
  () => pb.collection('messages').getOne(messageId.value, {
    expand: 'gifts,message_replies(message)'
  }),
  {
    refetchOnWindowFocus: false,
    retry: 0
  }
);

const message = query.data;
const imageUrl = computed(() => pb.buildUrl(`/messages/${messageId.value}/image`));
const permalink = computed(() => import.meta.env.VITE_FRONTEND_URL + route.fullPath);
const displayedGiftLabels = computed(() => {
  if (!message.value || !message.value.expand.gifts) return '';
  return message.value.expand.gifts.map(generateDisplayGiftLabelString).join(
    message.value.expand.gifts.length > 2 ? ', ' : message.value.expand.gifts.length == 2 ? ' ' : '');
});

provide('message', message);
</script>