<template>
  <div class="bg-pink-200 min-h-screen flex">
    <div class="max-w-3xl w-full mx-auto pt-4 flex flex-col space-y-4 self-start">
      <div class="text-center">
        <router-link :to="{ name: 'home-page' }" class="text-6xl py-4 text-white font-bold">Valentine Wall</router-link>
      </div>
      
      <template v-if="!notFound && message">
        <div class="w-full bg-white rounded-lg divide-y-2">
          <div class="p-12">
            <p class="text-gray-500 text-xl mb-2">For {{ message.recipient_id }}</p>
            <p class="font-bold text-4xl mb-8">{{ message.content }}</p>
            <p class="text-gray-500">
              Posted {{ relativifyDate(message.created_at) }} ({{ formatDate(message.created_at, 'MMMM D, YYYY h:mm A') }})
            </p>
          </div>

          <div class="flex space-x-2 px-8 py-4">
            <template v-if="!$store.state.isAuthLoading && $store.getters.isLoggedIn">
              <button 
                v-if="message.has_replied" 
                disabled
                class="flex-1 normal-case btn btn-md border-none space-x-2 bg-white text-gray-900 hover:bg-gray-100">
                <icon-reply class="text-pink-500" />
                <span>Has replied</span>
              </button>

              <button 
                v-else-if="$store.state.user.associatedId == message.recipient_id" 
                @click="openReplyModal = true" 
                class="flex-1 normal-case btn btn-md border-none space-x-2 bg-white text-gray-900 hover:bg-gray-100">
                <icon-reply class="text-gray-500" />
                <span>Reply</span>
              </button>
            </template>

            <!-- TODO: rework share button -->
            <button @click="copyURL" :class="[hasLinkCopied ? 'hover:bg-green-100' : 'hover:bg-gray-100']" class="flex-1 normal-case btn btn-md border-none space-x-2 bg-white text-gray-900">
              <icon-link :class="[hasLinkCopied ? 'text-green-600' : 'text-gray-500']" />
              <span>{{ hasLinkCopied ? 'Copied!' : 'Copy URL' }}</span>
            </button>
          </div>
        </div>

        <!-- TODO: Add reply box -->
        <div v-if="message.has_replied && reply" class="w-full bg-white rounded-lg p-12">
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
    <reply-message-modal @update:hasReplied="message.has_replied ?? false" v-model:open="openReplyModal" :message="message" />
  </teleport>
</template>

<script lang="ts">
import IconFacebook from '~icons/uil/facebook-f';
import IconTwitter from '~icons/uil/twitter';
import IconLink from '~icons/uil/link';
import IconReply from '~icons/uil/comment-heart';
import IconConfused from '~icons/uil/confused';
import ReplyMessageModal from '../components/ReplyMessageModal.vue';

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';

dayjs.extend(relativeTime);

export default {
  components: {
    IconFacebook, 
    IconTwitter, 
    IconLink, 
    IconReply,
    IconConfused,
    ReplyMessageModal,
  },
  mounted() {
    this.loadMessage();
  },
  data() {
    return {
      message: null as never as Record<string, any>,
      reply: null as never as Record<string, any>,
      notFound: false,
      openReplyModal: false,
      hasLinkCopied: false
    }
  },
  methods: {
    async loadMessage() {
      const resp = await fetch(import.meta.env.VITE_BACKEND_URL + `/messages/${this.$route.params.recipientId}/${this.$route.params.messageId}`, {
        headers: this.$store.getters.headers
      });

      if (resp.status == 200) {
        const json = await resp.json();
        this.message = json['message'];
        this.reply = json['reply'];
      } else if (resp.status == 404) {
        this.notFound = true;
      }
      logEvent(analytics, 'retrieve_message', { status_code: resp.status });
    },
    relativifyDate(date: Date) {
      return dayjs(date).fromNow();
    },
    formatDate(date: Date, format: string) {
      return dayjs(date).format(format);
    },
    copyURL() {
      this.hasLinkCopied = true;
      navigator.clipboard.writeText(window.location.toString());
      logEvent(analytics, 'share', { method: 'copy-url', item_id: this.message.id });
      setTimeout(() => {
        this.hasLinkCopied = false;
      }, 1500);
    }
  }
}
</script>