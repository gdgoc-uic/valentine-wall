<template>
  <modal title="Set an ID" :open="$store.state.isIDModalOpen" @update:open="commit('SET_ID_MODAL_OPEN', $event)">
    <form @submit.prevent="submitIDForm" class="flex flex-col">
      <div class="form-control" name="associated_id_field">
        <label class="label">
          <span class="label-text">Recipient student ID</span>
        </label>
        <input class="input input-bordered" type="text" name="associated_id" pattern="[0-9]{12}" placeholder="12-digit student ID">
      </div>
      <button class="self-end px-12 btn bg-rose-500 hover:bg-rose-600 border-none mt-4" type="submit">Send</button>
    </form>
  </modal>
</template>

<script lang="ts">
import { catchAndNotifyError, notify } from '../notify';
import Modal from './Modal.vue';

export default {
  components: { Modal },
  methods: {
    async submitIDForm(e: SubmitEvent) {
      if (!e.target || !(e.target instanceof HTMLFormElement)) return;
      const associatedIdField = e.target.children.namedItem('associated_id_field');
      if (!associatedIdField) return;
      const associatedIdTextField = associatedIdField.children.namedItem('associated_id');
      if (!associatedIdTextField || !(associatedIdTextField instanceof HTMLInputElement)) return;

      try {
        const resp = await fetch(import.meta.env.VITE_BACKEND_URL + '/user/connect_id', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...this.$store.getters.headers
          },
          body: JSON.stringify({
            associated_id: associatedIdTextField.value,
          })
        });

        const json = await resp.json();
        if (resp.status >= 200 && resp.status <= 299) {
          notify(this, { type: 'success', text: json['message'] });
          this.$store.commit('SET_USER_ASSOCIATED_ID', associatedIdTextField.value);
          this.$store.commit('SET_ID_MODAL_OPEN', false);
        } else if ('error_message' in json) {
          throw new Error(json['error_message']);
        } else {
          throw new Error('There are errors in your submission');
        }
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>