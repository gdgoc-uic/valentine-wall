<template>
  <modal title="Setup Account" :open="$store.state.isSetupModalOpen" @update:open="commit('SET_SETUP_MODAL_OPEN', $event)">
    <form ref="assocIdForm" @submit.prevent="shouldProceed" :class="[proceedToTerms ? 'hidden' : 'flex']" class="flex flex-col">
      <div class="form-control">
        <label class="label">
          <span class="label-text">Enter your student ID</span>
        </label>
        <input class="input input-bordered" type="text" name="associated_id" id="associated_id_field" pattern="[0-9]{6,12}" placeholder="6 to 12-digit Student ID (e.g. 200xxxxxxxxx)">
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text">College Deparment</span>
        </label>
        <select name="department" class="select select-bordered">
          <option value="none" selected>None</option>
          <option :value="dept.id" :key="dept.id" v-for="dept in departments">{{ dept.label }} ({{ dept.id }})</option>
        </select>
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text">Gender</span>
        </label>
        <select name="gender" class="select select-bordered">
          <option :value="g.value" :key="g.value" v-for="g in genderList">{{ g.label }}</option>
        </select>
      </div>
      <button class="self-end px-12 btn bg-rose-500 hover:bg-rose-600 border-none mt-4" type="submit">Next</button>
    </form>
    <form @submit.prevent="submitSetupForm" :class="[proceedToTerms ? 'flex' : 'hidden']" class="flex-col">
      <div class="p-4">
        <h2 class="text-center text-2xl mb-5">Terms and Conditions</h2>
        <ul class="list-disc pl-2">
          <li>
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. A, nisi!</p>
          </li>
          <li>
            <p>Eius error sed tempora hic ratione, culpa vitae dicta nostrum?</p>
          </li>
          <li>
            <p>Sed laboriosam unde iste in eum nam doloribus aspernatur delectus.</p>
          </li>
          <li>
            <p>Saepe quo beatae nobis doloremque, odio unde asperiores quaerat ipsa!</p>
          </li>
          <li>
            <p>Odit mollitia beatae dolorum. Neque aliquam dicta nihil iusto eos?</p>
          </li>
        </ul>
        <p class="mt-4">By clicking "Accept", you have agreed to the terms and conditions of this site. Should you violate any of the text above will result to account termination.</p>
      </div>
      <div class="w-full space-x-2 mt-4 flex">
        <button @click.prevent="proceedToTerms = false" class="px-6 btn mr-auto">Go back</button>
        <button @click="termsAgreed = true" class="px-6 btn bg-rose-500 hover:bg-rose-600 border-none" type="submit">Agree</button>
        <button @click="termsAgreed = false" class="btn px-6">Disagree</button>
      </div>
    </form>
  </modal>
</template>

<script lang="ts">
import { APIResponseError } from '../client';
import { catchAndNotifyError, notify } from '../notify';
import Modal from './Modal.vue';

const emailRegex = /^[a-z]+_([0-9]+)@uic.edu.ph$/;

export default {
  components: { Modal },
  data() {
    return {
      termsAgreed: false,
      proceedToTerms: false,
      departments: []
    };
  },
  mounted() {
    setTimeout(() => {
      const associatedIdField = document.getElementById('associated_id_field');
      if (associatedIdField && associatedIdField instanceof HTMLInputElement) {
        const extractedId = this.getIdFromEmail(this.$store.state.user.email);
        associatedIdField.value = extractedId;
      }
    }, 500);
  },
  computed: {
    genderList() {
      return [
        {
          label: 'Male',
          value: 'male'
        },
        {
          label: 'Female',
          value: 'female'
        }
      ];
    }
  },
  methods: {
    getIdFromEmail(input: string): string {
      const matches = emailRegex.exec(input)
      return matches?.[1] ?? '';
    },

    shouldProceed(e: SubmitEvent) {
      if (!e.target || !(e.target instanceof HTMLFormElement)) return;
      const formData = new FormData(e.target);
      if (!formData.get("associated_id")) {
        this.$notify({
          type: 'error',
          text: 'Please input your ID.'
        });
        return;
      } else if (formData.get('associated_id')?.toString() !== this.getIdFromEmail(this.$store.state.user.email)) {
        this.$notify({
          type: 'error',
          text: 'Your ID from e-mail does not match with the one you have inputted.'
        });
        return;
      }

      if (!formData.get('department') || formData.get('department')?.toString() == 'none') {
        this.$notify({
          type: 'error',
          text: 'Please select your department.'
        });
        return;
      }

      this.proceedToTerms = true;
    },

    async submitSetupForm(e: SubmitEvent) {
      if (!e.target || !(e.target instanceof HTMLFormElement)) return;
      const assocIdForm = this.$refs.assocIdForm as HTMLFormElement;
      if (!assocIdForm || !(assocIdForm instanceof HTMLFormElement)) return;
      const formData = new FormData(assocIdForm);

      try {
        try {
          const { data: json } = await this.$client.postJson('/user/setup', {
            associated_id: formData.get('associated_id')?.toString(),
            department: formData.get('department')?.toString(),
            gender: formData.get('gender')?.toString(),
            terms_agreed: this.termsAgreed
          });

          notify(this, { type: 'success', text: json['message'] });
          this.$store.commit('SET_USER_ASSOCIATED_ID', json['associated_id']);
          this.$store.commit('SET_SETUP_MODAL_OPEN', false);
        } catch (e) {
          if (e instanceof APIResponseError && e.rawResponse.status == 403 && e.message == 'Access to the service is denied.') {
            await this.$store.dispatch('logout');
            return;
          }
          throw e;
        }
      } catch(e) {
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>