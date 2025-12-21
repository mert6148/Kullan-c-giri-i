# ADMIN & DEVELOPER API - Hızlı Referans
## Quick Start & Command Reference

---

## 🚀 İlk Kurulum (5 dakika)

### 1. Python Developer API Sunucusu

```bash
# Başlat
python developer_api_server.py

# Test et
curl http://localhost:5001/api/v2/health

# Output:
# {"status": "ok", "version": "1.0.0", ...}
```

### 2. Admin Panel

```bash
# Oturum aç
http://localhost/admin/login

# Test Kullanıcı
Username: admin
Password: admin123
Role: admin
```

### 3. Dosya Kontrol

```bash
# Oluşturulan dosyalar
ls -la src/Service/AdminService.php      # ✅ RBAC
ls -la src/Controller/AdminController.php # ✅ Routes
ls -la developer_api_server.py             # ✅ API Server
```

---

## 🔐 Admin Login & Permissions

### Login Yap
```bash
curl -X POST http://localhost/admin/login \
  -d "username=admin&password=admin123&role=admin"
```

### Session Info Al
```bash
curl -X GET http://localhost/admin/session/info \
  -H "Cookie: PHPSESSID=abc123..."
```

### Session Uzat (30 dk timeout)
```bash
curl -X POST http://localhost/admin/session/extend \
  -H "Cookie: PHPSESSID=abc123..."
```

### İzinleri Kontrol Et
```bash
curl -X GET http://localhost/admin/permissions \
  -H "Cookie: PHPSESSID=abc123..."

# Response:
# {"role": "admin", "permissions": [...]}
```

---

## 🔑 Developer API - OAuth2

### 1. Authorization Code Al (Browser)
```
https://github.com/login/oauth/authorize?
  client_id=YOUR_CLIENT_ID&
  scope=admin:repo_hook,repo:status&
  redirect_uri=http://localhost/callback
```

### 2. Token Oluştur
```bash
curl -X POST http://localhost:5001/api/v2/oauth2/authorize \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "github",
    "code": "gho_16C7e42F292c6912E7710..."
  }'
```

### 3. Token Kullan
```bash
curl -X GET http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Token İptal Et
```bash
curl -X POST http://localhost:5001/api/v2/oauth2/token/revoke \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🔐 API Key Yönetimi

### Yeni Key Oluştur
```bash
curl -X POST http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "CI/CD Key",
    "permissions": ["read:user", "read:admin"]
  }'

# Response:
# {
#   "key_id": "dev_a1b2c3d4...",
#   "key_secret": "secret_xyz789...",
# }
```

### Keys Listele
```bash
curl -X GET http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer TOKEN"
```

### Key İptal Et
```bash
curl -X DELETE "http://localhost:5001/api/v2/developer/keys/dev_a1b2c3d4" \
  -H "Authorization: Bearer TOKEN"
```

### API Kullanım İstatistikleri
```bash
curl -X GET http://localhost:5001/api/v2/developer/usage \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc"
```

---

## 🪝 Webhook Yönetimi

### Webhook Oluştur
```bash
curl -X POST http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-app.com/github-webhook",
    "events": ["push", "pull_request", "issues"],
    "secret": "your_webhook_secret"
  }'
```

### Webhooks Listele
```bash
curl -X GET http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer TOKEN"
```

### Webhook Doğrulamak (Sunucunuzda)
```python
import hmac
import hashlib

def verify_webhook(payload, signature, secret):
    expected_sig = hmac.new(
        secret.encode(),
        payload,
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(signature, expected_sig)
```

---

## 🔗 Admin Entegrasyonları

### GitHub Entegrasyonu Oluştur
```bash
curl -X POST http://localhost:5001/api/v2/admin/integrations \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc" \
  -H "Content-Type: application/json" \
  -d '{
    "admin_id": "admin_mertd",
    "provider": "github",
    "config": {
      "repository": "mert6148/User-login",
      "branch": "main",
      "auto_sync": true
    }
  }'
```

### Entegrasyonları Listele
```bash
curl -X GET "http://localhost:5001/api/v2/admin/integrations?admin_id=admin_mertd" \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc"
```

---

## 📊 Admin Routes

### Kullanıcılar (requires: user:manage)
```bash
# Listele
curl -X GET http://localhost/admin/users \
  -H "Cookie: PHPSESSID=..."

# Düzenle
curl -X GET http://localhost/admin/users/123/edit \
  -H "Cookie: PHPSESSID=..."

# Sil
curl -X DELETE http://localhost/admin/users/123 \
  -H "Cookie: PHPSESSID=..."
```

### Loglar (requires: logs:view)
```bash
# Listele
curl -X GET "http://localhost/admin/logs?limit=50" \
  -H "Cookie: PHPSESSID=..."

# CSV Dışa Aktar
curl -X POST http://localhost/admin/logs/export \
  -d "limit=1000" \
  -H "Cookie: PHPSESSID=..." > audit.csv
```

### Güvenlik (requires: security:manage)
```bash
curl -X GET http://localhost/admin/security \
  -H "Cookie: PHPSESSID=..."
```

### Veritabanı (requires: database:manage)
```bash
curl -X GET http://localhost/admin/database \
  -H "Cookie: PHPSESSID=..."

# Yedek al
curl -X POST http://localhost/admin/database/backup \
  -H "Cookie: PHPSESSID=..."
```

---

## 🛡️ Güvenlik Kontrolleri

