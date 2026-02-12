import { describe, it, expect } from 'vitest'
import fs from 'fs'
import path from 'path'

/**
 * Tests for v3 fixes:
 * 1. Navbar/search bar mobile overlap fix
 * 2. Sent Messages tab on Wall
 * 3. Reply system fixes (SSE cleanup, backend typo)
 * 4. Currency symbol (₱ instead of ღ)
 * 5. Text overflow (break-words)
 * 6. Sidebar redesign
 */

function readComponent(relativePath: string): string {
  return fs.readFileSync(path.resolve(__dirname, '..', relativePath), 'utf-8')
}

function readBackend(relativePath: string): string {
  return fs.readFileSync(path.resolve(__dirname, '../../../backend', relativePath), 'utf-8')
}

// ── 1: Navbar mobile overlap fix ──────────────────────────────────────

describe('Navbar mobile responsiveness', () => {
  it('logo image uses explicit dimensions instead of h-full', () => {
    const src = readComponent('components/Navbar.vue')
    // Should NOT use h-full which causes unpredictable sizing
    expect(src).not.toMatch(/class="h-full w-14/)
    // Should use explicit small dimensions on mobile
    expect(src).toContain('w-8 h-8')
  })

  it('hamburger button uses compact btn-sm on mobile', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('btn-sm')
  })

  it('logo link is a flex container on all breakpoints', () => {
    const src = readComponent('components/Navbar.vue')
    // Should have "flex items-center" in base classes
    expect(src).toMatch(/class="flex-none flex items-center/)
  })

  it('search form has min-w-0 to prevent overflow', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('min-w-0')
  })
})

// ── 2: Sent Messages tab ──────────────────────────────────────────────

describe('Sent Messages feature on Wall', () => {
  it('Wall.vue has isSentView ref', () => {
    const src = readComponent('pages/Wall.vue')
    expect(src).toContain('isSentView')
  })

  it('Wall.vue has a Sent tab button', () => {
    const src = readComponent('pages/Wall.vue')
    expect(src).toContain('>Sent</span>')
  })

  it('Query key includes isSentView for proper cache separation', () => {
    const src = readComponent('pages/Wall.vue')
    expect(src).toContain("['wall', recipient, hasGift, isSentView]")
  })

  it('Sent tab filters by user (sender) field', () => {
    const src = readComponent('pages/Wall.vue')
    // When in sent view, filter by user = sender's details ID
    expect(src).toContain('user = "')
    expect(src).toContain('authState.user!.details')
  })

  it('SSE subscription handles sent view for real-time updates', () => {
    const src = readComponent('pages/Wall.vue')
    expect(src).toContain('isSentView.value')
    expect(src).toContain('data.record.user === authState.user!.details')
  })

  it('All/Messages/Gifts tabs reset isSentView to false', () => {
    const src = readComponent('pages/Wall.vue')
    // All tab
    expect(src).toContain('isSentView = false; hasGift = null')
    // Messages tab
    expect(src).toContain('isSentView = false; hasGift = false')
    // Gifts tab
    expect(src).toContain('isSentView = false; hasGift = true')
  })

  it('Tabs use responsive sizing for mobile', () => {
    const src = readComponent('pages/Wall.vue')
    expect(src).toContain('tab-md sm:tab-lg')
  })
})

// ── 3: Reply system fixes ─────────────────────────────────────────────

describe('Reply system fixes', () => {
  it('ReplyThread SSE subscription is properly assigned for cleanup', () => {
    const src = readComponent('components/ReplyThread.vue')
    // Should assign the subscribe result to unsubscribeFunc
    expect(src).toContain('unsubscribeFunc.value = await pb.collection')
  })

  it('ReplyThread uses async onMounted for proper await', () => {
    const src = readComponent('components/ReplyThread.vue')
    expect(src).toContain('onMounted(async ()')
  })

  it('Backend message_replies hook uses correct collection name', () => {
    const src = readBackend('main.go')
    // Should use "message_replies" not "message_repies"
    expect(src).not.toContain('message_repies')
    // The before-create hook should reference the correct name
    expect(src).toMatch(/case "message_replies":\s*\n\s*return onBeforeAddMessageReply/)
  })

  it('Reply enabled check allows message sender to see replies', () => {
    const src = readComponent('components/ReplyThread.vue')
    // The enabled computed should check message.user == authState.user.id
    expect(src).toContain('message.value.user == authState.user.id')
  })
})

