<template>
  <div class="space-y-2">
    <div class="flex flex-col p-8 border border-gray-500 rounded-md">
      <h3 class="text-xl font-bold">Archive</h3>

      <div class="flex flex-col md:flex-row space-y-2 items-center">
        <div class="space-y-3">
          <p>You have <b>{{ messagesCount }}</b> messages received under your ID. To download the messages, click the "download archive" button.</p>
          <p class="text-sm text-gray-600">This will download a ZIP file containing the image versions of your messages. The information of the message senders won't be included upon download.</p>
        </div>

        <archive-dialog @error="catchAndNotifyError" v-slot="{ openDialog }">
          <button @click="openDialog" class="btn btn-lg">
            Download archive
          </button>
        </archive-dialog>
      </div>
    </div>
    <form @submit.prevent="proceedDelete" class="flex flex-col p-8 border border-red-500 rounded-md">
      <h3 class="text-xl font-bold">Delete Account</h3>
      <div class="mb-8">
        <div class="space-y-2">
          <p>By deleting your account:</p>
          <ul class="list-disc pl-5">
            <li>Data such as your messages and your replies sent to others won't be removed.</li>
            <li>Any existing connections to third-party accounts (e.g. Twitter) will be removed.</li>
            <li>Access to your current data will be lost upon re-registering.</li>
          </ul>
        </div>
      </div>
      <div class="form-control">
        <p class="text-xl mb-2">To proceed, please enter your student ID.</p>
        <div class="flex space-x-2">
          <input
            @input="recipientId = $event.target.value"
            pattern="[0-9]{6,12}"
            :placeholder="user!.expand.details?.student_id"
            class="w-full input input-bordered"
            name="recipient_id" type="text">
          <button :disabled="!shouldDelete"
            class="btn btn-error bg-red-500 border-red-600 hover:bg-red-600 hover:border-red-700">Delete</button>
        </div>
      </div>
    </form>
  </div>
</template>

<script lang="ts" setup>
import { useMutation, useQuery } from '@tanstack/vue-query';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { pb } from '../../client';
import { catchAndNotifyError, notify } from '../../notify';
import { useAuth } from '../../store_new';
import ArchiveDialog from '../../components/ArchiveDialog.vue';

const { state: { user }, methods: { logout } } = useAuth();
const router = useRouter();
const recipientId = ref('');
const shouldDelete = computed(() => {
  return recipientId.value === user!.expand.details?.student_id;
});

const { mutateAsync: deleteAccount } = useMutation(() => {
  return pb.collection('users').delete(user!.id);
}, {
  onSuccess(data) {
    notify({ type: 'success', text: 'Your account was deleted succesfully.' });
  },
  onSettled() {
    recipientId.value = '';
  }
})

const { data: messagesCount, isLoading } = useQuery(['message_count'], async () => {
  const resp = await pb.collection('messages').getList(1, 1, { filter: `recipient = "${user.expand.details.student_id}"` });
  return resp.totalItems;
}, {
  initialData: 0
})

function proceedDelete(e: SubmitEvent) {
  if (!e.target || !(e.target instanceof HTMLFormElement)) return;

  deleteAccount()
    .then(() => {
      router.replace({ name: 'home-page' });
      return logout();
    });
}
</script>