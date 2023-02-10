import { ClientResponseError } from 'pocketbase';
import { computed, inject, InjectionKey, onMounted, onUnmounted, reactive, watch } from 'vue';
import { popupCenter } from './auth';
import { backendUrl, pb } from './client';
import { Gift, User, UserDetails, VirtualWallet } from './types';

interface StoreMethods {
  loadGiftsAndDepartments(): Promise<void>
  checkFirstTimeVisitor(): void
  toggleWelcomeModal(): void
}

export interface State {
  isSendMessageModalOpen: boolean,
  isSetupModalOpen: boolean,
  isWelcomeModalOpen: boolean,
  giftList: Gift[],
  departmentList: Record<string, string>[]
}

export const storeKey = Symbol() as InjectionKey<State & StoreMethods>;

export const useStore = () => inject(storeKey)!;

const firstTimeKey = '__vw_13042021';

export const createStore = () => {
  const store = reactive<State>({
    giftList: [],
    isSendMessageModalOpen: false,
    isSetupModalOpen: false,
    isWelcomeModalOpen: false,
    departmentList: []
  });

  const checkFirstTimeVisitor = () => {
    if (import.meta.env.VITE_READ_ONLY === "true") return;
    if (!localStorage.getItem(firstTimeKey)) {
      store.isWelcomeModalOpen = true;
    }
  }

  const loadGiftsAndDepartments = async () => {
    const depts = await pb.send('/departments', {});
    store.departmentList.push(...depts);

    const gifts = await pb.send('/gifts', {});
    store.giftList.push(...gifts);
  }

  const toggleWelcomeModal = () => {
    localStorage.setItem(firstTimeKey, '1');
    store.isWelcomeModalOpen = false;
  }

  // onMounted(() => {
  //   if (import.meta.env.VITE_READ_ONLY !== "true") {
  //     loadGiftsAndDepartments()
  //       .catch((err) => {
  //         // catchAndNotifyError(this, err);
  //         console.error(err);
  //       });
  //   }
  // })
  
  return {
    ...store,
    loadGiftsAndDepartments,
    checkFirstTimeVisitor,
    toggleWelcomeModal
  };
}

// Auth Store
export interface AuthStore {
  user: User
  isAuthLoading: boolean
  isLoggedIn: boolean
  login(): Promise<void>
  logout(): void
}

export const useAuth = () => inject(authStore)!;

export const authStore: InjectionKey<AuthStore> = Symbol();

export function createAuthStore(): AuthStore {
  const store = reactive({
    user: null!,
    isAuthLoading: false,
    isLoggedIn: computed(() => !!pb.authStore.model),
  }) as Omit<AuthStore, 'login' | 'logout'>;

  const mainStore = useStore();

  async function login() {
    try {
      store.isAuthLoading = true;
  
      const authMethods = await pb.collection('users').listAuthMethods();
      const googleProvider = authMethods.authProviders.find(p => p.name === 'google');
      if (!googleProvider) return;
  
      const redirectUrl = pb.buildUrl('/user_auth/callback');
  
      // TODO: change url on production!
      const connectUrl = `${googleProvider.authUrl}${redirectUrl}&hd=uic.edu.ph`;
      const loginWindow = popupCenter({ url: connectUrl, title: 'twitter_login_window', w: 800, h: 500 });
      if (!loginWindow) {
        throw new Error('Failed to open window.');
      }
  
      const handleFn = function (this: Window, e: MessageEvent) {
        if (e.origin === backendUrl && typeof e.data === 'object' && 'state' in e.data) {
          const data = e.data;
          if (googleProvider.state !== data['state']) {
            throw new Error();
          }
  
          pb.collection('users').authWithOAuth2(
            googleProvider.name,
            data['code'],
            googleProvider.codeVerifier,
            redirectUrl,
          ).then((authData) => {
            console.log(authData);
            
            window.removeEventListener('message', handleFn);
            loginWindow.close();
          });
        }
      }
  
      window.addEventListener('message', handleFn);
      if (mainStore.isSetupModalOpen) {
        // logEvent(analytics!, 'sign_up');
      } else {
        // logEvent(analytics!, 'login');
      }
    } catch (e) {
      throw e;
    } finally {
      store.isAuthLoading = false;
    }
  }

  async function onReceiveUser(user: User | null): Promise<void> {  
    try {
      if (!user) {
        return;
      }

      // TODO: use vue reactive state
      try {
        const wallet = await pb.collection('virtual_wallets').getFirstListItem(`user="${user.id}"`);
        user!.expand.wallet = wallet as VirtualWallet;
      } catch {}

      try {
        if (user.details) {
          const result = await pb.collection('user_details').getOne(user.details);
          user.expand.details = result as UserDetails;
        } else {
          const result = await pb.collection('user_details').getFirstListItem(`user="${user.id}"`);
          user.expand.details = result as UserDetails;
        }

        store.user = user;
      } catch (e) {
        if (e instanceof ClientResponseError && e.status === 404) {
          mainStore.isSetupModalOpen = true;
        }
      }
      
      if (import.meta.env.VITE_READ_ONLY !== 'true') {  
        // it should not affect the whole flow just in case
        // updateLastActiveAt won't go through
        try {
          await pb.collection('user_details').update(user!.details, {
            last_active: new Date()
          });
        } catch (e) {
          console.error(e);
        }
      }

      // setUserId(analytics!, user.uid);
      // setUserProperties(analytics!, { account_type: 'student' });
    } catch (e) {
      console.error(e);
      logout();
    } finally {
      store.isAuthLoading = false;
    }
  }

  function logout() {
    try {
      if (import.meta.env.VITE_READ_ONLY !== 'true') {
        // await getters.apiClient.post('/user/logout_callback');
      }

      pb.authStore.clear();
    } catch (e) {
      // if (e instanceof APIResponseError) {
      //   throw Error('Unable to logout user.');
      // }
      throw e;
    }
  }

  const unwatchUser = watch(store.user, (newUser) => {
    if (!newUser) return;

    // TODO: auto-update last active at
  });

  // onMounted(() => {
  //   pb.authStore.onChange((_, user) => {
  //     onReceiveUser(user as User | null);
  //   }, true);
  // });

  // onUnmounted(() => {
  //   unwatchUser();
  // });

  return { ...store, login, logout } as unknown as AuthStore;
}