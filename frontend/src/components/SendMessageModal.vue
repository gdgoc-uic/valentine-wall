<template>
  <modal 
    :open="open" 
    title="Submit a Message"
    @update:open="$emit('update:open', $event)"
    with-closing-button>
    <template #modal-box>
      <div class="bg-white max-w-5xl w-full rounded-xl flex p-8 flex-col">
        <div class="relative">
          <h2 class="text-left text-3xl font-bold mb-4">Send a Message</h2>
          <button @click="$emit('update:open', false)"
                  class="bg-rose-500 hover:bg-rose-600 transition-colors text-white p-2 rounded-full absolute top-0 right-0">
            <icon-close />
          </button>
        </div>

        <div class="flex flex-col-reverse lg:flex-row max-h-[80vh] overflow-y-scroll md:overflow-y-hidden md:max-h-full">
          <div class="flex flex-col lg:w-2/3 lg:pr-8 overflow-y-none lg:overflow-y-auto md:max-h-[80vh]">
            <send-message-form
              @success="handleSendSuccess"
              :existing-recipient="route.params.recipientId?.toString()" />
          </div>        

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
import SendMessageForm from './SendMessageForm.vue';
import Modal from './Modal.vue';
import IconClose from '~icons/uil/multiply';
import { VueComponent as RulesContent } from '../assets/texts/rules.md';
import { useRoute, useRouter } from 'vue-router';
const emit = defineEmits(['update:open']);
const props = defineProps({
  open: {
    type: Boolean
  }
});

const router = useRouter();
const route = useRoute();

function handleSendSuccess(recipientId: string, messageId: string) {
  // logEvent(analytics!, 'post-message');
  emit('update:open', false);
  // this.$router.push(json['route']);
  router.push({
    name: 'message-page',
    params: { recipientId, messageId }
  });
}
</script>
