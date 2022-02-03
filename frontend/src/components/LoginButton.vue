<template>
    <button v-if="!$store.getters.isLoggedIn" @click="login" class="btn text-gray-500 bg-white hover:bg-gray-100 hover:text-gray-500 normal-case border-0 shadow-md space-x-3">
        <google-icon />
        <span>Sign in with UIC Google</span>
    </button>
</template>

<script>
import GoogleIcon from '~icons/logos/google-icon';

export default {
    emits: ['click'],
    components: {
        GoogleIcon
    },
    methods: {
        async login() {
            try {
                await this.$store.dispatch('login');
                this.$emit('click');
            } catch (e) {
                catchAndNotifyError(this, e);
            }
        },
    }
}
</script>