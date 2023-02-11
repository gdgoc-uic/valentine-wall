<template>
    <button 
        v-if="!authState.isLoggedIn" @click="googleLogin" 
        class="btn text-gray-500 bg-white hover:bg-gray-100 hover:text-gray-500 normal-case border border-gray-300 shadow-md space-x-3">
        <google-icon />
        <span>Sign in with UIC Google</span>
    </button>
</template>

<script lang="ts" setup>
import GoogleIcon from '~icons/logos/google-icon';
import { catchAndNotifyError } from '../notify';
import { useAuth, useStore } from '../store_new';

const emit = defineEmits(['click']);
const { state: authState, methods: { login } } = useAuth();
const mainStore = useStore();

async function googleLogin() {
    try {
        await login();
        if (mainStore.state.isSetupModalOpen) {
            // logEvent(analytics!, 'sign_up');
        } else {
            // logEvent(analytics!, 'login');
        }
        emit('click');
    } catch (e) {
        catchAndNotifyError(e);
    }
}
</script>