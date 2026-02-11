# Valentine Wall - Quick Reference Card

Quick reference for common development tasks. Keep this handy! 📌

---

## 🚀 Starting & Stopping

```bash
# Start development environment
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Start in background
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Stop all services
docker-compose down

# Restart a specific service
docker-compose restart backend
docker-compose restart frontend

# Rebuild after dependency changes
docker-compose up --build
```

---

## 🌐 Service URLs

| Service | URL | Purpose |
|---------|-----|---------|
| **Frontend** | http://localhost:3000 | Main application |
| **Backend API** | http://localhost:8090 | API endpoints |
| **Admin UI** | http://localhost:8090/_/ | PocketBase dashboard |
| **Chrome DevTools** | http://localhost:5000 | Headless browser |

---

## 📁 Key Files & Locations

### Backend (Go)
```
backend/
├── main.go                    → App entry point
├── server.go                  → Custom API routes
├── message_hooks.go           → Message business logic
├── user_hooks.go              → User business logic
├── virtual_wallet_hooks.go    → Wallet/economy logic
├── config.go                  → Environment config
├── image_gen.go               → Image generation
├── models/                    → Data models
├── migrations/                → Database migrations
└── templates/mail/            → Email templates
```

### Frontend (Vue 3 + TypeScript)
```
frontend/src/
├── main.ts                    → App entry point
├── App.vue                    → Root component
├── router.ts                  → Routes
├── client.ts                  → PocketBase client
├── auth.ts                    → Auth utilities
├── store_new.ts               → State management
├── components/                → Reusable components
├── pages/                     → Page components
└── assets/                    → Styles & images
```

---

## 🔧 Backend Quick Tasks

### Add API Endpoint

**File:** `backend/server.go`

```go
app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
    e.Router.GET("/my-endpoint", func(c echo.Context) error {
        return c.JSON(200, map[string]interface{}{
            "status": "success",
            "data":   "Hello!",
        })
    })
    return nil
})
```

### Add Before/After Hook

**File:** `backend/message_hooks.go` (or similar)

```go
// Before create - validation
app.OnRecordBeforeCreateRequest("messages").Add(func(e *core.RecordCreateEvent) error {
    record := e.Record
    // Add validation logic
    return nil
})

// After create - side effects
app.OnRecordAfterCreateRequest("messages").Add(func(e *core.RecordCreateEvent) error {
    record := e.Record
    // Send email, update stats, etc.
    return nil
})
```

### Access PocketBase Collections in Code

```go
// Get all records
records, err := app.Dao().FindRecordsByExpr("messages",
    dbx.HashExp{"user": userId},
)

// Get single record
record, err := app.Dao().FindRecordById("messages", messageId)

// Create record
collection, _ := app.Dao().FindCollectionByNameOrId("messages")
record := models.NewRecord(collection)
record.Set("content", "Hello!")
app.Dao().SaveRecord(record)

// Update record
record.Set("content", "Updated!")
app.Dao().SaveRecord(record)

// Delete record
app.Dao().DeleteRecord(record)
```

---

## 🎨 Frontend Quick Tasks

### Create New Page

**File:** `frontend/src/pages/MyPage.vue`

```vue
<script setup lang="ts">
import { ref } from 'vue';

const title = ref('My New Page');
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold">{{ title }}</h1>
    <p>Content goes here</p>
  </div>
</template>
```

**Add Route:** `frontend/src/router.ts`

```typescript
{
  path: '/my-page',
  name: 'MyPage',
  component: () => import('./pages/MyPage.vue'),
  meta: { title: 'My Page' },
}
```

### Fetch Data from Backend

```typescript
import { useQuery } from '@tanstack/vue-query';
import { pb } from '@/client';

// Fetch list
const { data, isLoading, error } = useQuery({
  queryKey: ['messages'],
  queryFn: () => pb.collection('messages').getList(1, 20),
  enabled: !import.meta.env.SSR,
});

// Fetch single
const message = useQuery({
  queryKey: ['message', messageId],
  queryFn: () => pb.collection('messages').getOne(messageId),
});
```

### Create/Update Data

