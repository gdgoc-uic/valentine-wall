<template>
    <!-- TODO: handle disappear-on-loading -->
    <slot v-if="(!disappearOnLoading || query.isFetched.value) && typeof query.data.value !== 'undefined'" :data="query.data"></slot>
    <slot v-if="query.isError.value" name="error" :error="query.error.value"></slot>
    <div v-else-if="query.isLoading.value" class="py-12 w-full h-full flex-col items-center justify-center text-center">
        <loading class="mx-auto" />
        <p class="mt-4 font-bold text-gray-700">Loading...</p>
    </div>
</template>

<script lang="ts" setup>
import { UseQueryReturnType } from '@tanstack/vue-query';
import { PropType, provide } from 'vue';
import { catchAndNotifyError } from '../notify';
import Loading from './Loading.vue';

defineEmits(['success', 'error']);

// TODO: read async data in SSR
const props = defineProps({
    query: {
        type: Object as PropType<UseQueryReturnType<unknown, unknown>>,
        required: true
    },
    disappearOnLoading: {
        type: Boolean,
        default: false
    }
});

provide('isLoading', props.query.isLoading);
</script>