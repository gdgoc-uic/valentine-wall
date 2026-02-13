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
  // Helper that simulates the notification logic used in Home.vue and Wall.vue
  function simulateIncomingMessage(
    authState: { isLoggedIn: boolean; user: any },
    record: { recipient: string; gifts?: string[] }
  ) {
    if (authState.isLoggedIn && authState.user?.expand?.details?.student_id &&
        record.recipient === authState.user.expand.details.student_id) {
      const hasGifts = record.gifts && record.gifts.length > 0
      notify({
        type: 'success',
        text: hasGifts
          ? '🎁 You received a new message with a gift!'
          : '💌 You received a new message!'
      } as any)
    }
  }

  it('should notify user when they receive a new message', () => {
    const authState = {
      isLoggedIn: true,
      user: { expand: { details: { student_id: '202012345678' } } }
    }

    simulateIncomingMessage(authState, { recipient: '202012345678', gifts: [] })

    expect(vi.mocked(notiwind.notify)).toHaveBeenCalledWith({
      type: 'success',
      text: '💌 You received a new message!'
    })
  })

  it('should notify with gift text when message has gifts', () => {
    const authState = {
      isLoggedIn: true,
      user: { expand: { details: { student_id: '202012345678' } } }
    }

    simulateIncomingMessage(authState, { recipient: '202012345678', gifts: ['gift1'] })

    expect(vi.mocked(notiwind.notify)).toHaveBeenCalledWith({
      type: 'success',
      text: '🎁 You received a new message with a gift!'
    })
  })

  it('should not notify for "everyone" messages', () => {
    const authState = {
      isLoggedIn: true,
      user: { expand: { details: { student_id: '202012345678' } } }
    }

    vi.clearAllMocks()
    simulateIncomingMessage(authState, { recipient: 'everyone', gifts: [] })

    expect(vi.mocked(notiwind.notify)).not.toHaveBeenCalled()
  })

  it('should not notify when message is for different recipient', () => {
    const authState = {
      isLoggedIn: true,
      user: { expand: { details: { student_id: '202012345678' } } }
    }

    vi.clearAllMocks()
    simulateIncomingMessage(authState, { recipient: '202087654321', gifts: [] })

    expect(vi.mocked(notiwind.notify)).not.toHaveBeenCalled()
  })

  it('should not notify when user is not logged in', () => {
    const authState = { isLoggedIn: false, user: null }

    vi.clearAllMocks()
    simulateIncomingMessage(authState, { recipient: '202012345678', gifts: [] })

    expect(vi.mocked(notiwind.notify)).not.toHaveBeenCalled()
  })
})
