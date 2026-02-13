# 🚀 Valentine Wall - CloudPanel Deployment Guide

**Step-by-step guide for deploying Valentine Wall using CloudPanel on your VPS**

---

## What You Need

- ✅ **VPS with CloudPanel installed** (Your Contabo VPS - 194.233.73.243)
- ✅ **Domain in Cloudflare** (gdgocuic.org)
- ✅ **Docker installed** (Already have Docker 29.2.1)
- ✅ **SSH access** to your VPS

---

## Key Differences with CloudPanel

**Instead of using Caddy (like the standard guide), we'll use:**
- ✅ CloudPanel's **Nginx** as reverse proxy
- ✅ CloudPanel's **SSL certificate management** (Let's Encrypt)
- ✅ Docker Compose **without Caddy service**
- ✅ Containers expose ports to **localhost only** (127.0.0.1)

---

## 1. Access VPS and Prepare Directory

```bash
# SSH into VPS
ssh root@194.233.73.243

# Create app directory in CloudPanel's web root
mkdir -p /home/cloudpanel/htdocs/valentine-wall
cd /home/cloudpanel/htdocs/valentine-wall
```

---

## 2. Setup GitHub SSH Access (Already Done)

You already have SSH keys set up. Just verify:

```bash
# Check your SSH config
cat ~/.ssh/config

# Should show:
# Host github-valentine
#     HostName github.com
#     User git
#     IdentityFile ~/.ssh/github_deploy  # ← Make sure this is correct!
#     IdentitiesOnly yes

# Test connection
ssh -T git@github-valentine
```

**✅ Expected:** "Hi username! You've successfully authenticated..."

---

## 3. Clone Repository

```bash
cd /home/cloudpanel/htdocs/valentine-wall

# Clone using SSH
git clone git@github-valentine:yourusername/valentine-wall.git .
```

---

## 4. Configure Environment

Create `.env` file:

```bash
nano .env
```

Add this configuration:

```env
# Domain Configuration
BASE_DOMAIN=valentine-wall.gdgocuic.org
BACKEND_URL=https://valentine-wall.gdgocuic.org/pb
FRONTEND_URL=https://valentine-wall.gdgocuic.org

# Environment
ENV=production
NODE_ENV=production

# Profanity filter
PROFANITY_JSON_FILE_NAME=profanities.json

# Firebase (Optional - leave empty if not using)
FIREBASE_API_KEY=
FIREBASE_AUTH_DOMAIN=
FIREBASE_PROJECT_ID=
FIREBASE_STORAGE_BUCKET=
FIREBASE_MESSAGING_SENDER_ID=
FIREBASE_APP_ID=
FIREBASE_MEASUREMENT_ID=

# Report API (Optional - leave empty if not using)
REPORT_API_URL=
REPORT_API_KEY=
REPORT_API_CATEGORY_ID_KEY=
```

Save: `Ctrl+X`, `Y`, `Enter`

---

## 5. Create CloudPanel-Specific Docker Compose

Create `docker-compose.cloudpanel.yml`:

```bash
nano docker-compose.cloudpanel.yml
```

Add this content:

```yaml
version: '3.8'

services:
  backend:
    ports:
      - "127.0.0.1:8090:8090"  # Only localhost - Nginx will proxy
    environment:
      - BASE_DOMAIN=${BASE_DOMAIN}
    volumes:
      - ./pb_data:/pb/pb_data
      - ./pb_public:/pb/pb_public
    restart: unless-stopped

  frontend:
    ports:
      - "127.0.0.1:3000:3000"  # Only localhost - Nginx will proxy
    environment:
      - BASE_DOMAIN=${BASE_DOMAIN}
      - BACKEND_URL=${BACKEND_URL}
      - VITE_FIREBASE_API_KEY=${VITE_FIREBASE_API_KEY:-}
      - VITE_FIREBASE_AUTH_DOMAIN=${VITE_FIREBASE_AUTH_DOMAIN:-}
      - VITE_FIREBASE_PROJECT_ID=${VITE_FIREBASE_PROJECT_ID:-}
      - VITE_FIREBASE_STORAGE_BUCKET=${VITE_FIREBASE_STORAGE_BUCKET:-}
      - VITE_FIREBASE_MESSAGING_SENDER_ID=${VITE_FIREBASE_MESSAGING_SENDER_ID:-}
      - VITE_FIREBASE_APP_ID=${VITE_FIREBASE_APP_ID:-}
    depends_on:
      - backend
    restart: unless-stopped

# Note: No Caddy service - CloudPanel's Nginx handles reverse proxy and SSL
```

Save: `Ctrl+X`, `Y`, `Enter`

---

## 6. Deploy Application with Docker

```bash
# Build and start containers
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml up -d --build
```

