<template>
  <modal title="Submit a Message" :open="open" @update:open="$emit('update:open', $event)">
    <!-- TODO: fix inconsistencies in forms -->
    <form @submit.prevent="submitForm" class="flex flex-col">
      <div class="form-control">
        <label class="label">
          <span class="label-text">Recipient</span>
        </label>
        <input class="input input-bordered" type="text" name="student_id" v-model="recipientId" pattern="[0-9]{12}" placeholder="12-digit student ID">
      </div>

      <div class="flex flex-col mt-4 -mx-2 justify-center items-center">
        <p>Select gift</p>
        <fieldset class="gift-list-checkboxes">
          <div class="gift-item">
            <div class="gift-item-btn-wrapper">
              <input v-model="giftId" class="absolute appearance-none top-0 left-0" type="radio" name="gift_id" :value="null" id="none">
              <label class="btn btn-checkbox h-full w-full" for="none">
                None
              </label>
            </div>
          </div>
          <div class="gift-item" :key="'gift_' + gift.uid" v-for="(gift, i) in $store.state.giftList">
            <div class="gift-item-btn-wrapper">
              <input v-model="giftId" class="absolute appearance-none top-0 left-0" type="radio" name="gift_id" :value="gift.id" :id="gift.uid">
              <label class="btn btn-checkbox p-1 flex flex-col text-center h-full w-full" :for="gift.uid">
                <component :is="emojis[gift.uid]" class="text-4xl" />
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

// Emojis
import IconSunflower from '~icons/twemoji/sunflower';
import IconRose from '~icons/twemoji/rose';
import IconBalloon from '~icons/twemoji/balloon';
import IconTeddyBear from '~icons/twemoji/teddy-bear';
import IconRing from '~icons/twemoji/ring';
import IconMoney from '~icons/twemoji/money-with-wings';
import IconHeart from '~icons/twemoji/heart-with-ribbon';
import IconChocolate from '~icons/twemoji/chocolate-bar';
import IconPizza from '~icons/twemoji/pizza';

export default {
  components: { 
    Modal, 
    ContentCounter,
    IconSunflower,
    IconRose,
    IconBalloon,
    IconTeddyBear,
    IconRing,
    IconMoney,
    IconHeart,
    IconChocolate,
    IconPizza,
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
  computed: {
    emojis() {
      return {
        'sunflower': IconSunflower,
        'rose': IconRose,
        'balloons': IconBalloon,
        'teddy-bear': IconTeddyBear,
        'ring': IconRing,
        'money': IconMoney,
        'heart': IconHeart,
        'chocolate': IconChocolate,
        'pizza': IconPizza,
      }
    }
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
              // TODO: log gift in analytics
              gift_id: this.giftId,
              uid: this.$store.state.user.id
            })
          });
  
          const json = await resp.json();
          if (resp.status >= 200 && resp.status <= 299) {
            this.recipientId = '';
            this.content = '';
            logEvent(analytics, 'post-message');
            notify(this, { type: 'success', text: json['message'] });
            this.$emit('update:open', false);
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

.gift-list-checkboxes .gift-item .gift-item-btn-wrapper input[type="radio"]:checked + label {
  @apply bg-rose-100 border-rose-600;
}

.btn.btn-checkbox {
  @apply border-gray-500 hover:border-gray-700 bg-gray-50 hover:bg-gray-100 text-gray-900 hover:text-gray-900 inline-flex items-center justify-center;
}
</style>