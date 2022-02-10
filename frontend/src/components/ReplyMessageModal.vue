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
import { connectToEmail, connectToTwitter } from '../auth';

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
    openTwitterLogin() {
      let success: boolean = false;
      connectToTwitter(this.$store, {
        onSuccess: () => { success = true; },
        onError: (e) => {
          success = false;
          catchAndNotifyError(this, e);
        },
        onFinally: () => {
          logEvent(analytics!, 'connect_twitter', { success });
        }
      });
    },
    async skipLogin() {
      let success: boolean = false;
      try {
        connectToEmail(this.$store);
        success = true;
      } catch (e) {
        success = false;
        catchAndNotifyError(this, e); 
      } finally {
        logEvent(analytics!, 'connect_email', { success });
      }
    },
    async submitReply() {
      try {
        this.isSending = true;
        const { data: json } = await this.$client.postJson(`/messages/${this.message.recipient_id}/${this.message.id}/reply`, {
          'content': this.content
        });

        notify(this, { type: 'success', text: json['message'] });
        this.$store.commit('SET_USER_WALLET_BALANCE', json['current_balance']);
        this.$emit('update:open', false);
        this.$emit('update:hasReplied', true);
        logEvent(analytics!, 'reply-message', { id: this.message.id });
      } catch (e) {
        this.isSending = false;
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>