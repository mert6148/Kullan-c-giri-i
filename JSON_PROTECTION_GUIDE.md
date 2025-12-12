# JSON Koruma Sistemi - Kapsamlı Rehber

**Versiyon**: 2.0.0  
**Tarih**: 10 Aralık 2025  
**Durum**: ✅ Production Ready

---

## 📋 İçindekiler

1. [JSON Veri Yapısı & Koruma](#json-veri-yapısı--koruma)
2. [Request Validation](#request-validation)
3. [Response Formatting](#response-formatting)
4. [Error Handling](#error-handling)
5. [Security Measures](#security-measures)
6. [Best Practices](#best-practices)

---

## JSON Veri Yapısı & Koruma

### 1. User JSON Format

```json
{
  "id": 1,
  "username": "alice",
  "full_name": "Alice Smith",
  "email": "alice@example.com",
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-25T21:52:07Z",
  "is_active": true,
  "is_admin": false
}
```

**Koruma Seviyeleri**:
- 🟢 **Public**: username, full_name
- 🟡 **Private**: email, created_at
- 🔴 **Secret**: password hash, salt (hiçbir zaman JSON'da döndürülmez)

### 2. Session JSON Format

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": 1,
  "username": "alice",
  "status": "active",
  "login_ts": "2025-11-25T21:52:07Z",
  "logout_ts": null,
  "duration_seconds": 3600,
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "created_at": "2025-11-25T21:52:07Z",
  "updated_at": "2025-11-25T21:52:07Z"
}
```

**Gizli Alanlar**: Asla JSON'da döndürülmez
- Şifre hash
- Şifre salt
- 2FA secret

### 3. Attribute JSON Format

```json
{
  "id": 1,
  "user_id": 1,
  "attribute_name": "department",
  "attribute_value": "IT",
  "attribute_type": "string",
  "category": "profile",
  "created_at": "2025-11-25T10:00:00Z",
  "updated_at": "2025-11-25T10:00:00Z"
}
```

**Koruma Kuralları**:
- Protected attributes sadece admin tarafından görülebilir
- Sensitive attributes (SSN, etc) encryption ile saklanır

### 4. Log JSON Format

```json
{
  "timestamp": "2025-12-10T15:30:45.123456",
  "endpoint": "/api/v1/auth/login",
  "method": "POST",
  "status": 200,
  "ip": "127.0.0.1",
  "user_agent": "curl/7.68.0",
  "username": "alice",
  "message": "Login successful"
}
```

**Gizli Alanlar**:
- Request body'deki passwords
- Authentication headers
- Sensitive query parameters

---

## Request Validation

### 1. Input Sanitization

```python
def sanitize_json_input(data: dict) -> dict:
    """Sanitize input JSON data"""
    if not isinstance(data, dict):
        return {}
    
    # Remove potential XSS vectors
    for key, value in data.items():
        if isinstance(value, str):
            # Remove HTML tags
            value = value.replace("<", "&lt;").replace(">", "&gt;")
            # Limit length
            value = value[:1000]
        data[key] = value
    
    return data
```

**Kontrol Edilen Alanlar**:
- ✅ String length limits (max 1000 chars)
- ✅ No HTML/script tags
- ✅ SQL injection protection (prepared statements)
- ✅ No executable code

### 2. Type Validation

```python
def validate_attribute_type(attr_type: str) -> bool:
    """Validate attribute type"""
    valid_types = ["string", "integer", "boolean", "json", "binary", "file"]
    return attr_type in valid_types

def validate_json_type(value, expected_type: str) -> bool:
    """Validate value matches expected JSON type"""
    type_map = {
        "string": str,
        "integer": int,
        "boolean": bool,
        "json": dict,
        "array": list
    }
    return isinstance(value, type_map.get(expected_type, str))
```

### 3. Size Validation

```python
def validate_request_size(request_size: int) -> bool:
    """Check if request size is within limit"""
    MAX_SIZE = 10485760  # 10MB
    return request_size <= MAX_SIZE

# API'de
if request.content_length and request.content_length > MAX_SIZE:
    return jsonify({"error": "Request too large"}), 413
```

**Limit Kuralları**:
- Max request size: 10MB
- Max field value: 1000 chars
- Max array items: 10000 items

### 4. Required Field Validation

```python
def validate_required_fields(data: dict, required: list) -> bool:
    """Validate all required fields are present"""
    for field in required:
        if field not in data or data[field] is None:
            return False
    return True

# Kullanım
if not validate_required_fields(data, ["username", "password"]):
    return jsonify({"error": "Missing required fields"}), 400
```

---

## Response Formatting

### 1. Success Response Format

```python
def format_success_response(data: dict, message: str = None, status_code: int = 200):
    """Format success response"""
    response = {
        "success": True,
        "data": data,
        "timestamp": datetime.now().isoformat()
    }
    if message:
        response["message"] = message
    
    return jsonify(response), status_code
```

**Örnek Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com"
  },
  "message": "User created",
  "timestamp": "2025-12-10T15:30:45.123456"
}
```

### 2. Error Response Format

```python
def format_error_response(error: str, error_code: str = None, status_code: int = 400):
    """Format error response"""
    response = {
        "success": False,
        "error": error,
        "timestamp": datetime.now().isoformat()
    }
    if error_code:
        response["error_code"] = error_code
    
    return jsonify(response), status_code
```

**Örnek Response**:
```json
{
  "success": false,
  "error": "User not found",
  "error_code": "USER_NOT_FOUND",
  "timestamp": "2025-12-10T15:30:45.123456"
}
```

### 3. List Response Format

```python
def format_list_response(items: list, total: int, limit: int = 50, offset: int = 0):
    """Format list response with pagination"""
    return jsonify({
        "success": True,
        "data": items,
        "pagination": {
            "total": total,
            "returned": len(items),
            "limit": limit,
            "offset": offset,
            "pages": (total + limit - 1) // limit
        },
        "timestamp": datetime.now().isoformat()
    }), 200
```

**Örnek Response**:
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "total": 100,
    "returned": 50,
    "limit": 50,
    "offset": 0,
    "pages": 2
  },
  "timestamp": "2025-12-10T15:30:45.123456"
}
```

---

## Error Handling

### 1. Standard Error Codes

```python
ERROR_CODES = {
    "INVALID_CREDENTIALS": {
        "status": 401,
        "message": "Invalid username or password"
    },
    "USER_NOT_FOUND": {
        "status": 404,
        "message": "User not found"
    },
    "USER_EXISTS": {
        "status": 409,
        "message": "User already exists"
    },
    "INVALID_REQUEST": {
        "status": 400,
        "message": "Invalid request data"
    },
    "UNAUTHORIZED": {
        "status": 401,
        "message": "Unauthorized access"
    },
    "FORBIDDEN": {
        "status": 403,
        "message": "Forbidden resource"
    },
    "SERVER_ERROR": {
        "status": 500,
        "message": "Internal server error"
    }
}
```

### 2. Error Response Examples

#### Missing Credentials
```bash
curl -X POST "http://localhost:5000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "alice"}'
```

**Yanıt (400)**:
```json
{
  "success": false,
  "error": "Missing credentials",
  "error_code": "INVALID_REQUEST",
  "timestamp": "2025-12-10T15:30:45.123456"
}
```

#### User Not Found
```bash
curl "http://localhost:5000/api/v1/users/nonexistent?key=12345"
```

**Yanıt (404)**:
```json
{
  "success": false,
  "error": "User not found",
  "error_code": "USER_NOT_FOUND",
  "timestamp": "2025-12-10T15:30:45.123456"
}
```

#### Invalid API Key
```bash
curl "http://localhost:5000/api/v1/users?key=wrongkey"
```

**Yanıt (401)**:
```json
{
  "success": false,
  "error": "Unauthorized",
  "error_code": "UNAUTHORIZED",
  "timestamp": "2025-12-10T15:30:45.123456"
}
```

---

## Security Measures

### 1. Password Protection in JSON

**Asla döndürülmez**:
```python
# ❌ WRONG - Never do this!
return jsonify({
    "user": user,
    "password": user.password  # SECURITY RISK!
}), 200
```

**Doğru yaklaşım**:
```python
# ✅ CORRECT
user_data = {
    "id": user.id,
    "username": user.username,
    "full_name": user.full_name,
    "email": user.email
    # Password hash is NOT included
}
return jsonify({"success": True, "user": user_data}), 200
```

### 2. Sensitive Data Masking

```python
def mask_email(email: str) -> str:
    """Mask email for privacy"""
    parts = email.split("@")
    if len(parts[0]) > 2:
        masked = parts[0][0] + "*" * (len(parts[0]) - 2) + parts[0][-1]
    else:
        masked = "*" + parts[0][1:]
    return masked + "@" + parts[1]

# Kullanım
masked_email = mask_email("alice@example.com")
# Sonuç: "a***e@example.com"
```

### 3. Encryption for Sensitive Attributes

```python
from cryptography.fernet import Fernet

def encrypt_sensitive_data(data: str, key: bytes) -> str:
    """Encrypt sensitive data"""
    f = Fernet(key)
    return f.encrypt(data.encode()).decode()

def decrypt_sensitive_data(encrypted: str, key: bytes) -> str:
    """Decrypt sensitive data"""
    f = Fernet(key)
    return f.decrypt(encrypted.encode()).decode()
```

### 4. Access Control in JSON

```python
def filter_user_response(user: dict, requesting_user: dict) -> dict:
    """Filter user data based on permissions"""
    filtered = {
        "id": user["id"],
        "username": user["username"],
        "full_name": user["full_name"]
    }
    
    # Only include email if user is self or admin
    if requesting_user["id"] == user["id"] or requesting_user["is_admin"]:
        filtered["email"] = user["email"]
    
    # Only show sensitive attributes to self and admins
    if requesting_user["id"] == user["id"] or requesting_user["is_admin"]:
        filtered["attributes"] = user.get("attributes", {})
    
    return filtered
```

---

## Best Practices

### 1. JSON Schema Validation

```python
from jsonschema import validate, ValidationError

USER_SCHEMA = {
    "type": "object",
    "properties": {
        "username": {"type": "string", "minLength": 1, "maxLength": 50},
        "password": {"type": "string", "minLength": 8},
        "email": {"type": "string", "format": "email"},
        "full_name": {"type": "string", "maxLength": 100}
    },
    "required": ["username", "password"]
}

def validate_user_json(data: dict) -> tuple[bool, str]:
    """Validate user JSON against schema"""
    try:
        validate(instance=data, schema=USER_SCHEMA)
        return True, "Valid"
    except ValidationError as e:
        return False, str(e)
```

### 2. Consistent Field Names

```python
# ✅ Use snake_case consistently
{
    "user_id": 1,
    "first_name": "Alice",
    "last_name": "Smith",
    "email_address": "alice@example.com",
    "is_active": true,
    "created_at": "2025-11-20T10:00:00Z"
}

# ❌ Avoid mixing conventions
{
    "userId": 1,           # camelCase
    "firstName": "Alice",  # camelCase
    "last_name": "Smith",  # snake_case (inconsistent!)
    "Email": "alice@example.com",  # PascalCase
}
```

### 3. Proper Timestamp Format

```python
# ✅ ISO 8601 format with timezone
"2025-12-10T15:30:45.123456Z"
"2025-12-10T15:30:45+00:00"

# ❌ Avoid
"12/10/2025 15:30:45"  # Local format
"2025-12-10 15:30:45"  # Missing timezone
1734093045             # Unix timestamp (harder to read)
```

### 4. Null vs Missing Fields

```python
# ✅ Include null for explicitly missing optional fields
{
    "username": "alice",
    "phone_number": null,      # Explicitly null
    "email": "alice@example.com"
}

# ❌ Avoid inconsistency
{
    "username": "alice",
    "email": "alice@example.com"
    // phone_number field is missing - ambiguous!
}
```

### 5. Error Messages

```python
# ✅ Clear, specific error messages
{
    "success": false,
    "error": "Invalid email format",
    "error_code": "INVALID_EMAIL",
    "field": "email"
}

# ❌ Avoid vague messages
{
    "error": "Bad request"  // Not helpful!
}
```

---

## 🔒 Koruma Kontrol Listesi

- ✅ Passwords hiçbir zaman JSON'da döndürülmez
- ✅ Sensitive attributes encrypted
- ✅ Input validation on all fields
- ✅ Size limits enforced
- ✅ XSS protection (HTML escaping)
- ✅ SQL injection prevention (prepared statements)
- ✅ CSRF tokens for state-changing operations
- ✅ API key validation
- ✅ Rate limiting enabled
- ✅ HTTPS enforced in production
- ✅ Access control checks
- ✅ Audit logging enabled
- ✅ Request/response compression
- ✅ Error messages sanitized

---

## 📊 JSON Size Optimization

```python
def compress_response(response: dict) -> str:
    """Compress JSON response"""
    import json
    return json.dumps(response, separators=(',', ':'))

# Örnek
# Uncompressed: 543 bytes
# Compressed:   456 bytes (16% smaller)
```

---

## 🧪 Test Örnekleri

### Valid User Creation
```bash
curl -X POST "http://localhost:5000/api/v1/users?key=12345" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "SecurePass123!",
    "full_name": "Alice Smith"
  }' | jq
```

### Invalid Request (Missing Password)
```bash
curl -X POST "http://localhost:5000/api/v1/users?key=12345" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "full_name": "Alice Smith"
  }' | jq
```

**Yanıt**:
```json
{
  "success": false,
  "error": "Missing username or password"
}
```

---

**Son Güncelleme**: 10 Aralık 2025  
**Sürüm**: 2.0.0  
**Durum**: ✅ Production Ready
