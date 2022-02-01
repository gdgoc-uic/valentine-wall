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
  },
  {
    name: 'about-page',
    path: '/about',
    component: () => import('./pages/About.vue'),
  },
  {
    name: 'settings-page',
    path: '/settings',
    component: () => import('./pages/Settings.vue'),
    redirect: {
      name: 'settings-basic-info-section'
    },
    meta: {
      pageTitle: 'Settings',
      requiresAuth: true
    },
    children: [
      {
        name: 'settings-basic-info-section',
        path: 'info',
        component: () => import('./pages/settings/BasicInformation.vue'),
        meta: {
          pageTitle: 'Basic Information | Settings',
          label: 'Basic Information'
        }
      },
      {
        name: 'settings-delete-account-section',
        path: 'delete-account',
        component: () => import('./pages/settings/DeleteAccount.vue'),
        meta: {
          pageTitle: 'Delete Account | Settings',
          label: 'Delete Account'
        }
      }
    ]
  }
];

export function createRouter() {
  return createVueRouter({
    history: import.meta.env.SSR ? createMemoryHistory(): createWebHistory(),
    routes
  });
}