import { logEvent, setUserId, setUserProperties } from 'firebase/analytics';
import { createStore as createVuexStore, Store } from 'vuex';
import { InjectionKey } from 'vue';
import { backendUrl, pb } from './client';
import { popupCenter } from './auth';
import { Record as PbRecord } from 'pocketbase';
// import { analytics, auth } from './firebase';

// NOTE: snake_case because JSON response is in snake_case
interface UserConnection {
  provider: string
}

export interface Gift {
  id: number
  uid: string
  label: string
}

export interface VirtualWallet {
  balance: number
}

export interface State {
  user: {
    email: string,
    id: string,
    associatedId: string,
    sex: string,
    department: string,
    accessToken: string,
    wallet: VirtualWallet | null,
    connections: UserConnection[]
  },
  isAuthLoading: boolean,
  isSendMessageModalOpen: boolean,
  isSetupModalOpen: boolean,
  isWelcomeModalOpen: boolean,
  giftList: Gift[],
  departmentList: Record<string, string>[]
}

export const storeKey = 'vuex-store' as unknown as InjectionKey<Store<State>>;

const firstTimeKey = '__vw_13042021';

export function createStore() {
  return createVuexStore<State>({
    state() {
      return {
        user: {
          email: '',
          id: '',
          associatedId: '',
          accessToken: '',
          sex: '',
          department: '',
          wallet: null,
          connections: []
        },
        isAuthLoading: true,
        isSendMessageModalOpen: false,
        isSetupModalOpen: false,
        isWelcomeModalOpen: false,
        giftList: [],
        departmentList: []
      }
    },
  
    mutations: {
      SET_USER_EMAIL(state, payload: string) {
        state.user.email = payload;
      },
      SET_USER_ID(state, payload: string) {
        state.user.id = payload;
      },
      SET_USER_ASSOCIATED_ID(state, payload: string) {
        state.user.associatedId = payload;
      },
      SET_USER_SEX(state, payload: string) {
        state.user.sex = payload;
      },
      SET_USER_DEPARTMENT(state, payload: string) {
        state.user.department = payload;
      },
      SET_USER_ACCESS_TOKEN(state, payload: string) {
        state.user.accessToken = payload;
      },
      SET_USER_WALLET(state, payload: VirtualWallet | null) {
        state.user.wallet = payload;
      },
      SET_USER_WALLET_BALANCE(state, payload: number) {
        if (state.user.wallet) {
          state.user.wallet.balance = payload;
        }
      },
      SET_SETUP_MODAL_OPEN(state, payload: boolean) {
        state.isSetupModalOpen = payload;
      },
      SET_AUTH_LOADING(state, payload: boolean) {
        state.isAuthLoading = payload;
      },
      SET_USER_CONNECTIONS(state, payload: UserConnection[]) {
        state.user.connections = payload;
      },
      SET_GIFT_LIST(state, payload: Gift[]) {
        state.giftList = payload;
      },
      SET_DEPARTMENT_LIST(state, payload: Record<string, string>[]) {
        state.departmentList = payload;
      },
      SET_SEND_MESSAGE_MODAL_OPEN(state, payload: boolean) {
        state.isSendMessageModalOpen = payload;
      },
      SET_WELCOME_MODAL_OPEN(state, payload: boolean) {
        state.isWelcomeModalOpen = payload;
      }
    },
  
    getters: {
      isLoggedIn(state) {
        return state.user.accessToken.length != 0 && state.user.id.length != 0;
      },
      hasConnections(state) {
        return typeof state.user.connections != 'undefined' && state.user.connections.length != 0;
      },
      headers(state, getters) {
        return getters.isLoggedIn ? {
          'Authorization': `Bearer ${state.user.accessToken}`
        } : {};
      },
      username(state, getters) {
        if (!getters.isLoggedIn) {
          return '';
        }
        const emailRegex = /(^[a-z]+)_[0-9]+@uic.edu.ph$/;
        const res = emailRegex.exec(state.user.email);
        if (!res || res.length < 2) return '';
        return res[1];
      },
      sexList() {
        return [
          {
            label: 'Male',
            value: 'male'
          },
          {
            label: 'Female',
            value: 'female'
          }
        ];
      }
    },
  
    actions: {
      async login({ commit, state }) {
        try {
          commit('SET_AUTH_LOADING', true);
          // TODO: improve auth
          // 1. create an endpoint in server that will receive the callback
          // 2. create an HTML template that will receive the context (html/auth_callback.html.tpl)
          // 3. using window.postMessage and the code from connectToTwitter in auth.ts, receive the payload
          // 4. verify the code on the client side
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
              })
            
              window.removeEventListener('message', handleFn);
              loginWindow.close();
            }
          }
  
          window.addEventListener('message', handleFn);
          if (state.isSetupModalOpen) {
            // logEvent(analytics!, 'sign_up');
          } else {
            // logEvent(analytics!, 'login');
          }
        } catch (e) {
          throw e;
        } finally {
          commit('SET_AUTH_LOADING', false);
        }
      },
      async onReceiveUser({ commit, dispatch, getters }, user: PbRecord | null): Promise<void> {  
        try {
          if (!user) {
            return;
          }
  
          commit('SET_USER_ID', user.id);
          commit('SET_USER_EMAIL', user.email ?? '<unknown e-mail>');
          commit('SET_USER_ACCESS_TOKEN', 'dummy TODO');

          if (import.meta.env.VITE_READ_ONLY !== 'true') {
            // await Promise.all([
            //   getters.apiClient.post('/user/session', { credentials: 'include' }),
            //   dispatch('getUserInfo')
            // ]);
  
            // it should not affect the whole flow just in case
            // updateLastActiveAt won't go through
            try {
              await dispatch('updateLastActiveAt');
            } catch (e) {
              console.error(e);
            }
          }

          // setUserId(analytics!, user.uid);
          // setUserProperties(analytics!, { account_type: 'student' });
        } catch (e) {
          console.error(e);
          await dispatch('logout');
        } finally {
          commit('SET_AUTH_LOADING', false);
        }
      },
      async logout({ commit, getters }) {
        try {
          if (import.meta.env.VITE_READ_ONLY !== 'true') {
            await getters.apiClient.post('/user/logout_callback');
          }

          // await auth.signOut();
          pb.authStore.clear();

          commit('SET_USER_ID', '');
          commit('SET_USER_EMAIL', '');
          commit('SET_USER_ACCESS_TOKEN', '');
          commit('SET_USER_ASSOCIATED_ID', '');
          commit('SET_USER_SEX', '');
          commit('SET_USER_DEPARTMENT', '');
          commit('SET_USER_WALLET', null);
        } catch (e) {
          // if (e instanceof APIResponseError) {
          //   throw Error('Unable to logout user.');
          // }
          throw e;
        }
      },
      async getUserInfo({ commit, getters }) {
        const { 
          data: { associated_id, department, sex, user_connections, wallet } 
        }: { 
          data: { 
            associated_id: string, 
            department: string,
            sex: string,
            user_connections: UserConnection[],
            wallet: VirtualWallet | null
          } 
        } = await getters.apiClient.get('/user/info');
        if (associated_id.length == 0) {
          commit('SET_SETUP_MODAL_OPEN', true);
        } else {
          commit('SET_USER_ASSOCIATED_ID', associated_id);
          commit('SET_USER_SEX', sex);
          commit('SET_USER_DEPARTMENT', department);
          commit('SET_USER_CONNECTIONS', user_connections);
          commit('SET_USER_WALLET', wallet);
        }
      },
      async getGiftList({ commit, getters }) {
        const { data: json } = await pb.send('/gifts', {});
        commit('SET_GIFT_LIST', json);
      },
      async getDepartmentList({ commit, getters }) {
        const { data: json } = await pb.send('/departments', {});
        commit('SET_DEPARTMENT_LIST', json);
      },
      async updateLastActiveAt({ getters }) {
        await getters.apiClient.patch('/user/session');
      },
      checkFirstTimeVisitor({ commit }) {
        if (import.meta.env.VITE_READ_ONLY === "true") return;
        if (!localStorage.getItem(firstTimeKey)) {
          commit('SET_WELCOME_MODAL_OPEN', true);
        }
      },
      toggleWelcomeModal({ commit }) {
        localStorage.setItem(firstTimeKey, '1');
        commit('SET_WELCOME_MODAL_OPEN', false);
      }
    }
  })
};