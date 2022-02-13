<template>
  <slot :openDialog="openDialog"></slot>

  <portal>
    <modal title="Share" v-model:open="open" modal-box-class="h-auto" with-closing-button>
      <div class="flex flex-col">
        <img v-if="imageUrl" :src="imageUrl" class="rounded-2xl" alt="Image" />
        <div class="form-control my-4">
          <label class="label">
            <span class="label-text">URL / Permalink</span>
          </label>
          <input @click="copyURL" ref="permalinkTextBox" type="text" placeholder="URL" :value="permalink" class="w-full input input-bordered" readonly>
        </div>
        <div class="flex flex-col">
          <div class="flex flex-row flex-wrap -mx-4">
            <div 
              :key="'_share_button_' + bi"
              v-for="(btn, bi) in buttonList"
              class="w-1/5 flex flex-col items-center p-2 space-y-2">
              <a 
                v-if="btn.type == 'link'"
                v-bind="btn.attrs"
                class="btn btn-circle border-0 bg-gray-200 hover:bg-gray-300 text-gray-700 w-16 h-16">
                <component :is="btn.icon" class="text-2xl" />
              </a>
              <button 
                v-else-if="btn.type == 'button'"
                v-bind="btn.attrs"
                class="btn btn-circle border-0 bg-gray-200 hover:bg-gray-300 text-gray-700 w-16 h-16">
                <component :is="btn.icon" class="text-2xl" />
              </button>

              <span class="text-sm text-center text-gray-500">{{ btn.label }}</span>
            </div>
          </div>
        </div>
      </div>
    </modal>
  </portal>
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
import Portal from './Portal.vue';
import { popupCenter } from '../auth';

export default {
  emits: ['success', 'error'],
  props: {
    title: {
      type: String
    },
    text: {
      type: String
    },
    imageUrl: {
      type: String
    },
    imageFileName: {
      type: String,
      default: 'image.png'
    },
    permalink: {
      type: String,
      required: true
    },
    hashtags: {
      type: Array,
      default: []
    },
  },
  components: { 
    Modal,
    IconLink,
    IconFacebook,
    IconTwitter,
    IconTelegram,
    IconMessenger,
    IconImageDownload,
    IconCopy,
    Portal
  },
  computed: {
    buttonList() {
      return [
        ...(this.imageUrl ? [{
          type: 'link',
          label: 'Download Image',
          attrs: {
            href: this.imageUrl,
            target: '_blank',
            download: this.imageFileName
          },
          icon: IconImageDownload
        }] : []),
        {
          type: 'button',
          label: 'Copy URL',
          attrs: {
            onClick: this.copyURL
          },
          icon: IconCopy
        },
        {
          type: 'button',
          label: 'Facebook',
          attrs: {
            onClick: () => this.shareTo('facebook')
          },
          icon: IconFacebook
        },
        {
          type: 'button',
          label: 'Twitter',
          attrs: {
            onClick: () => this.shareTo('twitter')
          },
          icon: IconTwitter
        },
        {
          type: 'button',
          label: 'Telegram',
          attrs: {
            onClick: () => this.shareTo('telegram')
          },
          icon: IconTelegram
        }
      ];
    }
  },
  data() {
    return {
      hasLinkCopied: false,
      open: false
    }
  },
  methods: {
    openDialog() {
      if (import.meta.env.SSR) return;

      if (!!navigator.share) {
        navigator.share({
          title: this.title,
          url: this.permalink
        })
        .then(() => {
          this.$emit('success');
        })
        .catch(err => {
          this.$emit('error', err);
        });
      } else {
        this.open = true;
      }
    },
    shareTo(provider: string) {
      let android = navigator.userAgent.match(/Android/i);
      let ios = navigator.userAgent.match(/iPhone|iPad|iPod/i);
      const isDesktop = !(ios || android);

      let url: string;
      let windowTitle: string;
      
      switch (provider) {
        case 'facebook':
          url = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(this.permalink)}&quote=${encodeURIComponent(this.text ?? '')}`;
          windowTitle = 'Share to Facebook';
          break;
        case 'twitter':
          url = `https://twitter.com/intent/tweet?url=${encodeURIComponent(this.permalink)}&text=${this.text ?? ''}&hashtags=${this.hashtags.join('')}`;
          windowTitle = 'Share to Twitter';
          break;
        case 'telegram':
          if (isDesktop) {
            url = `https://telegram.me/share/msg?url=${this.permalink}&text=${this.text ?? ''}`;
          } else {
            url = `tg://msg?text=${((this.text + ' ') ?? '') + this.permalink}`;
          }
          windowTitle = 'Share to Telegram';
          break;
        default:
          this.$emit('error', new Error(`unknown provider '${provider}'`));
          return;
      }

      const shareWindow = popupCenter({
        url: url!,
        title: windowTitle!,
        w: 800,
        h: 500
      });

      const vm = this;
      const onClose = function() {
        vm.$emit('success', provider);
        shareWindow?.removeEventListener('close', onClose);
      };

      shareWindow?.addEventListener('close', onClose);
    },
    copyURL() {
      const permalinkTextbox = this.$refs.permalinkTextBox as HTMLInputElement;
      this.hasLinkCopied = true;
      navigator.clipboard.writeText(permalinkTextbox.value);
      this.$emit('success', 'clipboard');
      setTimeout(() => {
        this.hasLinkCopied = false;
      }, 1500);
    }
  }
}
</script>