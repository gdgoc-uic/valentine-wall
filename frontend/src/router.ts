import { logEvent, setCurrentScreen } from 'firebase/analytics';
import { createRouter, createWebHistory } from 'vue-router';
import { analytics } from './firebase';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      name: 'home-page',
      path: '/',
      component: () => import('./pages/Home.vue'),
      meta: {
        pageTitle: 'UIC Valentines Wall'
      }
    },
    {
      name: 'message-wall-page',
      path: '/wall/:recipientId?',
      component: () => import('./pages/Wall.vue'),
      meta: {
        pageTitle: 'Message Wall | UIC Valentines Wall'
      }
    },
    {
      name: 'message-page',
      path: '/wall/:recipientId/:messageId',
      component: () => import('./pages/Message.vue'),
      meta: {
        pageTitle: 'Message Page | UIC Valentines Wall'
      }
    },
    {
      name: 'rankings-page',
      path: '/rankings',
      component: () => import('./pages/Rankings.vue'),
      meta: {
        pageTitle: 'Rankings | UIC Valentine Wall'
      }
    }
  ]
});

router.beforeEach((to, from, next) => {
  setCurrentScreen(analytics, to.meta.pageTitle as string);
  logEvent(analytics, 'page_view', {
    page_title: to.meta.pageTitle as string,
    page_location: window.location.href,
  });
  next();
});

export default router;