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

<script lang="ts" setup>
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
import { computed, ref } from 'vue';

const emit = defineEmits(['success', 'error']);
const props = defineProps({
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
});

const buttonList = computed(() => [
  ...(props.imageUrl ? [{
    type: 'link',
    label: 'Download Image',
    attrs: {
      href: props.imageUrl,
      target: '_blank',
      download: props.imageFileName
    },
    icon: IconImageDownload
  }] : []),
  {
    type: 'button',
    label: 'Copy URL',
    attrs: {
      onClick: copyURL
    },
    icon: IconCopy
  },
  {
    type: 'button',
    label: 'Facebook',
    attrs: {
      onClick: () => shareTo('facebook')
    },
    icon: IconFacebook
  },
  {
    type: 'button',
    label: 'Twitter',
    attrs: {
      onClick: () => shareTo('twitter')
    },
    icon: IconTwitter
  },
  {
    type: 'button',
    label: 'Telegram',
    attrs: {
      onClick: () => shareTo('telegram')
    },
    icon: IconTelegram
  }
]);

const permalinkTextBox = ref<HTMLInputElement | null>(null);
const hasLinkCopied = ref(false);
const open = ref(false);

function openDialog() {
  if (import.meta.env.SSR) return;
  open.value = true;
}

function shareTo(provider: string) {
  let android = navigator.userAgent.match(/Android/i);
  let ios = navigator.userAgent.match(/iPhone|iPad|iPod/i);
  const isDesktop = !(ios || android);

  let url: string;
  let windowTitle: string;
  
  switch (provider) {
    case 'facebook':
      url = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(props.permalink)}&quote=${encodeURIComponent(props.text ?? '')}`;
      windowTitle = 'Share to Facebook';
      break;
    case 'twitter':
      url = `https://twitter.com/intent/tweet?url=${encodeURIComponent(props.permalink)}&text=${props.text ?? ''}&hashtags=${props.hashtags.join('')}`;
      windowTitle = 'Share to Twitter';
      break;
    case 'telegram':
      if (isDesktop) {
        url = `https://telegram.me/share/msg?url=${props.permalink}&text=${props.text ?? ''}`;
      } else {
        url = `tg://msg?text=${((props.text + ' ') ?? '') + props.permalink}`;
      }
      windowTitle = 'Share to Telegram';
      break;
    default:
      emit('error', new Error(`unknown provider '${provider}'`));
      return;
  }

  const shareWindow = popupCenter({
    url: url!,
    title: windowTitle!,
    w: 800,
    h: 500
  });

  const onClose = function() {
    emit('success', provider);
    shareWindow?.removeEventListener('close', onClose);
  };

  shareWindow?.addEventListener('close', onClose);
}

function copyURL() {
  hasLinkCopied.value = true;
  navigator.clipboard.writeText(permalinkTextBox.value!.value);
  emit('success', 'clipboard');
  setTimeout(() => {
    hasLinkCopied.value = false;
  }, 1500);
}
</script>