### Cihaz Doğrula
```bash
curl -X POST http://localhost/admin/device/verify \
  -H "Cookie: PHPSESSID=..."

# Response:
# {"success": true, "message": "Cihaz doğrulandı"}
```

### API Key Doğrula
```bash
# Yanlış key
curl -X GET "http://localhost:5001/api/v2/developer/usage?key_id=xyz&key_secret=wrong" \
  
# Response: {"error": "Invalid API key"}, 403

# Doğru key
curl -X GET "http://localhost:5001/api/v2/developer/usage" \
  -H "X-API-Key-ID: dev_xyz" \
  -H "X-API-Key: secret_abc"

# Response: {"requests_this_hour": 25, ...}
```

---

## 📈 API Hata Kodları

| Kod | Açıklama | Çözüm |
|-----|----------|-------|
| 200 | Success | ✅ |
| 201 | Created | ✅ |
| 400 | Bad Request | JSON'u kontrol et |
| 401 | Unauthorized | Token/Key eksik |
| 403 | Forbidden | Token/Key yanlış |
| 404 | Not Found | URL kontrol et |
| 429 | Rate Limited | 1 saat bekle |
| 500 | Server Error | Logs kontrol et |

---

## 🔄 Workflow Örnekleri

### Örnek 1: Admin Login → Dashboard → Logları Görüntüle

```bash
# 1. Login
curl -c cookies.txt -X POST http://localhost/admin/login \
  -d "username=admin&password=admin123&role=admin"

# 2. Dashboard
curl -b cookies.txt -X GET http://localhost/admin

# 3. Logları görüntüle (logs:view izni gerekli)
curl -b cookies.txt -X GET "http://localhost/admin/logs?limit=20"
```

### Örnek 2: OAuth2 → API Key → Entegrasyon

```bash
# 1. OAuth2 token al
TOKEN=$(curl -s -X POST http://localhost:5001/api/v2/oauth2/authorize \
  -d '{"provider": "github", "code": "gho_xyz..."}' \
  | jq -r .token)

# 2. API key oluştur
KEY=$(curl -s -X POST http://localhost:5001/api/v2/developer/keys \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Prod"}' \
  | jq -r .key_id)

# 3. Admin entegrasyonu oluştur
curl -X POST http://localhost:5001/api/v2/admin/integrations \
  -H "X-API-Key-ID: $KEY" \
  -H "X-API-Key: ..." \
  -d '{
    "admin_id": "admin_mertd",
    "provider": "github",
    "config": {"repository": "..."}
  }'
```

### Örnek 3: Webhook Oluştur → Doğrula → Test

```bash
# 1. Webhook oluştur
WH=$(curl -s -X POST http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url": "...", "events": ["push"]}' \
  | jq -r .webhook_id)

# 2. Test payload gönder
curl -X POST http://your-app.com/webhook \
  -H "X-Webhook-Signature: sha256=..." \
  -d '{"event": "push", "...": "..."}'

# 3. Başarı kontrol
curl -X GET http://localhost:5001/api/v2/developer/webhooks \
  -H "Authorization: Bearer $TOKEN" | jq '.webhooks[] | select(.webhook_id == "'$WH'")'
```

---

## ⚙️ Konfigürasyon İpuçları

### Session Timeout'u Değiştir
```php
// src/Service/AdminService.php
const SESSION_TIMEOUT = 60;     // 60 dakika
const SESSION_WARNING = 55;     // 55 dakikada uyar
```

### Rate Limit'i Değiştir
```python
# developer_api_server.py
DEV_CONFIG['rate_limit_per_key'] = 5000  # Artır
DEV_CONFIG['rate_limit_window'] = 3600   # 1 saat
```

### Token Expiry'yi Değiştir
```python
DEV_CONFIG['token_expiry'] = 30 * 24 * 60 * 60  # 30 gün
```

---

## 🐛 Debugging

### Admin Audit Logları
```bash
tail -f admin_audit.log

# Örnek çıktı:
# {"timestamp":"2025-12-10 10:30:00","action":"LOGIN_SUCCESS","username":"admin",...}
```

### Developer API Logları
```bash
tail -f developer_api.log

# Örnek çıktı:
# 2025-12-10 10:30:00 - INFO - API key created: dev_abc123...
```

### Database Sorgula
```bash
sqlite3 developer_api.db "SELECT * FROM api_keys LIMIT 5;"
```

---

## 📋 Checklist

- [ ] AdminService.php aktif
- [ ] AdminController.php configure edildi
- [ ] developer_api_server.py çalışıyor (port 5001)
- [ ] OAuth2 providers eklenmiş
- [ ] Admin user oluşturulmuş
- [ ] İlk API key oluşturulmuş
- [ ] Webhook endpoint kurulmuş
- [ ] Audit logs kontrol edildi
- [ ] Rate limiting tested
- [ ] CORS configured

---

## 📞 Support

**Dokümantasyon:**
- `ADMIN_PROTECTION_GUIDE.md` - Detaylı admin rehberi
- `DEVELOPER_API_GUIDE.md` - Developer API detayları

**Log Dosyaları:**
- `admin_audit.log` - Admin operasyonları
- `developer_api.log` - API operasyonları

**Database:**
- `developer_api.db` - SQLite (api_keys, webhooks, vb.)

---

**Version**: 2.0.0  
**Last Updated**: 10 Aralık 2025  
**Status**: Production Ready ✅
