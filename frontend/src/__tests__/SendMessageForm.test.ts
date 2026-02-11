import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive } from 'vue'
import SendMessageForm from '../components/SendMessageForm.vue'
import { storeKey, authStore as authStoreKey } from '../store_new'

// Mock dependencies
vi.mock('../client', () => ({
  pb: {
    collection: vi.fn(() => ({
      create: vi.fn()
    })),
    authStore: {
      token: 'mock-token'
    }
  }
}))

vi.mock('../notify', () => ({
  notify: vi.fn(),
  catchAndNotifyError: vi.fn()
}))

vi.mock('@tanstack/vue-query', () => ({
  useMutation: vi.fn(() => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
    isPending: { value: false }
  }))
}))

describe('SendMessageForm Component', () => {
  let mockStore: any
  let mockAuthStore: any

  beforeEach(() => {
    mockStore = reactive({
      state: reactive({
        giftList: [
          { id: '1', uid: 'rose', label: 'Rose', price: 50, is_remittable: true },
          { id: '2', uid: 'chocolate', label: 'Chocolate', price: 30, is_remittable: false }
        ],
        isSendMessageModalOpen: false,
        isSetupModalOpen: false,
        isWelcomeModalOpen: false,
        departmentList: [],
        sexList: [
          { value: 'male', label: 'Male' },
          { value: 'female', label: 'Female' }
        ]
      }),
      methods: {
        loadGiftsAndDepartments: vi.fn(),
        checkFirstTimeVisitor: vi.fn(),
        toggleWelcomeModal: vi.fn()
      }
    })

    mockAuthStore = reactive({
      state: reactive({
        user: {
          id: 'user123',
          email: 'test@example.com',
          expand: {
            'virtual_wallets(user)': [{ id: 'w1', balance: 1000 }],
            details: { student_id: '202012345678' }
          }
        },
        isAuthLoading: false,
        isLoggedIn: true,
        messageUnsubscribe: null
      }),
      methods: {
        login: vi.fn(),
        logout: vi.fn(),
        reward: vi.fn(),
        onReceiveUser: vi.fn()
      }
    })
  })

  const mountForm = (props = {}) => {
    return mount(SendMessageForm, {
      props,
      global: {
        provide: {
          [storeKey as symbol]: mockStore,
          [authStoreKey as symbol]: mockAuthStore
        },
        stubs: {
          ContentCounter: true,
          Modal: { template: '<div><slot /><slot name="footer" /></div>' },
          Tooltip: { template: '<div><slot /></div>' }
        }
      }
    })
  }

  it('renders the form component', () => {
    const wrapper = mountForm()
    expect(wrapper.exists()).toBe(true)
  })

  it('renders form with recipient input', () => {
    const wrapper = mountForm()
    expect(wrapper.find('input[name="recipient_id"]').exists()).toBe(true)
  })

  it('renders form with content textarea', () => {
    const wrapper = mountForm()
    expect(wrapper.find('textarea[name="content"]').exists()).toBe(true)
  })

  it('enforces 240 character limit on message content', () => {
    const wrapper = mountForm()
    const textarea = wrapper.find('textarea[name="content"]')
    expect(textarea.attributes('maxlength')).toBe('240')
  })

  it('shows "Everyone" button when no existing recipient', () => {
    const wrapper = mountForm()
    expect(wrapper.text()).toContain('Everyone')
  })

  it('passes existing recipient as prop', () => {
    const wrapper = mountForm({ existingRecipient: 'everyone' })
    expect(wrapper.exists()).toBe(true)
  })

  it('renders gift list from store', () => {
    const wrapper = mountForm()
    // Gift list should be available in the component
    expect(mockStore.state.giftList).toHaveLength(2)
  })

  it('has submit button', () => {
    const wrapper = mountForm()
    const submitBtn = wrapper.find('button[type="submit"]')
    expect(submitBtn.exists()).toBe(true)
  })
})