```typescript
import { useMutation, useQueryClient } from '@tanstack/vue-query';
import { pb } from '@/client';

const queryClient = useQueryClient();

const createMessage = useMutation({
  mutationFn: (data) => pb.collection('messages').create(data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['messages'] });
  },
});

// Use it
function handleSubmit() {
  createMessage.mutate({
    recipient: '2021-12345',
    content: 'Hello!',
  });
}
```

### Real-time Subscription

```typescript
import { onMounted, onUnmounted } from 'vue';
import { pb } from '@/client';

onMounted(() => {
  pb.collection('messages').subscribe(messageId, (data) => {
    console.log('Updated:', data.record);
  });
});

onUnmounted(() => {
  pb.collection('messages').unsubscribe(messageId);
});
```

---

## 🗄️ Database (PocketBase Admin)

### Access Admin UI
```
http://localhost:8090/_/
```

### Common Collections

```
users               → User authentication
user_details        → Extended user info (student_id, department)
messages            → Valentine messages
message_replies     → Replies to messages
gifts               → Available gifts catalog
college_departments → Departments list
virtual_wallets     → User wallets
virtual_transactions → Transaction history
rankings            → Recipient rankings
```

### API Rules (Collection Settings)

- **List Rule:** Who can view list of records
- **View Rule:** Who can view single record
- **Create Rule:** Who can create records
- **Update Rule:** Who can update records
- **Delete Rule:** Who can delete records

**Common patterns:**
```
@request.auth.id != ""           → Authenticated users
@request.auth.id = user.id       → Owner only
""                               → Public (empty string)
null                             → No one (null)
```

---

## 🐳 Docker Commands

```bash
# View running containers
docker-compose ps

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs frontend --tail=50

# Execute commands in container
docker-compose exec backend sh
docker-compose exec frontend npm install <package>

# Rebuild specific service
docker-compose up -d --no-deps --build backend

# Remove all containers and volumes (fresh start)
docker-compose down -v

# Clean up Docker system
docker system prune -a
```

---

## 🔍 Debugging

### Backend Logs
```bash
docker-compose logs -f backend
```

**Look for:**
- API request logs
- Error messages
- Database queries
- Email sending status

### Frontend Console
```
F12 → Console Tab (in browser)
```

**Look for:**
- JavaScript errors
- Network requests (Network tab)
- Vue warnings
- API response errors

### Database Inspection
```
http://localhost:8090/_/
```

- View all records
- Check field values
- Verify relationships
- Test API rules

### Check API Directly

```bash
# Get all messages
curl http://localhost:8090/api/collections/messages/records

# Get specific message
curl http://localhost:8090/api/collections/messages/records/RECORD_ID

# With authentication
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8090/api/collections/messages/records
```

---

## 📧 Email Templates

**Location:** `backend/templates/mail/*.txt.tpl`

**Available Templates:**
- `welcome.txt.tpl` - Welcome email
- `message.txt.tpl` - New message notification
- `reply.txt.tpl` - New reply notification

**Template Variables:**
```
{{.user.name}}          → User name
{{.user.email}}         → User email
{{.message.content}}    → Message content
{{.message.id}}         → Message ID
{{.sender.name}}        → Sender name
```

**Send Email in Code:**
```go
templateData := map[string]interface{}{
    "user":    userRecord,
    "message": messageRecord,
}

err := sendEmail(
    app,
    recipient.Email(),
    "Subject Line",
    "mail/template_name.txt.tpl",
    templateData,
)
```

---

## 🎨 Styling (TailwindCSS)

### Common Classes

```html
<!-- Layout -->
<div class="container mx-auto px-4 py-8">
<div class="flex items-center justify-between">
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">

<!-- Spacing -->
<div class="p-4 m-4">          <!-- Padding & Margin -->
<div class="px-4 py-2">        <!-- Horizontal & Vertical -->
<div class="space-y-4">        <!-- Vertical spacing between children -->

<!-- Typography -->
<h1 class="text-3xl font-bold text-gray-900">
<p class="text-sm text-gray-600">

<!-- Colors -->
<div class="bg-blue-500 text-white">
<div class="bg-red-100 text-red-800">

<!-- Buttons (DaisyUI) -->
<button class="btn btn-primary">Primary</button>
<button class="btn btn-secondary">Secondary</button>
<button class="btn btn-outline">Outline</button>

<!-- Cards (DaisyUI) -->
<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Title</h2>
    <p>Content</p>
  </div>
</div>

<!-- Responsive -->
<div class="w-full md:w-1/2 lg:w-1/3">
<div class="hidden md:block">           <!-- Hide on mobile -->
<div class="block md:hidden">           <!-- Show only on mobile -->
```

