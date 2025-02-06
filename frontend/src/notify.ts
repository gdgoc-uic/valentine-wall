import { Component, h, Plugin } from "vue";
import { notify as _notify } from 'notiwind';
import type { NotificationItem } from "notiwind/dist/notify";
import { logEvent } from './analytics';

export function notify<T>(args: NotificationItem<T>) {
  if (!import.meta.env.SSR){
    _notify(args);
    logEvent('server_notifications', args);
  }
}

export function catchAndNotifyError(e: unknown) {
  if (e instanceof Error && e.message) {
    notify({ type: 'error', text: e.message });
  } else {
    notify({ type: 'error', text: `Unknown error.` });
  }
}

export function notiwindSSRShim(): Plugin {
  const notificationSSR: Component = {
    name: 'Notification',
    render() {
      return h('div', []);
    }
  }

  const notificationsGroupSSR: Component = {
    name: 'NotificationsGroup',
    render() {
      return h('div', []);
    }
  }

  return {
    install(app) {
      app.component('Notification', notificationSSR);
      app.component('NotificationGroup', notificationsGroupSSR);

      // Compatibility with the old component names
      app.component('notification', notificationSSR);
      app.component('notificationGroup', notificationsGroupSSR);
    }
  }
}