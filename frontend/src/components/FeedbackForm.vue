<template>
  <slot :openDialog="openDialog"></slot>

  <portal>
    <modal title="Feedback" v-model:open="open" with-closing-button>
      <p class="mb-4">Suggestions, problems, and any other feedback can be addressed here.</p>

      <form v-if="!isSuccess" @submit.prevent="submitFeedback">
        <fieldset :disabled="isProcessing">
          <div class="form-control">
            <label for="email" class="label">E-mail</label>
            <input type="email" name="Email" id="email" class="input input-bordered" required />
          </div>
          <div class="form-control">
            <label for="content" class="label">Content</label>
            <textarea name="Content" id="content" class="textarea textarea-bordered"></textarea>
          </div>
        </fieldset>

        <div class="flex space-x-2 w-full items-end mt-8">
          <button type="submit" @click="open = false" class="btn ml-auto">Cancel</button>
          <button type="submit" class="btn btn-success" :disabled="isProcessing">Submit</button>
        </div>
      </form>

      <div v-else class="flex flex-col justify-center items-center text-center py-20">
        <h2 class="text-4xl font-bold mb-4">Thank You!</h2>
        <p class="text-lg lg:px-8">Your feedback has been submitted.</p>
      </div>
    </modal>
  </portal>
</template>

<script lang="ts">
import { expandReportApiEndpoint, headers } from '../report-api';
import { catchAndNotifyError } from '../notify';
import Modal from './Modal.vue';
import Portal from './Portal.vue';

export default {
  components: { Modal, Portal },
  emits: ['success', 'error'],
  props: {
    email: {
      type: String,
    },
    messageId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      open: false,
      isProcessing: true,
      isSuccess: false
    }
  },
  watch: {
    open(newValue, oldValue) {
      if (!newValue) return;
      setTimeout(() => {
        if (this.email) {
          (<HTMLInputElement>document.querySelector('input[name="Email"]')).value = this.email;
        }
      }, 300);
    }
  },
  methods: {
    openDialog() { 
      this.open = true; 
    },

    async submitFeedback(e: SubmitEvent) {
      try {
        this.isProcessing = true;
        if (!e.target || !(e.target instanceof HTMLFormElement)) return;
        const formData = new FormData(e.target);
        const payload: Record<string, any> = {};

        formData.forEach((v, k) => {
          payload[k] = v;
        });

        const resp = await fetch(expandReportApiEndpoint('/Feedbacks'), {
          method: 'POST',
          headers,
          body: JSON.stringify(payload)
        });

        if (!resp.ok) {
          throw new Error('Unable to submit your feedback.');
        }

        this.$notify({ type: 'success', text: 'Feedback has been submitted successfully.' });
        this.isSuccess = true;

        setTimeout(() => {
          this.open = false;
          this.$emit('success');
        }, 3000);
      } catch (e) {
        catchAndNotifyError(this, e);
        this.$emit('error', e);
      } finally {
        this.isProcessing = false;
      }
    }
  }
}
</script>