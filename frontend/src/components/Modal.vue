<template>
  <div v-if="open" class="modal" :class="{ 'modal-open': open }" @click.self="$emit('update:open', false)">
    <div class="modal-box">
      <div class="relative">
        <h2 v-if="title" class="text-center text-2xl font-bold mb-4">{{ title }}</h2>
        <button v-if="withClosingButton" @click="$emit('update:open', false)" class="bg-gray-200 hover:bg-gray-300 transition-colors text-gray-900 p-2 rounded-full absolute top-0 right-0">
          <icon-close />
        </button>
      </div>
      <slot></slot>
      <div v-if="$slots.actions" class="modal-action">
        <slot name="actions"></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import IconClose from '~icons/uil/multiply';

export default {
  emits: ['update:open'],
  components: { IconClose },
  props: {
    open: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      required: false
    },
    withClosingButton: {
      type: Boolean
    }
  }
}
</script>