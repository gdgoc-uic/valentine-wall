<template>
  <promise-loader 
    :disappear-on-loading="disappearOnLoading"
    :promise="fetchPromise" 
    @resolve="handleResolve" 
    :fail-fn="failFn" 
    @reject="handleReject">
    <template #default>
        <slot :response="resp"></slot>
    </template>
    <template #pending>
        <div class="w-full h-full">
            <!-- TODO: loading component -->
            <p>Loading...</p>
        </div>
    </template>
    <template #error>
        <slot 
            name="error" 
            :error="err" 
            :isResponseError="isResponseError()"></slot>
    </template>
  </promise-loader>
</template>

<script lang="ts">
import { APIResponse, APIResponseError } from '../client';
import PromiseLoader from './PromiseLoader.vue';

// TODO: read async data in SSR
export default {
    emits: ['success', 'error'],
    components: {
        PromiseLoader
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
            resp: null as unknown as APIResponse,
            error: null as unknown
        }
    },
    computed: {
        fetchPromise() {
            return this.$client.get(this.endpoint);
        }
    },
    methods: {
        handleResolve(r: APIResponse) {
            this.resp = r;
            this.$emit('success', r);
        },
        handleReject(err: unknown) {
            this.error = err;
            this.$emit('error', err);
        },
        isResponseError(): boolean {
            return this.error instanceof APIResponseError;
        }
    }
}
</script>