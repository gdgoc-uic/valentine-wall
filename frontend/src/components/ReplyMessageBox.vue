<template>
  <div 
    v-if="!authState.isLoggedIn || authState.user!.expand.details.student_id != message.recipient" 
    class="flex flex-col md:flex-row items-center">
    <icon-reply-lock class="text-gray-500 text-6xl mb-4 md:mb-0" />
    <div class="flex flex-col text-center items-center md:text-left md:items-start md:ml-4">
      <h3 class="text-2xl font-bold">Reply feature is locked.</h3>
      <p class="text-gray-500 text-xl">Recipient can only reply to this message.</p>
    </div>
  </div>

  <div v-else-if="message.replies_count > 0" class="flex flex-col justify-center text-center items-center">
    <icon-reply class="text-pink-500 text-9xl mb-4" />
    <h3 class="text-2xl font-bold">Already replied!</h3>
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

<script lang="ts" setup>
import IconReplyLock from '~icons/uil/comment-lock';
import IconSend from '~icons/uil/message';
import IconReply from '~icons/uil/comment-heart';
import ContentCounter from './ContentCounter.vue';
import { logEvent } from '@firebase/analytics';
// import { analytics } from '../firebase';
import { catchAndNotifyError, notify } from '../notify';
import { ref, computed } from 'vue';
import { useAuth } from '../store_new';

const emit = defineEmits(['update:hasReplied']);
defineProps({
  message: {
    type: Object,
    default: {}
  }
});

const { state: authState } = useAuth();
const counter = ref<InstanceType<typeof ContentCounter> | null>(null);
const content = ref('');
const isSending = ref(false);
const shouldSend = computed(() => counter.value?.shouldSend(content.value) ?? false);

async function submitReply() {
  try {
    isSending.value = true;
    // TODO:
    // const { data: json } = await this.$client.postJson(`/messages/${this.message.recipient_id}/${this.message.id}/reply`, {
    //   reply: {
    //     content: this.content
    //   },
    // });

    // notify(this, { type: 'success', text: json['message'] });
    // this.$store.commit('SET_USER_WALLET_BALANCE', json['current_balance']);
    emit('update:hasReplied', true);
    // logEvent(analytics!, 'reply-message', { id: this.message.id });
  } catch (e) {
    isSending.value = false;
    // catchAndNotifyError(this, e);
  }
}
</script>