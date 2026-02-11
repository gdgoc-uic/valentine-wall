# Valentine Wall - Architecture Documentation

This document provides visual representations and detailed explanations of the Valentine Wall architecture.

---

## 🏗️ System Architecture Overview

### **High-Level Component Diagram**

```
┌─────────────────────────────────────────────────────────────────────┐
│                           User Browser                              │
│                    (Desktop, Mobile, Tablet)                        │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                               │ HTTPS/HTTP
                               ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    Reverse Proxy (Production)                       │
│                         Caddy Server                                │
│  ┌──────────────────┐              ┌─────────────────────────────┐  │
│  │  Route: /        │              │  Route: /pb/*               │  │
│  │  → Frontend:3000 │              │  → Backend:8090             │  │
│  └──────────────────┘              └─────────────────────────────┘  │
└──────────┬──────────────────────────────────┬──────────────────────┘
           │                                   │
           │ HTTP                              │ HTTP
           ↓                                   ↓
┌──────────────────────┐            ┌─────────────────────────────────┐
│   Frontend Service   │            │      Backend Service            │
│                      │            │                                 │
│  • Vue 3 + TypeScript│◄───────────│  • Go + PocketBase              │
│  • Vite SSR          │   REST API │  • SQLite Database              │
│  • Node.js Server    │  + WebSocket│  • Echo HTTP Router            │
│  • Port: 3000        │            │  • Port: 8090                   │
│                      │            │                                 │
│  Components:         │            │  Components:                    │
│  - Pages             │            │  - API Routes                   │
│  - Components        │            │  - Business Hooks               │
│  - Router            │            │  - Database Models              │
│  - State              │            │  - Email Service                │
│  - PocketBase Client │            │  - Image Generator              │
└──────────────────────┘            └─────────┬───────────────────────┘
                                              │
                                              │ WebSocket
                                              ↓
                                  ┌───────────────────────────┐
                                  │  Headless Chrome          │
                                  │  (Image Rendering)        │
                                  │                           │
                                  │  • Browserless Chrome     │
                                  │  • Port: 5000 (internal)  │
                                  │  • Renders HTML → PNG     │
                                  └───────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                      Persistent Storage                             │
│                                                                     │
│  • pb_data/ - SQLite database + file uploads                        │
│  • pb_public/ - Public static files                                 │
│  • backend/_data/ - Terms & conditions                              │
│  • backend/renderer_assets/ - Fonts, emojis, images                 │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Data Flow Diagrams

### **1. User Authentication Flow**

```
┌──────────┐
│  User    │
└────┬─────┘
     │
     │ 1. Click "Login with Google"
     ↓
┌─────────────────┐
│  Frontend       │
│  (Vue App)      │
└────┬────────────┘
     │
     │ 2. Redirect to Firebase
     ↓
┌─────────────────┐
│  Firebase Auth  │
│  (Google OAuth) │
└────┬────────────┘
     │
     │ 3. User authenticates
     │ 4. Returns ID token
     ↓
┌─────────────────┐
│  Frontend       │
└────┬────────────┘
     │
     │ 5. Send token to backend
     ↓
┌─────────────────┐
│  Backend        │
│  (PocketBase)   │
└────┬────────────┘
     │
     │ 6. Verify token with Firebase
     │ 7. Create/update user in SQLite
     │ 8. Generate PocketBase auth token
     ↓
┌─────────────────┐
│  Frontend       │
│  (Store token)  │
└─────────────────┘
     │
     │ 9. User authenticated
     │    Future requests include token
     ↓
```

### **2. Message Creation Flow**

```
┌──────────┐                 ┌─────────────┐
│   User   │                 │  Frontend   │
└────┬─────┘                 └──────┬──────┘
     │                              │
     │ 1. Fill form                 │
     │    - Recipient student ID    │
     │    - Message content         │
     │    - Select gifts            │
     ├──────────────────────────────►
     │                              │
     │                              │ 2. Submit via PocketBase SDK
     │                              │    pb.collection('messages').create(...)
     │                              ↓
     │                       ┌─────────────────┐
     │                       │  Backend        │
     │                       │  (PocketBase)   │
     │                       └──────┬──────────┘
     │                              │
     │                              │ 3. BEFORE CREATE HOOK
     │                              │    - Check duplicates
     │                              │    - Profanity filter
     │                              │    - Verify wallet balance
     │                              │
     │                              │ 4. Create message record
     │                              │    → SQLite INSERT
     │                              │
     │                              │ 5. AFTER CREATE HOOK
     │                              │    - Deduct coins from sender
     │                              │    - Add coins to recipient
     │                              │    - Update rankings
     │                              │    - Send email notification
     │                              │    - Log transactions
     │                              ↓
     │                       ┌─────────────┐
     │                       │  SQLite DB  │
     │                       └──────┬──────┘
     │                              │
     │                              │ 6. Return created message
     │                              ↓
     │                       ┌─────────────┐
     │                       │  Frontend   │
     │                       └──────┬──────┘
     │                              │
     │ 7. Show success              │ 8. Refetch wallet balance
     │     Update UI                │    Update message list
     ◄──────────────────────────────┤
     │                              │