---

## 🔐 Authentication

### Check Auth Status (Frontend)

```typescript
import { pb } from '@/client';

// Current user
const user = pb.authStore.model;

// Is authenticated
const isAuth = pb.authStore.isValid;

// Login
await pb.collection('users').authWithPassword(email, password);

// Logout
pb.authStore.clear();

// Listen to auth changes
pb.authStore.onChange((token, model) => {
  console.log('Auth changed:', model);
});
```

### Require Auth in Routes

```typescript
// frontend/src/router.ts
{
  path: '/profile',
  component: () => import('./pages/Profile.vue'),
  meta: { requiresAuth: true },
}
```

### Require Auth in Backend

```go
// In server.go
e.Router.GET("/protected", func(c echo.Context) error {
    // Your logic
    return c.JSON(200, data)
}, apis.RequireRecordAuth())
```

---

## 🧪 Testing

### Test API Endpoint

```bash
# GET request
curl http://localhost:8090/api/collections/messages/records

# POST request
curl -X POST http://localhost:8090/api/collections/messages/records \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"recipient":"2021-12345","content":"Test"}'

# With query params
curl "http://localhost:8090/api/collections/messages/records?filter=user='USER_ID'"
```

### Frontend Type Checking

```bash
cd frontend
npm run type-check
```

### Build Test

```bash
# Frontend
cd frontend && npm run build

# Backend
cd backend && go build
```

---

## 📦 Dependencies

### Add Backend Dependency

```bash
cd backend
go get github.com/some/package
go mod tidy

# Rebuild
docker-compose up --build backend
```

### Add Frontend Dependency

```bash
cd frontend
npm install package-name

# Or in Docker
docker-compose exec frontend npm install package-name

# Rebuild
docker-compose up --build frontend
```

---

## 🚀 Production Deployment

```bash
# 1. Update .env for production
ENV=production
BASE_DOMAIN=yourdomain.com
BACKEND_URL=https://yourdomain.com/pb
FRONTEND_URL=https://yourdomain.com

# 2. Create Caddy network (one-time)
docker network create caddy

# 3. Deploy
./deploy.sh

# 4. Check logs
docker-compose logs -f www
docker-compose logs -f backend
docker-compose logs -f frontend
```

---

## 💾 Backup & Restore

### Backup Database

```bash
# SQLite database
docker cp valentine-wall-backend-1:/pb_data/data.db ./backup_$(date +%Y%m%d).db

# Entire pb_data directory
docker cp valentine-wall-backend-1:/pb_data ./backup_pb_data
```

### Restore Database

```bash
# Stop backend
docker-compose stop backend

# Copy database
docker cp ./backup.db valentine-wall-backend-1:/pb_data/data.db

# Start backend
docker-compose start backend
```

---

## 🆘 Emergency Commands

```bash
# Container won't start
docker-compose down
docker-compose up --build

# Complete reset (⚠️ deletes data)
docker-compose down -v
rm -rf pb_data/*
docker-compose up --build

# Port already in use (Windows)
netstat -ano | findstr :8090
taskkill /PID <PID> /F

# Port already in use (Linux/Mac)
lsof -i :8090
kill -9 <PID>

# Out of disk space
docker system prune -a --volumes
```

---

## 📚 More Help

- **Full Guide:** See `DEVELOPMENT_GUIDE.md`
- **Setup Checklist:** See `SETUP_CHECKLIST.md`
- **PocketBase Docs:** https://pocketbase.io/docs/
- **Vue 3 Docs:** https://vuejs.org/
- **TailwindCSS Docs:** https://tailwindcss.com/

---

**💡 Pro Tip:** Keep this file open in a separate tab while developing!
