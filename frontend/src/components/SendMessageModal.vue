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

        <div class="flex flex-col-reverse lg:flex-row max-h-[80vh] overflow-y-scroll md:overflow-y-hidden md:max-h-full">
          <form ref="submitMessageForm" @submit.prevent="submitForm" class="flex flex-col lg:w-2/3 lg:pr-8 overflow-y-none lg:overflow-y-auto md:max-h-[80vh]">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Recipient</span>
              </label>
              <input
                class="input input-bordered" 
                type="text"
                name="recipient_id" 
                @input="recipientId = $event.target!.value" 
                pattern="[0-9]{6,12}" 
                placeholder="6 to 12-digit Student ID (e.g. 200xxxxxxxxx)">
            </div>
            <div class="flex flex-col mt-4 -mx-2">
              <p class="pl-3 text-gray-900 text-sm my-2">Select gift (Optional)</p>
              <fieldset class="gift-list-checkboxes">
                <div class="gift-item tooltip tooltip-top z-10" :data-tip="gift.label" :key="'gift_' + gift.uid" v-for="gift in store.state.giftList">
                  <div class="gift-item-btn-wrapper indicator">
                    <div class="indicator-bottom indicator-center indicator-item badge badge-primary">ღ{{ gift.price }}</div> 
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
                @input="content = $event.target!.value"
                maxlength="240"></textarea>
            </div>
            <div class="flex justify-between items-center mt-4">
              <content-counter ref="counter" :content="content" :newline-count="13" />
              <div class="indicator">
                <div class="indicator-item badge badge-primary">ღ150.0</div> 
                <button
                  class="self-end px-12 btn bg-rose-500 hover:bg-rose-600 border-none"
                  type="submit"
                  :disabled="!shouldSend">Send</button>
              </div>
            </div>
          </form>

          <div class="bg-white lg:w-1/3 border shadow-lg rounded-xl p-5">
            <h2 class="text-rose-500 text-2xl font-semibold">Rules</h2>
            <rules-content />
          </div>
        </div>
      </div>
    </template>
  </modal>
</template>

<script lang="ts" setup>
import ContentCounter from './ContentCounter.vue';
import Modal from './Modal.vue';
import GiftIcon from './GiftIcon.vue';
import IconClose from '~icons/uil/multiply';

import { logEvent } from '@firebase/analytics';
// import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';
import { VueComponent as RulesContent } from '../assets/texts/rules.md';
import { ref, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useMutation } from '@tanstack/vue-query';
import { pb } from '../client';
import { useAuth, useStore } from '../store_new';

const emit = defineEmits(['update:open']);
const props = defineProps({
  open: {
    type: Boolean
  }
});

const { state: authState } = useAuth();
const store = useStore();
const route = useRoute();
const submitMessageForm = ref<HTMLFormElement | null>(null);
const counter = ref<InstanceType<typeof ContentCounter> | null>(null);
const recipientId = ref('');
const content = ref('');

function checkShouldSend(recipientId: string, content: string) {
  if (!counter.value) {
    return false;
  }
  return recipientId.length > 0 && counter.value.shouldSend(content);
}

function injectRecipientId() {
  // sets the recipient id automatically if route
  // has recipient id
  if (props.open && route.params.recipientId) {
    setTimeout(() => {
      if (!submitMessageForm.value || !(submitMessageForm.value instanceof HTMLFormElement)) {
        return;
      }

      const recipientIdField = submitMessageForm.value.querySelector('input[name="recipient_id"]');
      if (!recipientIdField || !(recipientIdField instanceof HTMLInputElement)) {
        return;
      }

      recipientIdField.value = <string> route.params.recipientId;
      recipientId.value = recipientIdField.value;
    }, 50);
  }
}

const shouldSend = computed(() => checkShouldSend(recipientId.value, content.value));

async function submitForm(e: SubmitEvent) {
  e.preventDefault();

  if (!shouldSend || !e.target || !(e.target instanceof HTMLFormElement)) return;
  const formData = new FormData(e.target);
  const giftIdRegex = /^gift_ids\[(\w+)\]$/;

  let gifts: string[] = [];
  formData.forEach((value, key) => {
    if (key.startsWith('gift_ids') && value === 'on') {
      const matches = giftIdRegex.exec(key);
      if (!matches || matches.length == 0) return;
      gifts.push(matches[1]);
    }
  });

  await sendMessage({
    gifts,
    user: authState.user!.details,
    recipient: formData.get('recipient_id')?.toString() ?? '',
    content: formData.get('content')?.toString() ?? ''
  }, {
    onSuccess() {
      // logEvent(analytics!, 'post-message');
      // this.$store.state.commit('SET_USER_WALLET_BALANCE', json['current_balance'])  
      (e.target! as HTMLFormElement).reset();
      emit('update:open', false);
      // this.$router.push(json['route']);
    }
  });
}

const { mutateAsync: sendMessage } = useMutation((message: {
  gifts: string[],
  content: string,
  user: string,
  recipient: string
}) => {
  if (message.gifts.length > 3) {
    throw new Error('Maximum of 3 gifts is allowed.');
  }

  return pb.collection('messages').create(message);
}, {
  onSuccess(data, variables, context) {
    // notify(this, { type: 'success', text: json['message'] });
  },
  onError(e) {
    // catchAndNotifyError(this, e);
  }
});

watch(props, () => {
  injectRecipientId();
});
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