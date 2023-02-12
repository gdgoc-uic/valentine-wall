<template>
  <modal v-model:open="isGiftModalOpen" with-closing-button title="Select a Gift">
    <form @submit.prevent="saveVirtualGifts" class="flex flex-col">
      <p class="text-center text-gray-900 text-sm my-2">Up to three virtual gifts only!</p>
      <fieldset class="gift-list-checkboxes">
        <div class="gift-item tooltip tooltip-top z-10" :data-tip="gift.label" :key="'gift_' + gift.uid" v-for="gift in store.state.giftList">
          <div class="gift-item-btn-wrapper indicator">
            <div class="indicator-bottom indicator-center indicator-item badge badge-primary">ღ{{ gift.price }}</div> 
            <input class="absolute appearance-none top-0 left-0" type="checkbox" 
              :checked="gifts.includes(gift.id)" :name="'gift_ids['+gift.id+']'" :id="gift.uid">
            <label class="btn btn-checkbox rounded-xl p-1 flex flex-col text-center h-full w-full" :for="gift.uid">
              <gift-icon :uid="gift.uid" class="text-4xl" />
            </label>
          </div>
        </div>
      </fieldset>

      <button
        class="px-12 btn bg-rose-500 hover:bg-rose-600 border-none"
        type="submit">Save</button>
    </form>
  </modal>

  <modal v-model:open="isRulesModalOpen" with-closing-button title="Rules">
    <rules-content />
  </modal>

  <form ref="submitMessageForm" 
    @submit.prevent="submitForm">
    <div class="form-control">
      <label class="label">
        <span class="label-text">Recipient</span>
      </label>
      <input
        class="input input-bordered" 
        type="text"
        name="recipient_id"
        v-model="recipientId"
        pattern="[0-9]{6,12}" 
        placeholder="6 to 12-digit Student ID (e.g. 200xxxxxxxxx)">
    </div>
    <div class="form-control">
      <label class="label">
        <span class="label-text">Message</span>
      </label>
      <textarea
        name="content"
        class="textarea textarea-bordered h-48"
        v-model="content"
        maxlength="240"></textarea>
    </div>
    <div class="flex flex-col md:flex-row space-y-2 md:space-y-0 justify-between items-center mt-4">
      <div class="w-full md:w-auto flex space-x-2 items-stretch md:items-start">
        <button @click.prevent="isGiftModalOpen = true" class="flex-1 md:flex-auto btn btn-sm md:btn-md space-x-2 bg-white hover:bg-rose-200 border-rose-300 hover:border-rose-600 text-gray-800">
          <icon-plus class="text-rose-500" />          

          <span>Virtual Gift</span>
        </button>

        <button @click.prevent="isRulesModalOpen = true" class="flex-1 md:flex-auto btn btn-sm md:btn-md space-x-2 bg-white hover:bg-rose-200 border-rose-300 hover:border-rose-600 text-gray-800">
          <icon-rules class="text-rose-500" />          

          <span>Rules</span>
        </button>
      </div>

      <div class="w-full md:w-auto space-x-4 flex items-center justify-end">
        <content-counter ref="counter" :content="content" :newline-count="13" />
        <div class="indicator">
          <div v-if="shouldSend" class="indicator-item badge badge-primary">ღ{{ (SEND_PRICE + totalGiftPrice).toFixed(2) }}</div> 
          <button
            class="self-end px-12 btn bg-rose-500 hover:bg-rose-600 border-none"
            type="submit"
            :disabled="!shouldSend">Send</button>
        </div>
      </div>
    </div>
  </form>
</template>

<script lang="ts" setup>
import Modal from './Modal.vue';
import { useMutation } from '@tanstack/vue-query';
import { ref, computed } from 'vue';
import { pb } from '../client';
import { notify } from '../notify';
import { useAuth, useStore } from '../store_new';

import IconRules from '~icons/uil/list-ui-alt';
import IconPlus from '~icons/uil/plus';
import GiftIcon from './GiftIcon.vue';
import ContentCounter from './ContentCounter.vue';
import { VueComponent as RulesContent } from '../assets/texts/rules.md';

const SEND_PRICE = 150;
const emit = defineEmits(['success']);

const props = defineProps({
  existingRecipient: {
    type: String,
    required: false
  }
});

const { state: authState } = useAuth();
const store = useStore();
const submitMessageForm = ref<HTMLFormElement | null>(null);
const counter = ref<InstanceType<typeof ContentCounter> | null>(null);
const inputtedRecipient = ref('');
const isGiftModalOpen = ref(false);
const isRulesModalOpen = ref(false);
const content = ref('');
const shouldSend = computed(() => recipientId.value.length > 0 && counter.value?.shouldSend);
const gifts = ref<string[]>([]);
const recipientId = computed({
  get() {
    return props.existingRecipient ?? inputtedRecipient.value;
  },
  set(newValue) {
    inputtedRecipient.value = newValue;
  }
});
const totalGiftPrice = computed(() => store.state.giftList.filter(g => gifts.value.includes(g.id)).reduce((c, g) => c + g.price, 0));

function saveVirtualGifts(e: SubmitEvent) {
  e.preventDefault();
  gifts.value.splice(0, gifts.value.length);

  if (!shouldSend || !e.target || !(e.target instanceof HTMLFormElement)) return;

  const formData = new FormData(e.target);
  const giftIdRegex = /^gift_ids\[(\w+)\]$/;

  formData.forEach((value, key) => {
    if (key.startsWith('gift_ids') && value === 'on') {
      const matches = giftIdRegex.exec(key);
      if (!matches || matches.length == 0) return;
      gifts.value.push(matches[1]);
    }
  });

  e.target.reset();
  isGiftModalOpen.value = false;
}

async function submitForm(e: SubmitEvent) {
  e.preventDefault();

  if (!shouldSend || !e.target || !(e.target instanceof HTMLFormElement)) return;
  const formData = new FormData(e.target);

  await sendMessage({
    gifts: gifts.value,
    user: authState.user!.details,
    recipient: formData.get('recipient_id')?.toString() ?? '',
    content: formData.get('content')?.toString() ?? ''
  }, {
    onSuccess(record) {
      (e.target! as HTMLFormElement).reset();
      gifts.value.splice(0, gifts.value.length);
      content.value = '';
      recipientId.value = '';
      emit('success', record.recipient, record.id);
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
  onSuccess() {
    notify({ type: 'success', text: 'Message created successfully.' });
  }
});
</script>

<style lang="postcss" scoped>
.gift-list-checkboxes {
  @apply flex flex-wrap items-center justify-center;
}

.gift-list-checkboxes .gift-item {
  @apply p-2 w-1/3 md:w-1/5;
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