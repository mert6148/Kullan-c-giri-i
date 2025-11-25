# Varlıklar (Assets) - Hızlı Komut Referansı

## CLI Komutları

### Kullanıcı Yönetimi

```bash
# Kullanıcı oluştur
python print.py add-user <username> -p <password> -f "<full name>"

# Kullanıcı listele
python print.py list-users

# Kullanıcı sil
python print.py del-user <username>
```

### Giriş/Çıkış

```bash
# Giriş yap
python print.py login <username> -p <password>

# Çıkış yap
python print.py logout -u <username>

# Oturumları göster
python print.py show-sessions
```

### Varlıklar (YENÍÑ)

```bash
# Varlık ayarla
python print.py set-asset <username> <asset_name> <value> \
  [--type string|integer|boolean|json] \
  [--category profile|preferences|security|system|custom]

# Varlık al
python print.py get-asset <username> <asset_name>

# Tüm varlıkları göster
python print.py show-assets <username> [--category <category>]

# Varlık sil
python print.py delete-asset <username> <asset_name>
```

---

## Varlık Örnekleri

### Profil Varlıkları

```bash
python print.py set-asset john first_name "John" -c profile
python print.py set-asset john last_name "Doe" -c profile
python print.py set-asset john email "john@example.com" -c profile
python print.py set-asset john phone "+1-555-0123" -c profile
```

### Tercih Varlıkları

```bash
python print.py set-asset john theme "dark" -c preferences
python print.py set-asset john language "tr" -c preferences
python print.py set-asset john timezone "Europe/Istanbul" -c preferences
python print.py set-asset john font_size "14" -t integer -c preferences
```

### Güvenlik Varlıkları

```bash
python print.py set-asset john two_factor_enabled "true" -t boolean -c security
python print.py set-asset john last_password_change "2025-11-20" -c security
```

### Sistem Varlıkları

```bash
python print.py set-asset john login_count "42" -t integer -c system
python print.py set-asset john last_activity "2025-11-25 20:00:00" -c system
python print.py set-asset john total_sessions "156" -t integer -c system
```

---

## Sık Kullanılan Senaryolar

### Senaryo 1: Yeni Kullanıcı Kurulumu

```bash
#!/bin/bash

USERNAME=$1
PASSWORD=$2
FIRST_NAME=$3
LAST_NAME=$4
EMAIL=$5

# Kullanıcı oluştur
python print.py add-user $USERNAME -p $PASSWORD -f "$FIRST_NAME $LAST_NAME"

# Profili doldur
python print.py set-asset $USERNAME first_name "$FIRST_NAME" -c profile
python print.py set-asset $USERNAME last_name "$LAST_NAME" -c profile
python print.py set-asset $USERNAME email "$EMAIL" -c profile

# Varsayılan tercihleri ayarla
python print.py set-asset $USERNAME theme "light" -c preferences
python print.py set-asset $USERNAME language "tr" -c preferences
python print.py set-asset $USERNAME timezone "Europe/Istanbul" -c preferences

# Güvenlik ayarlarını başlat
python print.py set-asset $USERNAME two_factor_enabled "false" -t boolean -c security

# Sistem kayıtlarını başlat
python print.py set-asset $USERNAME login_count "0" -t integer -c system
python print.py set-asset $USERNAME last_activity "" -c system

echo "✓ Kullanıcı kurulumu tamamlandı: $USERNAME"
```

### Senaryo 2: Giriş Sonrası Güncelleme

```bash
#!/bin/bash

USERNAME=$1

# Giriş yap
python print.py login $USERNAME -p <password>

if [ $? -eq 0 ]; then
    # Login sayısını artır
    COUNT=$(python print.py get-asset $USERNAME login_count | grep -oP '"asset_value":\s*"\K[^"]*')
    COUNT=$((COUNT + 1))
    python print.py set-asset $USERNAME login_count "$COUNT" -t integer -c system
    
    # Son aktivite zamanını güncelle
    python print.py set-asset $USERNAME last_activity "$(date -u +'%Y-%m-%d %H:%M:%S')" -c system
    
    echo "✓ Giriş sayısı güncellendi: $COUNT"
fi
```

### Senaryo 3: Güvenlik Kontrolü

```bash
#!/bin/bash

USERNAME=$1

echo "=== Güvenlik Kontrolü ==="
python print.py show-assets $USERNAME -c security
```

### Senaryo 4: Varlıkları Dışa Aktar

