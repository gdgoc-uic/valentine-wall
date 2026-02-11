# Valentine Wall - Testing Guide

## Overview

This document provides comprehensive information about the test suite for the Valentine Wall application, including setup instructions, test coverage, and running tests.

## Test Structure

### Backend Tests (Go)
Located in: `backend/*_test.go`

- **message_hooks_test.go** - Tests for message-related hooks and operations
- **virtual_wallet_test.go** - Tests for virtual wallet and transaction system
- **user_hooks_test.go** - Tests for user lifecycle hooks
- **server_test.go** - Tests for API endpoints

### Frontend Tests (TypeScript/Vitest)
Located in: `frontend/src/__tests__/`

- **GiftIcon.test.ts** - Tests for gift emoji icon component
- **SendMessageForm.test.ts** - Tests for message sending form
- **MessageTiles.test.ts** - Tests for message display component
- **Wall.test.ts** - Tests for wall page component
- **store.test.ts** - Tests for Vuex store and state management
- **auth.test.ts** - Tests for authentication system
- **notify.test.ts** - Tests for notification system (including real-time notifications)
- **utils.test.ts** - Tests for utility functions
- **setup.ts** - Test environment setup and mocks

## Test Coverage by Feature

### ✅ Recently Implemented Features

#### 1. Gift Icons with Emoji Support
**File:** `frontend/src/__tests__/GiftIcon.test.ts`

Tests cover:
- Default emoji rendering for unknown gifts
- Specific emoji mapping for all gift types (rose, chocolate, teddy, etc.)
- Whitespace handling in UIDs
- Legacy gift UID support
- Style attribute application

**Run:** `cd frontend && npm run test -- GiftIcon`

#### 2. Real-time Message Notifications
**File:** `frontend/src/__tests__/notify.test.ts`

Tests cover:
- Notification display when user receives message
- Filtering out "everyone" messages
- Only notifying correct recipient
- Notification cleanup on logout
- SSR compatibility

**Run:** `cd frontend && npm run test -- notify`

### ✅ Core Feature Tests

#### 3. User Authentication Flow
**Files:** 
- `frontend/src/__tests__/auth.test.ts`
- `frontend/src/__tests__/store.test.ts`
- `backend/user_hooks_test.go`

Tests cover:
- Google OAuth login flow
- Popup window handling
- Session management
- User details creation
- Login/logout state management
- Message subscription on login
- Cleanup on logout

**Run Frontend:** `cd frontend && npm run test -- auth store`
**Run Backend:** `cd backend && go test -v -run TestOnAddUser`

#### 4. Message Sending with Gifts
**Files:**
- `frontend/src/__tests__/SendMessageForm.test.ts`
- `backend/message_hooks_test.go`

Tests cover:
- Form validation (recipient, content length)
- Gift selection (max 3 gifts)
- Cost calculation
- Duplicate message detection
- Profanity filtering
- Insufficient funds handling
- Message expansion

**Run Frontend:** `cd frontend && npm run test -- SendMessageForm`
**Run Backend:** `cd backend && go test -v -run TestOnBeforeAddMessage`

#### 5. Virtual Wallet Operations
**Files:**
- `backend/virtual_wallet_test.go`

Tests cover:
- Initial wallet creation (1000 coins)
- Sufficient funds checking
- Transaction creation
- Balance updates
- Transaction from user ID

**Run:** `cd backend && go test -v -run Virtual`

#### 6. Message Filtering & Display
**Files:**
- `frontend/src/__tests__/Wall.test.ts`
- `frontend/src/__tests__/MessageTiles.test.ts`

Tests cover:
- Wall page rendering
- Recipient filtering
- Tab filtering (All/Messages/Gifts)
- Message tiles display
- Pagination

**Run:** `cd frontend && npm run test -- Wall MessageTiles`

#### 7. Rankings System
**Files:**
- `backend/message_hooks_test.go`

Tests cover:
- Ranking creation on message send
- Coin accumulation
- Department and sex assignment

**Run:** `cd backend && go test -v -run TestUpdateRanking`

## Setup Instructions

### Backend (Go) Tests

#### Prerequisites
```bash
# Go 1.19 or higher
go version
```

#### Install Dependencies
```bash
cd backend
go mod download
```

#### Run All Tests
```bash
cd backend
go test -v ./...
```

#### Run Specific Test File
```bash
cd backend
go test -v -run TestName
```

