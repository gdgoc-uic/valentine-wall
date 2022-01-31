import { Store } from 'vuex'
import { Notifier } from './notify'
import { State } from './store'
import { APIClientPlugin } from './client'

declare global {
  interface Window {
    __INITIAL_STATE__: State;
  }
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties extends Notifier, APIClientPlugin {
    $store: Store<State>
  }
}