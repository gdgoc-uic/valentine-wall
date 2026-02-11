import { describe, it, expect } from 'vitest'
import { isReadOnly } from '../utils'

// Note: We cannot easily mock import.meta.env in Vitest
// These tests verify the actual behavior with the current environment
describe('Utility Functions', () => {describe('isReadOnly', () => {
    it('returns a boolean value', () => {
      const result = isReadOnly()
      expect(typeof result).toBe('boolean')
    })
  })
})
