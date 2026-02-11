import { vi, beforeEach } from 'vitest'
import { config } from '@vue/test-utils'
import { defineComponent, h } from 'vue'

// Mock .md imports (vite-plugin-markdown)
vi.mock('../assets/texts/rules.md', () => ({
  VueComponent: defineComponent({ name: 'RulesContent', render: () => h('div', 'Mocked rules content') })
}))

// Mock floating-vue
vi.mock('floating-vue', () => ({
  Tooltip: defineComponent({ name: 'Tooltip', render() { return h('div', (this as any).$slots.default?.()) } })
}))

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock IntersectionObserver
globalThis.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  takeRecords() {
    return []
  }
  unobserve() {}
} as any

// Mock localStorage
const storage: Record<string, string> = {}
const localStorageMock = {
  getItem: vi.fn((key: string) => storage[key] || null),
  setItem: vi.fn((key: string, value: string) => { storage[key] = value }),
  removeItem: vi.fn((key: string) => { delete storage[key] }),
  clear: vi.fn(() => { Object.keys(storage).forEach(key => delete storage[key]) }),
}
globalThis.localStorage = localStorageMock as any

// Global test configuration
config.global.stubs = {
  Teleport: true,
  ClientOnly: true
}

// Reset mocks before each test
beforeEach(() => {
  vi.clearAllMocks()
  Object.keys(storage).forEach(key => delete storage[key])
})
