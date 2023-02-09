<template>
    <slot :is-loading="isLoading">
        <button
            v-if="shouldGoNext && !isLoading"
            @click="$emit('click')"
            class="mt-8 btn text-gray-900 px-12 self-center bg-white hover:bg-gray-100 border-gray-300 hover:border-gray-500">
            Load More
        </button>
    </slot>
</template>

<script lang="ts" setup>
import { inject, Ref, ref } from 'vue';

const isLoading = inject<Ref<boolean>>('isLoading', ref(false));

defineEmits(['click']);

defineProps({
    link: {
        type: String
    },
    shouldGoNext: {
        type: Boolean,
        default: true
    }
});
</script>