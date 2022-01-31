<template>
  <modal title="Settings" :open="open" @update:open="$emit('update:open', $event)" with-closing-button>
    <div class="flex flex-col">
      <form @submit.prevent="proceedDelete" class="flex flex-col p-4 border border-red-500 rounded-md">
        <div class="mb-8">
          <h2 class="text-red-500 text-xl font-bold mb-2">Delete Account</h2>
          <div class="space-y-2">
            <p>By deleting your account:</p>
            <ul class="list-disc pl-5">
              <li>Data such as your messages and your replies sent to others won't be removed.</li>
              <li>Any existing social media accounts connected to this account will be removed.</li>
              <li>You will be permanently losing access to your previous data the next time you re-register your account.</li>
            </ul>
            <p>To proceed, please enter your student ID.</p>
          </div>
        </div>

        <div class="form-control">
          <div class="flex space-x-2">
            <input
              @input="recipientId = $event.target.value"
              pattern="[0-9]{6,12}"
              :placeholder="$store.state.user.associatedId"
              class="w-full input input-bordered"
              name="recipient_id" type="text">
            <button :disabled="!shouldDelete" class="btn btn-error bg-red-500 border-red-600 hover:bg-red-600 hover:border-red-700">Delete</button>
          </div>
        </div>
      </form>
    </div>
  </modal>
</template>

<script lang="ts">
import { catchAndNotifyError, notify } from '../notify';
import Modal from './Modal.vue';
export default {
  components: { Modal },
  emits: ['update:open'],
  props: {
    open: {
      type: Boolean
    }
  },
  data() {
    return {
      recipientId: '',
      shouldDelete: false
    }
  },
  watch: {
    recipientId(newV: string) {
      if (newV === this.$store.state.user.associatedId) {
        this.shouldDelete = true;
      } else {
        this.shouldDelete = false;
      }
    }
  },
  methods: {
    async proceedDelete(e: SubmitEvent) {
      if (this.shouldDelete) {
        this.shouldDelete = false;
      } else {
        return;
      }

      try {
        if (!e.target || !(e.target instanceof HTMLFormElement)) return;
        const formData = new FormData(e.target);
        const { data: json } = await this.$client.postJson('/messages', {
          input_sid: formData.get('recipient_id'),
          input_uid: this.$store.state.user.id
        });

        notify(this, { type: 'success', text: json['message'] });
        this.$router.replace({ name: 'home-page' });
        await this.$store.dispatch('logout');
      } catch(e) {
        catchAndNotifyError(this, e);
      } finally {
        this.shouldDelete = true;
      }
    }
  }
}
</script>