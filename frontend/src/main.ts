import { createApp } from 'vue'
//@ts-ignore
import Notifications from 'notiwind'

import App from './App.vue'
import router from './router'
import store from './store'
import './assets/index.css'

createApp(App)
  .use(Notifications)
  .use(router)
  .use(store)
  .mount('#app')
