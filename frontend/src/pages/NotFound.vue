<template>
  <div class="max-w-3xl w-full mx-auto pt-4 flex flex-col space-y-4 self-start">
    <div class="w-full bg-white p-14 rounded-lg flex flex-col items-center text-center shadow-lg">
      <icon-confused class="text-gray-500 text-9xl mb-2" />
      <h2 class="text-4xl font-bold mb-4">Page not found.</h2>
      <p class="text-xl text-gray-500">
        Double-check if your link is correct and try again.
      </p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import IconConfused from '~icons/uil/confused';
</script>