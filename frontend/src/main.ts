import { createSSRApp } from 'vue'
import { createHead } from '@vueuse/head'

import App from './App.vue'
import { createRouter } from './router'
import { authStore as authStoreKey, createAuthStore, createStore, storeKey } from './store_new'
import { VueQueryPlugin, VueQueryPluginOptions } from '@tanstack/vue-query'

export function createApp() {
    const app = createSSRApp(App);
    const head = createHead();
    const router = createRouter();
    const store = createStore();
    const authStore = createAuthStore();

    app
        .use(head)
        .use(router)
        .use(VueQueryPlugin, {
            queryClientConfig: {
                defaultOptions: {
                    queries: {
                        enabled: !import.meta.env.SSR
                    }
                }
            }
        } as VueQueryPluginOptions)
        .provide(storeKey, store)
        .provide(authStoreKey, authStore);
    return { app, router, head, store, authStore };
}