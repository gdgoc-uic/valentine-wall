<template>
  <p class="text-lg" :class="[charLeft <= 0 ? 'text-red-600' : 'text-gray-500']">{{ charLeft }}</p>
</template>

<script lang="ts">
export default {
  props: {
    max: {
      type: Number,
      default: 240
    },
    newlineCount: {
      type: Number,
      default: 1
    },
    content: {
      type: String,
      default: ''
    },
  },
  computed: {
    charLeft() {
      return this.getCharLeft(this.content);
    },
    shouldSend(): boolean {
      return this.charLeft != this.max && this.charLeft >= 0; 
    }
  },
  methods: {
    getCharLeft(content: string): number {
      const linesMinusOne =  (content.match(/\n/g) ?? '').length;
      const actualContentCount = (content.match(/[^\n]/g) ?? '').length;
      return this.max - ((linesMinusOne * this.newlineCount) + actualContentCount);
    }
  }
}
</script>