#### Run with Coverage
```bash
cd backend
go test -v -cover ./...
```

### Frontend (Vitest) Tests

#### Prerequisites
```bash
# Node.js 16+ and npm
node --version
npm --version
```

#### Install Dependencies
```bash
cd frontend
npm install
```

#### Install Test Dependencies
```bash
cd frontend
npm install -D vitest @vitest/ui jsdom @vue/test-utils happy-dom
```

#### Run All Tests
```bash
cd frontend
npm run test
```

#### Run Tests in Watch Mode
```bash
cd frontend
npm run test:watch
```

#### Run Tests with UI
```bash
cd frontend
npm run test:ui
```

#### Run Specific Test File
```bash
cd frontend
npm run test -- GiftIcon
```

#### Generate Coverage Report
```bash
cd frontend
npm run test:coverage
```

## Package.json Scripts

Add these scripts to `frontend/package.json`:

```json
{
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest run --coverage"
  }
}
```

## Test File Summaries

### Backend Tests

| File | Tests | Coverage |
|------|-------|----------|
| `message_hooks_test.go` | 7 tests | Message expansion, duplicate detection, profanity check, gift cost calculation, ranking updates, insufficient funds |
| `virtual_wallet_test.go` | 7 tests | Wallet creation, transactions, balance updates, sufficient funds checking |
| `user_hooks_test.go` | 3 tests | User creation, user details, user deletion |
| `server_test.go` | 4 tests | API endpoints (departments, gifts, terms, image generation) |

**Total Backend Tests:** 21

### Frontend Tests

| File | Tests | Coverage |
|------|-------|----------|
| `GiftIcon.test.ts` | 8 tests | Emoji rendering, mapping, legacy support, styling |
| `SendMessageForm.test.ts` | 8 tests | Form validation, gift selection, submission |
| `MessageTiles.test.ts` | 4 tests | Message display, limit handling |
| `Wall.test.ts` | 3 tests | Page rendering, recipient display |
| `store.test.ts` | 10 tests | Auth state, main store, gift/dept loading |
| `auth.test.ts` | 5 tests | OAuth flow, popup handling |
| `notify.test.ts` | 7 tests | Notifications, real-time message alerts |
| `utils.test.ts` | 3 tests | Read-only mode checking |

**Total Frontend Tests:** 48

## CI/CD Integration

### GitHub Actions Example

Create `.github/workflows/test.yml`:

```yaml
name: Run Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.19'
      - name: Run backend tests
        run: |
          cd backend
          go test -v ./...

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      - name: Run frontend tests
        run: |
          cd frontend
          npm run test
```

## Known Limitations

1. **Backend Image Rendering Tests** - Skipped due to Chrome headless setup complexity
2. **Full OAuth Flow** - Mocked due to external service dependencies
3. **Email Sending** - Not tested (requires SMTP configuration)
4. **WebSocket Subscriptions** - Partially mocked in tests

## Best Practices

1. **Run tests before committing** - Ensure all tests pass
2. **Write tests for new features** - Maintain coverage
3. **Update tests when changing code** - Keep tests in sync
4. **Use descriptive test names** - Make failures easy to understand
5. **Mock external dependencies** - Keep tests isolated and fast

## Troubleshooting

### Backend Tests Failing

```bash
# Clear test cache
go clean -testcache

# Update dependencies
go mod tidy && go mod download
```

### Frontend Tests Failing

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vitest cache
rm -rf node_modules/.vitest
```

### Import Errors

Ensure `vitest.config.ts` has correct path aliases matching your project structure.

## Future Test Additions

- [ ] E2E tests with Playwright/Cypress
- [ ] Visual regression tests
- [ ] Performance/load tests
- [ ] Accessibility tests
- [ ] Database migration tests
- [ ] API contract tests

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [Vue Test Utils](https://test-utils.vuejs.org/)
- [Go Testing](https://pkg.go.dev/testing)
- [PocketBase Testing](https://pocketbase.io/docs/testing/)

---

## Quick Reference

**Run all tests:**
```bash
# Backend
cd backend && go test -v ./...

# Frontend
cd frontend && npm run test
```

**Run recent feature tests:**
```bash
# Gift icons
cd frontend && npm run test -- GiftIcon

# Real-time notifications
cd frontend && npm run test -- notify
```

**Generate coverage:**
```bash
# Backend
cd backend && go test -v -cover ./...

# Frontend
cd frontend && npm run test:coverage
```
