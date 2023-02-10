<template>
  <masonry class="message-results -mx-2">
    <div :key="msg.id" v-for="msg in messageList" :class="boxClass" class="message-paper-wrapper">
      <router-link
        :to="{ name: 'message-page', params: { recipientId: msg.recipient, messageId: msg.id } }"
        class="message-paper"
        :class="[msg.gifts.length != 0 ? 'paper-variant-gift' : `paper-variant-${msg.paperColor}`]">
          <div v-if="msg.gifts.length == 0" class="px-6 pt-6">
            <p>{{ msg.content }}</p>
          </div>
          <div v-else class="flex w-full h-full items-center justify-center">
            <icon-gift class="text-white text-9xl" />
          </div>
          <div class="message-meta-info">
            <p :class="[msg.gifts.length != 0 ? 'text-white' : 'text-gray-500']" class="text-sm">{{ toNow(msg.created) }} ago</p>
            <icon-reply v-if="msg.replies_count > 0" class="text-pink-500" />
          </div>
      </router-link>
    </div>
  </masonry>
</template>

<script lang="ts" setup>
import Masonry from './Masonry.vue';
import IconReply from '~icons/uil/comment-heart';
import IconGift from '~icons/uil/gift';
import { toNow } from '../time_utils';
import { watch, ref } from 'vue';

const props = defineProps({
  prepend: {
    type: Boolean,
    default: false
  },
  replace: {
    type: Boolean,
    default: false
  },
  limit: {
    type: Number
  },
  messages: {
    type: Array,
    default: []
  },
  boxClass: {
    type: String,
    default: 'w-1/2 md:w-1/3 lg:w-1/4'
  }
});

const messageList = ref<any[]>([]);

function applyChanges(newMessages: any[]) {
  if (props.replace) {
    messageList.value = processResults(newMessages);
    return;
  } else if (props.limit && (messageList.value.length + newMessages.length) > props.limit) {
    // limit list to (limit - nm.len)
    messageList.value = messageList.value.splice(0, props.limit - newMessages.length);
  } else if (props.prepend) {
    messageList.value.unshift(...processResults(newMessages));
  } else {
    messageList.value.push(...processResults(newMessages));
  }
}

function processResults(data: any[]): any[] {
  if (data.length === 0) {
    return data;
  }

  const availablePaperColorId = [1,2,3,4];
  const quo = data.length / availablePaperColorId.length;
  let paperColorIds: number[] = [];
  let j = 0;
  let times = Math.floor(quo);
  if (times < 1) times = Math.ceil(quo);

  if (data.length > 1) {
    for (let i = 0; i < data.length; i++) {
      paperColorIds.push(availablePaperColorId[j % availablePaperColorId.length]);
      if (i != 0 && i % times == 0) {
        j++;
      }
    }

    // add a check to avoid repetitions
    for (let i = 0; i < times; i++) {
      paperColorIds = shuffle(paperColorIds);
    }
  } else {
    paperColorIds = shuffle(availablePaperColorId.slice());
  }

  // console.log(paperColorIds);
  return data.map((d, i) => {
    const paperColor = paperColorIds[i];
    return {
      ...d, 
      paperColor
    }
  });
}

function shuffle(array: any[]) {
  var i = array.length,
      j = 0,
      temp;
  while (i--) {
    j = Math.floor(Math.random() * (i+1));
    // swap randomly chosen element with current element
    temp = array[i];
    array[i] = array[j];
    array[j] = temp;
  }
  return array;
}

watch(props.messages, (newMessages, oldMessages) => {
  applyChanges(newMessages);
}, {
  immediate: true
});
</script>

<style lang="postcss">
.message-results > .message-paper-wrapper {
  @apply p-2 min-h-16;
}

.message-paper-wrapper > .message-paper {
  @apply rounded-lg flex flex-col justify-between shadow-lg h-64 hover:scale-110 transition-transform;
  background-size: 100% 22px;
}

.message-paper .message-meta-info {
  @apply  px-6 pb-6 flex justify-between items-center mt-4 rounded-b-lg;
}

.message-paper.paper-variant-gift {
  @apply bg-rose-400;
}

.message-paper.paper-variant-1 {
  background-image: linear-gradient(to bottom,rgb(254, 243, 199) 21px,#00b0d7 1px); 
}

.message-paper.paper-variant-1 .message-meta-info {
  background-color: rgb(254, 243, 199);
}

.message-paper.paper-variant-2 {
  background-image: linear-gradient(to bottom,rgb(152, 221, 255) 21px,#213381 1px); 
}

.message-paper.paper-variant-2 .message-meta-info {
  background-color: rgb(152, 221, 255);
}

.message-paper.paper-variant-3 {
  background-image: linear-gradient(to bottom,rgb(155, 255, 183) 21px,#00b0d7 1px); 
}

.message-paper.paper-variant-3 .message-meta-info {
  background-color: rgb(155, 255, 183);
}

.message-paper.paper-variant-4 {
  background-image: linear-gradient(to bottom,rgb(255, 194, 175) 21px,#213381 1px); 
}

.message-paper.paper-variant-4 .message-meta-info {
  background-color: rgb(255, 194, 175);
}
</style>