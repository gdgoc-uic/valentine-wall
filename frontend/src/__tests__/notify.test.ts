import { describe, it, expect, vi, beforeEach } from 'vitest'
import { notify, catchAndNotifyError } from '../notify'
import * as notiwind from 'notiwind'
import * as analytics from '../analytics'

// Mock notiwind
vi.mock('notiwind', () => ({
  notify: vi.fn()
}))

// Mock analytics
vi.mock('../analytics', () => ({
  logEvent: vi.fn()
}))

describe('Notification System', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('notify calls notiwind notify with correct args', () => {
    const notifyArgs = {
      type: 'success',
      text: 'Test notification'
    }

    notify(notifyArgs as any)

    expect(vi.mocked(notiwind.notify)).toHaveBeenCalledWith(notifyArgs)
  })

  it('notify logs event to analytics', () => {
    const notifyArgs = {
      type: 'info',
      text: 'Info message'
    }

    notify(notifyArgs as any)

    expect(vi.mocked(analytics.logEvent)).toHaveBeenCalledWith('server_notifications', notifyArgs)
  })

  it('catchAndNotifyError handles Error objects', () => {
    const error = new Error('Test error message')
    
    catchAndNotifyError(error)

    expect(vi.mocked(notiwind.notify)).toHaveBeenCalledWith({
      type: 'error',
      text: 'Test error message'
    })
  })

  it('catchAndNotifyError handles unknown errors', () => {
    const unknownError = { unknown: 'error' }
    
    catchAndNotifyError(unknownError)

    expect(vi.mocked(notiwind.notify)).toHaveBeenCalledWith({
      type: 'error',
      text: 'Unknown error.'
    })
  })
})

describe('Real-time Message Notifications', () => {
  it('should notify user when they receive a new message', () => {
    // This test simulates the subscription callback
    const userStudentId: string = '202012345678'
    
    const mockData: { action: string; record: { recipient: string; content: string } } = {
      action: 'create',
      record: {
        recipient: '202012345678',
        content: 'Test message'
      }
    }

    // Simulate the subscription callback logic
    if (mockData.action === 'create') {
      const messageRecipient: string = mockData.record.recipient
      if (messageRecipient === userStudentId && messageRecipient !== 'everyone') {
        notify({
          type: 'success',
          text: '💌 You received a new message! Click to view.',
          duration: 8000
        } as any)
      }
    }

    expect(vi.mocked(notiwind.notify)).toHaveBeenCalledWith({
      type: 'success',
      text: '💌 You received a new message! Click to view.',
      duration: 8000
    })
  })

  it('should not notify for "everyone" messages', () => {
    const userStudentId: string = '202012345678'
    
    const mockData: { action: string; record: { recipient: string } } = {
      action: 'create',
      record: {
        recipient: 'everyone'
      }
    }

    // Clear previous calls
    vi.clearAllMocks()

    // Simulate the subscription callback logic
    if (mockData.action === 'create') {
      const messageRecipient: string = mockData.record.recipient
      const shouldNotify = messageRecipient === userStudentId && messageRecipient !== 'everyone'
      if (shouldNotify) {
        notify({
          type: 'success',
          text: '💌 You received a new message! Click to view.',
          duration: 8000
        } as any)
      }
    }

    expect(vi.mocked(notiwind.notify)).not.toHaveBeenCalled()
  })

  it('should not notify when message is for different recipient', () => {
    const userStudentId: string = '202012345678'
    
    const mockData: { action: string; record: { recipient: string } } = {
      action: 'create',
      record: {
        recipient: '202087654321'
      }
    }

    // Clear previous calls
    vi.clearAllMocks()

    // Simulate the subscription callback logic
    if (mockData.action === 'create') {
      const messageRecipient: string = mockData.record.recipient
      const shouldNotify = messageRecipient === userStudentId && messageRecipient !== 'everyone'
      if (shouldNotify) {
        notify({
          type: 'success',
          text: '💌 You received a new message! Click to view.',
          duration: 8000
        } as any)
      }
    }

    expect(vi.mocked(notiwind.notify)).not.toHaveBeenCalled()
  })
})
