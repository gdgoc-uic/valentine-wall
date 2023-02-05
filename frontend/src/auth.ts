import { Store } from "vuex";
import { State } from "./store";

export function popupCenter({ url, title, w, h }: { url: string, title: string, w: number, h: number }): Window | null {
    // Fixes dual-screen position                             Most browsers      Firefox
    const dualScreenLeft = typeof window.screenLeft !== 'undefined' ? window.screenLeft : window.screenX;
    const dualScreenTop = typeof window.screenTop !== 'undefined' ? window.screenTop : window.screenY;

    const width = window.innerWidth ? window.innerWidth : document.documentElement.clientWidth ? document.documentElement.clientWidth : screen.width;
    const height = window.innerHeight ? window.innerHeight : document.documentElement.clientHeight ? document.documentElement.clientHeight : screen.height;

    const systemZoom = width / window.screen.availWidth;
    const left = ((width - w) / 2) + systemZoom + dualScreenLeft
    const top = ((height - h) / 2) + systemZoom + dualScreenTop
    const newWindow = window.open(url, title,
        `
      scrollbars=yes,
      width=${w}, 
      height=${h}, 
      top=${top}, 
      left=${left}
      `
    )

    newWindow?.focus();
    return newWindow;
};

interface ConnectCallbacks {
    onSuccess?: () => void,
    onError?: (e: unknown) => void,
    onFinally?: () => void
}

export function connectToTwitter(store: Store<State>, cbs?: ConnectCallbacks) {
    // const connectUrl = expandAPIEndpoint('/user/connect_twitter');
    const connectUrl = '';
    const loginWindow = popupCenter({ url: connectUrl, title: 'twitter_login_window', w: 800, h: 500 });
    if (!loginWindow) {
        cbs?.onError?.(new Error('Failed to open window.'));
        cbs?.onFinally?.();
        return;
    }

    const handleFn = function (this: Window, e: MessageEvent) {
        if (e.origin === import.meta.env.VITE_BACKEND_URL && typeof e.data === 'object' && 'message' in e.data) {
            const data = e.data;
            if (data['message'] !== 'twitter connect success') {
                cbs?.onError?.(new Error());
                cbs?.onFinally?.();
                return 
            }
            store.commit('SET_USER_CONNECTIONS', data['user_connections']);
            cbs?.onSuccess?.();
            cbs?.onFinally?.();
            window.removeEventListener('message', handleFn);
            loginWindow.close();
        }
    }

    window.addEventListener('message', handleFn);
}

export async function connectToEmail(store: Store<State>) {
    const { data: json } = await store.getters.apiClient.post('/user/connect_email');
    store.commit('SET_USER_CONNECTIONS', json['user_connections']);
}