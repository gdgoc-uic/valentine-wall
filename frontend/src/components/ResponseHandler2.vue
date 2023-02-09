<template>
    <!-- TODO: handle disappear-on-loading -->
    <slot v-if="!disappearOnLoading" :data="query.data"></slot>
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

// TODO: review
// export default {
//     emits: ['success', 'error'],
//     components: {
//         PromiseLoader,
//         Loading
//     },
//     props: {
//         disappearOnLoading: {
//             type: Boolean,
//             default: false
//         },
//         endpoint: {
//             type: String,
//             required: true
//         },
//         failFn: {
//             type: Function
//         }
//     },
//     data() {
//         return {
//             resp: { data: null } as unknown,
//             error: null as unknown
//         }
//     },
//     computed: {
//         fetchPromise() {
//             // return !this.endpoint
//             //     ? 
//                 return Promise.reject(new Error('Not loaded.'))
//                 // : this.$client.get(this.endpoint);
//         }
//     },
//     methods: {
//         handleResolve(r: any) {
//             this.resp = r;
//             this.$emit('success', r);
//         },
//         handleReject(err: unknown) {
//             this.error = err;
//             this.$emit('error', err);
//             // if (err instanceof APIResponseError && err.rawResponse.status !== 404) {
//             //     catchAndNotifyError(this, err);
//             // }
//         },
//         isResponseError(): boolean {
//             // return this.error instanceof APIResponseError;
//             return true;
//         }
//     }
// }
</script>