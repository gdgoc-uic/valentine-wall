import { logEvent } from "firebase/analytics";
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
  nt.$notify(args);
  logEvent(analytics, 'server_notifications', args);
}

export function catchAndNotifyError(nt: Notifier, e: unknown) {
  if (e instanceof Error) {
    notify(nt, { type: 'error', text: e.message });
  } else {
    notify(nt, { type: 'error', text: `${e}` });
  }
}