**⏱️ First build takes 10-15 minutes.**

Monitor the build:
```bash
# Watch logs
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml logs -f

# Check running containers
docker ps
```

**✅ Expected:** Two containers running:
- `valentine-wall-backend-1` on 127.0.0.1:8090
- `valentine-wall-frontend-1` on 127.0.0.1:3000

---

## 7. Create Site in CloudPanel

### 7.1 Login to CloudPanel

1. Open browser: `https://194.233.73.243:8443`
2. Login with your CloudPanel credentials

### 7.2 Add New Site

1. Click **"Sites"** → **"Add Site"**
2. Configure:
   - **Domain Name:** `valentine-wall.gdgocuic.org`
   - **Site Type:** Select **"Generic"** or **"Node.js"** (doesn't matter, we're using custom Nginx config)
   - **PHP Version:** N/A (not needed)
3. Click **"Create"**

---

## 8. Configure SSL Certificate

### 8.1 Issue Let's Encrypt Certificate

1. In CloudPanel, go to your site **"valentine-wall.gdgocuic.org"**
2. Click **"SSL/TLS"** tab
3. Click **"New Let's Encrypt Certificate"**
4. Select: ✅ `valentine-wall.gdgocuic.org`
5. Click **"Create and Install"**

**✅ Wait 30 seconds** for certificate issuance.

---

## 9. Configure Nginx Reverse Proxy

### 9.1 Edit Nginx Configuration

```bash
# Edit the vhost file
nano /etc/nginx/sites-enabled/valentine-wall.gdgocuic.org.conf
```

