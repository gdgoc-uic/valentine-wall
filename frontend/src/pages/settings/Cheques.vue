<template>
  <div class="space-y-3">
    <basic-alert
      type="info"
      message="Cheques are claimable e-certificates that you can use to earn money from events." />
    <div class="flex flex-col divide-y space-y-3">
      <div class="pb-6 space-y-3">
        <div class="space-y-2">
          <h2 class="text-2xl font-bold">Deposit</h2>
          <p>Enter the cheque ID you have received so you can deposit them and be credited to your wallet.</p>
        </div>
        <form @submit.prevent="depositCheck">
          <div class="form-control">
            <div class="flex space-x-2">
              <input placeholder="Cheque ID" name="cheque_id" :disabled="isProcessing" type="text" class="input input-bordered flex-1" />
              <button type="submit" class="btn btn-success" :disabled="isProcessing">Submit</button>
            </div>
          </div>
        </form>
      </div>
      <div class="py-6 space-y-2">
        <h2 class="text-2xl font-bold">History</h2>
        <p>
          To check if your cheque has been deposited to your wallet, you may visit the
          <router-link :to="{ name: 'settings-transactions-section' }" class="link">Transactions</router-link> page.
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import BasicAlert from '../../components/BasicAlert.vue';
import { catchAndNotifyError } from '../../notify';
export default {
  components: { BasicAlert },
  data() {
    return {
      isProcessing: false
    }
  },
  methods: {
    async depositCheck(e: SubmitEvent) {
      try {
        this.isProcessing = true;
        // if (!e.target || !(e.target instanceof HTMLFormElement)) return;
        // const formData = new FormData(e.target);
        // const { data: json } = await this.$client.postJson('/user/cheque/deposit', { 
        //   cheque_id: formData.get('cheque_id')
        // });

        // this.$notify({ type: 'success', text: json['message'] });
        // this.$router.go(0);
      } catch (e) {
        catchAndNotifyError(this, e);
      } finally {
        this.isProcessing = false;
      }
    }
  }
}
</script>