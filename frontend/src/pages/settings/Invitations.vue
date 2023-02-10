<template>
  <div class="flex flex-col space-y-3">
    <div class="flex flex-row w-full justify-end">
      <button class="btn btn-sm btn-success" @click="openGenerateModal = true">
        <icon-plus />
        <span>Generate Invitation</span>
      </button>

      <portal>
        <generate-invitation-modal 
          :key="modalKey"
          :open="openGenerateModal"
          @update:open="handleModalOpen"
          @success="iResponseKey++" />
      </portal>
    </div>

    <response-handler :key="iResponseKey" endpoint="/user/invitations">
      <template #default="{ response: { data: invitations } }">
        <table class="table w-full">
          <thead>
            <tr>
              <th class="w-2/4 text-md normal-case text-red-400">
                <span>ID</span>
              </th>
              <th class="w-1/6 text-md normal-case text-red-400">
                <span>Max Users</span>
              </th>
              <th class="w-1/6 text-md normal-case text-red-400">
                <span>User Count</span>
              </th>
              <th class="w-1/6 text-md text-center normal-case text-red-400">
                <span>Created At</span>
              </th>
              <th class="w-1/6 text-md text-center normal-case text-red-400">
                <span>Expires At</span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr :key="iv.id" v-for="(iv, i) in invitations" :class="{'border-b-2': i < invitations.length - 1}">
              <td class="text-md font-bold text-gray-700">{{iv.id}}</td>
              <td class="text-md font-semibold text-gray-500">{{ iv.max_users }}</td>
              <td class="text-md font-semibold text-gray-500">{{ iv.user_count }}</td>
              <td class="text-md text-gray-500 text-center">{{ prettifyDate(iv.created_at) }}</td>
              <td class="text-md text-gray-500 text-center">{{ prettifyDate(iv.expires_at) }}</td>
            </tr>
          </tbody>
        </table>
      </template>
      <template #error="{ error }">
        <p>{{ error.message || error }}</p>
      </template>
    </response-handler>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { prettifyDateTime } from '../../time_utils';

import ResponseHandler from '../../components/ResponseHandler.vue';
import IconPlus from '~icons/uil/plus'
import Portal from '../../components/Portal.vue';
import GenerateInvitationModal from '../../components/GenerateInvitationModal.vue';


// TODO: user invitations

export default defineComponent({
  components: { 
    ResponseHandler, 
    IconPlus, 
    Portal,
    GenerateInvitationModal,
  },
  data() {
    return {
      openGenerateModal: false,
      maxUsers: 1,
      iResponseKey: 1,
      modalKey: 1,
    }
  },
  methods: {
    prettifyDate(date: Date) {
      return prettifyDateTime(date);
    },
    handleModalOpen(newOpen: boolean) {
      this.openGenerateModal = newOpen;
      if (!newOpen) {
        this.modalKey++;
      }
    }
  }
})
</script>