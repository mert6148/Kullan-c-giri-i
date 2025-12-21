# Network API - REST Dokümantasyon

**Versiyon**: 2.0  
**Tarih**: 10 Aralık 2025  
**Durum**: ✅ Production Ready

---

## 📋 Genel Bakış

Network API, sistem ağ profillerini yönetmek için RESTful endpoints sunmaktadır. API key ile korunmuş ve JSON tabanlı yanıtlar verir.

---

## 🔐 Kimlik Doğrulama

Tüm istekler `key` parametresi ile API key göndermelidir:

```
GET /api/network/list?key=YOUR_API_KEY
POST /api/network/switch?key=YOUR_API_KEY
```

**Default API Key**: `12345` (Prodüksiyonda değiştirin!)

**Hash Algoritması**: `hash_equals()` ile timing attack koruması

---

## 📡 Endpoints

### 1. Profilleri Listele
```
GET /api/network/list?key=YOUR_API_KEY
```

**Açıklama**: Tüm ağ profillerini listeler

**Parametreler**:
- `key` (required): API key

**Response (200)**:
```json
{
  "status": "success",
  "count": 3,
  "profiles": {
    "local": {
      "name": "Local Network",
      "ip": "192.168.1.0/24",
      "dns": "8.8.8.8"
    },
    "remote": {
      "name": "Remote Network",
      "ip": "10.0.0.0/8",
      "dns": "1.1.1.1"
    },
    "vpn": {
      "name": "VPN Network",
      "ip": "172.16.0.0/12",
      "dns": "9.9.9.9"
    }
  }
}
```

**Error (401)**:
```json
{
  "error": "Unauthorized"
}
```

---

### 2. Aktif Profili Görüntüle
```
GET /api/network/active?key=YOUR_API_KEY
```

**Açıklama**: Şu anda aktif olan ağ profilini döndürür

**Response (200)**:
```json
{
  "status": "success",
  "active": "local",
  "data": {
    "name": "Local Network",
    "ip": "192.168.1.0/24",
    "dns": "8.8.8.8"
  }
}
```

**Hiç profil seçilmemişse**:
```json
{
  "status": "success",
  "active": null
}
```

---

### 3. Profili Değiştir
```
POST /api/network/switch?key=YOUR_API_KEY
Content-Type: application/json

{
  "profile": "remote"
}
```

**Açıklama**: Aktif ağ profilini değiştirir

**Request Body**:
```json
{
  "profile": "profile_name"
}
```

**Response (200)**:
```json
{
  "status": "success",
  "message": "Active network switched to 'remote'",
  "active": "remote",
  "data": {
    "name": "Remote Network",
    "ip": "10.0.0.0/8",
    "dns": "1.1.1.1"
  }
}
```

**Errors**:

Profil parametresi eksikse (400):
```json
{
  "error": "Missing profile parameter"
}
```

Profil bulunamazsa (404):
```json
{
  "error": "Profile 'invalid' not found"
}
```

Dosya yazma hatası (500):
```json
{
  "error": "Failed to write active profile"
}
```

---

### 4. Profili Doğrula
```
GET /api/network/validate?key=YOUR_API_KEY&profile=local
```

**Açıklama**: Bir profil adının geçerli olup olmadığını kontrol eder

**Parametreler**:
- `key` (required): API key
- `profile` (required): Kontrol edilecek profil adı

**Response (200)**:
```json
{
  "status": "success",
  "profile": "local",
  "valid": true,
  "message": "Profile exists"
}
```

**Profil bulunamazsa**:
```json
{
  "status": "success",
  "profile": "invalid",
  "valid": false,
  "message": "Profile not found"
}
```

---

## 🛠️ Utility Functions

### buildNetworkQuery()

Ağ parametrelerini URL query string'e dönüştürür.

```php
function buildNetworkQuery(
    array $data,
    string $numeric_prefix = "",
    string $arg_separator = "&",
    int $encoding_type = PHP_QUERY_RFC1738
): string
```

**Örnek**:
```php
$params = [
    "ip" => "192.168.1.0",
    "dns" => ["primary" => "8.8.8.8", "secondary" => "8.8.4.4"]
];

$query = buildNetworkQuery($params);
// Sonuç: ip=192.168.1.0&dns%5Bprimary%5D=8.8.8.8&dns%5Bsecondary%5D=8.8.4.4
```

### encodeNetworkProfile()

Ağ profilini URL-safe parametrelere dönüştürür.

```php
function encodeNetworkProfile(array $profile): string
```

**Örnek**:
```php
$profile = [
    "name" => "Local Network",
    "ip" => "192.168.1.0/24"
];

$encoded = encodeNetworkProfile($profile);
// Sonuç: name=Local+Network&ip=192.168.1.0%2F24
```

### parseNetworkQuery()

Query string'i array'e dönüştürür.

```php
function parseNetworkQuery(string $query_string): array
```

**Örnek**:
```php
$query = "ip=192.168.1.0&dns[primary]=8.8.8.8";
$parsed = parseNetworkQuery($query);
// Sonuç:
// [
//   "ip" => "192.168.1.0",
//   "dns" => ["primary" => "8.8.8.8"]
// ]
```

### validateNetworkAccess()

API key doğrulaması yapar (timing attack korumalı).

```php
function validateNetworkAccess(string $key, string $valid_key): bool
```

### setNetworkHeaders()

HTTP başlıklarını ayarlar.

```php
function setNetworkHeaders(
    string $content_type = "application/json",
    ?int $content_length = null
): void
```

