<template>
  <div 
    v-if="!$store.getters.isLoggedIn || $store.state.user.associatedId != message.recipient_id" 
    class="flex flex-col md:flex-row items-center">
    <icon-reply-lock class="text-gray-500 text-6xl mb-4 md:mb-0" />
    <div class="flex flex-col text-center items-center md:text-left md:items-start md:ml-4">
      <h3 class="text-2xl font-bold">Reply feature is locked.</h3>
      <p class="text-gray-500 text-xl">Recipient can only reply to this message.</p>
    </div>
  </div>

  <div v-else-if="message.has_replied" class="flex flex-col justify-center text-center items-center">
    <icon-reply class="text-pink-500 text-9xl mb-4" />
    <h3 class="text-2xl font-bold">Already replied!</h3>
  </div>

  <div v-else-if="$store.getters.isLoggedIn && !$store.getters.hasConnections" class="flex flex-col justify-center text-center items-center space-y-8">
    <p class="text-gray-500 text-lg">Connect your Twitter account and gain more coins!</p>
    
    <div class="w-full md:w-2/3 flex flex-col md:flex-row space-y-2 md:space-y-0 md:space-x-2 items-stretch justify-center">
      <button class="w-full md:w-1/2 space-x-2 btn bg-twitter-500 hover:bg-twitter-600 rounded-lg border-none hover:border-none" @click="useTwitter">
        <icon-twitter />
        <span>Login to Twitter</span>
      </button>
      <button class="w-full md:w-1/2 btn font-normal" @click="useEmail">Skip</button>
    </div>
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
      <div class="indicator">
        <div class="indicator-item badge badge-primary">ღ150.0</div> 
        <button @click="submitReply" class="space-x-2 btn bg-rose-500 hover:bg-rose-600 border-none hover:border-none" :disabled="!shouldSend || isSending">
          <icon-send />
          <span>Send</span>
        </button>
      </div>
    </div>
  </div>
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
import ContentCounter from './ContentCounter.vue';
import { logEvent } from '@firebase/analytics';
// import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';
import { connectToEmail, connectToTwitter } from '../auth';
import {pb} from '../client';

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
  emits: ['update:hasReplied'],
  props: {
    message: {
      type: Object,
      default: {}
    }
  },
  setup() {
    return {
      pb
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
    useTwitter() {
      let success: boolean = false;
      connectToTwitter(this.$store, {
        onSuccess: () => { success = true; },
        onError: (e) => {
          success = false;
          catchAndNotifyError(this, e);
        },
        onFinally: () => {
          // logEvent(analytics!, 'connect_twitter', { success });
        }
      });
    },
    async useEmail() {
      let success: boolean = false;
      this.$notify({ type: 'info', text: 'You can connect it later ' })
      try {
        connectToEmail(this.$store);
        success = true;
      } catch (e) {
        success = false;
        catchAndNotifyError(this, e); 
      } finally {
        // logEvent(analytics!, 'connect_email', { success });
      }
    },
    async submitReply() {
      try {
        this.isSending = true;
        // const { data: json } = await this.$client.postJson(`/messages/${this.message.recipient_id}/${this.message.id}/reply`, {
        //   reply: {
        //     content: this.content
        //   },
        //   options: {
        //     post_to_email: true,
        //     post_to_twitter: true
        //   }
        // });

        // notify(this, { type: 'success', text: json['message'] });
        // this.$store.commit('SET_USER_WALLET_BALANCE', json['current_balance']);
        this.$emit('update:hasReplied', true);
        // logEvent(analytics!, 'reply-message', { id: this.message.id });
      } catch (e) {
        this.isSending = false;
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>