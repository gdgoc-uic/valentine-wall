<template>
  <main class="min-h-[60vh]">
    <div class="max-w-3xl w-full mx-auto pt-4 flex flex-col space-y-4 self-start">
      <response-handler 
        ref="messageResponse"
        @success="handleResponse"
        @error="handleResponseError"
        disappear-on-loading
        :endpoint="endpoint">
        <template #default="{ response: { data: { message, reply } } }">
          <div class="w-full bg-white rounded-xl divide-y-2 shadow-lg">
            <div class="p-6 lg:p-12">
              <div class="flex flex-col items-center text-center" v-if="hasGifts">
                <div class="flex flex-row space-x-2 items-center justify-center">
                  <div :key="gift.uid" v-for="gift in gifts" class="p-8 bg-white border-gray-200 border rounded-full shadow-md">
                    <gift-icon :uid="gift.uid" class="text-4xl" />
                  </div>
                </div>
                <p class="mt-4 text-gray-500 text-xl mb-2">Someone gifted {{ getDisplayName(message.recipient_id) }}</p>
                <p class="text-3xl font-bold">{{ displayedGiftLabels }}</p>
              </div>
              <p v-else class="text-gray-500 text-xl mb-2">For {{ getDisplayName(message.recipient_id) }}</p>
              <div class="mb-8" :class="{ 'mt-8 bg-amber-100 rounded-lg text-center': hasGifts, 'px-8 py-16': revealContent && hasGifts, 'mt-2': !hasGifts }">
                <button v-if="hasGifts && !revealContent" class="w-full p-4 hover:bg-amber-200 rounded-lg" @click="revealContent = true">Reveal note</button>
                <p v-if="revealContent" class="font-bold text-4xl">{{ message.content }}</p>
              </div>
              <p class="text-gray-500" :class="{ 'text-center': hasGifts }">
                Posted {{ relativifyDate(message.created_at) }} ({{ prettifyDate(message.created_at) }})
              </p>
            </div>

            <div class="flex space-x-2 px-8 py-4">
              <share-dialog 
                :image-url="imageUrl"
                :image-file-name="`${message.recipient_id}-${message.id}.png`"
                :permalink="permalink"
                :title="$route.meta.pageTitle($route)"
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
                :email="$store.getters.isLoggedIn ? $store.state.user.email : null" 
                :message-id="message.id">
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

          <div v-if="message.has_replied && (reply && reply.content)" class="shadow-lg w-full bg-white rounded-lg p-6 md:p-12">
            <p class="text-gray-500 mb-2">{{ message.recipient_id }} replied</p>
            <p class="text-2xl">{{ reply.content }}</p>

            <div class="flex items-end w-full">
              <delete-dialog
                v-if="$store.state.user.associatedId == message.recipient_id"
                @confirm="onMessageReplyDeletion" v-slot="{ openDialog }">
                <button
                  @click="openDialog" class="btn btn-error space-x-2 mt-8 self-end flex items-center">
                  <icon-trash />
                  <span>Delete</span>
                </button>
              </delete-dialog>
            </div>
          </div>
          <div v-else-if="!message.has_replied || !$store.getters.isLoggedIn" class="p-6 lg:p-8 bg-white rounded-xl shadow-lg">
            <div v-if="$store.state.user.associatedId == message.recipient_id"
              class="flex space-x-2 items-center text-2xl">
              <icon-reply class="text-pink-500 mb-4" />
              <h2 class="font-bold mb-4">Your reply</h2>
            </div>
            <reply-message-box @update:hasReplied="handleHasReplied" :message="message" />
          </div>
        </template>

        <template #error="{ error: { rawResponse: resp } }">
          <div class="w-full bg-white p-14 rounded-lg flex flex-col items-center text-center">
            <icon-confused class="text-gray-500 text-9xl mb-2" />
            <h2 class="text-4xl font-bold mb-4">{{ resp.status == 404 ? 'Message not found.' : 'Something went wrong.' }}</h2>
            <p class="text-xl text-gray-500">{{ resp.status == 404 ? 'Double-check if your link is correct and try again.' : 'Might be an error on our side. Please try again.' }}</p>
          </div>
        </template>
      </response-handler>
    </div>
  </main>
</template>

