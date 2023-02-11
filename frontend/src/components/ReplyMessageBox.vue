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
        <div v-if="shouldSend" class="indicator-item badge badge-primary">ღ150.0</div> 
        <button @click="() => submitReply()" 
          class="space-x-2 btn bg-rose-500 hover:bg-rose-600 border-none hover:border-none" :disabled="!shouldSend || isSending">
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
import ContentCounter from './ContentCounter.vue';
import { logEvent } from '@firebase/analytics';
// import { analytics } from '../firebase';
import { notify } from '../notify';
import { ref, computed, Ref, inject } from 'vue';
import { useAuth } from '../store_new';
import { useMutation } from '@tanstack/vue-query';
import { pb } from '../client';
import { Record as PbRecord } from 'pocketbase';

const emit = defineEmits(['update:hasReplied']);
const message = inject<Ref<PbRecord>>('message')!;
const { state: authState } = useAuth();
const counter = ref<InstanceType<typeof ContentCounter> | null>(null);
const content = ref('');
const shouldSend = computed(() => counter.value?.shouldSend);

const { mutate: submitReply, isLoading: isSending } = useMutation(() => {
  return pb.collection('message_replies').create({
    content: content.value,
    sender: authState.user.details,
    message: message.value.id
  })
}, {
  onSuccess() {
    notify({ type: 'success', text: 'Reply was sent successfully.' });
    emit('update:hasReplied', true);
    // logEvent(analytics!, 'reply-message', { id: this.message.id });
  },
})
</script>