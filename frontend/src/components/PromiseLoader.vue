<template>
  <slot v-if="!disappearOnLoading || status == 1"></slot>
  <slot v-if="status == 0" name="error"></slot>
  <slot v-else-if="status == 2" name="pending"></slot>
</template>

<script lang="ts">
import { computed } from '@vue/reactivity';
import { catchAndNotifyError } from '../notify';

const STATUS_REJECT = 0;
const STATUS_RESOLVE = 1;
const STATUS_PENDING = 2;

export default {
    emits: ['resolve', 'reject'],
    props: {
        promise: {
            type: Promise,
            required: true
        },
        disappearOnLoading: {
            type: Boolean,
            default: false
        },
        failFn: {
            type: Function
        }
    },
    provide() {
        return {
            isLoading: computed(() => this.status === STATUS_PENDING)
        }
    },
    data() {
        return {
            status: STATUS_PENDING // 2 for loading, 0 for reject, 1 for resolve
        }
    },
    watch: {
        promise: {
            handler() {
                this.status = STATUS_PENDING;
                this.promise.then(d => {
                    this.failFn?.(d);
                    this.status = STATUS_RESOLVE;
                    this.$emit('resolve', d);
                }).catch(err => {
                    this.status = STATUS_REJECT;
                    this.$emit('reject', err);
                    catchAndNotifyError(this, err);
                });
            },
            immediate: true
        }
    }
}
</script>