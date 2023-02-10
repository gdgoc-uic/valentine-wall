<template>
  <slot :openDialog="openDialog"></slot>
  
  <portal to="body">
    <modal title="Archive" v-model:open="open" with-closing-button>
      <div class="max-h-[40rem] h-[40rem]">
        <div v-if="status === 'done'" class="flex flex-col items-center justify-center h-full w-full">
          <h1 class="text-2xl font-bold">Archive success!</h1>
          <p class="text-gray-600 text-xl">{{ total }} messages were archived.</p>
          <button @click="downloadFile(endpoint)" class="mt-8 btn btn-primary w-2/3">Download</button>
        </div>
        <div v-else-if="status !== 'error'" class="flex flex-col space-y-2 items-center justify-center h-full w-full">
          <loading />
          <p class="text-xl font-bold">{{ statusText }}</p>
        </div>
        <div v-else class="flex flex-col text-center space-y-2 items-center justify-center h-full w-full">
          <h1 class="text-2xl font-bold">Something went wrong.</h1>
          <p class="text-gray-600 text-xl">{{ error }}</p>
        </div>
      </div>
    </modal>
  </portal>
</template>

<script lang="ts">
import { EventSourcePolyfill } from 'event-source-polyfill';
import Loading from './Loading.vue';
import Modal from './Modal.vue';
import Portal from './Portal.vue'
export default {
  components: { Portal, Modal, Loading },
  emits: ['success', 'error'],
  beforeUnmount() {
    if (this.sse)
      this.sse.close();
  },
  data() {
    return {
      open: false,
      processed: 0,
      total: 0,
      error: null as unknown as string | null,
      status: 'processing',
      endpoint: '',
      sse: null as unknown as EventSourcePolyfill | null
    }
  },
  watch: {
    open(newValue, oldValue) {
      if (!newValue) {
        this.error = null;
        this.total = 0;
        this.processed = 0;
        this.sse = null;
        this.endpoint = '';
        return;
      }
      this.initSse();
    }
  },
  methods: {
    openDialog() { 
      this.open = true; 
    },
    downloadFile(url: string) {
      const anchor = document.createElement('a');
      anchor.href = url;
      anchor.click();
    },
    initSse() {
      // TODO:
      // this.sse = new EventSourcePolyfill(expandAPIEndpoint('/user/messages/archive'), {
      //   headers: this.$store.getters.headers
      // });

      // this.sse.onmessage = (event) => {
      //   if (event.data === 'null') return;
      //   const { status, data } = JSON.parse(event.data);
      //   this.status = status;

      //   switch (status) {
      //     case 'starting':
      //       return;
      //     case 'processing':
      //       this.processed += data['len'];
      //       return;
      //     case 'set_file_count':
      //       this.total = data['count'];
      //       return;
      //     case 'error':
      //       this.error = data['message'];
      //       this.$emit('error', new Error(data['message'] ?? 'unknown error'));
      //       return;
      //     case 'done':
      //       this.endpoint = expandAPIEndpoint(data['endpoint'] + '?__auth_token=' + this.$store.state.user.accessToken);
      //       this.sse!.close();
      //       return;
      //   }
      // }

      // this.sse.onerror = (ev) => {
      //   this.status = 'error';
      //   this.error = 'Unknown error';
      //   this.$emit('error', new Error('Unknown error'));
      //   this.sse!.close();
      // }
    }
  },
  computed: {
    statusText(): string {
      switch (this.status) {
        case 'starting':
          return 'Starting...';
        case 'set_file_count':
        case 'processing':
          return `Processing ${this.processed}/${this.total}`;
        default:
          return '';
      }
    }
  }
}
</script>

<style>

</style>