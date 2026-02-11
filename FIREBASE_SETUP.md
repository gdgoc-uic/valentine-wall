# Firebase Setup Guide - Valentine Wall

## What You're Setting Up
Firebase Analytics only - for tracking app usage in production.

**Note:** Authentication is handled by PocketBase, not Firebase.

---

## Quick Setup Steps

### 1. Create Firebase Project
1. Go to https://console.firebase.google.com/
2. Click **Add project**
3. Enter project name (e.g., "Valentine Wall")
4. **Disable** Google Analytics if prompted (or enable if you want it)
5. Click **Create project**

### 2. Register Web App
1. In your Firebase project, click **Web** icon (`</>`) to add a web app
2. Enter app nickname (e.g., "Valentine Wall Web")
3. **Skip** Firebase Hosting setup
4. Click **Register app**
5. You'll see a config object - keep this page open

### 3. Enable Google Analytics (Required for Analytics SDK)
1. In Firebase console, go to **Project settings** (gear icon)
2. Click **Integrations** tab
3. Find Google Analytics, click **Enable**
4. Follow prompts to create or link Analytics account
5. Return to **Project settings** > **General** tab to get your Measurement ID

### 4. Copy Credentials to .env File

From the Firebase config object, copy these values to your `.env` file:

```env
# Firebase Analytics (uncomment and fill in)
FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXX
FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_STORAGE_BUCKET=your-project.appspot.com
FIREBASE_MESSAGING_SENDER_ID=123456789012
FIREBASE_APP_ID=1:123456789012:web:abcdef123456
FIREBASE_MEASUREMENT_ID=G-XXXXXXXXXX
```

**Where to find these:**
- Firebase Console → Project Settings → General → Your apps → Web app → Config

### 5. Update Frontend .env

The frontend needs these with `VITE_` prefix. Edit `frontend/.env` or your root `.env`:

```env
VITE_FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXX
VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=123456789012
VITE_FIREBASE_APP_ID=1:123456789012:web:abcdef123456
VITE_FIREBASE_MEASUREMENT_ID=G-XXXXXXXXXX
```

### 6. Restart Frontend

```powershell
# If using Docker
docker-compose restart frontend

# If running locally
cd frontend
npm run dev
```

---

## Verify Setup

Analytics only runs in production. To test:

```bash
# Build for production
cd frontend
npm run build

# Serve production build
npm run preview
```

Open browser DevTools → Network tab → Filter by "google-analytics" - you should see analytics requests when navigating the app.

---

## Troubleshooting

**No analytics events?**
- Analytics is disabled in development mode (intentional)
- Check browser console for Firebase errors
- Verify all VITE_FIREBASE_* variables are set in `.env`
- Ensure you're running production build (not dev server)

**Missing Measurement ID?**
- Make sure Google Analytics is enabled in Firebase Console → Integrations

**Config not loading?**
- Check that `.env` file has the VITE_ prefix for frontend variables
- Restart the frontend development server

---

## Skip Firebase?

**You can skip Firebase entirely if:**
- ❌ You're just developing locally
- ❌ You don't need analytics
- ❌ You're not deploying to production

The app will work fine without Firebase - it's only for production analytics tracking.

---

**That's it!** Firebase Analytics will automatically track user events in production.
