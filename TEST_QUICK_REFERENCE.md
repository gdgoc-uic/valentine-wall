# Valentine Wall - Test Quick Reference

## 🚀 Quick Commands

### Run All Tests
```bash
# Backend
cd backend && go test -v ./...

# Frontend
cd frontend && npm test
```

### Run Specific Tests
```bash
# Backend - Message hooks
cd backend && go test -v -run Message

# Backend - Virtual wallet
cd backend && go test -v -run Virtual

# Frontend - Gift icons
cd frontend && npm test -- GiftIcon

# Frontend - Notifications
cd frontend && npm test -- notify

# Frontend - Auth
cd frontend && npm test -- auth
```

### Watch Mode
```bash
cd frontend && npm run test:watch
```

### Coverage Reports
```bash
# Backend
cd backend && go test -v -cover ./...

# Frontend
cd frontend && npm run test:coverage
```

---

## 📋 Test Files by Feature

### Recent Features

| Feature | Test File | Command |
|---------|-----------|---------|
| Gift Icons | `frontend/src/__tests__/GiftIcon.test.ts` | `cd frontend && npm test -- GiftIcon` |
| Real-time Notifications | `frontend/src/__tests__/notify.test.ts` | `cd frontend && npm test -- notify` |

### Core Features

| Feature | Test File | Command |
|---------|-----------|---------|
| User Auth | `frontend/src/__tests__/auth.test.ts` | `cd frontend && npm test -- auth` |
| Auth Store | `frontend/src/__tests__/store.test.ts` | `cd frontend && npm test -- store` |
| Message Form | `frontend/src/__tests__/SendMessageForm.test.ts` | `cd frontend && npm test -- SendMessageForm` |
| Wall Page | `frontend/src/__tests__/Wall.test.ts` | `cd frontend && npm test -- Wall` |
| Message Hooks | `backend/message_hooks_test.go` | `cd backend && go test -v -run Message` |
| Virtual Wallet | `backend/virtual_wallet_test.go` | `cd backend && go test -v -run Virtual` |
| User Hooks | `backend/user_hooks_test.go` | `cd backend && go test -v -run User` |

---

## 📊 Test Statistics

### Backend
- **Files:** 4
- **Tests:** 21
- **Coverage:** Message hooks, wallet, user lifecycle, API endpoints

### Frontend
- **Files:** 8
- **Tests:** 48
- **Coverage:** Components, stores, auth, notifications, utils

**Total:** 69 tests across 12 files

---

## 🔧 First-Time Setup

### Backend
```bash
cd backend
go mod download
go test -v ./...
```

### Frontend
```bash
cd frontend
npm install
npm test
```

---

## 🐛 Troubleshooting

### Tests Not Running?

**Backend:**
```bash
go clean -testcache
go mod tidy
```

**Frontend:**
```bash
rm -rf node_modules package-lock.json
npm install
```

### Import Errors?
Check `vitest.config.ts` path aliases

### Timeout Errors?
Increase timeout in individual test files

---

## 📚 More Information

- **Full Guide:** [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **Feature List:** [FEATURE_LIST.md](FEATURE_LIST.md)
- **Summary:** [TEST_IMPLEMENTATION_SUMMARY.md](TEST_IMPLEMENTATION_SUMMARY.md)
