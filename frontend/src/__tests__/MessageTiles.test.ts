import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MessageTiles from '../components/MessageTiles.vue'

describe('MessageTiles Component', () => {
  const mockMessages = [
    {
      id: 'msg1',
      content: 'Test message 1',
      recipient: '202012345678',
      gifts: [],
      paperColor: 1,
      created: new Date().toISOString(),
      expand: {
        gifts: []
      }
    },
    {
      id: 'msg2',
      content: 'Test message 2',
      recipient: 'everyone',
      gifts: ['gift1'],
      paperColor: 2,
      created: new Date().toISOString(),
      expand: {
        gifts: [
          { id: 'gift1', uid: 'rose', label: 'Rose', price: 50 }
        ]
      }
    }
  ]

  const globalConfig = {
    stubs: {
      'router-link': { template: '<a><slot/></a>' },
      'masonry': { template: '<div><slot/></div>' },
      'gift-icon': { template: '<span />' }
    }
  }

  it('renders message tiles', () => {
    const wrapper = mount(MessageTiles, {
      props: {
        messages: mockMessages,
        limit: 10
      },
      global: globalConfig
    })

    expect(wrapper.exists()).toBe(true)
  })

  it('respects limit prop', () => {
    const wrapper = mount(MessageTiles, {
      props: {
        messages: mockMessages,
        limit: 1
      },
      global: globalConfig
    })

    const papers = wrapper.findAll('.message-paper-wrapper')
    expect(papers.length).toBe(1)
  })

  it('handles empty messages array', () => {
    const wrapper = mount(MessageTiles, {
      props: {
        messages: [],
        limit: 10
      },
      global: globalConfig
    })

    expect(wrapper.exists()).toBe(true)
  })

  it('displays message content', () => {
    const wrapper = mount(MessageTiles, {
      props: {
        messages: [mockMessages[0]],
        limit: 10
      },
      global: globalConfig
    })

    expect(wrapper.text()).toContain('Test message 1')
  })
})