```

### **3. Image Generation Flow**

```
┌──────────┐
│   User   │
└────┬─────┘
     │
     │ 1. Click "Share" on message
     ↓
┌─────────────────┐
│  Frontend       │
│  Generates URL: │
│  /messages/     │
│   {id}/image    │
└────┬────────────┘
     │
     │ 2. Request image
     ↓
┌─────────────────────────────────────┐
│  Backend - Image Generator          │
└────┬────────────────────────────────┘
     │
     │ 3. Check in-memory cache (10min TTL)
     ├─► [Cache Hit] ──► Return cached PNG
     │
     │ [Cache Miss]
     ↓
     │ 4. Load message from database
     │ 5. Render HTML template with data
     ↓
┌─────────────────┐
│  Send HTML via  │
│  WebSocket to   │
│  Chrome         │
└────┬────────────┘
     │
     ↓
┌─────────────────────────────┐
│  Headless Chrome            │
│                             │
│  1. Load HTML               │
│  2. Render page             │
│  3. Take screenshot         │
│  4. Return PNG (1200x675)   │
└────┬────────────────────────┘
     │
     ↓
┌─────────────────┐
│  Backend        │
│  - Cache result │
│  - Return PNG   │
└────┬────────────┘
     │
     ↓
┌─────────────────┐
│  User sees      │
│  social media   │
│  card image     │
└─────────────────┘
```

---

## 🗄️ Database Schema

### **Entity Relationship Diagram**

```
┌──────────────────┐
│     users        │
│  (PocketBase)    │
├──────────────────┤
│ id (PK)          │
│ email            │
│ verified         │
│ created          │
│ updated          │
└────┬─────────────┘
     │
     │ 1:1
     ↓
┌──────────────────┐
│  user_details    │
├──────────────────┤
│ id (PK)          │
│ user (FK) ───────┼─► users.id
│ student_id       │
│ college_dept ────┼─► college_departments.id
│ sex              │
└──────────────────┘
     │
     │ 1:1
     ↓
┌──────────────────────┐
│  virtual_wallets     │
├──────────────────────┤
│ id (PK)              │
│ user (FK) ───────────┼─► users.id
│ balance (default:    │
│          1000)       │
└──────┬───────────────┘
       │
       │ 1:N
       ↓
┌──────────────────────────┐
│  virtual_transactions    │
├──────────────────────────┤
│ id (PK)                  │
│ wallet (FK) ─────────────┼─► virtual_wallets.id
│ amount                   │
│ description              │
│ created                  │
└──────────────────────────┘

┌──────────────────┐
│     users        │
└────┬─────────────┘
     │
     │ 1:N (sender)
     ↓
┌──────────────────────┐
│    messages          │
├──────────────────────┤
│ id (PK)              │
│ user (FK) ───────────┼─► users.id (sender)
│ recipient            │  (student_id reference)
│ content (max 240)    │
│ gifts[] (array)  ────┼─► gifts.uid[]
│ deleted              │
│ created              │
└──────┬───────────────┘
       │
       │ 1:N
       ↓
┌──────────────────────┐
│  message_replies     │
├──────────────────────┤
│ id (PK)              │
│ message (FK) ────────┼─► messages.id
│ user (FK) ───────────┼─► users.id
│ content              │
│ created              │
└──────────────────────┘

┌──────────────────┐
│     gifts        │
├──────────────────┤
│ id (PK)          │
│ uid (unique)     │
│ label            │
│ price            │
│ is_remittable    │
└──────────────────┘

┌────────────────────────┐
│  college_departments   │
├────────────────────────┤
│ id (PK)                │
│ uid (unique)           │
│ label                  │
└────────────────────────┘

