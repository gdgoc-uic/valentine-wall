<template>
  <div v-if="open" class="modal" :class="{ 'modal-open': open }" @click.self="$emit('update:open', false)">
    <slot name="modal-box">
      <div class="modal-box" :class="modalBoxClass">
        <div class="relative">
          <h2 v-if="title" class="text-center text-2xl font-bold mb-4">{{ title }}</h2>
          <button v-if="withClosingButton" @click="$emit('update:open', false)"
                  class="bg-rose-500 hover:bg-rose-600 transition-colors text-white p-2 rounded-full absolute top-0 right-0">
            <icon-close />
          </button>
        </div>
        <slot></slot>
        <div v-if="$slots.actions" class="modal-action">
          <slot name="actions"></slot>
        </div>
      </div>
    </slot>
  </div>
</template>

<script lang="ts">
import IconClose from '~icons/uil/multiply';

export default {
  emits: ['update:open'],
  components: { IconClose },
  props: {
    modalBoxClass: {
      type: String,
    },
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