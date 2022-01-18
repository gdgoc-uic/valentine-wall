<template>
  <div v-if="open" class="modal" :class="{ 'modal-open': open }" @click.self="$emit('update:open', false)">
    <div class="modal-box">
      <h2 v-if="title" class="text-center text-2xl font-bold mb-4">{{ title }}</h2>
      <slot></slot>
      <div v-if="$slots.actions" class="modal-action">
        <slot name="actions"></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  emits: ['update:open'],
  props: {
    open: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      required: false
    }
  }
}
</script>