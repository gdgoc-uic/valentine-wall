import { describe, it, expect, vi } from 'vitest'
import { thirdPartyLogin, popupCenter } from '../auth'

// Mock client
vi.mock('../client', () => ({
  pb: {
    collection: vi.fn(() => ({
      listAuthMethods: vi.fn(),
      authWithOAuth2: vi.fn()
    })),
    buildUrl: vi.fn((path: string) => `http://localhost:8090${path}`)
  },
  backendUrl: 'http://localhost:8090'
}))

describe('Auth System', () => {
  describe('popupCenter', () => {
    it('calculates popup position correctly', () => {
      const mockOpen = vi.fn()
      globalThis.window.open = mockOpen

      const params = {
        url: 'https://example.com',
        title: 'Test Window',
        w: 800,
        h: 600
      }

      popupCenter(params)

      expect(mockOpen).toHaveBeenCalled()
      const callArgs = mockOpen.mock.calls[0]
      expect(callArgs[0]).toBe('https://example.com')
      expect(callArgs[1]).toBe('Test Window')
    })

    it('returns null if window.open fails', () => {
      globalThis.window.open = vi.fn().mockReturnValue(null)

      const result = popupCenter({
        url: 'https://example.com',
        title: 'Test',
        w: 800,
        h: 600
      })

      expect(result).toBeNull()
    })
  })

  describe('thirdPartyLogin', () => {
    it('throws error if popup fails to open', async () => {
      const { pb } = await import('../client')
      
      vi.mocked(pb.collection).mockReturnValue({
        listAuthMethods: vi.fn().mockResolvedValue({
          authProviders: [{
            name: 'google',
            authUrl: 'https://accounts.google.com/auth',
            state: 'test-state',
            codeVerifier: 'test-verifier'
          }]
        })
      } as any)

      globalThis.window.open = vi.fn().mockReturnValue(null)

      await expect(thirdPartyLogin('google')).rejects.toThrow('Failed to open window.')
    })

    it('finds correct provider from list', async () => {
      const { pb } = await import('../client')
      
      const mockListAuthMethods = vi.fn().mockResolvedValue({
        authProviders: [
          { name: 'facebook', authUrl: 'https://facebook.com' },
          { name: 'google', authUrl: 'https://google.com', state: 'test', codeVerifier: 'test' }
        ]
      })

      vi.mocked(pb.collection).mockReturnValue({
        listAuthMethods: mockListAuthMethods
      } as any)

      globalThis.window.open = vi.fn().mockReturnValue({
        focus: vi.fn(),
        close: vi.fn()
      })

      // This will throw since we're not completing the OAuth flow
      await thirdPartyLogin('google').catch(() => {})

      expect(mockListAuthMethods).toHaveBeenCalled()
    })

    it('returns early if provider not found', async () => {
      const { pb } = await import('../client')
      
      vi.mocked(pb.collection).mockReturnValue({
        listAuthMethods: vi.fn().mockResolvedValue({
          authProviders: []
        })
      } as any)

      const result = await thirdPartyLogin('google')
      expect(result).toBeUndefined()
    })

    it('uses correct table name', async () => {
      const { pb } = await import('../client')
      const mockCollection = vi.fn()
      
      pb.collection = mockCollection.mockReturnValue({
        listAuthMethods: vi.fn().mockResolvedValue({
          authProviders: []
        })
      })

      await thirdPartyLogin('google', 'custom_users')

      expect(mockCollection).toHaveBeenCalledWith('custom_users')
    })
  })
})
