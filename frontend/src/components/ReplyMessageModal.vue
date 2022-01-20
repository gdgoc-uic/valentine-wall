<template>
  <modal title="Reply" :open="open" @update:open="$emit('update:open', $event)">
    <div v-if="!$store.getters.isLoggedIn || $store.state.user.associatedId != message.recipient_id" class="flex flex-col justify-center text-center items-center py-8">
      <icon-reply-lock class="text-gray-500 text-9xl mb-4" />
      <h3 class="text-2xl font-bold">Reply is locked.</h3>
      <p class="text-gray-500 text-xl">Recipient can only reply to this message.</p>
    </div>

    <div v-else-if="$store.getters.isLoggedIn && message.has_replied" class="flex flex-col justify-center text-center items-center py-8">
      <icon-reply class="text-pink-500 text-9xl mb-4" />
      <h3 class="text-2xl font-bold">Already replied!</h3>
    </div>

    <div v-else-if="$store.getters.isLoggedIn && !$store.getters.hasConnections" class="flex flex-col justify-center text-center items-center py-8">
      <icon-annoyed class="text-gray-500 text-9xl mb-4" />
      <h3 class="text-2xl font-bold">One more step!</h3>
      <p class="text-gray-500 text-md">Connect your Twitter account in order to reply to this message.</p>
    
      <button class="mt-8 w-full space-x-2 btn bg-twitter-500 hover:bg-twitter-600 rounded-lg border-none hover:border-none" @click="openTwitterLogin">
        <icon-twitter />
        <span>Login to Twitter</span>
      </button>

      <button class="btn btn-ghost btn-sm mt-2 w-full font-normal text-gray-500" @click="skipLogin">Skip</button>
    </div>

    <div v-else>
      <div class="form-control">
        <textarea 
          :disabled="isSending"
          v-model="content"
          class="textarea textarea-bordered h-48" 
          placeholder="Reply something..."></textarea>
      </div>

      <div class="flex items-center space-x-2 mt-4">
        <content-counter ref="counter" :content="content" class="mr-auto" />
        <button @click="submitReply" class="space-x-2 btn bg-rose-500 hover:bg-rose-600 border-none hover:border-none" :disabled="!shouldSend || isSending">
          <icon-send />
          <span>Send</span>
        </button>
        <button @click="$emit('update:open', false)" class="btn" :disabled="isSending">Cancel</button>
      </div>
    </div>
  </modal>
</template>

<script lang="ts">
import Modal from './Modal.vue';
import IconReplyLock from '~icons/uil/comment-lock';
import IconUserLocation from '~icons/uil/user-location';
import IconSend from '~icons/uil/message';
import IconFacebook from '~icons/uil/facebook-f';
import IconTwitter from '~icons/uil/twitter';
import IconAnnoyed from '~icons/uil/annoyed';
import IconReply from '~icons/uil/comment-heart';
import ContentCounter from '../components/ContentCounter.vue';
import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';
import client from '../client';

export default {
  components: { 
    Modal, 
    ContentCounter,
    IconReplyLock,
    IconUserLocation,
    IconSend,
    IconFacebook,
    IconTwitter,
    IconReply,
    IconAnnoyed
  },
  emits: ['update:open', 'update:hasReplied'],
  props: {
    open: {
      type: Boolean,
    },
    message: {
      type: Object,
      default: {}
    }
  },
  data() {
    return {
      content: '',
      shouldSend: false,
      isSending: false
    }
  },
  watch: {
    content(newV: string, oldV: string) {
      const counter = (this.$refs.counter as typeof ContentCounter);
      if (!counter || typeof counter == 'undefined') {
        this.shouldSend = false;
      } else {
        this.shouldSend =  counter.shouldSend(newV);
      }
    }
  },
  methods: {
    popupCenter({url, title, w, h}: { url: string, title: string, w: number, h: number }): Window | null {
      // Fixes dual-screen position                             Most browsers      Firefox
      const dualScreenLeft = typeof window.screenLeft !==  'undefined' ? window.screenLeft : window.screenX;
      const dualScreenTop = typeof window.screenTop !==  'undefined'   ? window.screenTop  : window.screenY;

      const width = window.innerWidth ? window.innerWidth : document.documentElement.clientWidth ? document.documentElement.clientWidth : screen.width;
      const height = window.innerHeight ? window.innerHeight : document.documentElement.clientHeight ? document.documentElement.clientHeight : screen.height;

      const systemZoom = width / window.screen.availWidth;
      const left = ((width - w) / 2) + systemZoom + dualScreenLeft
      const top = ((height - h) / 2) + systemZoom + dualScreenTop
      const newWindow = window.open(url, title, 
        `
        scrollbars=yes,
        width=${w}, 
        height=${h}, 
        top=${top}, 
        left=${left}
        `
      )

      newWindow?.focus();
      return newWindow;
    },
    openTwitterLogin() {
      const connectUrl = import.meta.env.VITE_BACKEND_URL + '/user/connect_twitter';
      const loginWindow = this.popupCenter({ url: connectUrl, title: 'twitter_login_window', w: 800, h: 500});
      if (!loginWindow) {
        logEvent(analytics, 'connect_twitter', { success: false });
        this.$notify({ type: 'error', text: 'Failed to open window.' });
        return;
      }
      const vm = this;
      const handleFn = function(this: Window, e: MessageEvent) {
        if (e.origin !== import.meta.env.VITE_BACKEND_URL) return;
        if (typeof e.data === 'object' && 'message' in e.data) {
          const data = e.data;
          if (data['message'] !== 'twitter connect success') {
            logEvent(analytics, 'connect_twitter', { success: false });
            return;
          }
          logEvent(analytics, 'connect_twitter', { success: true });
          vm.$store.commit('SET_USER_CONNECTIONS', data['user_connections']);
          window.removeEventListener('message', handleFn);
          loginWindow.close();
        }
      }
      window.addEventListener('message', handleFn);
    },
    async skipLogin() {
      try {
        const resp = await client.get('/user/connect_email');
        const json = await resp.json();
        if (resp.status < 200 || resp.status > 299) {
          throw new Error(json['error_message']);
        }

        logEvent(analytics, 'connect_email', { success: true });
        this.$store.commit('SET_USER_CONNECTIONS', json['user_connections']);
      } catch (e) {
        logEvent(analytics, 'connect_email', { success: false });
        catchAndNotifyError(this, e); 
      }
    },
    async submitReply() {
      try {
        this.isSending = true;
        const resp = await client.postJson(`/messages/${this.message.recipient_id}/${this.message.id}/reply`, {
          'content': this.content
        });

        const json = await resp.json();
        if (resp.status < 200 || resp.status > 299) {
          if ('error_message' in json) {
            throw new Error(json['error_message']);
          } else if ('message' in json) {
            throw new Error(json['message']);
          } else {
            throw new Error('Something went wrong.');
          }
        }

        notify(this, { type: 'success', text: json['message'] });
        this.$emit('update:open', false);
        this.$emit('update:hasReplied', true);
        logEvent(analytics, 'reply-message', { id: this.message.id });
      } catch (e) {
        this.isSending = false;
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>