### 9.2 Replace Content with This Configuration

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name valentine-wall.gdgocuic.org;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name valentine-wall.gdgocuic.org;

    # SSL Certificate (managed by CloudPanel)
    ssl_certificate /etc/letsencrypt/live/valentine-wall.gdgocuic.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/valentine-wall.gdgocuic.org/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # PocketBase API (needed for admin panel JS calls to /api/*)
    location /api/ {
        proxy_pass http://127.0.0.1:8090/api/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # PocketBase admin panel (rewrite /pb/_/ → /_/)
    location /pb/_/ {
        proxy_pass http://127.0.0.1:8090/_/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Backend admin panel (direct /_/ access)
    location /_/ {
        proxy_pass http://127.0.0.1:8090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Backend API (PocketBase) - /pb prefix
    location /pb {
        proxy_pass http://127.0.0.1:8090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Increase limits for file uploads
        client_max_body_size 10M;
        proxy_read_timeout 300;
        proxy_connect_timeout 300;
        proxy_send_timeout 300;
    }

    # Frontend (Root)
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Logs
    access_log /home/cloudpanel/htdocs/valentine-wall.gdgocuic.org/logs/access.log;
    error_log /home/cloudpanel/htdocs/valentine-wall.gdgocuic.org/logs/error.log;
}
```

Save: `Ctrl+X`, `Y`, `Enter`

### 9.3 Test and Reload Nginx

```bash
# Test configuration
nginx -t

# If OK, reload Nginx
systemctl reload nginx
```

**✅ Expected:** "syntax is ok" and "test is successful"

---

## 10. DNS Configuration (Already Done ✅)

Your Cloudflare DNS is already set up correctly:
- **Type:** A
- **Name:** valentine-wall
- **IPv4:** 194.233.73.243
- **Proxy:** Enabled (orange cloud) ✅
- **TTL:** Auto

### Verify Cloudflare SSL Mode

1. Login to Cloudflare
2. Select domain: **gdgocuic.org**
3. Go to **SSL/TLS** tab
4. Set to: **Full (strict)** ✅

---

## 11. Configure PocketBase

### 11.1 Access Admin Panel

Open browser: `https://valentine-wall.gdgocuic.org/pb/_/`

**First-time setup:**
1. Create admin account
2. Set admin email and password
3. Click "Create admin"

### 11.2 Configure SMTP (Email Sending)

1. In admin panel, go to **"Settings"** → **"Mail settings"**
2. Configure your SMTP provider:

**Example for Gmail:**
```
SMTP Host: smtp.gmail.com
Port: 587
Username: your-email@gmail.com
Password: your-app-password
TLS: Enabled
From Address: your-email@gmail.com
From Name: Valentine Wall
```

3. Click **"Send test email"** to verify
4. Click **"Save changes"**

---

## 12. Test Everything

### 12.1 Test Frontend

```bash
# Visit the site
open https://valentine-wall.gdgocuic.org
```

**✅ Expected:** Valentine Wall homepage loads with SSL (green padlock)

### 12.2 Test Backend API

```bash
curl https://valentine-wall.gdgocuic.org/pb/api/health
```

**✅ Expected:** `{"code":200,"data":{}}`

### 12.3 Test User Registration

1. Go to: `https://valentine-wall.gdgocuic.org`
2. Click **"Sign Up"**
3. Fill form and submit
4. **✅ Check email** for verification link
5. Click verification link
6. Login with credentials

### 12.4 Test Message Sending

1. Login to app
2. Create a new message
3. **✅ Verify:** Message appears on wall
4. **✅ Verify:** Recipient gets email notification

---

## 13. Maintenance Commands

### View Logs

```bash
# Docker logs
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml logs -f

# Nginx logs
tail -f /home/cloudpanel/htdocs/valentine-wall.gdgocuic.org/logs/error.log
tail -f /home/cloudpanel/htdocs/valentine-wall.gdgocuic.org/logs/access.log
```

### Restart Application

```bash
cd /home/cloudpanel/htdocs/valentine-wall
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml restart
```

### Update Application

```bash
cd /home/cloudpanel/htdocs/valentine-wall
git pull origin main
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml up -d --build
```

### Stop Application

```bash
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml down
```

---

## 14. Backup Strategy

### 14.1 Create Backup Script

```bash
nano /root/backup-valentine-wall.sh
```

Add:

```bash
#!/bin/bash
BACKUP_DIR="/root/backups/valentine-wall"
DATE=$(date +%Y%m%d_%H%M%S)
APP_DIR="/home/cloudpanel/htdocs/valentine-wall"

mkdir -p $BACKUP_DIR

# Backup PocketBase data
tar -czf $BACKUP_DIR/pb_data_$DATE.tar.gz -C $APP_DIR pb_data

# Backup configuration
cp $APP_DIR/.env $BACKUP_DIR/.env_$DATE

# Keep only last 7 backups
find $BACKUP_DIR -name "pb_data_*.tar.gz" -mtime +7 -delete
find $BACKUP_DIR -name ".env_*" -mtime +7 -delete

echo "Backup completed: $DATE"
```

Make executable:
```bash
chmod +x /root/backup-valentine-wall.sh
```

### 14.2 Schedule Daily Backups

```bash
# Add to crontab
crontab -e

# Add this line (runs daily at 3am)
0 3 * * * /root/backup-valentine-wall.sh >> /var/log/valentine-wall-backup.log 2>&1
```

---

## 15. Security Checklist

**CloudPanel handles most security, but verify:**

- [x] **Firewall:** Only ports 22, 80, 443, 8443 open
- [x] **SSL:** Let's Encrypt certificate active
- [x] **Nginx:** Reverse proxy configured securely
- [x] **Docker:** Containers only expose to localhost
- [x] **SSH:** Key-based authentication enabled
- [x] **Backups:** Automated daily backups configured
- [ ] **MongoDB:** Close port 27017 if exposed (security issue from CliniConnect)

### Close MongoDB Port (Security Fix)

```bash
# Check what's listening
netstat -tlnp | grep 27017

# If MongoDB is exposed, update firewall
ufw delete allow 27017/tcp
ufw reload
```

---

## 16. Troubleshooting

### Site Not Loading

```bash
# Check containers
docker ps

# Check Nginx
systemctl status nginx
nginx -t

# Check logs
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml logs
tail -f /var/log/nginx/error.log
```

### SSL Certificate Issues

```bash
# Renew certificate in CloudPanel
# Or manually:
certbot renew --dry-run
systemctl reload nginx
```

### Container Not Starting

```bash
# Check specific container
docker logs valentine-wall-backend-1
docker logs valentine-wall-frontend-1

# Rebuild
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml up -d --build --force-recreate
```

### 502 Bad Gateway

```bash
# Usually means containers aren't running or ports wrong
docker ps  # Should show both containers

# Check if ports are listening
netstat -tlnp | grep 8090  # Backend
netstat -tlnp | grep 3000  # Frontend

# Restart containers
docker-compose -f docker-compose.yml -f docker-compose.cloudpanel.yml restart
```

---

## 17. Summary

**You've successfully deployed Valentine Wall with CloudPanel! 🎉**

**Your setup:**
- ✅ **Domain:** valentine-wall.gdgocuic.org
- ✅ **SSL:** Let's Encrypt via CloudPanel
- ✅ **Reverse Proxy:** Nginx (CloudPanel)
- ✅ **Backend:** PocketBase on localhost:8090
- ✅ **Frontend:** Vite on localhost:3000
- ✅ **DNS:** Cloudflare with proxy enabled
- ✅ **Backups:** Automated daily backups

**Access Points:**
- Frontend: https://valentine-wall.gdgocuic.org
- Backend API: https://valentine-wall.gdgocuic.org/pb
- Admin Panel: https://valentine-wall.gdgocuic.org/pb/_/
- CloudPanel: https://194.233.73.243:8443

**Need help?** Check troubleshooting section or review logs.
