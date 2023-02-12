
//@ts-ignore
import Notifications from 'notiwind'

import { createApp } from './main'
import { analytics } from './firebase'
import { logEvent, setCurrentScreen } from 'firebase/analytics'
import 'floating-vue/dist/style.css'

const { app, router, authStore } = createApp();
const clientApp = app.use(Notifications);

router.beforeEach((to, from, next) => {
  setCurrentScreen(analytics!, to.meta.pageTitle as string);

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  if (requiresAuth && !authStore.state.isLoggedIn) {
    next('/');
    return;
  }

  logEvent(analytics!, 'page_view', {
    page_title: to.meta.pageTitle as string,
    page_location: window.location.href,
  });
  next();
});

router.isReady().then(() => {
  clientApp.mount('#app');
});
