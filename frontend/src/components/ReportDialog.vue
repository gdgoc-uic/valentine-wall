<template>
  <slot :openDialog="openDialog"></slot>

  <portal>
    <modal title="Report" v-model:open="open" with-closing-button>
      <form v-if="!isSuccess" @submit.prevent="submitReport">
        <fieldset :disabled="isProcessing">
          <div class="form-control">
            <label for="email" class="label">E-mail</label>
            <input type="email" name="Email" id="email" class="input input-bordered" required />
          </div>
          <div class="form-control">
            <label for="message_id" class="label">Message ID</label>
            <input type="text" name="MessageID" :readonly="messageId" class="input input-bordered" id="message_id" required />
          </div>
          <div class="form-control">
            <label for="category_id" class="label">Category</label>
            <select
              required
              id="category_id"
              class="select select-bordered"
              :name="categoryIdKey">
              <option v-if="categories.length === 0">Loading ...</option>
              <option
                :value="c.id" :key="c.id"
                v-for="c in categories">{{ c.Title }}</option>
            </select>
          </div>
          <div class="form-control">
            <label for="additional_details" class="label">Additional Details</label>
            <textarea name="AdditionalDetails" id="additional_details" class="textarea textarea-bordered"></textarea>
          </div>
        </fieldset>

        <div class="flex space-x-2 w-full items-end mt-8">
          <button type="submit" @click="open = false" class="btn ml-auto">Cancel</button>
          <button type="submit" class="btn btn-success" :disabled="isProcessing">Submit</button>
        </div>
      </form>

      <div v-else class="flex flex-col justify-center items-center text-center py-20">
        <h2 class="text-4xl font-bold mb-4">Thank You!</h2>
        <p class="text-lg lg:px-8">Your report has been submitted. Our moderators will look into it shortly.</p>
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
      categories: [],
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

        if (this.messageId) {
          (<HTMLInputElement>document.querySelector('input[name="MessageID"]')).value = this.messageId;
        }
      }, 300);

      this.listReportCategories();
    }
  },
  computed: {
    categoryIdKey() {
      return import.meta.env.VITE_REPORT_API_CATEGORY_ID_KEY;
    }
  },
  methods: {
    async listReportCategories() {
      try {
        this.isProcessing = true;

        if (this.categories.length === 0) {
          const resp = await fetch(expandReportApiEndpoint('/ReportCategories'), { headers });
          if (!resp.ok) {
            throw new Error('Unable to load categories. Please try again later.');
          }
          
          this.categories = await resp.json();
        }

        this.isProcessing = false;
      } catch (e) {
        catchAndNotifyError(this, e);
        this.$emit('error', e);
      }
    },

    openDialog() { this.open = true },

    async submitReport(e: SubmitEvent) {
      try {
        this.isProcessing = true;
        if (!e.target || !(e.target instanceof HTMLFormElement)) return;
        const formData = new FormData(e.target);
        const payload: Record<string, any> = {};

        formData.forEach((v, k) => {
          if (k === import.meta.env.VITE_REPORT_API_CATEGORY_ID_KEY) {
            payload[k] = parseInt(v.toString());
          } else {
            payload[k] = v;
          }
        });

        const resp = await fetch(expandReportApiEndpoint('/Reports'), {
          method: 'POST',
          headers,
          body: JSON.stringify(payload)
        });

        if (!resp.ok) {
          throw new Error('Unable to submit your report.');
        }

        this.$notify({ type: 'success', text: 'Report has been submitted successfully.' });
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