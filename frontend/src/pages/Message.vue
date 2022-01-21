<template>
  <div class="bg-pink-200 min-h-screen flex">
    <div class="max-w-3xl w-full mx-auto pt-4 flex flex-col space-y-4 self-start">
      <div class="text-center">
        <router-link :to="{ name: 'home-page' }" class="text-6xl py-4 text-white font-bold">Valentine Wall</router-link>
      </div>

      <div v-if="isLoading" class="w-full bg-white rounded-lg p-8">
        <p class="text-center">Loading...</p>
      </div>
      <template v-else-if="!notFound && message">
        <div class="w-full bg-white rounded-lg divide-y-2">
          <div class="p-12">
            <div class="flex flex-col items-center text-center" v-if="hasGift">
              <div class="p-8 bg-white border-gray-200 border rounded-full shadow-md">
                <gift-icon :uid="gift.uid" class="text-4xl" />
              </div>
              <p class="mt-4 text-gray-500 text-xl mb-2">Someone gifted {{ displayName }}</p>
              <p class="text-3xl font-bold">{{ gift.label }}</p>
            </div>
            <p v-else class="text-gray-500 text-xl mb-2">For {{ displayName }}</p>
            <div class="mb-8" :class="{ 'mt-8 bg-amber-100 rounded-lg text-center': hasGift, 'px-8 py-16': revealContent && hasGift, 'mt-12': !hasGift }">
              <button v-if="hasGift && !revealContent" class="w-full p-4 hover:bg-amber-200 rounded-lg" @click="revealContent = true">Reveal note</button>
              <p v-if="revealContent" class="font-bold text-4xl">{{ message.content }}</p>
            </div>
            <p class="text-gray-500" :class="{ 'text-center': hasGift }">
              Posted {{ relativifyDate(message.created_at) }} ({{ formatDate(message.created_at, 'MMMM D, YYYY h:mm A') }})
            </p>
          </div>

          <div class="flex space-x-2 px-8 py-4">
            <template v-if="!$store.state.isAuthLoading && $store.getters.isLoggedIn">
              <button 
                v-if="message.has_replied || $store.state.user.associatedId == message.recipient_id" 
                @click="openReplyModal = true" 
                :disabled="message.has_replied"
                class="flex-1 normal-case btn btn-md border-none space-x-2 bg-white text-gray-900 hover:bg-gray-100">
                <icon-reply :class="[message.has_replied ? 'text-pink-500' : 'text-gray-500']" />
                <span>{{ message.has_replied ? 'Has replied' : 'Reply' }}</span>
              </button>
            </template>

            <button @click="openShareModal = true" class="hover:bg-gray-100 flex-1 normal-case btn btn-md border-none space-x-2 bg-white text-gray-900">
              <icon-share class="text-gray-500" />
              <span>Share</span>
            </button>

            <!-- TODO: add delete button -->
          </div>
        </div>

        <div v-if="message.has_replied && (reply && reply.content)" class="w-full bg-white rounded-lg p-12">
          <p class="text-gray-500 mb-2">{{ message.recipient_id }} replied</p>
          <p class="text-2xl">{{ reply.content }}</p>
        </div>
      </template>
      <template v-else>
        <div class="w-full bg-white p-14 rounded-lg flex flex-col items-center text-center">
          <icon-confused class="text-gray-500 text-9xl mb-2" />
          <h2 class="text-4xl font-bold mb-4">{{ notFound ? 'Message not found.' : 'Something went wrong.' }}</h2>
          <p class="text-xl text-gray-500">{{ notFound ? 'Double-check if your link is correct and try again.' : 'Might be an error on our side. Please try again.' }}</p>
        </div>
      </template>
    </div>
  </div>

  <teleport to="body">
    <share-modal v-model:open="openShareModal" :recipient-id="$route.params.recipientId" :message-id="$route.params.messageId" :permalink="permalink" />
    <reply-message-modal @update:hasReplied="message.has_replied ?? false" v-model:open="openReplyModal" :message="message" />
  </teleport>
</template>

<script lang="ts">
import IconFacebook from '~icons/uil/facebook-f';
import IconTwitter from '~icons/uil/twitter';
import IconLink from '~icons/uil/link';
import IconReply from '~icons/uil/comment-heart';
import IconConfused from '~icons/uil/confused';
import IconShare from '~icons/uil/share-alt';

import ReplyMessageModal from '../components/ReplyMessageModal.vue';
import ShareModal from '../components/ShareModal.vue';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';
import client from '../client';
import { catchAndNotifyError } from '../notify';
import { Gift } from '../store';
import GiftIcon from '../components/GiftIcon.vue';

dayjs.extend(relativeTime);

export default {
  components: {
    IconFacebook, 
    IconTwitter, 
    IconLink, 
    IconReply,
    IconConfused,
    IconShare,
    ReplyMessageModal,
    ShareModal,
    GiftIcon,
  },
  mounted() {
    this.loadMessage();
  },
  data() {
    return {
      isLoading: true,
      message: null as unknown as Record<string, any>,
      reply: null as unknown as Record<string, any>,
      notFound: false,
      openReplyModal: false,
      openShareModal: false,
      revealContent: false
    }
  },
  methods: {
    async loadMessage() {
      try {
        const resp = await client.get(`/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}`);
        const json = await resp.json();
        if (resp.status == 200) {
          this.message = json['message'];
          this.reply = json['reply'];

          if (!this.hasGift) {
            this.revealContent = true;
          }
        } else if (resp.status == 404) {
          this.notFound = true;
        } else {
          logEvent(analytics, 'retrieve_message', { status_code: resp.status });
          throw new Error(json['error_message']);
        }
        logEvent(analytics, 'retrieve_message', { status_code: resp.status });
      } catch(e) {
        catchAndNotifyError(this, e);
      } finally {
        this.isLoading = false;
      }
    },
    relativifyDate(date: Date) {
      return dayjs(date).fromNow();
    },
    formatDate(date: Date, format: string) {
      return dayjs(date).format(format);
    },
  },
  computed: {
    hasGift(): boolean {
      return this.message.gift_id && this.gift;
    },
    gift(): Gift | null {
      if (!this.message) return null;
      return this.$store.state.giftList.find(g => g.id === this.message.gift_id) ?? null;
    },
    displayName(): string {
      if (this.message.recipient_id === this.$store.state.user.associatedId) {
        return "you";
      } else {
        return this.message.recipient_id;
      }
    },
    permalink(): string {
      // return import.meta.env.BASE_URL + this.$route.fullPath;
      return window.location.href;
    }
  },
}
</script>