import { Store } from 'vuex'
import { Notifier } from './notify'
import { State } from './store'

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties extends Notifier {
    $store: Store<State>
  }
}