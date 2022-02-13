<template>
  <slot 
    :openDialog="openDialog" 
    :closeDialog="closeDialog" 
    :open="open"></slot>

  <portal>
    <modal v-model:open="open">
      <slot name="dialog-content">
        <p>Are you sure you want to delete this?</p>
      </slot>
      <template #actions>
        <button class="btn" @click="$emit('confirm', false, closeDialog);">Cancel</button>
        <button class="btn btn-error" @click="$emit('confirm', true, closeDialog);">Delete</button>
      </template>
    </modal>
  </portal>
</template>

<script>
import Modal from './Modal.vue'
import Portal from './Portal.vue'
export default {
  components: { Portal, Modal },
  emits: ['confirm'],
  data() {
    return {
      open: false
    }
  },
  methods: {
    openDialog() {
      this.open = true;
    },
    closeDialog() {
      this.open = false;
    }
  }
}
</script>