<template>
  <promise-loader 
    :disappear-on-loading="disappearOnLoading"
    :promise="fetchPromise" 
    :resolve-delay="800"
    :fail-fn="failFn" 
    @resolve="handleResolve" 
    @reject="handleReject">
    <template #default>
        <slot :response="resp"></slot>
    </template>
    <template #pending>
        <div class="py-12 w-full h-full flex-col items-center justify-center text-center">
            <loading class="mx-auto" />
            <p class="mt-4 font-bold text-gray-700">Loading...</p>
        </div>
    </template>
    <template #error>
        <slot 
            name="error" 
            :error="error" 
            :isResponseError="isResponseError()"></slot>
    </template>
  </promise-loader>
</template>

<script lang="ts">
import { catchAndNotifyError } from '../notify';
import Loading from './Loading.vue';
import PromiseLoader from './PromiseLoader.vue';

// TODO: read async data in SSR
export default {
    emits: ['success', 'error'],
    components: {
        PromiseLoader,
        Loading
    },
    props: {
        disappearOnLoading: {
            type: Boolean,
            default: false
        },
        endpoint: {
            type: String,
            required: true
        },
        failFn: {
            type: Function
        }
    },
    data() {
        return {
            resp: { data: null } as unknown,
            error: null as unknown
        }
    },
    computed: {
        fetchPromise() {
            // return !this.endpoint
            //     ? 
                return Promise.reject(new Error('Not loaded.'))
                // : this.$client.get(this.endpoint);
        }
    },
    methods: {
        handleResolve(r: any) {
            this.resp = r;
            this.$emit('success', r);
        },
        handleReject(err: unknown) {
            this.error = err;
            this.$emit('error', err);
            // if (err instanceof APIResponseError && err.rawResponse.status !== 404) {
            //     catchAndNotifyError(this, err);
            // }
        },
        isResponseError(): boolean {
            // return this.error instanceof APIResponseError;
            return true;
        }
    }
}
</script>