import dayjs from 'dayjs';
import { computed, inject, InjectionKey, reactive, readonly } from 'vue';
import { thirdPartyLogin } from './auth';
import { pb } from './client';
import { catchAndNotifyError, notify } from './notify';
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
  reward(amount: number, reason: string): Promise<void>
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
      await thirdPartyLogin('google');
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
          const now = new Date();
          if (user.expand.details && user.expand.details.last_active) {
            const before = dayjs(user.expand.details.last_active);
            const after = dayjs(now);
            const diff = after.diff(before, 'second');

            // rewarded for idle (0.05 coins per second / rounded)
            const idleReward = Math.round(0.05 * diff);
            if (idleReward > 0) {
              await reward(0.05 * diff, 'Idle time', user.expand.wallet.id);
            }
          }

          await pb.collection('user_details').update(user!.details, {
            last_active: now
          });
        } catch (e) {
          console.error(e);
        }
      }

      // setUserId(analytics!, user.uid);
      // setUserProperties(analytics!, { account_type: 'student' });
      state.user = user;

      await pb.collection('virtual_wallets').subscribe(user.expand.wallet.id, (data) => {
        if (data.action === 'update') {
          state.user.expand.wallet.balance = (data.record as VirtualWallet).balance;
        }
      });
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

      state.user = null!;
      pb.authStore.clear();
    } catch (e) {
      // if (e instanceof APIResponseError) {
      //   throw Error('Unable to logout user.');
      // }
      throw e;
    }
  }

  async function reward(amount: number, reason: string, walletId?: string) {
    if (amount <= 0) {
      return;
    }

    try {
      if (!walletId) {
        if (!state.user || !state.user.expand.wallet) {
          throw new Error('Unable to add money to your wallet.');
        } 
      }

      await pb.collection('virtual_transactions').create({
        wallet: walletId ?? state.user.expand.wallet.id,
        description: `Rewarded for: ${reason}`,
        amount,
      }, {
        magicword: 'V3L3NT1L3'
      });

      notify({ type: 'success', text: `You have been credited ${amount} into your account. Enjoy!` });
    } catch(e) {
      catchAndNotifyError(e);
    }
  }

  return reactive({
    state,
    methods: {
      login, 
      logout,
      reward,
      onReceiveUser
    }
  }) as Store<AuthState, AuthMethods>;
}