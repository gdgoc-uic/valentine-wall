<template>
  <div class="h-[70vh] flex flex-col items-center justify-center">
    <loading class="mx-auto" />
    <p class="mt-4 font-bold text-gray-700">Verifying invitation...</p>
  </div>
</template>

<script lang="ts">
import Loading from '../components/Loading.vue';
import { catchAndNotifyError } from '../notify';

export default {
  components: { Loading },
  mounted() {
    this.processInvitationCode(this.$route.params.invitationCode as string);
  },
  methods: {
    async processInvitationCode(invCode: string) {
      setTimeout(async () => {
        try {
          // await this.$client.get(`/invite/${invCode}`, { credentials: 'include' });
        } catch (e) {
          // catchAndNotifyError(this, e);
        } finally {
          this.$router.replace({ name: 'home-page' });
        }
      }, 1500);
    }
  }
}
</script>Loading