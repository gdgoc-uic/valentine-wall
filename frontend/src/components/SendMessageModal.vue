<template>
  <modal 
    :open="open" 
    title="Submit a Message"
    @update:open="$emit('update:open', $event)"
    with-closing-button>
    <template #modal-box>
      <div class="bg-white max-w-7xl w-full rounded-xl flex p-8 flex-col">
        <div class="relative">
          <h2 class="text-left text-3xl font-bold mb-4">Send a Message</h2>
          <button @click="$emit('update:open', false)"
                  class="bg-rose-500 hover:bg-rose-600 transition-colors text-white p-2 rounded-full absolute top-0 right-0">
            <icon-close />
          </button>
        </div>

        <div class="flex">
          <form @submit.prevent="submitForm" class="flex flex-col w-2/3 pr-8 overflow-y-auto"  style="max-height: 80vh">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Recipient</span>
              </label>
              <input 
                :value="$route.params.recipientId || ''"
                class="input input-bordered" 
                type="text"
                name="recipient_id" 
                @input="recipientId = $event.target.value" 
                pattern="[0-9]{6,12}" 
                placeholder="6 to 12-digit Student ID (e.g. 200xxxxxxxxx)">
            </div>
            <div class="flex flex-col mt-4 -mx-2">
              <p class="pl-3 text-gray-900 text-sm my-2">Select gift (A maximum of 3 may be attached)</p>
              <fieldset class="gift-list-checkboxes">
                <div class="gift-item tooltip tooltip-bottom" :data-tip="gift.label" :key="'gift_' + gift.uid" v-for="gift in $store.state.giftList">
                  <div class="gift-item-btn-wrapper">
                    <input class="absolute appearance-none top-0 left-0" type="checkbox" :name="'gift_ids['+gift.id+']'" :id="gift.uid">
                    <label class="btn btn-checkbox rounded-xl p-1 flex flex-col text-center h-full w-full" :for="gift.uid">
                      <gift-icon :uid="gift.uid" class="text-4xl" />
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

          <div class="bg-white w-1/3 border shadow-lg rounded-xl p-5">
            <h2 class="text-rose-500 text-2xl font-semibold">Rules</h2>
            <rules-content />
          </div>
        </div>
      </div>
    </template>
  </modal>
</template>

<script lang="ts">
import ContentCounter from './ContentCounter.vue';
import Modal from './Modal.vue';
import GiftIcon from './GiftIcon.vue';
import IconClose from '~icons/uil/multiply';

import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';
import { VueComponent as RulesContent } from '../assets/texts/rules.md';

export default {
  components: { 
    Modal, 
    ContentCounter,
    GiftIcon,
    IconClose,
    RulesContent
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

        const { data: json } = await this.$client.postJson('/messages', {
          recipient_id: formData.get('recipient_id'),
          content: formData.get('content'),
          gift_ids: giftIds,
          uid: this.$store.state.user.id
        }, {
          headers: this.$store.getters.headers
        });

        logEvent(analytics!, 'post-message');
        notify(this, { type: 'success', text: json['message'] });
        e.target.reset();
        this.$emit('update:open', false);
        this.$router.push(json['route']);
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
  @apply p-2 w-1/5;
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