<script lang="ts">
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
import { analytics } from '../firebase';
import { APIResponse, APIResponseError, expandAPIEndpoint } from '../client';
import { catchAndNotifyError } from '../notify';
import { Gift } from '../store';
import { WatchStopHandle } from '@vue/runtime-core';
import Portal from '../components/Portal.vue';
import ResponseHandler from '../components/ResponseHandler.vue';
import IconReport from '~icons/uil/exclamation-circle';
import { fromNow, prettifyDateTime } from '../time_utils';
import DeleteDialog from '../components/DeleteDialog.vue';
import ReportDialog from '../components/ReportDialog.vue';

export default {
  components: {
    IconFacebook, 
    IconTwitter, 
    IconLink, 
    IconReply,
    IconConfused,
    IconShare,
    IconSend,
    IconTrash,
    IconReport,
    ReplyMessageBox,
    ShareDialog,
    GiftIcon,
    Modal,
    Portal,
    ResponseHandler,
    DeleteDialog,
    ReportDialog,
  },
  mounted() {
    if (this.$route.query.from) {
      logEvent(analytics!, "traffic_source", { from: this.$route.query.from });
    }
  },
  data() {
    return {
      isDeletable: false,
      message: null as unknown as Record<string, any>,
      reply: null as unknown as Record<string, any>,
      openReportModal: false,
      revealContent: false,
      reportDialogKey: 1,
      authLoadingWatcher: null as unknown as WatchStopHandle
    }
  },
  methods: {
    onReportSuccess() {
      this.reportDialogKey++;
    },
    onShareSuccess(provider: string) {
      // share success
      if (provider === 'clipboard') {
        logEvent(analytics!, 'share', { method: 'copy-url', item_id: this.message.id });
      }
    },
    onShareFailed(err: unknown) {
      console.error(err);
    },
    async onMessageDeletion(confirmed: boolean, closeDialog: Function) {
      if (confirmed) {
        await this.deleteMessage();
      }
      closeDialog();
    },
    async onMessageReplyDeletion(confirmed: boolean, closeDialog: Function) {
      if (confirmed) {
        await this.deleteReply();
      }
      closeDialog();
    },
    async deleteMessage() {
      try {
        const { data: json } = await this.$client.delete(`/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}`);
        this.$notify({ type: 'success', text: json['message'] });
        this.$router.replace({ name: 'home-page' });
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    },
    async deleteReply() {
      try {
        const { data: json } = await this.$client.delete(`/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}/reply`);
        this.$notify({ type: 'success', text: json['message'] });
        this.message.has_replied = false;
        this.$router.go(0);
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    },
    handleHasReplied(hasReplied: boolean) {
      this.message.has_replied = hasReplied;
      if (hasReplied) {
        this.$router.go(0);
      }
    },
    handleResponse(r: APIResponse) {
      const { rawResponse: resp, data: json } = r;
      this.isDeletable = json['is_deletable'] ?? false;
      this.message = json['message'];
      if (!this.hasGifts) {
        this.revealContent = true;
      }
      logEvent(analytics!, 'retrieve_message', { status_code: resp.status });
    },
    handleResponseError(e: unknown) {
      if (e instanceof APIResponseError) {
        logEvent(analytics!, 'retrieve_message', { status_code: e.rawResponse.status });
      }
    },
    relativifyDate(date: Date) {
      return fromNow(date);
    },
    prettifyDate(date: Date) {
      return prettifyDateTime(date);
    },
    generateDisplayGiftLabelString(g: Gift, i: number, arr: Gift[]) {
      let displayStr = g.label;
      if (arr.length > 1 && i === arr.length - 1) {
        displayStr = 'and ' + displayStr;
      }
      return displayStr;
    },
    getDisplayName(recipientId: string): string {
      if (recipientId === this.$store.state.user.associatedId) {
        return "you";
      } else {
        return recipientId;
      }
    },
  },
  computed: {
    hasGifts(): boolean {
      return this.message.gift_ids && this.gifts?.length;
    },
    gifts(): Gift[] | null {
      if (!this.message) return null;
      return this.$store.state.giftList.filter(g => this.message.gift_ids?.includes(g.id) ?? false) ?? null;
    },
    displayedGiftLabels(): string {
      if (!this.hasGifts || !this.gifts) return '';
      return this.gifts.map(this.generateDisplayGiftLabelString).join(this.gifts.length > 2 ? ', ' : this.gifts.length == 2 ? ' ' : '');
    },
    permalink(): string {
      if (import.meta.env.SSR) {
        return import.meta.env.VITE_FRONTEND_URL + this.$route.fullPath;
      } else {
        return window.location.href;
      }
    },
    endpoint(): string {
      return `/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}`;
    },
    imageUrl(): string {
      return expandAPIEndpoint(this.endpoint + '?image');
    },
  },
}
</script>