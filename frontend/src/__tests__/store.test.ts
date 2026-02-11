import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createAuthStore, createStore } from '../store_new'
import { pb } from '../client'

// Mock dependencies
vi.mock('../client', () => ({
  pb: {
    collection: vi.fn(() => ({
      getOne: vi.fn(),
      subscribe: vi.fn(),
      update: vi.fn()
    })),
    authStore: {
      clear: vi.fn(),
      token: 'mock-token'
    },
    send: vi.fn()
  }
}))

vi.mock('../auth', () => ({
  thirdPartyLogin: vi.fn()
}))

vi.mock('../notify', () => ({
  notify: vi.fn(),
  catchAndNotifyError: vi.fn()
}))

describe('Auth Store', () => {
  let authStore: any

  beforeEach(() => {
    authStore = createAuthStore()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('initializes with correct default state', () => {
    expect(authStore.state.user).toBeFalsy()
    expect(authStore.state.isAuthLoading).toBe(false)
    expect(authStore.state.isLoggedIn).toBe(false)
    expect(authStore.state.messageUnsubscribe).toBe(null)
  })

  it('sets isLoggedIn to true when user exists', async () => {
    const mockUser = {
      id: 'user123',
      email: 'test@example.com',
      expand: {
        details: {
          student_id: '202012345678',
          last_active: new Date().toISOString()
        },
        wallet: {
          balance: 1000
        }
      }
    }

    authStore.state.user = mockUser
    expect(authStore.state.isLoggedIn).toBe(true)
  })

  it('login method sets loading state', async () => {
    const { login } = authStore.methods

    const loginPromise = login()
    expect(authStore.state.isAuthLoading).toBe(true)

    await loginPromise.catch(() => {})
    expect(authStore.state.isAuthLoading).toBe(false)
  })

  it('logout clears user and auth store', () => {
    const mockUser = { id: 'user123', email: 'test@example.com' }
    authStore.state.user = mockUser

    authStore.methods.logout()

    expect(authStore.state.user).toBeFalsy()
    expect(pb.authStore.clear).toHaveBeenCalled()
  })

  it('logout unsubscribes from message notifications', () => {
    const mockUnsubscribe = vi.fn()
    authStore.state.messageUnsubscribe = mockUnsubscribe

    authStore.methods.logout()

    expect(mockUnsubscribe).toHaveBeenCalled()
    expect(authStore.state.messageUnsubscribe).toBe(null)
  })

  it('onReceiveUser subscribes to messages', async () => {
    const mockSubscribe = vi.fn()
    const mockUser = {
      id: 'user123',
      email: 'test@example.com',
      details: 'details123',
      expand: {
        'virtual_wallets(user)': { id: 'wallet123', balance: 1000 },
        details: {
          student_id: '202012345678',
          last_active: new Date().toISOString()
        }
      }
    }

    vi.mocked(pb.collection).mockReturnValue({
      getOne: vi.fn().mockResolvedValue(mockUser),
      subscribe: mockSubscribe.mockResolvedValue(vi.fn()),
      update: vi.fn().mockResolvedValue({})
    } as any)

    const mainStore = createStore()
    await authStore.methods.onReceiveUser(mockUser, mainStore.state)

    expect(mockSubscribe).toHaveBeenCalled()
  })
})

describe('Main Store', () => {
  let store: any

  beforeEach(() => {
    store = createStore()
    vi.clearAllMocks()
  })

  it('initializes with correct default state', () => {
    expect(store.state.giftList).toEqual([])
    expect(store.state.departmentList).toEqual([])
    expect(store.state.isSendMessageModalOpen).toBe(false)
    expect(store.state.isSetupModalOpen).toBe(false)
    expect(store.state.isWelcomeModalOpen).toBe(false)
  })

  it('has correct sex list options', () => {
    expect(store.state.sexList).toHaveLength(2)
    expect(store.state.sexList[0]).toEqual({ value: 'male', label: 'Male' })
    expect(store.state.sexList[1]).toEqual({ value: 'female', label: 'Female' })
  })

  it('loadGiftsAndDepartments populates lists', async () => {
    const mockDepts = [
      { id: '1', label: 'Computer Science' },
      { id: '2', label: 'Engineering' }
    ]
    const mockGifts = [
      { id: '1', uid: 'rose', label: 'Rose', price: 50 },
      { id: '2', uid: 'chocolate', label: 'Chocolate', price: 30 }
    ]

    vi.mocked(pb.send)
      .mockResolvedValueOnce(mockDepts)
      .mockResolvedValueOnce(mockGifts)

    await store.methods.loadGiftsAndDepartments()

    expect(store.state.departmentList).toHaveLength(2)
    expect(store.state.giftList).toHaveLength(2)
  })

  it('checkFirstTimeVisitor shows welcome modal for new users', () => {
    localStorage.clear()
    
    store.methods.checkFirstTimeVisitor()
    
    expect(store.state.isWelcomeModalOpen).toBe(true)
  })

  it('toggleWelcomeModal closes modal and sets localStorage', () => {
    store.state.isWelcomeModalOpen = true
    
    store.methods.toggleWelcomeModal()
    
    expect(store.state.isWelcomeModalOpen).toBe(false)
    expect(localStorage.getItem('__vw_13042021')).toBe('1')
  })
})
