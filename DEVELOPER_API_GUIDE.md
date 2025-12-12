# Developer API Server v1.0.0
## GitHub Integration & Admin Protection

**Status**: ✅ Production Ready  
**Version**: 1.0.0  
**Port**: 5001  
**Security Level**: Enterprise Grade

---

## 📋 İçindekiler

1. [Genel Bakış](#genel-bakış)
2. [Kimlik Doğrulama](#kimlik-doğrulama)
3. [API Anahtarları](#api-anahtarları)
4. [OAuth2 Entegrasyonu](#oauth2-entegrasyonu)
5. [Webhook Yönetimi](#webhook-yönetimi)
6. [Admin Entegrasyonları](#admin-entegrasyonları)
7. [Güvenlik](#güvenlik)
8. [Hız Sınırlaması](#hız-sınırlaması)
9. [Örnekler](#örnekler)

---

## Genel Bakış

Developer API Server, GitHub, GitLab ve Bitbucket entegrasyonları için güvenli bir API sunucusudur. Admin panelinde kullanılan gelişmiş korumaları sağlar.

### Özellikleri
- ✅ OAuth2 Authentication (GitHub, GitLab, Bitbucket)
- ✅ API Key Management
- ✅ Webhook Subscriptions & Signing
- ✅ Admin Integration Management
- ✅ Rate Limiting (1000 req/hour)
- ✅ Audit Logging
- ✅ Token Management (7-day expiry)
- ✅ CORS Protection
- ✅ Timing Attack Protection

---

## Kimlik Doğrulama

### OAuth2 Flow

**1. Autorization Kodu Al**
```bash
# GitHub'a yönlendir
https://github.com/login/oauth/authorize?client_id=YOUR_CLIENT_ID&scope=admin:repo_hook%20repo:status
```

**2. Authorization Token Al**
```bash
curl -X POST http://localhost:5001/api/v2/oauth2/authorize \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "github",
    "code": "AUTHORIZATION_CODE"
  }'
```

**Response:**
```json
{
  "token": "abcdef1234567890...",
  "token_type": "Bearer",
  "expires_in": 604800,
  "scope": "admin:repo_hook repo:status repo:deployment public_repo"
}
```

**3. API Çağrılarında Kullan**
```bash
curl -X GET http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### API Key Authentication

Bazı endpoints API Key ile korunur:

```bash
curl -X GET "http://localhost:5001/api/v2/admin/integrations?key_id=dev_xyz&key_secret=secret_abc"

# Veya header ile
curl -X GET http://localhost:5001/api/v2/admin/integrations \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc"
```

---

## API Anahtarları

### Yeni API Anahtarı Oluştur

```bash
curl -X POST http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production Key",
    "permissions": ["read:user", "read:admin", "write:webhook"]
  }'
```

**Response:**
```json
{
  "success": true,
  "key_id": "dev_a1b2c3d4e5f6g7h8",
  "key_secret": "secret_1234567890abcdef1234567890abcdef",
  "message": "Keep your key secret safe!"
}
```

### API Anahtarlarını Listele

```bash
curl -X GET http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
  "keys": [
    {
      "key_id": "dev_a1b2c3d4e5f6g7h8",
      "name": "Production Key",
      "created_at": "2025-12-10T10:30:00",
      "last_used": "2025-12-10T15:45:23",
      "active": true,
      "rate_limit": 1000
    }
  ]
}
```

### API Anahtarını İptal Et

```bash
curl -X DELETE "http://localhost:5001/api/v2/developer/keys/dev_a1b2c3d4e5f6g7h8" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## OAuth2 Entegrasyonu

### Desteklenen Sağlayıcılar
- ✅ GitHub
- ✅ GitLab
- ✅ Bitbucket

### Token Oluştur

```bash
curl -X POST http://localhost:5001/api/v2/oauth2/authorize \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "github",
    "code": "gho_16C7e42F292c6912E7710..."
  }'
```

### Token İptal Et

```bash
curl -X POST http://localhost:5001/api/v2/oauth2/token/revoke \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Özellikler:**
- Token Expiry: 7 gün (configurable)
- Scope Support: Repo hooks, deployments, status
- Revocation: Anında
- Refresh: Manuel (yeni token oluştur)

---

## Webhook Yönetimi

### Webhook Oluştur

```bash
curl -X POST http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-app.com/webhook",
    "events": ["push", "pull_request", "issues"],
    "secret": "your_webhook_secret"
  }'
```

**Response:**
```json
{
  "webhook_id": "wh_1a2b3c4d5e6f7g8h",
  "secret": "whsec_1234567890abcdef...",
  "message": "Keep your secret safe!"
}
```

### Webhook'ları Listele

```bash
curl -X GET http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
  "webhooks": [
    {
      "webhook_id": "wh_1a2b3c4d5e6f7g8h",
      "url": "https://your-app.com/webhook",
      "events": ["push", "pull_request", "issues"],
      "active": true,
      "created_at": "2025-12-10T10:30:00",
      "last_triggered": "2025-12-10T15:45:23",
      "failures": 0
    }
  ]
}
```

### Webhook Payload İmzalamak

Server otomatik olarak HMAC-SHA256 ile imzalar:

```python
import hmac
import hashlib

# Serverda
payload = '{"push": true}'
signature = hmac.new(
    secret.encode(),
    payload.encode(),
    hashlib.sha256
).hexdigest()

# Header'da gönderilir
# X-Webhook-Signature: sha256=...
```

---

## Admin Entegrasyonları

### Entegrasyon Oluştur

```bash
curl -X POST http://localhost:5001/api/v2/admin/integrations \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc" \
  -H "Content-Type: application/json" \
  -d '{
    "admin_id": "admin_12345",
    "provider": "github",
    "config": {
      "repository": "myrepo/awesome-project",
      "branch": "main",
      "auto_deploy": true
    }
  }'
```

**Response:**
```json
{
  "integration_id": "int_a1b2c3d4e5f6g7h8",
  "provider": "github",
  "status": "created"
}
```

### Entegrasyonları Listele

```bash
# Belirli admin için
curl -X GET "http://localhost:5001/api/v2/admin/integrations?admin_id=admin_12345" \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc"

# Tüm entegrasyonlar
curl -X GET http://localhost:5001/api/v2/admin/integrations \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc"
```

---

## Güvenlik

### 1. Timing Attack Protection

```python
# Güvenli karşılaştırma
hash_equals($provided_key, $stored_key)
hmac.compare_digest(key1, key2)
```

### 2. HMAC Webhook Signing

Tüm webhooklar HMAC-SHA256 ile imzalanır:
```
X-Webhook-Signature: sha256=abcdef1234567890...
```

### 3. Token Expiry

- OAuth2 Tokens: 7 gün (default)
- Configurable via `DEV_CONFIG['token_expiry']`
- Automatic cleanup of expired tokens

### 4. Audit Logging

Tüm işlemler kaydedilir:
```json
{
  "timestamp": "2025-12-10T10:30:00",
  "action": "KEY_CREATED",
  "username": "developer@example.com",
  "ip": "192.168.1.100",
  "user_agent": "curl/7.68.0",
  "message": "API key created"
}
```

### 5. Rate Limiting

- Default: 1000 requests/hour per API key
- Per-key tracking
- HTTP 429 on exceed

### 6. IP Whitelisting (Optional)

```python
# AdminService'de
if !in_array($ip, $api_key['ip_whitelist']) {
    return forbidden();
}
```

### 7. Device Fingerprinting

Session'larda device fingerprint kontrol:
```python
fingerprint = hash('sha256', user_agent + accept_language + ip)
```

### 8. CORS Protection

```python
CORS(app)  # Configured in developer_api_server.py
```

---

## Hız Sınırlaması

### Yapılandırma

```python
DEV_CONFIG = {
    "enable_rate_limiting": True,
    "rate_limit_per_key": 1000,  # Requests per hour
    "rate_limit_window": 3600,   # seconds
}
```

### Headers

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 995
X-RateLimit-Reset: 1702276200
```

### Aşıldı Yanıtı

```bash
HTTP/1.1 429 Too Many Requests

{
  "error": "Rate limit exceeded",
  "remaining_requests": 0,
  "reset_at": "2025-12-10T11:30:00"
}
```

---

## Örnekler

### Senaryo 1: GitHub Entegrasyonu

**1. Admin GitHub hesabını bağla**
```bash
# Step 1: OAuth2 authorization code al (browser'dan)
# GitHub'a yönlendir:
https://github.com/login/oauth/authorize\
  ?client_id=Iv1.abc123...\
  &redirect_uri=http://localhost:5000/github-callback\
  &scope=admin:repo_hook,repo:status
```

**2. Token al**
```bash
curl -X POST http://localhost:5001/api/v2/oauth2/authorize \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "github",
    "code": "gho_16C7e42F292c6912E7710..."
  }'
```

**3. Admin entegrasyonu oluştur**
```bash
curl -X POST http://localhost:5001/api/v2/admin/integrations \
  -H "X-API-Key-ID: dev_prod_key" \
  -H "X-API-Key: secret_key_12345" \
  -H "Content-Type: application/json" \
  -d '{
    "admin_id": "admin_mertd",
    "provider": "github",
    "config": {
      "repository": "mert6148/User-login",
      "branch": "main",
      "auto_sync": true,
      "webhook_events": ["push", "pull_request"]
    }
  }'
```

### Senaryo 2: Webhook Subscription

**1. Webhook oluştur**
```bash
curl -X POST http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer YOUR_OAUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-app.com/github-webhook",
    "events": ["push", "pull_request", "release"],
    "secret": "your_webhook_secret_key"
  }'
```

**2. Webhook'u doğrula (sizin sunucunuzda)**
```python
import hmac
import hashlib

def verify_webhook(request):
    signature = request.headers.get('X-Webhook-Signature')
    payload = request.get_data()
    
    computed_sig = hmac.new(
        SECRET.encode(),
        payload,
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(signature.split('=')[1], computed_sig)
```

### Senaryo 3: API Key Yönetimi

**1. Yeni production key oluştur**
```bash
curl -X POST http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production - CI/CD",
    "permissions": [
      "read:user",
      "read:admin",
      "write:integrations"
    ]
  }'
```

**2. Key'i kullan**
```bash
# .env dosyasında sakla
export DEV_API_KEY_ID="dev_abc123..."
export DEV_API_KEY="secret_xyz789..."

# Requests'te kullan
curl -X GET http://localhost:5001/api/v2/developer/usage \
  -H "X-API-Key-ID: $DEV_API_KEY_ID" \
  -H "X-API-Key: $DEV_API_KEY"
```

**3. Key'i rotate et (eski key'i iptal et)**
```bash
# Yeni key oluştur
curl -X POST http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"name": "Production - CI/CD (Rotated)"}'

# Eski key'i iptal et
curl -X DELETE "http://localhost:5001/api/v2/developer/keys/dev_old_key" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## API Referans

### GET /api/v2
API bilgisi ve endpoints listesi

### POST /api/v2/oauth2/authorize
OAuth2 token oluştur

### POST /api/v2/oauth2/token/revoke
Token iptal et (Bearer required)

### POST /api/v2/developer/keys
Yeni API key oluştur (Bearer required)

### GET /api/v2/developer/keys
API keys'i listele (Bearer required)

### DELETE /api/v2/developer/keys/<key_id>
API key'i iptal et (Bearer required)

### POST /api/v2/developer/webhooks
Webhook oluştur (Bearer required)

### GET /api/v2/developer/webhooks
Webhook'ları listele (Bearer required)

### POST /api/v2/admin/integrations
Admin entegrasyonu oluştur (API Key required)

### GET /api/v2/admin/integrations
Admin entegrasyonlarını listele (API Key required)

### GET /api/v2/developer/usage
API kullanım istatistikleri (API Key required)

### GET /api/v2/health
Sağlık kontrol

---

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized (Missing credentials) |
| 403 | Forbidden (Invalid credentials) |
| 404 | Not Found |
| 429 | Rate Limit Exceeded |
| 500 | Internal Server Error |

---

## Database Şeması

### api_keys
```sql
CREATE TABLE api_keys (
    id INTEGER PRIMARY KEY,
    key_id TEXT UNIQUE,        -- dev_abc123...
    key_hash TEXT,             -- SHA256(secret)
    developer_id TEXT,         -- User ID
    name TEXT,                 -- "Production Key"
    created_at TIMESTAMP,
    last_used TIMESTAMP,
    rate_limit INTEGER,        -- 1000
    active BOOLEAN,            -- true/false
    permissions TEXT,          -- JSON array
    ip_whitelist TEXT          -- Optional
)
```

### oauth2_tokens
```sql
CREATE TABLE oauth2_tokens (
    id INTEGER PRIMARY KEY,
    token TEXT UNIQUE,
    provider TEXT,             -- github/gitlab/bitbucket
    user_id TEXT,
    scope TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP,
    revoked BOOLEAN
)
```

### webhooks
```sql
CREATE TABLE webhooks (
    id INTEGER PRIMARY KEY,
    webhook_id TEXT UNIQUE,
    developer_id TEXT,
    url TEXT,
    events TEXT,               -- JSON array
    secret TEXT,               -- HMAC secret
    active BOOLEAN,
    created_at TIMESTAMP,
    last_triggered TIMESTAMP,
    failures INTEGER
)
```

### admin_integrations
```sql
CREATE TABLE admin_integrations (
    id INTEGER PRIMARY KEY,
    integration_id TEXT UNIQUE,
    admin_id TEXT,
    provider TEXT,             -- github/gitlab
    config TEXT,               -- JSON config
    enabled BOOLEAN,
    created_at TIMESTAMP,
    last_sync TIMESTAMP
)
```

---

## Başlangıç

```bash
# 1. Dependencies yükle
pip install flask flask-cors

# 2. Server başlat
python developer_api_server.py

# 3. Health check yap
curl http://localhost:5001/api/v2/health

# 4. API info'yu al
curl http://localhost:5001/api/v2
```

---

## Gelişmiş Konfigürasyon

**developer_api_server.py içinde:**

```python
DEV_CONFIG = {
    "enable_logging": True,
    "enable_rate_limiting": True,
    "rate_limit_per_key": 1000,
    "rate_limit_window": 3600,
    "max_payload_size": 50 * 1024 * 1024,  # 50MB
    "allowed_github_events": [...],
    "token_expiry": 7 * 24 * 60 * 60,  # 7 days
    "webhook_timeout": 30,
    "enable_oauth2": True,
    "oauth2_providers": ["github", "gitlab", "bitbucket"]
}
```

---

## Loglama

Tüm işlemler `developer_api.log` ve `admin_audit.log` dosyalarına kaydedilir:

```json
{
  "timestamp": "2025-12-10 10:30:00",
  "level": "INFO",
  "message": "API key created: dev_abc123",
  "action": "KEY_CREATED",
  "developer_id": "user_123"
}
```

---

## Destek & Sorun Giderme

**Probleem:** "Invalid API key"
- Çözüm: key_id ve key_secret'u doğrula, whitelist'e IP'yi ekle

**Problem:** "Rate limit exceeded"
- Çözüm: 1 saat bekle veya rate limit'i artır

**Problem:** "Token expired"
- Çözüm: Yeni token oluştur (OAuth2 yeniden doğrula)

---

**Dokümantasyon**: v1.0  
**Son Güncelleme**: 10 Aralık 2025  
**Yazar**: GitHub Copilot
