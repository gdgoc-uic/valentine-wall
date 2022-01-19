import { logEvent, setUserId, setUserProperties } from 'firebase/analytics';
import { GoogleAuthProvider, signInWithPopup, User } from 'firebase/auth';
import { createStore } from 'vuex';
import { analytics, auth } from './firebase';

// NOTE: snake_case because JSON response is in snake_case
interface UserConnection {
  provider: string,
  refresh_token: string,
  expires_at: Date
}

export interface State {
  user: {
    email: string,
    id: string,
    associatedId: string,
    accessToken: string,
    connections: UserConnection[]
  },
  isAuthLoading: boolean,
  isIDModalOpen: boolean
}

export default createStore<State>({
  state() {
    return {
      user: {
        email: '',
        id: '',
        associatedId: '',
        accessToken: '',
        connections: []
      },
      isAuthLoading: true,
      isIDModalOpen: false
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
    SET_USER_ACCESS_TOKEN(state, payload: string) {
      state.user.accessToken = payload;
    },
    SET_ID_MODAL_OPEN(state, payload: boolean) {
      state.isIDModalOpen = payload;
    },
    SET_AUTH_LOADING(state, payload: boolean) {
      state.isAuthLoading = payload;
    },
    SET_USER_CONNECTIONS(state, payload: UserConnection[]) {
      state.user.connections = payload;
    }
  },

  getters: {
    isLoggedIn(state) {
      return state.user.accessToken.length != 0 && state.user.id.length != 0;
    },
    hasConnections(state) {
      return state.user.connections.length != 0;
    },
    headers(state, getters) {
      return getters.isLoggedIn ? {
        'Authorization': `Bearer ${state.user.accessToken}`
      } : {};
    }
  },

  actions: {
    async login({ commit, dispatch }) {
      const provider = new GoogleAuthProvider();
      provider.addScope('email');
      provider.addScope('profile');
      provider.setCustomParameters({
        'hd': 'uic.edu.ph'
      });

      try {
        commit('SET_AUTH_LOADING', true);
        await signInWithPopup(auth, provider);
        commit('SET_AUTH_LOADING', false);
      } catch (e) {
        commit('SET_AUTH_LOADING', false);
        throw e;
      }
    },
    async onReceiveUser({ commit, getters, dispatch }, user: User | null): Promise<void> {
      if (!user) {
        commit('SET_AUTH_LOADING', false);
        return;
      }

      try {
        commit('SET_USER_ID', user.uid);
        commit('SET_USER_EMAIL', user.email ?? '<unknown e-mail>');
        commit('SET_USER_ACCESS_TOKEN', await user.getIdToken(false));
        const idResp = await fetch(import.meta.env.VITE_BACKEND_URL + '/user/login_callback', {
          method: 'POST',
          headers: getters.headers,
          credentials: 'include'
        });

        const { associated_id, user_connections }: { associated_id: string, user_connections: UserConnection[] } = await idResp.json();
        commit('SET_USER_CONNECTIONS', user_connections);

        if (associated_id.length == 0) {
          commit('SET_ID_MODAL_OPEN', true);
        } else {
          commit('SET_USER_ASSOCIATED_ID', associated_id);
        }

        setUserId(analytics, user.uid);
        setUserProperties(analytics, { account_type: 'student' });
      } catch (e) {
        console.error(e);
        await dispatch('logout');
      } finally {
        commit('SET_AUTH_LOADING', false);
      }
    },
    async logout({ getters, commit }) {
      const logoutResp = await fetch(import.meta.env.VITE_BACKEND_URL + '/user/logout_callback', {
        method: 'POST',
        headers: getters.headers
      });

      if (logoutResp.status != 200) {
        throw Error('Unable to logout user.');
      }

      await auth.signOut();
      commit('SET_USER_ID', '');
      commit('SET_USER_EMAIL', '');
      commit('SET_USER_ACCESS_TOKEN', '');
      commit('SET_USER_ASSOCIATED_ID', '');
    }
  }
});