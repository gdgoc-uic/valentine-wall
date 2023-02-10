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
              <div class="mb-8" :class="{ 'mt-8 bg-amber-100 rounded-lg text-center': hasGifts, 'px-8 py-16': revealContent && hasGifts, 'mt-2': !hasGifts }">
                <button v-if="hasGifts && !revealContent" class="w-full p-4 hover:bg-amber-200 rounded-lg" @click="revealContent = true">Reveal note</button>
                <p v-if="revealContent" class="font-bold text-4xl">{{ message!.content }}</p>
              </div>
              <p class="text-gray-500" :class="{ 'text-center': hasGifts }">
                Posted {{ relativifyDate(message!.created_at) }} ({{ prettifyDate(message!.created_at) }})
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
                    v-if="isDeletable"
                    @click="openDialog"
                    class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-red-500">
                    <icon-trash class="text-red-500" />
                    <span>Delete</span>
                  </button>
                </template>
              </delete-dialog>

              <report-dialog 
                :key="reportDialogKey"
                @success="onReportSuccess"
                :email="isLoggedIn ? user!.email : undefined" 
                :message-id="message!.id">
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

          <div v-if="message!.expand.message_replies" class="shadow-lg w-full bg-white rounded-lg p-6 md:p-12">
            <p class="text-gray-500 mb-2">{{ message!.recipient }} replied</p>
            <!-- TODO: multi-reply support -->
            <p class="text-2xl">{{ message!.expand.message_replies.toString() }}</p>

            <div class="flex items-end w-full">
              <delete-dialog
                v-if="user!.expand.details?.student_id == message!.recipient"
                @confirm="onMessageReplyDeletion" v-slot="{ openDialog }">
                <button
                  @click="openDialog" class="btn btn-error space-x-2 mt-8 self-end flex items-center">
                  <icon-trash />
                  <span>Delete</span>
                </button>
              </delete-dialog>
            </div>
          </div>
          <div v-else-if="!message!.expand.message_replies || !isLoggedIn" class="p-6 lg:p-8 bg-white rounded-xl shadow-lg">
            <div v-if="user!.expand.details?.student_id == message!.recipient"
              class="flex space-x-2 items-center text-2xl">
              <icon-reply class="text-pink-500 mb-4" />
              <h2 class="font-bold mb-4">Your reply</h2>
            </div>
            <reply-message-box @update:hasReplied="handleHasReplied" :message="message" />
          </div>
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
import IconFacebook from '~icons/uil/facebook-f';
import IconTwitter from '~icons/uil/twitter';
import IconSend from '~icons/uil/message';
import IconLink from '~icons/uil/link';
import IconReply from '~icons/uil/comment-heart';
import IconConfused from '~icons/uil/confused';
import IconTrash from '~icons/uil/trash-alt';
import IconShare from '~icons/uil/share-alt';
import GiftIcon from '../components/GiftIcon.vue';

import ReplyMessageBox from '../components/ReplyMessageBox.vue';
import ShareDialog from '../components/ShareDialog.vue';
import Modal from '../components/Modal.vue';

import { logEvent } from '@firebase/analytics';
// import { analytics } from '../firebase';
import { catchAndNotifyError } from '../notify';
import { WatchStopHandle } from '@vue/runtime-core';
import Portal from '../components/Portal.vue';
import ResponseHandler from '../components/ResponseHandler2.vue';
import IconReport from '~icons/uil/exclamation-circle';
import { fromNow, prettifyDateTime } from '../time_utils';
import DeleteDialog from '../components/DeleteDialog.vue';
import ReportDialog from '../components/ReportDialog.vue';
import { ref, computed } from 'vue';
import { pb } from '../client';
import { useMutation, useQuery } from '@tanstack/vue-query';
import { useRoute, useRouter } from 'vue-router';
import { ClientResponseError } from 'pocketbase';
import { useAuth, useStore } from '../store_new';
import { Gift } from '../types';

const { state: {isLoggedIn, user} } = useAuth();
const router = useRouter();
const route = useRoute();
const isDeletable = ref(false);
const openReportModal = ref(false);
const revealContent = ref(false);
const reportDialogKey = ref(1);
const authLoadingWatcher = ref<WatchStopHandle>(null!);

function isClientError(err: unknown): err is ClientResponseError {
  return err instanceof ClientResponseError;
}

function onReportSuccess() {
  reportDialogKey.value++;
}

function onShareSuccess(provider: string) {
  // share success
  if (provider === 'clipboard') {
    // logEvent(analytics!, 'share', { method: 'copy-url', item_id: this.message.id });
  }
}

function onShareFailed(err: unknown) {
  console.error(err);
}

async function onMessageDeletion(confirmed: boolean, closeDialog: Function) {
  if (confirmed) {
    await deleteMessage();
  }
  closeDialog();
}

async function onMessageReplyDeletion(confirmed: boolean, closeDialog: Function) {
  if (confirmed) {
    // TODO:
    // await deleteReply();
  }
  closeDialog();
}

function handleHasReplied(hasReplied: boolean) {
  // this.message.has_replied = hasReplied;
  if (hasReplied) {
    router.go(0);
  }
}

function relativifyDate(date: Date) {
  return fromNow(date);
}

function prettifyDate(date: Date) {
  return prettifyDateTime(date);
}

function generateDisplayGiftLabelString(g: Gift, i: number, arr: Gift[]) {
  let displayStr = g.label;
  if (arr.length > 1 && i === arr.length - 1) {
    displayStr = 'and ' + displayStr;
  }
  return displayStr;
}

function getDisplayName(recipientId: string): string {
  if (recipientId === user!.expand.details?.student_id) {
    return "you";
  } else {
    return recipientId;
  }
}

// TODO(backend): add safe guards when deleting message
const { mutateAsync: deleteMessage } = useMutation(
  () => pb.collection('messages').delete(messageId.value.toString()),
  {
    onSuccess() {
       // this.$notify({ type: 'success', text: json['message'] });
      router.replace({ name: 'home-page' });
    },
    onError() {
      // catchAndNotifyError(this, e);
    }
  }
);

// TODO(backend): add safe guards when deleting message reply
const { mutateAsync: deleteReply } = useMutation(
  (id: string) => pb.collection('message_replies').delete(id),
  {
    onSuccess() {
      // this.$notify({ type: 'success', text: json['message'] });
      // this.message.has_replied = false;
      router.go(0);
    },
    onError() {
      // catchAndNotifyError(this, e);
    }
  }
);

// const hasGifts = computed(() => )
const hasGifts = computed(() => message.value?.expand.gifts?.length !== 0);
const recipient = computed(() => route.params.recipientId ?? ''); 
const messageId = computed(() => route.params.messageId ?? '');

const query = useQuery(
  ['message', recipient, messageId],
  () => pb.collection('messages').getFirstListItem(
    `id = "${messageId.value}" && recipient = "${recipient.value}"`,
    {
      expand: 'gifts,message_replies(message)'
    }
  ),
  {
    onSuccess(d) {
      console.log(d);
    }
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
</script>