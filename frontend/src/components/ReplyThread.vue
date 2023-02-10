<template>
  <div class="space-y-2">
    <div class="p-6 lg:p-8 bg-white rounded-xl shadow-lg">
      <div v-if="authState.user!.expand.details?.student_id == message!.recipient"
        class="flex space-x-2 items-center text-2xl">
        <icon-reply class="text-pink-500 mb-4" />
        <h2 class="font-bold mb-4">Your reply</h2>
      </div>
      <reply-message-box @update:hasReplied="handleHasReplied" />
    </div>

    <div class="p-6 lg:p-8 bg-white rounded-xl shadow-lg">
      <response-handler :query="repliesQuery">
        <template #default>
          <div :id="`reply_` + reply.id" v-for="reply in replies">
            <!-- TODO: improve UI -->
            {{ reply.content }}
          </div>
        </template>
      </response-handler>
    </div>    
  </div>
</template>

<script lang="ts" setup>
import ResponseHandler from './ResponseHandler2.vue';
import ReplyMessageBox from './ReplyMessageBox.vue';

import { inject, Ref } from 'vue';
import { useAuth } from '../store_new';
import { Record as PbRecord } from 'pocketbase';
import { useRouter } from 'vue-router';
import { useQuery } from '@tanstack/vue-query';
import { pb } from '../client';

const router = useRouter();
const message = inject<Ref<PbRecord>>('message')!;
const { state: authState } = useAuth();

function handleHasReplied(hasReplied: boolean) {
  message.value.replies_count++;
  if (hasReplied) {
    router.go(0);
  }
}

const repliesQuery = useQuery(['replies', message.value.id], () => {
  return pb.collection('message_replies').getFullList(undefined, {
    sort: '-created'
  });
});

const replies = repliesQuery.data;
</script>