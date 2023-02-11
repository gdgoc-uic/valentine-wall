import { backendUrl, pb } from "./client";

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

export async function thirdPartyLogin(provider: string, table = 'users') {
  const authMethods = await pb.collection(table).listAuthMethods();
  const chosenProvider = authMethods.authProviders.find(p => p.name === provider);
  if (!chosenProvider) return;

  const redirectUrl = pb.buildUrl('/user_auth/callback');
  const connectUrl = `${chosenProvider.authUrl}${redirectUrl}&hd=uic.edu.ph`;
  const loginWindow = popupCenter({ url: connectUrl, title: provider + '_login_window', w: 800, h: 500 });
  if (!loginWindow) {
    throw new Error('Failed to open window.');
  }

  const expectedOrigin = (new URL(backendUrl)).origin;

  const handleFn = function (this: Window, e: MessageEvent) {
    if (e.origin === expectedOrigin && typeof e.data === 'object' && 'state' in e.data) {
      const data = e.data;
      if (chosenProvider.state !== data['state']) {
        throw new Error();
      }

      pb.collection(table).authWithOAuth2(
        chosenProvider.name,
        data['code'],
        chosenProvider.codeVerifier,
        redirectUrl,
      ).then(() => {
        window.removeEventListener('message', handleFn);
        loginWindow.close();
      });
    }
  }

  window.addEventListener('message', handleFn);
}