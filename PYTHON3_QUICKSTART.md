# Python3 Şema Validatörü - Quickstart Guide

## 📋 İçindekiler
1. [Kurulum](#kurulum)
2. [Temel Kullanım](#temel-kullanım)
3. [If Döngüsü Örnekleri](#if-döngüsü-örnekleri)
4. [CLI Komutları](#cli-komutları)

---

## Kurulum

### Gereksinimler
- Python 3.8+
- SQLite3

### Basit Kurulum
```bash
# Dizine gir
cd assests

# Paketi yükle
python -m pip install -r requirements.txt
```

---

## Temel Kullanım

### 1️⃣ Manager Oluştur
```python
from assest import UserAssetManager

manager = UserAssetManager()
```

### 2️⃣ Varlık Doğrula
```python
from assest import validate_schema_with_conditions, ASSET_CATEGORY_PROFILE

is_valid, error = validate_schema_with_conditions(
    manager,
    ASSET_CATEGORY_PROFILE,
    "email",
    "user@example.com"
)

print("✓ Geçerli" if is_valid else f"✗ {error}")
```

### 3️⃣ Toplu Doğrula ve Kaydet
```python
from assest import (
    validate_user_assets_batch,
    process_validated_assets,
    ASSET_CATEGORY_PROFILE,
    ASSET_CATEGORY_PREFERENCES
)

# Doğrulanacak varlıklar
assets = {
    ASSET_CATEGORY_PROFILE: {
        "email": "user@example.com",
        "first_name": "Ahmet"
    },
    ASSET_CATEGORY_PREFERENCES: {
        "theme": "dark"
    }
}

# Doğrula
valid, errors = validate_user_assets_batch(manager, assets)

if valid:
    # Kaydet
    success, failed = process_validated_assets(manager, user_id=1, assets_dict=assets)
    print("✓ Kaydedildi" if success else f"✗ Hata: {failed}")
else:
    print(f"✗ Hata: {errors}")
```

---

## If Döngüsü Örnekleri

### Desen 1: Kategoriler Üzerinde Döngü
```python
from assest import VALID_ASSET_CATEGORIES

for category in VALID_ASSET_CATEGORIES:
    assets = manager.get_assets_by_category(user_id=1, category=category)
    
    if assets:  # Kategoride varlık var mı?
        print(f"\n{category}:")
        for name, asset in assets.items():
            print(f"  {name}: {asset.asset_value}")
```

### Desen 2: Koşullu Doğrulama
```python
# Geçerli ve geçersiz varlıkları ayır
valid_assets = []
invalid_assets = []

for category, assets in data.items():
    if not isinstance(assets, dict):
        continue  # Geçersiz tipi atla
    
    for name, value in assets.items():
        is_valid, error = validate_schema_with_conditions(
            manager, category, name, value
        )
        
        if is_valid:
            valid_assets.append((category, name, value))
        else:
            invalid_assets.append((category, name, error))

# Sonuçları işle
print(f"✓ {len(valid_assets)} geçerli")
print(f"✗ {len(invalid_assets)} geçersiz")
```

### Desen 3: Hata Kontrolü
```python
errors = []

# Tüm varlıkları kontrol et
for category in VALID_ASSET_CATEGORIES:
    if category not in data:
        continue
    
    for asset_name, asset_value in data[category].items():
        is_valid, error = validate_schema_with_conditions(
            manager, category, asset_name, asset_value
        )
        
        if not is_valid:
            errors.append(f"{category}.{asset_name}: {error}")

# Hataları raporla
if errors:
    print(f"✗ {len(errors)} hata:")
    for err in errors:
        print(f"  {err}")
else:
    print("✓ Tüm varlıklar geçerli")
```

### Desen 4: Veri Temizleme ve Geri Yükleme
```python
# Eski varlıkları sil
manager.delete_all_assets(user_id)

# Yenilerini kaydet
for category, assets in new_data.items():
    for name, value in assets.items():
        success, _ = manager.set_asset(
            user_id, name, value, category=category, validate=True
        )
        
        if not success:
            print(f"✗ {name} kaydedilemedi")
```

### Desen 5: Erken Çıkış
```python
# İlk hatayı bulunca dur
found_error = False

for category in VALID_ASSET_CATEGORIES:
    if found_error:
        break
    
    for asset_name, asset_value in data.get(category, {}).items():
        is_valid, error = validate_schema_with_conditions(
            manager, category, asset_name, asset_value
        )
        
        if not is_valid:
            print(f"✗ Hata: {error}")
            found_error = True
            break  # İç döngüyü kır
```

---

## CLI Komutları

### Tek Varlığı Doğrula
```bash
python3 -c "
from assets.assest import *
manager = UserAssetManager()
is_valid, err = validate_schema_with_conditions(
    manager, 'profile', 'email', 'test@example.com'
)
print('✓ Valid' if is_valid else f'✗ Error: {err}')
"
```

### Toplu Doğrulama
```bash
python3 -c "
from assests.assest import *

manager = UserAssetManager()
assets = {
    'profile': {'email': 'test@example.com'},
    'preferences': {'theme': 'dark'}
}

valid, errors = validate_user_assets_batch(manager, assets)
if valid:
    print('✓ Tüm varlıklar geçerli')
else:
    for err in errors:
        print(f'✗ {err}')
"
```

### Test Script Çalıştır
```bash
python test_validation_workflow.py
```

---

## Varlık Kategorileri

| Kategori | Açıklama | Örnek Varlıklar |
|----------|----------|-----------------|
| `profile` | Profil Bilgisi | email, first_name, phone |
| `preferences` | Tercihler | theme, language, timezone |
| `security` | Güvenlik | two_factor_enabled, login_attempts |
| `system` | Sistem | login_count, ip_address |
| `custom` | Özel | Kullanıcı tanımlı |

---

## Varlık Tipleri

| Tip | Açıklama | Örnek |
|-----|----------|--------|
| `string` | Metin | "Ahmet" |
| `integer` | Sayı | 5 |
| `boolean` | True/False | true, false |
| `json` | JSON | {"key": "value"} |
| `binary` | İkili | b"bytes" |
| `file` | Dosya | "path/to/file" |

---

## Hata Yönetimi

### Yaygın Hatalar

| Hata | Çözüm |
|------|-------|
| "Validatör aktif değil" | Manager'ı kontrol et |
| "Geçersiz kategori" | Kategoriyi VALID_ASSET_CATEGORIES'den seç |
| "Değer pattern'e uymuyor" | Formatı düzelt (email vs.) |
| "String uzunluğu maksimum X" | Stringi kısalt |
| "Değer maksimum X olmalı" | Değeri azalt |

### Try-Except Kullanımı
```python
try:
    success, failed = process_validated_assets(
        manager, user_id, assets_dict
    )
    
    if not success:
        for fail in failed:
            print(f"✗ {fail}")
            
except Exception as e:
    print(f"✗ Beklenmeyen hata: {e}")
```

---

## Tam Örnek: Kullanıcı Kaydı

```python
#!/usr/bin/env python3

from assest import (
    UserAssetManager,
    validate_user_assets_batch,
    process_validated_assets,
    ASSET_CATEGORY_PROFILE,
    ASSET_CATEGORY_SECURITY,
)


def register_user(user_id, user_data):
    """Yeni kullanıcı kaydı"""
    
    manager = UserAssetManager()
    
    # Varlıkları hazırla
    assets = {
        ASSET_CATEGORY_PROFILE: {
            "first_name": user_data["first_name"],
            "last_name": user_data["last_name"],
            "email": user_data["email"],
        },
        ASSET_CATEGORY_SECURITY: {
            "two_factor_enabled": "false",
            "login_attempts": "0",
        }
    }
    
    # Adım 1: Doğrula
    print("1️⃣  Doğrulanıyor...")
    valid, errors = validate_user_assets_batch(manager, assets)
    
    if not valid:
        print("❌ Doğrulama başarısız:")
        for err in errors:
            print(f"   {err}")
        return False
    
    print("✅ Doğrulama başarılı")
    
    # Adım 2: Kaydet
    print("\n2️⃣  Kaydediliyor...")
    success, failed = process_validated_assets(
        manager, user_id, assets
    )
    
    if not success:
        print("❌ Kayıt başarısız:")
        for fail in failed:
            print(f"   {fail}")
        return False
    
    print("✅ Varlıklar kaydedildi")
    
    # Adım 3: Doğrula
    print("\n3️⃣  Doğrulanıyor...")
    saved = manager.get_all_assets(user_id)
    
    if saved:
        print("✅ Varlıklar geri alındı")
        return True
    else:
        print("❌ Varlıklar geri alınamadı")
        return False


if __name__ == "__main__":
    # Test
    user_data = {
        "first_name": "Ahmet",
        "last_name": "Yılmaz",
        "email": "ahmet@example.com",
    }
    
    if register_user(1, user_data):
        print("\n✅ Kullanıcı başarıyla kaydedildi!")
    else:
        print("\n❌ Kullanıcı kaydedilemedi!")
```

---

## 🚀 Hızlı Başlangıç

### 3 Satırda Başla:
```python
from assests.assest import UserAssetManager, validate_schema_with_conditions
manager = UserAssetManager()
is_valid, err = validate_schema_with_conditions(manager, "profile", "email", "test@example.com")
```

### Sonra If Döngüsü Ekle:
```python
for category in ["profile", "preferences"]:
    for name, value in data[category].items():
        is_valid, err = validate_schema_with_conditions(manager, category, name, value)
        print(f"{'✓' if is_valid else '✗'} {name}: {value}")
```

### Kaydet:
```python
process_validated_assets(manager, user_id=1, assets_dict=data)
```

---

## 📚 Daha Fazla

- `VALIDATION_WORKFLOW.md` - Detaylı dokümantasyon
- `test_validation_workflow.py` - Tam test örnekleri
- `assest.py` - Kaynak kodu

---

**Son Güncelleme**: 10 Aralık 2025  
**Python**: 3.8+  
**Durum**: ✅ Prodüksiyon Hazırı
