# 💌 Valentine Wall

A full-stack web application for sending anonymous Valentine's Day messages with virtual gifts, built with modern technologies and containerized for easy deployment.

![Tech Stack](https://img.shields.io/badge/Frontend-Vue%203-4FC08D?style=flat&logo=vue.js&logoColor=white)
![Tech Stack](https://img.shields.io/badge/Backend-Go%20%2B%20PocketBase-00ADD8?style=flat&logo=go&logoColor=white)
![Tech Stack](https://img.shields.io/badge/Database-SQLite-003B57?style=flat&logo=sqlite&logoColor=white)
![Tech Stack](https://img.shields.io/badge/Container-Docker-2496ED?style=flat&logo=docker&logoColor=white)

---

## 📖 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Quick Start](#-quick-start)
- [Documentation](#-documentation)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Production Deployment](#-production-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

- **✉️ Anonymous Messaging** - Send Valentine's messages to anyone in the system
- **🎁 Virtual Gifts** - Attach gifts to messages using virtual currency
- **💰 Virtual Economy** - Users start with 1000 coins, spend to send messages
- **🏆 Ranking System** - Track most popular recipients by coins received
- **📧 Email Notifications** - Recipients get notified of new messages
- **🖼️ Social Media Cards** - Automatically generated shareable images for each message
- **🔐 Secure Authentication** - Email verification + optional Firebase OAuth (Google/Facebook)
- **📱 Responsive Design** - Works on desktop, tablet, and mobile
- **⚡ Real-time Updates** - WebSocket subscriptions for live data updates
- **🗂️ Message Archiving** - Download all received messages as a ZIP file
- **🚫 Profanity Filter** - Multi-language content moderation
- **🎨 Dark Mode Ready** - DaisyUI theming support

---

## 🛠️ Tech Stack

### Backend
- **Framework:** PocketBase 0.11.3 (Go-based BaaS)
- **Language:** Go 1.19+
- **Database:** SQLite (embedded)
- **HTTP Router:** Echo v5
- **Authentication:** PocketBase Auth + Firebase Admin SDK
- **Image Rendering:** Chrome DevTools Protocol (headless Chrome)
- **Email:** SMTP via PocketBase

### Frontend
- **Framework:** Vue 3.2.45 (Composition API with `<script setup>`)
- **Language:** TypeScript 4.9+
- **Build Tool:** Vite 4.0.0
- **Routing:** Vue Router 4.1.6
- **State Management:** Vuex 4 + Custom Reactive Stores
- **Data Fetching:** TanStack Query (Vue Query) 4.24.4
- **Backend SDK:** PocketBase JS SDK 0.10.1
- **UI Framework:** TailwindCSS 3.2.4 + DaisyUI 2.49.0
- **SSR:** Yes (Node.js + Polka server)

### Infrastructure
- **Containerization:** Docker + Docker Compose
- **Reverse Proxy:** Caddy 2 (auto-SSL for production)
- **Image Renderer:** Browserless Chrome

---

## 🚀 Quick Start

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop) (includes Docker Compose)
- Git
- Text editor (VS Code recommended)

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/valentine-wall.git
cd valentine-wall
```

### 2. Create Environment File

Create a `.env` file in the project root:

```bash
ENV=development
BASE_DOMAIN=localhost
BACKEND_URL=http://localhost:8090
FRONTEND_URL=http://localhost:3000
READ_ONLY=false
PROFANITY_JSON_FILE_NAME=profanities.json

# Optional: Firebase OAuth credentials
FIREBASE_API_KEY=
FIREBASE_AUTH_DOMAIN=
FIREBASE_PROJECT_ID=
FIREBASE_STORAGE_BUCKET=
FIREBASE_MESSAGING_SENDER_ID=
FIREBASE_APP_ID=
FIREBASE_MEASUREMENT_ID=
```

### 3. Prepare Data Directory

```bash
# Linux/Mac
cp -r backend/_data.please_copy backend/_data

# Windows PowerShell
Copy-Item -Recurse backend/_data.please_copy backend/_data
```

### 4. Start Development Environment

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

**First startup takes 3-5 minutes to build images.**

### 5. Access the Application

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8090
- **Admin Dashboard:** http://localhost:8090/_/

**Create an admin account** when first accessing the dashboard.

---

## 📚 Documentation

We've created comprehensive documentation to help you get started and develop effectively:

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **[SETUP_CHECKLIST.md](SETUP_CHECKLIST.md)** | Step-by-step setup verification | Setting up for the first time |
| **[DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md)** | Complete development guide | Learning the codebase, making changes |
| **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** | Common tasks & commands | Daily development, quick lookups |

### 📖 What Each Guide Contains

**SETUP_CHECKLIST.md** - Your First 30 Minutes
- ✅ Prerequisites checklist
- ✅ Step-by-step initial setup
- ✅ Service verification
- ✅ Sample data setup
- ✅ Common issues & fixes

**DEVELOPMENT_GUIDE.md** - Your Complete Reference (40+ pages)
- 🏗️ Architecture & how components interact
- 🔧 Making changes in backend (hooks, endpoints, models)
- 🎨 Making changes in frontend (pages, components, API calls)
- 📝 Common development tasks with examples
- 🚀 Production deployment guide
- 🐛 Troubleshooting section
- 🎓 Learning path for future projects

**QUICK_REFERENCE.md** - Your Daily Companion
- ⚡ Quick command reference
- 📁 File locations cheat sheet
- 🔍 Debugging tips
- 🎨 Common TailwindCSS patterns
- 🗄️ Database access patterns
- 💡 Code snippets for common tasks

---

## 📁 Project Structure

```
valentine-wall/
├── backend/                   # Go + PocketBase backend
│   ├── main.go               # Entry point
│   ├── server.go             # Custom routes
│   ├── *_hooks.go            # Business logic (messages, users, wallets)
│   ├── models/               # Database models
│   ├── migrations/           # Database migrations
│   ├── templates/            # Email & HTML templates
│   └── renderer_assets/      # Fonts, emojis, images
│
├── frontend/                  # Vue 3 + TypeScript frontend
│   ├── src/
│   │   ├── main.ts           # Client entry
│   │   ├── App.vue           # Root component
│   │   ├── router.ts         # Routes
│   │   ├── client.ts         # PocketBase client
│   │   ├── components/       # Reusable components
│   │   ├── pages/            # Page components
│   │   └── assets/           # Styles, images
│   ├── vite.config.ts        # Vite configuration
│   └── package.json          # Dependencies
│
├── docker-compose.yml         # Base Docker configuration
├── docker-compose.dev.yml     # Development overrides
├── docker-compose.prod.yml    # Production overrides
├── deploy.sh                  # Production deployment script
├── .env                       # Environment variables (you create this)
├── profanities.json           # Profanity filter list
│
├── DEVELOPMENT_GUIDE.md       # 📖 Complete development guide
├── SETUP_CHECKLIST.md         # ✅ Setup verification checklist
├── QUICK_REFERENCE.md         # ⚡ Quick reference card
└── README.md                  # 👈 You are here
```

---

## 💻 Development

### Development Workflow

**Option 1: Docker Compose (Recommended for beginners)**
```bash
# Start all services
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Restart a service
docker-compose restart backend
```

**Option 2: Local Scripts (Faster hot reload)**
```bash
# Install dotenv CLI
npm install -g dotenv-cli

# Terminal 1 - Backend
./start-backend.sh serve --http=0.0.0.0:8090 --dir=./pb_data

# Terminal 2 - Frontend
./start-frontend.sh
```

### Making Changes

**Backend (Go):**
- Custom routes → `backend/server.go`
- Business logic → `backend/*_hooks.go`
- Email templates → `backend/templates/mail/*.txt.tpl`

**Frontend (Vue 3):**
- New page → `frontend/src/pages/YourPage.vue` + add to `router.ts`
- Component → `frontend/src/components/YourComponent.vue`
- API calls → Use PocketBase SDK via `client.ts`

**See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) for detailed examples.**

### Admin Dashboard

Access PocketBase admin at http://localhost:8090/_/

- View/edit all collections
- Manage users
- Test API endpoints
- Configure collection rules
- View logs

---

## 🚀 Production Deployment

### Prerequisites

```bash
# Create external Caddy network (one-time)
docker network create caddy

# Update .env for production
ENV=production
BASE_DOMAIN=yourdomain.com
BACKEND_URL=https://yourdomain.com/pb
FRONTEND_URL=https://yourdomain.com
```

### Deploy

```bash
# Deploy all services
./deploy.sh

# Deploy specific services
./deploy.sh backend frontend
```

### Production Architecture

```
Internet → Caddy (port 80/443, auto-SSL)
            ├─→ Frontend (/)
            └─→ Backend (/pb/*)
                  └─→ Headless Chrome (internal)
```

**See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md#-production-deployment) for full deployment guide.**

---

## 🧪 Testing

### Manual Testing

```bash
# Test API endpoint
curl http://localhost:8090/api/collections/messages/records

# Test with authentication
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8090/api/collections/messages/records
```

### Frontend Type Checking

```bash
cd frontend
npm run type-check
```

### Build Testing

```bash
# Frontend production build
cd frontend && npm run build

# Backend build
cd backend && go build
```

---

## 🐛 Troubleshooting

**Services won't start:**
```bash
# Check Docker is running
docker ps

# View logs
docker-compose logs

# Clean restart
docker-compose down
docker-compose up --build
```

**Port already in use:**
```bash
# Windows
netstat -ano | findstr :8090
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8090
kill -9 <PID>
```

**More troubleshooting:** See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md#-troubleshooting)

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Commit: `git commit -m 'Add amazing feature'`
5. Push: `git push origin feature/amazing-feature`
6. Open a Pull Request

**Please read [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) to understand the codebase structure.**

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 🙏 Acknowledgments

- [PocketBase](https://pocketbase.io/) - Amazing Go backend framework
- [Vue.js](https://vuejs.org/) - Progressive JavaScript framework
- [Vite](https://vitejs.dev/) - Next generation frontend tooling
- [TailwindCSS](https://tailwindcss.com/) - Utility-first CSS framework
- [DaisyUI](https://daisyui.com/) - Tailwind component library

---

## 📞 Support

- **Documentation:** Check [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md), [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
- **Setup Issues:** See [SETUP_CHECKLIST.md](SETUP_CHECKLIST.md)
- **Questions:** Open an issue on GitHub

---

<div align="center">

**Made with ❤️ for Valentine's Day**

⭐ Star this repo if you find it helpful!

[Report Bug](https://github.com/yourusername/valentine-wall/issues) · [Request Feature](https://github.com/yourusername/valentine-wall/issues)

</div>
