# Valentine Wall - Development Guide

> **Last Updated:** February 11, 2026  
> **Target Audience:** Developers setting up or modifying the Valentine Wall project

---

## 📋 Table of Contents

1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Architecture & How Components Interact](#architecture--how-components-interact)
4. [Initial Setup](#initial-setup)
5. [Development Workflow](#development-workflow)
6. [Making Changes - Backend](#making-changes---backend)
7. [Making Changes - Frontend](#making-changes---frontend)
8. [Common Development Tasks](#common-development-tasks)
9. [Production Deployment](#production-deployment)
10. [Troubleshooting](#troubleshooting)

---

## 🎯 Project Overview

**Valentine Wall** is a full-stack web application that allows users to send anonymous Valentine's Day messages with virtual gifts. It features:

- User authentication with email verification
- Virtual currency wallet system (users start with 1000 coins)
- Message sending with virtual gifts (150 coins per message + gift prices)
- Real-time ranking system
- Social media card generation for messages
- Message archiving and downloading
- Profanity filtering in multiple languages

---

## 🔧 Technology Stack

### **Backend**
| Component | Technology | Version |
|-----------|-----------|---------|
| **Framework** | PocketBase (Go) | 0.11.3 |
| **Language** | Go | 1.19+ |
| **Database** | SQLite (embedded) | - |
| **HTTP Router** | Echo v5 | - |
| **Authentication** | PocketBase Auth + Firebase | 4.10.0 |
| **Image Rendering** | Chrome DevTools Protocol | chromedp 0.7.6 |
| **Email** | PocketBase SMTP | - |

**Key Backend Libraries:**
- `fogleman/gg` - 2D graphics rendering (fallback)
- `go-away` - Profanity detection
- `patrickmn/go-cache` - In-memory caching
- `go-playground/validator` - Data validation

### **Frontend**
| Component | Technology | Version |
|-----------|-----------|---------|
| **Framework** | Vue 3 (Composition API) | 3.2.45 |
| **Language** | TypeScript | 4.9+ |
| **Build Tool** | Vite | 4.0.0 |
| **Routing** | Vue Router | 4.1.6 |
| **State Management** | Vuex 4 + Custom Stores | 4.0.2 |
| **Data Fetching** | TanStack Query (Vue Query) | 4.24.4 |
| **Backend SDK** | PocketBase JS SDK | 0.10.1 |
| **UI Framework** | TailwindCSS + DaisyUI | 3.2.4 / 2.49.0 |
| **SSR** | Yes (Node.js + Polka) | - |

**Key Frontend Libraries:**
- `firebase` - Authentication (Google/Facebook)
- `@vueuse/head` - Document head management
- `floating-vue` - Tooltips and popovers
- `dayjs` - Date manipulation
- `masonry-layout` - Grid layouts

### **Infrastructure**
| Component | Technology |
|-----------|-----------|
| **Containerization** | Docker + Docker Compose |
| **Reverse Proxy (Prod)** | Caddy 2 (auto-SSL) |
| **Image Renderer** | Browserless Chrome |
| **Process Manager** | Docker |

---

## 🏗️ Architecture & How Components Interact

### **High-Level Architecture**

```
┌─────────────────────────────────────────────────────────────┐
│                        User Browser                         │
└────────────────┬────────────────────────────────────────────┘
                 │
                 │ HTTPS (Production) / HTTP (Development)
                 ↓
┌────────────────────────────────────────────────────────────┐
│                    Caddy Reverse Proxy                      │
│                    (Production Only)                        │
│  ┌──────────────────────┐    ┌──────────────────────────┐  │
│  │  / → Frontend:3000   │    │ /pb/* → Backend:8090     │  │
│  └──────────────────────┘    └──────────────────────────┘  │
└────────┬──────────────────────────────┬────────────────────┘
         │                               │
         │                               │
         ↓                               ↓
┌────────────────────┐         ┌──────────────────────────┐
│   Frontend (SSR)   │←────────│  Backend (PocketBase)    │
│   Vue 3 + Vite     │  REST   │  Go + SQLite             │
│   Node.js Server   │  + WS   │  Port: 8090              │
│   Port: 3000       │         │                          │
└────────────────────┘         └───────────┬──────────────┘
                                           │
                                           │ WebSocket
                                           ↓
                               ┌────────────────────────┐
                               │  Headless Chrome       │
                               │  Image Rendering       │
                               │  Port: 5000 (internal) │
                               └────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  Persistent Storage                         │
│  • pb_data/ (SQLite database + uploads)                     │
│  • pb_public/ (public static files)                         │
│  • backend/_data/ (terms & conditions)                      │
└─────────────────────────────────────────────────────────────┘
```

### **Frontend ↔ Backend Communication**

#### **1. REST API (Primary)**

The frontend uses the PocketBase JavaScript SDK to communicate with the backend:

```typescript
// frontend/src/client.ts
import PocketBase from 'pocketbase';

export const pb = new PocketBase(import.meta.env.VITE_BACKEND_URL);
```

**Standard PocketBase Endpoints (Auto-generated):**
- `POST /api/collections/users/auth-with-password` - Login
- `POST /api/collections/users/records` - Register
- `GET /api/collections/messages/records` - Fetch messages
- `POST /api/collections/messages/records` - Create message
- `GET /api/collections/gifts/records` - Get gifts
- `PATCH /api/collections/{collection}/records/{id}` - Update record
- `DELETE /api/collections/{collection}/records/{id}` - Delete record

**Custom Backend Endpoints:**
- `GET /departments` - List college departments
- `GET /gifts` - List available gifts
- `GET /messages/:id/image` - Generate social media card image
- `GET /user_messages/archive` - Stream ZIP archive creation (SSE)
- `GET /user_messages/download_archive/:userId` - Download archive
- `GET /terms-and-conditions` - Terms and conditions text

#### **2. Real-time Subscriptions (WebSocket)**

PocketBase provides real-time updates via WebSocket:

```typescript
// Example: Listen to wallet balance changes
pb.collection('virtual_wallets').subscribe(walletId, (data) => {
  // Update UI when balance changes
  wallet.value = data.record;
});
```

#### **3. Authentication Flow**

```
User Action → Firebase Auth → Get Token → PocketBase Verify → Session
```

The frontend authenticates with Firebase (for Google/Facebook login), then uses the Firebase token to authenticate with PocketBase:

```typescript
// Simplified flow
const firebaseUser = await signInWithGoogle();
const token = await firebaseUser.getIdToken();
await pb.collection('users').authWithProvider('firebase', token);
```

### **Backend Business Logic Flow**

#### **Message Creation Example:**

```
1. User submits message form (frontend)
   ↓
2. Frontend calls: pb.collection('messages').create(data)
   ↓
3. Backend receives request → PocketBase validation
   ↓
4. BEFORE CREATE HOOK in message_hooks.go:
   - Check for duplicate messages
   - Run profanity filter on content
   - Verify user has sufficient funds
   ↓
5. Create message record in SQLite
   ↓
6. AFTER CREATE HOOK in message_hooks.go:
   - Deduct coins from sender's wallet (150 + gift prices)
   - Add coins to recipient's wallet (gift values marked remittable)
   - Update recipient's ranking
   - Send email notification to recipient
   - Create virtual transaction records
   ↓
7. Return created message to frontend
   ↓
8. Frontend updates UI and refetches wallet balance
```

### **Data Models & Relationships**

```
users (PocketBase Auth)
  ├─→ user_details (one-to-one)
  │     • student_id (unique)
  │     • email
  │     • college_department
  │     • sex
  │
  ├─→ virtual_wallets (one-to-one)
  │     • balance (default: 1000)
  │     └─→ virtual_transactions (one-to-many)
  │
  ├─→ messages (one-to-many) as sender
  │     • recipient (student_id reference)
  │     • content (max 240 chars)
  │     • gifts[] (array of gift UIDs)
  │     └─→ message_replies (one-to-many)
  │
  └─→ rankings (computed)
        • Groups by: recipient, department, sex
        • total_coins (sum of gift values received)
```

### **Image Generation Flow**

When a user shares a message on social media:

```
1. Frontend generates URL: /messages/{id}/image
   ↓
2. Backend checks cache (10-minute TTL)
   ↓
3. If not cached:
   a. Render HTML template with message data
   b. Send HTML to headless Chrome via WebSocket
   c. Chrome renders page and takes screenshot
   d. Resize to 1200x675 (Twitter card size)
   e. Cache result in memory
   ↓
4. Return PNG image with proper headers
```

---

## 🚀 Initial Setup

### **Prerequisites**

Install the following on your development machine:

- **Docker Desktop** (includes Docker Compose)
- **Git**
- **(Optional) Go 1.19+** - for local backend development without Docker
- **(Optional) Node.js 16.13.2+** - for local frontend development without Docker

### **1. Clone the Repository**

```bash
git clone https://github.com/yourusername/valentine-wall.git
cd valentine-wall
```

### **2. Create Environment File**

Create `.env` file in the project root:

```bash
# Environment
ENV=development

# URLs
BASE_DOMAIN=localhost
BACKEND_URL=http://localhost:8090
FRONTEND_URL=http://localhost:3000

# Feature Flags
READ_ONLY=false

# Profanity Filter
PROFANITY_JSON_FILE_NAME=profanities.json

# Firebase (Optional - for OAuth)
FIREBASE_API_KEY=your_api_key
FIREBASE_AUTH_DOMAIN=your_project.firebaseapp.com
FIREBASE_PROJECT_ID=your_project_id
FIREBASE_STORAGE_BUCKET=your_project.appspot.com
FIREBASE_MESSAGING_SENDER_ID=123456789
FIREBASE_APP_ID=1:123456789:web:abcdef
FIREBASE_MEASUREMENT_ID=G-ABCDEFG

# Reporting API (Optional)
REPORT_API_URL=
REPORT_API_KEY=
REPORT_API_CATEGORY_ID_KEY=
```

### **3. Prepare Data Directory**

Copy template data files:

```bash
cp -r backend/_data.please_copy backend/_data
```

### **4. Start Development Environment**

```bash
# Start all services with Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

**Services will be available at:**
- 🌐 **Frontend:** http://localhost:3000
- 🔌 **Backend API:** http://localhost:8090
- 🎨 **Backend Admin UI:** http://localhost:8090/_/
- 🖥️ **Headless Chrome:** http://localhost:5000 (internal only)

### **5. Access PocketBase Admin**

1. Open http://localhost:8090/_/
2. Create admin account on first launch
3. Import initial data (departments, gifts) if needed

---

## 💻 Development Workflow

### **Recommended Development Setup**

**Option A: Docker Compose (Simplest)**
- All services run in containers
- No local dependencies needed
- Code changes auto-reload via volume mounts

**Option B: Hybrid (Most Flexible)**
- Backend in Docker
- Frontend runs locally with `npm run dev`
- Faster frontend hot reload

**Option C: Fully Local (Advanced)**
- All services run on host machine
- Requires Go, Node.js, and Chrome installed
- Best for debugging

### **Using Docker Compose (Recommended for Beginners)**

```bash
# Start all services
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Start specific service
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up frontend

# Rebuild after dependency changes
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# View logs
docker-compose logs -f frontend
docker-compose logs -f backend

# Stop all services
docker-compose down
```

### **Using Local Development Scripts**

First, install `dotenv-cli`:
```bash
npm install -g dotenv-cli
```

**Terminal 1 - Backend:**
```bash
./start-backend.sh serve --http=0.0.0.0:8090 --dir=./pb_data
```

**Terminal 2 - Frontend:**
```bash
./start-frontend.sh
```

This runs:
- Backend with Go (PocketBase)
- Frontend with Vite dev server
- Auto-reload on file changes

---

## 🛠️ Making Changes - Backend

### **File Structure Overview**

```
backend/
├── main.go                    # Entry point, PocketBase initialization
├── server.go                  # Custom routes and middleware
├── config.go                  # Environment configuration
├── message_hooks.go           # Business logic for messages
├── user_hooks.go              # Business logic for users
├── virtual_wallet_hooks.go    # Virtual economy logic
├── image_gen.go               # Message image generation
├── mail.go                    # Email utilities
├── utils.go                   # Helper functions
├── errors.go                  # Custom error types
├── virtual_models.go          # Data models for virtual items
├── virtual_transactions.go    # Transaction helpers
├── models/                    # Database models
│   ├── user.go
│   ├── message.go
│   ├── gift.go
│   └── college_department.go
├── migrations/                # Database migrations
│   └── *.go
└── templates/                 # Email & HTML templates
    ├── mail/*.txt.tpl
    └── html/*.html.tpl
```

### **Common Backend Modifications**

#### **1. Adding a New API Endpoint**

**File:** `backend/server.go`

```go
func addCustomRoutes(app *pocketbase.PocketBase) error {
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        // Add your custom route
        e.Router.GET("/your-new-endpoint", func(c echo.Context) error {
            // Your logic here
            return c.JSON(200, map[string]interface{}{
                "message": "Hello from custom endpoint",
            })
        }, apis.RequireRecordAuth()) // Optional: require authentication
        
        return nil
    })
    return nil
}
```

**Test your endpoint:**
```bash
curl http://localhost:8090/your-new-endpoint
```

#### **2. Adding Business Logic to Existing Collection**

**File:** `backend/message_hooks.go` (or create similar file for your collection)

```go
// Add BEFORE CREATE validation
func setupMessageHooks(app *pocketbase.PocketBase) error {
    app.OnRecordBeforeCreateRequest("messages").Add(func(e *core.RecordCreateEvent) error {
        // Validate or modify data before creation
        record := e.Record
        
        // Example: Add server-side timestamp
        record.Set("server_timestamp", time.Now().Unix())
        
        // Example: Validate custom business rule
        if record.GetString("content") == "" {
            return errors.New("content cannot be empty")
        }
        
        return nil
    })
    
    // Add AFTER CREATE side effects
    app.OnRecordAfterCreateRequest("messages").Add(func(e *core.RecordCreateEvent) error {
        // Perform actions after successful creation
        record := e.Record
        
        // Example: Send notification, update stats, etc.
        log.Println("New message created:", record.Id)
        
        return nil
    })
    
    return nil
}
```

#### **3. Creating a New Database Collection**

**Option A: Via PocketBase Admin UI (Recommended for development)**
1. Go to http://localhost:8090/_/
2. Click "New Collection"
3. Define fields and validation rules
4. Save and test

**Option B: Via Migration (Recommended for production)**

Create: `backend/migrations/TIMESTAMP_create_my_collection.go`

```go
package migrations

import (
    "github.com/pocketbase/dbx"
    "github.com/pocketbase/pocketbase/daos"
    m "github.com/pocketbase/pocketbase/migrations"
    "github.com/pocketbase/pocketbase/models"
    "github.com/pocketbase/pocketbase/models/schema"
)

func init() {
    m.Register(func(db dbx.Builder) error {
        dao := daos.New(db)

        collection := &models.Collection{
            Name:       "my_collection",
            Type:       models.CollectionTypeBase,
            ListRule:   types.Pointer("@request.auth.id != ''"),
            ViewRule:   types.Pointer("@request.auth.id != ''"),
            CreateRule: types.Pointer("@request.auth.id != ''"),
            UpdateRule: types.Pointer("@request.auth.id = user.id"),
            DeleteRule: types.Pointer("@request.auth.id = user.id"),
            Schema: schema.NewSchema(
                &schema.SchemaField{
                    Name:     "user",
                    Type:     schema.FieldTypeRelation,
                    Required: true,
                    Options: &schema.RelationOptions{
                        CollectionId:  "_pb_users_auth_",
                        CascadeDelete: true,
                    },
                },
                &schema.SchemaField{
                    Name:     "title",
                    Type:     schema.FieldTypeText,
                    Required: true,
                },
                &schema.SchemaField{
                    Name: "content",
                    Type: schema.FieldTypeText,
                },
            ),
        }

        return dao.SaveCollection(collection)
    }, func(db dbx.Builder) error {
        dao := daos.New(db)
        collection, _ := dao.FindCollectionByNameOrId("my_collection")
        return dao.DeleteCollection(collection)
    })
}
```

#### **4. Modifying Email Templates**

**File:** `backend/templates/mail/your_email.txt.tpl`

```text
Hello {{.user.name}},

Your message has been sent!

Message ID: {{.message.id}}
Content: {{.message.content}}

Best regards,
Valentine Wall Team
```

**Send email in code:**

```go
import (
    "github.com/pocketbase/pocketbase/tools/mailer"
)

err := app.NewMailClient().Send(
    &mailer.Message{
        From: mail.Address{
            Address: "noreply@yourapp.com",
            Name:    "Valentine Wall",
        },
        To:      []mail.Address{{Address: user.Email()}},
        Subject: "Your Message Was Sent",
        Text:    renderTemplate("mail/your_email.txt.tpl", data),
    },
)
```

#### **5. Testing Backend Changes**

```bash
# Run backend tests (if available)
cd backend
go test ./...

# Manual API testing with curl
curl -X POST http://localhost:8090/api/collections/messages/records \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"recipient": "2021-12345", "content": "Test message"}'

# Check logs
docker-compose logs -f backend
```

---

## 🎨 Making Changes - Frontend

### **File Structure Overview**

```
frontend/
├── src/
│   ├── main.ts                # Client entry point
│   ├── main-client.ts         # Client-side hydration
│   ├── main-server.js         # SSR entry point
│   ├── App.vue                # Root component
│   ├── router.ts              # Route definitions
│   ├── store_new.ts           # State management
│   ├── client.ts              # PocketBase client
│   ├── auth.ts                # Authentication utilities
│   ├── firebase.ts            # Firebase config
│   ├── types.ts               # TypeScript type definitions
│   ├── components/            # Reusable Vue components
│   │   ├── LoginButton.vue
│   │   ├── MessageCard.vue
│   │   ├── GiftIcon.vue
│   │   └── ...
│   ├── pages/                 # Page components (used by router)
│   │   ├── Home.vue
│   │   ├── Messages.vue
│   │   ├── Profile.vue
│   │   └── ...
│   └── assets/                # Static assets, styles, images
│       ├── index.css          # Global styles
│       └── images/
├── public/                    # Public static files
├── index.html                 # HTML template
├── vite.config.ts             # Vite configuration
├── tailwind.config.js         # TailwindCSS config
└── package.json               # Dependencies
```

### **Common Frontend Modifications**

#### **1. Adding a New Page**

**Step 1:** Create page component `frontend/src/pages/MyNewPage.vue`

```vue
<script setup lang="ts">
import { ref } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import { pb } from '@/client';

// Fetch data
const { data: items, isLoading } = useQuery({
  queryKey: ['my-items'],
  queryFn: async () => {
    return await pb.collection('my_collection').getFullList();
  },
  enabled: !import.meta.env.SSR, // Disable during SSR
});

const message = ref('Hello from my new page!');
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-4">My New Page</h1>
    <p>{{ message }}</p>
    
    <div v-if="isLoading">Loading...</div>
    <ul v-else>
      <li v-for="item in items" :key="item.id">
        {{ item.title }}
      </li>
    </ul>
  </div>
</template>

<style scoped>
/* Component-specific styles */
</style>
```

**Step 2:** Add route in `frontend/src/router.ts`

```typescript
const routes: RouteRecordRaw[] = [
  // ... existing routes ...
  {
    path: '/my-page',
    name: 'MyPage',
    component: () => import('./pages/MyNewPage.vue'),
    meta: {
      title: 'My New Page',
      requiresAuth: true, // Optional: require login
    },
  },
];
```

**Step 3:** Add navigation link in layout/component

```vue
<router-link to="/my-page" class="btn">
  Go to My Page
</router-link>
```

#### **2. Creating a Reusable Component**

**File:** `frontend/src/components/MyComponent.vue`

```vue
<script setup lang="ts">
import { computed } from 'vue';

// Props
interface Props {
  title: string;
  count?: number;
}

const props = withDefaults(defineProps<Props>(), {
  count: 0,
});

// Emits
const emit = defineEmits<{
  (e: 'click', value: number): void;
}>();

// Computed
const displayText = computed(() => {
  return `${props.title}: ${props.count}`;
});

// Methods
function handleClick() {
  emit('click', props.count + 1);
}
</script>

<template>
  <div class="my-component">
    <p>{{ displayText }}</p>
    <button @click="handleClick" class="btn btn-primary">
      Increment
    </button>
  </div>
</template>

<style scoped>
.my-component {
  @apply p-4 bg-gray-100 rounded;
}
</style>
```

**Usage in parent component:**

```vue
<script setup>
import MyComponent from '@/components/MyComponent.vue';

function handleIncrement(newValue) {
  console.log('New value:', newValue);
}
</script>

<template>
  <MyComponent 
    title="Counter" 
    :count="5" 
    @click="handleIncrement" 
  />
</template>
```

#### **3. Fetching Data from Backend**

**Using Vue Query (Recommended):**

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { pb } from '@/client';

// Fetch list
const { data, isLoading, error } = useQuery({
  queryKey: ['messages', userId],
  queryFn: async () => {
    return await pb.collection('messages').getList(1, 50, {
      filter: `user="${userId}"`,
      sort: '-created',
    });
  },
  enabled: !import.meta.env.SSR && !!userId,
});

// Create mutation
const queryClient = useQueryClient();
const createMessage = useMutation({
  mutationFn: async (data) => {
    return await pb.collection('messages').create(data);
  },
  onSuccess: () => {
    // Refetch messages list
    queryClient.invalidateQueries({ queryKey: ['messages'] });
  },
});

// Use mutation
function sendMessage() {
  createMessage.mutate({
    recipient: '2021-12345',
    content: 'Hello!',
    gifts: ['gift-1'],
  });
}
```

**Using direct PocketBase SDK:**

```typescript
import { pb } from '@/client';

// Get list
const messages = await pb.collection('messages').getList(1, 20, {
  filter: 'recipient="2021-12345"',
  expand: 'user',
});

// Get single record
const message = await pb.collection('messages').getOne(messageId, {
  expand: 'user,message_replies',
});

// Create
const newMessage = await pb.collection('messages').create({
  recipient: '2021-12345',
  content: 'Hello!',
  gifts: [],
});

// Update
await pb.collection('messages').update(messageId, {
  content: 'Updated content',
});

// Delete
await pb.collection('messages').delete(messageId);

// Real-time subscription
pb.collection('messages').subscribe(messageId, (e) => {
  console.log('Message updated:', e.record);
});
```

#### **4. Adding Global State**

**File:** `frontend/src/store_new.ts`

```typescript
import { reactive, readonly, inject, InjectionKey } from 'vue';

export interface MyStore {
  counter: number;
  items: any[];
  increment: () => void;
  addItem: (item: any) => void;
}

export const MyStoreKey: InjectionKey<MyStore> = Symbol('MyStore');

export function createMyStore(): MyStore {
  const state = reactive({
    counter: 0,
    items: [] as any[],
  });

  function increment() {
    state.counter++;
  }

  function addItem(item: any) {
    state.items.push(item);
  }

  return {
    ...readonly(state),
    increment,
    addItem,
  };
}

// Usage in main.ts
import { createMyStore, MyStoreKey } from './store_new';
const myStore = createMyStore();
app.provide(MyStoreKey, myStore);

// Usage in component
const store = inject(MyStoreKey)!;
store.increment();
```

#### **5. Styling with TailwindCSS**

**Utility Classes (Preferred):**

```vue
<template>
  <!-- Spacing, colors, typography -->
  <div class="p-4 bg-blue-500 text-white font-bold">
    Hello World
  </div>
  
  <!-- Responsive design -->
  <div class="w-full md:w-1/2 lg:w-1/3">
    Responsive width
  </div>
  
  <!-- DaisyUI components -->
  <button class="btn btn-primary">Primary Button</button>
  <div class="card bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title">Card Title</h2>
      <p>Card content</p>
    </div>
  </div>
</template>
```

**Custom CSS (When needed):**

```vue
<style scoped>
/* Component-scoped styles */
.custom-class {
  @apply px-4 py-2 bg-gradient-to-r from-purple-500 to-pink-500;
}

/* Non-Tailwind custom styles */
.special-effect {
  background: linear-gradient(45deg, #ff00ff, #00ffff);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
```

#### **6. Testing Frontend Changes**

```bash
# Start dev server
cd frontend
npm run dev

# Type checking
npm run type-check

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

---

## 📝 Common Development Tasks

### **1. Add a New Gift Type**

**Backend:** Add to database via PocketBase Admin UI

1. Go to http://localhost:8090/_/
2. Navigate to "gifts" collection
3. Click "New record"
4. Fill in:
   - `uid`: unique identifier (e.g., "rose-gold")
   - `label`: display name (e.g., "Golden Rose")
   - `price`: coin cost (e.g., 200)
   - `is_remittable`: true/false
5. Save

**Frontend:** Gift icon will auto-appear if using dynamic fetch

```typescript
// Already implemented in most components
const { data: gifts } = useQuery({
  queryKey: ['gifts'],
  queryFn: () => pb.send('/gifts', {}),
});
```

### **2. Change Message Cost or Initial Balance**

**File:** `backend/message_hooks.go`

```go
// Change sending cost
const sendPrice = 150 // Change this value

// Change initial balance
// File: backend/user_hooks.go
const initialBalance = 1000 // Change this value
```

After making changes:
```bash
docker-compose restart backend
```

### **3. Add Email Notification for New Event**

**Step 1:** Create template in `backend/templates/mail/my_event.txt.tpl`

```text
Hello {{.recipient.name}},

A new event has occurred!

Event: {{.event.type}}
Time: {{.event.time}}

Best regards,
Valentine Wall
```

**Step 2:** Send email in hook

```go
// In appropriate hook file
app.OnRecordAfterCreateRequest("my_collection").Add(func(e *core.RecordCreateEvent) error {
    recipient, _ := app.Dao().FindRecordById("users", recipientId)
    
    templateData := map[string]interface{}{
        "recipient": recipient,
        "event": map[string]interface{}{
            "type": "New Event",
            "time": time.Now().Format("January 2, 2006 3:04 PM"),
        },
    }
    
    err := sendEmail(app, recipient.Email(), "Event Notification", "mail/my_event.txt.tpl", templateData)
    return err
})
```

### **4. Update Profanity Filter**

**File:** `profanities.json` (root directory)

```json
{
  "words": [
    "badword1",
    "badword2",
    "inappropriate"
  ]
}
```

After updating:
```bash
docker-compose restart backend
```

### **5. Modify Message Character Limit**

**Backend:** `backend/migrations/*_updated_messages.go`

Find and modify:
```go
&schema.SchemaField{
    Name:     "content",
    Type:     schema.FieldTypeText,
    Required: true,
    Options: &schema.TextOptions{
        Min: types.Pointer(1),
        Max: types.Pointer(240), // Change this value
    },
},
```

**Frontend:** `frontend/src/components/ContentCounter.vue` or form validation

```vue
<script setup>
const MAX_LENGTH = 240; // Change this value
</script>
```

### **6. Add Department or Update Departments**

Via PocketBase Admin UI:
1. http://localhost:8090/_/
2. Go to "college_departments" collection
3. Add/edit records with:
   - `uid`: unique code (e.g., "CS")
   - `label`: full name (e.g., "Computer Science")

---

## 🚀 Production Deployment

### **1. Prepare Production Environment**

```bash
# Create external Caddy network (one-time setup)
docker network create caddy

# Update .env for production
ENV=production
BASE_DOMAIN=yourdomain.com
BACKEND_URL=https://yourdomain.com/pb
FRONTEND_URL=https://yourdomain.com
```

### **2. Deploy Using Script**

```bash
# Deploy all services
./deploy.sh

# Deploy specific services
./deploy.sh backend frontend
```

**What happens:**
- Builds Docker images without cache
- Starts services in detached mode
- Caddy automatically provisions SSL certificates
- Services restart automatically on failure

### **3. Verify Deployment**

```bash
# Check running containers
docker ps

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f www

# Check Caddy SSL status
docker exec valentine-wall-www-1 caddy list-certs
```

### **4. Backup Database**

```bash
# Backup SQLite database
docker cp valentine-wall-backend-1:/pb_data/data.db ./backups/data_$(date +%Y%m%d).db

# Backup uploaded files
docker cp valentine-wall-backend-1:/pb_data/storage ./backups/storage_$(date +%Y%m%d)
```

### **5. Update Production**

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
./deploy.sh

# Or for zero-downtime deployment (manual)
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --no-deps --build backend
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --no-deps --build frontend
```

---

## 🐛 Troubleshooting

### **Backend Issues**

**Problem:** Backend won't start or crashes

```bash
# Check logs
docker-compose logs backend

# Common issues:
# 1. Port 8090 already in use
sudo lsof -i :8090
kill -9 <PID>

# 2. Database locked
rm pb_data/data.db-shm pb_data/data.db-wal

# 3. Permissions issues (Linux)
sudo chown -R $USER:$USER pb_data/
```

**Problem:** Images not generating

```bash
# Check headless Chrome is running
docker-compose ps headless_chrome

# Test Chrome connection
curl http://localhost:5000/json/version

# Restart Chrome
docker-compose restart headless_chrome
```

### **Frontend Issues**

**Problem:** Frontend shows blank page

```bash
# Check logs
docker-compose logs frontend

# Common causes:
# 1. Backend not reachable - check VITE_BACKEND_URL
# 2. Build failed - check for TypeScript errors
# 3. SSR hydration mismatch - check console in browser

# Clear node_modules and rebuild
docker-compose down
docker volume rm valentine-wall_node_modules (if using volumes)
docker-compose up --build frontend
```

**Problem:** Authentication not working

```bash
# Verify Firebase config in .env
# Check browser console for errors
# Verify backend URL is correct
# Clear browser cookies and try again
```

### **Docker Issues**

**Problem:** Out of disk space

```bash
# Clean up Docker
docker system prune -a --volumes

# Remove unused images
docker image prune -a
```

**Problem:** Port conflicts

```bash
# Find process using port
# Windows:
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Linux/Mac:
lsof -i :3000
kill -9 <PID>
```

---

## 📚 Additional Resources

### **Documentation**
- [PocketBase Documentation](https://pocketbase.io/docs/)
- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [TailwindCSS Documentation](https://tailwindcss.com/)
- [Docker Documentation](https://docs.docker.com/)

### **Useful Commands**

```bash
# Backend
cd backend && go fmt ./...          # Format Go code
cd backend && go mod tidy           # Clean up dependencies
cd backend && go build              # Build binary

# Frontend
cd frontend && npm run dev          # Start dev server
cd frontend && npm run build        # Build for production
cd frontend && npm run type-check   # Check TypeScript
cd frontend && npm run lint         # Lint code

# Docker
docker-compose up -d                # Start in background
docker-compose down                 # Stop all services
docker-compose restart <service>    # Restart specific service
docker-compose exec backend sh      # Shell into container
docker-compose logs -f --tail=100   # Follow logs with tail
```

---

## 🎓 Learning Path for Future Projects

### **If you want to build something similar:**

1. **Choose your backend approach:**
   - **PocketBase (like this project):** Best for rapid development, built-in auth, real-time features
   - **Traditional API:** Express.js (Node), FastAPI (Python), Spring Boot (Java)
   - **Serverless:** Firebase, Supabase, AWS Amplify

2. **Choose your frontend framework:**
   - **Vue 3 (like this project):** Progressive, easy to learn, great ecosystem
   - **React:** Most popular, huge ecosystem, more verbose
   - **Svelte:** Compiler-based, fastest runtime, smaller bundles
   - **Solid.js:** React-like syntax, better performance

3. **Essential patterns to understand:**
   - Component-based architecture
   - State management (global vs local)
   - API integration and data fetching
   - Authentication and authorization
   - Form handling and validation
   - Real-time data with WebSockets
   - Server-side rendering (SSR) vs Client-side rendering (CSR)
   - Responsive design with TailwindCSS
   - Docker containerization
   - CI/CD deployment

4. **Key files to study in this project:**
   - `backend/main.go` - Backend initialization
   - `backend/message_hooks.go` - Business logic patterns
   - `frontend/src/router.ts` - Routing setup
   - `frontend/src/client.ts` - API client configuration
   - `docker-compose.yml` - Container orchestration
   - Component examples in `frontend/src/components/`

---

**Happy Coding! 🚀**

For questions or issues, check the troubleshooting section or review the inline code comments.
