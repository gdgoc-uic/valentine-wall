
//@ts-ignore
import Notifications from 'notiwind'

import { createApp } from './main'
import { analytics } from './firebase'
import { logEvent, setCurrentScreen } from 'firebase/analytics'

const { app, router, store } = createApp();
const clientApp = app.use(Notifications);

if (window.__INITIAL_STATE__) {
  store.replaceState(window.__INITIAL_STATE__);
} 

router.beforeEach((to, from, next) => {
  setCurrentScreen(analytics!, to.meta.pageTitle as string);
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  if (requiresAuth && !store.getters.isLoggedIn) {
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
