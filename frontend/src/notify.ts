import { logEvent } from "firebase/analytics";
import { Component, h, Plugin } from "vue";
import { analytics } from "./firebase";

export interface NotifierArguments {
  title?: string,
  text: string,
  type: string,
  group?: string,
  duration?: number,
  speed?: number,
  data?: any,
  clean?: boolean
}

export interface Notifier {
  $notify: (args: NotifierArguments, duration?: number) => void
}

export function notify(nt: Notifier, args: NotifierArguments) {
  if (!import.meta.env.SSR){
    nt.$notify(args);
    logEvent(analytics!, 'server_notifications', args);
  }
}

export function catchAndNotifyError(nt: Notifier, e: unknown) {
  if (e instanceof Error) {
    notify(nt, { type: 'error', text: e.message });
  } else {
    notify(nt, { type: 'error', text: `${e}` });
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
      app.config.globalProperties.$notify = (n: any, t: any) => {};
      app.component('Notification', notificationSSR);
      app.component('NotificationGroup', notificationsGroupSSR);

      // Compatibility with the old component names
      app.component('notification', notificationSSR);
      app.component('notificationGroup', notificationsGroupSSR);
    }
  }
}