### sendJsonResponse()

JSON yanıtını güvenli bir şekilde gönderir.

```php
function sendJsonResponse(
    array $data,
    int $http_code = 200,
    int $json_options = JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE
): void
```

---

## 🔒 Güvenlik Özellikleri

### API Key Yönetimi
```php
// Hash equality check (timing attack koruması)
if (!validateNetworkAccess($_GET["key"], $API_KEY)) {
    sendJsonResponse(["error" => "Unauthorized"], 401);
}
```

### Encoding İşlemleri
- **RFC1738**: `urlencode()` - Uyumluluk
- **RFC3986**: `rawurlencode()` - Daha kesin

### Error Handling
- Geçersiz JSON: UTF-8 dönüşümü
- Dosya yazma hataları: HTTP 500
- Bulunamayan endpoint: HTTP 404

### Response Headers
```
Content-Type: application/json; charset=utf-8
X-API-Version: 2.0
Cache-Control: no-cache, no-store, must-revalidate
```

---

## 📝 cURL Örnekleri

### Listeyi Al
```bash
curl -X GET "http://localhost/api/network/list?key=12345"
```

### Aktif Profili Görüntüle
```bash
curl -X GET "http://localhost/api/network/active?key=12345"
```

### Profili Değiştir
```bash
curl -X POST "http://localhost/api/network/switch?key=12345" \
  -H "Content-Type: application/json" \
  -d '{"profile": "remote"}'
```

### Profili Doğrula
```bash
curl -X GET "http://localhost/api/network/validate?key=12345&profile=local"
```

---

## 🧪 PHP Örnekleri

### GET Request
```php
<?php
$api_key = "12345";
$url = "http://localhost/api/network/list?key=" . urlencode($api_key);

$response = file_get_contents($url);
$data = json_decode($response, true);

echo "Profil Sayısı: " . $data['count'];
?>
```

### POST Request
```php
<?php
$api_key = "12345";
$url = "http://localhost/api/network/switch?key=" . urlencode($api_key);

$options = [
    "http" => [
        "method" => "POST",
        "header" => "Content-Type: application/json\r\n",
        "content" => json_encode(["profile" => "remote"])
    ]
];

$context = stream_context_create($options);
$response = file_get_contents($url, false, $context);
$data = json_decode($response, true);

echo "Aktif Profil: " . $data['active'];
?>
```

---

## 📊 HTTP Status Kodları

| Kod | Anlamı | Örnek |
|-----|--------|-------|
| 200 | OK | Başarılı istek |
| 400 | Bad Request | Eksik parametre |
| 401 | Unauthorized | Hatalı API key |
| 404 | Not Found | Profil/Endpoint bulunamadı |
| 500 | Server Error | Dosya yazma hatası |

---

## 🔄 İş Akışı

```
Client Request
    ↓
API Key Doğrulama
    ↓
URL Parsing
    ↓
Endpoint Eşleştir
    ↓
Parametreleri Al
    ↓
Veri İşle
    ↓
JSON Response
    ↓
HTTP Headers
    ↓
Client Response
```

---

## 🚀 Gelişmiş Kullanım

### Profil Yönetimi

```php
<?php
// Tüm profilleri getir
$profiles = json_decode(
    file_get_contents('http://localhost/api/network/list?key=12345'),
    true
)['profiles'];

// Her profilin özelliklerini kontrol et
foreach ($profiles as $name => $config) {
    echo "Profil: $name\n";
    echo "  IP: " . $config['ip'] . "\n";
    echo "  DNS: " . $config['dns'] . "\n";
}
?>
```

### Batch Operations

```php
<?php
$profiles = ['local', 'remote', 'vpn'];

foreach ($profiles as $profile) {
    $url = "http://localhost/api/network/validate?key=12345&profile=$profile";
    $result = json_decode(file_get_contents($url), true);

    if ($result['valid']) {
        echo "✓ $profile\n";
    } else {
        echo "✗ $profile\n";
    }
}
?>
```

---

## 🐛 Troubleshooting

### API Yanıt Vermiyor
1. API key'in doğru olduğunu kontrol edin
2. `network.php` dosya yolunu kontrol edin
3. Server logs'u kontrol edin

### JSON Parsing Hatası
1. Response'ın geçerli JSON olduğunu doğrulayın
2. Character encoding'i kontrol edin (UTF-8)
3. PHP's `json_last_error()` kullanın

### Profil Bulunamıyor
1. Profil adının doğru yazıldığını kontrol edin
2. `/api/network/list` ile mevcut profilleri listeleyin
3. `network.php` config dosyasını kontrol edin

---

## 📚 Referanslar

- **RFC 1738**: https://tools.ietf.org/html/rfc1738
- **RFC 3986**: https://tools.ietf.org/html/rfc3986
- **HTTP Status Codes**: https://httpwg.org/specs/rfc7231.html
- **JSON Encoding**: https://www.json.org/

---

## ✅ Version History

### v2.0 (10 Aralık 2025)
- ✅ buildNetworkQuery() fonksiyonu refactor
- ✅ Yeni utility functions eklendi
- ✅ HTTP status code'ları standardize edildi
- ✅ Error handling iyileştirildi
- ✅ Timing attack protection eklendi
- ✅ Comprehensive documentation

### v1.0
- Temel API endpoints

---

**Son Güncelleme**: 10 Aralık 2025  
**Desteklenen PHP**: 7.4+  
**Lisans**: MIT
