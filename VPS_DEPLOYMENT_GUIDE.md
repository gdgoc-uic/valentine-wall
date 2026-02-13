# 🚀 Valentine Wall - VPS Deployment Guide

**Simple, step-by-step guide for deploying Valentine Wall to Contabo VPS with Docker & Cloudflare**

---

## What You Need

- ✅ **Contabo VPS** (4GB RAM, 2 CPU, 50GB storage, Ubuntu 22.04)
- ✅ **Domain** registered and added to Cloudflare
- ✅ **SSH access** to your VPS
- ✅ **SMTP credentials** (Gmail, SendGrid, etc.) for emails

---

## Table of Contents

1. [VPS Initial Setup](#1-vps-initial-setup) (10 min)
2. [Install Docker](#2-install-docker) (5 min)
3. [Deploy Application](#3-deploy-application) (15 min)
4. [Configure DNS & SSL](#4-configure-dns--ssl) (10 min)
5. [Configure PocketBase](#5-configure-pocketbase) (10 min)
6. [Test Everything](#6-test-everything) (5 min)
7. [Basic Security](#7-basic-security) (10 min)
8. [Backups](#8-backups) (5 min)
9. [Troubleshooting](#9-troubleshooting)

**Total Time: ~1 hour**

---

## 1. VPS Initial Setup

SSH into your VPS and run these commands:

```bash
# SSH into VPS
ssh root@your-vps-ip

# Update system
apt update && apt upgrade -y

# Install essentials
apt install -y curl wget git nano ufw

# Configure firewall
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw enable

# Verify
ufw status
```

**✅ Test:** Run `ufw status` - should show rules active.

---

## 2. Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Start Docker
systemctl enable docker
systemctl start docker

# Create network for app
docker network create caddy

# Verify
docker ps
docker network ls
```

**✅ Test:** Run `docker --version` - should show Docker version 24+.

---

## 3. Deploy Application

### Setup SSH Key for GitHub (for Private Repos)

If your repository is private, set up SSH keys:

```bash
# Generate SSH key for GitHub
ssh-keygen -t ed25519 -C "your-email@example.com" -f ~/.ssh/github_deploy

# Start SSH agent
eval "$(ssh-agent -s)"

# Add key to agent
ssh-add ~/.ssh/github_deploy

# Display public key (copy this)
cat ~/.ssh/github_deploy.pub
```

**Add to GitHub:**
1. Copy the public key output
2. Go to GitHub → Settings → SSH and GPG keys
3. Click "New SSH key"
4. Paste key and save

**Configure SSH:**
```bash
# Create/edit SSH config
nano ~/.ssh/config
```

Add:
```
Host github.com
  HostName github.com
  User git
  IdentityFile ~/.ssh/github_deploy
  IdentitiesOnly yes
```

Save: `Ctrl+X`, `Y`, `Enter`

**✅ Test:** Run `IdentityFile ~/.ssh/valentine_wall` - should see "Hi username! You've successfully authenticated"

### Clone Repository

```bash
# Create app directory
mkdir -p ~/apps
cd ~/apps

# Clone via SSH (for private repos)
git clone git@github.com:yourusername/valentine-wall.git

# OR clone via HTTPS (for public repos)
# git clone https://github.com/yourusername/valentine-wall.git

cd valentine-wall
```

### Create Environment File

```bash
nano .env
```

Paste this (update values with your domain/credentials):

```env
# Environment
ENV=production

# Domain (replace with your subdomain)
BASE_DOMAIN=valentine.yourdomain.com
BACKEND_URL=https://valentine.yourdomain.com/pb
FRONTEND_URL=https://valentine.yourdomain.com

# App Settings
READ_ONLY=false
PROFANITY_JSON_FILE_NAME=profanities.json

# Firebase (optional - leave empty to skip OAuth)
FIREBASE_API_KEY=
FIREBASE_AUTH_DOMAIN=
FIREBASE_PROJECT_ID=
FIREBASE_STORAGE_BUCKET=
FIREBASE_MESSAGING_SENDER_ID=
FIREBASE_APP_ID=
FIREBASE_MEASUREMENT_ID=

# Report API (optional)
REPORT_API_URL=
REPORT_API_KEY=
REPORT_API_CATEGORY_ID_KEY=
```

Save: `Ctrl+X`, `Y`, `Enter`

### Deploy

```bash
# Make deploy script executable
chmod +x deploy.sh

# Deploy (takes 10-15 min first time)
./deploy.sh
```

**✅ Test:** Run `docker ps` - should see 4 containers running (backend, frontend, www, headless_chrome).

---

## 4. Configure DNS & SSL

### Add DNS Record in Cloudflare

1. Login to [Cloudflare](https://dash.cloudflare.com)
2. Select your domain
3. Go to **DNS** > **Records**
4. Click **Add record**:
   - Type: `A`
   - Name: `valentine` (or your subdomain)
   - IPv4 address: `your-vps-ip`
   - Proxy status: **Proxied** 🟧 (orange cloud - IMPORTANT!)
   - TTL: Auto
5. Click **Save**

### Configure SSL in Cloudflare

1. Go to **SSL/TLS** > **Overview**
2. Set mode to: **Full (strict)**
3. Go to **SSL/TLS** > **Edge Certificates**
4. Enable:
   - ✅ Always Use HTTPS
   - ✅ Automatic HTTPS Rewrites
   - ✅ Minimum TLS 1.2

**✅ Test:** Wait 2-5 minutes, then visit `https://valentine.yourdomain.com` - should load with HTTPS.

---

## 5. Configure PocketBase

### Access Admin Panel

Visit: `https://valentine.yourdomain.com/pb/_/`

**First time:** Create admin account (use a strong password!)

### Configure Email (SMTP)

1. Go to **Settings** > **Mail settings**
2. Configure SMTP:

**For Gmail:**
```
SMTP server: smtp.gmail.com
Port: 587
Username: your-email@gmail.com
Password: [App Password - see below]
Sender: your-email@gmail.com
```

**Gmail App Password:**
1. Enable 2FA on Google account
2. Visit: https://myaccount.google.com/apppasswords
3. Create app password for "Mail"
4. Use that 16-character password (not your regular password)

3. Click **Test email** - check your inbox
4. Click **Save**

### Add Initial Data (Optional)

**Gifts:**
- Collections > gifts > New record
- Add: `uid: rose, label: Red Rose, price: 50`

**Departments:**
- Collections > college_departments > New record
- Add your organization's departments

**✅ Test:** Send test email from PocketBase - should arrive in inbox.

---

## 6. Test Everything

### Quick Tests

```bash
# Containers running?
docker ps

# Frontend loads?
curl -I https://valentine.yourdomain.com

# Backend accessible?
curl -I https://valentine.yourdomain.com/pb/api/health
```

### Full Test (Browser)

1. **Visit:** `https://valentine.yourdomain.com`
2. **Register:** Create new account
3. **Verify:** Check email for verification link
4. **Login:** Login with verified account
5. **Send Message:** Send a message to another user
6. **Check Notification:** Recipient should get email

**✅ All working?** You're live! 🎉

---

## 8. Production Security Hardening

### Step 1: Create Deployment User

```bash
# Create dedicated deployment user
adduser deploy
usermod -aG sudo deploy
usermod -aG docker deploy

# Set strong password (use password generator)
```

### Step 2: Setup SSH Key Authentication

**On your local machine:**
```bash
# Generate SSH key pair (if you don't have one)
ssh-keygen -t ed25519 -C "your-email@example.com"

# Copy public key to server
ssh-copy-id deploy@your-vps-ip

# Test key-based login
ssh deploy@your-vps-ip
```

**On VPS (as root):**
```bash
# Copy root's authorized_keys to deploy user
mkdir -p /home/deploy/.ssh
cp ~/.ssh/authorized_keys /home/deploy/.ssh/
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh
chmod 600 /home/deploy/.ssh/authorized_keys
```

**✅ Test:** Open new terminal, `ssh deploy@your-vps-ip` should work without password

### Step 3: Harden SSH Configuration

```bash
# Backup SSH config
cp /etc/ssh/sshd_config /etc/ssh/sshd_config.backup

# Edit SSH config
nano /etc/ssh/sshd_config
```

**Configure these settings:**
```bash
# Disable root login
PermitRootLogin no

# Disable password authentication (use keys only)
PasswordAuthentication no
ChallengeResponseAuthentication no
PubkeyAuthentication yes

# Disable empty passwords
PermitEmptyPasswords no

# Change default SSH port (optional but recommended)
Port 2222

# Limit authentication attempts
MaxAuthTries 3
MaxSessions 2

# Disable X11 forwarding
X11Forwarding no

# Disable host-based authentication
HostbasedAuthentication no

# Set login grace time
LoginGraceTime 30s

# Allow only specific user
AllowUsers deploy
```

Save and test configuration:
```bash
# Test SSH config for errors
sshd -t

# If OK, restart SSH
systemctl restart sshd
```

**⚠️ IMPORTANT:** If you changed SSH port to 2222:
```bash
# Update firewall BEFORE closing current session
ufw allow 2222/tcp
ufw delete allow 22/tcp  # Remove old port after confirming new port works
```

**✅ Test:** Open NEW terminal, connect: `ssh -p 2222 deploy@your-vps-ip`

### Step 4: Install & Configure Fail2Ban

```bash
# Install Fail2Ban
apt install -y fail2ban

# Create local configuration
cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
nano /etc/fail2ban/jail.local
```

**Configure protection:**
```ini
[DEFAULT]
# Ban for 1 hour
bantime = 3600
# Check last 10 minutes
findtime = 600
# Max 3 retries
maxretry = 3
# Email alerts (optional)
destemail = your-email@example.com
sendername = Fail2Ban-VPS
action = %(action_mwl)s

[sshd]
enabled = true
port = 2222  # Change if you modified SSH port
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600

[docker-auth]
enabled = true
filter = docker-auth
logpath = /var/log/docker.log
maxretry = 3
bantime = 3600
```

Start Fail2Ban:
```bash
systemctl enable fail2ban
systemctl start fail2ban

# Check status
fail2ban-client status
fail2ban-client status sshd
```

**✅ Test:** Run `fail2ban-client status` - should show jails active

### Step 5: Configure Firewall (UFW) with Strict Rules

```bash
# Reset firewall to start fresh
ufw --force reset

# Default policies: deny all incoming, allow outgoing
ufw default deny incoming
ufw default allow outgoing

# Allow SSH (use your custom port if changed)
ufw allow 2222/tcp comment 'SSH'

# Allow HTTP/HTTPS only
ufw allow 80/tcp comment 'HTTP'
ufw allow 443/tcp comment 'HTTPS'

# Enable rate limiting on SSH
ufw limit 2222/tcp

# Enable firewall
ufw enable

# Check status
ufw status numbered
```

**✅ Test:** `ufw status` should show only necessary ports

### Step 6: Secure Docker

```bash
# Configure Docker daemon security
nano /etc/docker/daemon.json
```

Add:
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "live-restore": true,
  "userland-proxy": false,
  "no-new-privileges": true,
  "icc": false
}
```

```bash
# Restart Docker
systemctl restart docker
```

### Step 7: Secure Application Files

```bash
# Secure environment file
chmod 600 ~/apps/valentine-wall/.env
chown deploy:deploy ~/apps/valentine-wall/.env

# Secure data directories
chmod 700 ~/apps/valentine-wall/pb_data_fresh
chown -R deploy:deploy ~/apps/valentine-wall

# Verify permissions
ls -la ~/apps/valentine-wall/.env
ls -la ~/apps/valentine-wall/pb_data_fresh
```

### Step 8: Enable Automatic Security Updates

```bash
# Install unattended-upgrades
apt install -y unattended-upgrades apt-listchanges

# Configure automatic updates
dpkg-reconfigure -plow unattended-upgrades

# Edit config for security updates only
nano /etc/apt/apt.conf.d/50unattended-upgrades
```

Ensure these lines are uncommented:
```bash
Unattended-Upgrade::Allowed-Origins {
    "${distro_id}:${distro_codename}-security";
};

# Auto reboot if needed (optional)
Unattended-Upgrade::Automatic-Reboot "true";
Unattended-Upgrade::Automatic-Reboot-Time "03:00";
```

### Step 9: Install AIDE (Intrusion Detection)

```bash
# Install AIDE
apt install -y aide

# Initialize database (takes 5-10 min)
aideinit

# Move database
mv /var/lib/aide/aide.db.new /var/lib/aide/aide.db

# Schedule daily checks
echo '0 4 * * * root /usr/bin/aide --check | mail -s "AIDE Report" your-email@example.com' >> /etc/crontab
```

### Step 10: Setup Log Monitoring

```bash
# Install logwatch
apt install -y logwatch

# Configure daily reports
nano /etc/cron.daily/00logwatch
```

Add:
```bash
#!/bin/bash
/usr/sbin/logwatch --output mail --mailto your-email@example.com --detail high
```

```bash
chmod +x /etc/cron.daily/00logwatch
```

### Step 11: Harden PocketBase Admin

1. **Strong Admin Password:**
   - Minimum 20 characters
   - Mix of uppercase, lowercase, numbers, symbols
   - Use password manager

2. **IP Whitelist (Optional):**
```bash
# Add to UFW - allow admin access only from your IP
ufw allow from YOUR_IP_ADDRESS to any port 443 proto tcp comment 'Admin access'
```

3. **Regular Security Audits:**
   - Review user accounts weekly
   - Check PocketBase logs for suspicious activity
   - Monitor failed login attempts

### Step 12: Additional Hardening

```bash
# Disable unused services
systemctl list-unit-files --state=enabled
# Disable any unnecessary services

# Secure shared memory
echo "tmpfs /run/shm tmpfs defaults,noexec,nosuid 0 0" >> /etc/fstab

# Protect against IP spoofing
cat >> /etc/host.conf << EOF
order bind,hosts
nospoof on
EOF

# Kernel hardening
cat >> /etc/sysctl.conf << EOF
# IP Spoofing protection
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1

# Ignore ICMP ping requests
net.ipv4.icmp_echo_ignore_all = 1

# Ignore Directed pings
net.ipv4.icmp_echo_ignore_broadcasts = 1

# Disable source packet routing
net.ipv4.conf.all.accept_source_route = 0
net.ipv6.conf.all.accept_source_route = 0

# Ignore send redirects
net.ipv4.conf.all.send_redirects = 0

# Block SYN attacks
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.tcp_synack_retries = 2
net.ipv4.tcp_syn_retries = 5

# Log Martians
net.ipv4.conf.all.log_martians = 1
net.ipv4.icmp_ignore_bogus_error_responses = 1

# Ignore ICMP redirects
net.ipv4.conf.all.accept_redirects = 0
net.ipv6.conf.all.accept_redirects = 0

# Disable IPv6 (if not used)
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
EOF

# Apply kernel parameters
sysctl -p
```

### Step 13: Setup Security Monitoring Script

```bash
nano ~/security-check.sh
```

Add:
```bash
#!/bin/bash

echo "=== Security Status Check ==="
echo "Date: $(date)"
echo ""

# Check SSH attempts
echo "Failed SSH Attempts (last 24h):"
grep "Failed password" /var/log/auth.log | grep "$(date +%b\ %d)" | wc -l

# Check Fail2Ban status
echo ""
echo "Fail2Ban Jails:"
fail2ban-client status

# Check firewall
echo ""
echo "Firewall Rules:"
ufw status numbered | head -20

# Check running services
echo ""
echo "Listening Services:"
netstat -tulpn | grep LISTEN

# Check last logins
echo ""
echo "Recent Logins:"
lastlog | head -10

# Check container status
echo ""
echo "Docker Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}"

# Check disk usage
echo ""
echo "Disk Usage:"
df -h / | tail -1

# Check for updates
echo ""
echo "Security Updates Available:"
apt list --upgradable 2>/dev/null | grep -i security | wc -l
```

```bash
chmod +x ~/security-check.sh

# Run weekly
echo "0 9 * * 1 /root/security-check.sh | mail -s 'Weekly Security Report' your-email@example.com" >> /etc/crontab
```

**✅ Test:** Run `~/security-check.sh` - should display security status

### Security Checklist

After completing all steps:

- [ ] ✅ SSH key authentication enabled
- [ ] ✅ Password authentication disabled
- [ ] ✅ Root login disabled
- [ ] ✅ SSH port changed (optional)
- [ ] ✅ Fail2Ban active and configured
- [ ] ✅ UFW firewall configured with minimal rules
- [ ] ✅ Docker daemon secured
- [ ] ✅ Environment files protected (chmod 600)
- [ ] ✅ Automatic security updates enabled
- [ ] ✅ AIDE intrusion detection installed
- [ ] ✅ Log monitoring configured
- [ ] ✅ Kernel hardening applied
- [ ] ✅ PocketBase admin password strong (20+ chars)
- [ ] ✅ Security monitoring script scheduled
- [ ] ✅ Regular backups configured (see next section)

**Your VPS is now production-hardened! 🔒**

---

## 8. Backups

### Create Backup Script

```bash
nano ~/backup.sh
```

Paste:

```bash
#!/bin/bash
BACKUP_DIR=~/backups
APP_DIR=~/apps/valentine-wall
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup database
cp -r $APP_DIR/pb_data_fresh $BACKUP_DIR/pb_data_$DATE

# Backup files
tar -czf $BACKUP_DIR/valentine_$DATE.tar.gz -C $APP_DIR pb_data_fresh pb_public .env

# Keep last 7 days only
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "Backup completed: valentine_$DATE.tar.gz"
```

```bash
chmod +x ~/backup.sh

# Test
~/backup.sh

# Schedule daily backups (2 AM)
crontab -e
# Add: 0 2 * * * ~/backup.sh >> ~/backup.log 2>&1
```

**✅ Test:** Run `~/backup.sh` - should create backup in `~/backups/`.

---

## 10. Troubleshooting

### Containers Won't Start

```bash
# Check logs
docker logs backend
docker logs frontend

# Common fixes:
docker compose -f docker-compose.yml -f docker-compose.prod.yml down
docker network create caddy  # If needed
cd ~/apps/valentine-wall && ./deploy.sh
```

### Site Not Loading

```bash
# Check DNS
ping valentine.yourdomain.com  # Should resolve to your IP

# Check firewall
ufw status  # Ports 80, 443 should be allowed

# Check Cloudflare
# Ensure orange cloud is ON and SSL is "Full (strict)"

# Check Caddy
docker logs www
```

### SSL Certificate Issues

```bash
# Restart Caddy
docker restart www

# Wait 2-3 minutes for Let's Encrypt  
docker logs www -f  # Watch for cert generation

# Check Cloudflare SSL mode (should be "Full (strict)")
```

### Emails Not Sending

```bash
# Check backend logs
docker logs backend | grep -i mail

# Verify SMTP in PocketBase admin
# Test email should arrive in inbox

# For Gmail: ensure using App Password, not regular password
```

### Site Slow/Down

```bash
# Check resources
free -h      # Memory
df -h        # Disk
docker stats # Container usage

# Restart app
cd ~/apps/valentine-wall && ./deploy.sh
```

### Update Application

```bash
cd ~/apps/valentine-wall

# Pull latest code (SSH)
git pull

# OR if using HTTPS with credentials
# git pull origin main

# Redeploy
./deploy.sh

# Check logs
docker logs backend -f
```

### Rotate SSH Keys (Quarterly)

```bash
# Generate new GitHub deploy key
ssh-keygen -t ed25519 -C "your-email@example.com" -f ~/.ssh/github_deploy_new

# Add to GitHub (replace old key)
cat ~/.ssh/github_deploy_new.pub

# Update SSH config
sed -i 's/github_deploy/github_deploy_new/g' ~/.ssh/config

# Test
ssh -T git@github.com

# Remove old key
rm ~/.ssh/github_deploy ~/.ssh/github_deploy.pub
mv ~/.ssh/github_deploy_new ~/.ssh/github_deploy
```

---

## Quick Reference

```bash
# View logs
docker logs backend -f
docker logs frontend -f

# Restart app
cd ~/apps/valentine-wall && ./deploy.sh

# Stop app
docker compose -f docker-compose.yml -f docker-compose.prod.yml down

# Start app
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Check status
docker ps
docker stats

# Run backup
~/backup.sh

# Update system
apt update && apt upgrade -y
```

---

## Post-Deployment Checklist

- [ ] VPS accessible via SSH
- [ ] Docker installed and running
- [ ] Firewall configured (22, 80, 443 open)
- [ ] Application deployed (4 containers running)
- [ ] DNS pointing to VPS in Cloudflare
- [ ] SSL certificate working (HTTPS loads)
- [ ] PocketBase admin account created
- [ ] SMTP configured and tested
- [ ] User registration tested end-to-end
- [ ] Message sending tested
- [ ] Email notifications working
- [ ] Backups configured and tested
- [ ] Basic security implemented
- [ ] Environment file secured (chmod 600)

---

## Done! 🎉

Your Valentine Wall is now live at: **https://valentine.yourdomain.com**

**Maintenance:**
- Check logs weekly: `docker logs backend --tail 100`
- Update monthly: `git pull && ./deploy.sh`
- Test backups monthly: `~/backup.sh`
- Monitor disk space: `df -h`

**Need help?**
- Check logs: `docker logs backend`
- Review PocketBase docs: https://pocketbase.io/docs/
- Check container status: `docker ps`

---

**Guide Version:** 2.0 (Simplified)  
**Last Updated:** February 2026
