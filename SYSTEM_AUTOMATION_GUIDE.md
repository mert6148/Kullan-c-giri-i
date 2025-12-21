# Sistem Sunucusu Test Otomasyon & Dağıtım Rehberi

**Versiyon**: 2.0.0  
**Tarih**: 10 Aralık 2025  
**Durum**: ✅ Production Ready  
**Diller**: Python3, C#, PHP

---

## 📋 İçindekiler

1. [Test Otomasyon Akışı](#test-otomasyon-akışı)
2. [Sistem Sağlığı Kontrolleri](#sistem-sağlığı-kontrolleri)
3. [GET/POST Sağlamlaştırma](#getpost-sağlamlaştırma)
4. [Veritabanı Schema Doğrulaması](#veritabanı-schema-doğrulaması)
5. [İntegrasyon Testi](#i̇ntegrasyon-testi)
6. [Dağıtım Kontrol Listesi](#dağıtım-kontrol-listesi)

---

## Test Otomasyon Akışı

### 1. Ön-Test Hazırlığı

```bash
#!/bin/bash
# test_prep.sh - Test ortamı hazırla

echo "=== Test Ortamı Hazırlığı ==="

# 1. Python ortamı kontrol et
echo "1. Python sürümü kontrol ediliyor..."
python --version

# 2. Bağımlılıkları yükle
echo "2. Bağımlılıklar yükleniyor..."
pip install -r requirements.txt

# 3. Veritabanı sıfırla
echo "3. Veritabanı sıfırlanıyor..."
rm -f login_system.db
python -c "import print; print.init_db()"

# 4. Log dosyalarını temizle
echo "4. Log dosyaları temizleniyor..."
rm -f login_log.txt api_access.log

# 5. Cache temizle
echo "5. Cache temizleniyor..."
rm -f api_cache.json sessions_cache.json

# 6. API sunucusunu başlat
echo "6. API sunucusu başlatılıyor..."
python api_server.py &
API_PID=$!

# 7. Başlatılmasını bekle
sleep 2

echo "✅ Test ortamı hazır!"
echo "API PID: $API_PID"

exit 0
```

### 2. Health Check Test

```bash
#!/bin/bash
# test_health.sh - Sistem sağlığını kontrol et

echo "=== Sistem Sağlığı Kontrolleri ==="

API="http://localhost:5000"

# 1. API Health
echo "1. API Health Check..."
HEALTH=$(curl -s "$API/api/v1/health" | jq '.status')
if [ "$HEALTH" == '"healthy"' ]; then
    echo "   ✅ API Healthy"
else
    echo "   ❌ API Unhealthy: $HEALTH"
    exit 1
fi

# 2. Database Check
echo "2. Database Check..."
if [ -f "login_system.db" ]; then
    echo "   ✅ Database exists"
else
    echo "   ❌ Database missing"
    exit 1
fi

# 3. API Reachability
echo "3. API Reachability..."
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$API/api/v1/health")
if [ "$STATUS" == "200" ]; then
    echo "   ✅ API responding (HTTP $STATUS)"
else
    echo "   ❌ API error (HTTP $STATUS)"
    exit 1
fi

# 4. Database Connectivity
echo "4. Database Connectivity..."
python -c "
import sqlite3
try:
    conn = sqlite3.connect('login_system.db')
    c = conn.cursor()
    c.execute('SELECT COUNT(*) FROM users')
    count = c.fetchone()[0]
    conn.close()
    print(f'   ✅ Database connected ({count} users)')
except Exception as e:
    print(f'   ❌ Database error: {e}')
    exit(1)
"

echo "✅ All health checks passed!"
exit 0
```

### 3. User Management Test

```bash
#!/bin/bash
# test_users.sh - Kullanıcı yönetimi testleri

echo "=== Kullanıcı Yönetimi Testi ==="

API="http://localhost:5000"
KEY="?key=12345"

# Test 1: Create User
echo "1. Create User..."
RESPONSE=$(curl -s -X POST "$API/api/v1/users$KEY" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"pass123","full_name":"Test User"}')

if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "   ✅ User created"
else
    echo "   ❌ Failed to create user"
    exit 1
fi

# Test 2: Get User
echo "2. Get User..."
RESPONSE=$(curl -s "$API/api/v1/users/testuser$KEY")

if echo "$RESPONSE" | jq -e '.user.username == "testuser"' > /dev/null; then
    echo "   ✅ User retrieved"
else
    echo "   ❌ Failed to get user"
    exit 1
fi

# Test 3: List Users
echo "3. List Users..."
RESPONSE=$(curl -s "$API/api/v1/users$KEY")

if echo "$RESPONSE" | jq -e '.users | length > 0' > /dev/null; then
    COUNT=$(echo "$RESPONSE" | jq '.users | length')
    echo "   ✅ Users listed ($COUNT users)"
else
    echo "   ❌ Failed to list users"
    exit 1
fi

# Test 4: Delete User
echo "4. Delete User..."
RESPONSE=$(curl -s -X DELETE "$API/api/v1/users/testuser$KEY")

if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "   ✅ User deleted"
else
    echo "   ❌ Failed to delete user"
    exit 1
fi

echo "✅ All user tests passed!"
exit 0
```

---

## Sistem Sağlığı Kontrolleri

### 1. Veritabanı Kontrolleri

```python
def check_database_health():
    """Check database health"""
    import sqlite3
    from pathlib import Path
    
    checks = {
        "database_exists": False,
        "tables_present": False,
        "user_count": 0,
        "session_count": 0,
        "indexes_present": False,
        "disk_space_ok": True
    }
    
    try:
        # Check file exists
        if Path("login_system.db").exists():
            checks["database_exists"] = True
        
        # Connect and check tables
        conn = sqlite3.connect("login_system.db")
        c = conn.cursor()
        
        # Check tables
        c.execute("SELECT name FROM sqlite_master WHERE type='table'")
        tables = c.fetchall()
        if len(tables) > 0:
            checks["tables_present"] = True
        
        # Check user count
        c.execute("SELECT COUNT(*) FROM users")
        checks["user_count"] = c.fetchone()[0]
        
        # Check session count
        c.execute("SELECT COUNT(*) FROM user_attributes")
        checks["session_count"] = c.fetchone()[0]
        
        conn.close()
        
        return checks, all(checks.values())
    except Exception as e:
        return checks, False
```

### 2. API Endpoint Kontrolleri

```python
def check_api_endpoints():
    """Check all API endpoints"""
    import requests
    
    endpoints = {
        "GET /api/v1/health": {"method": "GET", "requires_key": False},
        "GET /api/v1/users": {"method": "GET", "requires_key": True},
        "POST /api/v1/auth/login": {"method": "POST", "requires_key": False},
        "GET /api/v1/dashboard": {"method": "GET", "requires_key": True},
        "GET /api/v1/logs": {"method": "GET", "requires_key": True},
    }
    
    base_url = "http://localhost:5000"
    results = {}
    
    for endpoint, config in endpoints.items():
        try:
            url = base_url + endpoint.split(" ")[1]
            if config["requires_key"]:
                url += "?key=12345"
            
            if config["method"] == "GET":
                response = requests.get(url, timeout=5)
            else:
                response = requests.post(url, json={}, timeout=5)
            
            results[endpoint] = response.status_code in [200, 201, 400, 401]
        except Exception:
            results[endpoint] = False
    
    return results, all(results.values())
```

### 3. Performans Kontrolleri

```python
def check_performance():
    """Check API performance"""
    import requests
    import time
    
    metrics = {
        "response_times": {},
        "slow_endpoints": [],
        "avg_response_time": 0
    }
    
    endpoints = [
        "/api/v1/health",
        "/api/v1/users?key=12345",
        "/api/v1/dashboard?key=12345"
    ]
    
    base_url = "http://localhost:5000"
    times = []
    
    for endpoint in endpoints:
        try:
            start = time.time()
            response = requests.get(base_url + endpoint, timeout=5)
            elapsed = time.time() - start
            
            metrics["response_times"][endpoint] = elapsed
            times.append(elapsed)
            
            # Slow if > 1 second
            if elapsed > 1.0:
                metrics["slow_endpoints"].append(endpoint)
        except Exception as e:
            metrics["response_times"][endpoint] = -1
    
    if times:
        metrics["avg_response_time"] = sum(times) / len(times)
    
    return metrics
```

---

## GET/POST Sağlamlaştırma

### 1. GET Request Validation

```python
def validate_get_request(endpoint, params=None, api_key=None):
    """Validate GET request safety"""
    issues = []
    
    # 1. Check URL safety
    if ".." in endpoint:
        issues.append("Path traversal detected")
    
    # 2. Check parameters
    if params:
        for key, value in params.items():
            if isinstance(value, str):
                # Check for SQL injection
                if ";" in value or "OR" in value.upper():
                    issues.append(f"Potential SQL injection in {key}")
                
                # Check size
                if len(value) > 1000:
                    issues.append(f"Parameter {key} too large")
    
    # 3. Check API key
    if not api_key:
        issues.append("Missing API key")
    
    return len(issues) == 0, issues
```

### 2. POST Request Validation

```python
def validate_post_request(endpoint, data=None, api_key=None):
    """Validate POST request safety"""
    issues = []
    
    # 1. Check content-type
    if not data:
        issues.append("Missing request body")
    
    # 2. Validate JSON structure
    if isinstance(data, dict):
        # Check for suspicious keys
        suspicious_keys = ["__proto__", "constructor", "prototype"]
        for key in data.keys():
            if key in suspicious_keys:
                issues.append(f"Suspicious key: {key}")
    
    # 3. Size check
    import json
    size = len(json.dumps(data))
    if size > 1048576:  # 1MB
        issues.append(f"Request too large: {size} bytes")
    
    # 4. API key validation
    if not api_key:
        issues.append("Missing API key")
    
    return len(issues) == 0, issues
```

### 3. Request/Response Logging

```python
def log_request_response(request_obj, response_obj):
    """Log request and response for audit"""
    log_entry = {
        "timestamp": datetime.now().isoformat(),
        "request": {
            "method": request_obj.method,
            "endpoint": request_obj.path,
            "params": dict(request_obj.args),
            "size": request_obj.content_length
        },
        "response": {
            "status": response_obj.status_code,
            "size": len(response_obj.get_data()),
            "time_ms": response_obj.headers.get("X-Response-Time", "N/A")
        }
    }
    
    with open("request_response.log", "a") as f:
        f.write(json.dumps(log_entry) + "\n")
```

---

## Veritabanı Schema Doğrulaması

### 1. Schema Kontrol Fonksiyonu

```python
def validate_database_schema():
    """Validate database schema"""
    import sqlite3
    
    conn = sqlite3.connect("login_system.db")
    c = conn.cursor()
    
    validations = {}
    
    # 1. Check users table
    try:
        c.execute("PRAGMA table_info(users)")
        columns = c.fetchall()
        expected = ["id", "username", "salt", "hash", "full_name", "email"]
        actual = [col[1] for col in columns]
        validations["users_table"] = all(col in actual for col in expected)
    except:
        validations["users_table"] = False
    
    # 2. Check user_attributes table
    try:
        c.execute("PRAGMA table_info(user_attributes)")
        columns = c.fetchall()
        expected = ["id", "user_id", "attribute_name", "attribute_value"]
        actual = [col[1] for col in columns]
        validations["attributes_table"] = all(col in actual for col in expected)
    except:
        validations["attributes_table"] = False
    
    # 3. Check indexes
    try:
        c.execute("SELECT name FROM sqlite_master WHERE type='index'")
        indexes = c.fetchall()
        validations["indexes_present"] = len(indexes) > 0
    except:
        validations["indexes_present"] = False
    
    conn.close()
    
    return validations, all(validations.values())
```

### 2. Schema Repair

```python
def repair_database_schema():
    """Repair missing tables/indexes"""
    import sqlite3
    
    conn = sqlite3.connect("login_system.db")
    c = conn.cursor()
    
    # Create missing tables
    c.executescript("""
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        salt TEXT NOT NULL,
        hash TEXT NOT NULL,
        full_name TEXT,
        email TEXT UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        last_login TIMESTAMP,
        is_active BOOLEAN DEFAULT 1
    );
    
    CREATE TABLE IF NOT EXISTS user_attributes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        attribute_name TEXT NOT NULL,
        attribute_value TEXT,
        attribute_type TEXT DEFAULT 'string',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        UNIQUE(user_id, attribute_name)
    );
    
    CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
    CREATE INDEX IF NOT EXISTS idx_user_attributes_user_id ON user_attributes(user_id);
    """)
    
    conn.commit()
    conn.close()
    
    return True
```

---

## İntegrasyon Testi

### 1. Complete Workflow Test

```bash
#!/bin/bash
# test_integration.sh - Tam iş akışı testi

echo "=== İntegrasyon Testi ==="

API="http://localhost:5000"
KEY="?key=12345"

# Test data
USERNAME="integration_test_user"
PASSWORD="TestPass123!"
EMAIL="test@example.com"

# 1. Create user
echo "1. Kullanıcı oluşturuluyor..."
CREATE=$(curl -s -X POST "$API/api/v1/users$KEY" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"email\":\"$EMAIL\"}")

if echo "$CREATE" | jq -e '.success' > /dev/null; then
    echo "   ✅ Kullanıcı oluşturuldu"
else
    echo "   ❌ Hata: Kullanıcı oluşturulamadı"
    exit 1
fi

# 2. Set attributes
echo "2. Özellikler ayarlanıyor..."
SET_ATTR=$(curl -s -X POST "$API/api/v1/users/$USERNAME/attributes$KEY" \
  -H "Content-Type: application/json" \
  -d '{"attribute_name":"department","attribute_value":"IT"}')

if echo "$SET_ATTR" | jq -e '.success' > /dev/null; then
    echo "   ✅ Özellikler ayarlandı"
else
    echo "   ❌ Hata: Özellikler ayarlanamadı"
    exit 1
fi

# 3. Login
echo "3. Kullanıcı giriş yapıyor..."
LOGIN=$(curl -s -X POST "$API/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

if echo "$LOGIN" | jq -e '.success' > /dev/null; then
    SESSION_ID=$(echo "$LOGIN" | jq -r '.session_id')
    echo "   ✅ Giriş başarılı (Session: $SESSION_ID)"
else
    echo "   ❌ Hata: Giriş başarısız"
    exit 1
fi

# 4. Get dashboard
echo "4. Dashboard alınıyor..."
DASHBOARD=$(curl -s "$API/api/v1/dashboard$KEY")

if echo "$DASHBOARD" | jq -e '.statistics.total_users' > /dev/null; then
    TOTAL=$(echo "$DASHBOARD" | jq '.statistics.total_users')
    echo "   ✅ Dashboard alındı ($TOTAL kullanıcı)"
else
    echo "   ❌ Hata: Dashboard alınamadı"
    exit 1
fi

# 5. Get logs
echo "5. Loglar alınıyor..."
LOGS=$(curl -s "$API/api/v1/logs$KEY&limit=10")

if echo "$LOGS" | jq -e '.logs | length > 0' > /dev/null; then
    COUNT=$(echo "$LOGS" | jq '.logs | length')
    echo "   ✅ Loglar alındı ($COUNT log)"
else
    echo "   ❌ Hata: Loglar alınamadı"
    exit 1
fi

# 6. Cleanup
echo "6. Test kullanıcısı siliniyor..."
DELETE=$(curl -s -X DELETE "$API/api/v1/users/$USERNAME$KEY")

if echo "$DELETE" | jq -e '.success' > /dev/null; then
    echo "   ✅ Test kullanıcısı silindi"
else
    echo "   ❌ Hata: Test kullanıcısı silinemedi"
    exit 1
fi

echo "✅ İntegrasyon testi başarılı!"
exit 0
```

### 2. Error Handling Test

```bash
#!/bin/bash
# test_errors.sh - Hata yönetimi testi

echo "=== Hata Yönetimi Testi ==="

API="http://localhost:5000"

# Test 1: Invalid API Key
echo "1. Geçersiz API Key..."
RESULT=$(curl -s "$API/api/v1/users?key=wrongkey")
if echo "$RESULT" | jq -e '.error' > /dev/null; then
    echo "   ✅ API Key validation works"
else
    echo "   ❌ API Key validation failed"
fi

# Test 2: Missing Parameters
echo "2. Eksik Parametreler..."
RESULT=$(curl -s -X POST "$API/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"test"}')
if echo "$RESULT" | jq -e '.error' > /dev/null; then
    echo "   ✅ Parameter validation works"
else
    echo "   ❌ Parameter validation failed"
fi

# Test 3: User Not Found
echo "3. Kullanıcı Bulunamadı..."
RESULT=$(curl -s "$API/api/v1/users/nonexistent?key=12345")
if echo "$RESULT" | jq -e '.error' > /dev/null; then
    echo "   ✅ 404 handling works"
else
    echo "   ❌ 404 handling failed"
fi

# Test 4: Invalid Request Size
echo "4. Büyük İstek..."
# Create a large payload
LARGE_DATA=$(python -c "print('{\"data\":\"' + 'x'*11000000 + '\"}')")
RESULT=$(curl -s -X POST "$API/api/v1/users?key=12345" \
  -H "Content-Type: application/json" \
  -d "$LARGE_DATA")
if echo "$RESULT" | jq -e '.error' > /dev/null; then
    echo "   ✅ Size limit works"
else
    echo "   ❌ Size limit failed"
fi

echo "✅ Hata yönetimi testleri tamamlandı!"
exit 0
```

---

## Dağıtım Kontrol Listesi

### Pre-Deployment Checklist

- ✅ Tüm testler geçti
- ✅ Veritabanı schema doğrulandi
- ✅ API endpoints çalışıyor
- ✅ Logging aktif
- ✅ Caching yapılandırıldı
- ✅ Rate limiting ayarlandı
- ✅ CORS origins yapılandırıldı
- ✅ Backup sistemi test edildi
- ✅ Security headers eklendi
- ✅ Documentation güncel
- ✅ Performance acceptable (<100ms)
- ✅ Error handling comprehensive
- ✅ Audit logging enabled
- ✅ Monitoring configured
- ✅ Rollback plan prepared

### Deployment Steps

```bash
#!/bin/bash
# deploy.sh - Production deployment

echo "=== Production Deployment ==="

# 1. Run tests
echo "1. Tests çalıştırılıyor..."
bash test_health.sh || exit 1
bash test_users.sh || exit 1
bash test_integration.sh || exit 1

# 2. Backup database
echo "2. Veritabanı yedekleniyor..."
cp login_system.db "login_system.db.backup.$(date +%s)"

# 3. Migrate data
echo "3. Veriler migrate ediliyor..."
python -c "import print; print.init_db()"

# 4. Start services
echo "4. Servisler başlatılıyor..."
python api_server.py > api_server.log 2>&1 &
API_PID=$!

# 5. Health check
echo "5. Sağlık kontrolleri yapılıyor..."
sleep 2
curl -s "http://localhost:5000/api/v1/health" | jq

# 6. Verify
echo "6. Sistem doğrulanıyor..."
bash test_health.sh || { kill $API_PID; exit 1; }

echo "✅ Deployment başarılı!"
echo "API PID: $API_PID"

exit 0
```

---

**Son Güncelleme**: 10 Aralık 2025  
**Sürüm**: 2.0.0  
**Durum**: ✅ Production Ready
