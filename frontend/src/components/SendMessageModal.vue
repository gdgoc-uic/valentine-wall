<template>
  <modal title="Submit a Message" :open="open" @update:open="$emit('update:open', $event)">
    <form @submit.prevent="submitForm" class="flex flex-col">
      <div class="form-control">
        <label class="label">
          <span class="label-text">Recipient student ID</span>
        </label>
        <input class="input input-bordered" type="text" name="student_id" v-model="recipientId" pattern="[0-9]{12}" placeholder="12-digit student ID">
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text">Message</span>
        </label>
        <textarea class="textarea textarea-bordered h-48" v-model="content" maxlength="240"></textarea>
      </div>
      <div class="flex justify-between items-center mt-4">
        <content-counter ref="counter" :content="content" :newline-count="13" />
        <button 
          class="self-end px-12 btn bg-rose-500 hover:bg-rose-600 border-none" 
          type="submit" 
          :disabled="!shouldSend">Send</button>
      </div>
    </form>
  </modal>
</template>

<script lang="ts">
import { logEvent } from '@firebase/analytics';
import ContentCounter from './ContentCounter.vue';
import Modal from './Modal.vue';
import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';

export default {
  components: { Modal, ContentCounter },
  emits: ['update:open'],
  props: {
    open: {
      type: Boolean,
    }
  },
  data() {
    return {
      recipientId: '',
      content: '',
      shouldSend: false
    };
  },
  watch: {
    content(newV: string, oldV: string) {
      const counter = (this.$refs.counter as typeof ContentCounter);
      if (!counter || typeof counter == 'undefined') {
        this.shouldSend = false;
      } else {
        this.shouldSend = this.recipientId.length > 0 && counter.shouldSend(newV);
      }
    }
  },
  methods: {
    async submitForm(e: SubmitEvent) {
      if (!this.shouldSend) return;
      let target = e.target;
      if (target && target instanceof HTMLFormElement) {
        try {
          const resp = await fetch(import.meta.env.VITE_BACKEND_URL + '/messages', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              ...this.$store.getters.headers,
            },
            body: JSON.stringify({
              recipient_id: this.recipientId,
              content: this.content,
              uid: this.$store.state.user.id
            })
          });
  
          const json = await resp.json();
          if (resp.status >= 200 && resp.status <= 299) {
            this.recipientId = '';
            this.content = '';
            logEvent(analytics, 'post-message');
            notify(this, { type: 'success', text: json['message'] });
          } else if ('error_message' in json) {
            throw new Error(json['error_message']);
          } else {
            throw new Error('There are errors in your submission');
          }
        } catch(e) {
          catchAndNotifyError(this, e);
        }
      }
    }
  }
}
</script>