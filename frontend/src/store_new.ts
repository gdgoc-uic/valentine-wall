import { ClientResponseError } from 'pocketbase';
import { computed, inject, InjectionKey, onMounted, onUnmounted, reactive, readonly, watch } from 'vue';
import { popupCenter } from './auth';
import { backendUrl, pb } from './client';
import { CollegeDepartment, Gift, User, VirtualWallet } from './types';

interface StoreMethods {
  loadGiftsAndDepartments(): Promise<void>
  checkFirstTimeVisitor(): void
  toggleWelcomeModal(): void
}

export interface Store<T, U = {}> {
  state: T,
  methods: U
}

export interface State {
  isSendMessageModalOpen: boolean,
  isSetupModalOpen: boolean,
  isWelcomeModalOpen: boolean,
  giftList: Gift[],
  departmentList: CollegeDepartment[],
  sexList: readonly { value: string, label: string }[]
}

export const storeKey = Symbol() as InjectionKey<Store<State, StoreMethods>>;

export const useStore = () => inject(storeKey)!;

const firstTimeKey = '__vw_13042021';

export const createStore = (): Store<State, StoreMethods> => {
  const state = reactive<State>({
    giftList: [],
    isSendMessageModalOpen: false,
    isSetupModalOpen: false,
    isWelcomeModalOpen: false,
    departmentList: [],
    sexList: readonly([
      { value: 'male', label: 'Male' },
      { value: 'female', label: 'Female' },
      // { value: 'unknown', label: 'Prefer not to say' }
    ])
  });

  const checkFirstTimeVisitor = () => {
    if (import.meta.env.VITE_READ_ONLY === "true") return;
    if (!localStorage.getItem(firstTimeKey)) {
      state.isWelcomeModalOpen = true;
    }
  }

  const loadGiftsAndDepartments = async () => {
    const depts = await pb.send('/departments', {});
    state.departmentList.push(...depts);

    const gifts = await pb.send('/gifts', {});
    state.giftList.push(...gifts);
  }

  const toggleWelcomeModal = () => {
    localStorage.setItem(firstTimeKey, '1');
    state.isWelcomeModalOpen = false;
  }
  
  return reactive({
    state,
    methods: {
      loadGiftsAndDepartments,
      checkFirstTimeVisitor,
      toggleWelcomeModal
    }
  });
}

// Auth Store
export interface AuthState {
  user: User
  isAuthLoading: boolean
  isLoggedIn: boolean
}

export interface AuthMethods {
  login(): Promise<void>
  logout(): void
  onReceiveUser(user: User | null, state: State): Promise<void>
}

export const useAuth = () => inject(authStore)!;

export const authStore: InjectionKey<Store<AuthState, AuthMethods>> = Symbol();

export function createAuthStore(): Store<AuthState, AuthMethods> {
  const state = reactive({
    user: null!,
    isAuthLoading: false,
    isLoggedIn: computed(() => !!state.user),
  }) as AuthState;

  async function login() {
    try {
      state.isAuthLoading = true;
  
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
          ).then(() => {
            window.removeEventListener('message', handleFn);
            loginWindow.close();
          });
        }
      }
  
      window.addEventListener('message', handleFn);
    } catch (e) {
      throw e;
    } finally {
      state.isAuthLoading = false;
    }
  }

  async function onReceiveUser(receivedUser: User | null, mainStore: State): Promise<void> {  
    try {
      if (!receivedUser) {
        return;
      }

      // verify user
      const user = await pb.collection('users').getOne<User>(receivedUser.id, {
        expand: 'virtual_wallets(user),details'
      });

      user.expand.wallet = user.expand['virtual_wallets(user)'] as VirtualWallet;
      delete user.expand['virtual_wallets(user)'];

      if (import.meta.env.VITE_READ_ONLY !== 'true') {  
        if (!user.details) {
          mainStore.isSetupModalOpen = true;
        }

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
      state.user = user;
    } catch (e) {
      console.error(e);
      logout();
    } finally {
      state.isAuthLoading = false;
    }
  }

  function logout() {
    try {
      if (import.meta.env.VITE_READ_ONLY !== 'true') {
        // await getters.apiClient.post('/user/logout_callback');
      }

      console.log('LOGOUT');
      pb.authStore.clear();
    } catch (e) {
      // if (e instanceof APIResponseError) {
      //   throw Error('Unable to logout user.');
      // }
      throw e;
    }
  }

  return reactive({
    state,
    methods: {
      login, 
      logout,
      onReceiveUser
    }
  }) as Store<AuthState, AuthMethods>;
}