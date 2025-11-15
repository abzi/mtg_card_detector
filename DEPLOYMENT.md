# Deployment Guide

This guide covers deploying the MTG Card Detector application to production.

## Table of Contents

- [Backend Deployment](#backend-deployment)
- [Android APK Distribution](#android-apk-distribution)
- [Security Considerations](#security-considerations)
- [Monitoring](#monitoring)

---

## Backend Deployment

### Option 1: Docker (Recommended)

#### Prerequisites
- Docker and Docker Compose installed
- Domain name (optional but recommended)
- SSL certificate (Let's Encrypt recommended)

#### Steps

1. **Clone repository**:
   ```bash
   git clone https://github.com/abzi/mtg_card_detector.git
   cd mtg_card_detector
   ```

2. **Set environment variables**:
   ```bash
   export JWT_SECRET=$(openssl rand -base64 32)
   ```

3. **Build and run**:
   ```bash
   docker-compose up -d
   ```

4. **Verify**:
   ```bash
   curl http://localhost:8080/health
   ```

#### With Nginx Reverse Proxy

Create `nginx.conf`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable HTTPS with Certbot:
```bash
sudo certbot --nginx -d your-domain.com
```

### Option 2: Manual Deployment

#### On Ubuntu/Debian

1. **Install Go**:
   ```bash
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **Build application**:
   ```bash
   cd backend
   go build -o /usr/local/bin/mtg-server ./cmd/server
   ```

3. **Create systemd service** (`/etc/systemd/system/mtg-detector.service`):
   ```ini
   [Unit]
   Description=MTG Card Detector API
   After=network.target

   [Service]
   Type=simple
   User=www-data
   WorkingDirectory=/var/lib/mtg-detector
   ExecStart=/usr/local/bin/mtg-server
   Restart=always
   RestartSec=5

   Environment="PORT=8080"
   Environment="DATABASE_PATH=/var/lib/mtg-detector/data/mtg_cards.db"
   Environment="JWT_SECRET=CHANGE_THIS_TO_SECURE_SECRET"
   Environment="MIGRATIONS_PATH=/var/lib/mtg-detector/migrations"

   [Install]
   WantedBy=multi-user.target
   ```

4. **Setup directories**:
   ```bash
   sudo mkdir -p /var/lib/mtg-detector/{data,migrations}
   sudo cp backend/migrations/* /var/lib/mtg-detector/migrations/
   sudo chown -R www-data:www-data /var/lib/mtg-detector
   ```

5. **Start service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable mtg-detector
   sudo systemctl start mtg-detector
   sudo systemctl status mtg-detector
   ```

### Option 3: Cloud Platforms

#### Heroku

1. Create `Procfile`:
   ```
   web: cd backend && ./server
   ```

2. Deploy:
   ```bash
   heroku create your-app-name
   heroku config:set JWT_SECRET=$(openssl rand -base64 32)
   git push heroku main
   ```

#### DigitalOcean App Platform

1. Create app via web interface or doctl
2. Set environment variables in dashboard
3. Connect GitHub repository
4. Deploy automatically

#### AWS EC2

1. Launch EC2 instance (Ubuntu 22.04)
2. SSH and follow manual deployment steps
3. Configure security groups (port 80, 443)
4. Set up Elastic IP for static address

---

## Android APK Distribution

### Building Release APK

1. **Generate signing key** (one-time):
   ```bash
   keytool -genkey -v -keystore release.keystore \
     -alias mtg_detector -keyalg RSA -keysize 2048 \
     -validity 10000
   ```

2. **Configure signing** in `app/build.gradle`:
   ```gradle
   android {
       signingConfigs {
           release {
               storeFile file("release.keystore")
               storePassword "your_password"
               keyAlias "mtg_detector"
               keyPassword "your_password"
           }
       }
       buildTypes {
           release {
               signingConfig signingConfigs.release
           }
       }
   }
   ```

3. **Update API URL** for production:
   ```gradle
   buildConfigField "String", "API_BASE_URL", "\"https://your-api.com/api/v1\""
   ```

4. **Build**:
   ```bash
   cd android
   ./gradlew assembleRelease
   ```

5. **APK location**:
   ```
   android/app/build/outputs/apk/release/app-release.apk
   ```

### Distribution Options

#### 1. Direct Download
- Host APK on your website
- Users download and install manually
- Enable "Install from Unknown Sources"

#### 2. Google Play Store
- Create Developer account ($25 one-time fee)
- Create app listing
- Upload AAB (App Bundle):
  ```bash
  ./gradlew bundleRelease
  ```
- Submit for review

#### 3. F-Droid (Open Source)
- Submit to F-Droid repository
- Must be fully open source
- Automatic updates for users

---

## Security Considerations

### Backend

1. **JWT Secret**:
   - Use strong random secret (32+ characters)
   - Never commit to version control
   - Rotate periodically

2. **HTTPS**:
   - Always use HTTPS in production
   - Use Let's Encrypt for free certificates
   - Enable HSTS headers

3. **Database**:
   - Regular backups
   - Restrict file permissions
   - Consider encryption at rest

4. **Rate Limiting**:
   - Implement rate limiting middleware
   - Prevent abuse and DoS

5. **CORS**:
   - Restrict allowed origins in production
   - Don't use wildcard (*) in production

### Android

1. **API Keys**:
   - Never hardcode secrets in APK
   - Use BuildConfig or remote config

2. **Certificate Pinning** (advanced):
   - Pin SSL certificates
   - Prevent MITM attacks

3. **ProGuard**:
   - Enable code obfuscation
   - Reduce APK size

4. **Secure Storage**:
   - Already using EncryptedSharedPreferences
   - Don't log sensitive data

---

## Monitoring

### Backend Monitoring

#### 1. Application Logs

View logs:
```bash
# Docker
docker-compose logs -f backend

# Systemd
sudo journalctl -u mtg-detector -f
```

#### 2. Health Checks

```bash
curl https://your-api.com/health
```

#### 3. Metrics (Optional)

Integrate Prometheus + Grafana:
- Request rate
- Response times
- Error rates
- Database connections

### Database Backups

**Automated backup script** (`/usr/local/bin/backup-mtg-db.sh`):

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/mtg-detector"
DB_PATH="/var/lib/mtg-detector/data/mtg_cards.db"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR
sqlite3 $DB_PATH ".backup '$BACKUP_DIR/mtg_cards_$DATE.db'"

# Keep only last 7 days
find $BACKUP_DIR -name "*.db" -mtime +7 -delete
```

**Cron job**:
```bash
0 2 * * * /usr/local/bin/backup-mtg-db.sh
```

### Monitoring Services

- **Uptime monitoring**: UptimeRobot, Pingdom
- **Error tracking**: Sentry
- **APM**: New Relic, Datadog

---

## Scaling Considerations

### When to Scale

- Database grows beyond SQLite limits (>100GB)
- Concurrent users exceed single server capacity
- Request latency increases

### Scaling Options

1. **Vertical Scaling**:
   - Increase server resources (CPU, RAM)
   - Easy but has limits

2. **Horizontal Scaling**:
   - Multiple backend instances
   - Load balancer (nginx, HAProxy)
   - Shared database (PostgreSQL recommended)

3. **Database Migration**:
   - Migrate from SQLite to PostgreSQL
   - Minimal code changes needed
   - Better concurrent performance

---

## Troubleshooting

### Backend Won't Start

1. Check logs:
   ```bash
   docker-compose logs backend
   ```

2. Verify environment variables:
   ```bash
   docker-compose config
   ```

3. Check database permissions:
   ```bash
   ls -la /var/lib/mtg-detector/data/
   ```

### Android App Connection Issues

1. Verify API URL is correct
2. Check CORS settings on backend
3. Test with curl:
   ```bash
   curl -X POST https://your-api.com/api/v1/auth/anonymous \
     -H "Content-Type: application/json" \
     -d '{"device_id":"test"}'
   ```

### Database Locked Errors

- SQLite doesn't handle high concurrency well
- Consider migrating to PostgreSQL
- Reduce simultaneous requests

---

## Rollback Procedure

If deployment fails:

1. **Docker**:
   ```bash
   docker-compose down
   git checkout previous-commit
   docker-compose up -d
   ```

2. **Systemd**:
   ```bash
   sudo systemctl stop mtg-detector
   # Replace binary with previous version
   sudo systemctl start mtg-detector
   ```

3. **Database**:
   ```bash
   # Restore from backup
   cp /var/backups/mtg-detector/mtg_cards_YYYYMMDD.db \
      /var/lib/mtg-detector/data/mtg_cards.db
   ```

---

## Post-Deployment Checklist

- [ ] Backend health check responds
- [ ] Authentication endpoint works
- [ ] Scan endpoint works
- [ ] Inventory endpoint works
- [ ] HTTPS configured (production)
- [ ] Backups configured
- [ ] Monitoring set up
- [ ] Android APK tested with production API
- [ ] Documentation updated
- [ ] Users notified of update

---

For questions or issues, consult the main README or open an issue on GitHub.
