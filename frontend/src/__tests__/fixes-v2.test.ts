import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive, ref } from 'vue'
import GiftIcon from '../components/GiftIcon.vue'

/**
 * Tests for the 5 bug fixes:
 * 1. Gift icons are bigger (text-5xl, 7.5rem wrapper)
 * 2. Recent messages only show "everyone" messages (private filtered out)
 * 3. Mobile responsiveness - recipient input & Everyone button don't overlap
 * 4. Sidebar is narrower on mobile (70vw max-w-300px)
 * 5. Wall.vue recent view filters by recipient="everyone"
 */

// ── Fix 1: Gift icon display ──────────────────────────────────────────

describe('Fix 1: Gift icons are visible and properly sized', () => {
  it('renders emoji with inline-block display for proper sizing', () => {
    const wrapper = mount(GiftIcon, { props: { uid: 'rose' } })
    const span = wrapper.find('span')
    expect(span.attributes('style')).toContain('display: inline-block')
    expect(span.attributes('style')).toContain('font-size: inherit')
  })

  it('renders all standard gift emojis correctly', () => {
    const gifts = [
      { uid: 'rose', emoji: '🌹' },
      { uid: 'chocolate', emoji: '🍫' },
      { uid: 'teddy', emoji: '🧸' },
      { uid: 'flowers', emoji: '💐' },
      { uid: 'candy', emoji: '🍬' },
      { uid: 'card', emoji: '💌' },
      { uid: 'balloon', emoji: '🎈' },
      { uid: 'money_100', emoji: '💵' },
      { uid: 'money_500', emoji: '💰' },
    ]

    gifts.forEach(({ uid, emoji }) => {
      const wrapper = mount(GiftIcon, { props: { uid } })
      expect(wrapper.text()).toBe(emoji)
    })
  })

  it('defaults to gift box emoji for unknown UIDs', () => {
    const wrapper = mount(GiftIcon, { props: { uid: 'nonexistent' } })
    expect(wrapper.text()).toBe('🎁')
  })
})

// ── Fix 2: Message visibility - everyone vs private ───────────────────

describe('Fix 2: Recent messages only show "everyone" messages', () => {
  it('Home.vue filter uses recipient = "everyone" not gifts filter', async () => {
    // Read the actual Home.vue source to verify the filter
    const fs = await import('fs')
    const path = await import('path')
    const homeSource = fs.readFileSync(
      path.resolve(__dirname, '../pages/Home.vue'), 'utf-8'
    )
    
    // Should filter by recipient = "everyone"
    expect(homeSource).toContain('recipient = "everyone"')
    // Should NOT use the old gifts-only filter for recent messages
    expect(homeSource).not.toContain("filter: 'gifts = \"[]\"'")
  })

  it('Home.vue SSE subscription filters by recipient "everyone"', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const homeSource = fs.readFileSync(
      path.resolve(__dirname, '../pages/Home.vue'), 'utf-8'
    )
    
    // SSE should check recipient is everyone
    expect(homeSource).toContain("e.record.recipient !== 'everyone'")
  })

  it('Wall.vue recent wall uses recipient = "everyone" filter', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const wallSource = fs.readFileSync(
      path.resolve(__dirname, '../pages/Wall.vue'), 'utf-8'
    )
    
    // The fallback filter for no-recipient (recent wall) should be everyone
    expect(wallSource).toContain('recipient = "everyone"')
    // Should not use old gifts:length filter for recent view
    expect(wallSource).not.toContain("'gifts:length = 0'")
  })

  it('Wall.vue SSE filters non-everyone from recent wall', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const wallSource = fs.readFileSync(
      path.resolve(__dirname, '../pages/Wall.vue'), 'utf-8'
    )
    
    // SSE should filter out non-everyone for recent wall
    expect(wallSource).toContain("data.record.recipient !== 'everyone'")
  })
})

// ── Fix 3: Mobile responsiveness ─────────────────────────────────────

describe('Fix 3: Recipient input and Everyone button responsiveness', () => {
  it('SendMessageForm uses responsive flex layout for recipient row', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const formSource = fs.readFileSync(
      path.resolve(__dirname, '../components/SendMessageForm.vue'), 'utf-8'
    )
    
    // Should use flex-col on small screens, flex-row on sm+
    expect(formSource).toContain('flex-col sm:flex-row')
    // Input should have min-w-0 to prevent overflow
    expect(formSource).toContain('min-w-0')
    // Everyone button should be smaller on mobile
    expect(formSource).toContain('btn-sm sm:btn-md')
  })

  it('Search bar has min-w-0 to prevent overflow', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const navSource = fs.readFileSync(
      path.resolve(__dirname, '../components/Navbar.vue'), 'utf-8'
    )
    
    expect(navSource).toContain('min-w-0')
  })
})

// ── Fix 4: Sidebar is narrower on mobile ─────────────────────────────

describe('Fix 4: Mobile sidebar is compact', () => {
  it('Sidebar uses constrained width', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const navSource = fs.readFileSync(
      path.resolve(__dirname, '../components/Navbar.vue'), 'utf-8'
    )
    
    // Should use constrained width
    expect(navSource).toContain('max-w-[320px]')
    // Should NOT have old 85vw
    expect(navSource).not.toContain('w-[85vw]')
  })

  it('Sidebar has proper navigation links', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const navSource = fs.readFileSync(
      path.resolve(__dirname, '../components/Navbar.vue'), 'utf-8'
    )
    
    // Sidebar should have nav links with hover states
    expect(navSource).toContain('hover:bg-rose-50')
  })

  it('Sidebar has user section with avatar', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const navSource = fs.readFileSync(
      path.resolve(__dirname, '../components/Navbar.vue'), 'utf-8'
    )
    
    // User section should show username with avatar
    expect(navSource).toContain('authState.user.username')
    expect(navSource).toContain('charAt(0)')
  })
})

// ── Fix 5: Gift icon wrapper height ──────────────────────────────────

describe('Fix 5: Gift selector has proper sizing', () => {
  it('Gift wrapper height is 7.5rem for larger icons', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const formSource = fs.readFileSync(
      path.resolve(__dirname, '../components/SendMessageForm.vue'), 'utf-8'
    )
    
    expect(formSource).toContain('height: 7.5rem')
    expect(formSource).not.toContain('height: 6.3rem')
  })

  it('Gift icon uses text-5xl class', async () => {
    const fs = await import('fs')
    const path = await import('path')
    const formSource = fs.readFileSync(
      path.resolve(__dirname, '../components/SendMessageForm.vue'), 'utf-8'
    )
    
    // In the gift modal
    expect(formSource).toContain('class="text-5xl"')
  })
})
