<template>
  <modal :open="open" @update:open="$emit('update:open', $event)" title="Generate Invitation" with-closing-button>
    <div v-if="gotInvitationCode" class="bg-green-50 border border-green-500 rounded-md p-3">
      <p class="font-bold">Congratulations, a link has been generated!</p>

      <div class="form-control my-4">
        <label class="label">
          <span class="label-text">URL / Permalink</span>
        </label>
        <div class="tooltip tooltip-top" :data-tip="hasLinkCopied ? 'Copied!' : 'Click the link to copy.'">
          <input
            @click="copyURL"
            ref="permalinkTextBox"
            type="text" placeholder="URL"
            :value="invitationLink"
            class="w-full input input-bordered"
            readonly />
        </div>
      </div>

      <p>
        Your invitation link will only last for {{ gotExpirationHrs }} hours 
        and is usable for {{ gotMaxUsers }} successful sign-ups. For every successful invite, 
        you will be given {{ gotRewardCoins }} coins both for you and the user.
      </p>
    </div>

    <form ref="generateInvForm">
      <div class="flex flex-col md:flex-row justify-between items-center md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="max_users">Max Users</label>
        <select :disabled="isProcessing" class="select select-bordered" name="max_users" id="max_users">
          <option :key="i" v-for="i in 5" :value="i">{{ i }}</option>
        </select>
      </div>
    </form>

    <template #actions>
      <button @click="generateInvitationLink" class="btn" :disabled="isProcessing">Generate</button>
    </template>
  </modal>
</template>

<script>
import { catchAndNotifyError } from '../notify';
import Modal from './Modal.vue';

export default {
  emits: ['update:open', 'success'],
  props: {
    open: {
      type: Boolean,
      default: false
    }
  },
  components: {
    Modal,
  },
  data() {
    return {
      isProcessing: false,
      hasLinkCopied: false,
      gotInvitationCode: null,
      gotExpirationHrs: null,
      gotMaxUsers: null,
      gotRewardCoins: null
    }
  },
  computed: {
    invitationLink() {
      if (!this.gotInvitationCode) return '';
      return import.meta.env.VITE_FRONTEND_URL + '/invite/' + this.gotInvitationCode;
    }
  },
  methods: {
    copyURL() {
      const permalinkTextbox = this.$refs.permalinkTextBox;
      if (!permalinkTextbox || !(permalinkTextbox instanceof HTMLInputElement)) return;

      this.hasLinkCopied = true;
      navigator.clipboard.writeText(permalinkTextbox.value);
      setTimeout(() => {
        this.hasLinkCopied = false;
      }, 1500);
    },
    async generateInvitationLink() {
      // try {
      //   this.isProcessing = true;
      //   this.gotInvitationCode = null;
      //   this.gotExpirationHrs = null;
      //   this.gotMaxUsers = null;
      //   this.gotRewardCoins = null;

      //   const form = this.$refs.generateInvForm;
      //   if (!form || !(form instanceof HTMLFormElement)) return;
  
      //   const formData = new FormData(form);
      //   const maxUsers = formData.get('max_users');
      //   const { data: json } = await this.$client.post('/user/invitations/generate?max_users=' + maxUsers);
      //   this.$notify({ type: 'success', text: json['message'] });
      //   this.$emit('success', json);
      //   this.gotInvitationCode = json['invitation_code'];
      //   this.gotExpirationHrs = json['expires_in_hrs'];
      //   this.gotMaxUsers = json['max_users'];
      //   this.gotRewardCoins = json['reward_coins'];
      // } catch (e) {
      //   catchAndNotifyError(this, e);
      // } finally {
      //   this.isProcessing = false;
      // }
    }
  }
}
</script>