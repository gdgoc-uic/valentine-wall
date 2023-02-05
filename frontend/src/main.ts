import { createSSRApp } from 'vue'
import { createHead } from '@vueuse/head'

import App from './App.vue'
import { createRouter } from './router'
import { createStore, storeKey } from './store'

export function createApp() {
    const app = createSSRApp(App);
    const head = createHead();
    const router = createRouter();
    const store = createStore();
    app
        .use(head)
        .use(router)
        .use(store, storeKey);
    return { app, router, head, store };
}