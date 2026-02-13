import { describe, it, expect } from 'vitest'
import * as fs from 'fs'
import * as path from 'path'

const homeSrc = fs.readFileSync(
  path.resolve(__dirname, '../pages/Home.vue'),
  'utf-8'
)

const wallSrc = fs.readFileSync(
  path.resolve(__dirname, '../pages/Wall.vue'),
  'utf-8'
)

describe('Message Toast Notifications - Source Integration', () => {
  describe('Home.vue', () => {
    it('should import notify', () => {
      expect(homeSrc).toContain("import { notify } from '../notify'")
    })

    it('should check if user is logged in before notifying', () => {
      expect(homeSrc).toContain('authState.isLoggedIn')
      expect(homeSrc).toContain('authState.user?.expand?.details?.student_id')
    })

    it('should show different toast text for messages with gifts', () => {
      expect(homeSrc).toContain('You received a new message with a gift!')
      expect(homeSrc).toContain('You received a new message!')
    })

    it('should check for gifts in the record', () => {
      expect(homeSrc).toContain('e.record.gifts && e.record.gifts.length > 0')
    })
  })

  describe('Wall.vue', () => {
    it('should import notify', () => {
      expect(wallSrc).toContain("import { notify } from '../notify'")
    })

    it('should check if user is logged in before notifying', () => {
      expect(wallSrc).toContain('authState.isLoggedIn')
      expect(wallSrc).toContain('authState.user?.expand?.details?.student_id')
    })

    it('should show different toast text for messages with gifts', () => {
      expect(wallSrc).toContain('You received a new message with a gift!')
      expect(wallSrc).toContain('You received a new message!')
    })

    it('should check for gifts in the record', () => {
      expect(wallSrc).toContain('data.record.gifts && data.record.gifts.length > 0')
    })
  })
})