┌──────────────────┐
│    rankings      │
│  (computed)      │
├──────────────────┤
│ id (PK)          │
│ recipient        │  (student_id)
│ college_dept     │
│ sex              │
│ total_coins      │  (sum of gifts received)
└──────────────────┘
```

### **Collection Details**

#### **users** (PocketBase Auth Collection)
- Managed by PocketBase authentication system
- Stores credentials, email verification status
- Auto-created on registration

#### **user_details**
- Extended user information
- Links to `college_departments` for department info
- `student_id` must be unique (used for message recipient)

#### **messages**
- Core message entity
- `gifts` is JSON array of gift UIDs
- `recipient` references `student_id` (not user.id)
- Soft delete with `deleted` timestamp

#### **virtual_wallets**
- One wallet per user
- Initial balance: 1000 coins
- Updated via transactions

#### **virtual_transactions**
- Immutable transaction log
- Types: "Initial balance", "Message sent", "Gift received", etc.
- Amount can be positive (receive) or negative (spend)

#### **rankings**
- Computed/aggregated data
- Groups by recipient, department, sex
- `total_coins` = sum of all remittable gift values received

---

## 🔌 API Architecture

### **REST API Endpoints**

#### **PocketBase Auto-Generated (CRUD)**

```
Base URL: http://localhost:8090/api

Authentication:
POST   /users/auth-with-password          → Login
POST   /users/auth-refresh                → Refresh token
POST   /users/request-verification        → Request email verification
POST   /users/confirm-verification        → Confirm email

Collections (same pattern for all):
GET    /collections/{name}/records        → List records
GET    /collections/{name}/records/{id}   → Get one record
POST   /collections/{name}/records        → Create record
PATCH  /collections/{name}/records/{id}   → Update record
DELETE /collections/{name}/records/{id}   → Delete record

Available collections:
- users
- user_details
- messages
- message_replies
- gifts
- college_departments
- virtual_wallets
- virtual_transactions
- rankings
```

#### **Custom Endpoints**

```
GET    /terms-and-conditions              → Get T&C markdown
GET    /departments                       → List all departments
GET    /gifts                             → List all gifts
GET    /messages/:id/image                → Generate message image (PNG)
GET    /messages/:id/image?template_image → Get HTML template
GET    /user_messages/archive             → Generate ZIP archive (SSE)
GET    /user_messages/download_archive/:userId → Download archive
GET    /user_auth/callback                → OAuth callback
```

### **Real-time Subscriptions (WebSocket)**

```javascript
// Connect
const pb = new PocketBase('http://localhost:8090');

// Subscribe to record changes
pb.collection('messages').subscribe(messageId, (data) => {
  console.log(data.action); // 'create', 'update', 'delete'
  console.log(data.record); // Updated record
});

// Subscribe to all records in collection
pb.collection('virtual_wallets').subscribe('*', (data) => {
  // Handle any wallet update
});

// Unsubscribe
pb.collection('messages').unsubscribe(messageId);
```

---

## 🏃 Application Flow Sequences

### **Complete User Journey: Send Message**

```
User Action                  Frontend                 Backend                  Database
─────────────────────────────────────────────────────────────────────────────────────

1. Visit homepage ──────► Load Vue app
                          Router: /
                                │
                                │ Fetch departments
                                │ & gifts on mount
                                └──────────────────► GET /departments
                                                     GET /gifts
                                                                                  │
                                                     Query collections ──────────►
                                                     Return data ◄────────────────
                       ◄─────────────────────────── JSON response
   Display form

2. Fill out form
   - Recipient ID
   - Message text
   - Select gifts

3. Click "Send" ────────► Validate locally
                          Check balance
                                │
                                │ Submit message
                                └──────────────────► POST /collections/
                                                            messages/records
                                                     {
                                                       recipient: "...",
                                                       content: "...",
                                                       gifts: [...]
                                                     }
                                                                                  │
                                                     BEFORE CREATE HOOK:          │
                                                     1. Check duplicates ────────►│
                                                     2. Profanity filter          │
                                                     3. Verify balance ──────────►│
                                                                                  │
                                                     Create record ──────────────►│
                                                     INSERT INTO messages         │
                                                                                  │
                                                     AFTER CREATE HOOK:           │
                                                     1. Deduct sender coins ─────►│
                                                        UPDATE wallets            │
                                                     2. Add recipient coins ─────►│
                                                     3. Log transactions ────────►│
                                                        INSERT INTO transactions  │
                                                     4. Update rankings ─────────►│
                                                     5. Send email (async)        │
                                                                                  │
                       ◄─────────────────────────── Return created message ◄──────
   Show success
   notification