// ── 4: Currency symbol ₱ ──────────────────────────────────────────────

describe('Currency symbol is PHP Peso (₱)', () => {
  it('Navbar uses ₱ not ღ', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('₱')
    expect(src).not.toContain('ღ')
  })

  it('SendMessageForm uses ₱', () => {
    const src = readComponent('components/SendMessageForm.vue')
    expect(src).toContain('₱')
    expect(src).not.toContain('ღ')
  })

  it('ReplyMessageBox uses ₱', () => {
    const src = readComponent('components/ReplyMessageBox.vue')
    expect(src).toContain('₱')
    expect(src).not.toContain('ღ')
  })

  it('RankingsBox uses ₱', () => {
    const src = readComponent('components/RankingsBox.vue')
    expect(src).toContain('₱')
    expect(src).not.toContain('ღ')
  })

  it('Rankings page uses ₱', () => {
    const src = readComponent('pages/Rankings.vue')
    expect(src).toContain('₱')
    expect(src).not.toContain('ღ')
  })
})

// ── 5: Text overflow fix ──────────────────────────────────────────────

describe('Text overflow prevention (break-words)', () => {
  it('Message.vue content has break-words class', () => {
    const src = readComponent('pages/Message.vue')
    expect(src).toContain('text-4xl break-words')
  })

  it('ReplyThread reply text has break-words class', () => {
    const src = readComponent('components/ReplyThread.vue')
    expect(src).toContain('text-lg break-words')
  })

  it('MessageTiles content has overflow-hidden and break-words', () => {
    const src = readComponent('components/MessageTiles.vue')
    expect(src).toContain('overflow-hidden')
    expect(src).toContain('break-words')
  })
})

// ── 6: Sidebar redesign ──────────────────────────────────────────────

describe('Sidebar redesign', () => {
  it('Sidebar has gradient header with branding', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('bg-gradient-to-r from-rose-500 to-rose-400')
    expect(src).toContain('Valentine Wall')
  })

  it('Sidebar has user avatar initial', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain("charAt(0)")
    expect(src).toContain('toUpperCase()')
  })

  it('Sidebar has slide animation', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('sidebar-slide')
    expect(src).toContain('translateX(-100%)')
  })

  it('Sidebar overlay uses backdrop blur', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('backdrop-blur-sm')
  })

  it('Sidebar nav links have hover states', () => {
    const src = readComponent('components/Navbar.vue')
    expect(src).toContain('hover:bg-rose-50 hover:text-rose-600')
  })
})

// ── 7: Horizontal overflow prevention ─────────────────────────────────

describe('Horizontal overflow prevention', () => {
  it('App.vue has overflow-x hidden on html and body', () => {
    const src = readComponent('App.vue')
    expect(src).toContain('overflow-x: hidden')
  })

  it('App.vue has overscroll-behavior-x: none', () => {
    const src = readComponent('App.vue')
    expect(src).toContain('overscroll-behavior-x: none')
  })

  it('App.vue #app has overflow-x: clip', () => {
    const src = readComponent('App.vue')
    expect(src).toContain('overflow-x: clip')
  })

  it('index.html has viewport meta tag', () => {
    const indexHtml = fs.readFileSync(
      path.resolve(__dirname, '../../index.html'), 'utf-8'
    )
    expect(indexHtml).toContain('viewport')
    expect(indexHtml).toContain('width=device-width')
  })
})
