import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import GiftIcon from '../components/GiftIcon.vue'

describe('GiftIcon Component', () => {
  it('renders default gift emoji for unknown UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'unknown_gift' }
    })
    
    expect(wrapper.text()).toBe('🎁')
  })

  it('renders rose emoji for rose UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'rose' }
    })
    
    expect(wrapper.text()).toBe('🌹')
  })

  it('renders chocolate emoji for chocolate UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'chocolate' }
    })
    
    expect(wrapper.text()).toBe('🍫')
  })

  it('renders teddy emoji for teddy UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'teddy' }
    })
    
    expect(wrapper.text()).toBe('🧸')
  })

  it('renders money emoji for money_100 UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'money_100' }
    })
    
    expect(wrapper.text()).toBe('💵')
  })

  it('renders money bag emoji for money_500 UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'money_500' }
    })
    
    expect(wrapper.text()).toBe('💰')
  })

  it('handles whitespace in UID', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: '  rose  ' }
    })
    
    expect(wrapper.text()).toBe('🌹')
  })

  it('renders legacy gift UIDs correctly', () => {
    const testCases = [
      { uid: 'sigenapls', emoji: '🌹' },
      { uid: 'isforu', emoji: '💝' },
      { uid: 'timberlake', emoji: '🌹' },
      { uid: 'mukuha', emoji: '😺' },
      { uid: 'rizzler', emoji: '😎' }
    ]

    testCases.forEach(({ uid, emoji }) => {
      const wrapper = mount(GiftIcon, {
        props: { uid }
      })
      expect(wrapper.text()).toBe(emoji)
    })
  })

  it('applies custom styles', () => {
    const wrapper = mount(GiftIcon, {
      props: { uid: 'rose' }
    })
    
    const span = wrapper.find('span')
    expect(span.attributes('style')).toContain('font-size: inherit')
    expect(span.attributes('style')).toContain('line-height: 1')
  })
})