4. UI updates ──────────► Refetch wallet
                          balance (React Query
                          invalidation)
                                └──────────────────► GET /collections/
                                                            virtual_wallets/
                                                            records
                       ◄─────────────────────────── Updated balance

   Display new
   balance
```

---

## 🔐 Security Architecture

### **Authentication Layers**

```
┌─────────────────────────────────────────────────────────┐
│                    Request Flow                         │
└─────────────────────────────────────────────────────────┘

1. User Login
   ↓
2. Frontend gets JWT from PocketBase
   ↓
3. Store token in browser (localStorage via PocketBase SDK)
   ↓
4. Include token in all requests:
   Authorization: Bearer <token>
   ↓
5. Backend validates token on each request
   ↓
6. Check collection-level permissions (API rules)
   ↓
7. Execute business logic hooks
   ↓
8. Return response
```

### **PocketBase API Rules**

Collection rules use a custom filter syntax:

```javascript
// Examples of rule expressions:

// Public read (anyone can view)
listRule: ""

// Authenticated users only
listRule: "@request.auth.id != ''"

// Owner only (user must match record's user field)
updateRule: "@request.auth.id = user.id"

// Specific field check
createRule: "@request.auth.verified = true"

// Relation check (user owns the message being replied to)
createRule: "@request.auth.id = message.user.id"

// No one (not even admins via API)
deleteRule: null
```

### **Security Features**

1. **Email Verification Required**
   - Users must verify email before full access
   - Configurable per collection

2. **Rate Limiting**
   - Built into PocketBase
   - Prevents spam and abuse

3. **Profanity Filter**
   - Server-side content moderation
   - Multi-language support
   - Blocks inappropriate content

4. **SQL Injection Protection**
   - PocketBase uses parameterized queries
   - SQLite prepared statements

5. **XSS Protection**
   - Vue 3 auto-escapes template output
   - Sanitization for user content

6. **CORS Configuration**
   - Configured in backend
   - Restricts allowed origins

---

## 📦 Deployment Architecture

### **Development Environment**

```
┌─────────────────────────────────────────────┐
│          Docker Network: backend            │
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │  headless_chrome                     │  │
│  │  Port: 5000 (internal only)          │  │
│  └──────────────────────────────────────┘  │
│                    ▲                        │
│                    │ ws://                  │
│  ┌─────────────────┴────────────────────┐  │
│  │  backend                             │  │
│  │  Port: 8090 → 8090 (exposed)         │  │
│  │  Volumes:                             │  │
│  │    - ./pb_data                        │  │
│  │    - ./backend/_data                  │  │
│  └──────────────────────────────────────┘  │
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │  frontend                            │  │
│  │  Port: 3000 → 3000 (exposed)         │  │
│  │  Vite dev server with HMR            │  │
│  └──────────────────────────────────────┘  │
└─────────────────────────────────────────────┘

