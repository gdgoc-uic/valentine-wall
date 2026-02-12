<template>
  <div v-if="message" class="space-y-2">
    <div v-if="!isReadOnly()" class="p-6 lg:p-8 bg-white rounded-xl shadow-lg">
      <div v-if="authState.isLoggedIn && (message.recipient === 'everyone' || authState.user!.expand.details?.student_id == message.recipient)"
        class="flex space-x-2 items-center text-2xl">
        <icon-reply class="text-pink-500 mb-4" />
        <h2 class="font-bold mb-4">Your reply</h2>
      </div>
      <reply-message-box @update:hasReplied="handleHasReplied" />
    </div>

    <div v-if="replies && replies.length != 0" class="bg-white rounded-xl shadow-lg">
      <response-handler :query="repliesQuery">
        <template #default>
          <section class="flex flex-col text-gray-800">
            <div :id="`reply_` + reply.id" v-for="reply in replies" class="flex p-6 lg:p-8 border-b">
              <div class="flex-1 flex flex-col">
                <div class="flex items-start space-x-2">
                  <span
                    :class="{ 'text-rose-500': reply.expand.sender.student_id == message.recipient }"
                    class="font-bold text-lg mb-2">{{ reply.expand.sender.student_id  }}</span>
                  <span class="text-gray-500">{{ fromNow(reply.created) }}</span>
                </div>
                <p class="text-lg break-words">{{ reply.content }}</p>
                <div class="space-x-1 text-sm flex mt-6 items-center" v-if="reply.liked && (authState.isLoggedIn && reply.sender === authState.user.details)">
                  <icon-heart class="text-rose-500" />

                  <span class="text-gray-700">
                    Your reply was liked by <b>{{ message.recipient }}</b>
                  </span>
                </div>
              </div>

              <div v-if="!isReadOnly() && authState.isLoggedIn" class="flex flex-col space-y-4 items-start">
                <button 
                  v-if="authState.user.expand.details.student_id === message.recipient"
                  :class="[reply.liked ? 'text-white bg-rose-500 hover:bg-rose-700 border-rose-500' : 'bg-white border-gray-200 text-rose-500 hover:border-rose-500 hover:bg-rose-500']"
                  @click="like(reply as PbRecord)"
                  class="btn-sm btn btn-circle shadow-md hover:text-white">
                  <icon-heart />
                </button>

                <delete-dialog @confirm="(confirmed: boolean) => {if (confirmed){ deleteReply(reply.id)}}">
                  <template #default="{ openDialog }">
                    <button 
                      v-if="message.recipient === authState.user.expand.details.student_id"
                      @click="openDialog"
                      class="btn-sm btn btn-circle shadow-md !text-gray-900 bg-white border-gray-200 hover:bg-gray-200">
                      <icon-trash />
                    </button>
                  </template>
                </delete-dialog>
              </div>
            </div>
          </section>
        </template>
      </response-handler>
    </div>    
  </div>
</template>

<script lang="ts" setup>
import DeleteDialog from '../components/DeleteDialog.vue';
import IconTrash from '~icons/uil/trash-alt';
import IconHeart from '~icons/uil/heart';
import IconReply from '~icons/uil/comment-heart';
import ResponseHandler from './ResponseHandler2.vue';
import ReplyMessageBox from './ReplyMessageBox.vue';

import { inject, Ref, ref, computed, onMounted, onUnmounted } from 'vue';
import { useAuth } from '../store_new';
import { Record as PbRecord, UnsubscribeFunc } from 'pocketbase';
import { useMutation, useQuery } from '@tanstack/vue-query';
import { pb } from '../client';
import { fromNow } from '../time_utils';
import { isReadOnly } from '../utils';

const message = inject<Ref<PbRecord>>('message')!;
const { state: authState } = useAuth();

function handleHasReplied() {
  if (!message?.value) return;
  message.value.replies_count++;
  repliesQuery.refetch();
}

const repliesQuery = useQuery(
  computed(() => ['replies', message?.value?.id]),
  () => {
    if (!message?.value) return Promise.reject('No message');
    return pb.collection('message_replies').getFullList(undefined, {
      filter: `message="${message.value.id}"`,
      sort: '-created',
      expand: 'sender'
    });
  }, {
  onSuccess(data) {
    replies.value = data;
  },
  enabled: computed(() => {
    if (!message?.value) return false;
    return (
      message.value.recipient == 'everyone' ||
      (authState.isLoggedIn && (
        message.value.recipient == authState.user!.expand.details?.student_id || 
        message.value.user == authState.user.id
      ))
    );
  })
});

const replies = ref<PbRecord[]>([]);

const { mutate: like } = useMutation((r: PbRecord) => {
  return pb.collection('message_replies').update(r.id, { liked: !r.liked });
});

const { mutate: deleteReply } = useMutation((id: string) => {
  return pb.collection('message_replies').delete(id);
});

const unsubscribeFunc = ref<UnsubscribeFunc | null>(null);

onMounted(() => {
  pb.collection('message_replies').subscribe('*', (data) => {
    if (!message?.value || data.record.message !== message.value.id) {
      return;
    }

    if (data.action === 'delete') {
      replies.value = replies.value.filter(r => r.id !== data.record.id);
    } else {
      pb.collection(data.record.collectionName).getOne(data.record.id, {
        expand: 'sender'
      }).then(record => {
        if (data.action === 'create') {
          replies.value.unshift(record);
        } else if (data.action === 'update') {
          replies.value = replies.value.map(r => {
            if (r.id === record.id) {
              return {
                ...record,
                expand: r.expand
              } as PbRecord;
            } else {
              return r;
            }
          });
        }
      });
    }
  });
});

onUnmounted(() => {
  unsubscribeFunc.value?.();
});
</script>