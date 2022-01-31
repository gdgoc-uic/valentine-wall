<template>
  <modal title="Share" :open="open" @update:open="$emit('update:open', false)" with-closing-button>
    <div class="flex flex-col">
      <img :src="imageUrl" class="rounded-2xl" alt="Image" />

      <div class="form-control my-4">
        <label class="label">
          <span class="label-text">URL / Permalink</span>
        </label>
        <input @click="copyURL" ref="permalinkTextBox" type="text" placeholder="URL" :value="permalink" class="w-full input input-bordered" readonly>
      </div>

      <div class="flex flex-row space-x-2">
        <button @click="copyURL" class="btn btn-primary flex-1">
          <icon-copy class="text-lg" />
          <span class="ml-2">{{ hasLinkCopied ? 'Copied!' : 'Copy Link' }}</span>
        </button>
        <a :href="imageUrl" target="_blank" download="image.png" class="btn btn-success flex-1 space-x-2">
          <icon-image-download class="text-lg" />
          <span>Download Image</span>
        </a>
      </div>

      <div class="flex flex-col mt-4">
        <p class="text-sm">Share via </p>

        <div class="flex flex-row space-x-4 mt-2 h-16">
          <button class="btn btn-circle flex-1 h-full"><icon-facebook class="text-lg" /></button>
          <button class="btn btn-circle flex-1 h-full"><icon-messenger class="text-lg" /></button>
          <button class="btn btn-circle flex-1 h-full"><icon-twitter class="text-lg" /></button>
          <button class="btn btn-circle flex-1 h-full"><icon-telegram class="text-lg" /></button>
        </div>
      </div>
    </div>
  </modal>
</template>

<script lang="ts">
import IconLink from '~icons/uil/link';
import IconFacebook from '~icons/uil/facebook-f';
import IconTwitter from '~icons/uil/twitter';
import IconTelegram from '~icons/uil/telegram-alt';
import IconMessenger from '~icons/uil/facebook-messenger';
import IconImageDownload from '~icons/uil/image-download';
import IconCopy from '~icons/uil/copy';

import Modal from './Modal.vue';
import { logEvent } from '@firebase/analytics';
import { analytics } from '../firebase';

export default {
  emits: ['update:open'],
  props: {
    open: {
      type: Boolean
    },
    recipientId: {
      type: String,
      required: true
    },
    messageId: {
      type: String,
      required: true
    },
    permalink: {
      type: String
    }
  },
  components: { 
    Modal,
    IconLink,
    IconFacebook,
    IconTwitter,
    IconTelegram,
    IconMessenger,
    IconImageDownload,
    IconCopy
  },
  computed: {
    imageUrl() {
      return `${import.meta.env.VITE_BACKEND_URL}/messages/${this.recipientId}/${this.messageId}?image`
    }
  },
  data() {
    return {
      hasLinkCopied: false
    }
  },
  methods: {
    copyURL() {
      const permalinkTextbox = this.$refs.permalinkTextBox as HTMLInputElement;
      this.hasLinkCopied = true;
      navigator.clipboard.writeText(permalinkTextbox.value);
      logEvent(analytics!, 'share', { method: 'copy-url', item_id: this.messageId });
      setTimeout(() => {
        this.hasLinkCopied = false;
      }, 1500);
    }
  }
}
</script>