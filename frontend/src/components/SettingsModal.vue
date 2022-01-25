<template>
  <modal title="Settings" :open="open" @update:open="$emit('update:open', $event)" with-closing-button>
    <!-- TODO: -->
    <form @submit.prevent="saveSettings" class="flex flex-col">
      <!-- Info -->
      <p class="text-center my-4 text-2xl text-gray-400">Coming soon!</p>

      <!-- Posting method -->

      <!-- Delete account -->
    </form>
  </modal>
</template>

<script>
import Modal from './Modal.vue';
export default {
  components: { Modal },
  emits: ['update:open'],
  props: {
    open: {
      type: Boolean
    }
  },
  methods: {
    saveSettings() {
      
    }
  }
}
</script>