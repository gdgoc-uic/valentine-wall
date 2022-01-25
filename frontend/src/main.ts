import { createApp } from 'vue'
//@ts-ignore
import Notifications from 'notiwind'
import { createHead } from '@vueuse/head'

import App from './App.vue'
import router from './router'
import store from './store'
import './assets/index.css'

const head = createHead()

createApp(App)
  .use(Notifications)
  .use(head)
  .use(router)
  .use(store)
  .mount('#app')
