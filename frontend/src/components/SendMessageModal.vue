<template>
  <modal title="Submit a Message" :open="open" @update:open="$emit('update:open', $event)" with-closing-button>
    <form @submit.prevent="submitForm" class="flex flex-col">
      <div class="form-control">
        <label class="label">
          <span class="label-text">Recipient</span>
        </label>
        <input class="input input-bordered" type="text" name="recipient_id" @input="recipientId = $event.target.value" pattern="[0-9]{6,12}" placeholder="12-digit student ID">
      </div>

      <div class="flex flex-col mt-4 -mx-2 justify-center items-center">
        <p>Select gift</p>
        <fieldset class="gift-list-checkboxes">
          <div class="gift-item" :key="'gift_' + gift.uid" v-for="(gift, i) in $store.state.giftList">
            <div class="gift-item-btn-wrapper">
              <input class="absolute appearance-none top-0 left-0" type="checkbox" :name="'gift_ids['+gift.id+']'" :id="gift.uid">
              <label class="btn btn-checkbox p-1 flex flex-col text-center h-full w-full" :for="gift.uid">
                <gift-icon :uid="gift.uid" class="text-4xl" />
                <span class="mt-2">{{ gift.label }}</span>
              </label>
            </div>
          </div>
        </fieldset>
      </div>

      <div class="form-control">
        <label class="label">
          <span class="label-text">Message</span>
        </label>
        <textarea 
          name="content"
          class="textarea textarea-bordered h-48" 
          @input="content = $event.target.value" 
          maxlength="240"></textarea>
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
import ContentCounter from './ContentCounter.vue';
import Modal from './Modal.vue';
import GiftIcon from './GiftIcon.vue';

import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';
import client from '../client';

export default {
  components: { 
    Modal, 
    ContentCounter,
    GiftIcon
  },
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
      giftId: null as unknown as number,
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
      if (this.shouldSend) {
        this.shouldSend = false;
      } else {
        return;
      }

      if (!e.target || !(e.target instanceof HTMLFormElement)) return;
      const formData = new FormData(e.target);
      const giftIdRegex = /^gift_ids\[([0-9]+)\]$/;

      try {
        let giftIds: number[] = [];
        formData.forEach((value, key) => {
          if (key.startsWith('gift_ids') && value === 'on') {
            const matches = giftIdRegex.exec(key);
            if (!matches || matches.length == 0) return;
            giftIds.push(parseInt(matches[1]));
          }
        });

        if (giftIds.length > 3) {
          throw new Error('Maximum of 3 gifts is allowed.');
        }

        const resp = await client.postJson('/messages', {
          recipient_id: formData.get('recipient_id'),
          content: formData.get('content'),
          gift_ids: giftIds,
          uid: this.$store.state.user.id
        });

        const json = await resp.json();
        if (resp.status >= 200 && resp.status <= 299) {
          logEvent(analytics, 'post-message');
          notify(this, { type: 'success', text: json['message'] });
          e.target.reset();
          this.$emit('update:open', false);
          this.$router.push(json['route']);
        } else if ('error_message' in json) {
          throw new Error(json['error_message']);
        } else {
          throw new Error('There are errors in your submission.');
        }
      } catch(e) {
        catchAndNotifyError(this, e);
      } finally {
        this.shouldSend = true;
      }
    }
  }
}
</script>

<style lang="postcss" scoped>
.gift-list-checkboxes {
  @apply flex flex-wrap items-center justify-center;
}

.gift-list-checkboxes .gift-item {
  @apply p-2 w-1/4;
}

.gift-list-checkboxes .gift-item .gift-item-btn-wrapper {
  height: 6.3rem;
  @apply w-full relative;
}

.gift-list-checkboxes .gift-item .gift-item-btn-wrapper input[type="checkbox"]:checked + label {
  @apply bg-rose-100 border-rose-600;
}

.btn.btn-checkbox {
  @apply border-gray-500 hover:border-gray-700 bg-gray-50 hover:bg-gray-100 text-gray-900 hover:text-gray-900 inline-flex items-center justify-center;
}
</style>