```bash
#!/bin/bash

USERNAME=$1
OUTPUT_FILE="${USERNAME}_assets.json"

python print.py show-assets $USERNAME > $OUTPUT_FILE
echo "✓ Varlıklar dışa aktarıldı: $OUTPUT_FILE"
```

---

## REST API Komutları

### Varlık Ekle

```bash
curl -X POST http://localhost:5000/api/v1/users/john/attributes \
  -H "Content-Type: application/json" \
  -d '{
    "asset_name": "theme",
    "asset_value": "dark",
    "asset_type": "string",
    "category": "preferences"
  }'
```

### Varlık Al

```bash
curl http://localhost:5000/api/v1/users/john/attributes/theme
```

### Tüm Varlıkları Al

```bash
curl http://localhost:5000/api/v1/users/john/attributes
```

### Varlık Sil

```bash
curl -X DELETE http://localhost:5000/api/v1/users/john/attributes/theme
```

---

## Python API Örnekleri

### UserAssetManager Kullanımı

```python
from assets.assest import UserAssetManager, ASSET_CATEGORY_PREFERENCES

# Manager başlat
manager = UserAssetManager("login_system.db")

# Varlık ayarla
manager.set_asset(user_id, "theme", "dark", category="preferences")

# Varlık al
asset = manager.get_asset(user_id, "theme")
print(asset.asset_value)  # "dark"

# Kategoriye göre al
prefs = manager.get_assets_by_category(user_id, "preferences")
for name, asset in prefs.items():
    print(f"{name}: {asset.asset_value}")

# Varlık sil
manager.delete_asset(user_id, "theme")
```

---

## Veri Türleri Referansı

```bash
# String (varsayılan)
python print.py set-asset john department "Engineering"

# Integer
python print.py set-asset john age "30" -t integer

# Boolean
python print.py set-asset john active "true" -t boolean

# JSON
python print.py set-asset john preferences '{"color":"red","size":"large"}' -t json

# File (dosya yolu)
python print.py set-asset john avatar_path "/path/to/avatar.png" -t file
```

---

## Kategori Referansı

```bash
# Profile (Profil) - Kişisel bilgiler
python print.py set-asset john first_name "John" -c profile

# Preferences (Tercihler) - Kullanıcı seçimleri
python print.py set-asset john theme "dark" -c preferences

# Security (Güvenlik) - Kimlik doğrulama bilgileri
python print.py set-asset john two_factor_enabled "true" -t boolean -c security

# System (Sistem) - Sistem meta-veri
python print.py set-asset john login_count "42" -t integer -c system

# Custom (Özel) - Uygulamaya özel (varsayılan)
python print.py set-asset john department "Engineering"
```

---

## Hata Giderme

| Hata | Sebep | Çözüm |
|------|-------|-------|
| "Kullanıcı bulunamadı" | Kullanıcı adı yanlış | `python print.py list-users` ile kontrol et |
| "Varlık bulunamadı" | Varlık silinmiş | `python print.py show-assets <username>` |
| "Geçersiz tür" | Yanlış --type | Desteklenenleri kontrol et |
| "Assets modülü yok" | assets/assest.py eksik | Dosya mevcut mu kontrol et |
| JSON parse hatası | Yanlış JSON formatı | JSON'u valide et |

---

## Faydalı Komut Kombinasyonları

### Kullanıcı Profili Yedekle

```bash
python print.py show-assets <username> > profile_backup.json
```

### Tüm Varlıkları Sıfırla

```bash
# Her kategoriyi kontrol et ve sil
python print.py show-assets <username>
# Sonra gerekli olanları delete-asset ile silin
```

### Kategoriye Göre Rapor Al

```bash
python print.py show-assets <username> -c profile
python print.py show-assets <username> -c preferences
python print.py show-assets <username> -c security
python print.py show-assets <username> -c system
```

---

## Kaynak Dokumentasyon

- **Kapsamlı Rehber**: ASSETS_GUIDE.md
- **Uygulama Raporu**: ASSETS_IMPLEMENTATION_REPORT.md
- **API Dokümantasyonu**: API_UI_SETUP.md
- **Veritabanı Şeması**: DATABASE.md
- **Test Paketi**: test_assets.py

---

**Versiyon**: 1.0  
**Tarih**: 25 Kasım 2025  
**Dil**: Türkçe  
**Proje**: Kullanıcı Girişi Sistemi
