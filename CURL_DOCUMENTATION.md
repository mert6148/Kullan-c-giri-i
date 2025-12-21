# API Server - cURL Komutları Rehberi

**Versiyon**: 2.0.0  
**Tarih**: 10 Aralık 2025  
**API URL**: `http://localhost:5000`  
**API Key**: `?key=12345`

---

## 📋 İçindekiler

1. [Health Check](#health-check)
2. [User Management](#user-management)
3. [Authentication](#authentication)
4. [Attributes](#attributes)
5. [Sessions](#sessions)
6. [Monitoring & Logs](#monitoring--logs)
7. [Dashboard](#dashboard)

---

## Health Check

### 🏥 Health Status Check

```bash
curl -X GET "http://localhost:5000/api/v1/health"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "status": "healthy",
  "api_version": "2.0.0",
  "timestamp": "2025-12-10T15:30:45.123456",
  "checks": {
    "database": true,
    "logging": true,
    "caching": true
  }
}
```

---

## User Management

### 📋 Get All Users

```bash
# GET /api/v1/users?key=12345
curl -X GET "http://localhost:5000/api/v1/users?key=12345"

# Alternative: with header
curl -X GET "http://localhost:5000/api/v1/users" \
  -H "X-API-Key: 12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "users": [
    {
      "id": 1,
      "username": "admin",
      "full_name": "System Administrator",
      "email": "admin@example.com"
    },
    {
      "id": 2,
      "username": "alice",
      "full_name": "Alice Smith",
      "email": "alice@example.com"
    }
  ]
}
```

### ➕ Create New User

```bash
curl -X POST "http://localhost:5000/api/v1/users?key=12345" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "password": "password123",
    "full_name": "Bob Johnson"
  }'
```

**Yanıt (201)**:
```json
{
  "success": true,
  "message": "User created",
  "username": "bob"
}
```

**Hataları**:
- `400`: Missing username or password
- `409`: User already exists
- `413`: Request too large
- `500`: Server error

### 👤 Get User by Username

```bash
curl -X GET "http://localhost:5000/api/v1/users/alice?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "user": {
    "id": 2,
    "username": "alice",
    "full_name": "Alice Smith",
    "email": "alice@example.com"
  }
}
```

**Hataları**:
- `404`: User not found
- `401`: Unauthorized (invalid key)

### 🗑️ Delete User

```bash
curl -X DELETE "http://localhost:5000/api/v1/users/bob?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "message": "User deleted"
}
```

---

## Authentication

### 🔐 Login

```bash
curl -X POST "http://localhost:5000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123"
  }'
```

**Yanıt (200)**:
```json
{
  "success": true,
  "message": "Login successful",
  "username": "alice",
  "session_id": "123e4567-e89b-12d3-a456-426614174000",
  "expires_at": "2025-12-11T15:30:45.123456"
}
```

**Hataları**:
- `400`: Missing credentials
- `401`: Invalid credentials
- `413`: Request too large

### 🚪 Logout

```bash
curl -X POST "http://localhost:5000/api/v1/auth/logout" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice"
  }'
```

**Yanıt (200)**:
```json
{
  "success": true,
  "message": "Logout successful"
}
```

---

## Attributes

### 📝 Get All User Attributes

```bash
curl -X GET "http://localhost:5000/api/v1/users/alice/attributes?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "username": "alice",
  "attributes": {
    "department": "IT",
    "role": "developer",
    "theme": "dark",
    "timezone": "Europe/Istanbul"
  }
}
```

### 🔍 Get Specific Attribute

```bash
curl -X GET "http://localhost:5000/api/v1/users/alice/attributes/department?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "username": "alice",
  "attribute_name": "department",
  "value": "IT"
}
```

**Hataları**:
- `404`: Attribute not found

### ✏️ Set User Attribute

```bash
curl -X POST "http://localhost:5000/api/v1/users/alice/attributes?key=12345" \
  -H "Content-Type: application/json" \
  -d '{
    "attribute_name": "department",
    "attribute_value": "Sales",
    "attribute_type": "string"
  }'
```

**Yanıt (200)**:
```json
{
  "success": true,
  "message": "Attribute set",
  "username": "alice",
  "attribute_name": "department"
}
```

**Desteklenen Tipler**:
- `string` (default)
- `integer`
- `boolean`
- `json`

### ❌ Delete Attribute

```bash
curl -X DELETE "http://localhost:5000/api/v1/users/alice/attributes/department?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "message": "Attribute deleted"
}
```

---

## Sessions

### 📊 Get All Sessions

```bash
curl -X GET "http://localhost:5000/api/v1/sessions?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "sessions": [
    {
      "id": "02a31bdf-aad7-4868-afd0-9f2d48ae2b47",
      "user_id": 2,
      "username": "alice",
      "status": "active",
      "login_ts": "2025-11-25T21:52:07Z"
    }
  ],
  "count": 1
}
```

### ➕ Create Session

```bash
curl -X POST "http://localhost:5000/api/v1/sessions?key=12345" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice"
  }'
```

**Yanıt (201)**:
```json
{
  "success": true,
  "session_id": "123e4567-e89b-12d3-a456-426614174000",
  "username": "alice",
  "expires_at": "2025-12-11T15:30:45.123456"
}
```

### 🔍 Get Session Details

```bash
curl -X GET "http://localhost:5000/api/v1/sessions/02a31bdf-aad7-4868-afd0-9f2d48ae2b47?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "session": {
    "id": "02a31bdf-aad7-4868-afd0-9f2d48ae2b47",
    "username": "alice",
    "status": "active",
    "login_ts": "2025-11-25T21:52:07Z"
  }
}
```

### 🛑 End Session

```bash
curl -X POST "http://localhost:5000/api/v1/sessions/02a31bdf-aad7-4868-afd0-9f2d48ae2b47?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "message": "Session ended"
}
```

---

## Monitoring & Logs

### 📋 Get API Logs

```bash
# Get last 50 logs (default)
curl -X GET "http://localhost:5000/api/v1/logs?key=12345"

# Get last 20 logs
curl -X GET "http://localhost:5000/api/v1/logs?key=12345&limit=20"

# Get last 100 logs
curl -X GET "http://localhost:5000/api/v1/logs?key=12345&limit=100"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "total": 50,
  "limit": 50,
  "logs": [
    {
      "timestamp": "2025-12-10T15:30:45.123456",
      "endpoint": "/api/v1/auth/login",
      "method": "POST",
      "status": 200,
      "ip": "127.0.0.1",
      "user_agent": "curl/7.68.0",
      "username": "alice",
      "message": "Login successful"
    },
    {
      "timestamp": "2025-12-10T15:30:40.654321",
      "endpoint": "/api/v1/users",
      "method": "GET",
      "status": 200,
      "ip": "127.0.0.1",
      "user_agent": "curl/7.68.0",
      "username": null,
      "message": "From cache"
    }
  ]
}
```

**Parametreler**:
- `limit`: 1-1000 (default: 50)

---

## Dashboard

### 📊 API Dashboard

```bash
curl -X GET "http://localhost:5000/api/v1/dashboard?key=12345"
```

**Yanıt (200)**:
```json
{
  "success": true,
  "timestamp": "2025-12-10T15:30:45.123456",
  "api_version": "2.0.0",
  "statistics": {
    "total_users": 4,
    "total_sessions": 2,
    "active_sessions": 1,
    "total_api_calls": 156,
    "successful_calls": 148,
    "failed_calls": 8,
    "cache_enabled": true,
    "logging_enabled": true
  },
  "configuration": {
    "rate_limit": 200,
    "rate_window": 60,
    "cache_ttl": 300,
    "max_request_size": 10485760
  }
}
```

---

## 🔐 Security Features

### API Key Authentication

Tüm `?key=12345` gerektiren endpoints'a API key gerekir:

```bash
# Method 1: Query parameter
curl "http://localhost:5000/api/v1/users?key=12345"

# Method 2: Header
curl "http://localhost:5000/api/v1/users" \
  -H "X-API-Key: 12345"
```

### Request Size Limits

Maximum 10MB request size. Daha büyük istekler 413 Payload Too Large döner:

```bash
curl -X POST "http://localhost:5000/api/v1/users?key=12345" \
  -H "Content-Type: application/json" \
  -d "$(head -c 11000000 /dev/zero)"  # Will fail with 413
```

### Rate Limiting

200 istekleri 60 saniyede izin verilir:

```bash
# Will work
for i in {1..200}; do
  curl -s "http://localhost:5000/api/v1/health"
done

# Might fail with 429 if rate limit exceeded
curl "http://localhost:5000/api/v1/health"
```

### Response Caching

GET istekleri 5 dakika (300 saniye) cache'lenebilir:

```bash
# İlk istek - database'den
curl "http://localhost:5000/api/v1/users?key=12345"
# Response: "from_cache": false

# Sonraki istekler - cache'ten
curl "http://localhost:5000/api/v1/users?key=12345"
# Response: "from_cache": true
```

---

## 🧪 Test Senaryoları

### Complete User Workflow

```bash
#!/bin/bash

API="http://localhost:5000"
KEY="?key=12345"

echo "1. Create user..."
curl -X POST "$API/api/v1/users$KEY" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"pass123","full_name":"Test User"}'

echo -e "\n2. Get user..."
curl -X GET "$API/api/v1/users/testuser$KEY"

echo -e "\n3. Set attribute..."
curl -X POST "$API/api/v1/users/testuser/attributes$KEY" \
  -H "Content-Type: application/json" \
  -d '{"attribute_name":"department","attribute_value":"IT"}'

echo -e "\n4. Get attribute..."
curl -X GET "$API/api/v1/users/testuser/attributes/department$KEY"

echo -e "\n5. Login..."
curl -X POST "$API/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"pass123"}'

echo -e "\n6. Get logs..."
curl -X GET "$API/api/v1/logs$KEY"

echo -e "\n7. Dashboard..."
curl -X GET "$API/api/v1/dashboard$KEY"

echo -e "\n8. Delete user..."
curl -X DELETE "$API/api/v1/users/testuser$KEY"

echo -e "\nDone!"
```

### Performance Test

```bash
#!/bin/bash

API="http://localhost:5000"
KEY="?key=12345"

echo "Testing response times..."

for endpoint in "health" "users" "dashboard"; do
  start=$(date +%s%N)
  curl -s "$API/api/v1/$endpoint$KEY" > /dev/null
  end=$(date +%s%N)
  ms=$((($end - $start) / 1000000))
  echo "$endpoint: ${ms}ms"
done
```

---

## 📚 Örnek Komut Kitapçığı

### Tüm kullanıcıları listele
```bash
curl "http://localhost:5000/api/v1/users?key=12345" | jq
```

### Alice'ın profil bilgilerini al
```bash
curl "http://localhost:5000/api/v1/users/alice?key=12345" | jq
```

### Alice'ı sistemde oluştur
```bash
curl -X POST "http://localhost:5000/api/v1/users?key=12345" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secure123","full_name":"Alice Smith"}' | jq
```

### Alice olarak giriş yap
```bash
curl -X POST "http://localhost:5000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secure123"}' | jq
```

### Alice'ın son 100 API çağrısını gör
```bash
curl "http://localhost:5000/api/v1/logs?key=12345&limit=100" | jq '.logs[] | select(.username=="alice")'
```

### Sistem sağlığını kontrol et
```bash
curl "http://localhost:5000/api/v1/health" | jq
```

### API Dashboard'u aç
```bash
curl "http://localhost:5000/api/v1/dashboard?key=12345" | jq
```

---

## ⚠️ Hata Kodları

| Kod | Anlam | Çözüm |
|-----|-------|-------|
| 200 | OK | İstek başarılı |
| 201 | Created | Kaynak oluşturuldu |
| 400 | Bad Request | Geçersiz istek formatı |
| 401 | Unauthorized | API key yok veya geçersiz |
| 404 | Not Found | Kaynak bulunamadı |
| 409 | Conflict | Kaynak zaten var (örn: user) |
| 413 | Payload Too Large | İstek çok büyük |
| 500 | Server Error | Sunucu hatası |

---

## 🔗 Linkler

- **API Root**: GET http://localhost:5000/
- **Health**: GET http://localhost:5000/api/v1/health
- **Users**: GET http://localhost:5000/api/v1/users?key=12345
- **Dashboard**: GET http://localhost:5000/api/v1/dashboard?key=12345
- **Logs**: GET http://localhost:5000/api/v1/logs?key=12345&limit=50

---

**Son Güncelleme**: 10 Aralık 2025  
**Sürüm**: 2.0.0  
**Durum**: ✅ Production Ready
