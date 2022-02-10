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
          <div class="w-full bg-white rounded-lg divide-y-2 shadow-lg">
            <div class="p-12">
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
              <template v-if="!$store.state.isAuthLoading && $store.getters.isLoggedIn">
                <button 
                  v-if="message.has_replied || $store.state.user.associatedId == message.recipient_id" 
                  @click="openReplyModal = true" 
                  :disabled="message.has_replied"
                  class="flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-gray-900 hover:bg-gray-100">
                  <icon-reply :class="[message.has_replied ? 'text-pink-500' : 'text-gray-500']" />
                  <span>{{ message.has_replied ? 'Has replied' : 'Reply' }}</span>
                </button>
              </template>

              <button @click="openShareModal = true" class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-gray-900">
                <icon-share class="text-gray-500" />
                <span>Share</span>
              </button>

              <button v-if="isDeletable" @click="openDeleteModal = true" class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-red-500">
                <icon-trash class="text-red-500" />
                <span>Delete</span>
              </button>

              <!-- TODO: add report link -->
              <button class="hover:bg-gray-100 flex-1 lg:flex-none normal-case btn btn-md border-none space-x-2 bg-white text-gray-900">
                <icon-report class="text-gray-500" />
                <span>Report</span>
              </button>
            </div>
          </div>

          <div v-if="message.has_replied && (reply && reply.content)" class="shadow-lg w-full bg-white rounded-lg p-12">
            <p class="text-gray-500 mb-2">{{ message.recipient_id }} replied</p>
            <p class="text-2xl">{{ reply.content }}</p>
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

    <portal>
      <share-modal 
        v-model:open="openShareModal" 
        :recipient-id="$route.params.recipientId" 
        :message-id="$route.params.messageId" 
        :permalink="permalink" />
      <reply-message-modal 
        @update:hasReplied="handleHasReplied" 
        v-model:open="openReplyModal" 
        :message="message" />
      <modal v-model:open="openDeleteModal">
        <p>Are you sure you want to delete this?</p>
        <template #actions>
          <button class="btn" @click="openDeleteModal = false">Cancel</button>
          <button class="btn btn-error" @click="deleteMessage">Delete</button>
        </template>
      </modal>
    </portal>
  </main>
</template>

<script lang="ts">
import IconFacebook from '~icons/uil/facebook-f';
import IconTwitter from '~icons/uil/twitter';
import IconLink from '~icons/uil/link';
import IconReply from '~icons/uil/comment-heart';
import IconConfused from '~icons/uil/confused';
import IconTrash from '~icons/uil/trash-alt';
import IconShare from '~icons/uil/share-alt';
import GiftIcon from '../components/GiftIcon.vue';

import ReplyMessageModal from '../components/ReplyMessageModal.vue';
import ShareModal from '../components/ShareModal.vue';
import Modal from '../components/Modal.vue';

import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';
import { APIResponse, APIResponseError } from '../client';
import { catchAndNotifyError } from '../notify';
import { Gift } from '../store';
import { WatchStopHandle } from '@vue/runtime-core';
import Portal from '../components/Portal.vue';
import ResponseHandler from '../components/ResponseHandler.vue';
import IconReport from '~icons/uil/exclamation-circle';
import { fromNow, prettifyDateTime } from '../time_utils';

export default {
  components: {
    IconFacebook, 
    IconTwitter, 
    IconLink, 
    IconReply,
    IconConfused,
    IconShare,
    IconTrash,
    IconReport,
    ReplyMessageModal,
    ShareModal,
    GiftIcon,
    Modal,
    Portal,
    ResponseHandler,
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
      openReplyModal: false,
      openDeleteModal: false,
      openShareModal: false,
      revealContent: false,
      authLoadingWatcher: null as unknown as WatchStopHandle
    }
  },
  methods: {
    async deleteMessage() {
      try {
        const { data: json } = await this.$client.delete(`/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}`);
        this.$notify({ type: 'success', text: json['message'] });
        this.$router.replace({ name: 'home-page' });
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
        return import.meta.env.BASE_URL + this.$route.fullPath;
      } else {
        return window.location.href;
      }
    },
    endpoint(): string {
      return `/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}`;
    }
  },
}
</script>