import { logEvent, setCurrentScreen } from 'firebase/analytics';
import { createRouter, createWebHistory, RouteLocationNormalizedLoaded } from 'vue-router';
import { expandAPIEndpoint } from './client';
import { analytics } from './firebase';

export function getPageTitle(route: RouteLocationNormalizedLoaded, pageSuffix: string): string {
  const pageTitle = route.meta.pageTitle;
  if (!pageTitle) return pageSuffix;

  if (pageTitle instanceof Function) {
    return `${pageTitle(route)} | ${pageSuffix}`
  }
  return `${pageTitle} | ${pageSuffix}`;
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      name: 'home-page',
      path: '/',
      component: () => import('./pages/Home.vue')
    },
    {
      name: 'message-wall-page',
      path: '/wall/:recipientId?',
      component: () => import('./pages/Wall.vue'),
      meta: {
        pageTitle: (route: RouteLocationNormalizedLoaded) => {
          return route.params.recipientId ? `Messages for ${route.params.recipientId}` : 'Recent messages';
        }
      }
    },
    {
      name: 'message-page',
      path: '/wall/:recipientId/:messageId',
      component: () => import('./pages/Message.vue'),
      meta: {
        pageTitle: (route: RouteLocationNormalizedLoaded) => {
          return `Message for ${route.params.recipientId}`;
        },
        metaTags: (route: RouteLocationNormalizedLoaded) => [
          {
            name: 'og:image',
            content: expandAPIEndpoint(`/messages/${route.params.recipientId}/${route.params.messageId}?image`)
          },
          {
            name: 'og:image:width',
            content: '1200'
          },
          {
            name: 'og:image:height',
            content: '675'
          }
        ]
      }
    },
    {
      name: 'rankings-page',
      path: '/rankings',
      component: () => import('./pages/Rankings.vue'),
      meta: {
        pageTitle: 'Rankings'
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