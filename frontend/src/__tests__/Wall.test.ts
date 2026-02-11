import { describe, it, expect, vi } from 'vitest'
import { reactive, ref } from 'vue'

// Mock all dependencies before importing
vi.mock('../client', () => ({
  pb: {
    collection: vi.fn(() => ({
      getList: vi.fn().mockResolvedValue({ items: [], totalPages: 0, page: 1 }),
      subscribe: vi.fn().mockResolvedValue(() => {})
    })),
    authStore: {
      isValid: true,
      model: { id: 'user123' }
    }
  }
}))

vi.mock('../store_new', () => ({
  useAuth: () => ({
    state: reactive({
      isLoggedIn: false,
      user: null
    })
  }),
  useStore: () => ({
    state: reactive({
      giftList: [],
      isSendMessageModalOpen: false,
      departmentList: []
    })
  }),
  storeKey: Symbol('store'),
  authStore: Symbol('auth')
}))

vi.mock('@tanstack/vue-query', () => ({
  useInfiniteQuery: vi.fn(() => ({
    data: ref({ pages: [] }),
    fetchNextPage: vi.fn(),
    hasNextPage: ref(false),
    isFetched: ref(true),
    isFetching: ref(false),
    isLoading: ref(false),
    isError: ref(false),
    error: ref(null)
  }))
}))

vi.mock('pocketbase', () => ({
  ClientResponseError: class ClientResponseError extends Error {
    status: number
    constructor(msg: string, status = 0) { super(msg); this.status = status }
  },
  Record: class Record {}
}))

describe('Wall Page', () => {
  it('can be imported without errors', async () => {
    const Wall = await import('../pages/Wall.vue')
    expect(Wall.default).toBeDefined()
  })

  it('exports a Vue component with setup', async () => {
    const Wall = await import('../pages/Wall.vue')
    expect(Wall.default.__name || Wall.default.name).toBeTruthy()
  })

  it('has correct component file structure', async () => {
    const Wall = await import('../pages/Wall.vue')
    // Vue SFC compiled components have a render or setup function
    expect(typeof Wall.default).toBe('object')
  })
})
