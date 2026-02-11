# Valentine Wall - Setup Checklist

Use this checklist to verify your development environment is properly configured.

## ✅ Prerequisites

- [ ] Docker Desktop installed and running
- [ ] Git installed
- [ ] Code editor (VS Code recommended)
- [ ] Terminal/Command Prompt access

## ✅ Initial Setup Steps

### 1. Repository Setup
- [ ] Repository cloned: `git clone <repo-url>`
- [ ] Changed to project directory: `cd valentine-wall`
- [ ] Checked current branch: `git branch`

### 2. Environment Configuration
- [ ] Created `.env` file in project root
- [ ] Copied template configuration from `DEVELOPMENT_GUIDE.md`
- [ ] Updated `BACKEND_URL=http://localhost:8090`
- [ ] Updated `FRONTEND_URL=http://localhost:3000`
- [ ] Set `ENV=development`
- [ ] (Optional) Added Firebase credentials if using OAuth

### 3. Data Directory
- [ ] Copied data templates: `cp -r backend/_data.please_copy backend/_data`
  - **Windows PowerShell:** `Copy-Item -Recurse backend/_data.please_copy backend/_data`
- [ ] Verified `backend/_data/terms-and-conditions.md` exists

### 4. Docker Setup
- [ ] Docker Desktop is running
- [ ] Checked Docker version: `docker --version`
- [ ] Checked Docker Compose version: `docker-compose --version`

## ✅ Starting Development Environment

### 5. Start All Services
- [ ] Run: `docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build`
- [ ] Wait for all services to start (may take 3-5 minutes first time)
- [ ] Check for error messages in terminal

### 6. Verify Services

**Backend (PocketBase)**
- [ ] Open: http://localhost:8090
- [ ] Should see "404 page not found" or PocketBase API response (this is normal)
- [ ] Open Admin UI: http://localhost:8090/_/
- [ ] Create admin account when prompted
- [ ] Admin dashboard loads successfully

**Frontend (Vue + Vite)**
- [ ] Open: http://localhost:3000
- [ ] Homepage loads without errors
- [ ] Check browser console for errors (F12 → Console tab)
- [ ] Navigation works

**Headless Chrome**
- [ ] Check running: `docker-compose ps`
- [ ] Should show `headless_chrome` as "Up"

## ✅ Post-Setup Verification

### 7. PocketBase Admin Configuration

- [ ] Logged into http://localhost:8090/_/
- [ ] Verified collections exist:
  - [ ] users
  - [ ] user_details
  - [ ] messages
  - [ ] message_replies
  - [ ] gifts
  - [ ] college_departments
  - [ ] virtual_wallets
  - [ ] virtual_transactions
  - [ ] rankings

### 8. Add Sample Data (Optional)

**College Departments:**
- [ ] Go to "college_departments" collection
- [ ] Add sample records:
  ```
  uid: CS, label: Computer Science
  uid: IT, label: Information Technology
  uid: ENG, label: Engineering
  ```

**Gifts:**
- [ ] Go to "gifts" collection
- [ ] Add sample gifts:
  ```
  uid: rose, label: Red Rose, price: 50, is_remittable: true
  uid: chocolate, label: Chocolate, price: 75, is_remittable: true
  uid: teddy, label: Teddy Bear, price: 150, is_remittable: true
  ```

### 9. Test Basic Functionality

**User Registration:**
- [ ] Go to http://localhost:3000
- [ ] Find register/signup button
- [ ] Create test account
- [ ] Verify email verification email sent (check backend logs)
- [ ] Check user appears in Admin UI → users collection

**Send a Message:**
- [ ] Login with test account
- [ ] Navigate to send message page
- [ ] Fill out form (recipient student ID, content)
- [ ] Select a gift
- [ ] Send message
- [ ] Verify message created in Admin UI → messages collection
- [ ] Check virtual wallet balance decreased

**View Messages:**
- [ ] Navigate to messages/inbox page
- [ ] Verify messages display
- [ ] Check message details

### 10. Test Image Generation (Advanced)

- [ ] Send a message
- [ ] Get message ID from Admin UI
- [ ] Open: `http://localhost:8090/messages/<message-id>/image`
- [ ] Should generate and display a PNG image
- [ ] Check backend logs for Chrome rendering status

## ✅ Development Workflow Verification

### 11. Hot Reload Testing

**Frontend:**
- [ ] With frontend running, edit `frontend/src/App.vue`
- [ ] Add a comment or change text
- [ ] Save file
- [ ] Browser auto-refreshes with changes

**Backend:**
- [ ] Edit `backend/server.go`
- [ ] Add a log statement
- [ ] Save file
- [ ] Restart backend: `docker-compose restart backend`
- [ ] Check logs for your message

### 12. Database Inspection

- [ ] Open PocketBase Admin: http://localhost:8090/_/
- [ ] Can view collections
- [ ] Can create/edit/delete records
- [ ] Changes reflect immediately in frontend

## ✅ Common Issues Resolved

- [ ] If port 8090 in use: Stop conflicting process or change port
- [ ] If port 3000 in use: Stop conflicting process or change port
- [ ] If Docker containers won't start: Check Docker Desktop is running
- [ ] If database locked: Remove `.db-shm` and `.db-wal` files
- [ ] If frontend blank: Check browser console, verify VITE_BACKEND_URL
- [ ] If No data displays: Add sample data via Admin UI

## ✅ Ready for Development

If all items are checked, your environment is ready! 🎉

**Next Steps:**
1. Read `DEVELOPMENT_GUIDE.md` for detailed development instructions
2. Explore the codebase
3. Make your first change (try adding a new page or component)
4. Test your changes
5. Commit to git

---

## Quick Commands Reference

```bash
# Start everything
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Start in background
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop everything
docker-compose down

# Rebuild after changes
docker-compose up --build

# Access backend container shell
docker-compose exec backend sh

# Access frontend container shell
docker-compose exec frontend sh

# Remove all containers and volumes (fresh start)
docker-compose down -v
```

---

**Stuck?** Check the Troubleshooting section in `DEVELOPMENT_GUIDE.md`
