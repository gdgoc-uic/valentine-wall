<template>
  <masonry class="message-results -mx-2">
    <!-- TODO: make card widths and --colors-- different -->
    <div :key="msg.id" v-for="msg in messageList" class="w-1/2 md:w-1/3 lg:w-1/4 message-paper-wrapper">
      <router-link
        :to="{ name: 'message-page', params: { recipientId: msg.recipient_id, messageId: msg.id } }"
        class="message-paper"
        :class="[msg.has_gifts ? 'paper-variant-gift' : `paper-variant-${msg.paperColor}`]">
          <div v-if="!msg.has_gifts" class="px-6 pt-6">
            <p>{{ msg.content }}</p>
          </div>
          <div v-else class="flex w-full h-full items-center justify-center">
            <icon-gift class="text-white text-9xl" />
          </div>
          <div class="message-meta-info">
            <p :class="[msg.has_gifts ? 'text-white' : 'text-gray-500']" class="text-sm">{{ humanizeTime(msg.created_at) }} ago</p>
            <icon-reply v-if="msg.has_replied" class="text-pink-500" />
          </div>
      </router-link>
    </div>
  </masonry>
</template>

<script lang="ts">
import Masonry from './Masonry.vue';
import IconReply from '~icons/uil/comment-heart';
import IconGift from '~icons/uil/gift';
import { toNow } from '../time_utils';

export default {
  components: { 
    Masonry,
    IconReply,
    IconGift,
  },
  props: {
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
      default :[]
    }
  },
  mounted() {
    this.applyChanges(this.processResults(this.messages));
  },
  data() {
    return {
      messageList: [] as any[]
    }
  },
  watch: {
    messages(newMessages: any[], oldMessages: any[]) {
      this.applyChanges(newMessages);
    }
  },
  methods: {
    applyChanges(newMessages: any[]) {
      if (!newMessages) return;

      if (this.replace) {
        this.messageList = this.processResults(newMessages);
        return;
      } 

      if (this.limit && (this.messageList.length + newMessages.length) > this.limit) {
        // limit list to (limit - nm.len)
        this.messageList = this.messageList.splice(0, this.limit - newMessages.length);
      }
      
      if (this.prepend) {
        this.messageList.unshift(...this.processResults(newMessages));
      } else {
        this.messageList.push(...this.processResults(newMessages));
      }
    },
    processResults(data: any[]): any[] {
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
          paperColorIds = this.shuffle(paperColorIds);
        }
      } else {
        paperColorIds = this.shuffle(availablePaperColorId.slice());
      }

      // console.log(paperColorIds);
      return data.map((d, i) => {
        const paperColor = paperColorIds[i];
        return {
          ...d, 
          paperColor
        }
      });
    },
    shuffle(array: any[]) {
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
    },
    humanizeTime(date: Date): string {
      return toNow(date);
    }
  }
}
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