Developer accesses:
- Frontend: http://localhost:3000
- Backend: http://localhost:8090
```

### **Production Environment**

```
┌─────────────────────────────────────────────────────┐
│       Docker Network: caddy (external)              │
│                                                     │
│  ┌──────────────────────────────────────────────┐  │
│  │  www (Caddy Reverse Proxy)                   │  │
│  │  Ports: 80:80, 443:443                       │  │
│  │  Auto SSL via Let's Encrypt                  │  │
│  │                                               │  │
│  │  Routes:                                      │  │
│  │    / → frontend:3000                         │  │
│  │    /pb/* → backend:8090                      │  │
│  └────────┬─────────────────────┬────────────────┘  │
└───────────┼─────────────────────┼────────────────────┘
            │                     │
            │                     │
┌───────────┼─────────────────────┼────────────────────┐
│           │  Docker Network: backend                 │
│           │                     │                    │
│  ┌────────▼─────────┐  ┌────────▼─────────────────┐ │
│  │  frontend        │  │  backend                 │ │
│  │  Port: 3000      │  │  Port: 8090              │ │
│  │  (internal)      │  │  (internal)              │ │
│  └──────────────────┘  └───────┬──────────────────┘ │
│                                │                    │
│                                │ ws://              │
│                       ┌────────▼──────────────────┐ │
│                       │  headless_chrome          │ │
│                       │  Port: 5000 (internal)    │ │
│                       └───────────────────────────┘ │
└─────────────────────────────────────────────────────┘

Users access:
- https://yourdomain.com → Frontend
- https://yourdomain.com/pb → Backend API
```

---

## 🔄 State Management

### **Frontend State Architecture**

```
┌──────────────────────────────────────────────────────┐
│                   Vue Application                    │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  Global State (Vuex)                           │ │
│  │  - Legacy store (minimal usage)                │ │
│  └────────────────────────────────────────────────┘ │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  Custom Reactive Stores (Composition API)      │ │
│  │                                                 │ │
│  │  Store (store_new.ts):                        │ │
│  │    - modals (open/close states)               │ │
│  │    - gifts (cached list)                      │ │
│  │    - departments (cached list)                │ │
│  │                                                 │ │
│  │  AuthStore (auth.ts):                         │ │
│  │    - currentUser                              │ │
│  │    - isAuthenticated                          │ │
│  │    - login/logout methods                     │ │
│  └────────────────────────────────────────────────┘ │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  TanStack Query (Vue Query)                    │ │
│  │  - Server state caching                        │ │
│  │  - Automatic refetching                        │ │
│  │  - Optimistic updates                          │ │
│  │                                                 │ │
│  │  Query Keys:                                   │ │
│  │    ['messages', userId]                       │ │
│  │    ['message', messageId]                     │ │
│  │    ['wallet', userId]                         │ │
│  │    ['gifts']                                   │ │
│  │    ['departments']                             │ │
│  └────────────────────────────────────────────────┘ │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  Component Local State                         │ │
│  │  - ref(), reactive() in <script setup>        │ │
│  │  - Props and emits                             │ │
│  └────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

### **When to Use Each State Type**

| State Type | Use Case | Example |
|------------|----------|---------|
| **Vue Query** | Server data, auto-refetch | Messages, wallet balance |
| **Custom Store** | App-wide UI state | Modals, cached reference data |
| **Vuex** | Legacy compatibility | Minimal use, consider migrating |
| **Local State** | Component-specific | Form inputs, local toggles |
| **Props/Emits** | Parent-child communication | Component data flow |

---

## 🎯 Performance Optimizations

### **Frontend Optimizations**

1. **Code Splitting**
   - Lazy-loaded routes: `() => import('./pages/Home.vue')`
   - Reduces initial bundle size

2. **SSR (Server-Side Rendering)**
   - Faster First Contentful Paint (FCP)
   - Better SEO
   - Hydration on client

3. **Vue Query Caching**
   - Reduces redundant API calls
   - Configurable stale time
   - Background refetching

4. **Component Lazy Loading**
   - Heavy components loaded on demand
   - `defineAsyncComponent()`

5. **Asset Optimization**
   - Vite automatic image optimization
   - Font subsetting
   - Tree-shaking for unused code

### **Backend Optimizations**

1. **In-Memory Caching**
   - Image generation results (10min TTL)
   - Reduces Chrome rendering load

2. **Database Indexes**
   - PocketBase auto-indexes relations
   - Custom indexes on frequently queried fields

3. **SQLite Optimizations**
   - WAL mode for better concurrency
   - Prepared statements

4. **Efficient Hooks**
   - Minimal database queries in hooks
   - Batch operations where possible

5. **Static Asset CDN**
   - Fonts and images can be served via CDN
   - Reduce backend load

---

## 📊 Monitoring & Logging

### **Application Logs**

```bash
# Backend logs (PocketBase)
docker-compose logs -f backend

# Look for:
# - API requests (GET, POST, PATCH, DELETE)
# - Hook executions
# - Email delivery status
# - Image generation timing
# - Database errors

# Frontend logs (Vue + Vite)
docker-compose logs -f frontend

# Look for:
# - Build errors
# - SSR rendering errors
# - Server startup
```

### **Browser Console**

```
F12 → Console Tab

# Frontend logging:
# - API errors
# - Vue warnings
# - Navigation events
# - State changes
```

### **PocketBase Admin Logs**

Access via: http://localhost:8090/_/logs

- View all API requests
- Filter by collection
- See request/response details
- Export logs

---

## 🔮 Future Enhancements

Potential architecture improvements:

1. **Caching Layer**
   - Add Redis for distributed caching
   - Cache API responses
   - Session storage

2. **Message Queue**
   - Use RabbitMQ or Redis for async tasks
   - Email sending queue
   - Image generation queue

3. **PostgreSQL Migration**
   - For horizontal scaling
   - Better full-text search
   - More advanced queries

4. **CDN Integration**
   - CloudFlare or similar
   - Serve static assets
   - Image optimization

5. **Monitoring**
   - Prometheus + Grafana
   - Application metrics
   - Performance tracking

6. **Testing**
   - Unit tests (Vitest)
   - E2E tests (Playwright)
   - API tests (Go testing)

---

**For implementation details, see [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md)**
