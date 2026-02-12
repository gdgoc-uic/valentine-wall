import { describe, it, expect, vi, beforeEach } from 'vitest'
import fs from 'fs'
import path from 'path'

// =============================================================================
// Bug Fix Tests
// =============================================================================
// Tests for the following fixes:
// 1. Message.vue: suppress duplicate toast when #error slot handles the error
// 2. Wall.vue: fix 'receipientId' typo → 'recipientId'
// 3. Wall.vue: suppress duplicate toast when #error slot handles the error
// 4. Lemuel.jpg: exists in public/about/ for the footer team section
// =============================================================================

vi.mock('notiwind', () => ({
  notify: vi.fn()
}))

vi.mock('../analytics', () => ({
  logEvent: vi.fn()
}))

const readSource = (relPath: string) =>
  fs.readFileSync(path.resolve(__dirname, relPath), 'utf-8')

// =========================================================
// Test 1: Message.vue suppresses global error toast
// =========================================================

describe('Bug Fix: Message.vue duplicate error toast', () => {
  const source = readSource('../pages/Message.vue')

  it('useQuery has onError override to suppress global toast', () => {
    // The useQuery call should include onError: () => {} so the global
    // catchAndNotifyError from main.ts does NOT fire alongside the
    // inline #error template that already shows "Message not found."
    expect(source).toContain('onError: () => {}')
  })

  it('still renders inline error for 404 in the template', () => {
    expect(source).toContain('Message not found.')
    expect(source).toContain('Double-check if your link is correct and try again.')
  })

  it('has retry: 0 to fail fast on missing messages', () => {
    expect(source).toContain('retry: 0')
  })
})

// =========================================================
// Test 2: Wall.vue fixes
// =========================================================

describe('Bug Fix: Wall.vue receipientId typo', () => {
  const source = readSource('../pages/Wall.vue')

  it('uses correct "recipientId" (not "receipientId")', () => {
    // The typo "receipientId" should NOT exist anywhere in the source
    expect(source).not.toContain('receipientId')

    // The correct param name should be used for the hasGift initialization
    expect(source).toContain('route.params.recipientId')
  })
})

describe('Bug Fix: Wall.vue duplicate error toast', () => {
  const source = readSource('../pages/Wall.vue')

  it('useInfiniteQuery has onError override to suppress global toast', () => {
    // Wall.vue already renders "Nothing to see here!" inline for empty results
    // so the global catchAndNotifyError toast should be suppressed
    expect(source).toContain('onError: () => {}')
  })

  it('still renders inline error for empty results', () => {
    expect(source).toContain('Nothing to see here!')
  })
})

// =========================================================
// Test 3: Lemuel.jpg existence in public/about/
// =========================================================

describe('Bug Fix: Lemuel.jpg exists in public/about/', () => {
  it('Lemuel.jpg file exists and is a valid non-empty file', () => {
    const imagePath = path.resolve(__dirname, '../../public/about/Lemuel.jpg')
    expect(fs.existsSync(imagePath)).toBe(true)

    const stats = fs.statSync(imagePath)
    expect(stats.size).toBeGreaterThan(0)
  })

  it('about.json includes Lemuel in members list', async () => {
    const aboutData = await import('../assets/about.json')
    const members = (aboutData as any).default?.members || (aboutData as any).members
    expect(members).toBeDefined()

    const lemuel = members.find((m: any) => m.name.includes('Lemuel'))
    expect(lemuel).toBeDefined()
    expect(lemuel.name.split(' ')[0]).toBe('Lemuel')
  })

  it('all team member images exist on disk', async () => {
    const aboutData = await import('../assets/about.json')
    const members = (aboutData as any).default?.members || (aboutData as any).members

    const missingImages: string[] = []
    for (const member of members) {
      const firstName = member.name.split(' ')[0]
      const imagePath = path.resolve(__dirname, `../../public/about/${firstName}.jpg`)
      if (!fs.existsSync(imagePath)) {
        missingImages.push(`${firstName}.jpg`)
      }
    }

    expect(missingImages).toEqual([])
  })
})

// =========================================================
// Test 4: App.vue router-view has :key for proper navigation
// =========================================================

describe('Bug Fix: App.vue navigation between routes', () => {
  const source = readSource('../App.vue')

  it('router-view destructures route from slot', () => {
    // Must have v-slot with route destructuring for keying
    expect(source).toMatch(/router-view\s+v-slot\s*=\s*"\{[^}]*route[^}]*\}"/)
  })

  it('component inside suspense has :key bound to route path', () => {
    // The <component> must have :key="...route.path" or :key="viewRoute.path"
    // so Vue destroys and recreates the component when the route changes
    expect(source).toMatch(/:key\s*=\s*"[^"]*\.path"/)
  })

  it('uses suspense to wrap async page components', () => {
    expect(source).toContain('<suspense>')
    expect(source).toContain('</suspense>')
  })
})

// =========================================================
// Test 5: Error handling behavior
// =========================================================

describe('Bug Fix: Global vs component error handling', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('catchAndNotifyError shows toast notification for errors', async () => {
    const { catchAndNotifyError } = await import('../notify')
    const { notify: notiwindNotify } = await import('notiwind')
    vi.mocked(notiwindNotify).mockClear()

    catchAndNotifyError(new Error('Some API error'))

    expect(vi.mocked(notiwindNotify)).toHaveBeenCalledWith({
      type: 'error',
      text: 'Some API error'
    })
  })

  it('Message.vue has both inline error handling and toast suppression', () => {
    const source = readSource('../pages/Message.vue')

    // Has inline error template
    expect(source).toContain('Message not found.')
    expect(source).toContain('Something went wrong.')

    // Has toast suppression
    expect(source).toContain('onError: () => {}')

    // Has retry: 0
    expect(source).toContain('retry: 0')
  })
})
