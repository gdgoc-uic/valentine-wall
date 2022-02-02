<template>
    <slot :is-loading="isLoading.value">
        <button
            v-if="link"
            :disabled="isLoading.value"
            @click="$emit('click')"
            class="mt-8 btn text-gray-900 px-12 self-center bg-white hover:bg-gray-100 border-gray-300 hover:border-gray-500">
            Load More
        </button>
    </slot>
</template>

<script>
export default {
    inject: ['isLoading'],
    emits: ['click'],
    props: {
        link: {
            type: String
        }
    }
}
</script>