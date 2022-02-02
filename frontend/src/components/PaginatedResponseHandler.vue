<template>
  <response-handler 
    :fail-fn="failFn"
    :disappear-on-loading="false"
    :endpoint="endpoint"
    @success="handleSuccessResponse" 
    @error="$emit('error', $event)">
    <template #default>
        <slot :data="data" :links="links" :goto="goto"></slot>
    </template>
    <template #error="{ error, isResponseError }">
        <slot name="error" 
            :error="error" 
            :isResponseError="isResponseError"></slot>
    </template>
  </response-handler>
</template>

<script lang="ts">
import { APIResponse } from '../client';
import ResponseHandler from './ResponseHandler.vue'
export default {
    emits: ['success', 'error'],
    components: {
        ResponseHandler
    },
    props: {
        originEndpoint: {
            type: String,
            required: true
        },
        failFn: {
            type: Function
        }
    },
    created() {
        this.endpoint = this.originEndpoint;
    },
    data() {
        return {
            links: { 
                first: null as string | null, 
                last: null as string | null, 
                next: null as string | null, 
                previous: null as string | null
            },
            page: 1,
            perPage: 10,
            pageCount: 1,
            data: [],
            endpoint: '',
            merge: false
        }
    },
    watch: {
        originEndpoint(newVal, oldVal) {
            if (newVal == oldVal) return;
            this.merge = false;
            this.data = [];
        }
    },
    methods: {
        handleSuccessResponse(r: APIResponse) {
            const json = r.data;
            this.links = json['links'];
            this.page = json['page'];
            this.perPage = json['per_page'];
            this.pageCount = json['page_count'];
            this.data = this.merge ? this.data.concat(...json['data']) : json['data'];
            this.$emit('success', r);
        },
        goto(endpoint: string, merge: boolean = false) {
            this.endpoint = endpoint;
            this.merge = merge;
        }
    }
}
</script>