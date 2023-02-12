import { createRouter as createVueRouter, createWebHistory, createMemoryHistory, RouteLocationNormalizedLoaded, RouteRecordRaw } from 'vue-router';
import { pb } from './client';

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

const readOnlyRoutes: RouteRecordRaw[] = [
  {
    name: 'home-page',
    path: '/',
    component: () => import('./pages/Home_ReadOnly.vue'),
    meta: {
      disableAppHeader: true
    }
  },
];

const routes: RouteRecordRaw[] = [
  {
    name: 'home-page',
    path: '/',
    component: () => import('./pages/Home.vue')
  },
  {
    name: 'recent-wall-page',
    path: '/wall',
    component: () => import('./pages/Wall.vue'),
    meta: {
      pageTitle: 'Recent Messages'
    }
  },
  {
    path: '/wall/everyone',
    redirect: {
      name: 'recent-wall-page',
    }
  },
  {
    name: 'message-wall-page',
    path: '/wall/:recipientId',
    component: () => import('./pages/Wall.vue'),
    meta: {
      pageTitle: (route: RouteLocationNormalizedLoaded) => {
        return route.params.recipientId ? `Messages for ${route.params.recipientId}` : 'Recent messages';
      }
    }
  },
  {
    name: 'old-message-page',
    path: '/messages/:recipientId/:messageId',
    redirect: to => ({
      name: 'message-page',
      params: to.params
    })
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
          content: pb.buildUrl(`/messages/${route.params.messageId}/image`)
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
        name: 'settings-transactions-section',
        path: 'transactions',
        component: () => import('./pages/settings/Transactions.vue'),
        meta: {
          pageTitle: 'Transactions | Settings',
          label: 'Transactions'
        }
      },
      {
        name: 'settings-delete-account-section',
        path: 'delete-account',
        component: () => import('./pages/settings/DeleteAccount.vue'),
        meta: {
          pageTitle: 'Archive / Delete Account | Settings',
          label: 'Archive / Delete Account'
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('./pages/NotFound.vue'),
    meta: {
      pageTitle: 'Page not found.'
    }
  }
];

export function createRouter() {
  return createVueRouter({
    history: import.meta.env.SSR ? createMemoryHistory(): createWebHistory(),
    scrollBehavior(to, from, savedPosition) {
      // always scroll to top
      return { top: 0 }
    },
    routes: import.meta.env.VITE_READ_ONLY === 'true' ? readOnlyRoutes : routes
  });
}