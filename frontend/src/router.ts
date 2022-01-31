import { createRouter as createVueRouter, createWebHistory, createMemoryHistory, RouteLocationNormalizedLoaded, RouteRecordRaw } from 'vue-router';
import { expandAPIEndpoint } from './client';

const pageSuffix = 'UIC Valentine Wall';
const pageSeparator = ' | ';

export function getPageTitle(route: RouteLocationNormalizedLoaded): string {
  const pageTitle = route.meta.pageTitle;
  if (!pageTitle) return pageSuffix;

  if (pageTitle instanceof Function) {
    return `${pageTitle(route)}${pageSeparator}${pageSuffix}`
  }
  return `${pageTitle}${pageSeparator}${pageSuffix}`;
}

const routes: RouteRecordRaw[] = [
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
];

export function createRouter() {
  return createVueRouter({
    history: import.meta.env.SSR ? createMemoryHistory(): createWebHistory(),
    routes
  });
}