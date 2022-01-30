<template>
  <div class="alert rounded-md" :class="{ 
    'bg-orange-50 border-red-500 text-red-600': type == 'error', 
    'bg-blue-50 border-blue-500 text-blue-600': type == 'info', 
    'bg-green-50 border-teal-500 text-teal-600': type == 'success', 
    'bg-amber-50 border-orange-500 text-orange-600': type == 'warning' 
  }">
    <slot>
      <p class="mx-auto text-center">{{ message }}</p>
    </slot>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    type: {
      type: String,
      default: 'error'
    },
    message: {
      type: String,
      default: ''
    }
  }
}
</script>