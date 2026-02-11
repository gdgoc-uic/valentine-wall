# 🚀 Valentine Wall - Quick Start

Get up and running in 5 minutes!

---

## ✅ Prerequisites Done

- ✅ `.env` file created with development defaults
- ✅ DBeaver connection guide ready
- ✅ Firebase setup guide ready

---

## 🏃 Start the System (3 Steps)

### **Step 1: Prepare Data Directory**

```powershell
# Copy the template data
Copy-Item -Recurse backend\_data.please_copy backend\_data
```

### **Step 2: Start All Services**

```powershell
# Start with Docker Compose (recommended)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

**What's happening:**
- ⏳ First build takes 3-5 minutes
- 🌐 Backend will start on http://localhost:8090
- 🎨 Frontend will start on http://localhost:3000
- 🖥️ Headless Chrome for image rendering

**Watch for:**
```
backend_1   | Server started at http://0.0.0.0:8090
frontend_1  | ➜  Local:   http://localhost:3000/
```

### **Step 3: Access the Application**

Open your browser:

**Main App:**
- 🌐 http://localhost:3000

**Admin Dashboard:**
- 🔐 http://localhost:8090/_/
- **First time:** Create admin account
- Use this to manage data, users, settings

---

## 🔥 Firebase Setup (Optional - For Analytics)

Firebase is **optional** and only used for analytics in production.

### Skip Firebase if:
- ❌ You're just developing/testing locally
- ❌ You don't need analytics

### Set up Firebase if:
- ✅ You want production analytics
- ✅ You're deploying to production

**To set up Firebase:**
1. See [FIREBASE_SETUP.md](FIREBASE_SETUP.md)
2. Follow the 5 steps
3. Add credentials to `.env` file
4. Restart frontend: `docker-compose restart frontend`

---

## 💾 Connect DBeaver to SQLite

**When to do this:**
After starting the system (database is auto-created on first run)

**How to connect:**
1. See [DATABASE_CONNECTION.md](DATABASE_CONNECTION.md)
2. Database is at: `d:\valentine-wall\pb_data\data.db`
3. Use SQLite driver in DBeaver
4. View all tables and data

**Quick DBeaver Setup:**
1. New Database Connection → SQLite
2. Path: `d:\valentine-wall\pb_data\data.db`
3. Click "Download" driver if needed
4. Test Connection → Finish
5. Browse tables in Database Navigator

---

## 🎯 Next Steps

### 1. Create Admin Account
- Go to http://localhost:8090/_/
- Set up admin email and password
- This lets you access PocketBase admin UI

### 2. Add Sample Data

**Departments:**
1. Open Admin UI → Collections → `college_departments`
2. Click "New record"
3. Add:
   - uid: `CS`, label: `Computer Science`
   - uid: `IT`, label: `Information Technology`
   - uid: `ENG`, label: `Engineering`

**Gifts:**
1. Collections → `gifts`
2. Add:
   - uid: `rose`, label: `Red Rose`, price: `50`, is_remittable: `true`
   - uid: `chocolate`, label: `Chocolate Box`, price: `75`, is_remittable: `true`
   - uid: `teddy`, label: `Teddy Bear`, price: `150`, is_remittable: `true`

### 3. Create Test User
1. Go to frontend: http://localhost:3000
2. Click Sign Up
3. Fill in details (use any student ID for testing)
4. Check backend logs for verification email
5. In Admin UI, manually verify the user

### 4. Send a Test Message
1. Login with test user
2. Send a message to another student ID
3. Attach a gift
4. Check wallet balance decreased

---

## 🛑 Stop the System

```powershell
# Stop all services (Ctrl+C in terminal, then:)
docker-compose down

# Or in new terminal:
docker-compose down
```

---

## 🐛 Troubleshooting

### Backend won't start?
```powershell
# Check if port 8090 is in use
netstat -ano | findstr :8090

# Kill process if needed
taskkill /PID <PID> /F

# Restart
docker-compose restart backend
```

### Frontend won't start?
```powershell
# Check port 3000
netstat -ano | findstr :3000

# Kill process if needed
taskkill /PID <PID> /F

# Restart
docker-compose restart frontend
```

### Database not showing in DBeaver?
- Make sure backend has started at least once
- Check file exists: `d:\valentine-wall\pb_data\data.db`
- Wait for backend to fully initialize
- Refresh DBeaver connection

### Can't login/register?
- Admin UI: http://localhost:8090/_/ → Users collection
- Manually set `verified: true` on user
- Or check backend logs for verification link

---

## 📁 Important Files

```
valentine-wall/
├── .env                          ← Your configuration (created ✅)
├── pb_data/
│   └── data.db                   ← SQLite database (auto-created)
├── backend/_data/
│   └── terms-and-conditions.md   ← Terms text (from template)
├── docker-compose.yml            ← Main compose file
└── docker-compose.dev.yml        ← Development overrides
```

---

## 📚 More Help

- **Development:** See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md)
- **Commands:** See [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
- **DBeaver:** See [DATABASE_CONNECTION.md](DATABASE_CONNECTION.md)
- **Firebase:** See [FIREBASE_SETUP.md](FIREBASE_SETUP.md)
- **Architecture:** See [ARCHITECTURE.md](ARCHITECTURE.md)

---

## ✅ You're Ready!

The system is now running with:
- ✅ Backend (PocketBase) on port 8090
- ✅ Frontend (Vue 3) on port 3000
- ✅ SQLite database ready for DBeaver
- ✅ Firebase (optional, add later)

**Start developing or just explore the app!** 🎉
