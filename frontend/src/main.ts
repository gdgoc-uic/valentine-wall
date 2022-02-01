import { createSSRApp } from 'vue'
import { createHead } from '@vueuse/head'

import App from './App.vue'
import { createRouter } from './router'
import { createStore, storeKey } from './store'
import { installClient } from './client'

export function createApp() {
    const app = createSSRApp(App);
    const head = createHead();
    const router = createRouter();
    const store = createStore();
    app.use(head);
    app.use(router);
    app.use(store, storeKey);
    app.use(installClient(store.getters.apiClient));
    return { app, router, head, store };
}